package main

//go:generate abigen --sol cpuHeavyTC/CPUHeavy.sol --pkg cpuHeavyTC --out cpuHeavyTC/CPUHeavy.go
//go:generate abigen --sol userStorageTC/UserStorage.sol --pkg userStorageTC --out userStorageTC/UserStorage.go

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/blockbench/analyticTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/blockbench/doNothingTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/blockbench/ioHeavyTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/blockbench/smallBankTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/blockbench/ycsbTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/bridgeRequestTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/bridgeSubmitTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/cpuHeavyTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/erc20TransferTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/erc721TransferTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/ethereumTxAccessListTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/ethereumTxDynamicFeeTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/ethereumTxLegacyTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/internalTxTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/largeMemoTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newAccountCreationTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newAccountUpdateTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newCancelTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newEthereumAccessListTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newEthereumDynamicFeeTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newFeeDelegatedAccountUpdateTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newFeeDelegatedAccountUpdateWithRatioTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newFeeDelegatedCancelTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newFeeDelegatedCancelWithRatioTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newFeeDelegatedSmartContractDeployTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newFeeDelegatedSmartContractDeployWithRatioTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newFeeDelegatedSmartContractExecutionTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newFeeDelegatedSmartContractExecutionWithRatioTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newFeeDelegatedValueTransferMemoTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newFeeDelegatedValueTransferMemoWithRatioTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newFeeDelegatedValueTransferTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newFeeDelegatedValueTransferWithRatioTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newSmartContractDeployTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newSmartContractExecutionTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newValueTransferLargeMemoTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newValueTransferMemoTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newValueTransferSmallMemoTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newValueTransferTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/newValueTransferWithCancelTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/readApiCallContractTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/readApiCallTC"
	receiptCheckTc "github.com/klaytn/klaytn-load-tester/klayslave/receiptCheckTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/storageTrieWriteTC"
	"github.com/klaytn/klaytn-load-tester/klayslave/transferSignedTc"
	"github.com/klaytn/klaytn-load-tester/klayslave/transferSignedWithCheckTc"
	"github.com/klaytn/klaytn-load-tester/klayslave/transferUnsignedTc"
	"github.com/klaytn/klaytn-load-tester/klayslave/userStorageTC"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	klay "github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/myzhan/boomer"
)

// sets build options from ldflags.
var (
	Version   = "1.0.0"
	Commit    string
	Branch    string
	Tag       string
	BuildDate string
	BuildUser string
)

var (
	coinbasePrivatekey = ""
	gCli               *klay.Client
	gEndpoint          string

	coinbase    *account.Account
	newCoinbase *account.Account

	nUserForUnsigned    = 5 //number of virtual user account for unsigned tx
	accGrpForUnsignedTx []*account.Account

	nUserForSigned    = 5
	accGrpForSignedTx []*account.Account

	nUserForNewAccounts  = 5
	accGrpForNewAccounts []*account.Account

	activeUserPercent = 100

	SmartContractAccount *account.Account
	//GovParamContractAddress common.Address

	tcStr     string
	tcStrList []string

	bridgeInfoFilePath string

	chargeValue *big.Int

	gasPrice *big.Int
	baseFee  *big.Int
)

type ExtendedTask struct {
	Name    string
	Weight  int
	Fn      func()
	Init    func(accs []*account.Account, endpoint string, gp *big.Int)
	AccGrp  []*account.Account
	EndPint string
}

func Create(endpoint string) *klay.Client {
	c, err := klay.Dial(endpoint)
	if err != nil {
		log.Fatalf("Failed to connect RPC: %v", err)
	}
	return c
}

func inTheTCList(tcName string) bool {
	for _, tc := range tcStrList {
		if tcName == tc {
			return true
		}
	}
	return false
}

// Dedicated and fixed private key used to deploy a smart contract for ERC20 and ERC721 value transfer performance test.
var ERC20DeployPrivateKeyStr = "eb2c84d41c639178ff26a81f488c196584d678bb1390cc20a3aeb536f3969a98"
var BridgeDeployPrivateKeyStr = "ed233e94fa04d1487080b773efbba4654dccc9371ac7225b790b641a42426bbf"
var ERC721DeployPrivateKeyStr = "45c40d95c9b7898a21e073b5bf952bcb05f2e70072e239a8bbd87bb74a53355e"

type BridgeInfo struct {
	VerifierContractName  string
	GovStateContractName  string
	BridgeLibContractName string
	BridgeLibAddr         string
	VerifierAddr          string
	GovStateAddr          string
	TokenContractName     string
	TokenAddr             string
}

type ChainInfo struct {
	Http          string
	SenderPrivKey string
	FeeReceiver   string
}

type SettingJson struct {
	BridgeFee     string
	HdContribFee  string
	GovContribFee string
	ChainA        ChainInfo
	ChainB        ChainInfo
}

