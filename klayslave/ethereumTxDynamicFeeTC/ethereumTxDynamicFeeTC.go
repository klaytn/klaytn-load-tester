package ethereumTxDynamicFeeTC

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/clipool"
	"github.com/klaytn/klaytn-load-tester/klayslave/ethereumTxLegacyTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/task"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/myzhan/boomer"
)

const Name = "ethereumTxDynamicFeeTC"

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

	executablePath string

	maxRetryCount int

	SmartContractAccount *account.Account
	code                 string
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

	// Path to executable file that generates ethereum tx.
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println("exPath: ", exPath)

	executablePath = exPath + "/ethTxGenerator"
	log.Println("executablePath: ", executablePath)

	// Retry to get transaction receipt in checkResult
	maxRetryCount = 30

	code = "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"
}

func Run() {
	cli := cliPool.Alloc().(*client.Client)

	from := accGrp[rand.Int()%nAcc]
	to, value, input, reqType, err := ethereumTxLegacyTC.CreateRandomArguments(from.GetAddress())
	if err != nil {
		fmt.Printf("Failed to creat arguments to send Dynamic Fee Tx: %v\n", err.Error())
		return
	}

	start := boomer.Now()

	txHash, _, err := from.TransferNewEthDynamicFeeTxWithEth(cli, endPoint, to, value, input, executablePath)

	elapsed := boomer.Now() - start

	if err != nil {
		boomer.Events.Publish("request_failure", "http", "TransferNewEthDynamicFeeTx"+" to "+endPoint, elapsed, err.Error())
	}

	cliPool.Free(cli)

	// Check test result with CheckResult function
	go func(transactionHash common.Hash) {
		ret, err := CheckResult(transactionHash, reqType)
		if ret == false || err != nil {
			boomer.Events.Publish("request_failure", "http", "TransferNewEthDynamicFeeTx"+" to "+endPoint, elapsed, err.Error())
			return
		}

		boomer.Events.Publish("request_success", "http", "TransferNewEthDynamicFeeTx"+" to "+endPoint, elapsed, int64(10))
	}(txHash)
}

// CheckResult returns true and nil error, if expected results are observed, otherwise returns false and error.
func CheckResult(txHash common.Hash, reqType int) (bool, error) {
	cli := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(cli)

	receipt := ethereumTxLegacyTC.GetReceipt(cli, txHash, maxRetryCount)

	if receipt == nil {
		return false, errors.New("failed to get transaction receipt")
	}

	status, _ := receipt["status"].(string)
	if status != hexutil.Uint(types.ReceiptStatusSuccessful).String() {
		fmt.Printf("[FAILED] TxHash=%v, Receipt status=%v, Tx error msg=%v\n", txHash.String(), status, receipt["txError"])
		return false, errors.New("transaction status in receipt is fail")
	}

	// Check smart contract related fields
	if reqType == 1 {
		contractAddress, _ := receipt["contractAddress"]
		_, ok := contractAddress.(string)
		if !ok {
			return false, errors.New("failed to get contract address from the receipt")
		}
	} else if reqType == 2 {
		toFromReceipt, ok := receipt["to"].(string)
		if !ok || strings.ToLower(toFromReceipt) != strings.ToLower(SmartContractAccount.GetAddress().String()) {
			fmt.Printf("[FAILED] sca=%v, to=%v\n", SmartContractAccount.GetAddress().String(), toFromReceipt)
			return false, errors.New("mismatched to address in the receipt and smart contract address")
		}
	}

	typeString, ok := receipt["type"].(string)
	if !ok {
		return false, errors.New("failed to get type from the receipt")
	}
	if typeString != types.TxTypeEthereumDynamicFee.String() {
		return false, errors.New("type mismatched in transaction receipt")
	}

	typeInt, ok := receipt["typeInt"].(float64)
	if !ok {
		return false, errors.New("failed to get typeInt from the receipt")
	}
	if types.TxType(typeInt) != types.TxTypeEthereumDynamicFee {
		return false, errors.New("typeInt mismatched in transaction receipt")
	}

	return true, nil
}
