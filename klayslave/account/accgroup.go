package account

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/crypto"
)

type AccLoader func(*AccGroup)

type AccGroup struct {
	listAcc []*Account
}

func (a *AccGroup) Load(loader AccLoader) {
	loader(a)
}

func (a *AccGroup) Get(idx int) *Account {
	return a.listAcc[idx]
}

func (a *AccGroup) GetCount() int {
	return len(a.listAcc)
}

func (a *AccGroup) AddAcc(acc *Account) {
	a.listAcc = append(a.listAcc, acc)
}

func newAccGroup() *AccGroup {
	ag := &AccGroup{}
	return ag
}

type AccMgr struct {
	listAccGroup []*AccGroup
}

func (a *AccMgr) GetGroup(idx int) *AccGroup {
	return a.listAccGroup[idx]
}

// for temporary

func (a *AccMgr) ChargeAccounts(ctx context.Context, cli *client.Client, coinAcc *ecdsa.PrivateKey, acc *Account, nonce *uint64) {
	tx := types.NewTransaction(
		*nonce,
		acc.GetAddress(),
		big.NewInt(1e6),
		100000,
		gasPrice,
		nil)

	signTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), coinAcc)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	hash, err := cli.SendRawTransaction(ctx, signTx)
	if err != nil {
		log.Fatalf("Failed to charge account: %v", err)
	} else {
		*nonce++
	}
	fmt.Printf("%v is charged %v peb. tx_hash=%v\n", acc.GetAddress().Hex(), 1e6, hash.String())
}

func (a *AccMgr) BuildAccount(userCnt int, cnt int) {
	id := 0
	for i := 0; i < cnt; i++ {
		ag := newAccGroup()

		ld := func(g *AccGroup) {
			for j := 0; j < userCnt; j++ {
				id++
				//acc := NewAccount(int64(id))
				//// TODO charge balance
				////a.ChargeAccounts(ctx, cli, account, &nonce)
				//
				//g.AddAcc(acc)
			}
		}
		ag.Load(ld)
		a.listAccGroup = append(a.listAccGroup, ag)
	}
}

func NewAccountOnly(id int) *Account {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic("couldn't generate key: " + err.Error())
	}
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return &Account{id: id, privateKey: []*ecdsa.PrivateKey{key},
		address: addr}
}

var once sync.Once
var accMgrInst *AccMgr

func GetAccMgr() *AccMgr {
	once.Do(func() {
		accMgrInst = &AccMgr{}
	})
	return accMgrInst
}