// prepareBridgeRequestTransfer sets up Bridge request transfer performance test.
func prepareBridgeRequestTransfer(accGrp map[common.Address]*account.Account) {
	if !inTheTCList(bridgeRequestTC.Name) {
		return
	}

	// parse privkey from settings.json
	d, err := os.Open(bridgeInfoFilePath + "/test/integration_test/settings.json")
	if err != nil {
		fmt.Println(err)
	}

	readall, _ := ioutil.ReadAll(d)

	var setting SettingJson

	json.Unmarshal(readall, &setting)
	BridgeDeployPrivateKeyStr = setting.ChainA.SenderPrivKey

	erc20DeployAcc := account.GetAccountFromKey(0, BridgeDeployPrivateKeyStr)
	account.HardhatPath = bridgeInfoFilePath

	bridgeRequestTC.Erc20DeployAccount = erc20DeployAcc
	log.Printf("prepareBridgeTransfer", "addr", erc20DeployAcc.GetAddress().String())
	chargeKLAYToTestAccounts(map[common.Address]*account.Account{erc20DeployAcc.GetAddress(): erc20DeployAcc})

	// A ERC20 contract for bridge transfer performance TC.
	bridgeRequestTC.SmartContractAccount = deploySingleSmartContract(erc20DeployAcc, erc20DeployAcc.DeployERC20, "Bridge Performance Test Contract")
	newCoinBaseAccountMap := map[common.Address]*account.Account{newCoinbase.GetAddress(): newCoinbase}
	firstChargeTokenToTestAccounts(newCoinBaseAccountMap, bridgeRequestTC.SmartContractAccount.GetAddress(), erc20DeployAcc.TransferERC20, big.NewInt(1e11))

	chargeTokenToTestAccounts(accGrp, bridgeRequestTC.SmartContractAccount.GetAddress(), newCoinbase.TransferERC20, big.NewInt(1e10))

	// GovParm contract
	//bridgeRequestTC.GovParamContractAccount = deploySingleSmartContract(erc20DeployAcc, erc20DeployAcc.DeployGovParam, "Bridge Performance Test GovParam Contract")

	// Parse already deployed verifier contract addr
	data, err := os.Open(bridgeInfoFilePath + "/chainA_deployedAddrs.json")
	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(data)
	fmt.Println(byteValue)

	var bridgeInfo BridgeInfo

	json.Unmarshal(byteValue, &bridgeInfo)

	bridgeLibAddr, _ := hex.DecodeString(string(bridgeInfo.BridgeLibAddr[2:]))
	govStateAddr, _ := hex.DecodeString(string(bridgeInfo.GovStateAddr[2:]))
	verifierAddr, _ := hex.DecodeString(string(bridgeInfo.VerifierAddr[2:]))
	account.BridgeLibraryContractAddr = common.BytesToAddress(bridgeLibAddr)
	account.GovStateContractAddr = common.BytesToAddress(govStateAddr)
	account.VerifierContractAddr = common.BytesToAddress(verifierAddr)

	bridgeRequestTC.BridgeLibraryContractAccount = account.NewKlaytnAccountWithAddr(1, account.BridgeLibraryContractAddr)
	bridgeRequestTC.GovStateContractAccount = account.NewKlaytnAccountWithAddr(1, account.GovStateContractAddr)
	bridgeRequestTC.VerifierContractAccount = account.NewKlaytnAccountWithAddr(1, account.VerifierContractAddr)

	deploySingleSmartContract(erc20DeployAcc, erc20DeployAcc.InitVerifier, "Bridge Performance Test Verifier Contract")

	// Get list of committee
	ctx := context.Background()
	blockNum, err := gCli.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("[PrepareBridge] Cannot get blocknum")
		return
	}
	header, err := gCli.HeaderByNumber(ctx, blockNum)
	if err != nil {
		log.Fatalf("[PrepareBridge] Cannot get header")
		return
	}
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		log.Fatalf("[PrepareBridge] Cannot get istanbul extra")
		return
	}

	bridgeRequestTC.OriginalCommittee = make([]common.Address, len(istanbulExtra.Validators))
	bridgeRequestTC.CurrentCommittee = make([]common.Address, len(istanbulExtra.Validators))

	copy(bridgeRequestTC.OriginalCommittee, istanbulExtra.Validators)
	copy(bridgeRequestTC.CurrentCommittee, istanbulExtra.Validators)

	bridgeRequestTC.IsAdding = false
}

// prepareBridgeSubmitTransfer sets up Bridge submit transfer performance test.
func prepareBridgeSubmitTransfer(accGrp map[common.Address]*account.Account) {
	if !inTheTCList(bridgeSubmitTC.Name) {
		return
	}

	// Get current blockNum for initial search
	ctx := context.Background()
	blockNum, err := gCli.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("[PrepareBridgeSubmit] Cannot get blocknum")
		return
	}
	if !inTheTCList(bridgeRequestTC.Name) {
		data, err := os.Open(bridgeInfoFilePath + "/chainA_deployedAddrs.json")
		if err != nil {
			fmt.Println(err)
		}

		byteValue, _ := ioutil.ReadAll(data)
		fmt.Println(byteValue)

		var bridgeInfo BridgeInfo

		json.Unmarshal(byteValue, &bridgeInfo)
		bridgeLibAddr, _ := hex.DecodeString(string(bridgeInfo.BridgeLibAddr[2:]))
		govStateAddr, _ := hex.DecodeString(string(bridgeInfo.GovStateAddr[2:]))
		verifierAddr, _ := hex.DecodeString(string(bridgeInfo.VerifierAddr[2:]))
		account.BridgeLibraryContractAddr = common.BytesToAddress(bridgeLibAddr)
		account.GovStateContractAddr = common.BytesToAddress(govStateAddr)
		account.VerifierContractAddr = common.BytesToAddress(verifierAddr)
	}
	account.MaxTxCount = 1000
	bridgeSubmitTC.StartBlock = blockNum.Int64()
	bridgeSubmitTC.BridgeLibraryContractAccount = account.NewKlaytnAccountWithAddr(1, account.BridgeLibraryContractAddr)
	bridgeSubmitTC.GovStateContractAccount = account.NewKlaytnAccountWithAddr(1, account.GovStateContractAddr)
	bridgeSubmitTC.VerifierContractAccount = account.NewKlaytnAccountWithAddr(1, account.VerifierContractAddr)
}

