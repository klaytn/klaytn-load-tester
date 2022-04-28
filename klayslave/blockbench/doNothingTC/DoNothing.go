// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package doNothingTC

import (
	"strings"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

// DoNothingABI is the input ABI used to generate the binding from.
const DoNothingABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"nothing\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// DoNothingBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const DoNothingBinRuntime = `0x608060405260043610603e5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663448f30a381146043575b600080fd5b348015604e57600080fd5b5060556057565b005b5600a165627a7a723058200745444e555b8b16ea8d72738c5d1bd7972daef0924dd2f7aee967e348bff9f50029`

// DoNothingBin is the compiled bytecode used for deploying new contracts.
const DoNothingBin = `0x6080604052348015600f57600080fd5b5060858061001e6000396000f300608060405260043610603e5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663448f30a381146043575b600080fd5b348015604e57600080fd5b5060556057565b005b5600a165627a7a723058200745444e555b8b16ea8d72738c5d1bd7972daef0924dd2f7aee967e348bff9f50029`

// DeployDoNothing deploys a new GXP contract, binding an instance of DoNothing to it.
func DeployDoNothing(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DoNothing, error) {
	parsed, err := abi.JSON(strings.NewReader(DoNothingABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DoNothingBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DoNothing{DoNothingCaller: DoNothingCaller{contract: contract}, DoNothingTransactor: DoNothingTransactor{contract: contract}, DoNothingFilterer: DoNothingFilterer{contract: contract}}, nil
}

// DoNothing is an auto generated Go binding around an GXP contract.
type DoNothing struct {
	DoNothingCaller     // Read-only binding to the contract
	DoNothingTransactor // Write-only binding to the contract
	DoNothingFilterer   // Log filterer for contract events
}

// DoNothingCaller is an auto generated read-only Go binding around an GXP contract.
type DoNothingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DoNothingTransactor is an auto generated write-only Go binding around an GXP contract.
type DoNothingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DoNothingFilterer is an auto generated log filtering Go binding around an GXP contract events.
type DoNothingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DoNothingSession is an auto generated Go binding around an GXP contract,
// with pre-set call and transact options.
type DoNothingSession struct {
	Contract     *DoNothing        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DoNothingCallerSession is an auto generated read-only Go binding around an GXP contract,
// with pre-set call options.
type DoNothingCallerSession struct {
	Contract *DoNothingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// DoNothingTransactorSession is an auto generated write-only Go binding around an GXP contract,
// with pre-set transact options.
type DoNothingTransactorSession struct {
	Contract     *DoNothingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// DoNothingRaw is an auto generated low-level Go binding around an GXP contract.
type DoNothingRaw struct {
	Contract *DoNothing // Generic contract binding to access the raw methods on
}

// DoNothingCallerRaw is an auto generated low-level read-only Go binding around an GXP contract.
type DoNothingCallerRaw struct {
	Contract *DoNothingCaller // Generic read-only contract binding to access the raw methods on
}

// DoNothingTransactorRaw is an auto generated low-level write-only Go binding around an GXP contract.
type DoNothingTransactorRaw struct {
	Contract *DoNothingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDoNothing creates a new instance of DoNothing, bound to a specific deployed contract.
func NewDoNothing(address common.Address, backend bind.ContractBackend) (*DoNothing, error) {
	contract, err := bindDoNothing(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DoNothing{DoNothingCaller: DoNothingCaller{contract: contract}, DoNothingTransactor: DoNothingTransactor{contract: contract}, DoNothingFilterer: DoNothingFilterer{contract: contract}}, nil
}

// NewDoNothingCaller creates a new read-only instance of DoNothing, bound to a specific deployed contract.
func NewDoNothingCaller(address common.Address, caller bind.ContractCaller) (*DoNothingCaller, error) {
	contract, err := bindDoNothing(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DoNothingCaller{contract: contract}, nil
}

// NewDoNothingTransactor creates a new write-only instance of DoNothing, bound to a specific deployed contract.
func NewDoNothingTransactor(address common.Address, transactor bind.ContractTransactor) (*DoNothingTransactor, error) {
	contract, err := bindDoNothing(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DoNothingTransactor{contract: contract}, nil
}

// NewDoNothingFilterer creates a new log filterer instance of DoNothing, bound to a specific deployed contract.
func NewDoNothingFilterer(address common.Address, filterer bind.ContractFilterer) (*DoNothingFilterer, error) {
	contract, err := bindDoNothing(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DoNothingFilterer{contract: contract}, nil
}

// bindDoNothing binds a generic wrapper to an already deployed contract.
func bindDoNothing(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DoNothingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DoNothing *DoNothingRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DoNothing.Contract.DoNothingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DoNothing *DoNothingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DoNothing.Contract.DoNothingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DoNothing *DoNothingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DoNothing.Contract.DoNothingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DoNothing *DoNothingCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DoNothing.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DoNothing *DoNothingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DoNothing.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DoNothing *DoNothingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DoNothing.Contract.contract.Transact(opts, method, params...)
}

// Nothing is a paid mutator transaction binding the contract method 0x448f30a3.
//
// Solidity: function nothing() returns()
func (_DoNothing *DoNothingTransactor) Nothing(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DoNothing.contract.Transact(opts, "nothing")
}

// Nothing is a paid mutator transaction binding the contract method 0x448f30a3.
//
// Solidity: function nothing() returns()
func (_DoNothing *DoNothingSession) Nothing() (*types.Transaction, error) {
	return _DoNothing.Contract.Nothing(&_DoNothing.TransactOpts)
}

// Nothing is a paid mutator transaction binding the contract method 0x448f30a3.
//
// Solidity: function nothing() returns()
func (_DoNothing *DoNothingTransactorSession) Nothing() (*types.Transaction, error) {
	return _DoNothing.Contract.Nothing(&_DoNothing.TransactOpts)
}
