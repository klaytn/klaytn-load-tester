//go:generate abigen --sol kvstore.sol --pkg ycsbTC --out kvstore.go

// Package ycsbTC implements a simplified YCSB benchmark.
// See README.md for more details.
package ycsbTC

import (
	"context"
	"errors"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strconv"
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

const Name = "ycsbTC"

var (
	endPoint string

	nAcc   int
	accGrp []*account.Account

	gasPrice *big.Int

	cliPool clipool.ClientPool

	mutex       sync.Mutex
	initialized = false

	gKVstore *KVstore

	// maxNumKeys determines the maximum number of keys used in tests.
	maxNumKeys = 100000
)

const (
	testSet = iota
	testGet
	testLast
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
				log.Fatalf("[YCSB] Failed to connect to %s, err=%v", endPoint, err)
			}
			return c
		}

		cliPool.Init(20, 300, cliCreate)

		for _, acc := range params.AccGrp {
			accGrp = append(accGrp, acc)
		}

		nAcc = len(accGrp)

		deployContract(params.AccGrp[0])

		// Change the maximum number of keys if the environment variable YCSB_MAX_NUM_KEYS has been set
		if v := os.Getenv("YCSB_MAX_NUM_KEYS"); v != "" {
			if envVal, err := strconv.Atoi(v); err == nil {
				maxNumKeys = envVal
			}
		}
		log.Printf("[YCSB] maxNumKeys=%d\n", maxNumKeys)
	}
}

func deployContract(coinbase *account.Account) {
	conn, ok := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(conn)
	if !ok {
		log.Fatal("[YCSB] conn is not client.Client")
		return
	}

	auth := bind.NewKeyedTransactor(coinbase.GetKey())
	auth.GasLimit = 999999
	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))

	var address common.Address
	var tx *types.Transaction
	log.Println("[YCSB] Deploying a new smart contract")

	for {
		var err error
		address, tx, gKVstore, err = DeployKVstore(auth, conn)
		if err == nil {
			coinbase.UpdateNonce()
			break
		}

		log.Printf("[YCSB] Failed to deploy the contract, err=%v\n", err)
		auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))
		time.Sleep(1 * time.Second)
	}
	log.Printf("[YCSB] Contract address : 0x%x\n", address)
	log.Printf("[YCSB] Transaction waiting to be mined: 0x%x\n", tx.Hash())

	ctx := context.Background()
	defer ctx.Done()
	for {
		time.Sleep(500 * time.Millisecond)
		receipt, err := conn.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			log.Printf("[YCSB] Failed to check receipt: %v\n", err)
			continue
		}
		log.Println("[YCSB] Received the receipt")
		if receipt.Status == types.ReceiptStatusSuccessful {
			log.Println("[YCSB] Contract deployment was successful")
			break
		} else {
			log.Fatalf("[YCSB] Contract Receipt Status: %v\n", receipt.Status)
		}
	}
}

func toString(target int) string {
	switch target {
	case testSet:
		return "Set"
	case testGet:
		return "Get"
	default:
		return "Unknown"
	}
}

// Run randomly calls one test case.
func Run() {
	target := rand.Int() % testLast
	log.Printf("[YCSB] calling %s()...\n", toString(target))
	callFunc(target)
}

// Set tests the set function in the KVstore contract.
func Set() {
	callFunc(testSet)
}

// Get tests the get function in the KVstore contract.
func Get() {
	callFunc(testGet)
}

func callFunc(target int) {
	// Check if target is valid
	if target < testSet || target >= testLast {
		log.Printf("[YCSB] Unknown target: %d\n", target)
		boomer.Events.Publish("request_failure", "contract", "ycsb/"+string(target)+" to "+endPoint, 0, "Unknown target")
		return
	}

	// Get the function name as a string
	funcName := toString(target)

	// Prepare function parameters
	user := "user" + strconv.Itoa(rand.Int()%maxNumKeys)
	val := "val" + strconv.Itoa(rand.Int())

	conn := cliPool.Alloc().(*client.Client)

	fromAccount := accGrp[rand.Int()%nAcc]
	from := fromAccount.GetKey()

	var err error
	var elapsed int64

	switch target {
	case testSet:
		auth := bind.NewKeyedTransactor(from)
		auth.GasLimit = 9999999
		auth.GasPrice = gasPrice

		fromAccount.Lock()

		nonce := fromAccount.GetNonce(conn)
		auth.Nonce = big.NewInt(int64(nonce))

		log.Printf("[YCSB] from=%s nonce=%d %s()\n", fromAccount.GetAddress().String(), nonce, funcName)
		log.Printf("[YCSB] %s(%s, %s)\n", funcName, user, val)

		var tx *types.Transaction

		start := boomer.Now()
		tx, err = gKVstore.Set(auth, user, val)
		elapsed = boomer.Now() - start

		if err != nil {
			log.Printf("[YCSB] Failed to call %s(), err=%v\n", funcName, err)
			fromAccount.GetNonceFromBlock(conn)
		} else {
			log.Printf("[YCSB] %s tx=%s\n", funcName, tx.Hash().String())
			fromAccount.UpdateNonce()
		}

		fromAccount.UnLock()

		// Uncomment the below for debugging
		//if err == nil {
		//	utils.CheckReceipt(conn, tx.Hash())
		//}

	case testGet:
		var value string
		callOpts := &bind.CallOpts{Pending: false, From: fromAccount.GetAddress(), Context: nil}

		start := boomer.Now()
		value, err = gKVstore.Get(callOpts, user)
		if err == nil {
			log.Printf("[YCSB] %s(%s)=%v\n", funcName, user, value)
		} else {
			log.Printf("[YCSB] Failed to call %s(), err=%v\n", funcName, err)
		}
		elapsed = boomer.Now() - start

	default:
		log.Printf("[YCSB] target %d (%s) is not handled.\n", target, funcName)
		err = errors.New("unhandled target")
	}

	msg := "ycsb/" + funcName + " to " + endPoint
	if err == nil {
		boomer.Events.Publish("request_success", "contract", msg, elapsed, int64(10))
		cliPool.Free(conn)
	} else {
		boomer.Events.Publish("request_failure", "contract", msg, elapsed, err.Error())
		conn.Close()
	}
}
