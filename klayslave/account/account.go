package account

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/contracts/bridge"
	"github.com/klaytn/klaytn/contracts/sc_erc20"
	"github.com/klaytn/klaytn/contracts/sc_erc721"
	"github.com/klaytn/klaytn/crypto"
	"log"
	"math/big"
	"os"
	"sync"
	"time"
	"unsafe"
)

const (
	MagicGasLimit = 9999999999
)

type Backend struct {
	*client.Client
	chainID  *big.Int
	gasPrice *big.Int
}

func LockAccounts(aAcc, bAcc *Account) {
	// This pointer comparing code is for double lock without deadlock.
	if uintptr(unsafe.Pointer(aAcc)) >= uintptr(unsafe.Pointer(bAcc)) {
		aAcc.Lock()
		bAcc.Lock()
	} else {
		bAcc.Lock()
		aAcc.Lock()
	}
}

func UnlockAccounts(aAcc, bAcc *Account) {
	if uintptr(unsafe.Pointer(aAcc)) >= uintptr(unsafe.Pointer(bAcc)) {
		bAcc.UnLock()
		aAcc.UnLock()
	} else {
		aAcc.UnLock()
		bAcc.UnLock()
	}
}

type Account struct {
	privateKey []*ecdsa.PrivateKey
	key        []string
	address    common.Address
	nonce      uint64
	balance    *big.Int
	mutex      sync.Mutex

	backend  *Backend
	tokenBal *big.Int        // TODO-Klaytn need to support mutiple token.  map[common.Address]*big.Int
	nftBal   map[uint64]bool // TODO-Klaytn need to consider uint256 token id.

	ctAccount *Account // counter part account is on counter part chain (backend) with same address and key.
}

func (acc *Account) Backend() *Backend {
	return acc.backend
}

func (acc *Account) Lock() {
	acc.mutex.Lock()
}

func (acc *Account) UnLock() {
	acc.mutex.Unlock()
}

func (acc *Account) CounterPartAccount() *Account {
	return acc.ctAccount
}

func NewBackend(ep string) *Backend {
	cli, err := client.Dial(ep)
	if err != nil {
		log.Fatalf("Failed to connect RPC: %v", err)
	}
	ctx := context.Background()
	chainID, err := cli.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get ChainID: %v", err)
	}

	gasPrice, err := cli.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("Failed to get SuggestGasPrice: %v", err)
	}

	backend := &Backend{
		cli,
		chainID,
		gasPrice,
	}

	return backend
}

func RequestTokenTransfer(from *Account, to *Account, value *big.Int, targetToken common.Address) (*types.Transaction, error) {
	if from.tokenBal.Cmp(value) < 0 {
		log.Println("Not enough ERC20 balance of the from", "balance", from.tokenBal.String())
		return nil, errors.New("not enough ERC20 balance")
	}
	token, err := sctoken.NewServiceChainToken(targetToken, from.backend.Client)
	if err != nil {
		log.Println("Failed to get ERC20 object", "err", err)
		return nil, err
	}

	auth := from.GetTransactOpts(MagicGasLimit)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	auth.Context = ctx

	tx, err := token.RequestValueTransfer(auth, value, to.address, common.Big0, nil)
	if err != nil {
		log.Println("Failed to RequestValueTransfer of Token", "err", err)
		return nil, err
	}

	//fmt.Printf("Success to RequestValueTransfer ERC20. txhash(%v)\n", tx.Hash().String())

	from.SubTokenBalance(value)
	defer to.AddTokenBalance(value)

	from.nonce++
	return tx, nil
}

