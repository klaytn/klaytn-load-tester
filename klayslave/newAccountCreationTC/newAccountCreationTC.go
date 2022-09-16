package newAccountCreationTC

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/clipool"
	"github.com/klaytn/klaytn-load-tester/klayslave/task"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
	"github.com/myzhan/boomer"
)

const Name = "newAccountCreationTC"

var (
	endPoint string
	nAcc     int
	accGrp   []*account.Account
	cliPool  clipool.ClientPool
	gasPrice *big.Int

	// multinode tester
	transferedValue *big.Int
	expectedFee     *big.Int

	fromAccount     *account.Account
	prevBalanceFrom *big.Int

	toAccount     *account.Account
	prevBalanceTo *big.Int
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
}

func Run() {
	cli := cliPool.Alloc().(*client.Client)

	from := accGrp[rand.Int()%nAcc]
	to := account.NewKlaytnAccount(0)
	value := big.NewInt(int64(rand.Int() % 3))

	start := boomer.Now()
	_, _, err := from.TransferNewAccountCreationTx(cli, to, value)
	elapsed := boomer.Now() - start

	if err == nil {
		boomer.Events.Publish("request_success", "http", "transferNewAccountCreationTx"+" to "+endPoint, elapsed, int64(10))
		cliPool.Free(cli)
	} else {
		boomer.Events.Publish("request_failure", "http", "transferNewAccountCreationTx"+" to "+endPoint, elapsed, err.Error())
	}
}

func RunSingle() (txHash common.Hash, err error) {
	cli := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(cli)

	fromIdx := rand.Int() % nAcc

	from := accGrp[fromIdx]
	to := account.NewKlaytnAccount(0)
	value := big.NewInt(int64(rand.Int() % 3))
	fmt.Printf("[TC] transferSignedTc: %v, from:%v, to:%v, value:%v\n", endPoint, from.GetAddress().String(), to.GetAddress().String(), value)
	transferedValue = big.NewInt(value.Int64())
	expectedFee = big.NewInt(0).Mul(big.NewInt(25*params.Ston), big.NewInt(21000))

	balance, err := from.GetBalance(cli)
	if err != nil {
		return common.Hash{}, err
	}
	fromAccount = from
	prevBalanceFrom = big.NewInt(balance.Int64())
	fmt.Printf("From:%v, balance:%v\n", fromAccount.GetAddress().String(), prevBalanceFrom.Int64())

	balance, err = to.GetBalance(cli)
	if err != nil {
		return common.Hash{}, err
	}
	toAccount = to
	prevBalanceTo = big.NewInt(balance.Int64())
	fmt.Printf("To:%v, balance:%v\n", toAccount.GetAddress().String(), prevBalanceTo.Int64())

	txHash, _, err = from.TransferNewAccountCreationTx(cli, to, value)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, err
}

// CheckResult returns true and nil error, if expected results are observed.
// Otherewise returns false and error.
func CheckResult() (bool, error) {
	cli := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(cli)

	balance, err := fromAccount.GetBalance(cli)
	if err != nil {
		return false, err
	}
	expectedBalance := big.NewInt(0)
	expectedBalance.Sub(prevBalanceFrom, transferedValue)
	expectedBalance.Sub(expectedBalance, expectedFee)
	// fmt.Printf("prevBalanceFrom=%v, transferedValue=%v, expectedFee=%v\n", prevBalanceFrom.Int64(), transferedValue.Int64(), expectedFee.Int64())

	if expectedBalance.Int64() != balance.Int64() {
		fmt.Printf("[FAILED] From account address=%v, Expected balance=%v, Actual balance=%v\n", fromAccount.GetAddress().String(), expectedBalance.Int64(), balance.Int64())
		return false, errors.New("Balance mismatched!")
	} else {
		fmt.Printf("[PASSED] From account address=%v, Expected balance=%v, Actual balance=%v\n", fromAccount.GetAddress().String(), expectedBalance.Int64(), balance.Int64())
	}

	balance, err = toAccount.GetBalance(cli)
	if err != nil {
		return false, err
	}
	expectedBalance = big.NewInt(0)
	expectedBalance.Add(prevBalanceTo, transferedValue)
	fmt.Printf("prevBalanceTo=%v, transferedValue=%v\n", prevBalanceTo.Int64(), transferedValue.Int64())
	if expectedBalance.Int64() != balance.Int64() {
		fmt.Printf("[FAILED] To account address=%v, Expected balance=%v, Actual balance=%v\n", toAccount.GetAddress().String(), expectedBalance.Int64(), balance.Int64())
		return false, errors.New("Balance mismatched!")
	} else {
		fmt.Printf("[PASSED] To account address=%v, Expected balance=%v, Actual balance=%v\n", toAccount.GetAddress().String(), expectedBalance.Int64(), balance.Int64())
	}

	return true, err
}
