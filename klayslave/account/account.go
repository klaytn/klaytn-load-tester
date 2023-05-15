package account

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
)

const Letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var (
	gasPrice       *big.Int
	chainID        *big.Int
	baseFee        *big.Int
	HardhatPath    string
	GlobalGovMutex sync.Mutex
	MaxTxCount     uint
)

type Proof struct {
	Proof [][]byte
	Key   []byte
}

type Account struct {
	id         int
	privateKey []*ecdsa.PrivateKey
	key        []string
	address    common.Address
	nonce      uint64
	balance    *big.Int
	mutex      sync.Mutex
}

func init() {
	gasPrice = big.NewInt(0)
	chainID = big.NewInt(2018)
	baseFee = big.NewInt(0)
}

func SetGasPrice(gp *big.Int) {
	gasPrice = gp
}

func SetBaseFee(bf *big.Int) {
	baseFee = bf
}

func SetChainID(id *big.Int) {
	chainID = id
}

func (acc *Account) Lock() {
	acc.mutex.Lock()
}

func (acc *Account) UnLock() {
	acc.mutex.Unlock()
}

func GetAccountFromKey(id int, key string) *Account {
	acc, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Fatalf("Key(%v): Failed to HexToECDSA %v", key, err)
	}

	tAcc := Account{
		0,
		[]*ecdsa.PrivateKey{acc},
		[]string{key},
		crypto.PubkeyToAddress(acc.PublicKey),
		0,
		big.NewInt(0),
		sync.Mutex{},
		//make(TransactionMap),
	}

	return &tAcc
}

func (account *Account) ImportUnLockAccount(endpoint string) {
	key := account.key[0]
	acc, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Fatalf("Key(%v): Failed to HexToECDSA %v", err)
	}

	testAddr := crypto.PubkeyToAddress(acc.PublicKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := client.Dial(endpoint)
	if err != nil {
		log.Fatalf("ImportUnLockAccount(): Create Client %v", err)
	}

	addr, err := c.ImportRawKey(ctx, key, "")
	if err != nil {
		log.Fatalf("Account(%v) : Failed to import => %v\n", account.address, err)
	} else {
		if testAddr != addr {
			log.Fatalf("origial:%v, imported: %v\n", testAddr.String(), addr.String())
		}
	}

	res, err := c.UnlockAccount(ctx, account.address, "", 0)
	if err != nil {
		log.Fatalf("Account(%v) : Failed to Unlock: %v\n", account.address.String(), err)
	} else {
		log.Printf("Wallet UnLock Result: %v", res)
	}
}

func NewAccount(id int) *Account {
	acc, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}

	testKey := hex.EncodeToString(crypto.FromECDSA(acc))

	tAcc := Account{
		0,
		[]*ecdsa.PrivateKey{acc},
		[]string{testKey},
		crypto.PubkeyToAddress(acc.PublicKey),
		0,
		big.NewInt(0),
		sync.Mutex{},
		//make(TransactionMap),
	}

	return &tAcc
}

func NewAccountOnNode(id int, endpoint string) *Account {

	tAcc := NewAccount(id)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := client.Dial(endpoint)
	if err != nil {
		log.Fatalf("NewAccountOnNode() : Failed to create client %v", err)
	}

	addr, err := c.ImportRawKey(ctx, tAcc.key[0], "")
	if err != nil {
		//log.Printf("Account(%v) : Failed to import\n", tAcc.address, err)
	} else {
		if tAcc.address != addr {
			log.Fatalf("origial:%v, imported: %v\n", tAcc.address, addr.String())
		}
		//log.Printf("origial:%v, imported:%v\n", tAcc.address, addr.String())
	}

	_, err = c.UnlockAccount(ctx, tAcc.GetAddress(), "", 0)
	if err != nil {
		log.Printf("Account(%v) : Failed to Unlock: %v\n", tAcc.GetAddress().String(), err)
	}

	//log.Printf("Wallet UnLock Result: %v", flag)

	return tAcc
}

func NewKlaytnAccount(id int) *Account {
	acc, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}

	testKey := hex.EncodeToString(crypto.FromECDSA(acc))

	randomAddr := common.BytesToAddress(crypto.Keccak256([]byte(testKey))[12:])

	tAcc := Account{
		0,
		[]*ecdsa.PrivateKey{acc},
		[]string{testKey},
		randomAddr,
		0,
		big.NewInt(0),
		sync.Mutex{},
		//make(TransactionMap),
	}

	return &tAcc
}

func NewKlaytnAccountWithAddr(id int, addr common.Address) *Account {
	acc, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}

	testKey := hex.EncodeToString(crypto.FromECDSA(acc))

	tAcc := Account{
		0,
		[]*ecdsa.PrivateKey{acc},
		[]string{testKey},
		addr,
		0,
		big.NewInt(0),
		sync.Mutex{},
		//make(TransactionMap),
	}

	return &tAcc
}

func NewKlaytnMultisigAccount(id int) *Account {
	k1, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}
	k2, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}
	k3, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}

	testKey := hex.EncodeToString(crypto.FromECDSA(k1))

	randomAddr := common.BytesToAddress(crypto.Keccak256([]byte(testKey))[12:])

	tAcc := Account{
		0,
		[]*ecdsa.PrivateKey{k1, k2, k3},
		[]string{testKey},
		randomAddr,
		0,
		big.NewInt(0),
		sync.Mutex{},
		//make(TransactionMap),
	}

	return &tAcc
}

func UnlockAccount(c *client.Client, addr common.Address, pwd string) {
	ctx := context.Background()
	defer ctx.Done()

	_, e := c.UnlockAccount(ctx, addr, pwd, 0)
	if e == nil {
	} else {
		fmt.Println(e)
	}
}

func (acc *Account) GetKey() *ecdsa.PrivateKey {
	return acc.privateKey[0]
}

func (acc *Account) GetAddress() common.Address {
	return acc.address
}

func (acc *Account) GetPrivateKey() string {
	return acc.key[0]
}

func (acc *Account) GetNonce(c *client.Client) uint64 {
	if acc.nonce != 0 {
		return acc.nonce
	}
	ctx := context.Background()
	nonce, err := c.NonceAt(ctx, acc.GetAddress(), nil)
	if err != nil {
		log.Printf("GetNonce(): Failed to NonceAt() %v\n", err)
		return acc.nonce
	}
	acc.nonce = nonce

	//fmt.Printf("account= %v  nonce = %v\n", acc.GetAddress().String(), nonce)
	return acc.nonce
}

func (acc *Account) GetNonceFromBlock(c *client.Client) uint64 {
	ctx := context.Background()
	nonce, err := c.NonceAt(ctx, acc.GetAddress(), nil)
	if err != nil {
		log.Printf("GetNonce(): Failed to NonceAt() %v\n", err)
		return acc.nonce
	}

	acc.nonce = nonce

	fmt.Printf("%v: account= %v  nonce = %v\n", os.Getpid(), acc.GetAddress().String(), nonce)
	return acc.nonce
}

func (acc *Account) UpdateNonce() {
	acc.nonce++
}

func (a *Account) GetReceipt(c *client.Client, txHash common.Hash) (*types.Receipt, error) {
	ctx := context.Background()
	return c.TransactionReceipt(ctx, txHash)
}

func (a *Account) GetBalance(c *client.Client) (*big.Int, error) {
	ctx := context.Background()
	balance, err := c.BalanceAt(ctx, a.GetAddress(), nil)
	if err != nil {
		return nil, err
	}
	return balance, err
}

func (self *Account) TransferSignedTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, gasPrice, err := self.TransferSignedTxReturnTx(true, c, to, value)
	return tx.Hash(), gasPrice, err
}

func (self *Account) TransferSignedTxWithoutLock(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	tx, gasPrice, err := self.TransferSignedTxReturnTx(false, c, to, value)
	return tx.Hash(), gasPrice, err
}

func (self *Account) TransferSignedTxReturnTx(withLock bool, c *client.Client, to *Account, value *big.Int) (*types.Transaction, *big.Int, error) {
	if withLock {
		self.mutex.Lock()
		defer self.mutex.Unlock()
	}

	nonce := self.GetNonce(c)

	//fmt.Printf("account=%v, nonce = %v\n", self.GetAddress().String(), nonce)

	tx := types.NewTransaction(
		nonce,
		to.GetAddress(),
		value,
		21000,
		gasPrice,
		nil)
	gasPrice := tx.GasPrice()
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), self.privateKey[0])
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	_, err = c.SendRawTransaction(ctx, signTx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return signTx, gasPrice, err
	}

	self.nonce++

	//fmt.Printf("%v transferSignedTx %v klay to %v klay.\n", self.GetAddress().Hex(), to.GetAddress().Hex(), value)

	return signTx, gasPrice, nil
}

func (self *Account) TransferNewValueTransferWithCancelTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	var txList []*types.Transaction
	nonce := self.GetNonce(c)

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyFrom:     self.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	txList = append(txList, tx)

	cancelTx, err := types.NewTransactionWithMap(types.TxTypeCancel, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyGasLimit: uint64(100000000),
		types.TxValueKeyGasPrice: gasPrice,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = cancelTx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	txList = append(txList, cancelTx)

	var hash common.Hash
	for _, tx := range txList {
		hash, err := c.SendRawTransaction(ctx, tx)
		if err != nil {
			if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
				fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
				fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
				self.nonce++
			} else {
				fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			}
			return hash, gasPrice, err
		}
	}

	self.nonce++
	return hash, gasPrice, nil
}

