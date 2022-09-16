//go:generate abigen --sol IOHeavy.sol --pkg ioHeavyTC --out IOHeavy.go

// Package ioHeavyTC implements BlockBench's IOHeavy benchmark for Klaytn and Locust.
// See README.md for more details.
package ioHeavyTC

import (
	"context"
	"errors"
	"log"
	"math/big"
	"math/rand"
	"sync"
	"sync/atomic"
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

const Name = "ioHeavyTC"

var (
	endPoint string

	nAcc   int
	accGrp []*account.Account

	gasPrice *big.Int

	cliPool clipool.ClientPool

	mutex       sync.Mutex
	initialized = false

	gIOHeavy *IOHeavy
	gSig     int64 // should be updated atomically
)

const maxKey = 100000
const writeSize = 100 // TODO: fixed size vs. random size
const scanSize = 100  // TODO: fixed size vs. random size

const (
	testWrite = iota
	testScan
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
				log.Fatalf("[IOHeavy] Failed to connect to %s, err=%v", endPoint, err)
			}
			return c
		}

		cliPool.Init(20, 300, cliCreate)

		for _, acc := range params.AccGrp {
			accGrp = append(accGrp, acc)
		}

		nAcc = len(accGrp)

		deployContract(params.AccGrp[0])
		gSig = 0
	}
}

func deployContract(coinbase *account.Account) {
	conn, ok := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(conn)
	if !ok {
		log.Fatal("[IOHeavy] conn is not client.Client")
		return
	}

	auth := bind.NewKeyedTransactor(coinbase.GetKey())
	auth.GasLimit = 9999999
	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))

	var address common.Address
	var tx *types.Transaction
	log.Println("[IOHeavy] Deploying a new smart contract")

	for {
		var err error
		address, tx, gIOHeavy, err = DeployIOHeavy(auth, conn)
		if err == nil {
			coinbase.UpdateNonce()
			break
		}

		log.Printf("[IOHeavy] Failed to deploy the contract, err=%v\n", err)
		auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))
		time.Sleep(1 * time.Second)
	}
	log.Printf("[IOHeavy] Contract address : 0x%x\n", address)
	log.Printf("[IOHeavy] Transaction waiting to be mined: 0x%x\n", tx.Hash())

	ctx := context.Background()
	defer ctx.Done()
	for {
		time.Sleep(500 * time.Millisecond)
		receipt, err := conn.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			log.Printf("[IOHeavy] Failed to check receipt: %v\n", err)
			continue
		}
		log.Println("[IOHeavy] Received the receipt")
		if receipt.Status == types.ReceiptStatusSuccessful {
			log.Println("[IOHeavy] Contract deployment was successful")
			break
		} else {
			log.Fatalf("[IOHeavy] Contract Receipt Status: %v\n", receipt.Status)
		}
	}
}

func toString(target int) string {
	switch target {
	case testWrite:
		return "Write"
	case testScan:
		return "Scan"
	default:
		return "Unknown"
	}
}

// Run randomly calls Write() or Scan().
func Run() {
	target := rand.Int() % testLast
	log.Printf("[IOHeavy] calling %s()...\n", toString(target))
	callFunc(target)
}

// Write tests the write function in the IOHeavy contract.
func Write() {
	callFunc(testWrite)
}

// Scan tests the scan function in the IOHeavy contract.
func Scan() {
	callFunc(testScan)
}

func callFunc(target int) {
	var size int64

	// Check if target is valid
	switch target {
	case testWrite:
		size = writeSize
	case testScan:
		size = scanSize
	default:
		log.Printf("[IOHeavy] Unknown target: %d\n", target)
		boomer.Events.Publish("request_failure", "contract", "ioHeavy/"+string(target)+" to "+endPoint, 0, "Unknown target")
		return
	}

	// Get the function name as a string
	funcName := toString(target)

	// Choose the start key randomly
	startKey := rand.Int63() % maxKey

	// Signature to distinguish txs
	sig := atomic.AddInt64(&gSig, 1)

	conn := cliPool.Alloc().(*client.Client)

	fromAccount := accGrp[rand.Int()%nAcc]
	from := fromAccount.GetKey()

	auth := bind.NewKeyedTransactor(from)
	auth.GasLimit = 99999999
	auth.GasPrice = gasPrice

	fromAccount.Lock()

	nonce := fromAccount.GetNonce(conn)
	auth.Nonce = big.NewInt(int64(nonce))

	log.Printf("[IOHeavy] from=%s nonce=%d %s(startKey=%d, size=%d, sig=%d)\n",
		fromAccount.GetAddress().String(), nonce, funcName, startKey, size, sig)

	var tx *types.Transaction
	var err error

	start := boomer.Now()
	switch target {
	case testWrite:
		tx, err = gIOHeavy.Write(auth, big.NewInt(startKey), big.NewInt(size), big.NewInt(sig))
	case testScan:
		tx, err = gIOHeavy.Scan(auth, big.NewInt(startKey), big.NewInt(size), big.NewInt(sig))
	default:
		log.Printf("[IOHeavy] target %d (%s) is not handled.\n", target, funcName)
		err = errors.New("unhandled target")
	}
	elapsed := boomer.Now() - start

	if err != nil {
		log.Printf("[IOHeavy] Failed to call %s(), err=%v\n", funcName, err)
		fromAccount.GetNonceFromBlock(conn)
	} else {
		log.Printf("[IOHeavy] %s tx=%s\n", funcName, tx.Hash().String())
		fromAccount.UpdateNonce()
	}

	fromAccount.UnLock()

	// Uncomment the below for debugging
	//if err == nil {
	//	utils.CheckReceipt(conn, tx.Hash())
	//}

	msg := "ioHeavy/" + funcName + " to " + endPoint
	if err == nil {
		boomer.Events.Publish("request_success", "contract", msg, elapsed, int64(10))
		cliPool.Free(conn)
	} else {
		boomer.Events.Publish("request_failure", "contract", msg, elapsed, err.Error())
		conn.Close()
	}
}
