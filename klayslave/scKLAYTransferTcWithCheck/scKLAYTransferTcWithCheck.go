package scKLAYTransferTcWithCheck

import (
	"context"
	"fmt"
	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/myzhan/boomer"
	"math/big"
	"math/rand"
	"time"
)

const Name = "scKLAYTransferTcWithCheck"

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
	var scBackend *account.Backend
	var hTxHash common.Hash
	var err error
	var rTx *types.Transaction

	var from, to *account.Account
	toServiceChain := false
	var targetBridgeAddress common.Address

	if (rand.Int() % 2) == 0 {
		from = mcAccGrp[rand.Int()%nMcAcc]
		to = scAccGrp[rand.Int()%nScAcc]
		toServiceChain = true
		scBackend = to.Backend()

		targetBridgeAddress = mcBridgeAddr
	} else {
		from = scAccGrp[rand.Int()%nScAcc]
		to = mcAccGrp[rand.Int()%nMcAcc]
		scBackend = from.Backend()

		targetBridgeAddress = scBridgeAddr
	}

	account.LockAccounts(from, to)
	defer account.UnlockAccounts(from, to)

	value := big.NewInt(int64(rand.Int()%2 + 1))
	start := boomer.Now()

	rTx, err = account.RequestKlayTransferReturnTx(from, to, value, targetBridgeAddress, false)
	if err == nil {
		// Wait for the corresponded handle transaction hash
		for i := 0; i < 1000; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			hTxHash, err = scBackend.BridgeConvertRequestTxHashToHandleTxHash(ctx, rTx.Hash())
			if err == nil && hTxHash != (common.Hash{}) {
				//fmt.Printf("Found HTX(%v) of RTX(%v)\n", hTxHash.String(), rTx.Hash().String())

				goto CHECK_HTX_RECEIPT
			}
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Printf("Not found HTX(%v) of RTX(%v)\n", hTxHash.String(), rTx.Hash().String())
		err = fmt.Errorf("can not find a HTX of the RTX(%v)", rTx.Hash().String)
		goto FINISH
	}
	fmt.Printf("Fail to RequestKlayTransferReturnTx RTX(%v)\n", rTx.Hash().String())
	goto FINISH

CHECK_HTX_RECEIPT:
	for i := 0; i <= 1000; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		receipt, _ := to.Backend().TransactionReceipt(ctx, hTxHash)
		if receipt != nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				err = fmt.Errorf("HTX Receipt status : %v", receipt.Status)
				fmt.Printf("HTX Receipt status HTX(%v) receipt.Status(%v)", hTxHash.String(), receipt.Status)
			}
			//fmt.Printf("Success KLAY transfer. RTX(%v), HTX(%v), elapsed = %v\n", rTx.Hash().String(), hTxHash.String(), boomer.Now()-start)
			goto FINISH
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("fail to KLAY transfer. RTX(%v), HTX(%v), elapsed = %v\n", rTx.Hash().String(), hTxHash.String(), boomer.Now()-start)
	err = fmt.Errorf("timeout to find HTX(%v)'s receipt", hTxHash.String())

FINISH:
	elapsed := boomer.Now() - start

	var comment string
	if toServiceChain {
		comment = "Main -> Service chain."
	} else {
		comment = "Service -> Main chain."
	}

	if err == nil {
		//boomer.Events.Publish("request_success", "http", "signedtransfer"+" to "+endPoint, elapsed, int64(10))
		boomer.Events.Publish("request_success", "http", "Transfer KLAY with check"+comment, elapsed, int64(10))

	} else {
		//boomer.Events.Publish("request_failure", "http", "signedtransfer"+" to "+endPoint, elapsed, err.Error())
		boomer.Events.Publish("request_failure", "http", "Transfer KLAY with check"+comment, elapsed, err.Error())
	}
}
