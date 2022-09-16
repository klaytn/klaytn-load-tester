package cpuHeavyTC

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"time"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/clipool"
	"github.com/klaytn/klaytn-load-tester/klayslave/task"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/myzhan/boomer"
)

const Name = "cpuHeavyTx"

var (
	endPoint string

	nAcc   int
	accGrp []*account.Account

	gCpuHeavy *CPUHeavy

	cliPool clipool.ClientPool

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

	deployContract(params.AccGrp[0], endPoint)
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
	println("Deploying new smart contract")

	for {
		addr, tTx, cpuHeavy, err := DeployCPUHeavy(auth, conn)
		address = addr
		tx = tTx
		if err != nil {
			coinbase.UpdateNonce()
		}
		gCpuHeavy = cpuHeavy

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

func Run() {
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
	_, err := gCpuHeavy.Sort(auth, big.NewInt(100), big.NewInt(1))
	elapsed := boomer.Now() - start

	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() ||
			err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", fromAccount.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", fromAccount.GetAddress().String(), nonce+1)
			fromAccount.UpdateNonce()
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", fromAccount.GetAddress().String(), nonce, err)
		}

		fmt.Printf("Failed to retrieve pending tx: %v\n", err)
		//fromAccount.GetNonceFromBlock(conn)
	} else {
		//fmt.Println("Pending tx:", res.Hash().String())
		fromAccount.UpdateNonce()
	}

	fromAccount.UnLock()

	if err == nil {
		boomer.Events.Publish("request_success", "contract", "cpuHeavy"+" to "+endPoint, elapsed, int64(10))
		cliPool.Free(conn)
	} else {
		boomer.Events.Publish("request_failure", "contract", "cpuHeavy"+" to "+endPoint, elapsed, err.Error())
	}
}

func RunSingle() (tx *types.Transaction, err error) {
	conn := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(conn)

	fromAccount := accGrp[rand.Int()%nAcc]
	from := fromAccount.GetKey()

	fromAccount.Lock()
	defer fromAccount.UnLock()
	auth := bind.NewKeyedTransactor(from)

	// nonce := fromAccount.GetNonce(conn)
	nonce := fromAccount.GetNonceFromBlock(conn)
	fmt.Printf("[TC] cpuHeavyTC/sortSingle(): %v, fromAccount=%v, nonce=%v\n", endPoint, fromAccount, nonce)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = 999999
	auth.GasPrice = gasPrice

	tx, err = gCpuHeavy.SortSingle(auth, big.NewInt(1))
	fmt.Printf("[TC] cpuHeavyTC: %v, tx:%v\n", endPoint, tx)

	return
}

// CheckResult returns true and nil error, if expected results are observed.
// Otherewise returns false and error.
func CheckResult() (result bool, err error) {
	conn := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(conn)

	fromAccount := accGrp[rand.Int()%nAcc]

	fromAccount.Lock()
	defer fromAccount.UnLock()

	nonce := fromAccount.GetNonceFromBlock(conn)
	fmt.Printf("[TC] cpuHeavyTC/checkResult(): fromAccount=%v, nonce=%v\n", fromAccount, nonce)

	var callopts bind.CallOpts
	callopts.Pending = false
	callopts.From = fromAccount.GetAddress()
	result, err = gCpuHeavy.CheckResult(&callopts)
	if err != nil {
		return false, err
	}

	return
}
