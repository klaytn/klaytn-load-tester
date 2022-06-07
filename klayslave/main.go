package main

import (
	"context"
	"flag"
	"fmt"
	"runtime"
	"syscall"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/scKLAYTransferFallbackTc"
	"github.com/klaytn/klaytn-load-tester/klayslave/scKLAYTransferTc"
	"github.com/klaytn/klaytn-load-tester/klayslave/scKLAYTransferTcWithCheck"
	"github.com/klaytn/klaytn-load-tester/klayslave/scNFTTransfer2StepTcWithCheck"
	"github.com/klaytn/klaytn-load-tester/klayslave/scNFTTransferTcWithCheck"
	"github.com/klaytn/klaytn-load-tester/klayslave/scTokenTransfer2StepTcWithCheck"
	"github.com/klaytn/klaytn-load-tester/klayslave/scTokenTransferTc"
	"github.com/klaytn/klaytn-load-tester/klayslave/scTokenTransferTcWithCheck"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/sc_erc20"
	"github.com/klaytn/klaytn/contracts/sc_erc721"
	"github.com/myzhan/boomer"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

var (
	mcCoinbasePrivatekey = ""
	mcBackend            *account.Backend
	mcEndpoint           string
	mcIP                 string

	scCoinbasePrivatekey = ""
	scBackend            *account.Backend
	scEndpoint           string

	subBridgeBackends []*account.Backend
	subBridges        []string

	mcCoinbase    *account.Account
	mcNewCoinbase *account.Account

	scCoinbase    *account.Account
	scNewCoinbase *account.Account

	nUser             = 100
	accGrpMc          []*account.Account
	mcTokenNFTAddress *common.Address

	operatorThreshold uint8

	nUserSc  = 100
	accGrpSc []*account.Account

	tcStr     string
	tcStrList []string

	KlayForAcc  *big.Int // KLAY balance for each test account
	TokenForAcc *big.Int // Token balance for each test account
)

type ExtendedTask struct {
	Weight int
	Fn     func()
	Name   string

	// For Initializing the task
	Init     func(mcAccs []*account.Account, scAccs []*account.Account, mcBridgeAddress, scBridgeAddress, mcTokenAddress, scTokenAddress, mcNFTAddress, scNFTAddress common.Address)
	AccGrpMc []*account.Account
	AccGrpSc []*account.Account
}

func chargeTestAccounts(coinBase *account.Account, accGrp map[common.Address]*account.Account) {
	println("Test Account Charge Start...")

	numChargedAcc := 0
	lastFailedNum := 0
	for _, acc := range accGrp {
		numChargedAcc++
		for {
			_, _, err := coinBase.TransferSignedTx(acc, KlayForAcc)
			if err == nil {
				acc.AddBalance(KlayForAcc)
				break // Success, move to next account.
			}

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
		}
		//bar.Increment()
	}
	//bar.Finish()	//bar.FinishPrint("Completed.")
}

func chargeTokenTestAccounts(coinBase *account.Account, accGrp map[common.Address]*account.Account, tokenAddr common.Address) {
	println("Test Account Token Charge Start...")

	numChargedAcc := 0
	lastFailedNum := 0

	var lastTx *types.Transaction

	for _, acc := range accGrp {
		numChargedAcc++
		for {
			tx, err := coinBase.TransferTokenToAccount(tokenAddr, acc, TokenForAcc)
			if err == nil {
				lastTx = tx
				break // Success, move to next account.
			}

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
		}
		//bar.Increment()
	}

	bind.WaitMined(context.Background(), coinBase.Backend(), lastTx)
}

func chargeNFTTestAccounts(coinBase *account.Account, accGrp map[common.Address]*account.Account, tokenAddr common.Address, startIDX uint64, numNFTPerAccount uint64) {
	println("Test Account NFT Charge Start...")

	numChargedAcc := 0
	lastFailedNum := 0

	var lastTx *types.Transaction

	for _, acc := range accGrp {
		numChargedAcc++
		for {
			start := startIDX + uint64(numChargedAcc-1)*numNFTPerAccount
			end := startIDX + uint64(numChargedAcc)*numNFTPerAccount
			tx, err := coinBase.RegisterNFTToAccount(tokenAddr, acc, start, end)
			if err == nil {
				lastTx = tx
				break // Success, move to next account.
			}

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
		}
		//bar.Increment()
	}

	bind.WaitMined(context.Background(), coinBase.Backend(), lastTx)
}

func prepareAccounts() {
	totalKLAY := new(big.Int)
	totalKLAY.Mul(KlayForAcc, big.NewInt(int64(nUser+nUserSc)))
	totalKLAY.Mul(totalKLAY, big.NewInt(3))

	// Import coinbase Account
	mcCoinbase = account.GetAccountFromKey(mcCoinbasePrivatekey, mcBackend)
	mcNewCoinbase = account.NewAccount(mcBackend)

	for {
		mcCoinbase.GetNonceFromBlock()
		tx, _, err := mcCoinbase.TransferSignedTx(mcNewCoinbase, totalKLAY)
		if err != nil {
			log.Printf("%v: charge newCoinbase fail: %v\n", os.Getpid(), err)
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		log.Printf("%v : charge newCoinbase: %v, Txhash=%v\n", os.Getpid(), mcNewCoinbase.GetAddress().String(), tx.Hash().String())

		getReceipt := false
		// After this loop waiting for 10 sec, It will retry to charge with new nonce.
		// it means another node stole the nonce.
		for i := 0; i < 5; i++ {
			time.Sleep(2000 * time.Millisecond)

			val, err := mcNewCoinbase.GetBalance()
			if err == nil {
				if val.Cmp(big.NewInt(0)) == 1 {
					getReceipt = true
					log.Printf("%v : charge newCoinbase success: %v, balance=%v peb\n", os.Getpid(), mcNewCoinbase.GetAddress().String(), val.String())
					break
				}
				log.Printf("%v : charge newCoinbase waiting: %v\n", os.Getpid(), mcNewCoinbase.GetAddress().String())
			} else {
				log.Printf("%v : check banalce err: %v\n", os.Getpid(), err)
			}
		}

		if getReceipt {
			break
		}
	}

	scCoinbase = account.GetAccountFromKey(scCoinbasePrivatekey, scBackend)
	scNewCoinbase = account.NewAccount(scBackend)

	for {
		scCoinbase.GetNonceFromBlock()
		tx, _, err := scCoinbase.TransferSignedTx(scNewCoinbase, totalKLAY)
		if err != nil {
			log.Printf("%v: charge newCoinbase fail: %v\n", os.Getpid(), err)
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		log.Printf("%v : charge newCoinbase: %v, Txhash=%v\n", os.Getpid(), scNewCoinbase.GetAddress().String(), tx.Hash().String())

		getReceipt := false
		// After this loop waiting for 10 sec, It will retry to charge with new nonce.
		// it means another node stole the nonce.
		for i := 0; i < 5; i++ {
			time.Sleep(2000 * time.Millisecond)

			val, err := scNewCoinbase.GetBalance()
			if err == nil {
				if val.Cmp(big.NewInt(0)) == 1 {
					getReceipt = true
					log.Printf("%v : charge newCoinbase success: %v, balance=%v peb\n", os.Getpid(), scNewCoinbase.GetAddress().String(), val.String())
					break
				}
				log.Printf("%v : charge newCoinbase waiting: %v\n", os.Getpid(), scNewCoinbase.GetAddress().String())
			} else {
				log.Printf("%v : check banalce err: %v\n", os.Getpid(), err)
			}
		}

		if getReceipt {
			break
		}
	}

	println("Parent/Child chain Account Group Preparation...")
	for i := 0; i < nUser; i++ {
		cAcc, pAcc := account.NewPairAccount(scBackend, mcBackend)
		accGrpMc = append(accGrpMc, pAcc)
		accGrpSc = append(accGrpSc, cAcc)
		fmt.Printf("%v\n", accGrpMc[i].GetAddress().String())
	}
}

func initArgs(tcNames string) {
	mcEndpointPtr := flag.String("mcEndpoint", "http://localhost:8545", "Target Main chain EndPoint")
	mcIPPtr := flag.String("mcIP", "127.0.0.1", "Target Main chain EndPoint's IP")
	scEndpointPtr := flag.String("scEndpoint", "http://localhost:7545", "Target Service chain EndPoint")

	nUserPtr := flag.Int("numUser", nUser, "num of test accounts")
	subBridgesPtr := flag.String("subbridges", "http://localhost:", "sub-bridge node EndPoint")

	mcKeyPtr := flag.String("mcKey", "", "privatekey of main chain coinbase")
	scKeyPtr := flag.String("scKey", "", "privatekey of service chain coinbase")

	operatorThresholdPtr := flag.Uint("threshold", 1, "operator threshold of bridge contracts")

	flag.StringVar(&tcStr, "tc", tcNames, "tasks which user want to run, multiple tasks are separated by comma.")

	flag.Parse()

	if *mcKeyPtr == "" {
		log.Fatal("mcKey argument is not defined. You should set the key for coinbase.\n example) klayslave -mcKeyPtr='2ef07640fd8d3f568c23185799ee92e0154bf08ccfe5c509466d1d40baca3430'")
	}

	if *scKeyPtr == "" {
		log.Fatal("scKey argument is not defined. You should set the key for coinbase.\n example) klayslave -scKeyPtr='2ef07640fd8d3f568c23185799ee92e0154bf08ccfe5c509466d1d40baca3430'")
	}

	// for TC Selection
	if tcStr != "" {
		// Run tasks without connecting to the master.
		tcStrList = strings.Split(tcStr, ",")
	}

	mcEndpoint = *mcEndpointPtr
	mcIP = *mcIPPtr
	scEndpoint = *scEndpointPtr

	nUser = *nUserPtr
	subBridges = strings.Split(*subBridgesPtr, ",")

	mcCoinbasePrivatekey = *mcKeyPtr
	scCoinbasePrivatekey = *scKeyPtr

	operatorThreshold = uint8(*operatorThresholdPtr)

	fmt.Println("Arguments are set like the following:")
	fmt.Printf("- Target mcEndPoint = %v\n", mcEndpoint)
	fmt.Printf("- Target scEndPoint = %v\n", scEndpoint)
	fmt.Printf("- nUser Main chain = %v\n", nUser)
	fmt.Printf("- subBridges = %v\n", *subBridgesPtr)
	fmt.Printf("- nUser Service chain = %v\n", nUserSc)
	fmt.Printf("- Main chain Key = %v\n", mcCoinbasePrivatekey)
	fmt.Printf("- Service chain Key = %v\n", scCoinbasePrivatekey)
	fmt.Printf("- tc = %v\n", tcStr)
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
		Name:     "scKLAYTransferTc",
		Weight:   10,
		Fn:       scKLAYTransferTc.Run,
		Init:     scKLAYTransferTc.Init,
		AccGrpMc: accGrpMc,
		AccGrpSc: accGrpSc,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:     "scKLAYTransferFallbackTc",
		Weight:   10,
		Fn:       scKLAYTransferFallbackTc.Run,
		Init:     scKLAYTransferFallbackTc.Init,
		AccGrpMc: accGrpMc,
		AccGrpSc: accGrpSc,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:     "scTokenTransferTc",
		Weight:   10,
		Fn:       scTokenTransferTc.Run,
		Init:     scTokenTransferTc.Init,
		AccGrpMc: accGrpMc,
		AccGrpSc: accGrpSc,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:     "scKLAYTransferTcWithCheck",
		Weight:   10,
		Fn:       scKLAYTransferTcWithCheck.Run,
		Init:     scKLAYTransferTcWithCheck.Init,
		AccGrpMc: accGrpMc,
		AccGrpSc: accGrpSc,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:     "scTokenTransferTcWithCheck",
		Weight:   10,
		Fn:       scTokenTransferTcWithCheck.Run,
		Init:     scTokenTransferTcWithCheck.Init,
		AccGrpMc: accGrpMc,
		AccGrpSc: accGrpSc,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:     "scTokenTransfer2StepTcWithCheck",
		Weight:   10,
		Fn:       scTokenTransfer2StepTcWithCheck.Run,
		Init:     scTokenTransfer2StepTcWithCheck.Init,
		AccGrpMc: accGrpMc,
		AccGrpSc: accGrpSc,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:     "scNFTTransferTcWithCheck",
		Weight:   10,
		Fn:       scNFTTransferTcWithCheck.Run,
		Init:     scNFTTransferTcWithCheck.Init,
		AccGrpMc: accGrpMc,
		AccGrpSc: accGrpSc,
	})

	taskSet = append(taskSet, &ExtendedTask{
		Name:     "scNFTTransfer2StepTcWithCheck",
		Weight:   10,
		Fn:       scNFTTransfer2StepTcWithCheck.Run,
		Init:     scNFTTransfer2StepTcWithCheck.Init,
		AccGrpMc: accGrpMc,
		AccGrpSc: accGrpSc,
	})

	return taskSet
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

	// Create main/service chain client
	mcBackend = account.NewBackend(mcEndpoint)
	scBackend = account.NewBackend(scEndpoint)
	for _, subBridge := range subBridges {
		subBridgeBackends = append(subBridgeBackends, account.NewBackend(subBridge))
	}

	// Set coinbase & Create Test Account
	KlayForAcc = new(big.Int)
	KlayForAcc.SetString("1000000000000000000000000", 10) //1e24 = 1,000,000 KLAY

	TokenForAcc = new(big.Int)
	TokenForAcc.SetString("1000000000000", 10)

	prepareAccounts()

	taskSet = initTCList()

	var filteredTask []*ExtendedTask

	println("Adding & Initializing tasks")
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

	if len(filteredTask) == 0 {
		log.Fatal("There is no valid TC.")
	}

	///////////////////////////
	// Temporal code /////////
	// TODO-Klaytn the following procedure will be omitted.

	ctx := context.Background()
	// deploy bridge
	scBridgeAddr, mcBridgeAddr, err := subBridgeBackends[0].Client.BridgeDeployBridge(ctx)
	if err != nil {
		log.Fatal("Failed to deploy bridge contract", err)
	}

	err = subBridgeBackends[0].BridgeSubscribeBridge(ctx, scBridgeAddr, mcBridgeAddr)
	if err != nil {
		log.Fatal("Failed to SubscribeEventBridge", "err", err)
	}

	_, err = subBridgeBackends[0].BridgeSetValueTransferOperatorThreshold(ctx, scBridgeAddr, operatorThreshold)
	if err != nil {
		log.Fatal("Failed to BridgeSetValueTransferOperatorThreshold", "err", err)
	}

	_, err = subBridgeBackends[0].BridgeSetValueTransferOperatorThreshold(ctx, mcBridgeAddr, operatorThreshold)
	if err != nil {
		log.Fatal("Failed to BridgeSetValueTransferOperatorThreshold", "err", err)
	}

	for i := 1; i < len(subBridgeBackends); i++ {
		err := subBridgeBackends[i].Client.BridgeRegisterBridge(ctx, scBridgeAddr, mcBridgeAddr)
		if err != nil {
			log.Fatal("Failed to BridgeRegisterBridge", "err", err)
		}

		err = subBridgeBackends[i].Client.BridgeSubscribeBridge(ctx, scBridgeAddr, mcBridgeAddr)
		if err != nil {
			log.Fatal("Failed to BridgeSubscribeBridge", "err", err)
		}

		pOperator, err := subBridgeBackends[i].Client.BridgeGetParentOperatorAddr(ctx)
		if err != nil {
			log.Fatal("Failed to BridgeGetParentOperatorAddr", "err", err)
		}

		cOperator, err := subBridgeBackends[i].Client.BridgeGetChildOperatorAddr(ctx)
		if err != nil {
			log.Fatal("Failed to BridgeGetChildOperatorAddr", "err", err)
		}

		// Register other sub-bridge's operators
		_, err = subBridgeBackends[0].BridgeRegisterOperator(ctx, mcBridgeAddr, pOperator)
		if err != nil {
			log.Fatal("Failed to BridgeRegisterOperator", "err", err)
		}

		_, err = subBridgeBackends[0].BridgeRegisterOperator(ctx, scBridgeAddr, cOperator)
		if err != nil {
			log.Fatal("Failed to BridgeRegisterOperator", "err", err)
		}
	}

	// deploy token contract
	mcTokenAddr, tx, mcToken, err := sctoken.DeployServiceChainToken(mcNewCoinbase.GetTransactOpts(account.MagicGasLimit), mcBackend, mcBridgeAddr)
	if err != nil {
		log.Fatal("Failed to DeployServiceChainToken on service chain", "err", err)
	}
	_, err = bind.WaitDeployed(ctx, mcBackend, tx)
	if err != nil {
		log.Fatal("Failed to WaitDeployed the token contract on service chain", "err", err)
	}
	mcNewCoinbase.UpdateNonce()

	scTokenAddr, tx, scToken, err := sctoken.DeployServiceChainToken(scNewCoinbase.GetTransactOpts(account.MagicGasLimit), scBackend, scBridgeAddr)
	if err != nil {
		log.Fatal("Failed to DeployServiceChainToken on service chain", "err", err)
	}
	_, err = bind.WaitDeployed(ctx, scBackend, tx)
	if err != nil {
		log.Fatal("Failed to WaitDeployed the token contract on service chain", "err", err)
	}
	scNewCoinbase.UpdateNonce()

	tx, err = scToken.AddMinter(scNewCoinbase.GetTransactOpts(account.MagicGasLimit), scBridgeAddr)
	r, err := bind.WaitMined(ctx, scBackend, tx)
	if r.Status != types.ReceiptStatusSuccessful {
		log.Fatal("Failed to WaitMined for AddMinter. Reverted Receipt")
	}
	if err != nil {
		log.Fatal("Failed to WaitMined for AddMinter", "err", err)
	}
	scNewCoinbase.UpdateNonce()

	// Register the pair of token on service chain bridge.
	err = subBridgeBackends[0].BridgeRegisterTokenContract(ctx, scBridgeAddr, mcBridgeAddr, scTokenAddr, mcTokenAddr)
	if err != nil {
		log.Fatal("Failed to BridgeRegisterTokenContract", "err", err)
	}

	// deploy NFT contract
	mcNFTAddr, tx, mcNFT, err := scnft.DeployServiceChainNFT(mcNewCoinbase.GetTransactOpts(account.MagicGasLimit), mcBackend, mcBridgeAddr)
	if err != nil {
		log.Fatal("Failed to DeployServiceChainNFT on service chain", "err", err)
	}

	_, err = bind.WaitDeployed(ctx, mcBackend, tx)
	if err != nil {
		log.Fatal("Failed to WaitDeployed the NFT contract on service chain", "err", err)
	}
	mcNewCoinbase.UpdateNonce()

	scNFTAddr, tx, scNFT, err := scnft.DeployServiceChainNFT(scNewCoinbase.GetTransactOpts(account.MagicGasLimit), scBackend, scBridgeAddr)
	if err != nil {
		log.Fatal("Failed to DeployServiceChainNFT on service chain", "err", err)
	}
	_, err = bind.WaitDeployed(ctx, scBackend, tx)
	if err != nil {
		log.Fatal("Failed to WaitDeployed the NFT contract on service chain", "err", err)
	}
	scNewCoinbase.UpdateNonce()

	tx, err = scNFT.AddMinter(scNewCoinbase.GetTransactOpts(account.MagicGasLimit), scBridgeAddr)
	r, err = bind.WaitMined(ctx, scBackend, tx)
	if r.Status != types.ReceiptStatusSuccessful {
		log.Fatal("Failed to WaitMined for AddMinter. Reverted Receipt")
	}
	if err != nil {
		log.Fatal("Failed to WaitDeployed for AddMinter", "err", err)
	}
	scNewCoinbase.UpdateNonce()

	// Register the pair of NFT on service chain bridge.
	err = subBridgeBackends[0].BridgeRegisterTokenContract(ctx, scBridgeAddr, mcBridgeAddr, scNFTAddr, mcNFTAddr)
	if err != nil {
		log.Fatal("Failed to BridgeRegisterContract", "err", err)
	}

	// Charge KLAY to bridge
	println("Charge KLAY to bridge")
	klayForBridge := new(big.Int)
	klayForBridge.Mul(KlayForAcc, big.NewInt(int64(nUser+nUserSc+1)))

	mcNewCoinbase.ChargeBridge(mcBridgeAddr, klayForBridge)
	scNewCoinbase.ChargeBridge(scBridgeAddr, klayForBridge)

	// Charge Token to bridge
	println("Charge Token to bridge")
	tokenForBridge := new(big.Int).Mul(TokenForAcc, big.NewInt(int64(nUser+nUserSc+1)))
	mcToken.Transfer(mcNewCoinbase.GetTransactOpts(account.MagicGasLimit), mcBridgeAddr, tokenForBridge)
	mcNewCoinbase.UpdateNonce()
	scToken.Transfer(scNewCoinbase.GetTransactOpts(account.MagicGasLimit), scBridgeAddr, tokenForBridge)
	scNewCoinbase.UpdateNonce()

	// Charge NFT to bridge
	println("Charge NFT to bridge")
	numNFTPerAccount := uint64(100)
	mcNFTStartIDX := big.NewInt(0)
	mcNFTEndIDX := new(big.Int).SetUint64(uint64(nUser) * numNFTPerAccount)
	scNFTStartIDX := mcNFTEndIDX
	scNFTEndIDX := new(big.Int).SetUint64(uint64(nUser+nUserSc) * numNFTPerAccount)

	chunk := numNFTPerAccount

	var txs []*types.Transaction
	for start := scNFTStartIDX.Uint64(); start+chunk <= scNFTEndIDX.Uint64(); start += chunk {
		st := new(big.Int).SetUint64(start)
		en := new(big.Int).SetUint64(start + chunk)
		tx, err := mcNFT.RegisterBulk(mcNewCoinbase.GetTransactOpts(account.MagicGasLimit), mcBridgeAddr, st, en)
		if err != nil {
			log.Fatal("Failed to RegisterBulk", "err", err)
		}
		txs = append(txs, tx)
		log.Printf("Charge NFT to mcbridge hash(%v), start(%v), end(%v)\n", tx.Hash().String(), st.String(), en.String())
		mcNewCoinbase.UpdateNonce()
	}

	for _, tx := range txs {
		receipt, err := bind.WaitMined(context.Background(), mcBackend, tx)
		if err != nil || receipt.Status != types.ReceiptStatusSuccessful {
			log.Fatal("Failed to RegisterBulk", "err", err, "txHash", tx.Hash().String())
		}
	}

	// This code is not necessary for mint-burn mode.
	//for start := mcNFTStartIDX.Uint64(); start+chunk <= mcNFTEndIDX.Uint64(); start += chunk {
	//	st := new(big.Int).SetUint64(start)
	//	en := new(big.Int).SetUint64(start + chunk)
	//	tx, err := scNFT.RegisterBulk(scNewCoinbase.GetTransactOpts(account.MagicGasLimit), scBridgeAddr, st, en)
	//	if err != nil {
	//		log.Fatal("Failed to RegisterBulk", "err", err)
	//	}
	//	log.Printf("Charge NFT to scbridge hash(%v), start(%v), end(%v)\n", tx.Hash().String(), st.String(), en.String())
	//	scNewCoinbase.UpdateNonce()
	//}

	// Charge Accounts
	mcAccGrp := make(map[common.Address]*account.Account)
	scAccGrp := make(map[common.Address]*account.Account)
	for _, task := range filteredTask {
		for _, acc := range task.AccGrpMc {
			_, exist := mcAccGrp[acc.GetAddress()]
			if !exist {
				mcAccGrp[acc.GetAddress()] = acc
			}
		}

		for _, acc := range task.AccGrpSc {
			_, exist := scAccGrp[acc.GetAddress()]
			if !exist {
				scAccGrp[acc.GetAddress()] = acc
			}
		}

	}

	chargeTestAccounts(mcNewCoinbase, mcAccGrp)
	chargeTestAccounts(scNewCoinbase, scAccGrp)

	chargeTokenTestAccounts(mcNewCoinbase, mcAccGrp, mcTokenAddr)
	chargeTokenTestAccounts(scNewCoinbase, scAccGrp, scTokenAddr)

	chargeNFTTestAccounts(mcNewCoinbase, mcAccGrp, mcNFTAddr, mcNFTStartIDX.Uint64(), numNFTPerAccount)
	chargeNFTTestAccounts(scNewCoinbase, scAccGrp, scNFTAddr, scNFTStartIDX.Uint64(), numNFTPerAccount)

	if len(filteredTask) == 0 {
		log.Fatal("No Tc is set. Please set TcList. \nExample argument) -tc='" + tcNames + "'")
	}

	println("Initializing tasks")
	var filteredBoomerTask []*boomer.Task
	for _, task := range filteredTask {
		task.Init(task.AccGrpMc, task.AccGrpSc, mcBridgeAddr, scBridgeAddr, mcTokenAddr, scTokenAddr, mcNFTAddr, scNFTAddr)
		filteredBoomerTask = append(filteredBoomerTask, &boomer.Task{task.Weight, task.Fn, task.Name})
		println("=> " + task.Name + " task is initialized.")
	}

	setRLimit(syscall.RLIMIT_NOFILE, 1024*400)

	// Locust Slave Run
	boomer.Run(filteredBoomerTask...)
	//boomer.Run(cpuHeavyTx)
}
