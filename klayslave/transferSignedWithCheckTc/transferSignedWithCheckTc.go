package transferSignedWithCheckTc

import (
	"context"
	"errors"
	"log"
	"math/big"
	"math/rand"
	"time"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/clipool"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/myzhan/boomer"
)

const Name = "transferSignedWithCheckTx"

var (
	endPoint string
	nAcc     int
	accGrp   []*account.Account
	cliPool  clipool.ClientPool
	gasPrice *big.Int
)

func Init(accs []*account.Account, endpoint string, gp *big.Int) {
	gasPrice = gp

	endPoint = endpoint

	cliCreate := func() interface{} {
		c, err := client.Dial(endPoint)
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

func doubleLock(to *account.Account, from *account.Account) {
	if from.GetAddress().String() == to.GetAddress().String() {
		from.Lock()
	} else if from.GetAddress().String() > to.GetAddress().String() {
		from.Lock()
		to.Lock()
	} else {
		to.Lock()
		from.Lock()
	}
}

func doubleUnlock(to *account.Account, from *account.Account) {
	if from.GetAddress().String() == to.GetAddress().String() {
		from.UnLock()
	} else if from.GetAddress().String() > to.GetAddress().String() {
		from.UnLock()
		to.UnLock()
	} else {
		to.UnLock()
		from.UnLock()
	}
}

func TransferAndCheck(cli *client.Client, to *account.Account, from *account.Account, value *big.Int) error {
	ctx := context.Background()

	doubleLock(to, from)
	defer doubleUnlock(to, from)
	// The reason of saving balance of current accounts is to comparing with later balance.
	fromFormerBalance, _ := from.GetBalance(cli)
	toFormerBalance, _ := to.GetBalance(cli)

	hash, gasPrice, err := from.TransferSignedTxWithoutLock(cli, to, value)
	if err != nil {
		return err
	}
	startTime := time.Now().Unix()
	var receipt *types.Receipt
	for {
		receipt, _ = cli.TransactionReceipt(ctx, hash)
		if receipt != nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
		if time.Now().Unix()-startTime > 100 {
			return errors.New("Time out : It took more than 100 seconds to make a block ")
		}
	}

	if to.GetAddress() == from.GetAddress() {
		value.SetUint64(0)
	}

	fromFormerBalance.Sub(fromFormerBalance, value)
	gasUsed := big.NewInt((int64)(receipt.GasUsed))
	fee := new(big.Int).Mul(gasUsed, gasPrice)
	fromFormerBalance.Sub(fromFormerBalance, fee)
	toFormerBalance.Add(toFormerBalance, value)

	startTime = time.Now().Unix()
	for {
		errFrom := from.CheckBalance(fromFormerBalance, cli)
		if errFrom != nil {
			log.Printf("from account : %s", errFrom.Error())
			time.Sleep(100 * time.Millisecond)
			if time.Now().Unix()-startTime > 10 {
				return errors.New("Time out (from) : It took more than 10 seconds to retrieve the correct receipt ")
			}
		} else {
			break
		}
	}

	if from.GetAddress() == to.GetAddress() {
		return nil
	}

	startTime = time.Now().Unix()
	for {
		errTo := to.CheckBalance(toFormerBalance, cli)
		if errTo != nil {
			log.Printf("to account : %s", errTo.Error())
			time.Sleep(100 * time.Millisecond)
			if time.Now().Unix()-startTime > 10 {
				return errors.New("Time out (to) : It took more than 10 seconds to retrieve the correct receipt ")
			}
		} else {
			break
		}
	}

	return nil
}

func Run() {
	cli := cliPool.Alloc().(*client.Client)

	from := accGrp[rand.Int()%nAcc]
	to := accGrp[rand.Int()%nAcc]

	value := big.NewInt(int64(rand.Int() % 3))
	start := boomer.Now()

	err := TransferAndCheck(cli, to, from, value)

	elapsed := boomer.Now() - start

	if err == nil {
		boomer.Events.Publish("request_success", "http", "signedtransfer_with_check"+" to "+endPoint, elapsed, int64(10))
		cliPool.Free(cli)
	} else {
		boomer.Events.Publish("request_failure", "http", "signedtransfer_with_check"+" to "+endPoint, elapsed, err.Error())
	}
}
