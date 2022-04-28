// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ycsbTC

import (
	"strings"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

// KVstoreABI is the input ABI used to generate the binding from.
const KVstoreABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"}],\"name\":\"get\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"},{\"name\":\"value\",\"type\":\"string\"}],\"name\":\"set\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// KVstoreBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const KVstoreBinRuntime = `0x60806040526004361061004b5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663693ec85e8114610050578063e942b5161461011e575b600080fd5b34801561005c57600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100a99436949293602493928401919081908401838280828437509497506101b79650505050505050565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100e35781810151838201526020016100cb565b50505050905090810190601f1680156101105780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561012a57600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526101b594369492936024939284019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497506102ab9650505050505050565b005b60606000826040518082805190602001908083835b602083106101eb5780518252601f1990920191602091820191016101cc565b518151600019602094850361010090810a820192831692199390931691909117909252949092019687526040805197889003820188208054601f600260018316159098029095011695909504928301829004820288018201905281875292945092505083018282801561029f5780601f106102745761010080835404028352916020019161029f565b820191906000526020600020905b81548152906001019060200180831161028257829003601f168201915b50505050509050919050565b806000836040518082805190602001908083835b602083106102de5780518252601f1990920191602091820191016102bf565b51815160209384036101000a6000190180199092169116179052920194855250604051938490038101909320845161031f9591949190910192509050610324565b505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061036557805160ff1916838001178555610392565b82800160010185558215610392579182015b82811115610392578251825591602001919060010190610377565b5061039e9291506103a2565b5090565b6103bc91905b8082111561039e57600081556001016103a8565b905600a165627a7a72305820557c473a2a7ac0d4df8f90b0b7cde3307059a5c0e8da71d33c70f7079aa80e020029`

// KVstoreBin is the compiled bytecode used for deploying new contracts.
const KVstoreBin = `0x608060405234801561001057600080fd5b506103eb806100206000396000f30060806040526004361061004b5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663693ec85e8114610050578063e942b5161461011e575b600080fd5b34801561005c57600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100a99436949293602493928401919081908401838280828437509497506101b79650505050505050565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100e35781810151838201526020016100cb565b50505050905090810190601f1680156101105780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561012a57600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526101b594369492936024939284019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497506102ab9650505050505050565b005b60606000826040518082805190602001908083835b602083106101eb5780518252601f1990920191602091820191016101cc565b518151600019602094850361010090810a820192831692199390931691909117909252949092019687526040805197889003820188208054601f600260018316159098029095011695909504928301829004820288018201905281875292945092505083018282801561029f5780601f106102745761010080835404028352916020019161029f565b820191906000526020600020905b81548152906001019060200180831161028257829003601f168201915b50505050509050919050565b806000836040518082805190602001908083835b602083106102de5780518252601f1990920191602091820191016102bf565b51815160209384036101000a6000190180199092169116179052920194855250604051938490038101909320845161031f9591949190910192509050610324565b505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061036557805160ff1916838001178555610392565b82800160010185558215610392579182015b82811115610392578251825591602001919060010190610377565b5061039e9291506103a2565b5090565b6103bc91905b8082111561039e57600081556001016103a8565b905600a165627a7a72305820557c473a2a7ac0d4df8f90b0b7cde3307059a5c0e8da71d33c70f7079aa80e020029`

// DeployKVstore deploys a new GXP contract, binding an instance of KVstore to it.
func DeployKVstore(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *KVstore, error) {
	parsed, err := abi.JSON(strings.NewReader(KVstoreABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(KVstoreBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KVstore{KVstoreCaller: KVstoreCaller{contract: contract}, KVstoreTransactor: KVstoreTransactor{contract: contract}, KVstoreFilterer: KVstoreFilterer{contract: contract}}, nil
}

// KVstore is an auto generated Go binding around an GXP contract.
type KVstore struct {
	KVstoreCaller     // Read-only binding to the contract
	KVstoreTransactor // Write-only binding to the contract
	KVstoreFilterer   // Log filterer for contract events
}

// KVstoreCaller is an auto generated read-only Go binding around an GXP contract.
type KVstoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KVstoreTransactor is an auto generated write-only Go binding around an GXP contract.
type KVstoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KVstoreFilterer is an auto generated log filtering Go binding around an GXP contract events.
type KVstoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KVstoreSession is an auto generated Go binding around an GXP contract,
// with pre-set call and transact options.
type KVstoreSession struct {
	Contract     *KVstore          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KVstoreCallerSession is an auto generated read-only Go binding around an GXP contract,
// with pre-set call options.
type KVstoreCallerSession struct {
	Contract *KVstoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// KVstoreTransactorSession is an auto generated write-only Go binding around an GXP contract,
// with pre-set transact options.
type KVstoreTransactorSession struct {
	Contract     *KVstoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// KVstoreRaw is an auto generated low-level Go binding around an GXP contract.
type KVstoreRaw struct {
	Contract *KVstore // Generic contract binding to access the raw methods on
}

// KVstoreCallerRaw is an auto generated low-level read-only Go binding around an GXP contract.
type KVstoreCallerRaw struct {
	Contract *KVstoreCaller // Generic read-only contract binding to access the raw methods on
}

// KVstoreTransactorRaw is an auto generated low-level write-only Go binding around an GXP contract.
type KVstoreTransactorRaw struct {
	Contract *KVstoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKVstore creates a new instance of KVstore, bound to a specific deployed contract.
func NewKVstore(address common.Address, backend bind.ContractBackend) (*KVstore, error) {
	contract, err := bindKVstore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KVstore{KVstoreCaller: KVstoreCaller{contract: contract}, KVstoreTransactor: KVstoreTransactor{contract: contract}, KVstoreFilterer: KVstoreFilterer{contract: contract}}, nil
}

// NewKVstoreCaller creates a new read-only instance of KVstore, bound to a specific deployed contract.
func NewKVstoreCaller(address common.Address, caller bind.ContractCaller) (*KVstoreCaller, error) {
	contract, err := bindKVstore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KVstoreCaller{contract: contract}, nil
}

// NewKVstoreTransactor creates a new write-only instance of KVstore, bound to a specific deployed contract.
func NewKVstoreTransactor(address common.Address, transactor bind.ContractTransactor) (*KVstoreTransactor, error) {
	contract, err := bindKVstore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KVstoreTransactor{contract: contract}, nil
}

// NewKVstoreFilterer creates a new log filterer instance of KVstore, bound to a specific deployed contract.
func NewKVstoreFilterer(address common.Address, filterer bind.ContractFilterer) (*KVstoreFilterer, error) {
	contract, err := bindKVstore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KVstoreFilterer{contract: contract}, nil
}

// bindKVstore binds a generic wrapper to an already deployed contract.
func bindKVstore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KVstoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KVstore *KVstoreRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KVstore.Contract.KVstoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KVstore *KVstoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KVstore.Contract.KVstoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KVstore *KVstoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KVstore.Contract.KVstoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KVstore *KVstoreCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KVstore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KVstore *KVstoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KVstore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KVstore *KVstoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KVstore.Contract.contract.Transact(opts, method, params...)
}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(key string) constant returns(string)
func (_KVstore *KVstoreCaller) Get(opts *bind.CallOpts, key string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _KVstore.contract.Call(opts, out, "get", key)
	return *ret0, err
}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(key string) constant returns(string)
func (_KVstore *KVstoreSession) Get(key string) (string, error) {
	return _KVstore.Contract.Get(&_KVstore.CallOpts, key)
}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(key string) constant returns(string)
func (_KVstore *KVstoreCallerSession) Get(key string) (string, error) {
	return _KVstore.Contract.Get(&_KVstore.CallOpts, key)
}

// Set is a paid mutator transaction binding the contract method 0xe942b516.
//
// Solidity: function set(key string, value string) returns()
func (_KVstore *KVstoreTransactor) Set(opts *bind.TransactOpts, key string, value string) (*types.Transaction, error) {
	return _KVstore.contract.Transact(opts, "set", key, value)
}

// Set is a paid mutator transaction binding the contract method 0xe942b516.
//
// Solidity: function set(key string, value string) returns()
func (_KVstore *KVstoreSession) Set(key string, value string) (*types.Transaction, error) {
	return _KVstore.Contract.Set(&_KVstore.TransactOpts, key, value)
}

// Set is a paid mutator transaction binding the contract method 0xe942b516.
//
// Solidity: function set(key string, value string) returns()
func (_KVstore *KVstoreTransactorSession) Set(key string, value string) (*types.Transaction, error) {
	return _KVstore.Contract.Set(&_KVstore.TransactOpts, key, value)
}