func (self *Account) TransferNewValueTransferTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyFrom:     self.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedValueTransferTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransfer, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyFeePayer: to.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = tx.SignFeePayerWithKeys(signer, to.privateKey)
	if err != nil {
		log.Fatalf("Failed to fee payer sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedValueTransferWithRatioTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferWithRatio, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:              nonce,
		types.TxValueKeyTo:                 to.GetAddress(),
		types.TxValueKeyAmount:             value,
		types.TxValueKeyGasLimit:           uint64(100000),
		types.TxValueKeyGasPrice:           gasPrice,
		types.TxValueKeyFrom:               self.address,
		types.TxValueKeyFeePayer:           to.address,
		types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = tx.SignFeePayerWithKeys(signer, to.privateKey)
	if err != nil {
		log.Fatalf("Failed to fee payer sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewValueTransferMemoTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	data := []byte("hello")

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransferMemo, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyData:     data,
		types.TxValueKeyFrom:     self.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// increase memo size from 5 bytes to between 50 bytes and 2,000 bytes

func (self *Account) TransferNewValueTransferBigRandomStringMemoTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)
	minBytes := 50
	maxBytes := 2000

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	data := randomString(randInt(minBytes, maxBytes))
	// data := []byte("hello")

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransferMemo, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyData:     data,
		types.TxValueKeyFrom:     self.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

// create 200 strings of memo
func (self *Account) TransferNewValueTransferSmallMemoTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)
	length := 200

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	data := randomString(length)
	// data := []byte("hello")

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransferMemo, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyData:     data,
		types.TxValueKeyFrom:     self.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

// create 2000 strings of memo
func (self *Account) TransferNewValueTransferLargeMemoTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)
	length := 2000

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	data := randomString(length)
	// data := []byte("hello")

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransferMemo, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyData:     data,
		types.TxValueKeyFrom:     self.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedValueTransferMemoTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	data := []byte("hello")

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferMemo, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyData:     data,
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyFeePayer: to.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = tx.SignFeePayerWithKeys(signer, to.privateKey)
	if err != nil {
		log.Fatalf("Failed to fee payer sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedValueTransferMemoWithRatioTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	data := []byte("hello")

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferMemoWithRatio, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:              nonce,
		types.TxValueKeyTo:                 to.GetAddress(),
		types.TxValueKeyAmount:             value,
		types.TxValueKeyGasLimit:           uint64(100000),
		types.TxValueKeyGasPrice:           gasPrice,
		types.TxValueKeyData:               data,
		types.TxValueKeyFrom:               self.address,
		types.TxValueKeyFeePayer:           to.address,
		types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = tx.SignFeePayerWithKeys(signer, to.privateKey)
	if err != nil {
		log.Fatalf("Failed to fee payer sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewAccountCreationTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeAccountCreation, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:         nonce,
		types.TxValueKeyFrom:          self.address,
		types.TxValueKeyTo:            to.GetAddress(),
		types.TxValueKeyAmount:        value,
		types.TxValueKeyGasLimit:      uint64(1000000),
		types.TxValueKeyGasPrice:      gasPrice,
		types.TxValueKeyHumanReadable: false,
		types.TxValueKeyAccountKey:    accountkey.NewAccountKeyPublicWithValue(&to.privateKey[0].PublicKey),
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewAccountUpdateTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:      nonce,
		types.TxValueKeyFrom:       self.address,
		types.TxValueKeyGasLimit:   uint64(100000),
		types.TxValueKeyGasPrice:   gasPrice,
		types.TxValueKeyAccountKey: accountkey.NewAccountKeyPublicWithValue(&self.privateKey[0].PublicKey),
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedAccountUpdateTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdate, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:      nonce,
		types.TxValueKeyFrom:       self.address,
		types.TxValueKeyGasLimit:   uint64(100000),
		types.TxValueKeyGasPrice:   gasPrice,
		types.TxValueKeyAccountKey: accountkey.NewAccountKeyPublicWithValue(&self.privateKey[0].PublicKey),
		types.TxValueKeyFeePayer:   to.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = tx.SignFeePayerWithKeys(signer, to.privateKey)
	if err != nil {
		log.Fatalf("Failed to fee payer sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedAccountUpdateWithRatioTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdateWithRatio, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:              nonce,
		types.TxValueKeyFrom:               self.address,
		types.TxValueKeyGasLimit:           uint64(100000),
		types.TxValueKeyGasPrice:           gasPrice,
		types.TxValueKeyAccountKey:         accountkey.NewAccountKeyPublicWithValue(&self.privateKey[0].PublicKey),
		types.TxValueKeyFeePayer:           to.address,
		types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = tx.SignFeePayerWithKeys(signer, to.privateKey)
	if err != nil {
		log.Fatalf("Failed to fee payer sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewSmartContractDeployTx(c *client.Client, to *Account, value *big.Int) (common.Address, *types.Transaction, *big.Int, error) {
	return self.TransferNewSmartContractDeployTxHumanReadable(c, to, value, false)
}

func (self *Account) DeployStorageTrieWrite(c *client.Client, to *Account, value *big.Int, humanReadable bool) (common.Address, *types.Transaction, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	if nonce != 0 {
		fmt.Println("Contract seems to already have been deployed!", "nonce", nonce)
		return common.Address{}, nil, nil, AlreadyDeployedErr
	}

	gaslimit := uint64(10000000)
	if humanReadable {
		gaslimit = uint64(4100000000)
	}

	contractABI := `[{"constant":true,"inputs":[],"name":"rootCaCertificate","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_serialNumber","type":"string"}],"name":"getIdentity","outputs":[{"name":"","type":"string"},{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_caKey","type":"string"}],"name":"deleteCaCertificate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_caKey","type":"string"},{"name":"_caCert","type":"string"}],"name":"insertCaCertificate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_serialNumber","type":"string"},{"name":"_publicKey","type":"string"},{"name":"_hash","type":"string"}],"name":"insertIdentity","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_serialNumber","type":"string"}],"name":"deleteIdentity","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_caKey","type":"string"}],"name":"getCaCertificate","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]`
	parsed, err := abi.JSON(strings.NewReader(contractABI))
	byteCode := common.FromHex("0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550610f76806100606000396000f30060806040526004361061008e576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806301c0ae49146100935780630a29ae6f146101235780631fde075b146102715780636bda98c3146102da5780638da5cb5b14610389578063b912b308146103e0578063bf951c68146104d5578063f09fdbef1461053e575b600080fd5b34801561009f57600080fd5b506100a8610620565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100e85780820151818401526020810190506100cd565b50505050905090810190601f1680156101155780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561012f57600080fd5b5061018a600480360381019080803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506106be565b604051808060200180602001838103835285818151815260200191508051906020019080838360005b838110156101ce5780820151818401526020810190506101b3565b50505050905090810190601f1680156101fb5780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b83811015610234578082015181840152602081019050610219565b50505050905090810190601f1680156102615780820380516001836020036101000a031916815260200191505b5094505050505060405180910390f35b34801561027d57600080fd5b506102d8600480360381019080803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506108af565b005b3480156102e657600080fd5b50610387600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610994565b005b34801561039557600080fd5b5061039e610a93565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156103ec57600080fd5b506104d3600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610ab8565b005b3480156104e157600080fd5b5061053c600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610baa565b005b34801561054a57600080fd5b506105a5600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610ca6565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156105e55780820151818401526020810190506105ca565b50505050905090810190601f1680156106125780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156106b65780601f1061068b576101008083540402835291602001916106b6565b820191906000526020600020905b81548152906001019060200180831161069957829003601f168201915b505050505081565b6060806106c9610dc3565b600084511115156106d957600080fd5b6003846040518082805190602001908083835b60208310151561071157805182526020820191506020810190506020830392506106ec565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020604080519081016040529081600082018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107e85780601f106107bd576101008083540402835291602001916107e8565b820191906000526020600020905b8154815290600101906020018083116107cb57829003601f168201915b50505050508152602001600182018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561088a5780601f1061085f5761010080835404028352916020019161088a565b820191906000526020600020905b81548152906001019060200180831161086d57829003601f168201915b5050505050815250509050806000015181602001518191508090509250925050915091565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561090a57600080fd5b6000815111151561091a57600080fd5b6002816040518082805190602001908083835b602083101515610952578051825260208201915060208101905060208303925061092d565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902060006109919190610ddd565b50565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156109ef57600080fd5b600082511115156109ff57600080fd5b60008151111515610a0f57600080fd5b806002836040518082805190602001908083835b602083101515610a485780518252602082019150602081019050602083039250610a23565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390209080519060200190610a8e929190610e25565b505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008351111515610ac857600080fd5b60008251111515610ad857600080fd5b60008151111515610ae857600080fd5b6040805190810160405280838152602001828152506003846040518082805190602001908083835b602083101515610b355780518252602082019150602081019050602083039250610b10565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390206000820151816000019080519060200190610b84929190610ea5565b506020820151816001019080519060200190610ba1929190610ea5565b50905050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610c0557600080fd5b60008151111515610c1557600080fd5b6003816040518082805190602001908083835b602083101515610c4d5780518252602082019150602081019050602083039250610c28565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902060008082016000610c919190610ddd565b600182016000610ca19190610ddd565b505050565b606060008251111515610cb857600080fd5b6002826040518082805190602001908083835b602083101515610cf05780518252602082019150602081019050602083039250610ccb565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390208054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610db75780601f10610d8c57610100808354040283529160200191610db7565b820191906000526020600020905b815481529060010190602001808311610d9a57829003601f168201915b50505050509050919050565b604080519081016040528060608152602001606081525090565b50805460018160011615610100020316600290046000825580601f10610e035750610e22565b601f016020900490600052602060002090810190610e219190610f25565b5b50565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610e6657805160ff1916838001178555610e94565b82800160010185558215610e94579182015b82811115610e93578251825591602001919060010190610e78565b5b509050610ea19190610f25565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610ee657805160ff1916838001178555610f14565b82800160010185558215610f14579182015b82811115610f13578251825591602001919060010190610ef8565b5b509050610f219190610f25565b5090565b610f4791905b80821115610f43576000816000905550600101610f2b565b5090565b905600a165627a7a7230582089a867aeaa08bec696937a378160fadb7e3ffe65cc89c1e648dec0b1359cd4e00029")

	if err != nil {
		fmt.Println("Error while parsing contractABI", "err", err)
	}
	txOpts := &bind.TransactOpts{
		From: self.address, Nonce: big.NewInt(int64(nonce)),
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != self.address {
				return nil, errors.New("not authorized to sign this account")
			}
			return types.SignTx(tx, signer, self.privateKey[0])
		}, Value: common.Big0,
		GasPrice: gasPrice, GasLimit: gaslimit, Context: ctx}
	contractAddr, contractTx, _, err := bind.DeployContract(txOpts, parsed, byteCode, c)
	if err != nil {
		log.Printf("Failed to deploy storage trie write performance test contract, err: %v, account: %v", err, self.address.String())
		return common.Address{}, nil, nil, err
	}

	self.nonce++
	return contractAddr, contractTx, gasPrice, nil
}

var AlreadyDeployedErr = errors.New("contract seems to already have been deployed")

func (self *Account) TransferNewSmartContractDeployTxHumanReadable(c *client.Client, to *Account, value *big.Int, humanReadable bool) (common.Address, *types.Transaction, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	code := "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"

	gaslimit := uint64(10000000)
	if humanReadable {
		gaslimit = uint64(4100000000)
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:         nonce,
		types.TxValueKeyFrom:          self.address,
		types.TxValueKeyTo:            (*common.Address)(nil),
		types.TxValueKeyAmount:        common.Big0,
		types.TxValueKeyGasLimit:      gaslimit,
		types.TxValueKeyGasPrice:      gasPrice,
		types.TxValueKeyHumanReadable: humanReadable,
		types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		types.TxValueKeyData:          common.FromHex(code),
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	_, err = c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return common.Address{}, tx, gasPrice, err
	}

	contractAddr := crypto.CreateAddress(self.address, self.nonce)

	self.nonce++

	return contractAddr, tx, gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedSmartContractDeployTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	code := "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeploy, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:         nonce,
		types.TxValueKeyFrom:          self.address,
		types.TxValueKeyTo:            &to.address,
		types.TxValueKeyAmount:        common.Big0,
		types.TxValueKeyGasLimit:      uint64(10000000),
		types.TxValueKeyGasPrice:      gasPrice,
		types.TxValueKeyHumanReadable: false,
		types.TxValueKeyData:          common.FromHex(code),
		types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		types.TxValueKeyFeePayer:      self.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = tx.SignFeePayerWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to fee payer sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedSmartContractDeployWithRatioTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	code := "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeployWithRatio, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:              nonce,
		types.TxValueKeyFrom:               self.address,
		types.TxValueKeyTo:                 &to.address,
		types.TxValueKeyAmount:             common.Big0,
		types.TxValueKeyGasLimit:           uint64(10000000),
		types.TxValueKeyGasPrice:           gasPrice,
		types.TxValueKeyHumanReadable:      false,
		types.TxValueKeyData:               common.FromHex(code),
		types.TxValueKeyFeePayer:           self.address,
		types.TxValueKeyCodeFormat:         params.CodeFormatEVM,
		types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = tx.SignFeePayerWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to fee payer sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = Letters[r.Intn(len(Letters))]
	}
	return string(b)
}

func (self *Account) ExecuteStorageTrieStore(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	abiStr := `[{"constant":true,"inputs":[],"name":"rootCaCertificate","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_serialNumber","type":"string"}],"name":"getIdentity","outputs":[{"name":"","type":"string"},{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_caKey","type":"string"}],"name":"deleteCaCertificate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_caKey","type":"string"},{"name":"_caCert","type":"string"}],"name":"insertCaCertificate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_serialNumber","type":"string"},{"name":"_publicKey","type":"string"},{"name":"_hash","type":"string"}],"name":"insertIdentity","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_serialNumber","type":"string"}],"name":"deleteIdentity","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_caKey","type":"string"}],"name":"getCaCertificate","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]`

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}
	data, err := abii.Pack("insertIdentity", randomString(39), randomString(814), randomString(40))
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: uint64(5000000),
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   common.Big0,
		types.TxValueKeyTo:       to.address,
		types.TxValueKeyData:     data,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	// log.Printf("data %s", common.Bytes2Hex(data))
	// log.Printf("to.address %s", to.address.String())
	// log.Printf("tx %s\n", tx.String())

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func encodeHeader(block *types.Header) []byte {
	encoded, _ := rlp.EncodeToBytes([]interface{}{
		block.ParentHash,
		block.Rewardbase,
		block.Root,
		block.TxHash,
		block.ReceiptHash,
		block.Bloom,
		block.BlockScore,
		block.Number,
		block.GasUsed,
		block.Time,
		block.TimeFoS,
		block.Extra,
		block.Governance,
		block.Vote,
		block.BaseFee,
	})
	return encoded
}

func encodeReceipt(receipt *types.Receipt) ([]byte, error) {
	var logs []interface{}
	for j := 0; j < len(receipt.Logs); j++ {
		logs = append(logs, []interface{}{receipt.Logs[j].Address, receipt.Logs[j].Topics, receipt.Logs[j].Data})
	}
	encReceipt, err := rlp.EncodeToBytes([]interface{}{receipt.Status, receipt.GasUsed, receipt.Bloom, logs})
	return encReceipt, err
}

func (self *Account) TransferBridgeTx(c *client.Client, tokenContractAddr, verifierContractAddr common.Address, to *Account, value *big.Int, toChainId uint) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	//gaslimit := uint64(10000000)
	gaslimit := uint64(4100000000)

	// 1. approve
	var erc20PerformanceABI = `[{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"sender","type":"address"},{"name":"recipient","type":"address"},{"name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"account","type":"address"},{"name":"amount","type":"uint256"}],"name":"mint","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"account","type":"address"}],"name":"addMinter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"renounceMinter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"recipient","type":"address"},{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"isMinter","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"}],"name":"MinterAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"}],"name":"MinterRemoved","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"}]`
	contractABI := erc20PerformanceABI
	parsed, err := abi.JSON(strings.NewReader(contractABI))

	amount := &big.Int{}
	amount.SetInt64(100)
	data, err := parsed.Pack("approve", VerifierContractAddr, amount)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		//types.TxValueKeyGasLimit: uint64(5000000),
		types.TxValueKeyGasLimit: gaslimit,
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   value,
		//types.TxValueKeyTo:       to.address,
		types.TxValueKeyTo:   SmartContractAddr,
		types.TxValueKeyData: data,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Bridge approve Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Bridge approve Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Bridge approve Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return common.Hash{0}, big.NewInt(0), err
	}
	fmt.Printf("Verifier approve hash: 0x%v\n", hex.EncodeToString(hash[:]))

	self.nonce++
	nonce++

	// 2. erc20TransferRequest
	// Get blocknumber before sending erc20TransferRequest
	curBlockNum, err := c.BlockNumber(ctx)
	if err != nil {
		return common.Hash{0}, big.NewInt(0), err
	}

	// verifier
	abiStr := `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"requestTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"handleTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"from","type":"address"},{"indexed":false,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"bytes32","name":"reqHeaderHash","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"key","type":"bytes"}],"name":"BridgeERC20TransferHandle","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":true,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"address","name":"burnedTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"ctTokenAddr","type":"address"}],"name":"BridgeERC20TransferRequest","type":"event"},{"anonymous":false,"inputs":[{"components":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"indexed":true,"internalType":"struct IManager.FeeSchedule","name":"feeSchedule","type":"tuple"}],"name":"CancelScheduledFee","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"feeReceiver","type":"address"}],"name":"ChangedFeeReceiver","type":"event"},{"anonymous":false,"inputs":[{"components":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"indexed":true,"internalType":"struct IManager.FeeSchedule","name":"feeSchedule","type":"tuple"}],"name":"FeeScheduled","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"blockNumber","type":"uint256"}],"name":"Frozen","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"headerNumber","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"requiredQuorum","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"committeeSize","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"validatorSize","type":"uint256"}],"name":"HeaderIntegrityPassed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"feeReceiver","type":"address"},{"indexed":false,"internalType":"uint256","name":"bridgeFee","type":"uint256"},{"indexed":false,"internalType":"address","name":"headerContributor","type":"address"},{"indexed":false,"internalType":"uint256","name":"govContributionFee","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"headerContributionFee","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"verifierBalance","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"userSent","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"change","type":"uint256"}],"name":"PayBridgeFee","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"onChainToken","type":"address"},{"indexed":true,"internalType":"address","name":"counterpartChainToken","type":"address"},{"indexed":true,"internalType":"bytes32","name":"transferEvent","type":"bytes32"}],"name":"TokenPairRegistered","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"governanceContributor","type":"address"},{"indexed":false,"internalType":"uint256","name":"governanceContributionFee","type":"uint256"}],"name":"TransferGovContrib","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"blockNumber","type":"uint256"}],"name":"Unfrozen","type":"event"},{"inputs":[],"name":"bridgeLibAddr","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"cancelFeeSchedule","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ctToken2TransferEvent","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ctToken2ocToken","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"address","name":"ocTokenAddr","type":"address"},{"internalType":"address","name":"ctTokenAddr","type":"address"}],"name":"erc20TransferRequest","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"feeReceiver","outputs":[{"internalType":"address payable","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"feeSchedule","outputs":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"freeze","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getStampBlockNumbers","outputs":[{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"}],"name":"getTransferReceipt","outputs":[{"components":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"address","name":"tokenAddr","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"internalType":"struct Typs.BridgeTransferStamp[]","name":"","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getVersion","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"govStateAddr","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_bridgeLibAddr","type":"address"},{"internalType":"address","name":"_govStateAddr","type":"address"},{"internalType":"bool","name":"_mintable","type":"bool"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"isVerifiedHeader","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"mintable","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ocToken2ctToken","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"prevFeeSchdule","outputs":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"effectAt","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"}],"name":"registerFeeSchedule","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"onChainToken","type":"address"},{"internalType":"address","name":"counterpartChainToken","type":"address"},{"internalType":"bytes32","name":"transferEventSig","type":"bytes32"}],"name":"registerTokenPair","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"verifier","type":"address"}],"name":"registerVerifier","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address payable","name":"_feeReceiver","type":"address"}],"name":"setFeeReceiver","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"headerHash","type":"bytes32"},{"internalType":"bytes","name":"rlpEncodedReceipt","type":"bytes"},{"internalType":"bytes[]","name":"MPTProof","type":"bytes[]"},{"internalType":"bytes","name":"key","type":"bytes"}],"name":"submitGov","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"rlpEncodedHeader","type":"bytes"}],"name":"submitHeader","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"headerHash","type":"bytes32"},{"internalType":"bytes","name":"rlpEncodedReceipt","type":"bytes"},{"internalType":"bytes[]","name":"MPTProof","type":"bytes[]"},{"internalType":"bytes","name":"key","type":"bytes"}],"name":"submitReceipt","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"temporalStop","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"transferStamps","outputs":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"address","name":"tokenAddr","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"unfreeze","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"verifiedHeaders","outputs":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"bytes32","name":"receiptHash","type":"bytes32"},{"internalType":"bytes","name":"vote","type":"bytes"},{"internalType":"bytes","name":"governance","type":"bytes"},{"internalType":"address","name":"contributor","type":"address"},{"internalType":"bool","name":"stored","type":"bool"}],"stateMutability":"view","type":"function"}]`

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}
	//log.Println("bridgeTx tokenContractAddr", tokenContractAddr.String(), "verifierContractAddr", verifierContractAddr.String(), "tokenRecipientAddr", to.address.String())

	data, err = abii.Pack("erc20TransferRequest", to.address, value, tokenContractAddr, tokenContractAddr)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer = types.NewEIP155Signer(chainID)
	tx, err = types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: uint64(5000000),
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   common.Big0,
		types.TxValueKeyTo:       verifierContractAddr,
		types.TxValueKeyData:     data,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err = c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("erc20TransferRequest Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("erc20TransferRequest Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("erc20TransferRequest Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return common.Hash{0}, big.NewInt(0), err
	}
	fmt.Printf("erc20TransferRequest hash: 0x%v\n", hex.EncodeToString(hash[:]))

	self.nonce++

	// 3. Make proof
	// - encoded header
	// - encoded receipt
	// - mpt proof and key

	var resTx *types.Transaction
	resTx, _, err = c.TransactionByHash(ctx, hash)
	bind.WaitMined(ctx, c, resTx)

	var blockHash common.Hash
	for i := curBlockNum; i.Cmp(big.NewInt(0)) == 1; i = i.Add(i, big.NewInt(1)) {
		time.Sleep(1 * time.Second)
		fmt.Printf("Searching tx in blockNum: %d\n", i.Int64())
		if i.Cmp(big.NewInt(0)) == 0 {
			return hash, gasPrice, fmt.Errorf("Failed to find block with out tx")
		}
		block, err := c.BlockByNumber(ctx, i)
		if err != nil {
			return hash, gasPrice, err
		}
		count, err := c.TransactionCount(ctx, block.Hash())
		if err != nil {
			return hash, gasPrice, err
		}
		found := false
		for j := 0; j < int(count); j++ {
			fmt.Printf("Comparing %dth tx in blockNum: %d\n", j, i.Int64())
			tmpTx, err := c.TransactionInBlock(ctx, block.Hash(), uint(j))
			if err != nil {
				return hash, gasPrice, err
			}
			if resTx.Hash() == tmpTx.Hash() {
				blockHash = block.Hash()
				found = true
				break
			}
		}
		if found {
			fmt.Println("Found")
			break
		}
	}

	header, err := c.HeaderByHash(ctx, blockHash)
	if err != nil {
		log.Fatalf("Failed to find header with hash: %v", err)
		return hash, gasPrice, err
	}
	encodedHeader := encodeHeader(header)

	count, err := c.TransactionCount(ctx, blockHash)
	if err != nil {
		log.Fatalf("Failed to get tx count: %v", err)
		return hash, gasPrice, err
	}

	var receipt *types.Receipt
	var encodedReceipt []byte
	var encReceipt []byte

	for i := 0; i < int(count); i++ {
		tx, err := c.TransactionInBlock(ctx, blockHash, uint(i))
		if err != nil {
			log.Fatalf("Failed to get tx in block: %v", err)
			return hash, gasPrice, err
		}
		receipt, err = c.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			log.Fatalf("Failed to get receipt: %v", err)
			return hash, gasPrice, err
		}
		encodedReceipt, err = encodeReceipt(receipt)
		if err != nil {
			log.Fatalf("Failed to encode receipt: %v", err)
			return hash, gasPrice, err
		}
		if tx.Hash() == resTx.Hash() {
			encReceipt = encodedReceipt
		}
	}

	// Get proof using hardhat
	cmd := exec.Command("npx", "hardhat", "getProof", "--tx", hash.Hex())
	cmd.Dir = HardhatPath
	stdout, err := cmd.Output()

	if err != nil {
		log.Fatalf("Failed to run npx: %v", err)
	}
	outputStr := string(stdout)
	outputStr = strings.ReplaceAll(outputStr, "'", "\"")
	var proofStr []string
	//fmt.Printf("output: %v\n", string(stdout))
	err = json.Unmarshal([]byte(outputStr), &proofStr)
	if err != nil {
		log.Fatalf("Failed to parse proof from npx output: %v", err)
	}
	var proof [][]byte
	for i := 0; i < len(proofStr); i++ {
		decoded, err := hex.DecodeString(proofStr[i][2:])
		if err != nil {
			log.Fatalf("Failed to decode proof[%d] into []byte: %v", i, err)
		}
		proof = append(proof, decoded)
	}

	// Get key from hardhat
	cmd = exec.Command("npx", "hardhat", "getIndexKey", "--tx", hash.Hex())
	cmd.Dir = HardhatPath
	stdout, err = cmd.Output()

	if err != nil {
		log.Fatalf("Failed to run npx: %v", err)
	}
	key := strings.TrimSpace(string(stdout))
	//fmt.Printf("key: %v\n", key)
	//fmt.Println("proof: %v\n", proof)

	// 4. Submit header
	nonce = self.GetNonce(c)

	data, err = abii.Pack("submitHeader", encodedHeader)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer = types.NewEIP155Signer(chainID)
	tx, err = types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: uint64(5000000000000),
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   common.Big0,
		types.TxValueKeyTo:       verifierContractAddr,
		types.TxValueKeyData:     data,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err = c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("SubmitHeader Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("SubmitHeader Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("SubmitHeader Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	bind.WaitMined(ctx, c, tx)

	// print result with related info
	receipt, err = c.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		log.Fatalf("Failed to get receipt for erc20TransferVerify: %v", err)
	}
	blockNum, err := c.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("failed to get blocknum: %v\n", err)
	}
	header, err = c.HeaderByNumber(ctx, blockNum)
	if err != nil {
		log.Fatalf("failed to get header: %v\n", err)
	}
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		log.Fatalf("failed to get istanbulextra: %v\n", err)
	}
	gasUsed := new(big.Int).SetUint64(receipt.GasUsed)
	fmt.Printf("submitHeader hash: 0x%v, status: %d, gasUsed: %v, validator len: %d\n", hex.EncodeToString(hash[:]), receipt.Status, gasUsed.String(), len(istanbulExtra.Validators))

	// 5. submitReceipt

	nonce = self.GetNonce(c)
	// 1 KLAY
	//fee := big.NewInt(uint64(1000000000000000000))
	fee := &big.Int{}
	// 9 klay
	fee.SetUint64(uint64(9000000000000000000))
	keyByte, err := hex.DecodeString(key[2:])
	if err != nil {
		log.Fatalf("Failed to decode key into []byte: %v", err)
	}
	data, err = abii.Pack("submitReceipt", blockHash, encReceipt, proof, keyByte)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer = types.NewEIP155Signer(chainID)
	tx, err = types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: uint64(5000000000000),
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   fee,
		types.TxValueKeyTo:       verifierContractAddr,
		types.TxValueKeyData:     data,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err = c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("SubmitReceipt Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("SubmitReceipt Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("SubmitReceipt Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	// Wait and log gasUsed
	bind.WaitMined(ctx, c, tx)
	receipt, err = c.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		log.Fatalf("Failed to get receipt for erc20TransferVerify: %v", err)
	}
	gasUsed = new(big.Int).SetUint64(receipt.GasUsed)
	fmt.Printf("submitReceipt hash: 0x%v, status: %d, gasUsed: %v, validator len: %d, tx count: %d\n", hex.EncodeToString(hash[:]), receipt.Status, gasUsed.String(), len(istanbulExtra.Validators), count)
	return hash, gasPrice, nil
}

func (self *Account) Approve(c *client.Client, tokenContractAddr, verifierContractAddr common.Address, to *Account, value *big.Int, toChainId uint) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	gaslimit := uint64(4100000000)

	// approve
	var erc20PerformanceABI = `[{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"sender","type":"address"},{"name":"recipient","type":"address"},{"name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"account","type":"address"},{"name":"amount","type":"uint256"}],"name":"mint","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"account","type":"address"}],"name":"addMinter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"renounceMinter","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"recipient","type":"address"},{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"isMinter","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"}],"name":"MinterAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"}],"name":"MinterRemoved","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"}]`
	contractABI := erc20PerformanceABI
	parsed, err := abi.JSON(strings.NewReader(contractABI))

	amount := &big.Int{}
	amount.SetInt64(1e10)
	data, err := parsed.Pack("approve", VerifierContractAddr, amount)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: gaslimit,
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   common.Big0,
		types.TxValueKeyTo:       tokenContractAddr,
		types.TxValueKeyData:     data,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("approve Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("approve Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("approve Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return common.Hash{0}, big.NewInt(0), err
	}
	fmt.Printf("approve hash: 0x%v\n", hex.EncodeToString(hash[:]))
	bind.WaitMined(ctx, c, tx)

	self.nonce++
	nonce++
	return hash, gasPrice, nil
}

func (self *Account) TransferBridgeErc20Request(c *client.Client, tokenContractAddr, verifierContractAddr common.Address, to *Account, value *big.Int, toChainId uint) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	// verifier
	abiStr := `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"requestTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"handleTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"from","type":"address"},{"indexed":false,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"bytes32","name":"reqHeaderHash","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"key","type":"bytes"}],"name":"BridgeERC20TransferHandle","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":true,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"address","name":"burnedTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"ctTokenAddr","type":"address"}],"name":"BridgeERC20TransferRequest","type":"event"},{"anonymous":false,"inputs":[{"components":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"indexed":true,"internalType":"struct IManager.FeeSchedule","name":"feeSchedule","type":"tuple"}],"name":"CancelScheduledFee","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"feeReceiver","type":"address"}],"name":"ChangedFeeReceiver","type":"event"},{"anonymous":false,"inputs":[{"components":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"indexed":true,"internalType":"struct IManager.FeeSchedule","name":"feeSchedule","type":"tuple"}],"name":"FeeScheduled","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"blockNumber","type":"uint256"}],"name":"Frozen","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"headerNumber","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"requiredQuorum","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"committeeSize","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"validatorSize","type":"uint256"}],"name":"HeaderIntegrityPassed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"feeReceiver","type":"address"},{"indexed":false,"internalType":"uint256","name":"bridgeFee","type":"uint256"},{"indexed":false,"internalType":"address","name":"headerContributor","type":"address"},{"indexed":false,"internalType":"uint256","name":"govContributionFee","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"headerContributionFee","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"verifierBalance","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"userSent","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"change","type":"uint256"}],"name":"PayBridgeFee","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"onChainToken","type":"address"},{"indexed":true,"internalType":"address","name":"counterpartChainToken","type":"address"},{"indexed":true,"internalType":"bytes32","name":"transferEvent","type":"bytes32"}],"name":"TokenPairRegistered","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"governanceContributor","type":"address"},{"indexed":false,"internalType":"uint256","name":"governanceContributionFee","type":"uint256"}],"name":"TransferGovContrib","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"blockNumber","type":"uint256"}],"name":"Unfrozen","type":"event"},{"inputs":[],"name":"bridgeLibAddr","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"cancelFeeSchedule","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ctToken2TransferEvent","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ctToken2ocToken","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"address","name":"ocTokenAddr","type":"address"},{"internalType":"address","name":"ctTokenAddr","type":"address"}],"name":"erc20TransferRequest","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"feeReceiver","outputs":[{"internalType":"address payable","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"feeSchedule","outputs":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"freeze","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getStampBlockNumbers","outputs":[{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"}],"name":"getTransferReceipt","outputs":[{"components":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"address","name":"tokenAddr","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"internalType":"struct Typs.BridgeTransferStamp[]","name":"","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getVersion","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"govStateAddr","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_bridgeLibAddr","type":"address"},{"internalType":"address","name":"_govStateAddr","type":"address"},{"internalType":"bool","name":"_mintable","type":"bool"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"isVerifiedHeader","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"mintable","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ocToken2ctToken","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"prevFeeSchdule","outputs":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"effectAt","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"}],"name":"registerFeeSchedule","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"onChainToken","type":"address"},{"internalType":"address","name":"counterpartChainToken","type":"address"},{"internalType":"bytes32","name":"transferEventSig","type":"bytes32"}],"name":"registerTokenPair","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"verifier","type":"address"}],"name":"registerVerifier","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address payable","name":"_feeReceiver","type":"address"}],"name":"setFeeReceiver","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"headerHash","type":"bytes32"},{"internalType":"bytes","name":"rlpEncodedReceipt","type":"bytes"},{"internalType":"bytes[]","name":"MPTProof","type":"bytes[]"},{"internalType":"bytes","name":"key","type":"bytes"}],"name":"submitGov","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"rlpEncodedHeader","type":"bytes"}],"name":"submitHeader","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"headerHash","type":"bytes32"},{"internalType":"bytes","name":"rlpEncodedReceipt","type":"bytes"},{"internalType":"bytes[]","name":"MPTProof","type":"bytes[]"},{"internalType":"bytes","name":"key","type":"bytes"}],"name":"submitReceipt","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"temporalStop","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"transferStamps","outputs":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"address","name":"tokenAddr","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"unfreeze","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"verifiedHeaders","outputs":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"bytes32","name":"receiptHash","type":"bytes32"},{"internalType":"bytes","name":"vote","type":"bytes"},{"internalType":"bytes","name":"governance","type":"bytes"},{"internalType":"address","name":"contributor","type":"address"},{"internalType":"bool","name":"stored","type":"bool"}],"stateMutability":"view","type":"function"}]`

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}

	data, err := abii.Pack("erc20TransferRequest", to.address, value, tokenContractAddr, tokenContractAddr)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: uint64(5000000),
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   common.Big0,
		types.TxValueKeyTo:       verifierContractAddr,
		types.TxValueKeyData:     data,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("erc20TransferRequest Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("erc20TransferReqeust Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("erc20TransferRequest Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferSubmitHeader(c *client.Client, verifierContractAddr common.Address, blockNum int64) (common.Hash, common.Hash, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	header, err := c.HeaderByNumber(ctx, big.NewInt(blockNum))
	if err != nil {
		log.Fatalf("failed to get header: %v\n", err)
	}
	encodedHeader := encodeHeader(header)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	abiStr := `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"requestTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"handleTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"from","type":"address"},{"indexed":false,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"bytes32","name":"reqHeaderHash","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"key","type":"bytes"}],"name":"BridgeERC20TransferHandle","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":true,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"address","name":"burnedTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"ctTokenAddr","type":"address"}],"name":"BridgeERC20TransferRequest","type":"event"},{"anonymous":false,"inputs":[{"components":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"indexed":true,"internalType":"struct IManager.FeeSchedule","name":"feeSchedule","type":"tuple"}],"name":"CancelScheduledFee","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"feeReceiver","type":"address"}],"name":"ChangedFeeReceiver","type":"event"},{"anonymous":false,"inputs":[{"components":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"indexed":true,"internalType":"struct IManager.FeeSchedule","name":"feeSchedule","type":"tuple"}],"name":"FeeScheduled","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"blockNumber","type":"uint256"}],"name":"Frozen","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"headerNumber","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"requiredQuorum","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"committeeSize","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"validatorSize","type":"uint256"}],"name":"HeaderIntegrityPassed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"feeReceiver","type":"address"},{"indexed":false,"internalType":"uint256","name":"bridgeFee","type":"uint256"},{"indexed":false,"internalType":"address","name":"headerContributor","type":"address"},{"indexed":false,"internalType":"uint256","name":"govContributionFee","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"headerContributionFee","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"verifierBalance","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"userSent","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"change","type":"uint256"}],"name":"PayBridgeFee","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"onChainToken","type":"address"},{"indexed":true,"internalType":"address","name":"counterpartChainToken","type":"address"},{"indexed":true,"internalType":"bytes32","name":"transferEvent","type":"bytes32"}],"name":"TokenPairRegistered","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"governanceContributor","type":"address"},{"indexed":false,"internalType":"uint256","name":"governanceContributionFee","type":"uint256"}],"name":"TransferGovContrib","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"blockNumber","type":"uint256"}],"name":"Unfrozen","type":"event"},{"inputs":[],"name":"bridgeLibAddr","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"cancelFeeSchedule","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ctToken2TransferEvent","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ctToken2ocToken","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"address","name":"ocTokenAddr","type":"address"},{"internalType":"address","name":"ctTokenAddr","type":"address"}],"name":"erc20TransferRequest","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"feeReceiver","outputs":[{"internalType":"address payable","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"feeSchedule","outputs":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"freeze","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getStampBlockNumbers","outputs":[{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"}],"name":"getTransferReceipt","outputs":[{"components":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"address","name":"tokenAddr","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"internalType":"struct Typs.BridgeTransferStamp[]","name":"","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getVersion","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"govStateAddr","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_bridgeLibAddr","type":"address"},{"internalType":"address","name":"_govStateAddr","type":"address"},{"internalType":"bool","name":"_mintable","type":"bool"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"isVerifiedHeader","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"mintable","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ocToken2ctToken","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"prevFeeSchdule","outputs":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"effectAt","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"}],"name":"registerFeeSchedule","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"onChainToken","type":"address"},{"internalType":"address","name":"counterpartChainToken","type":"address"},{"internalType":"bytes32","name":"transferEventSig","type":"bytes32"}],"name":"registerTokenPair","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"verifier","type":"address"}],"name":"registerVerifier","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address payable","name":"_feeReceiver","type":"address"}],"name":"setFeeReceiver","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"headerHash","type":"bytes32"},{"internalType":"bytes","name":"rlpEncodedReceipt","type":"bytes"},{"internalType":"bytes[]","name":"MPTProof","type":"bytes[]"},{"internalType":"bytes","name":"key","type":"bytes"}],"name":"submitGov","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"rlpEncodedHeader","type":"bytes"}],"name":"submitHeader","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"headerHash","type":"bytes32"},{"internalType":"bytes","name":"rlpEncodedReceipt","type":"bytes"},{"internalType":"bytes[]","name":"MPTProof","type":"bytes[]"},{"internalType":"bytes","name":"key","type":"bytes"}],"name":"submitReceipt","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"temporalStop","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"transferStamps","outputs":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"address","name":"tokenAddr","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"unfreeze","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"verifiedHeaders","outputs":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"bytes32","name":"receiptHash","type":"bytes32"},{"internalType":"bytes","name":"vote","type":"bytes"},{"internalType":"bytes","name":"governance","type":"bytes"},{"internalType":"address","name":"contributor","type":"address"},{"internalType":"bool","name":"stored","type":"bool"}],"stateMutability":"view","type":"function"}]`

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}

	data, err := abii.Pack("submitHeader", encodedHeader)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: uint64(5000000000000),
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   common.Big0,
		types.TxValueKeyTo:       verifierContractAddr,
		types.TxValueKeyData:     data,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("SubmitHeader Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("SubmitHeader Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("SubmitHeader Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return common.Hash{0}, common.Hash{0}, err
	}

	self.nonce++
	go func() {
		receipt, err := bind.WaitMined(ctx, c, tx)
		if err != nil {
			log.Fatalf("Failed to get receipt for submitHeader tx: %v", err)
			return
		}
		gasUsed := new(big.Int).SetUint64(receipt.GasUsed)
		fmt.Printf("submitHeader for blockNum %d hash: 0x%v, status: %d, gasUsed: %v\n", blockNum, hex.EncodeToString(hash[:]), receipt.Status, gasUsed.String())
	}()

	return header.Hash(), hash, nil
}

func (self *Account) TransferSubmitReceipt(c *client.Client, verifierContractAddr common.Address, blockNum int64, blockHash, submitHeaderTxHash common.Hash, endpoint string) (common.Hash, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	abiStr := `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"requestTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"handleTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"from","type":"address"},{"indexed":false,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"bytes32","name":"reqHeaderHash","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"key","type":"bytes"}],"name":"BridgeERC20TransferHandle","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":true,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"address","name":"burnedTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"ctTokenAddr","type":"address"}],"name":"BridgeERC20TransferRequest","type":"event"},{"anonymous":false,"inputs":[{"components":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"indexed":true,"internalType":"struct IManager.FeeSchedule","name":"feeSchedule","type":"tuple"}],"name":"CancelScheduledFee","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"feeReceiver","type":"address"}],"name":"ChangedFeeReceiver","type":"event"},{"anonymous":false,"inputs":[{"components":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"indexed":true,"internalType":"struct IManager.FeeSchedule","name":"feeSchedule","type":"tuple"}],"name":"FeeScheduled","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"blockNumber","type":"uint256"}],"name":"Frozen","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"headerNumber","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"requiredQuorum","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"committeeSize","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"validatorSize","type":"uint256"}],"name":"HeaderIntegrityPassed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"feeReceiver","type":"address"},{"indexed":false,"internalType":"uint256","name":"bridgeFee","type":"uint256"},{"indexed":false,"internalType":"address","name":"headerContributor","type":"address"},{"indexed":false,"internalType":"uint256","name":"govContributionFee","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"headerContributionFee","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"verifierBalance","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"userSent","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"change","type":"uint256"}],"name":"PayBridgeFee","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"onChainToken","type":"address"},{"indexed":true,"internalType":"address","name":"counterpartChainToken","type":"address"},{"indexed":true,"internalType":"bytes32","name":"transferEvent","type":"bytes32"}],"name":"TokenPairRegistered","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"governanceContributor","type":"address"},{"indexed":false,"internalType":"uint256","name":"governanceContributionFee","type":"uint256"}],"name":"TransferGovContrib","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"blockNumber","type":"uint256"}],"name":"Unfrozen","type":"event"},{"inputs":[],"name":"bridgeLibAddr","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"cancelFeeSchedule","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ctToken2TransferEvent","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ctToken2ocToken","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"address","name":"ocTokenAddr","type":"address"},{"internalType":"address","name":"ctTokenAddr","type":"address"}],"name":"erc20TransferRequest","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"feeReceiver","outputs":[{"internalType":"address payable","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"feeSchedule","outputs":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"freeze","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getStampBlockNumbers","outputs":[{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"}],"name":"getTransferReceipt","outputs":[{"components":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"address","name":"tokenAddr","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"internalType":"struct Typs.BridgeTransferStamp[]","name":"","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getVersion","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"govStateAddr","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_bridgeLibAddr","type":"address"},{"internalType":"address","name":"_govStateAddr","type":"address"},{"internalType":"bool","name":"_mintable","type":"bool"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"isVerifiedHeader","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"mintable","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ocToken2ctToken","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"prevFeeSchdule","outputs":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"effectAt","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"}],"name":"registerFeeSchedule","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"onChainToken","type":"address"},{"internalType":"address","name":"counterpartChainToken","type":"address"},{"internalType":"bytes32","name":"transferEventSig","type":"bytes32"}],"name":"registerTokenPair","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"verifier","type":"address"}],"name":"registerVerifier","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address payable","name":"_feeReceiver","type":"address"}],"name":"setFeeReceiver","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"headerHash","type":"bytes32"},{"internalType":"bytes","name":"rlpEncodedReceipt","type":"bytes"},{"internalType":"bytes[]","name":"MPTProof","type":"bytes[]"},{"internalType":"bytes","name":"key","type":"bytes"}],"name":"submitGov","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"rlpEncodedHeader","type":"bytes"}],"name":"submitHeader","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"headerHash","type":"bytes32"},{"internalType":"bytes","name":"rlpEncodedReceipt","type":"bytes"},{"internalType":"bytes[]","name":"MPTProof","type":"bytes[]"},{"internalType":"bytes","name":"key","type":"bytes"}],"name":"submitReceipt","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"temporalStop","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"transferStamps","outputs":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"address","name":"tokenAddr","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"unfreeze","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"verifiedHeaders","outputs":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"bytes32","name":"receiptHash","type":"bytes32"},{"internalType":"bytes","name":"vote","type":"bytes"},{"internalType":"bytes","name":"governance","type":"bytes"},{"internalType":"address","name":"contributor","type":"address"},{"internalType":"bool","name":"stored","type":"bool"}],"stateMutability":"view","type":"function"}]`

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}

	count, err := c.TransactionCount(ctx, blockHash)
	if err != nil {
		log.Fatalf("Failed to get tx count: %v", err)
		return common.Hash{0}, err
	}

	if count <= MaxTxCount {
		return blockHash, nil
	}
	MaxTxCount = count

	fee := &big.Int{}
	fee.SetUint64(uint64(9000000000000000000))
	bridgeTxCount := 0

	signature, _ := hex.DecodeString("7a1ddb9dc7e950a74cfed8b204dc38c5d21c9225625a585d6111d9578d02c5ba")

	txs := make([]*types.Transaction, 0)
	indexes := make([]int, 0)
	for i := 0; i < int(count); i++ {
		tx, err := c.TransactionInBlock(ctx, blockHash, uint(i))
		if err != nil {
			log.Fatalf("Failed to get tx in block: %v", err)
			return common.Hash{0}, err
		}
		receipt, err := c.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			log.Fatalf("Failed to get receipt: %v", err)
			return common.Hash{0}, err
		}

		sigFound := false
		for _, log := range receipt.Logs {
			if bytes.Equal(log.Topics[0][:], signature[:]) {
				sigFound = true
			}
		}

		if !sigFound {
			continue
		}
		fmt.Printf("found bridge tx with index: %d\n", i)
		bridgeTxCount++

		encodedReceipt, err := encodeReceipt(receipt)
		if err != nil {
			log.Fatalf("Failed to encode receipt: %v", err)
			return common.Hash{0}, err
		}

		// Get proof
		type ProofOutput struct {
			Proof []string `json:"proof"`
			Key   string   `json:"key"`
		}
		var proofOutput ProofOutput
		var cmd *exec.Cmd
		var stdout []byte
		succeed := false
		for i := 0; i < 5; i++ {
			cmd = exec.Command("node", "proof.js", endpoint, receipt.TxHash.Hex())
			stdout, err = cmd.Output()

			if err != nil {
				fmt.Printf("Failed %d times to run npx for proof: %v, txHash: %v\n", i, err, receipt.TxHash.Hex())
				fmt.Printf("output: %v\n", string(stdout))
				time.Sleep(1 * time.Second)
				continue
			}
			outputStr := string(stdout)
			outputStr = strings.ReplaceAll(outputStr, "'", "\"")
			err = json.Unmarshal([]byte(outputStr), &proofOutput)
			if err != nil {
				fmt.Printf("output: %v\n", string(stdout))
				log.Fatalf("Failed to parse proof from npx output: %v", err)
			}
			succeed = true
		}
		if !succeed {
			log.Fatalf("Failed to run npx for proof: %v, txHash: %v", err, receipt.TxHash.Hex())
		}
		var proof [][]byte
		for i := 0; i < len(proofOutput.Proof); i++ {
			decoded, err := hex.DecodeString(proofOutput.Proof[i][2:])
			if err != nil {
				log.Fatalf("Failed to decode proof[%d] into []byte: %v", i, err)
			}
			proof = append(proof, decoded)
		}

		keyByte, err := hex.DecodeString(proofOutput.Key[2:])

		// Make tx
		data, err := abii.Pack("submitReceipt", blockHash, encodedReceipt, proof, keyByte)
		if err != nil {
			log.Fatalf("failed to abi.Pack: %v", err)
		}

		signer := types.NewEIP155Signer(chainID)
		tx, err = types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    nonce,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyGasLimit: uint64(5000000000000),
			types.TxValueKeyFrom:     self.address,
			types.TxValueKeyAmount:   fee,
			types.TxValueKeyTo:       verifierContractAddr,
			types.TxValueKeyData:     data,
		})
		if err != nil {
			log.Fatalf("Failed to encode tx: %v", err)
		}

		err = tx.SignWithKeys(signer, self.privateKey)
		if err != nil {
			log.Fatalf("Failed to sign tx: %v", err)
		}
		txs = append(txs, tx)
		indexes = append(indexes, i)
		self.nonce++
		nonce++

		// To save time, we can just test one bridge tx in a block
		//break
	}

	if len(txs) == 0 {
		fmt.Printf("found no bridge tx in block: %d\n", blockNum)
		return common.Hash{0}, nil
	}

	submitHeaderTx, _, err := c.TransactionByHash(ctx, submitHeaderTxHash)
	if err != nil {
		log.Fatalf("Cannot find tx with given submitHeader tx hash: %v\n", err)
		return common.Hash{0}, err
	}
	bind.WaitMined(ctx, c, submitHeaderTx)

	for index, tx := range txs {
		hash, err := c.SendRawTransaction(ctx, tx)
		if err != nil {
			if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
				fmt.Printf("submitReceipt Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
				fmt.Printf("submitReceipt Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
				self.nonce++
			} else {
				fmt.Printf("submitReceipt Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			}
			return common.Hash{0}, err
		}

		go func(curIndex int) {
			receipt, err := bind.WaitMined(ctx, c, tx)
			if err != nil {
				log.Fatalf("Failed to get receipt for erc20TransferVerify: %v", err)
			}
			gasUsed := new(big.Int).SetUint64(receipt.GasUsed)
			fmt.Printf("submitReceipt hash: 0x%v, status: %d, gasUsed: %v, tx count: %d bridge tx count: %d index: %d\n", hex.EncodeToString(hash[:]), receipt.Status, gasUsed.String(), count, bridgeTxCount, curIndex)
		}(indexes[index])
	}

	return blockHash, nil
}

func (self *Account) TransferGovUpdate(c *client.Client, endpoint string, verifierContractAddr common.Address, voteKey, voteVal string, toChainId uint) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	GlobalGovMutex.Lock()
	defer GlobalGovMutex.Unlock()

	nonce := self.GetNonce(c)
	gaslimit := uint64(4100000000)

	// Get blocknumber
	curBlockNum, err := c.BlockNumber(ctx)
	if err != nil {
		return common.Hash{0}, big.NewInt(0), err
	}
	header, err := c.HeaderByNumber(ctx, curBlockNum)
	if err != nil {
		log.Fatalf("failed to get header: %v\n", err)
	}
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		log.Fatalf("failed to get istanbulextra: %v\n", err)
	}
	type request struct {
		To   string `json:"to"`
		Data string `json:"data"`
	}

	client := &http.Client{}
	d := strings.NewReader(`{"jsonrpc":"2.0","method":"governance_chainConfigAt","params":[],"id":83}`)
	req, err := http.NewRequest("POST", endpoint, d)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	fmt.Printf("%s\n", bodyText)
	d = strings.NewReader(`{"jsonrpc":"2.0","method":"governance_vote","params":["` + voteKey + `", "` + voteVal + `"],"id":83}`)
	req, err = http.NewRequest("POST", endpoint, d)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)

	time.Sleep(3 * time.Second)
	epoch := int64(1)
	afterBlockNum, err := c.BlockNumber(ctx)
	if err != nil {
		return common.Hash{0}, big.NewInt(0), err
	}

	//startBlockNum := big.NewInt(beforeBlockNum.Uint64())
	fmt.Printf("Start finding gov header from block: %d\n", curBlockNum.Int64())
	limitBlockNum := big.NewInt(afterBlockNum.Int64() + epoch*2 + 10)
	for {
		if curBlockNum.Cmp(limitBlockNum) == 1 {
			log.Fatalf("Cannot find gov header until blockNum: %d\n", curBlockNum.Int64())
			return common.Hash{0}, common.Big0, fmt.Errorf("could not find gov header")
		}
		header, err = c.HeaderByNumber(ctx, curBlockNum)
		if err != nil {
			//log.Fatalf("Failed to find header with num: %v", err)
			//return common.Hash{0}, common.Big0, err
			time.Sleep(1 * time.Second)
			continue
		}
		//if len(header.Governance) > 2 {
		if len(header.Vote) > 2 {
			break
		}
		curBlockNum.Add(curBlockNum, big.NewInt(1))
	}
	encodedHeader := encodeHeader(header)

	// verifier
	abiStr := `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"requestTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"handleTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"from","type":"address"},{"indexed":false,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"bytes32","name":"reqHeaderHash","type":"bytes32"},{"indexed":false,"internalType":"bytes","name":"key","type":"bytes"}],"name":"BridgeERC20TransferHandle","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":true,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"address","name":"burnedTokenAddr","type":"address"},{"indexed":false,"internalType":"address","name":"ctTokenAddr","type":"address"}],"name":"BridgeERC20TransferRequest","type":"event"},{"anonymous":false,"inputs":[{"components":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"indexed":true,"internalType":"struct IManager.FeeSchedule","name":"feeSchedule","type":"tuple"}],"name":"CancelScheduledFee","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"feeReceiver","type":"address"}],"name":"ChangedFeeReceiver","type":"event"},{"anonymous":false,"inputs":[{"components":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"indexed":true,"internalType":"struct IManager.FeeSchedule","name":"feeSchedule","type":"tuple"}],"name":"FeeScheduled","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"blockNumber","type":"uint256"}],"name":"Frozen","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"headerNumber","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"requiredQuorum","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"committeeSize","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"validatorSize","type":"uint256"}],"name":"HeaderIntegrityPassed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"feeReceiver","type":"address"},{"indexed":false,"internalType":"uint256","name":"bridgeFee","type":"uint256"},{"indexed":false,"internalType":"address","name":"headerContributor","type":"address"},{"indexed":false,"internalType":"uint256","name":"govContributionFee","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"headerContributionFee","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"verifierBalance","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"userSent","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"change","type":"uint256"}],"name":"PayBridgeFee","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"onChainToken","type":"address"},{"indexed":true,"internalType":"address","name":"counterpartChainToken","type":"address"},{"indexed":true,"internalType":"bytes32","name":"transferEvent","type":"bytes32"}],"name":"TokenPairRegistered","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"governanceContributor","type":"address"},{"indexed":false,"internalType":"uint256","name":"governanceContributionFee","type":"uint256"}],"name":"TransferGovContrib","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"blockNumber","type":"uint256"}],"name":"Unfrozen","type":"event"},{"inputs":[],"name":"bridgeLibAddr","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"cancelFeeSchedule","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ctToken2TransferEvent","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ctToken2ocToken","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"address","name":"ocTokenAddr","type":"address"},{"internalType":"address","name":"ctTokenAddr","type":"address"}],"name":"erc20TransferRequest","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"feeReceiver","outputs":[{"internalType":"address payable","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"feeSchedule","outputs":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"freeze","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getStampBlockNumbers","outputs":[{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"}],"name":"getTransferReceipt","outputs":[{"components":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"address","name":"tokenAddr","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"internalType":"struct Typs.BridgeTransferStamp[]","name":"","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getVersion","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"govStateAddr","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_bridgeLibAddr","type":"address"},{"internalType":"address","name":"_govStateAddr","type":"address"},{"internalType":"bool","name":"_mintable","type":"bool"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"isVerifiedHeader","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"mintable","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"ocToken2ctToken","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"prevFeeSchdule","outputs":[{"internalType":"uint256","name":"curBridgeFee","type":"uint256"},{"internalType":"uint256","name":"curHdContribFee","type":"uint256"},{"internalType":"uint256","name":"curGovContribFee","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"},{"internalType":"uint256","name":"enrolledAt","type":"uint256"},{"internalType":"uint256","name":"effectedAt","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"effectAt","type":"uint256"},{"internalType":"uint256","name":"futureBridgeFee","type":"uint256"},{"internalType":"uint256","name":"futureHdContribFee","type":"uint256"},{"internalType":"uint256","name":"futureGovContribFee","type":"uint256"}],"name":"registerFeeSchedule","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"onChainToken","type":"address"},{"internalType":"address","name":"counterpartChainToken","type":"address"},{"internalType":"bytes32","name":"transferEventSig","type":"bytes32"}],"name":"registerTokenPair","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"verifier","type":"address"}],"name":"registerVerifier","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address payable","name":"_feeReceiver","type":"address"}],"name":"setFeeReceiver","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"headerHash","type":"bytes32"},{"internalType":"bytes","name":"rlpEncodedReceipt","type":"bytes"},{"internalType":"bytes[]","name":"MPTProof","type":"bytes[]"},{"internalType":"bytes","name":"key","type":"bytes"}],"name":"submitGov","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"rlpEncodedHeader","type":"bytes"}],"name":"submitHeader","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"headerHash","type":"bytes32"},{"internalType":"bytes","name":"rlpEncodedReceipt","type":"bytes"},{"internalType":"bytes[]","name":"MPTProof","type":"bytes[]"},{"internalType":"bytes","name":"key","type":"bytes"}],"name":"submitReceipt","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"temporalStop","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"transferStamps","outputs":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"address","name":"tokenAddr","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"unfreeze","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"verifiedHeaders","outputs":[{"internalType":"uint256","name":"headerNumber","type":"uint256"},{"internalType":"bytes32","name":"receiptHash","type":"bytes32"},{"internalType":"bytes","name":"vote","type":"bytes"},{"internalType":"bytes","name":"governance","type":"bytes"},{"internalType":"address","name":"contributor","type":"address"},{"internalType":"bool","name":"stored","type":"bool"}],"stateMutability":"view","type":"function"}]`

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}

	data, err := abii.Pack("submitHeader", encodedHeader)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: uint64(5000000000000),
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   common.Big0,
		types.TxValueKeyTo:       verifierContractAddr,
		types.TxValueKeyData:     data,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("GovUpdate 1 Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("GovUpdate 1 Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("GovUpdate 1 Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++
	nonce++

	bind.WaitMined(ctx, c, tx)
	receipt, err := c.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		log.Fatalf("Failed to get receipt: %v", err)
		return hash, gasPrice, err
	}
	gasUsed := new(big.Int).SetUint64(receipt.GasUsed)
	fmt.Printf("submitHeader hash: 0x%v, status: %d, gasUsed: %v, validator len: %d\n", hex.EncodeToString(hash[:]), receipt.Status, gasUsed.String(), len(istanbulExtra.Validators))

	data, err = abii.Pack("submitGov", header.Hash(), []byte(""), [][]byte{}, []byte(""))
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer = types.NewEIP155Signer(chainID)
	tx, err = types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: gaslimit,
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   common.Big0,
		types.TxValueKeyTo:       SmartContractAddr,
		types.TxValueKeyData:     data,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err = c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("GovUpdate 2 Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("GovUpdate 2 Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("GovUpdate 2 Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return common.Hash{0}, big.NewInt(0), err
	}

	self.nonce++
	nonce++

	bind.WaitMined(ctx, c, tx)

	receipt, err = c.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		log.Fatalf("Failed to get receipt: %v", err)
		return hash, gasPrice, err
	}
	gasUsed = new(big.Int).SetUint64(receipt.GasUsed)
	fmt.Printf("submitGov hash: 0x%v, status: %d, gasUsed: %v, validator len: %d\n", hex.EncodeToString(hash[:]), receipt.Status, gasUsed.String(), len(istanbulExtra.Validators))
	return hash, gasPrice, nil
}

func (self *Account) TransferNewSmartContractExecutionTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	abiStr := `[{"constant":true,"inputs":[],"name":"totalAmount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"receiver","type":"address"}],"name":"reward","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"safeWithdrawal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"}]`

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}

	data, err := abii.Pack("reward", self.address)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: uint64(5000000),
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   value,
		types.TxValueKeyTo:       to.address,
		types.TxValueKeyData:     data,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedSmartContractExecutionTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	abiStr := `[{"constant":true,"inputs":[],"name":"totalAmount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"receiver","type":"address"}],"name":"reward","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"safeWithdrawal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"}]`

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}

	data, err := abii.Pack("reward", self.address)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: uint64(5000000),
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   value,
		types.TxValueKeyTo:       to.address,
		types.TxValueKeyData:     data,
		types.TxValueKeyFeePayer: self.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = tx.SignFeePayerWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to fee payer sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedSmartContractExecutionWithRatioTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)
	abiStr := `[{"constant":true,"inputs":[],"name":"totalAmount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"receiver","type":"address"}],"name":"reward","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"safeWithdrawal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"}]`

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}

	data, err := abii.Pack("reward", self.address)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecutionWithRatio, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:              nonce,
		types.TxValueKeyGasPrice:           gasPrice,
		types.TxValueKeyGasLimit:           uint64(5000000),
		types.TxValueKeyFrom:               self.address,
		types.TxValueKeyAmount:             value,
		types.TxValueKeyTo:                 to.address,
		types.TxValueKeyData:               data,
		types.TxValueKeyFeePayer:           self.address,
		types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = tx.SignFeePayerWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to fee payer sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewCancelTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeCancel, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyGasLimit: uint64(100000000),
		types.TxValueKeyGasPrice: gasPrice,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedCancelTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancel, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyGasLimit: uint64(100000000),
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyFeePayer: to.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = tx.SignFeePayerWithKeys(signer, to.privateKey)
	if err != nil {
		log.Fatalf("Failed to fee payer sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedCancelWithRatioTx(c *client.Client, to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancelWithRatio, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:              nonce,
		types.TxValueKeyFrom:               self.address,
		types.TxValueKeyGasLimit:           uint64(100000000),
		types.TxValueKeyGasPrice:           gasPrice,
		types.TxValueKeyFeePayer:           to.address,
		types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = tx.SignFeePayerWithKeys(signer, to.privateKey)
	if err != nil {
		log.Fatalf("Failed to fee payer sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewEthereumAccessListTx(c *client.Client, to *Account, value *big.Int, input []byte) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	gas := uint64(5000000)

	var toAddress *common.Address
	if to != nil {
		toAddress = &to.address
	}
	callMsg := klaytn.CallMsg{
		From:     self.address,
		To:       toAddress,
		Gas:      gas,
		GasPrice: gasPrice,
		Value:    value,
		Data:     input,
	}
	accessList, _, _, err := c.CreateAccessList(ctx, callMsg)
	if err != nil {
		log.Fatalf("Failed to get accessList: %v", err)
	}

	signer := types.LatestSignerForChainID(chainID)

	tx := types.NewTx(&types.TxInternalDataEthereumAccessList{
		ChainID:      chainID,
		AccountNonce: nonce,
		Recipient:    toAddress,
		GasLimit:     gas,
		Price:        gasPrice,
		Amount:       value,
		AccessList:   *accessList,
		Payload:      input,
	})

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewEthereumDynamicFeeTx(c *client.Client, to *Account, value *big.Int, input []byte) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	gas := uint64(5000000)

	var toAddress *common.Address
	if to != nil {
		toAddress = &to.address
	}
	callMsg := klaytn.CallMsg{
		From:     self.address,
		To:       toAddress,
		Gas:      gas,
		GasPrice: gasPrice,
		Value:    value,
		Data:     input,
	}
	accessList, _, _, err := c.CreateAccessList(ctx, callMsg)
	if err != nil {
		log.Fatalf("Failed to get accessList: %v", err)
	}

	signer := types.LatestSignerForChainID(chainID)

	tx := types.NewTx(&types.TxInternalDataEthereumDynamicFee{
		ChainID:      chainID,
		AccountNonce: nonce,
		Recipient:    toAddress,
		GasLimit:     gas,
		GasFeeCap:    gasPrice,
		GasTipCap:    gasPrice,
		Amount:       value,
		AccessList:   *accessList,
		Payload:      input,
	})

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := c.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, gasPrice, err
	}

	self.nonce++

	return hash, gasPrice, nil
}

func (self *Account) TransferNewLegacyTxWithEth(c *client.Client, endpoint string, to *Account, value *big.Int, input string, exePath string) (common.Hash, *big.Int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	// Ethereum LegacyTx
	txType := "0"
	gas := "100000"

	var toAddress string
	if to != nil {
		toAddress = to.GetAddress().String()
	} else {
		// When to is nil, smart contract deployment with legacyTx case.
		// To send as a command argument which has to be string type,
		// explicitly send "nil" string for deploying.
		toAddress = "nil"
		gas = "200000"
	}

	// To test this, you need to update submodule and build executable file.
	// ./ethTxGenerator endPoint txType chainID gasPrice gas baseFee value fromPrivateKey nonce to [data]
	cmd := exec.Command(exePath, endpoint, txType, chainID.String(), gasPrice.String(), gas, baseFee.String(), value.String(), self.GetPrivateKey(), strconv.FormatUint(nonce, 10), toAddress, input)
	result, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to create and send tx : %v", err)
	}

	strResult := string(result[:])
	// Executable file will return transaction hash or error string.
	// So if result does not include "0x" prefix, means something went wrong.
	if !strings.Contains(strResult, "0x") {
		err = errors.New(strResult)
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return common.Hash{0}, gasPrice, err
	}

	self.nonce++

	return common.HexToHash(strResult), gasPrice, nil
}

func (self *Account) TransferNewEthAccessListTxWithEth(c *client.Client, endpoint string, to *Account, value *big.Int, input string, exePath string) (common.Hash, *big.Int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	// Ethereum AccessListTx
	txType := "1"
	gas := "100000"

	var toAddress string
	if to != nil {
		toAddress = to.GetAddress().String()
	} else {
		// When to is nil, smart contract deployment with legacyTx case.
		// To send as a command argument which has to be string type,
		// explicitly send "nil" string for deploying.
		toAddress = "nil"
		gas = "200000"
	}

	// To test this, you need to update submodule and build executable file.
	// ./ethTxGenerator endPoint txType chainID gasPrice gas baseFee value fromPrivateKey nonce to [data]
	cmd := exec.Command(exePath, endpoint, txType, chainID.String(), gasPrice.String(), gas, baseFee.String(), value.String(), self.GetPrivateKey(), strconv.FormatUint(nonce, 10), toAddress, input)
	result, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to create and send tx : %v", err)
	}

	strResult := string(result[:])
	// Executable file will return transaction hash or error string.
	// So if result does not include "0x" prefix, means something went wrong.
	if !strings.Contains(strResult, "0x") {
		err = errors.New(strResult)
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return common.Hash{0}, gasPrice, err
	}

	self.nonce++

	return common.HexToHash(strResult), gasPrice, nil
}

