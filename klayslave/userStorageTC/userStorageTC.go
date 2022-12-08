package userStorageTC

import (
	"context"
	"errors"
	"fmt"
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

const Name = "userStorageTC"

var (
	endPoint string

	nAcc   int
	accGrp []*account.Account

	gUserStorage *UserStorage

	cliPool clipool.ClientPool

	mutex       sync.Mutex
	initialized = false

	gasPrice *big.Int
)

func Init(params *task.Params) {
	mutex.Lock()
	defer mutex.Unlock()
	if !initialized {
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
	}
}

func deployContract(coinbase *account.Account, endPoint string) {
	conn := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(conn)

	auth := bind.NewKeyedTransactor(coinbase.GetKey())
	auth.GasLimit = 999999
	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(coinbase.GetNonce(conn)))

	var address common.Address
	var tx *types.Transaction
	println("Deploying a new smart contract")

	for {
		addr, tTx, userStorage, err := DeployUserStorage(auth, conn)
		address = addr
		tx = tTx
		if err != nil {
			coinbase.UpdateNonce()
		}
		gUserStorage = userStorage

		if err != nil {
			//log.Printf("Failed to deploy new contract: %v\n", err)
		} else {
			break
		}
		auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))
		time.Sleep(1 * time.Second) // Avoiding Nonce corruption
	}
	fmt.Printf("=> Contract pending deploy: 0x%x\n", address)

	fmt.Printf("Transaction waiting to be mined: 0x%x\n", tx.Hash())
	ctx := context.Background()
	defer ctx.Done()
	for {
		time.Sleep(500 * time.Millisecond) // Allow it to be processed by the local node :P
		receipt, err := conn.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			//fmt.Printf("Failed to check receipt: %v\n", err)
			continue
		}
		fmt.Printf("=> Contract Receipt Status: %v\n", receipt.Status)
		break
	}

}

func RunSet() {
	conn := cliPool.Alloc().(*client.Client)

	fromAccount := accGrp[rand.Int()%nAcc]
	from := fromAccount.GetKey()

	fromAccount.Lock()
	auth := bind.NewKeyedTransactor(from)

	nonce := fromAccount.GetNonce(conn)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = 999999
	auth.GasPrice = gasPrice

	start := boomer.Now()
	_, err := gUserStorage.Set(auth, big.NewInt(100))
	elapsed := boomer.Now() - start

	if err != nil {
		fmt.Printf("Failed to retrieve pending tx: %v\n", err)
		fromAccount.GetNonceFromBlock(conn)
	} else {
		fromAccount.UpdateNonce()
	}

	fromAccount.UnLock()

	if err == nil {
		boomer.Events.Publish("request_success", "contract", "userStorage/Set"+" to "+endPoint, elapsed, int64(10))
		cliPool.Free(conn)
	} else {
		boomer.Events.Publish("request_failure", "contract", "userStorage/Set"+" to "+endPoint, elapsed, err.Error())
	}
}

func RunSetGet() {
	conn := cliPool.Alloc().(*client.Client)

	fromAccount := accGrp[rand.Int()%nAcc]
	from := fromAccount.GetKey()

	fromAccount.Lock()

	start := boomer.Now()

	// Call Set()
	auth := bind.NewKeyedTransactor(from)
	auth.Nonce = big.NewInt(int64(fromAccount.GetNonce(conn)))
	auth.GasLimit = 999999
	auth.GasPrice = gasPrice

	value := auth.Nonce

	tx, err := gUserStorage.Set(auth, value)
	elapsed := boomer.Now() - start

	if err != nil {
		fmt.Printf("Failed to retrieve pending tx: %v\n", err)
		fromAccount.GetNonceFromBlock(conn)

		fromAccount.UnLock()

		boomer.Events.Publish("request_failure", "contract", "userStorage/SetGet"+" to "+endPoint, elapsed, err.Error())

		return
	}

	start = boomer.Now()

	// Increment fromAccount's nonce
	fromAccount.UpdateNonce()

	ctx := context.Background()
	defer ctx.Done()

	// Wait for the receipt to be available
	for {
		time.Sleep(500 * time.Millisecond)
		receipt, receiptErr := conn.TransactionReceipt(ctx, tx.Hash())
		if receiptErr != nil {
			continue
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			err = errors.New(fmt.Sprintf("tx=%v: from=%v failed=%v", tx.Hash().String(), fromAccount.GetAddress().String(), receipt.Status))
			fmt.Println(err.Error())
		}
		break
	}

	if err == nil {
		time.Sleep(1500 * time.Millisecond)

		// Wait until tx is included in the block
		for {
			_, isPending, _ := conn.TransactionByHash(ctx, tx.Hash())
			if isPending {
				time.Sleep(5 * time.Millisecond)
			} else {
				break
			}
		}

		// Call Get() to retrieve the value set by Set()
		var callopts bind.CallOpts
		callopts.Pending = false
		callopts.From = fromAccount.GetAddress()
		result, getErr := gUserStorage.Get(&callopts)
		if getErr != nil {
			err = getErr
		} else if result.Cmp(value) != 0 {
			err = errors.New(fmt.Sprintf("tx=%v: from=%v, incorrect value (received=%v vs. expected=%v)", tx.Hash().String(), callopts.From.String(), result, value))
			fmt.Println(err.Error())
		}
	}

	elapsed += boomer.Now() - start

	fromAccount.UnLock()

	if err == nil {
		boomer.Events.Publish("request_success", "contract", "userStorage/SetGet"+" to "+endPoint, elapsed, int64(10))
		cliPool.Free(conn)
	} else {
		boomer.Events.Publish("request_failure", "contract", "userStorage/SetGet"+" to "+endPoint, elapsed, err.Error())
	}
}