func RequestTokenTransfer2Step(from *Account, to *Account, value *big.Int, erc20Addr common.Address, bridgeAddr common.Address) (*types.Transaction, error) {
	if from.tokenBal.Cmp(value) < 0 {
		log.Println("Not enough ERC20 balance of the from", "balance", from.tokenBal.String())
		return nil, errors.New("not enough ERC20 balance")
	}
	erc20, err := sctoken.NewServiceChainToken(erc20Addr, from.backend.Client)
	if err != nil {
		log.Println("Failed to get ERC20 object", "err", err)
		return nil, err
	}

	b, err := bridge.NewBridge(bridgeAddr, from.backend)
	if err != nil {
		log.Println("Failed to get bridge object", "err", err)
		return nil, err
	}

	auth := from.GetTransactOpts(MagicGasLimit)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	auth.Context = ctx

	tx, err := erc20.Approve(auth, bridgeAddr, value)
	if err != nil {
		log.Println("Failed to Approve ERC20", "err", err)
		return nil, err
	}
	from.nonce++

	auth = from.GetTransactOpts(MagicGasLimit)
	tx, err = b.RequestERC20Transfer(auth, erc20Addr, to.address, value, common.Big0, nil)
	if err != nil {
		log.Println("Failed to RequestValueTransfer of ERC20", "err", err)
		return nil, err
	}
	from.nonce++

	//fmt.Printf("Success to RequestValueTransfer2Step ERC20. txhash(%v)\n", tx.Hash().String())

	from.SubTokenBalance(value)
	defer to.AddTokenBalance(value)

	return tx, nil
}

func RequestNFTTransfer(from *Account, to *Account, targetToken common.Address) (*types.Transaction, error) {
	uid := from.GetAnyNFT()

	if !from.isOwnNFT(uid) {
		log.Println("not own the NFT", "from", from.address.String(), "uid", uid.String(), "balance", from.NFTBalance())
		return nil, errors.New("not own NFT")
	}

	nft, err := scnft.NewServiceChainNFT(targetToken, from.backend.Client)
	if err != nil {
		log.Println("Failed to get ERC721 object", "err", err)
		return nil, err
	}

	auth := from.GetTransactOpts(MagicGasLimit)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	auth.Context = ctx

	tx, err := nft.RequestValueTransfer(auth, uid, to.address, nil)
	if err != nil {
		log.Println("Failed to RequestValueTransfer of Token", "err", err)
		return nil, err
	}

	//fmt.Printf("Success to RequestValueTransfer ERC721. txhash(%v)\n", tx.Hash().String())

	from.SubNFT(uid)
	defer to.AddNFT(uid)

	from.nonce++
	return tx, nil
}

func RequestNFTTransfer2Step(from *Account, to *Account, erc721Addr common.Address, bridgeAddr common.Address) (*types.Transaction, error) {
	uid := from.GetAnyNFT()

	if !from.isOwnNFT(uid) {
		log.Println("not own the ERC721", "from", from.address.String(), "uid", uid.String(), "balance", from.NFTBalance())
		return nil, errors.New("not own NFT")
	}

	erc721, err := scnft.NewServiceChainNFT(erc721Addr, from.backend.Client)
	if err != nil {
		log.Println("Failed to get ERC721 object", "err", err)
		return nil, err
	}

	b, err := bridge.NewBridge(bridgeAddr, from.backend)
	if err != nil {
		log.Println("Failed to get bridge object", "err", err)
		return nil, err
	}

	auth := from.GetTransactOpts(MagicGasLimit)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	auth.Context = ctx

	_, err = erc721.Approve(auth, bridgeAddr, uid)
	if err != nil {
		log.Println("Failed to Approve of ERC721", "err", err)
		return nil, err
	}
	from.nonce++

	auth = from.GetTransactOpts(MagicGasLimit)
	tx, err := b.RequestERC721Transfer(auth, erc721Addr, to.address, uid, nil)
	if err != nil {
		log.Println("Failed to RequestValueTransfer of ERC721", "err", err)
		return nil, err
	}
	from.nonce++

	//fmt.Printf("Success to RequestValueTransfer2Step ERC721. txhash(%v)\n", tx.Hash().String())

	from.SubNFT(uid)
	defer to.AddNFT(uid)

	return tx, nil
}

