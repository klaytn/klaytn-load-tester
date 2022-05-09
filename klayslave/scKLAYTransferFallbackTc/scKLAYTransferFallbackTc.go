package scKLAYTransferFallbackTc

import (
	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn/common"
	"github.com/myzhan/boomer"
	"math/big"
	"math/rand"
)

const Name = "scKLAYTransferFallbackTc"

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
	var bridgeAddr common.Address

	pAcc := mcAccGrp[rand.Int()%nMcAcc]
	cAcc := pAcc.CounterPartAccount()

	if (rand.Int() % 2) == 0 {
		from = pAcc
		to = cAcc
		toServiceChain = true

		bridgeAddr = mcBridgeAddr
	} else {
		from = cAcc
		to = pAcc

		bridgeAddr = scBridgeAddr
	}

	account.LockAccounts(pAcc, cAcc)
	defer account.UnlockAccounts(pAcc, cAcc)

	value := big.NewInt(int64(rand.Int()%2 + 1))
	start := boomer.Now()

	_, err = account.RequestKlayTransferFallbackReturnTx(from, to, value, bridgeAddr, false)

	elapsed := boomer.Now() - start

	var comment string
	if toServiceChain {
		comment = "Main -> Service chain."
	} else {
		comment = "Service -> Main chain."
	}

	if err == nil {
		//boomer.Events.Publish("request_success", "http", "signedtransfer"+" to "+endPoint, elapsed, int64(10))
		boomer.Events.Publish("request_success", "http", "Transfer KLAY fallback "+comment, elapsed, int64(10))

	} else {
		//boomer.Events.Publish("request_failure", "http", "signedtransfer"+" to "+endPoint, elapsed, err.Error())
		boomer.Events.Publish("request_failure", "http", "Transfer KLAY fallback "+comment, elapsed, err.Error())
	}
}