// prepareERC20Transfer sets up ERC20 transfer performance test.
func prepareERC20Transfer(accGrp map[common.Address]*account.Account) {
	if !inTheTCList(erc20TransferTC.Name) {
		return
	}
	erc20DeployAcc := account.GetAccountFromKey(0, ERC20DeployPrivateKeyStr)
	log.Printf("prepareERC20Transfer", "addr", erc20DeployAcc.GetAddress().String())
	chargeKLAYToTestAccounts(map[common.Address]*account.Account{erc20DeployAcc.GetAddress(): erc20DeployAcc})

	// A smart contract for ERC20 value transfer performance TC.
	erc20TransferTC.SmartContractAccount = deploySingleSmartContract(erc20DeployAcc, erc20DeployAcc.DeployERC20, "ERC20 Performance Test Contract")
	newCoinBaseAccountMap := map[common.Address]*account.Account{newCoinbase.GetAddress(): newCoinbase}
	firstChargeTokenToTestAccounts(newCoinBaseAccountMap, erc20TransferTC.SmartContractAccount.GetAddress(), erc20DeployAcc.TransferERC20, big.NewInt(1e11))

	chargeTokenToTestAccounts(accGrp, erc20TransferTC.SmartContractAccount.GetAddress(), newCoinbase.TransferERC20, big.NewInt(1e4))
}

// prepareERC721Transfer sets up ERC721 transfer performance test.
func prepareERC721Transfer(accGrp []*account.Account) {
	if !inTheTCList(erc721TransferTC.Name) {
		return
	}
	erc721DeployAcc := account.GetAccountFromKey(0, ERC721DeployPrivateKeyStr)
	log.Printf("prepareERC721Transfer", "addr", erc721DeployAcc.GetAddress().String())
	chargeKLAYToTestAccounts(map[common.Address]*account.Account{erc721DeployAcc.GetAddress(): erc721DeployAcc})

	// A smart contract for ERC721 value transfer performance TC.
	erc721TransferTC.SmartContractAccount = deploySingleSmartContract(erc721DeployAcc, erc721DeployAcc.DeployERC721, "ERC721 Performance Test Contract")

	// Wait for reward tester to get started
	time.Sleep(30 * time.Second)
	newCoinbase.MintERC721ToTestAccounts(gCli, accGrp, erc721TransferTC.SmartContractAccount.GetAddress(), 5)
	log.Println("MintERC721ToTestAccounts", "len(accGrp)", len(accGrp))
}

// Dedicated and fixed private key used to deploy a smart contract for storage trie write performance test.
var storageTrieDeployPrivateKeyStr = "3737c381633deaaa4c0bdbc64728f6ef7d381b17e1d30bbb74665839cec942b8"

// prepareStorageTrieWritePerformance sets up ERC20 storage trie write performance test.
func prepareStorageTrieWritePerformance(accGrp map[common.Address]*account.Account) {
	if !inTheTCList(storageTrieWriteTC.Name) {
		return
	}
	storageTrieDeployAcc := account.GetAccountFromKey(0, storageTrieDeployPrivateKeyStr)
	log.Printf("prepareStorageTrieWritePerformance", "addr", storageTrieDeployAcc.GetAddress().String())
	chargeKLAYToTestAccounts(map[common.Address]*account.Account{storageTrieDeployAcc.GetAddress(): storageTrieDeployAcc})

	// A smart contract for storage trie store performance TC.
	storageTrieWriteTC.SmartContractAccount = deploySingleSmartContract(storageTrieDeployAcc, storageTrieDeployAcc.DeployStorageTrieWrite, "Storage Trie Performance Test Contract")
}

func prepareTestAccountsAndContracts(accGrp map[common.Address]*account.Account) {
	// First, charging KLAY to the test accounts.
	chargeKLAYToTestAccounts(accGrp)

	// Second, deploy contracts used for some TCs.
	// If the test case is not on the list, corresponding contract won't be deployed.
	prepareERC20Transfer(accGrp)
	prepareBridgeRequestTransfer(accGrp)
	prepareBridgeSubmitTransfer(accGrp)
	prepareStorageTrieWritePerformance(accGrp)

	// Third, deploy contracts for general tests.
	// A smart contract for general smart contract related TCs.
	GeneralSmartContract := deploySmartContract(newCoinbase.TransferNewSmartContractDeployTxHumanReadable, "General Purpose Test Smart Contract")
	newSmartContractExecutionTC.SmartContractAccount = GeneralSmartContract
	newFeeDelegatedSmartContractExecutionTC.SmartContractAccount = GeneralSmartContract
	newFeeDelegatedSmartContractExecutionWithRatioTC.SmartContractAccount = GeneralSmartContract
	ethereumTxLegacyTC.SmartContractAccount = GeneralSmartContract
	ethereumTxAccessListTC.SmartContractAccount = GeneralSmartContract
	ethereumTxDynamicFeeTC.SmartContractAccount = GeneralSmartContract
	newEthereumAccessListTC.SmartContractAccount = GeneralSmartContract
	newEthereumDynamicFeeTC.SmartContractAccount = GeneralSmartContract
}

func chargeKLAYToTestAccounts(accGrp map[common.Address]*account.Account) {
	log.Printf("Start charging KLAY to test accounts")

	numChargedAcc := 0
	lastFailedNum := 0
	for _, acc := range accGrp {
		for {
			_, _, err := newCoinbase.TransferSignedTxReturnTx(true, gCli, acc, chargeValue)
			if err == nil {
				break // Success, move to next account.
			}
			numChargedAcc, lastFailedNum = estimateRemainingTime(accGrp, numChargedAcc, lastFailedNum)
		}
		numChargedAcc++
	}

	log.Printf("Finished charging KLAY to %d test account(s), Total %d transactions are sent.\n", len(accGrp), numChargedAcc)
}

type tokenChargeFunc func(initialCharge bool, c *client.Client, tokenContractAddr common.Address, recipient *account.Account, value *big.Int) (*types.Transaction, *big.Int, error)

