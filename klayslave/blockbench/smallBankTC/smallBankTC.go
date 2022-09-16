//go:generate abigen --sol SmallBank.sol --pkg smallBankTC --out SmallBank.go

// Package smallBankTC implements BlockBench's SmallBank benchmark for Klaytn and Locust.
// See README.md for more details.
package smallBankTC

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

const Name = "smallBankTC"

var (
	endPoint string

	nAcc   int
	accGrp []*account.Account

	gasPrice *big.Int

	cliPool clipool.ClientPool

	mutex       sync.Mutex
	initialized = false

	gSmallBank *SmallBank

	// maxNumUsers determines the maximum number of users used in tests.
	maxNumUsers = 100000
)

const (
	testAlmagate = iota
	testGetBalance
	testUpdateBalance
	testUpdateSaving
	testSendPayment
	testWriteCheck
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
				log.Fatalf("[SmallBank] Failed to connect to %s, err=%v", endPoint, err)
			}
			return c
		}

		cliPool.Init(20, 300, cliCreate)

		for _, acc := range params.AccGrp {
			accGrp = append(accGrp, acc)
		}

		nAcc = len(accGrp)

		deployContract(params.AccGrp[0])

		// Change the maximum number of users if the environment variable SMALLBANK_MAX_NUM_USERS has been set
		if v := os.Getenv("SMALLBANK_MAX_NUM_USERS"); v != "" {
			if envVal, err := strconv.Atoi(v); err == nil {
				maxNumUsers = envVal
			}
		}
		log.Printf("[SmallBank] maxNumUsers=%d\n", maxNumUsers)
	}
}

func deployContract(coinbase *account.Account) {
	conn, ok := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(conn)
	if !ok {
		log.Fatal("[SmallBank] conn is not client.Client")
		return
	}

	auth := bind.NewKeyedTransactor(coinbase.GetKey())
	auth.GasLimit = 9999999
	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))

	var address common.Address
	var tx *types.Transaction
	log.Println("[SmallBank] Deploying a new smart contract")

	for {
		var err error
		address, tx, gSmallBank, err = DeploySmallBank(auth, conn)
		if err == nil {
			coinbase.UpdateNonce()
			break
		}

		log.Printf("[SmallBank] Failed to deploy the contract, err=%v\n", err)
		auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))
		time.Sleep(1 * time.Second)
	}
	log.Printf("[SmallBank] Contract address : 0x%x\n", address)
	log.Printf("[SmallBank] Transaction waiting to be mined: 0x%x\n", tx.Hash())

	ctx := context.Background()
	defer ctx.Done()
	for {
		time.Sleep(500 * time.Millisecond)
		receipt, err := conn.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			log.Printf("[SmallBank] Failed to check receipt: %v\n", err)
			continue
		}
		log.Println("[SmallBank] Received the receipt")
		if receipt.Status == types.ReceiptStatusSuccessful {
			log.Println("[SmallBank] Contract deployment was successful")
			break
		} else {
			log.Fatalf("[SmallBank] Contract Receipt Status: %v\n", receipt.Status)
		}
	}
}

func toString(target int) string {
	switch target {
	case testAlmagate:
		return "Almagate"
	case testGetBalance:
		return "GetBalance"
	case testUpdateBalance:
		return "UpdateBalance"
	case testUpdateSaving:
		return "UpdateSaving"
	case testSendPayment:
		return "SendPayment"
	case testWriteCheck:
		return "WriteCheck"
	default:
		return "Unknown"
	}
}

// Run randomly calls one test case.
func Run() {
	target := rand.Int() % testLast
	log.Printf("[SmallBank] calling %s()...\n", toString(target))
	callFunc(target)
}

// Almagate tests the almagate function in the SmallBank contract.
func Almagate() {
	callFunc(testAlmagate)
}

// GetBalance tests the getBalance function in the SmallBank contract.
func GetBalance() {
	callFunc(testGetBalance)
}

// UpdateBalance tests the updateBalance function in the SmallBank contract.
func UpdateBalance() {
	callFunc(testUpdateBalance)
}