func RequestKlayTransfer(from *Account, to *Account, value *big.Int, targetBridge common.Address) error {
	_, err := RequestKlayTransferReturnTx(from, to, value, targetBridge, false)
	return err
}

func RequestKlayTransferReturnTx(from *Account, to *Account, value *big.Int, targetBridge common.Address, withCheck bool) (*types.Transaction, error) {
	if from.balance.Cmp(value) < 0 {
		log.Println("Not enough KLAY balance of the from")
		return nil, errors.New("not enough KLAY balance")
	}

	bridgeObj, err := bridge.NewBridge(targetBridge, from.backend)
	if err != nil {
		log.Println("failed to get bridge obj", err)
	}

	auth := bind.NewKeyedTransactor(from.privateKey[0])
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	auth.Context = ctx
	auth.GasPrice = from.backend.gasPrice
	auth.GasLimit = MagicGasLimit
	auth.Nonce = big.NewInt(int64(from.GetNonce()))
	auth.Value = value
	auth.Signer = func(signer types.Signer, addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return types.SignTx(tx, types.NewEIP155Signer(from.backend.chainID), from.privateKey[0])
	}

	tx, err := bridgeObj.RequestKLAYTransfer(auth, to.address, value, nil)
	if err != nil {
		log.Println("failed to request transferKLAY", "err", err)
		return nil, err
	}
	//fmt.Printf("Success to RequestKLAYTransfer. txhash(%v)\n", tx.Hash().String())

	if withCheck {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		receipt, err := bind.WaitMined(ctx, from.backend.Client, tx)
		if err != nil {
			log.Println("WaitMined time out %v", err)
			return nil, err
		}
		fee := big.NewInt(0)
		fee.Mul(big.NewInt(int64(receipt.GasUsed)), from.backend.gasPrice)
		from.SubBalance(fee)
	}

	from.SubBalance(value)
	defer to.AddBalance(value)

	from.nonce++

	return tx, nil
}

func RequestKlayTransferFallbackReturnTx(from *Account, to *Account, value *big.Int, bridgeAddr common.Address, withCheck bool) (*types.Transaction, error) {
	if from.balance.Cmp(value) < 0 {
		log.Println("Not enough KLAY balance of the from")
		return nil, errors.New("not enough KLAY balance")
	}

	tx, _, err := from.TransferSignedTxToAddrWithoutLock(bridgeAddr, value)
	if err != nil {
		log.Println("failed to request klayTransferFallback", "err", err)
		return nil, err
	}
	//fmt.Printf("Success to klayTransferFallback. txhash(%v)\n", tx.Hash().String())

	if withCheck {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		receipt, err := bind.WaitMined(ctx, from.backend.Client, tx)
		if err != nil {
			log.Println("WaitMined time out %v", err)
			return nil, err
		}
		fee := big.NewInt(0)
		fee.Mul(big.NewInt(int64(receipt.GasUsed)), from.backend.gasPrice)
		from.SubBalance(fee)
	}

	from.SubBalance(value)
	defer to.AddBalance(value)

	return tx, nil
}

// GetTransactOpts generates the transactOpts from the account with the given gasLimit.
// If given gasLimit is zero, transactor will calculate the gasLimit by client call, or it will use the gasLimit.
func (self *Account) GetTransactOpts(gasLimit uint64) *bind.TransactOpts {
	opts := bind.NewKeyedTransactor(self.privateKey[0])
	opts.Nonce = new(big.Int).SetUint64(self.nonce)
	opts.GasLimit = gasLimit
	opts.GasPrice = self.backend.gasPrice
	opts.Signer = func(signer types.Signer, addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return types.SignTx(tx, types.NewEIP155Signer(self.backend.chainID), self.privateKey[0])
	}
	return opts
}

func (self *Account) AddBalance(value *big.Int) *big.Int {
	self.balance.Add(self.balance, value)
	return self.balance
}

