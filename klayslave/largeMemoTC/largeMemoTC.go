//go:generate abigen --sol LargeMemo.sol --pkg largeMemoTC --out LargeMemo.go

// Package largeMemoTC is used to test required network bandwidth for large block sizes.
// tries to simulate bots which exhausts resource
// See README.md for more details.
package largeMemoTC

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
	"github.com/myzhan/boomer"
)

const Name = "largeMemoTC"
const Letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var (
	endPoint string

	nAcc   int
	accGrp []*account.Account

	gasPrice *big.Int

	cliPool clipool.ClientPool

	mutex       sync.Mutex
	initialized = false

	gLargeMemo *LargeMemo
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
				log.Fatalf("[LargeMemo] Failed to connect to %s, err=%v", endPoint, err)
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
		log.Fatal("[LargeMemo] conn is not client.Client")
		return
	}

	auth := bind.NewKeyedTransactor(coinbase.GetKey())
	auth.GasLimit = 999999
	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))

	var tx *types.Transaction

	for {
		var err error
		_, tx, gLargeMemo, err = DeployLargeMemo(auth, conn)
		if err == nil {
			coinbase.UpdateNonce()
			break
		}

		log.Printf("[LargeMemo] Failed to deploy the contract, err=%v\n", err)
		auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))
		time.Sleep(1 * time.Second)
	}

	ctx := context.Background()
	defer ctx.Done()
	for {
		time.Sleep(500 * time.Millisecond)
		receipt, err := conn.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			log.Printf("[LargeMemo] Failed to check receipt: %v\n", err)
			continue
		}
		if receipt.Status == types.ReceiptStatusSuccessful {
			break
		} else {
			log.Fatalf("[LargeMemo] Contract Receipt Status: %v\n", receipt.Status)
		}
	}
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = Letters[r.Intn(len(Letters))]
	}
	return string(b)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func Run() {
	funcName := "largeMemo"

	conn := cliPool.Alloc().(*client.Client)

	fromAccount := accGrp[rand.Int()%nAcc]
	from := fromAccount.GetKey()

	auth := bind.NewKeyedTransactor(from)
	auth.GasLimit = 999999
	auth.GasPrice = gasPrice

	fromAccount.Lock()

	nonce := fromAccount.GetNonce(conn)
	auth.Nonce = big.NewInt(int64(nonce))

	var err error

	str := randomString(randInt(50, 2000))

	start := boomer.Now()
	_, err = gLargeMemo.SetName(auth, str)
	elapsed := boomer.Now() - start

	if err != nil {
		log.Printf("[LargeMemo] Failed to call %s(), err=%v\n", funcName, err)
		fromAccount.GetNonceFromBlock(conn)
	} else {
		fromAccount.UpdateNonce()
	}

	fromAccount.UnLock()

	// Uncomment the below for debugging
	//if err == nil {
	//	utils.CheckReceipt(conn, tx.Hash())
	//}

	msg := "LargeMemo" + " to " + endPoint
	if err == nil {
		boomer.Events.Publish("request_success", "contract", msg, elapsed, int64(10))
		cliPool.Free(conn)
	} else {
		boomer.Events.Publish("request_failure", "contract", msg, elapsed, err.Error())
		log.Printf("[LargeMemo] request_failure of msg %s, err=%v\n", msg, err)

		conn.Close()
	}
}
