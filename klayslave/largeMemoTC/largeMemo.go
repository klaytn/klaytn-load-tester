// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package largeMemoTC

import (
	"math/big"
	"strings"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = klaytn.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// LargeMemoABI is the input ABI used to generate the binding from.
const LargeMemoABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"run\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"str\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_str\",\"type\":\"string\"}],\"name\":\"setName\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// LargeMemoBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const LargeMemoBinRuntime = `6080604052600436106100565763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663c0406226811461005b578063c15bae84146100e5578063c47f0027146100fa575b600080fd5b34801561006757600080fd5b50610070610155565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100aa578181015183820152602001610092565b50505050905090810190601f1680156100d75780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156100f157600080fd5b506100706101ec565b34801561010657600080fd5b506040805160206004803580820135601f810184900484028501840190955284845261015394369492936024939284019190819084018382808284375094975061027a9650505050505050565b005b60008054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156101e15780601f106101b6576101008083540402835291602001916101e1565b820191906000526020600020905b8154815290600101906020018083116101c457829003601f168201915b505050505090505b90565b6000805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156102725780601f1061024757610100808354040283529160200191610272565b820191906000526020600020905b81548152906001019060200180831161025557829003601f168201915b505050505081565b805161028d906000906020840190610291565b5050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106102d257805160ff19168380011785556102ff565b828001600101855582156102ff579182015b828111156102ff5782518255916020019190600101906102e4565b5061030b92915061030f565b5090565b6101e991905b8082111561030b57600081556001016103155600a165627a7a72305820cb9945dc3b1f86fcbf14778927de89ef6badbc0ad78bfe86e770376894cd74330029`

// LargeMemoFuncSigs maps the 4-byte function signature to its string representation.
var LargeMemoFuncSigs = map[string]string{
	"c0406226": "run()",
	"c47f0027": "setName(string)",
	"c15bae84": "str()",
}

// LargeMemoBin is the compiled bytecode used for deploying new contracts.
var LargeMemoBin = "0x608060405234801561001057600080fd5b5060408051808201909152600c8082527f48656c6c6f2c20576f726c64000000000000000000000000000000000000000060209092019182526100559160009161005b565b506100f6565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061009c57805160ff19168380011785556100c9565b828001600101855582156100c9579182015b828111156100c95782518255916020019190600101906100ae565b506100d59291506100d9565b5090565b6100f391905b808211156100d557600081556001016100df565b90565b610355806101056000396000f3006080604052600436106100565763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663c0406226811461005b578063c15bae84146100e5578063c47f0027146100fa575b600080fd5b34801561006757600080fd5b50610070610155565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100aa578181015183820152602001610092565b50505050905090810190601f1680156100d75780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156100f157600080fd5b506100706101ec565b34801561010657600080fd5b506040805160206004803580820135601f810184900484028501840190955284845261015394369492936024939284019190819084018382808284375094975061027a9650505050505050565b005b60008054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156101e15780601f106101b6576101008083540402835291602001916101e1565b820191906000526020600020905b8154815290600101906020018083116101c457829003601f168201915b505050505090505b90565b6000805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156102725780601f1061024757610100808354040283529160200191610272565b820191906000526020600020905b81548152906001019060200180831161025557829003601f168201915b505050505081565b805161028d906000906020840190610291565b5050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106102d257805160ff19168380011785556102ff565b828001600101855582156102ff579182015b828111156102ff5782518255916020019190600101906102e4565b5061030b92915061030f565b5090565b6101e991905b8082111561030b57600081556001016103155600a165627a7a72305820cb9945dc3b1f86fcbf14778927de89ef6badbc0ad78bfe86e770376894cd74330029"

// DeployLargeMemo deploys a new Klaytn contract, binding an instance of LargeMemo to it.
func DeployLargeMemo(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LargeMemo, error) {
	parsed, err := abi.JSON(strings.NewReader(LargeMemoABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LargeMemoBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LargeMemo{LargeMemoCaller: LargeMemoCaller{contract: contract}, LargeMemoTransactor: LargeMemoTransactor{contract: contract}, LargeMemoFilterer: LargeMemoFilterer{contract: contract}}, nil
}

// LargeMemo is an auto generated Go binding around a Klaytn contract.
type LargeMemo struct {
	LargeMemoCaller     // Read-only binding to the contract
	LargeMemoTransactor // Write-only binding to the contract
	LargeMemoFilterer   // Log filterer for contract events
}

// LargeMemoCaller is an auto generated read-only Go binding around a Klaytn contract.
type LargeMemoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LargeMemoTransactor is an auto generated write-only Go binding around a Klaytn contract.
type LargeMemoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LargeMemoFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type LargeMemoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LargeMemoSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type LargeMemoSession struct {
	Contract     *LargeMemo        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LargeMemoCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type LargeMemoCallerSession struct {
	Contract *LargeMemoCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// LargeMemoTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type LargeMemoTransactorSession struct {
	Contract     *LargeMemoTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// LargeMemoRaw is an auto generated low-level Go binding around a Klaytn contract.
type LargeMemoRaw struct {
	Contract *LargeMemo // Generic contract binding to access the raw methods on
}

// LargeMemoCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type LargeMemoCallerRaw struct {
	Contract *LargeMemoCaller // Generic read-only contract binding to access the raw methods on
}

// LargeMemoTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type LargeMemoTransactorRaw struct {
	Contract *LargeMemoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLargeMemo creates a new instance of LargeMemo, bound to a specific deployed contract.
func NewLargeMemo(address common.Address, backend bind.ContractBackend) (*LargeMemo, error) {
	contract, err := bindLargeMemo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LargeMemo{LargeMemoCaller: LargeMemoCaller{contract: contract}, LargeMemoTransactor: LargeMemoTransactor{contract: contract}, LargeMemoFilterer: LargeMemoFilterer{contract: contract}}, nil
}

// NewLargeMemoCaller creates a new read-only instance of LargeMemo, bound to a specific deployed contract.
func NewLargeMemoCaller(address common.Address, caller bind.ContractCaller) (*LargeMemoCaller, error) {
	contract, err := bindLargeMemo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LargeMemoCaller{contract: contract}, nil
}

// NewLargeMemoTransactor creates a new write-only instance of LargeMemo, bound to a specific deployed contract.
func NewLargeMemoTransactor(address common.Address, transactor bind.ContractTransactor) (*LargeMemoTransactor, error) {
	contract, err := bindLargeMemo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LargeMemoTransactor{contract: contract}, nil
}

// NewLargeMemoFilterer creates a new log filterer instance of LargeMemo, bound to a specific deployed contract.
func NewLargeMemoFilterer(address common.Address, filterer bind.ContractFilterer) (*LargeMemoFilterer, error) {
	contract, err := bindLargeMemo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LargeMemoFilterer{contract: contract}, nil
}

// bindLargeMemo binds a generic wrapper to an already deployed contract.
func bindLargeMemo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LargeMemoABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LargeMemo *LargeMemoRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LargeMemo.Contract.LargeMemoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LargeMemo *LargeMemoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LargeMemo.Contract.LargeMemoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LargeMemo *LargeMemoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LargeMemo.Contract.LargeMemoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LargeMemo *LargeMemoCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LargeMemo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LargeMemo *LargeMemoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LargeMemo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LargeMemo *LargeMemoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LargeMemo.Contract.contract.Transact(opts, method, params...)
}

// Run is a free data retrieval call binding the contract method 0xc0406226.
//
// Solidity: function run() view returns(string)
func (_LargeMemo *LargeMemoCaller) Run(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LargeMemo.contract.Call(opts, out, "run")
	return *ret0, err
}

// Run is a free data retrieval call binding the contract method 0xc0406226.
//
// Solidity: function run() view returns(string)
func (_LargeMemo *LargeMemoSession) Run() (string, error) {
	return _LargeMemo.Contract.Run(&_LargeMemo.CallOpts)
}

// Run is a free data retrieval call binding the contract method 0xc0406226.
//
// Solidity: function run() view returns(string)
func (_LargeMemo *LargeMemoCallerSession) Run() (string, error) {
	return _LargeMemo.Contract.Run(&_LargeMemo.CallOpts)
}

// Str is a free data retrieval call binding the contract method 0xc15bae84.
//
// Solidity: function str() view returns(string)
func (_LargeMemo *LargeMemoCaller) Str(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _LargeMemo.contract.Call(opts, out, "str")
	return *ret0, err
}

// Str is a free data retrieval call binding the contract method 0xc15bae84.
//
// Solidity: function str() view returns(string)
func (_LargeMemo *LargeMemoSession) Str() (string, error) {
	return _LargeMemo.Contract.Str(&_LargeMemo.CallOpts)
}

// Str is a free data retrieval call binding the contract method 0xc15bae84.
//
// Solidity: function str() view returns(string)
func (_LargeMemo *LargeMemoCallerSession) Str() (string, error) {
	return _LargeMemo.Contract.Str(&_LargeMemo.CallOpts)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string _str) returns()
func (_LargeMemo *LargeMemoTransactor) SetName(opts *bind.TransactOpts, _str string) (*types.Transaction, error) {
	t, err := _LargeMemo.contract.Transact(opts, "setName", _str)
	return t, err
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string _str) returns()
func (_LargeMemo *LargeMemoSession) SetName(_str string) (*types.Transaction, error) {
	return _LargeMemo.Contract.SetName(&_LargeMemo.TransactOpts, _str)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string _str) returns()
func (_LargeMemo *LargeMemoTransactorSession) SetName(_str string) (*types.Transaction, error) {
	return _LargeMemo.Contract.SetName(&_LargeMemo.TransactOpts, _str)
}
