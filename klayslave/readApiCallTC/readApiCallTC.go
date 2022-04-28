package readApiCallTC

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/clipool"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/myzhan/boomer"
	"github.com/tidwall/gjson"
)

var (
	gasPrice *big.Int
	endPoint string
	cliPool  clipool.ClientPool

	mutex       sync.Mutex
	initialized = false

	nAcc   int
	accGrp []*account.Account

	latestBlockNumber *big.Int
	count             uint64
)

func Init(accs []*account.Account, ep string, gp *big.Int) {
	mutex.Lock()
	defer mutex.Unlock()

	if initialized {
		return
	}
	initialized = true

	latestBlockNumber = big.NewInt(0)
	count = 0
	gasPrice = gp
	endPoint = ep

	cliCreate := func() interface{} {
		c, err := rpc.Dial(endPoint)
		if err != nil {
			log.Fatalf("Failed to connect RPC: %v", err)
		}
		return c
	}
	cliPool.Init(20, 300, cliCreate)

	for _, acc := range accs {
		accGrp = append(accGrp, acc)
	}
	nAcc = len(accGrp)
}

func sendBoomerEvent(tcName string, logString string, elapsed int64, cli *rpc.Client, err error) {
	if err == nil {
		boomer.Events.Publish("request_success", "http", tcName+" to "+endPoint, elapsed, int64(10))
		cliPool.Free(cli)
	} else {
		log.Printf("[TC] %s: %s, err=%v\n", tcName, logString, err)
		boomer.Events.Publish("request_failure", "http", tcName+" to "+endPoint, elapsed, err.Error())
		cli.Close()
	}
}

func getRandomBlockNumber(cli *client.Client, ctx context.Context) *big.Int {
	mutex.Lock()
	defer mutex.Unlock()

	count %= 10000000
	if count%10000 == 0 {
		bn, err := cli.BlockNumber(ctx)
		if err != nil {
			log.Printf("Failed to update the current block number. err=%s\n", err)
		} else {
			log.Printf("Update the current block number. blockNumber=0x%s\n", bn.Text(16))
			latestBlockNumber.Set(bn)
		}
	}
	count += 1

	return big.NewInt(0).Rand(rand.New(rand.NewSource(time.Now().UnixNano())), latestBlockNumber)
}

func GasPrice() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	start := boomer.Now()
	gp, err := cli.SuggestGasPrice(ctx)
	elapsed := boomer.Now() - start
	if err == nil && gp.Cmp(gasPrice) != 0 {
		err = errors.New("wrong gas price: " + gp.String() + ", answer: " + gasPrice.String())
	}
	sendBoomerEvent("readGasPrice", "Failed to call klay_gasPrice", elapsed, rpcCli, err)
}

func BlockNumber() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	start := boomer.Now()

	bn, err := cli.BlockNumber(ctx)
	if err == nil && bn.Cmp(big.NewInt(0)) != 1 {
		err = errors.New("wrong block number: 0x" + bn.Text(16) + ", answer: smaller than 0")
	}

	elapsed := boomer.Now() - start
	sendBoomerEvent("readBlockNumber", "Failed to call klay_blockNumber", elapsed, rpcCli, err)
}

func GetBlockByNumber() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)
	cli := client.NewClient(rpcCli)

	ansBN := getRandomBlockNumber(cli, ctx)
	start := boomer.Now()

	block, err := cli.BlockByNumber(ctx, ansBN) //read the random block
	if err == nil && block.Header().Number.Cmp(ansBN) != 0 {
		err = errors.New("wrong block: 0x" + block.Header().Number.Text(16) + ", answer: 0x" + ansBN.Text(16))
	}

	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetBlockByNumber", "Failed to call klay_getBlockByNumber", elapsed, rpcCli, err)
}

func GetAccount() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)

	var j json.RawMessage
	fromAccount := accGrp[rand.Int()%nAcc]
	start := boomer.Now()

	err := rpcCli.CallContext(ctx, &j, "klay_getAccount", fromAccount.GetAddress(), "latest")
	if err == nil {
		ret := gjson.Get(string(j), "accType").String()
		if ret != "1" {
			err = errors.New("wrong account type: " + ret + ", answer: 1")
		}
	}

	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetAccount", "Failed to call klay_getAccount", elapsed, rpcCli, err)
}

func GetBlockWithConsensusInfoByNumber() {
	ctx := context.Background()
	rpcCli := cliPool.Alloc().(*rpc.Client)

	ansBN := getRandomBlockNumber(client.NewClient(rpcCli), ctx)
	start := boomer.Now()

	var j json.RawMessage
	err := rpcCli.CallContext(ctx, &j, "klay_getBlockWithConsensusInfoByNumber", "0x"+ansBN.Text(16))
	if err == nil {
		ret := gjson.Get(string(j), "number").String()
		if !strings.Contains(ret, "0x"+ansBN.Text(16)) {
			err = errors.New("wrong block: " + ret + ", answer: " + "0x" + ansBN.Text(16))
		}
	}

	elapsed := boomer.Now() - start
	sendBoomerEvent("readGetBlockWithConsensusInfoByNumber",
		"Failed to call klay_GetBlockWithConsensusInfoByNumber", elapsed, rpcCli, err)
}