// firstChargeTokenToTestAccounts charges initially generated tokens to newCoinbase account for further testing.
// As this work is done simultaneously by different slaves, this should be done in "try and check" manner.
func firstChargeTokenToTestAccounts(accGrp map[common.Address]*account.Account, tokenContractAddr common.Address, tokenChargeFn tokenChargeFunc, tokenChargeAmount *big.Int) {
	log.Printf("Start initial token charging to new coinbase")

	numChargedAcc := 0
	for _, recipientAccount := range accGrp {
		for {
			tx, _, err := tokenChargeFn(true, gCli, tokenContractAddr, recipientAccount, tokenChargeAmount)
			for err != nil {
				log.Printf("Failed to execute %s: err %s", tx.Hash().String(), err.Error())
				time.Sleep(1 * time.Second) // Mostly, the err is `txpool is full`, retry after a while.
				tx, _, err = tokenChargeFn(true, gCli, tokenContractAddr, recipientAccount, tokenChargeAmount)
			}
			ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
			receipt, err := bind.WaitMined(ctx, gCli, tx)
			cancelFn()
			if receipt != nil {
				break
			}
		}
		numChargedAcc++
	}

	log.Printf("Finished initial token charging to %d new coinbase account(s), Total %d transactions are sent.\n", len(accGrp), numChargedAcc)
}

// chargeTokenToTestAccounts charges default token to the test accounts for testing.
// As it is done independently among the slaves, it has simpler logic than firstChargeTokenToTestAccounts.
func chargeTokenToTestAccounts(accGrp map[common.Address]*account.Account, tokenContractAddr common.Address, tokenChargeFn tokenChargeFunc, tokenChargeAmount *big.Int) {
	log.Printf("Start charging tokens to test accounts")

	numChargedAcc := 0
	lastFailedNum := 0
	for _, recipientAccount := range accGrp {
		for {
			_, _, err := tokenChargeFn(false, gCli, tokenContractAddr, recipientAccount, tokenChargeAmount)
			if err == nil {
				break // Success, move to next account.
			}
			numChargedAcc, lastFailedNum = estimateRemainingTime(accGrp, numChargedAcc, lastFailedNum)
		}
		numChargedAcc++
	}

	log.Printf("Finished charging tokens to %d test account(s), Total %d transactions are sent.\n", len(accGrp), numChargedAcc)
}

func estimateRemainingTime(accGrp map[common.Address]*account.Account, numChargedAcc, lastFailedNum int) (int, int) {
	if lastFailedNum > 0 {
		// Not 1st failed cases.
		TPS := (numChargedAcc - lastFailedNum) / 5 // TPS of only this slave during `txpool is full` situation.
		lastFailedNum = numChargedAcc

		if TPS <= 5 {
			log.Printf("Retry to charge test account #%d. But it is too slow. %d TPS\n", numChargedAcc, TPS)
		} else {
			remainTime := (len(accGrp) - numChargedAcc) / TPS
			remainHour := remainTime / 3600
			remainMinute := (remainTime % 3600) / 60

			log.Printf("Retry to charge test account #%d. Estimated remaining time: %d hours %d mins later\n", numChargedAcc, remainHour, remainMinute)
		}
	} else {
		// 1st failed case.
		lastFailedNum = numChargedAcc
		log.Printf("Retry to charge test account #%d.\n", numChargedAcc)
	}
	time.Sleep(5 * time.Second) // Mostly, the err is `txpool is full`, retry after a while.
	return numChargedAcc, lastFailedNum
}

type contractDeployFunc func(c *client.Client, to *account.Account, value *big.Int, humanReadable bool) (common.Address, *types.Transaction, *big.Int, error)

// deploySmartContract deploys smart contracts by the number of locust slaves.
// In other words, each slave owns its own contract for testing.
func deploySmartContract(contractDeployFn contractDeployFunc, contractName string) *account.Account {
	addr, lastTx, _, err := contractDeployFn(gCli, SmartContractAccount, common.Big0, false)
	for err != nil {
		log.Printf("Failed to deploy a %s: err %s", contractName, err.Error())
		time.Sleep(5 * time.Second) // Mostly, the err is `txpool is full`, retry after a while.
		addr, lastTx, _, err = contractDeployFn(gCli, SmartContractAccount, common.Big0, false)
	}

	log.Printf("Start waiting the receipt of the %s tx(%v).\n", contractName, lastTx.Hash().String())
	bind.WaitMined(context.Background(), gCli, lastTx)

	deployedContract := account.NewKlaytnAccountWithAddr(1, addr)
	log.Printf("%s has been deployed to : %s\n", contractName, addr.String())
	return deployedContract
}

// deploySingleSmartContract deploys only one smart contract among the slaves.
// It the contract is already deployed by other slave, it just calculates the address of the contract.
func deploySingleSmartContract(erc20DeployAcc *account.Account, contractDeployFn contractDeployFunc, contractName string) *account.Account {
	addr, lastTx, _, err := contractDeployFn(gCli, SmartContractAccount, common.Big0, false)
	for err != nil {
		if err == account.AlreadyDeployedErr {
			erc20Addr := crypto.CreateAddress(erc20DeployAcc.GetAddress(), 0)
			return account.NewKlaytnAccountWithAddr(1, erc20Addr)
		}
		if strings.HasPrefix(err.Error(), "known transaction") {
			erc20Addr := crypto.CreateAddress(erc20DeployAcc.GetAddress(), 0)
			return account.NewKlaytnAccountWithAddr(1, erc20Addr)
		}
		log.Printf("Failed to deploy a %s: err %s", contractName, err.Error())
		time.Sleep(5 * time.Second) // Mostly, the err is `txpool is full`, retry after a while.
		addr, lastTx, _, err = contractDeployFn(gCli, SmartContractAccount, common.Big0, false)
	}

	log.Printf("Start waiting the receipt of the %s tx(%v).\n", contractName, lastTx.Hash().String())
	bind.WaitMined(context.Background(), gCli, lastTx)

	deployedContract := account.NewKlaytnAccountWithAddr(1, addr)
	log.Printf("%s has been deployed to : %s\n", contractName, addr.String())
	return deployedContract
}

