//go:generate abigen --sol DoNothing.sol --pkg doNothingTC --out DoNothing.go

// Package doNothingTC implements BlockBench's DoNothing benchmark for Klaytn and Locust.
// See README.md for more details.
package doNothingTC

import (
	"context"
	"log"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/clipool"
	"github.com/klaytn/klaytn-load-tester/klayslave/task"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/myzhan/boomer"
)

const Name = "doNothingTC"

var (
	endPoint string

	nAcc   int
	accGrp []*account.Account

	gasPrice *big.Int

	cliPool clipool.ClientPool

	mutex       sync.Mutex
	initialized = false

	gDoNothing *DoNothing
)

// Init initializes cliPool and accGrp; and also deploys the smart contract.
func Init(params *task.Params) {
	mutex.Lock()
	defer mutex.Unlock()

	if !initialized {
		initialized = true

		endPoint = params.Endpoint
		gasPrice = params.GasPrice

		cliCreate := func() interface{} {
			c, err := client.Dial(endPoint)
			if err != nil {
				log.Fatalf("[DoNothing] Failed to connect to %s, err=%v", endPoint, err)
			}
			return c
		}

		cliPool.Init(20, 300, cliCreate)

		for _, acc := range params.AccGrp {
			accGrp = append(accGrp, acc)
		}

		nAcc = len(accGrp)

		deployContract(params.AccGrp[0])
	}
}

func deployContract(coinbase *account.Account) {
	conn, ok := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(conn)
	if !ok {
		log.Fatal("[DoNothing] conn is not client.Client")
		return
	}

	auth := bind.NewKeyedTransactor(coinbase.GetKey())
	auth.GasLimit = 999999
	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))

	var address common.Address
	var tx *types.Transaction
	log.Println("[DoNothing] Deploying a new smart contract")

	for {
		var err error
		address, tx, gDoNothing, err = DeployDoNothing(auth, conn)
		if err == nil {
			coinbase.UpdateNonce()
			break
		}

		log.Printf("[DoNothing] Failed to deploy the contract, err=%v\n", err)
		auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))
		time.Sleep(1 * time.Second)
	}
	log.Printf("[DoNothing] Contract address : 0x%x\n", address)
	log.Printf("[DoNothing] Transaction waiting to be mined: 0x%x\n", tx.Hash())

	ctx := context.Background()
	defer ctx.Done()
	for {
		time.Sleep(500 * time.Millisecond)
		receipt, err := conn.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			log.Printf("[DoNothing] Failed to check receipt: %v\n", err)
			continue
		}
		log.Println("[DoNothing] Received the receipt")
		if receipt.Status == types.ReceiptStatusSuccessful {
			log.Println("[DoNothing] Contract deployment was successful")
			break
		} else {
			log.Fatalf("[DoNothing] Contract Receipt Status: %v\n", receipt.Status)
		}
	}
}

func Run() {
	funcName := "Nothing"

	conn := cliPool.Alloc().(*client.Client)

	fromAccount := accGrp[rand.Int()%nAcc]
	from := fromAccount.GetKey()

	auth := bind.NewKeyedTransactor(from)
	auth.GasLimit = 999999
	auth.GasPrice = gasPrice

	fromAccount.Lock()

	nonce := fromAccount.GetNonce(conn)
	auth.Nonce = big.NewInt(int64(nonce))

	log.Printf("[DoNothing] from=%s nonce=%d %s()\n", fromAccount.GetAddress().String(), nonce, funcName)

	var tx *types.Transaction
	var err error

	start := boomer.Now()
	tx, err = gDoNothing.Nothing(auth)
	elapsed := boomer.Now() - start

	if err != nil {
		log.Printf("[DoNothing] Failed to call %s(), err=%v\n", funcName, err)
		fromAccount.GetNonceFromBlock(conn)
	} else {
		log.Printf("[DoNothing] %s tx=%s\n", funcName, tx.Hash().String())
		fromAccount.UpdateNonce()
	}

	fromAccount.UnLock()

	// Uncomment the below for debugging
	//if err == nil {
	//	utils.CheckReceipt(conn, tx.Hash())
	//}

	msg := "doNothing" + " to " + endPoint
	if err == nil {
		boomer.Events.Publish("request_success", "contract", msg, elapsed, int64(10))
		cliPool.Free(conn)
	} else {
		boomer.Events.Publish("request_failure", "contract", msg, elapsed, err.Error())
		conn.Close()
	}
}
