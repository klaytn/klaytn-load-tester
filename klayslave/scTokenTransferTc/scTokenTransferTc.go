package scTokenTransferTc

import (
	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn/common"
	"math/big"
	"math/rand"
)

const Name = "scTokenTransferTc"

var (
	nMcAcc   int
	mcAccGrp []*account.Account

	nScAcc   int
	scAccGrp []*account.Account

	mcTokenAddr common.Address
	scTokenAddr common.Address
)

func Init(mcAccs []*account.Account, scAccs []*account.Account, _, _, mcTokenAddress, scTokenAddress, _, _ common.Address) {
	mcAccGrp = append(mcAccGrp, mcAccs...)
	nMcAcc = len(mcAccGrp)
	mcTokenAddr = mcTokenAddress

	scAccGrp = append(scAccGrp, scAccs...)
	nScAcc = len(scAccGrp)
	scTokenAddr = scTokenAddress
}

func Run() {
	var from, to *account.Account
	toServiceChain := false
	var targetTokenAddress common.Address

	if (rand.Int() % 2) == 0 {
		from = mcAccGrp[rand.Int()%nMcAcc]
		to = scAccGrp[rand.Int()%nScAcc]
		toServiceChain = true

		targetTokenAddress = mcTokenAddr
	} else {
		from = scAccGrp[rand.Int()%nScAcc]
		to = mcAccGrp[rand.Int()%nMcAcc]

		targetTokenAddress = scTokenAddr
	}

	account.LockAccounts(from, to)
	defer account.UnlockAccounts(from, to)

	value := big.NewInt(int64(rand.Int()%2) + 1)
	start := boomer.Now()

	_, err := account.RequestTokenTransfer(from, to, value, targetTokenAddress)

	elapsed := boomer.Now() - start

	var comment string
	if toServiceChain {
		comment = "Main -> Service chain."
	} else {
		comment = "Service -> Main chain."
	}

	if err == nil {
		//boomer.Events.Publish("request_success", "http", "signedtransfer"+" to "+endPoint, elapsed, int64(10))
		boomer.Events.Publish("request_success", "http", "Transfer ERC20 "+comment, elapsed, int64(10))

	} else {
		//boomer.Events.Publish("request_failure", "http", "signedtransfer"+" to "+endPoint, elapsed, err.Error())
		boomer.Events.Publish("request_failure", "http", "Transfer ERC20 "+comment, elapsed, err.Error())
	}
}