func RunSetSingle() (tx *types.Transaction, err error) {
	conn := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(conn)

	fromAccount := accGrp[rand.Int()%nAcc]
	from := fromAccount.GetKey()

	fromAccount.Lock()
	defer fromAccount.UnLock()
	auth := bind.NewKeyedTransactor(from)

	// nonce := fromAccount.GetNonce(conn)
	nonce := fromAccount.GetNonceFromBlock(conn)
	fmt.Printf("[TC] userStorageTC/Set: fromAccount=%v, nonce=%v\n", fromAccount, nonce)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = 999999
	auth.GasPrice = gasPrice

	tx, err = gUserStorage.Set(auth, big.NewInt(100))
	fmt.Printf("[TC] userStorage/Set: %v, tx:%v\n", endPoint, tx)

	if err != nil {
		return nil, err
	}

	return
}

func RunSetGetSingle() (err error) {
	conn := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(conn)

	fromAccount := accGrp[rand.Int()%nAcc]
	from := fromAccount.GetKey()

	fromAccount.Lock()
	defer fromAccount.UnLock()

	// Call Set()
	auth := bind.NewKeyedTransactor(from)
	nonce := fromAccount.GetNonceFromBlock(conn)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = 999999
	auth.GasPrice = gasPrice

	value := auth.Nonce

	fmt.Printf("[TC] userStorageTC/SetGet Set: fromAccount=%v, nonce=%v\n", fromAccount, nonce)

	tx, err := gUserStorage.Set(auth, value)
	fmt.Printf("[TC] userStorageTC/SetGet Set: %v, auth:%v, tx:%v\n", endPoint, auth, tx)

	if err != nil {
		return
	}

	// Increment fromAccount's nonce
	fromAccount.UpdateNonce()

	ctx := context.Background()
	defer ctx.Done()

	// Wait for the receipt to be available
	for {
		time.Sleep(500 * time.Millisecond)
		receipt, receiptErr := conn.TransactionReceipt(ctx, tx.Hash())
		if receiptErr != nil {
			continue
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			err = errors.New(fmt.Sprintf("tx=%v: from=%v failed=%v", tx.Hash().String(), fromAccount.GetAddress().String(), receipt.Status))
			fmt.Println(err.Error())
		}
		break
	}

	if err == nil {
		time.Sleep(1500 * time.Millisecond)

		// Wait until tx is included in the block
		for {
			_, isPending, _ := conn.TransactionByHash(ctx, tx.Hash())
			if isPending {
				time.Sleep(5 * time.Millisecond)
			} else {
				break
			}
		}

		// Call Get() to retrieve the value set by Set()
		var callopts bind.CallOpts
		callopts.Pending = false
		callopts.From = fromAccount.GetAddress()
		fmt.Printf("[TC] userStorageTC/SetGet Get: %v, From:%s\n", endPoint, callopts.From.String())
		result, getErr := gUserStorage.Get(&callopts)

		if getErr != nil {
			err = getErr
		} else if result.Cmp(value) != 0 {
			err = errors.New(fmt.Sprintf("tx=%v: from=%v, incorrect value (received=%v vs. expected=%v)", tx.Hash().String(), callopts.From.String(), result, value))
			fmt.Println(err.Error())
		}
	}

	return
}
