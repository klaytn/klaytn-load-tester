package scKLAYTransferTc

import (
	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn/common"
	"github.com/myzhan/boomer"
	"math/big"
	"math/rand"
)

const Name = "scKLAYTransferTc"

var (
	nMcAcc   int
	mcAccGrp []*account.Account

	nScAcc   int
	scAccGrp []*account.Account

	mcBridgeAddr common.Address
	scBridgeAddr common.Address
)

func Init(mcAccs []*account.Account, scAccs []*account.Account, mcBridgeAddress, scBridgeAddress, _, _, _, _ common.Address) {
	mcAccGrp = append(mcAccGrp, mcAccs...)
	nMcAcc = len(mcAccGrp)
	mcBridgeAddr = mcBridgeAddress

	scAccGrp = append(scAccGrp, scAccs...)
	nScAcc = len(scAccGrp)
	scBridgeAddr = scBridgeAddress
}

func Run() {
	var err error

	var from, to *account.Account
	toServiceChain := false
	var targetBridgeAddress common.Address

	if (rand.Int() % 2) == 0 {
		from = mcAccGrp[rand.Int()%nMcAcc]
		to = scAccGrp[rand.Int()%nScAcc]
		toServiceChain = true

		targetBridgeAddress = mcBridgeAddr
	} else {
		from = scAccGrp[rand.Int()%nScAcc]
		to = mcAccGrp[rand.Int()%nMcAcc]

		targetBridgeAddress = scBridgeAddr
	}

	account.LockAccounts(from, to)
	defer account.UnlockAccounts(from, to)

	value := big.NewInt(int64(rand.Int()%2 + 1))
	start := boomer.Now()

	_, err = account.RequestKlayTransferReturnTx(from, to, value, targetBridgeAddress, false)

	elapsed := boomer.Now() - start

	var comment string
	if toServiceChain {
		comment = "Main -> Service chain."
	} else {
		comment = "Service -> Main chain."
	}

	if err == nil {
		//boomer.Events.Publish("request_success", "http", "signedtransfer"+" to "+endPoint, elapsed, int64(10))
		boomer.Events.Publish("request_success", "http", "Transfer KLAY "+comment, elapsed, int64(10))

	} else {
		//boomer.Events.Publish("request_failure", "http", "signedtransfer"+" to "+endPoint, elapsed, err.Error())
		boomer.Events.Publish("request_failure", "http", "Transfer KLAY "+comment, elapsed, err.Error())
	}
}
