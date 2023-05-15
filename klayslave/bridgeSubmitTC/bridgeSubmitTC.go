package bridgeSubmitTC

import (
	"context"
	"log"
	"math/big"
	"math/rand"
	"sync"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/clipool"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/myzhan/boomer"
)

const Name = "bridgeSubmitTC"

var (
	endPointA   string
	nAccA       int
	accGrpA     []*account.Account
	accMu       sync.Mutex
	accApproved map[common.Address]bool
	cliPoolA    clipool.ClientPool
	gasPrice    *big.Int

	endPointB string
	nAccB     int
	accGrpB   []*account.Account
	cliPoolB  clipool.ClientPool

	// multinode tester
	transferedValue *big.Int
	expectedFee     *big.Int

	fromAccount     *account.Account
	prevBalanceFrom *big.Int

	toAccount     *account.Account
	prevBalanceTo *big.Int

	SmartContractAccount         *account.Account
	GovParamContractAccount      *account.Account
	BridgeContractAccount        *account.Account
	BridgeLibraryContractAccount *account.Account
	GovStateContractAccount      *account.Account
	VerifierContractAccount      *account.Account

	Erc20DeployAccount *account.Account

	blockMu      sync.Mutex
	highestBlock int64
	StartBlock   int64
)

func Init(accs []*account.Account, endpoint string, gp *big.Int) {

	highestBlock = 0

	gasPrice = gp

	endPointA = endpoint

	cliCreateA := func() interface{} {
		c, err := client.Dial(endPointA)
		if err != nil {
			log.Fatalf("Failed to connect RPC: %v", err)
		}
		return c
	}

	cliPoolA.Init(20, 300, cliCreateA)

	accApproved = make(map[common.Address]bool)
	for _, acc := range accs {
		accGrpA = append(accGrpA, acc)
		accApproved[acc.GetAddress()] = false
	}

	nAccA = len(accGrpA)
}

func Run() {
	cli := cliPoolA.Alloc().(*client.Client)

	from := accGrpA[rand.Int()%nAccA]

	start := boomer.Now()

	ctx := context.Background()
	blockNum, err := cli.BlockNumber(ctx)
	if err != nil {
		elapsed := boomer.Now() - start
		boomer.Events.Publish("request_failure", "http", "transferSubmitReceiptTx"+" to "+endPointA, elapsed, err.Error())
		return
	}
	blockMu.Lock()
	defer blockMu.Unlock()
	if blockNum.Int64() <= highestBlock {
		elapsed := boomer.Now() - start
		boomer.Events.Publish("request_success", "http", "transferSubmitReceiptTx"+" to "+endPointA, elapsed, int64(10))
		cliPoolA.Free(cli)
		return
	}
	blockHash, txHash, err := from.TransferSubmitHeader(cli, VerifierContractAccount.GetAddress(), blockNum.Int64())
	//bind.WaitMined(ctx, cli, txHash)
	if err != nil {
		elapsed := boomer.Now() - start
		boomer.Events.Publish("request_failure", "http", "transferSubmitReceiptTx"+" to "+endPointA, elapsed, err.Error())
		return
	}
	highestBlock = blockNum.Int64()

	_, err = from.TransferSubmitReceipt(cli, VerifierContractAccount.GetAddress(), blockNum.Int64(), blockHash, txHash, endPointA)
	elapsed := boomer.Now() - start
	if err == nil {
		boomer.Events.Publish("request_success", "http", "transferSubmitReceiptTx"+" to "+endPointA, elapsed, int64(10))
		cliPoolA.Free(cli)
	} else {
		boomer.Events.Publish("request_failure", "http", "transferSubmitReceiptTx"+" to "+endPointA, elapsed, err.Error())
	}

	return
}