func prepareAccounts() {
	totalChargeValue := new(big.Int)
	totalChargeValue.Mul(chargeValue, big.NewInt(int64(nUserForUnsigned+nUserForSigned+nUserForNewAccounts+1)))

	// Import coinbase Account
	coinbase = account.GetAccountFromKey(0, coinbasePrivatekey)
	newCoinbase = account.NewAccount(0)

	if len(chargeValue.Bits()) != 0 {
		for {
			coinbase.GetNonceFromBlock(gCli)
			hash, _, err := coinbase.TransferSignedTx(gCli, newCoinbase, totalChargeValue)
			if err != nil {
				log.Printf("%v: charge newCoinbase fail: %v\n", os.Getpid(), err)
				time.Sleep(1000 * time.Millisecond)
				continue
			}

			log.Printf("%v : charge newCoinbase: %v, Txhash=%v\n", os.Getpid(), newCoinbase.GetAddress().String(), hash.String())

			getReceipt := false
			// After this loop waiting for 10 sec, It will retry to charge with new nonce.
			// it means another node stole the nonce.
			for i := 0; i < 5; i++ {
				time.Sleep(2000 * time.Millisecond)
				ctx := context.Background()

				//_, err := gCli.TransactionReceipt(ctx, hash)
				//if err != nil {
				//	getReceipt = true
				//	log.Printf("%v : charge newCoinbase success: %v\n", os.Getpid(), newCoinbase.GetAddress().String())
				//	break
				//}
				//log.Printf("%v : charge newCoinbase waiting: %v\n", os.Getpid(), newCoinbase.GetAddress().String())

				val, err := gCli.BalanceAt(ctx, newCoinbase.GetAddress(), nil)
				if err == nil {
					if val.Cmp(big.NewInt(0)) == 1 {
						getReceipt = true
						log.Printf("%v : charge newCoinbase success: %v, balance=%v peb\n", os.Getpid(), newCoinbase.GetAddress().String(), val.String())
						break
					}
					log.Printf("%v : charge newCoinbase waiting: %v\n", os.Getpid(), newCoinbase.GetAddress().String())
				} else {
					log.Printf("%v : check balance err: %v\n", os.Getpid(), err)
				}
			}

			if getReceipt {
				break
			}
		}
	}

	println("Unsigned Account Group Preparation...")
	//bar := pb.StartNew(nUserForUnsigned)

	// Create test account pool
	for i := 0; i < nUserForUnsigned; i++ {
		accGrpForUnsignedTx = append(accGrpForUnsignedTx, account.NewAccount(i))
		fmt.Printf("%v\n", accGrpForUnsignedTx[i].GetAddress().String())
		//bar.Increment()
	}
	//bar.Finish()	//bar.FinishPrint("Completed.")
	//
	println("Signed Account Group Preparation...")
	//bar = pb.StartNew(nUserForSigned)

	for i := 0; i < nUserForSigned; i++ {
		accGrpForSignedTx = append(accGrpForSignedTx, account.NewAccount(i))
		fmt.Printf("%v\n", accGrpForSignedTx[i].GetAddress().String())
		//bar.Increment()
	}

	println("New account group preparation...")
	for i := 0; i < nUserForNewAccounts; i++ {
		accGrpForNewAccounts = append(accGrpForNewAccounts, account.NewKlaytnAccount(i))
	}
}

func initArgs(tcNames string) {
	chargeKLAYAmount := 1000000000
	gEndpointPtr := flag.String("endpoint", "http://localhost:8545", "Target EndPoint")
	nUserForSignedPtr := flag.Int("vusigned", nUserForSigned, "num of test account for signed Tx TC")
	nUserForUnsignedPtr := flag.Int("vuunsigned", nUserForUnsigned, "num of test account for unsigned Tx TC")
	activeUserPercentPtr := flag.Int("activepercent", activeUserPercent, "percent of active accounts")
	keyPtr := flag.String("key", "", "privatekey of coinbase")
	chargeKLAYAmountPtr := flag.Int("charge", chargeKLAYAmount, "charging amount for each test account in KLAY")
	versionPtr := flag.Bool("version", false, "show version number")
	httpMaxIdleConnsPtr := flag.Int("http.maxidleconns", 100, "maximum number of idle connections in default http client")
	flag.StringVar(&tcStr, "tc", tcNames, "tasks which user want to run, multiple tasks are separated by comma.")
	bridgeInfoFilePathPtr := flag.String("bridgeinfo", "", "path to bridge transfer contract info file")

	flag.Parse()

	if *versionPtr || (len(os.Args) >= 2 && os.Args[1] == "version") {
		printVersion()
		os.Exit(0)
	}

	if *keyPtr == "" {
		log.Fatal("key argument is not defined. You should set the key for coinbase.\n example) klaytc -key='2ef07640fd8d3f568c23185799ee92e0154bf08ccfe5c509466d1d40baca3430'")
	}

	// setup default http client.
	if tr, ok := http.DefaultTransport.(*http.Transport); ok {
		tr.MaxIdleConns = *httpMaxIdleConnsPtr
		tr.MaxIdleConnsPerHost = *httpMaxIdleConnsPtr
	}

	// for TC Selection
	if tcStr != "" {
		// Run tasks without connecting to the master.
		tcStrList = strings.Split(tcStr, ",")
	}

	bridgeInfoFilePath = *bridgeInfoFilePathPtr

	gEndpoint = *gEndpointPtr

	nUserForSigned = *nUserForSignedPtr
	nUserForUnsigned = *nUserForUnsignedPtr
	activeUserPercent = *activeUserPercentPtr
	coinbasePrivatekey = *keyPtr
	chargeKLAYAmount = *chargeKLAYAmountPtr
	chargeValue = new(big.Int)
	chargeValue.Set(new(big.Int).Mul(big.NewInt(int64(chargeKLAYAmount)), big.NewInt(params.KLAY)))

	fmt.Println("Arguments are set like the following:")
	fmt.Printf("- Target EndPoint = %v\n", gEndpoint)
	fmt.Printf("- nUserForSigned = %v\n", nUserForSigned)
	fmt.Printf("- nUserForUnsigned = %v\n", nUserForUnsigned)
	fmt.Printf("- activeUserPercent = %v\n", activeUserPercent)
	fmt.Printf("- coinbasePrivatekey = %v\n", coinbasePrivatekey)
	fmt.Printf("- charging KLAY Amount = %v\n", chargeKLAYAmount)
	fmt.Printf("- tc = %v\n", tcStr)
}

