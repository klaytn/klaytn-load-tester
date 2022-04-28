package account

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
)

var ERC721Ledger map[common.Address]chan *big.Int // each index indicates ERC721 tokens, which the corresponding account in accGrp owns

func (self *Account) DeployERC721(c *client.Client, to *Account, value *big.Int, humanReadable bool) (common.Address, *types.Transaction, *big.Int, error) {
	ctx := context.Background()

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

	contractABI := erc721PerformanceABI
	parsed, err := abi.JSON(strings.NewReader(contractABI))
	byteCode := erc721PerformanceByteCode

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
	contractAddr, contractTx, _, err := bind.DeployContract(txOpts, parsed, byteCode, c, "thisIsName", "thisIsSymbol")
	if err != nil {
		fmt.Printf("Failed to DeployContract for ERC721 err: %v \n", err)
		return common.Address{}, nil, nil, err
	}

	self.nonce++
	return contractAddr, contractTx, gasPrice, nil
}

func (self *Account) MintERC721ToTestAccounts(c *client.Client, accGrp []*Account, smartContractAddr common.Address, numInitialTokensPerAccount int) {
	ERC721Ledger = make(map[common.Address]chan *big.Int)

	// Initialize ERC721 ledger before minting
	for _, acc := range accGrp {
		ERC721Ledger[acc.address] = make(chan *big.Int, numInitialTokensPerAccount*5)
	}
	// Randomly assign start token id to each locust slave to avoid token id collision
	rand.Seed(time.Now().UnixNano())
	startTokenId := int64(rand.Intn(len(accGrp)) * 100 * numInitialTokensPerAccount)
	endTokenId := startTokenId + int64(numInitialTokensPerAccount)

	for _, tokenRecipient := range accGrp {
		for {
			_, err := self.mintERC721ToTestAccounts(c, smartContractAddr, tokenRecipient, startTokenId, endTokenId)
			if err != nil {
				log.Printf("Error while minting ERC721 to test account, err: %v, recipient: %v", err, tokenRecipient.address.String())
				time.Sleep(1 * time.Second) // Mostly the error happens due to full txpool, wait 1 second
				continue
			}
			log.Println("MintERC721", "from", self.address.String(), "to", tokenRecipient.address.String(),
				"startTokenId", startTokenId, "endTokenId", endTokenId)
			break
		}
		startTokenId = endTokenId
		endTokenId += int64(numInitialTokensPerAccount)
	}

	log.Println("End MintERC721ToTestAccounts")
}

func (self *Account) mintERC721ToTestAccounts(c *client.Client, smartContractAddr common.Address, tokenRecipient *Account, startTokenId, endTokenId int64) (*types.Transaction, error) {
	ctx := context.Background()

	self.mutex.Lock()
	defer self.mutex.Unlock()

	nonce := self.GetNonce(c)

	abiStr := erc721PerformanceABI

	abii, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		log.Fatalf("[mintERC721ToTestAccounts] failed to abi.JSON: %v", err)
	}

	data, err := abii.Pack("registerBulk", tokenRecipient.address, big.NewInt(startTokenId), big.NewInt(endTokenId))
	if err != nil {
		log.Fatalf("[mintERC721ToTestAccounts] failed to abi.Pack: %v", err)
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: uint64(50000000000),
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   big.NewInt(0),
		types.TxValueKeyTo:       smartContractAddr,
		types.TxValueKeyData:     data,
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
		return tx, err
	}

	self.nonce++
	// update erc721TransferTC.ERC721Ledger with newly minted tokens
	for tokenId := startTokenId; tokenId < endTokenId; tokenId++ {
		ERC721Ledger[tokenRecipient.address] <- big.NewInt(tokenId)
	}

	return tx, nil
}

func (self *Account) TransferERC721(initialCharge bool, c *client.Client, tokenContractAddr common.Address, tokenRecipient *Account, tokenId *big.Int) (*types.Transaction, *big.Int, error) {
	ctx := context.Background()

	self.mutex.Lock()
	defer self.mutex.Unlock()

	var nonce uint64
	if initialCharge {
		nonce = self.GetNonceFromBlock(c)
	} else {
		nonce = self.GetNonce(c)
	}

	abiStr := erc721PerformanceABI
	abii, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}
	data, err := abii.Pack("transferFrom", self.address, tokenRecipient.address, tokenId)
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: uint64(5000000),
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   big.NewInt(0),
		types.TxValueKeyTo:       tokenContractAddr,
		types.TxValueKeyData:     data,
	})
	if err != nil {
		log.Fatalf("Failed to encode tx: %v", err)
	}

	log.Println("TransferERC721", "from", self.address.String(), "to", tokenRecipient.address.String(), "tokenId", tokenId.String())

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
		return tx, gasPrice, err
	}

	self.nonce++

	return tx, gasPrice, nil
}

func (self *Account) AddMinter(c *client.Client, minterCandidate *Account, smartContractAddr common.Address) (*types.Transaction, error) {
	ctx := context.Background()

	self.mutex.Lock()
	defer self.mutex.Unlock()

	fmt.Printf("Start AddMinter, minterCandidateAddr: %v \n", minterCandidate.GetAddress().String())

	nonce := self.GetNonceFromBlock(c)

	abii := getERC721PerformanceABII()
	data, err := abii.Pack("addMinter", minterCandidate.GetAddress())
	if err != nil {
		log.Fatalf("failed to abi.Pack: %v", err)
	}

	signer := types.NewEIP155Signer(chainID)
	tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce,
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyGasLimit: uint64(5000000),
		types.TxValueKeyFrom:     self.address,
		types.TxValueKeyAmount:   big.NewInt(0),
		types.TxValueKeyTo:       smartContractAddr,
		types.TxValueKeyData:     data,
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
		return tx, err
	}

	self.nonce++
	fmt.Printf("End AddMinter, minterCandidateAddr: %v \n", minterCandidate.GetAddress().String())
	return tx, nil
}

func getERC721PerformanceABII() abi.ABI {
	abiStr := erc721PerformanceABI
	abii, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		log.Fatalf("failed to abi.JSON: %v", err)
	}
	return abii
}
