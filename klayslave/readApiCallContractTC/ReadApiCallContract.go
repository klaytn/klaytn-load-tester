// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package readApiCallContractTC

import (
	"math/big"
	"strings"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

// ReadApiCallContractABI is the input ABI used to generate the binding from.
const ReadApiCallContractABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"set\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ReadApiCallContractBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ReadApiCallContractBinRuntime = `0x60806040526004361060485763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416636d4ce63c8114604d578063b8e010de146071575b600080fd5b348015605857600080fd5b50605f6085565b60408051918252519081900360200190f35b348015607c57600080fd5b506083608b565b005b60005490565b60086001555600a165627a7a72305820895a5f786a9fee2bf16c4e6494b5ff67d1811d3ac3496ee4347005d251a9885d0029`

// ReadApiCallContractBin is the compiled bytecode used for deploying new contracts.
const ReadApiCallContractBin = `0x6080604052600460005534801561001557600080fd5b5060be806100246000396000f30060806040526004361060485763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416636d4ce63c8114604d578063b8e010de146071575b600080fd5b348015605857600080fd5b50605f6085565b60408051918252519081900360200190f35b348015607c57600080fd5b506083608b565b005b60005490565b60086001555600a165627a7a72305820895a5f786a9fee2bf16c4e6494b5ff67d1811d3ac3496ee4347005d251a9885d0029`

// DeployReadApiCallContract deploys a new Klaytn contract, binding an instance of ReadApiCallContract to it.
func DeployReadApiCallContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ReadApiCallContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ReadApiCallContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ReadApiCallContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ReadApiCallContract{ReadApiCallContractCaller: ReadApiCallContractCaller{contract: contract}, ReadApiCallContractTransactor: ReadApiCallContractTransactor{contract: contract}, ReadApiCallContractFilterer: ReadApiCallContractFilterer{contract: contract}}, nil
}

// ReadApiCallContract is an auto generated Go binding around a Klaytn contract.
type ReadApiCallContract struct {
	ReadApiCallContractCaller     // Read-only binding to the contract
	ReadApiCallContractTransactor // Write-only binding to the contract
	ReadApiCallContractFilterer   // Log filterer for contract events
}

// ReadApiCallContractCaller is an auto generated read-only Go binding around a Klaytn contract.
type ReadApiCallContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReadApiCallContractTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ReadApiCallContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReadApiCallContractFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ReadApiCallContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReadApiCallContractSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ReadApiCallContractSession struct {
	Contract     *ReadApiCallContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ReadApiCallContractCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ReadApiCallContractCallerSession struct {
	Contract *ReadApiCallContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ReadApiCallContractTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ReadApiCallContractTransactorSession struct {
	Contract     *ReadApiCallContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ReadApiCallContractRaw is an auto generated low-level Go binding around a Klaytn contract.
type ReadApiCallContractRaw struct {
	Contract *ReadApiCallContract // Generic contract binding to access the raw methods on
}

// ReadApiCallContractCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ReadApiCallContractCallerRaw struct {
	Contract *ReadApiCallContractCaller // Generic read-only contract binding to access the raw methods on
}

// ReadApiCallContractTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ReadApiCallContractTransactorRaw struct {
	Contract *ReadApiCallContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReadApiCallContract creates a new instance of ReadApiCallContract, bound to a specific deployed contract.
func NewReadApiCallContract(address common.Address, backend bind.ContractBackend) (*ReadApiCallContract, error) {
	contract, err := bindReadApiCallContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReadApiCallContract{ReadApiCallContractCaller: ReadApiCallContractCaller{contract: contract}, ReadApiCallContractTransactor: ReadApiCallContractTransactor{contract: contract}, ReadApiCallContractFilterer: ReadApiCallContractFilterer{contract: contract}}, nil
}

// NewReadApiCallContractCaller creates a new read-only instance of ReadApiCallContract, bound to a specific deployed contract.
func NewReadApiCallContractCaller(address common.Address, caller bind.ContractCaller) (*ReadApiCallContractCaller, error) {
	contract, err := bindReadApiCallContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReadApiCallContractCaller{contract: contract}, nil
}

// NewReadApiCallContractTransactor creates a new write-only instance of ReadApiCallContract, bound to a specific deployed contract.
func NewReadApiCallContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ReadApiCallContractTransactor, error) {
	contract, err := bindReadApiCallContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReadApiCallContractTransactor{contract: contract}, nil
}

// NewReadApiCallContractFilterer creates a new log filterer instance of ReadApiCallContract, bound to a specific deployed contract.
func NewReadApiCallContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ReadApiCallContractFilterer, error) {
	contract, err := bindReadApiCallContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReadApiCallContractFilterer{contract: contract}, nil
}

// bindReadApiCallContract binds a generic wrapper to an already deployed contract.
func bindReadApiCallContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ReadApiCallContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReadApiCallContract *ReadApiCallContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ReadApiCallContract.Contract.ReadApiCallContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReadApiCallContract *ReadApiCallContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReadApiCallContract.Contract.ReadApiCallContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReadApiCallContract *ReadApiCallContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReadApiCallContract.Contract.ReadApiCallContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReadApiCallContract *ReadApiCallContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ReadApiCallContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReadApiCallContract *ReadApiCallContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReadApiCallContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReadApiCallContract *ReadApiCallContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReadApiCallContract.Contract.contract.Transact(opts, method, params...)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() constant returns(uint256)
func (_ReadApiCallContract *ReadApiCallContractCaller) Get(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ReadApiCallContract.contract.Call(opts, out, "get")
	return *ret0, err
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() constant returns(uint256)
func (_ReadApiCallContract *ReadApiCallContractSession) Get() (*big.Int, error) {
	return _ReadApiCallContract.Contract.Get(&_ReadApiCallContract.CallOpts)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() constant returns(uint256)
func (_ReadApiCallContract *ReadApiCallContractCallerSession) Get() (*big.Int, error) {
	return _ReadApiCallContract.Contract.Get(&_ReadApiCallContract.CallOpts)
}

// Set is a paid mutator transaction binding the contract method 0xb8e010de.
//
// Solidity: function set() returns()
func (_ReadApiCallContract *ReadApiCallContractTransactor) Set(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReadApiCallContract.contract.Transact(opts, "set")
}

// Set is a paid mutator transaction binding the contract method 0xb8e010de.
//
// Solidity: function set() returns()
func (_ReadApiCallContract *ReadApiCallContractSession) Set() (*types.Transaction, error) {
	return _ReadApiCallContract.Contract.Set(&_ReadApiCallContract.TransactOpts)
}

// Set is a paid mutator transaction binding the contract method 0xb8e010de.
//
// Solidity: function set() returns()
func (_ReadApiCallContract *ReadApiCallContractTransactorSession) Set() (*types.Transaction, error) {
	return _ReadApiCallContract.Contract.Set(&_ReadApiCallContract.TransactOpts)
}