func (self *Account) SubBalance(value *big.Int) *big.Int {
	self.balance.Sub(self.balance, value)
	return self.balance
}

func (self *Account) Balance() *big.Int {
	return self.balance
}

func (self *Account) AddTokenBalance(value *big.Int) *big.Int {
	self.tokenBal.Add(self.tokenBal, value)
	return self.tokenBal
}

func (self *Account) SubTokenBalance(value *big.Int) *big.Int {
	self.tokenBal.Sub(self.tokenBal, value)
	return self.tokenBal
}

func (self *Account) TokenBalance() *big.Int {
	return self.tokenBal
}

func (self *Account) AddNFT(uid *big.Int) *big.Int {
	self.nftBal[uid.Uint64()] = true
	return uid
}

func (self *Account) SubNFT(uid *big.Int) *big.Int {
	delete(self.nftBal, uid.Uint64())
	return uid
}

func (self *Account) NFTBalance() int {
	return len(self.nftBal)
}

func (self *Account) GetAnyNFT() *big.Int {
	if len(self.nftBal) > 0 {
		for id, _ := range self.nftBal {
			return new(big.Int).SetUint64(id)
		}
	}
	return nil
}

func (self *Account) isOwnNFT(uid *big.Int) bool {
	_, exist := self.nftBal[uid.Uint64()]
	if exist {
		return true
	}
	return false
}

func GetAccountFromKey(key string, backend *Backend) *Account {
	acc, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Fatalf("Key(%v): Failed to HexToECDSA %v", key, err)
	}

	tAcc := Account{
		[]*ecdsa.PrivateKey{acc},
		[]string{key},
		crypto.PubkeyToAddress(acc.PublicKey),
		0,
		big.NewInt(0),
		sync.Mutex{},
		//make(TransactionMap),
		backend,
		big.NewInt(0), //make(map[common.Address]*big.Int),
		make(map[uint64]bool, 100),
		nil,
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

	addr, err := account.backend.Client.ImportRawKey(ctx, key, "")
	if err != nil {
		log.Fatalf("Account(%v) : Failed to import => %v\n", account.address, err)
	} else {
		if testAddr != addr {
			log.Fatalf("origial:%v, imported: %v\n", testAddr.String(), addr.String())
		}
	}

	res, err := account.backend.Client.UnlockAccount(ctx, account.address, "", 0)
	if err != nil {
		log.Fatalf("Account(%v) : Failed to Unlock: %v\n", account.address.String(), err)
	} else {
		log.Printf("Wallet UnLock Result: %v", res)
	}
}

func NewPairAccount(cBack, pBack *Backend) (*Account, *Account) {
	acc, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}

	testKey := hex.EncodeToString(crypto.FromECDSA(acc))

	pAcc := Account{
		[]*ecdsa.PrivateKey{acc},
		[]string{testKey},
		crypto.PubkeyToAddress(acc.PublicKey),
		0,
		big.NewInt(0),
		sync.Mutex{},
		//make(TransactionMap),
		pBack,
		big.NewInt(0), //make(map[common.Address]*big.Int),
		make(map[uint64]bool, 100),
		nil,
	}

	cAcc := Account{
		[]*ecdsa.PrivateKey{acc},
		[]string{testKey},
		crypto.PubkeyToAddress(acc.PublicKey),
		0,
		big.NewInt(0),
		sync.Mutex{},
		//make(TransactionMap),
		cBack,
		big.NewInt(0), //make(map[common.Address]*big.Int),
		make(map[uint64]bool, 100),
		nil,
	}

	pAcc.ctAccount = &cAcc
	cAcc.ctAccount = &pAcc

	return &cAcc, &pAcc
}

