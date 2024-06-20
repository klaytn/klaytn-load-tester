//go:generate abigen --sol ReadApiCallContract.sol --pkg readApiCallContract --out ReadApiCallContract.go
package readApiCallContractTC

import (
	"context"
	"errors"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/clipool"
	"github.com/klaytn/klaytn-load-tester/klayslave/task"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/myzhan/boomer"
	"golang.org/x/crypto/sha3"
)

var (
	endPoint string
	cliPool  clipool.ClientPool

	mutex       sync.Mutex
	initialized = false

	nAcc   int
	accGrp []*account.Account

	readApiCallContract *ReadApiCallContract
	contractAddr        common.Address
	gasPrice            *big.Int

	retValOfCall        *big.Int
	retValOfStorageAt   *big.Int
	retValOfEstimateGas uint64
)

func Init(params *task.Params) {
	mutex.Lock()
	defer mutex.Unlock()

	if initialized {
		return
	}
	initialized = true

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

	deployContract(params.AccGrp[0], endPoint)
	setAnswerVariables()
}

func deployContract(coinbase *account.Account, endPoint string) {
	conn := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(conn)

	auth := bind.NewKeyedTransactor(coinbase.GetKey())
	auth.GasLimit = 999999
	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(coinbase.GetNonce(conn)))

	var tx *types.Transaction
	log.Println("[TC] readApiCallContract: Deploying new smart contract")

	for {
		var err error
		contractAddr, tx, readApiCallContract, err = DeployReadApiCallContract(auth, conn)
		if err == nil {
			coinbase.UpdateNonce()
			break
		}
		log.Printf("[TC] readApiCallContract: Failed to deploy new contract: %v\n", err)
		auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))
		time.Sleep(1 * time.Second) // Avoiding Nonce corruption
	}
	log.Printf("[TC] readApiCallContract: Contract address: 0x%x\n", contractAddr)
	log.Printf("[TC] readApiCallContract: Transaction waiting to be mined: 0x%x\n", tx.Hash())

	ctx := context.Background()
	defer ctx.Done()
	for {
		time.Sleep(500 * time.Millisecond) // Allow it to be processed by the local node :P
		receipt, err := conn.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			//fmt.Printf("Failed to check receipt: %v\n", err)
			continue
		}
		log.Printf("=> Contract Receipt Status: %v\n", receipt.Status)
		break
	}
}

func getMethodId(str string) []byte {
	transferFnSignature := []byte(str)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	return methodID
}

//TODO-klaytn-load-tester: deleting for loop
func setAnswerVariables() {
	retValOfCall = big.NewInt(4)
	retValOfStorageAt = big.NewInt(4)
	for {
		ctx := context.Background()
		cli := cliPool.Alloc().(*client.Client)

		fromAccount := accGrp[rand.Int()%nAcc]
		callMsg := klaytn.CallMsg{
			From:     fromAccount.GetAddress(),
			To:       &contractAddr,
			Gas:      1100000,
			GasPrice: gasPrice,
			Value:    big.NewInt(0),
			Data:     getMethodId("set()"),
		}
		ret, err := cli.EstimateGas(ctx, callMsg)

		if err == nil {
			retValOfEstimateGas = ret
			cliPool.Free(cli)
			break
		} else {
			cli.Close()
		}
	}
}

func sendBoomerEvent(tcName string, logString string, elapsed int64, cli *client.Client, err error) {
	if err == nil {
		boomer.Events.Publish("request_success", "http", tcName+" to "+endPoint, elapsed, int64(10))
		cliPool.Free(cli)
	} else {
		log.Printf("[TC] %s: %s, err=%v\n", tcName, logString, err)
		boomer.Events.Publish("request_failure", "http", tcName+" to "+endPoint, elapsed, err.Error())
		cli.Close()
	}
}

func GetStorageAt() {
	ctx := context.Background()
	cli := cliPool.Alloc().(*client.Client)

	start := boomer.Now()
	ret, err := cli.StorageAt(ctx, contractAddr, common.Hash{}, nil)
	elapsed := boomer.Now() - start

	if err == nil && new(big.Int).SetBytes(ret).Cmp(retValOfStorageAt) != 0 {
		err = errors.New("wrong storage value: " + string(ret) + ", answer: " + retValOfStorageAt.String())
	}
	sendBoomerEvent("readGetStorageAt", "Failure to call klay_getStorageAt", elapsed, cli, err)
}

func Call() {
	cli := cliPool.Alloc().(*client.Client)

	fromAccount := accGrp[rand.Int()%nAcc]
	var callopts bind.CallOpts
	callopts.Pending = false
	callopts.From = fromAccount.GetAddress()

	start := boomer.Now()
	ret, err := readApiCallContract.ReadApiCallContractCaller.Get(&callopts)
	elapsed := boomer.Now() - start

	if err == nil && ret.Cmp(retValOfCall) != 0 {
		err = errors.New("wrong call: " + ret.String() + ", answer: " + retValOfCall.String())
	}
	sendBoomerEvent("readCall", "Failed to call klay_call", elapsed, cli, err)
}

func EstimateGas() {
	ctx := context.Background()
	cli := cliPool.Alloc().(*client.Client)

	fromAccount := accGrp[rand.Int()%nAcc]
	callMsg := klaytn.CallMsg{
		From:     fromAccount.GetAddress(),
		To:       &contractAddr,
		Gas:      1100000,
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
		Data:     getMethodId("set()"),
	}
	start := boomer.Now()
	ret, err := cli.EstimateGas(ctx, callMsg)
	elapsed := boomer.Now() - start

	if err == nil && ret != retValOfEstimateGas {
		err = errors.New("wrong estimate gas: " + strconv.Itoa(int(ret)) + ", answer: " + strconv.Itoa(int(retValOfEstimateGas)))
	}
	sendBoomerEvent("readEstimateGas", "Failed to call klay_estimateGas", elapsed, cli, err)
}