// UpdateSaving tests the updateSaving function in the SmallBank contract.
func UpdateSaving() {
	callFunc(testUpdateSaving)
}

// SendPayment tests the sendPayment function in the SmallBank contract.
func SendPayment() {
	callFunc(testSendPayment)
}

// WriteCheck tests the writeCheck function in the SmallBank contract.
func WriteCheck() {
	callFunc(testWriteCheck)
}

func callFunc(target int) {
	// Check if target is valid
	if target < testAlmagate || target >= testLast {
		log.Printf("[SmallBank] Unknown target: %d\n", target)
		boomer.Events.Publish("request_failure", "contract", "smallBank/"+string(target)+" to "+endPoint, 0, "Unknown target")
		return
	}

	// Get the function name as a string
	funcName := toString(target)

	// Prepare function parameters
	user1 := "user" + strconv.Itoa(rand.Int()%maxNumUsers)
	user2 := "user" + strconv.Itoa(rand.Int()%maxNumUsers)
	//log.Printf("[SmallBank] users: %s %s\n", user1, user2)

	conn := cliPool.Alloc().(*client.Client)

	fromAccount := accGrp[rand.Int()%nAcc]
	from := fromAccount.GetKey()

	var err error
	var elapsed int64

	if target == testGetBalance {
		var balance *big.Int
		callOpts := &bind.CallOpts{Pending: false, From: fromAccount.GetAddress(), Context: nil}

		start := boomer.Now()
		balance, err = gSmallBank.GetBalance(callOpts, user1)
		if err == nil {
			log.Printf("[SmallBank] %s(%s)=%v\n", funcName, user1, balance)
		} else {
			log.Printf("[SmallBank] Failed to call %s(), err=%v\n", funcName, err)
		}
		elapsed = boomer.Now() - start
	} else {
		auth := bind.NewKeyedTransactor(from)
		auth.GasLimit = 9999999
		auth.GasPrice = gasPrice

		fromAccount.Lock()

		nonce := fromAccount.GetNonce(conn)
		auth.Nonce = big.NewInt(int64(nonce))

		log.Printf("[SmallBank] from=%s nonce=%d %s()\n", fromAccount.GetAddress().String(), nonce, funcName)

		var tx *types.Transaction

		start := boomer.Now()
		switch target {
		case testAlmagate:
			tx, err = gSmallBank.Almagate(auth, user1, user2)
		case testUpdateBalance:
			// TODO: use a more meaningful value for the new balance
			tx, err = gSmallBank.UpdateBalance(auth, user1, big.NewInt(0))
		case testUpdateSaving:
			// TODO: use a more meaningful value for the new balance
			tx, err = gSmallBank.UpdateSaving(auth, user1, big.NewInt(0))
		case testSendPayment:
			// TODO: use a more meaningful send value
			tx, err = gSmallBank.SendPayment(auth, user1, user2, big.NewInt(0))
		case testWriteCheck:
			// TODO: use a more meaningful check value
			tx, err = gSmallBank.WriteCheck(auth, user1, big.NewInt(0))
		default:
			log.Printf("[SmallBank] target %d (%s) is not handled.\n", target, funcName)
			err = errors.New("unhandled target")
		}
		elapsed = boomer.Now() - start

		if err != nil {
			log.Printf("[SmallBank] Failed to call %s(), err=%v\n", funcName, err)
			fromAccount.GetNonceFromBlock(conn)
		} else {
			log.Printf("[SmallBank] %s tx=%s\n", funcName, tx.Hash().String())
			fromAccount.UpdateNonce()
		}

		fromAccount.UnLock()

		// Uncomment the below for debugging
		//if err == nil {
		//	utils.CheckReceipt(conn, tx.Hash())
		//}
	}

	msg := "smallBank/" + funcName + " to " + endPoint
	if err == nil {
		boomer.Events.Publish("request_success", "contract", msg, elapsed, int64(10))
		cliPool.Free(conn)
	} else {
		boomer.Events.Publish("request_failure", "contract", msg, elapsed, err.Error())
		conn.Close()
	}
}