func NewAccount(backend *Backend) *Account {
	acc, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}

	testKey := hex.EncodeToString(crypto.FromECDSA(acc))

	tAcc := Account{
		[]*ecdsa.PrivateKey{acc},
		[]string{testKey},
		crypto.PubkeyToAddress(acc.PublicKey),
		0,
		big.NewInt(0),
		sync.Mutex{},
		//make(TransactionMap),
		backend,
		big.NewInt(0), //make(map[common.Address]*big.Int),
		make(map[uint64]bool, 100),
		nil,
	}

	return &tAcc
}

func NewAccountOnNode(backend *Backend) *Account {
	tAcc := NewAccount(backend)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addr, err := backend.ImportRawKey(ctx, tAcc.key[0], "")
	if err != nil {
		//log.Printf("Account(%v) : Failed to import\n", tAcc.address, err)
	} else {
		if tAcc.address != addr {
			log.Fatalf("origial:%v, imported: %v\n", tAcc.address, addr.String())
		}
		//log.Printf("origial:%v, imported:%v\n", tAcc.address, addr.String())
	}

	_, err = backend.UnlockAccount(ctx, tAcc.GetAddress(), "", 0)
	if err != nil {
		log.Printf("Account(%v) : Failed to Unlock: %v\n", tAcc.GetAddress().String(), err)
	}

	//log.Printf("Wallet UnLock Result: %v", flag)

	return tAcc
}

func NewKlaytnAccount(backend *Backend) *Account {
	acc, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}

	testKey := hex.EncodeToString(crypto.FromECDSA(acc))

	randomAddr := common.BytesToAddress(crypto.Keccak256([]byte(testKey))[12:])

	tAcc := Account{
		[]*ecdsa.PrivateKey{acc},
		[]string{testKey},
		randomAddr,
		0,
		big.NewInt(0),
		sync.Mutex{},
		//make(TransactionMap),
		backend,
		big.NewInt(0), //make(map[common.Address]*big.Int),
		make(map[uint64]bool, 100),
		nil,
	}

	return &tAcc
}

func NewKlaytnAccountWithAddr(addr common.Address, backend *Backend) *Account {
	acc, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("crypto.GenerateKey() : Failed to generateKey %v", err)
	}

	testKey := hex.EncodeToString(crypto.FromECDSA(acc))

	tAcc := Account{
		[]*ecdsa.PrivateKey{acc},
		[]string{testKey},
		addr,
		0,
		big.NewInt(0),
		sync.Mutex{},
		//make(TransactionMap),
		backend,
		big.NewInt(0), //make(map[common.Address]*big.Int),
		make(map[uint64]bool, 100),
		nil,
	}

	return &tAcc
}