func (self *Account) TransferNewEthDynamicFeeTxWithEth(c *client.Client, endpoint string, to *Account, value *big.Int, input string, exePath string) (common.Hash, *big.Int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	// Ethereum DynamicFeeTx
	txType := "2"
	gas := "100000"

	var toAddress string
	if to != nil {
		toAddress = to.GetAddress().String()
	} else {
		// When to is nil, smart contract deployment with legacyTx case.
		// To send as a command argument which has to be string type,
		// explicitly send "nil" string for deploying.
		toAddress = "nil"
		gas = "200000"
	}

	// To test this, you need to update submodule and build executable file.
	// ./ethTxGenerator endPoint txType chainID gasPrice gas baseFee value fromPrivateKey nonce to [data]
	cmd := exec.Command(exePath, endpoint, txType, chainID.String(), gasPrice.String(), gas, baseFee.String(), value.String(), self.GetPrivateKey(), strconv.FormatUint(nonce, 10), toAddress, input)
	result, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("fromAddress: %v, strconv.FormatUint(nonce, 10): %v, to: %v input: %v gas: %v \n", self.GetAddress().String(), strconv.FormatUint(nonce, 10), toAddress, input, gas)
		log.Fatalf("Failed to create and send tx : %v", err)
	}

	strResult := string(result[:])
	// Executable file will return transaction hash or error string.
	// So if result does not include "0x" prefix, means something went wrong.
	if !strings.Contains(strResult, "0x") {
		err = errors.New(strResult)
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTransaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return common.Hash{0}, gasPrice, err
	}

	self.nonce++

	return common.HexToHash(strResult), gasPrice, nil
}

