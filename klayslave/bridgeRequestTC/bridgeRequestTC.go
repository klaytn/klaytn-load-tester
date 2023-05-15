package bridgeRequestTC

import (
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

const Name = "bridgeRequestTC"

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

	OriginalCommittee  []common.Address
	CurrentCommittee   []common.Address
	CandidateCommittee []common.Address
	IsAdding           bool
)

func Init(accs []*account.Account, endpoint string, gp *big.Int) {
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
	to := Erc20DeployAccount
	value := big.NewInt(int64(10))

	start := boomer.Now()

	var err error
	accMu.Lock()
	val, ok := accApproved[from.GetAddress()]
	accMu.Unlock()
	if !ok || !val {
		_, _, err = from.Approve(cli, SmartContractAccount.GetAddress(), VerifierContractAccount.GetAddress(), to, value, 1001)
		elapsed := boomer.Now() - start
		if err != nil {
			log.Fatal(err)
			boomer.Events.Publish("request_failure", "http", "transferBridgeRequestTx"+" to "+endPointA, elapsed, err.Error())
			return
		}
		accMu.Lock()
		accApproved[from.GetAddress()] = true
		accMu.Unlock()
		boomer.Events.Publish("request_success", "http", "transferBridgeRequestTx"+" to "+endPointA, elapsed, int64(10))
		cliPoolA.Free(cli)
	} else {
		_, _, err := from.TransferBridgeErc20Request(cli, SmartContractAccount.GetAddress(), VerifierContractAccount.GetAddress(), to, value, 1001)
		elapsed := boomer.Now() - start
		if err != nil {
			boomer.Events.Publish("request_failure", "http", "transferBridgeRequestTx"+" to "+endPointA, elapsed, err.Error())
			return
		}
		boomer.Events.Publish("request_success", "http", "transferBridgeRequestTx"+" to "+endPointA, elapsed, int64(10))
		cliPoolA.Free(cli)
	}

}