func updateChainID() {
	fmt.Println("Updating ChainID from RPC")
	for {
		ctx := context.Background()
		chainID, err := gCli.ChainID(ctx)

		if err == nil {
			fmt.Println("chainID :", chainID)
			account.SetChainID(chainID)
			break
		}
		fmt.Println("Retrying updating chainID... ERR: ", err)

		time.Sleep(2 * time.Second)
	}
}

func updateGasPrice() {
	// TODO: refactor to updating gasPrice with goverance.magma.upperboundbasefee
	gasPrice = big.NewInt(750000000000)

	/* Deprecated because of KIP-71 hardfork
	fmt.Println("Updating GasPrice from RPC")
	for {
		ctx := context.Background()
		gp, err := gCli.SuggestGasPrice(ctx)

		if err == nil {
			gasPrice = gp
			fmt.Println("gas price :", gasPrice.String())
			break
		}
		fmt.Println("Retrying updating GasPrice... ERR: ", err)

		time.Sleep(2 * time.Second)
	}
	*/
	account.SetGasPrice(gasPrice)
}

func updateBaseFee() {
	baseFee = big.NewInt(0)
	// TODO: Uncomment below when klaytn 1.8.0 is released.
	//for {
	//	ctx := context.Background()
	//	h, err := gCli.HeaderByNumber(ctx, nil)
	//
	//	if err == nil {
	//		baseFee = h.BaseFee
	//		fmt.Println("base fee :", baseFee.String())
	//		break
	//	}
	//	fmt.Println("Retrying updating BaseFee... ERR: ", err)
	//
	//	time.Sleep(2 * time.Second)
	//}
	account.SetBaseFee(baseFee)
}

func setRLimit(resourceType int, val uint64) error {
	if runtime.GOOS == "darwin" {
		return nil
	}
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(resourceType, &rLimit)
	if err != nil {
		return err
	}
	rLimit.Cur = val
	err = syscall.Setrlimit(resourceType, &rLimit)
	if err != nil {
		return err
	}
	return nil
}