func (self *Account) TransferUnsignedTx(c *client.Client, to *Account, value *big.Int) (common.Hash, error) {
	ctx := context.Background()

	fromAddr := self.GetAddress()
	toAddr := to.GetAddress()

	data := hexutil.Bytes{}
	input := hexutil.Bytes{}

	var err error
	hash, err := c.SendUnsignedTransaction(ctx, fromAddr, toAddr, 21000, gasPrice.Uint64(), value, data, input)
	if err != nil {
		log.Printf("Account(%v) : Failed to sendTransaction: %v\n", self.address[:5], err)
		return common.Hash{}, err
	}
	//log.Printf("Account(%v) : Success to sendTransaction: %v\n", self.address[:5], hash.String())
	return hash, nil
}

func TransferUnsignedTx(c *client.Client, from common.Address, to common.Address, value *big.Int) (common.Hash, error) {
	ctx := context.Background()

	data := hexutil.Bytes{}
	input := hexutil.Bytes{}

	var err error
	hash, err := c.SendUnsignedTransaction(ctx, from, to, 21000, gasPrice.Uint64(), value, data, input)
	if err != nil {
		log.Printf("Account(%v) : Failed to sendTransaction: %v\n", from[:5], err)
		return common.Hash{}, err
	}

	return hash, nil
}

func (a *Account) CheckBalance(expectedBalance *big.Int, cli *client.Client) error {
	balance, _ := a.GetBalance(cli)
	if balance.Cmp(expectedBalance) != 0 {
		fmt.Println(a.address.String() + " expected : " + expectedBalance.Text(10) + " actual : " + balance.Text(10))
		return errors.New("expected : " + expectedBalance.Text(10) + " actual : " + balance.Text(10))
	}

	return nil
}