func NewKlaytnMultisigAccount(backend *Backend) *Account {
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
		[]*ecdsa.PrivateKey{k1, k2, k3},
		[]string{testKey},
		randomAddr,
		0,
		big.NewInt(0),
		sync.Mutex{},
		//make(TransactionMap),
		backend,
		big.NewInt(0), //make(map[common.Address]*big.Int),
		make(map[uint64]bool, 100),
		nil,
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

func (acc *Account) GetNonce() uint64 {
	if acc.nonce != 0 {
		return acc.nonce
	}
	ctx := context.Background()
	nonce, err := acc.backend.Client.NonceAt(ctx, acc.GetAddress(), nil)
	if err != nil {
		log.Printf("GetNonce(): Failed to NonceAt() %v\n", err)
		return acc.nonce
	}
	acc.nonce = nonce

	//fmt.Printf("account= %v  nonce = %v\n", acc.GetAddress().String(), nonce)
	return acc.nonce
}

func (acc *Account) GetNonceFromBlock() uint64 {
	ctx := context.Background()
	nonce, err := acc.backend.Client.NonceAt(ctx, acc.GetAddress(), nil)
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

func (a *Account) GetReceipt(txHash common.Hash) (*types.Receipt, error) {
	ctx := context.Background()
	return a.backend.TransactionReceipt(ctx, txHash)
}

func (a *Account) GetBalance() (*big.Int, error) {
	ctx := context.Background()
	balance, err := a.backend.Client.BalanceAt(ctx, a.GetAddress(), nil)
	if err != nil {
		return nil, err
	}
	return balance, err
}

func (self *Account) ChargeBridge(bridgeAddr common.Address, value *big.Int) (common.Hash, *big.Int, error) {
	bridge, err := bridge.NewBridge(bridgeAddr, self.backend.Client)
	if err != nil {
		return common.Hash{}, nil, err
	}

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()
	gasPrice := self.backend.gasPrice

	auth := bind.NewKeyedTransactor(self.privateKey[0])
	auth.GasPrice = gasPrice
	auth.GasLimit = MagicGasLimit
	auth.Value = value
	auth.Nonce = big.NewInt(int64(nonce))
	tx, err := bridge.ChargeWithoutEvent(auth)
	if err != nil {
		return common.Hash{}, nil, err
	}

	fee := big.NewInt(0)
	//ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	//defer cancel()
	//receipt, err := bind.WaitMined(ctx, self.backend, tx)
	//if err != nil {
	//	log.Fatalf("WaitMined time out %v", err)
	//}

	//fee.Mul(big.NewInt(int64(receipt.GasUsed)), gasPrice)
	self.nonce++

	self.SubBalance(fee)
	self.SubBalance(value)

	fmt.Println("charging txHash=", tx.Hash().String())
	return tx.Hash(), fee, nil
}

func (self *Account) TransferTokenToAccount(tokenAddr common.Address, to *Account, value *big.Int) (*types.Transaction, error) {
	tx, err := self.TransferTokenToAddr(tokenAddr, to.address, value)
	if err != nil {
		log.Fatal("Failed to TransferTokenToAddr", "err", err)
		return nil, err
	}

	self.SubTokenBalance(value)
	to.AddTokenBalance(value)
	// TODO-Klaytn need to consider gas fee.

	self.UpdateNonce()
	return tx, nil
}

func (self *Account) TransferTokenToAddr(tokenAddr, to common.Address, value *big.Int) (*types.Transaction, error) {
	token, err := sctoken.NewServiceChainToken(tokenAddr, self.backend.Client)
	if err != nil {
		log.Fatal("Failed to get ERC20 object", "err", err)
	}

	tx, err := token.Transfer(self.GetTransactOpts(MagicGasLimit), to, value)

	return tx, err
}

func (self *Account) RegisterNFTToAccount(nftAddr common.Address, to *Account, start, end uint64) (*types.Transaction, error) {
	nft, err := scnft.NewServiceChainNFT(nftAddr, self.backend.Client)
	if err != nil {
		log.Fatal("Failed to get ERC721 object", "err", err)
		return nil, err
	}

	tx, err := nft.RegisterBulk(self.GetTransactOpts(MagicGasLimit), to.address, new(big.Int).SetUint64(start), new(big.Int).SetUint64(end))
	if err != nil {
		log.Fatal("Failed to RegisterBulk", "err", err)
		return nil, err
	}

	log.Printf("RegisterNFT, txhash=%v\n", tx.Hash().String())

	for i := start; i < end; i++ {
		to.AddNFT(new(big.Int).SetUint64(i))
	}

	self.UpdateNonce()
	return tx, nil
}

func (self *Account) TransferSignedTxToAddr(to common.Address, value *big.Int) (*types.Transaction, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	//fmt.Printf("account=%v, nonce = %v\n", self.GetAddress().String(), nonce)

	tx := types.NewTransaction(
		nonce,
		to,
		value,
		MagicGasLimit, //21000,
		self.backend.gasPrice,
		nil)
	gasPrice := tx.GasPrice()
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(self.backend.chainID), self.privateKey[0])
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	_, err = self.backend.Client.SendRawTransaction(ctx, signTx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return signTx, gasPrice, err
	}

	self.nonce++

	//fmt.Printf("%v transferSignedTx %v klay to %v klay.\n", self.GetAddress().Hex(), to.GetAddress().Hex(), value)

	return signTx, gasPrice, nil
}

func (self *Account) TransferSignedTxToAddrWithoutLock(to common.Address, value *big.Int) (*types.Transaction, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	nonce := self.GetNonce()

	//fmt.Printf("account=%v, nonce = %v\n", self.GetAddress().String(), nonce)

	tx := types.NewTransaction(
		nonce,
		to,
		value,
		MagicGasLimit, //21000,
		self.backend.gasPrice,
		nil)
	gasPrice := tx.GasPrice()
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(self.backend.chainID), self.privateKey[0])
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	_, err = self.backend.Client.SendRawTransaction(ctx, signTx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return signTx, gasPrice, err
	}

	self.nonce++

	//fmt.Printf("%v transferSignedTx %v klay to %v klay.\n", self.GetAddress().Hex(), to.GetAddress().Hex(), value)

	return signTx, gasPrice, nil
}