// initTCList initializes TCs and returns a slice of TCs.
func initTCList() (taskSet []*ExtendedTask) {

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "analyticTx",
		Weight:  10,
		Fn:      analyticTC.Run,
		Init:    analyticTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "analyticQueryLargestAccBalTx",
		Weight:  10,
		Fn:      analyticTC.QueryLargestAccBal,
		Init:    analyticTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "analyticQueryLargestTxValTx",
		Weight:  10,
		Fn:      analyticTC.QueryLargestTxVal,
		Init:    analyticTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "analyticQueryTotalTxValTx",
		Weight:  10,
		Fn:      analyticTC.QueryTotalTxVal,
		Init:    analyticTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "cpuHeavyTx",
		Weight:  10,
		Fn:      cpuHeavyTC.Run,
		Init:    cpuHeavyTC.Init,
		AccGrp:  accGrpForSignedTx, //[nUserForSigned/2:],
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "doNothingTx",
		Weight:  10,
		Fn:      doNothingTC.Run,
		Init:    doNothingTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    internalTxTC.Name,
		Weight:  10,
		Fn:      internalTxTC.Run,
		Init:    internalTxTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    internalTxTC.NameMintNFT,
		Weight:  10,
		Fn:      internalTxTC.RunMintNFT,
		Init:    internalTxTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "ioHeavyTx",
		Weight:  10,
		Fn:      ioHeavyTC.Run,
		Init:    ioHeavyTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "ioHeavyScanTx",
		Weight:  10,
		Fn:      ioHeavyTC.Scan,
		Init:    ioHeavyTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "ioHeavyWriteTx",
		Weight:  10,
		Fn:      ioHeavyTC.Write,
		Init:    ioHeavyTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "largeMemoTC",
		Weight:  10,
		Fn:      largeMemoTC.Run,
		Init:    largeMemoTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    receiptCheckTc.Name,
		Weight:  10,
		Fn:      receiptCheckTc.Run,
		Init:    receiptCheckTc.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "smallBankTx",
		Weight:  10,
		Fn:      smallBankTC.Run,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "smallBankAlmagateTx",
		Weight:  10,
		Fn:      smallBankTC.Almagate,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "smallBankGetBalanceTx",
		Weight:  10,
		Fn:      smallBankTC.GetBalance,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "smallBankSendPaymentTx",
		Weight:  10,
		Fn:      smallBankTC.SendPayment,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "smallBankUpdateBalanceTx",
		Weight:  10,
		Fn:      smallBankTC.UpdateBalance,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "smallBankUpdateSavingTx",
		Weight:  10,
		Fn:      smallBankTC.UpdateSaving,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "smallBankWriteCheckTx",
		Weight:  10,
		Fn:      smallBankTC.WriteCheck,
		Init:    smallBankTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "transferSignedTx",
		Weight:  10,
		Fn:      transferSignedTc.Run,
		Init:    transferSignedTc.Init,
		AccGrp:  accGrpForSignedTx, //[:nUserForSigned/2-1],
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newValueTransferTC",
		Weight:  10,
		Fn:      newValueTransferTC.Run,
		Init:    newValueTransferTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newValueTransferWithCancelTC",
		Weight:  10,
		Fn:      newValueTransferWithCancelTC.Run,
		Init:    newValueTransferWithCancelTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newFeeDelegatedValueTransferTC",
		Weight:  10,
		Fn:      newFeeDelegatedValueTransferTC.Run,
		Init:    newFeeDelegatedValueTransferTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newFeeDelegatedValueTransferWithRatioTC",
		Weight:  10,
		Fn:      newFeeDelegatedValueTransferWithRatioTC.Run,
		Init:    newFeeDelegatedValueTransferWithRatioTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newValueTransferMemoTC",
		Weight:  10,
		Fn:      newValueTransferMemoTC.Run,
		Init:    newValueTransferMemoTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newValueTransferLargeMemoTC",
		Weight:  10,
		Fn:      newValueTransferLargeMemoTC.Run,
		Init:    newValueTransferLargeMemoTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newValueTransferSmallMemoTC",
		Weight:  10,
		Fn:      newValueTransferSmallMemoTC.Run,
		Init:    newValueTransferSmallMemoTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newFeeDelegatedValueTransferMemoTC",
		Weight:  10,
		Fn:      newFeeDelegatedValueTransferMemoTC.Run,
		Init:    newFeeDelegatedValueTransferMemoTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newFeeDelegatedValueTransferMemoWithRatioTC",
		Weight:  10,
		Fn:      newFeeDelegatedValueTransferMemoWithRatioTC.Run,
		Init:    newFeeDelegatedValueTransferMemoWithRatioTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newAccountCreationTC",
		Weight:  10,
		Fn:      newAccountCreationTC.Run,
		Init:    newAccountCreationTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newAccountUpdateTC",
		Weight:  10,
		Fn:      newAccountUpdateTC.Run,
		Init:    newAccountUpdateTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newFeeDelegatedAccountUpdateTC",
		Weight:  10,
		Fn:      newFeeDelegatedAccountUpdateTC.Run,
		Init:    newFeeDelegatedAccountUpdateTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newFeeDelegatedAccountUpdateWithRatioTC",
		Weight:  10,
		Fn:      newFeeDelegatedAccountUpdateWithRatioTC.Run,
		Init:    newFeeDelegatedAccountUpdateWithRatioTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newSmartContractDeployTC",
		Weight:  10,
		Fn:      newSmartContractDeployTC.Run,
		Init:    newSmartContractDeployTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newFeeDelegatedSmartContractDeployTC",
		Weight:  10,
		Fn:      newFeeDelegatedSmartContractDeployTC.Run,
		Init:    newFeeDelegatedSmartContractDeployTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newFeeDelegatedSmartContractDeployWithRatioTC",
		Weight:  10,
		Fn:      newFeeDelegatedSmartContractDeployWithRatioTC.Run,
		Init:    newFeeDelegatedSmartContractDeployWithRatioTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newSmartContractExecutionTC",
		Weight:  10,
		Fn:      newSmartContractExecutionTC.Run,
		Init:    newSmartContractExecutionTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "bridgeRequestTC",
		Weight:  10,
		Fn:      bridgeRequestTC.Run,
		Init:    bridgeRequestTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "bridgeSubmitTC",
		Weight:  10,
		Fn:      bridgeSubmitTC.Run,
		Init:    bridgeSubmitTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    storageTrieWriteTC.Name,
		Weight:  10,
		Fn:      storageTrieWriteTC.Run,
		Init:    storageTrieWriteTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newFeeDelegatedSmartContractExecutionTC",
		Weight:  10,
		Fn:      newFeeDelegatedSmartContractExecutionTC.Run,
		Init:    newFeeDelegatedSmartContractExecutionTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newFeeDelegatedSmartContractExecutionWithRatioTC",
		Weight:  10,
		Fn:      newFeeDelegatedSmartContractExecutionWithRatioTC.Run,
		Init:    newFeeDelegatedSmartContractExecutionWithRatioTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newCancelTC",
		Weight:  10,
		Fn:      newCancelTC.Run,
		Init:    newCancelTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newFeeDelegatedCancelTC",
		Weight:  10,
		Fn:      newFeeDelegatedCancelTC.Run,
		Init:    newFeeDelegatedCancelTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newFeeDelegatedCancelWithRatioTC",
		Weight:  10,
		Fn:      newFeeDelegatedCancelWithRatioTC.Run,
		Init:    newFeeDelegatedCancelWithRatioTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "transferSignedWithCheckTx",
		Weight:  10,
		Fn:      transferSignedWithCheckTc.Run,
		Init:    transferSignedWithCheckTc.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "transferUnsignedTx",
		Weight:  10,
		Fn:      transferUnsignedTc.Run,
		Init:    transferUnsignedTc.Init,
		AccGrp:  accGrpForUnsignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "userStorageSetTx",
		Weight:  10,
		Fn:      userStorageTC.RunSet,
		Init:    userStorageTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "userStorageSetGetTx",
		Weight:  10,
		Fn:      userStorageTC.RunSetGet,
		Init:    userStorageTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "ycsbTx",
		Weight:  10,
		Fn:      ycsbTC.Run,
		Init:    ycsbTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "ycsbGetTx",
		Weight:  10,
		Fn:      ycsbTC.Get,
		Init:    ycsbTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "ycsbSetTx",
		Weight:  10,
		Fn:      ycsbTC.Set,
		Init:    ycsbTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    erc20TransferTC.Name,
		Weight:  10,
		Fn:      erc20TransferTC.Run,
		Init:    erc20TransferTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    erc721TransferTC.Name,
		Weight:  10,
		Fn:      erc721TransferTC.Run,
		Init:    erc721TransferTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "readGasPrice",
		Weight:  10,
		Fn:      readApiCallTC.GasPrice,
		Init:    readApiCallTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "readBlockNumber",
		Weight:  10,
		Fn:      readApiCallTC.BlockNumber,
		Init:    readApiCallTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "readGetBlockByNumber",
		Weight:  10,
		Fn:      readApiCallTC.GetBlockByNumber,
		Init:    readApiCallTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "readGetAccount",
		Weight:  10,
		Fn:      readApiCallTC.GetAccount,
		Init:    readApiCallTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "readGetBlockWithConsensusInfoByNumber",
		Weight:  10,
		Fn:      readApiCallTC.GetBlockWithConsensusInfoByNumber,
		Init:    readApiCallTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "readGetStorageAt",
		Weight:  10,
		Fn:      readApiCallContractTC.GetStorageAt,
		Init:    readApiCallContractTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "readCall",
		Weight:  10,
		Fn:      readApiCallContractTC.Call,
		Init:    readApiCallContractTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "readEstimateGas",
		Weight:  10,
		Fn:      readApiCallContractTC.EstimateGas,
		Init:    readApiCallContractTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "ethereumTxLegacyTC",
		Weight:  10,
		Fn:      ethereumTxLegacyTC.Run,
		Init:    ethereumTxLegacyTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "ethereumTxAccessListTC",
		Weight:  10,
		Fn:      ethereumTxAccessListTC.Run,
		Init:    ethereumTxAccessListTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "ethereumTxDynamicFeeTC",
		Weight:  10,
		Fn:      ethereumTxDynamicFeeTC.Run,
		Init:    ethereumTxDynamicFeeTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newEthereumAccessListTC",
		Weight:  10,
		Fn:      newEthereumAccessListTC.Run,
		Init:    newEthereumAccessListTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:    "newEthereumDynamicFeeTC",
		Weight:  10,
		Fn:      newEthereumDynamicFeeTC.Run,
		Init:    newEthereumDynamicFeeTC.Init,
		AccGrp:  accGrpForSignedTx,
		EndPint: gEndpoint,
	})

	return taskSet
}

