package receiptCheckTc

import (
	"context"
	"log"
	"math/big"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/clipool"
	"github.com/klaytn/klaytn-load-tester/klayslave/task"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/myzhan/boomer"
)

const Name = "receiptCheckTx"

var (
	endPoint string

	nAcc   int
	accGrp []*account.Account

	cliPool clipool.ClientPool

	hashPoolSize      = 100 * 5 * 60 // for init 5min, if input send tps is 100Txs/Sec
	defaultInitSendTx = 1000 * 10    // for init 10sec, if input send TPS is 1000Txs/Sec

	tail     = 0
	isFull   = false
	hashPool []common.Hash
	rwMutex  sync.RWMutex

	ratioReadPerSend = 9 // read:send = ratioReadPerSend:1

	gasPrice *big.Int
)

func Init(params *task.Params) {
	gasPrice = params.GasPrice

	endPoint = params.Endpoint

	cliCreate := func() interface{} {
		c, err := client.Dial(endPoint)
		if err != nil {
			log.Fatalf("Failed to connect RPC: %v", err)
		}
		return c
	}

	cliPool.Init(20, 300, cliCreate)

	for _, acc := range params.AccGrp {
		accGrp = append(accGrp, acc)
	}

	nAcc = len(accGrp)

	hashPool = make([]common.Hash, hashPoolSize, hashPoolSize)
}

func addHash(hash common.Hash) {
	rwMutex.Lock()
	hashPool[tail] = hash

	tail = (tail + 1) % hashPoolSize
	if tail == 0 {
		isFull = true
	}

	rwMutex.Unlock()
}

func getHash() common.Hash {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	if isFull {
		return hashPool[rand.Int()%hashPoolSize]
	}
	return hashPool[rand.Int()%tail]
}

func RunSendTx() {
	cli := cliPool.Alloc().(*client.Client)

	from := accGrp[rand.Int()%nAcc]
	to := accGrp[rand.Int()%nAcc]
	value := big.NewInt(int64(rand.Int() % 3))

	start := boomer.Now()
	hash, _, err := from.TransferSignedTx(cli, to, value)
	addHash(hash)
	elapsed := boomer.Now() - start

	if err == nil {
		boomer.Events.Publish("request_success", Name, "send tx"+" to "+endPoint, elapsed, int64(10))
		cliPool.Free(cli)
	} else {
		boomer.Events.Publish("request_failure", Name, "send tx"+" to "+endPoint, elapsed, err.Error())
	}
}

func RunSendTxSingle() {
	cli := cliPool.Alloc().(*client.Client)

	from := accGrp[rand.Int()%nAcc]
	from.GetNonceFromBlock(cli)
	to := accGrp[rand.Int()%nAcc]
	value := big.NewInt(int64(rand.Int() % 3))

	start := boomer.Now()
	hash, _, err := from.TransferSignedTx(cli, to, value)
	addHash(hash)
	elapsed := boomer.Now() - start

	if err == nil {
		boomer.Events.Publish("request_success", Name, "send tx"+" to "+endPoint, elapsed, int64(10))
		cliPool.Free(cli)
	} else {
		boomer.Events.Publish("request_failure", Name, "send tx"+" to "+endPoint, elapsed, err.Error())
	}
}

var (
	receiptDisplayCnt uint64
)

func RunReadTx() {
	cli := cliPool.Alloc().(*client.Client)

	ctx := context.Background()
	hash := getHash()

	start := boomer.Now()

	receipt, err := cli.TransactionReceipt(ctx, hash)
	if err == nil {
		if rand.Int()%(1000*60) == 0 {
			log.Printf("pid(%v) : hash(%v) receipt checked\n", os.Getpid(), hash.String())
			log.Printf("%v", receipt)
		}
	} else {
		log.Printf("pid(%v) : hash(%v) receipt check err : %v\n", os.Getpid(), hash.String(), err)
	}

	elapsed := boomer.Now() - start

	if err == nil {
		boomer.Events.Publish("request_success", Name, "read tx"+" to "+endPoint, elapsed, int64(10))
		cliPool.Free(cli)
	} else {
		boomer.Events.Publish("request_failure", Name, "read tx"+" to "+endPoint, elapsed, err.Error())
	}
}

var (
	cnt      uint32
	initFlag = false
)

func Run() {
	nc := atomic.AddUint32(&cnt, 1)

	if !initFlag && nc < uint32(defaultInitSendTx) {
		RunSendTx()
	} else {
		initFlag = true

		// following logic can control the ratio between send/read task
		nc = nc % uint32(ratioReadPerSend+1)

		if nc == uint32(ratioReadPerSend) {
			RunSendTx()
		} else {
			RunReadTx()
		}
	}
}

func RunSingle() {
	nc := atomic.AddUint32(&cnt, 1)

	if !initFlag && nc < uint32(defaultInitSendTx) {
		RunSendTx()
	} else {
		initFlag = true

		// following logic can control the ratio between send/read task
		nc = nc % uint32(ratioReadPerSend+1)

		if nc == uint32(ratioReadPerSend) {
			RunSendTxSingle()
		} else {
			RunReadTx()
		}
	}
}