func (self *Account) TransferSignedTx(to *Account, value *big.Int) (*types.Transaction, *big.Int, error) {
	return self.TransferSignedTxToAddr(to.address, value)
}

func (self *Account) TransferNewValueTransferTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: self.backend.gasPrice,
		types.TxValueKeyFrom:     self.address,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedValueTransferTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransfer, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: self.backend.gasPrice,
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

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedValueTransferWithRatioTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferWithRatio, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:              nonce,
		types.TxValueKeyTo:                 to.GetAddress(),
		types.TxValueKeyAmount:             value,
		types.TxValueKeyGasLimit:           uint64(100000),
		types.TxValueKeyGasPrice:           self.backend.gasPrice,
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

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewValueTransferMemoTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()
	data := []byte("hello")

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransferMemo, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: self.backend.gasPrice,
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

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedValueTransferMemoTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()
	data := []byte("hello")

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferMemo, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyTo:       to.GetAddress(),
		types.TxValueKeyAmount:   value,
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: self.backend.gasPrice,
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

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedValueTransferMemoWithRatioTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()
	data := []byte("hello")

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferMemoWithRatio, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:              nonce,
		types.TxValueKeyTo:                 to.GetAddress(),
		types.TxValueKeyAmount:             value,
		types.TxValueKeyGasLimit:           uint64(100000),
		types.TxValueKeyGasPrice:           self.backend.gasPrice,
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

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewAccountCreationTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeAccountCreation, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:         nonce,
		types.TxValueKeyFrom:          self.address,
		types.TxValueKeyTo:            to.GetAddress(),
		types.TxValueKeyAmount:        value,
		types.TxValueKeyGasLimit:      uint64(1000000),
		types.TxValueKeyGasPrice:      self.backend.gasPrice,
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

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewAccountUpdateTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:      nonce,
		types.TxValueKeyFrom:       self.address,
		types.TxValueKeyGasLimit:   uint64(100000),
		types.TxValueKeyGasPrice:   self.backend.gasPrice,
		types.TxValueKeyAccountKey: accountkey.NewAccountKeyPublicWithValue(&self.privateKey[0].PublicKey),
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedAccountUpdateTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdate, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:      nonce,
		types.TxValueKeyFrom:       self.address,
		types.TxValueKeyGasLimit:   uint64(100000),
		types.TxValueKeyGasPrice:   self.backend.gasPrice,
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

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedAccountUpdateWithRatioTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdateWithRatio, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:              nonce,
		types.TxValueKeyFrom:               self.address,
		types.TxValueKeyGasLimit:           uint64(100000),
		types.TxValueKeyGasPrice:           self.backend.gasPrice,
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

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewSmartContractDeployTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	code := "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:         nonce,
		types.TxValueKeyFrom:          self.address,
		types.TxValueKeyTo:            to.address,
		types.TxValueKeyAmount:        common.Big0,
		types.TxValueKeyGasLimit:      uint64(10000000),
		types.TxValueKeyGasPrice:      self.backend.gasPrice,
		types.TxValueKeyHumanReadable: false,
		types.TxValueKeyData:          common.FromHex(code),
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedSmartContractDeployTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	code := "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeploy, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:         nonce,
		types.TxValueKeyFrom:          self.address,
		types.TxValueKeyTo:            to.address,
		types.TxValueKeyAmount:        common.Big0,
		types.TxValueKeyGasLimit:      uint64(10000000),
		types.TxValueKeyGasPrice:      self.backend.gasPrice,
		types.TxValueKeyHumanReadable: false,
		types.TxValueKeyData:          common.FromHex(code),
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

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedSmartContractDeployWithRatioTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	code := "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeployWithRatio, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:              nonce,
		types.TxValueKeyFrom:               self.address,
		types.TxValueKeyTo:                 to.address,
		types.TxValueKeyAmount:             common.Big0,
		types.TxValueKeyGasLimit:           uint64(10000000),
		types.TxValueKeyGasPrice:           self.backend.gasPrice,
		types.TxValueKeyHumanReadable:      false,
		types.TxValueKeyData:               common.FromHex(code),
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

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewCancelTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeCancel, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyGasLimit: uint64(100000000),
		types.TxValueKeyGasPrice: self.backend.gasPrice,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	err = tx.SignWithKeys(signer, self.privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedCancelTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancel, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyGasLimit: uint64(100000000),
		types.TxValueKeyGasPrice: self.backend.gasPrice,
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

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferNewFeeDelegatedCancelWithRatioTx(to *Account, value *big.Int) (common.Hash, *big.Int, error) {
	ctx := context.Background() //context.WithTimeout(context.Background(), 100*time.Second)

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce()

	signer := types.NewEIP155Signer(self.backend.chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancelWithRatio, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:              nonce,
		types.TxValueKeyFrom:               self.address,
		types.TxValueKeyGasLimit:           uint64(100000000),
		types.TxValueKeyGasPrice:           self.backend.gasPrice,
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

	hash, err := self.backend.Client.SendRawTransaction(ctx, tx)
	if err != nil {
		if err.Error() == blockchain.ErrNonceTooLow.Error() || err.Error() == blockchain.ErrReplaceUnderpriced.Error() {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
			fmt.Printf("Account(%v) nonce is added to %v\n", self.GetAddress().String(), nonce+1)
			self.nonce++
		} else {
			fmt.Printf("Account(%v) nonce(%v) : Failed to sendTrasnaction: %v\n", self.GetAddress().String(), nonce, err)
		}
		return hash, self.backend.gasPrice, err
	}

	self.nonce++

	return hash, self.backend.gasPrice, nil
}

func (self *Account) TransferUnsignedTx(to *Account, value *big.Int) (common.Hash, error) {
	ctx := context.Background()

	fromAddr := self.GetAddress()
	toAddr := to.GetAddress()

	data := hexutil.Bytes{}
	input := hexutil.Bytes{}

	var err error
	hash, err := self.backend.Client.SendUnsignedTransaction(ctx, fromAddr, toAddr, 21000, self.backend.gasPrice.Uint64(), value, data, input)
	if err != nil {
		log.Printf("Account(%v) : Failed to sendTrasnaction: %v\n", self.address[:5], err)
		return common.Hash{}, err
	}
	//log.Printf("Account(%v) : Success to sendTrasnaction: %v\n", self.address[:5], hash.String())
	return hash, nil
}

func (a *Account) CheckBalance(expectedBalance *big.Int) error {
	balance, _ := a.GetBalance()
	if balance.Cmp(expectedBalance) != 0 {
		fmt.Println(a.address.String() + " expected : " + expectedBalance.Text(10) + " actual : " + balance.Text(10))
		return errors.New("expected : " + expectedBalance.Text(10) + " actual : " + balance.Text(10))
	}

	return nil
}