func printVersion() {
	version := Version
	if len(Commit) >= 7 {
		version += "-" + Commit[:7]
	}
	if Tag != "" && Tag != "undefined" {
		version = Tag
	}
	fmt.Printf("Version :\t%s\n", version)
	fmt.Printf("git.Branch :\t%s\n", Branch)
	fmt.Printf("git.Commit :\t%s\n", Commit)
	fmt.Printf("git.Tag :\t%s\n", Tag)
	fmt.Printf("build.Date :\t%s\n", BuildDate)
	fmt.Printf("build.User :\t%s\n", BuildUser)
}

func main() {
	// Call initTCList to get all TC names
	taskSet := initTCList()

	var tcNames string
	for i, task := range taskSet {
		if i != 0 {
			tcNames += ","
		}
		tcNames += task.Name
	}

	initArgs(tcNames)

	// Create Cli pool
	gCli = Create(gEndpoint)

	// Update chainID
	updateChainID()

	// Update gasPrice
	updateGasPrice()

	gasPrice = big.NewInt(750000000000)

	// Update baseFee
	updateBaseFee()

	// Set coinbase & Create Test Account
	prepareAccounts()

	// Call initTCList again to actually define all TCs
	taskSet = initTCList()

	var filteredTask []*ExtendedTask

	println("Adding tasks")
	for _, task := range taskSet {
		if task.Name == "" {
			continue
		} else {
			flag := false
			for _, name := range tcStrList {
				if name == task.Name {
					flag = true
					break
				}
			}
			if flag {
				filteredTask = append(filteredTask, task)
				println("=> " + task.Name + " task is added.")
			}
		}
	}

	// Import/Unlock Account on the node if there is a task to use unsigned account group.
	for _, task := range filteredTask {
		if task.AccGrp[0] == accGrpForUnsignedTx[0] {
			for _, acc := range task.AccGrp {
				acc.ImportUnLockAccount(gEndpoint)
			}
			break // to import/unlock once.
		}
	}

	// Charge Accounts
	accGrp := make(map[common.Address]*account.Account)
	for _, task := range filteredTask {
		for _, acc := range task.AccGrp {
			_, exist := accGrp[acc.GetAddress()]
			if !exist {
				accGrp[acc.GetAddress()] = acc
			}
		}

	}

	if len(chargeValue.Bits()) != 0 {
		prepareTestAccountsAndContracts(accGrp)
	}

	// After charging accounts, cut the slice to the desired length, calculated by ActiveAccountPercent.
	for _, task := range filteredTask {
		if activeUserPercent > 100 {
			log.Fatalf("ActiveAccountPercent should be less than or equal to 100, but it is %v", activeUserPercent)
		}
		numActiveAccounts := len(task.AccGrp) * activeUserPercent / 100
		// Not to assign 0 account for some cases.
		if numActiveAccounts == 0 {
			numActiveAccounts = 1
		}
		task.AccGrp = task.AccGrp[:numActiveAccounts]
		prepareERC721Transfer(task.AccGrp)
	}

	if len(filteredTask) == 0 {
		log.Fatal("No Tc is set. Please set TcList. \nExample argument) -tc='" + tcNames + "'")
	}

	println("Initializing tasks")
	var filteredBoomerTask []*boomer.Task
	for _, task := range filteredTask {
		task.Init(task.AccGrp, task.EndPint, gasPrice)
		filteredBoomerTask = append(filteredBoomerTask, &boomer.Task{task.Weight, task.Fn, task.Name})
		println("=> " + task.Name + " task is initialized.")
	}

	setRLimit(syscall.RLIMIT_NOFILE, 1024*400)

	// Locust Slave Run
	boomer.Run(filteredBoomerTask...)
	//boomer.Run(cpuHeavyTx)
}
