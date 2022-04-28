package internalTxTC

import (
	"context"
	"log"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/klaytn/klaytn-load-tester/klayslave/account"
	"github.com/klaytn/klaytn-load-tester/klayslave/clipool"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/myzhan/boomer"
)

const (
	Name        = "internalTxTC"
	NameMintNFT = "mintNFTTC"
)

var (
	endPoint string

	nAcc   int
	accGrp []*account.Account

	gasPrice *big.Int

	cliPool clipool.ClientPool

	mutex       sync.Mutex
	initialized = false

	KIP17Address  common.Address
	KIP17Contract *bind.BoundContract
	mainContract  *bind.BoundContract
)

// Init initializes cliPool and accGrp; and also deploys the smart contracts.
func Init(accs []*account.Account, endpoint string, gp *big.Int) {
	mutex.Lock()
	defer mutex.Unlock()

	if !initialized {
		initialized = true

		endPoint = endpoint
		gasPrice = gp

		cliCreate := func() interface{} {
			c, err := client.Dial(endPoint)
			if err != nil {
				log.Fatalf("Failed to connect to %s, err=%v", endPoint, err)
			}
			return c
		}
		cliPool.Init(20, 300, cliCreate)

		for _, acc := range accs {
			accGrp = append(accGrp, acc)
		}
		nAcc = len(accGrp)

		deployContracts(accs[0])
		log.Printf("TC initialized")
	}
}

func deployContracts(coinbase *account.Account) {
	conn, ok := cliPool.Alloc().(*client.Client)
	defer cliPool.Free(conn)
	if !ok {
		log.Fatal("[internalTxTC] conn is not client.Client")
		return
	}

	auth := bind.NewKeyedTransactor(coinbase.GetKey())
	auth.GasLimit = 9999999
	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))
	log.Println("[internalTxTC] Deploying two smart contracts")

	ctx := context.Background()
	defer ctx.Done()

	// Deploy Token Contract
	var KIP17Tx *types.Transaction
	for {
		var err error
		KIP17Address, KIP17Tx, KIP17Contract, err = DeployKIP17TokenContract(auth, conn)
		if err == nil {
			coinbase.UpdateNonce()
			break
		}

		log.Printf("[internalTxTC] Failed to deploy the KIP17 token mainContract, err=%v\n", err)
		auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))
		time.Sleep(1 * time.Second)
	}
	log.Printf("[internalTxTC] KIP17 token contract address: 0x%x\n", KIP17Address)
	log.Printf("[internalTxTC] Transaction waiting to be mined: 0x%x\n", KIP17Tx.Hash())

	for {
		time.Sleep(500 * time.Millisecond)
		receipt, err := conn.TransactionReceipt(ctx, KIP17Tx.Hash())
		if err != nil {
			log.Printf("[internalTxTC] Failed to check receipt: %v\n", err)
			continue
		}
		log.Println("[internalTxTC] Received the receipt")
		if receipt.Status == types.ReceiptStatusSuccessful {
			log.Println("[internalTxTC] Contract deployment was successful")
			break
		} else {
			log.Fatalf("[internalTxTC] Contract Receipt Status: %v\n", receipt.Status)
		}
	}

	// Deploy Main Contract
	var mainAddress common.Address
	var mainTx *types.Transaction
	auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))
	for {
		var err error
		mainAddress, mainTx, mainContract, err = DeployMainContract(auth, conn, KIP17Address)
		if err == nil {
			coinbase.UpdateNonce()
			break
		}

		log.Printf("[internalTxTC] Failed to deploy the main contract, err=%v\n", err)
		auth.Nonce = big.NewInt(int64(coinbase.GetNonceFromBlock(conn)))
		time.Sleep(1 * time.Second)
	}
	log.Printf("[internalTxTC] Main contract address : 0x%x\n", mainAddress)
	log.Printf("[internalTxTC] Transaction waiting to be mined: 0x%x\n", mainTx.Hash())

	for {
		time.Sleep(500 * time.Millisecond)
		receipt, err := conn.TransactionReceipt(ctx, mainTx.Hash())
		if err != nil {
			log.Printf("[internalTxTC] Failed to check receipt: %v\n", err)
			continue
		}
		log.Println("[internalTxTC] Received the receipt")
		if receipt.Status == types.ReceiptStatusSuccessful {
			log.Println("[internalTxTC] Contract deployment was successful")
			break
		} else {
			log.Fatalf("[internalTxTC] Contract Receipt Status: %v\n", receipt.Status)
		}
	}
}

// Run transfers txs calling sendRewards function of mainContract.
// During the execution of the function, four internal transactions are triggered also:
// Mint KIP17 token, Transfer KIP17 token, send KLAY to inviteeAccount and hostAccount.
func Run() {
	conn := cliPool.Alloc().(*client.Client)

	fromAccount := accGrp[rand.Int()%nAcc]
	fromAccount.Lock()
	nonce := fromAccount.GetNonce(conn)

	auth := bind.NewKeyedTransactor(fromAccount.GetKey())
	auth.GasLimit = 9999999
	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(100)

	inviteeAccount := accGrp[rand.Int()%nAcc]
	hostAccount := accGrp[rand.Int()%nAcc]

	start := boomer.Now()
	tx, err := mainContract.Transact(auth, "sendRewards", inviteeAccount.GetAddress(), hostAccount.GetAddress())
	if err != nil {
		log.Printf("[internalTxTC] Failed to execute contract, from=%s nonce=%d err=%v\n",
			fromAccount.GetAddress().String(), nonce, err)
		fromAccount.GetNonceFromBlock(conn)
	} else {
		log.Printf("[internalTxTC] tx=%s\n", tx.Hash().String())
		fromAccount.UpdateNonce()
	}
	elapsed := boomer.Now() - start
	fromAccount.UnLock()

	msg := "internalTxTC/" + " to " + endPoint
	if err == nil {
		boomer.Events.Publish("request_success", "contract", msg, elapsed, int64(10))
		cliPool.Free(conn)
	} else {
		boomer.Events.Publish("request_failure", "contract", msg, elapsed, err.Error())
		conn.Close()
	}
}

// RunMintNFT transfers txs calling mintCard function of KIP17Contract.
// The function mints a KIP17 token, NFT, for the sender.
func RunMintNFT() {
	conn := cliPool.Alloc().(*client.Client)

	fromAccount := accGrp[rand.Int()%nAcc]
	fromAccount.Lock()
	nonce := fromAccount.GetNonce(conn)

	auth := bind.NewKeyedTransactor(fromAccount.GetKey())
	auth.GasLimit = 9999999
	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(nonce))

	start := boomer.Now()
	tx, err := KIP17Contract.Transact(auth, "mintCard")
	if err != nil {
		log.Printf("[mintNFTTC] Failed to execute contract, from=%s nonce=%d err=%v\n",
			fromAccount.GetAddress().String(), nonce, err)
		fromAccount.GetNonceFromBlock(conn)
	} else {
		log.Printf("[mintNFTTC] tx=%s\n", tx.Hash().String())
		fromAccount.UpdateNonce()
	}
	elapsed := boomer.Now() - start
	fromAccount.UnLock()

	msg := "mintNFTTC/" + " to " + endPoint
	if err == nil {
		boomer.Events.Publish("request_success", "contract", msg, elapsed, int64(10))
		cliPool.Free(conn)
	} else {
		boomer.Events.Publish("request_failure", "contract", msg, elapsed, err.Error())
		conn.Close()
	}
}
