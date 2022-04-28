// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ioHeavyTC

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

// IOHeavyABI is the input ABI used to generate the binding from.
const IOHeavyABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"key\",\"type\":\"bytes20\"}],\"name\":\"get\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"start_key\",\"type\":\"uint256\"},{\"name\":\"size\",\"type\":\"uint256\"},{\"name\":\"signature\",\"type\":\"uint256\"}],\"name\":\"scan\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"start_key\",\"type\":\"uint256\"},{\"name\":\"size\",\"type\":\"uint256\"},{\"name\":\"signature\",\"type\":\"uint256\"}],\"name\":\"revert_scan\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"start_key\",\"type\":\"uint256\"},{\"name\":\"size\",\"type\":\"uint256\"},{\"name\":\"signature\",\"type\":\"uint256\"}],\"name\":\"write\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"bytes20\"},{\"name\":\"value\",\"type\":\"bytes\"}],\"name\":\"set\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"size\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"signature\",\"type\":\"uint256\"}],\"name\":\"finishWrite\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"size\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"signature\",\"type\":\"uint256\"}],\"name\":\"finishScan\",\"type\":\"event\"}]"

// IOHeavyBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IOHeavyBinRuntime = `0x60806040526004361061006c5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416635acecc7881146100715780636531695d1461010d578063c315d63e1461012d578063d4cd87901461014b578063d778e2da14610169575b600080fd5b34801561007d57600080fd5b506100986bffffffffffffffffffffffff19600435166101d6565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100d25781810151838201526020016100ba565b50505050905090810190601f1680156100ff5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561011957600080fd5b5061012b600435602435604435610285565b005b34801561013957600080fd5b5061012b6004356024356044356102f1565b34801561015757600080fd5b5061012b60043560243560443561031b565b34801561017557600080fd5b5060408051602060046024803582810135601f810185900485028601850190965285855261012b9583356bffffffffffffffffffffffff191695369560449491939091019190819084018382808284375094975061038d9650505050505050565b6bffffffffffffffffffffffff1981166000908152602081815260409182902080548351601f60026000196101006001861615020190931692909204918201849004840281018401909452808452606093928301828280156102795780601f1061024e57610100808354040283529160200191610279565b820191906000526020600020905b81548152906001019060200180831161025c57829003601f168201915b50505050509050919050565b606060005b838110156102af576102a56102a08287016103bf565b6101d6565b915060010161028a565b604080518581526020810185905281517f2e8128137e55a67bef5f6fa7e5c6722c5632e21b8c8bcf6df64bc32239dd6a3f929181900390910190a15050505050565b606060005b838110156102af576103116102a060018387890103036103bf565b91506001016102f6565b60005b8281101561034c576103446103348286016103bf565b61033f8387016103d0565b61038d565b60010161031e565b604080518481526020810184905281517fe849f68c74be0ec2d162615e7bc539b752b8e3e7db7ccb69f93eb19c85597f7e929181900390910190a150505050565b6bffffffffffffffffffffffff19821660009081526020818152604090912082516103ba9284019061050c565b505050565b60006103ca82610490565b92915050565b60408051606480825260a082019092526060916000919060208201610c8080388339019050509150600090505b606481101561048a5760c060405190810160405280609681526020016105a860969139805160328506830190811061043157fe5b90602001015160f860020a900460f860020a02828281518110151561045257fe5b9060200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506001016103fd565b50919050565b60008115156104c057507f3030303030303030303030303030303030303030000000000000000000000000610507565b6000821115610507576101006c010000000000000000000000008204046c01000000000000000000000000029050600a820660300160f860020a0217600a820491506104c0565b919050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061054d57805160ff191683800117855561057a565b8280016001018555821561057a579182015b8281111561057a57825182559160200191906001019061055f565b5061058692915061058a565b5090565b6105a491905b808211156105865760008155600101610590565b9056006162636465666768696a6b6c6d6e6f707172737475767778792324255e262a28295f2b5b5d7b7d7c3b3a2c2e2f3c3e3f607e6162636465666768696a6b6c6d6e6f707172737475767778792324255e262a28295f2b5b5d7b7d7c3b3a2c2e2f3c3e3f607e6162636465666768696a6b6c6d6e6f707172737475767778792324255e262a28295f2b5b5d7b7d7c3b3a2c2e2f3c3e3f607ea165627a7a7230582057064680bc1d92fb3eb1693cd3f18830fe3adc6f1a95e25a2798e3e44bfadba70029`

// IOHeavyBin is the compiled bytecode used for deploying new contracts.
const IOHeavyBin = `0x608060405234801561001057600080fd5b50610669806100206000396000f30060806040526004361061006c5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416635acecc7881146100715780636531695d1461010d578063c315d63e1461012d578063d4cd87901461014b578063d778e2da14610169575b600080fd5b34801561007d57600080fd5b506100986bffffffffffffffffffffffff19600435166101d6565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100d25781810151838201526020016100ba565b50505050905090810190601f1680156100ff5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561011957600080fd5b5061012b600435602435604435610285565b005b34801561013957600080fd5b5061012b6004356024356044356102f1565b34801561015757600080fd5b5061012b60043560243560443561031b565b34801561017557600080fd5b5060408051602060046024803582810135601f810185900485028601850190965285855261012b9583356bffffffffffffffffffffffff191695369560449491939091019190819084018382808284375094975061038d9650505050505050565b6bffffffffffffffffffffffff1981166000908152602081815260409182902080548351601f60026000196101006001861615020190931692909204918201849004840281018401909452808452606093928301828280156102795780601f1061024e57610100808354040283529160200191610279565b820191906000526020600020905b81548152906001019060200180831161025c57829003601f168201915b50505050509050919050565b606060005b838110156102af576102a56102a08287016103bf565b6101d6565b915060010161028a565b604080518581526020810185905281517f2e8128137e55a67bef5f6fa7e5c6722c5632e21b8c8bcf6df64bc32239dd6a3f929181900390910190a15050505050565b606060005b838110156102af576103116102a060018387890103036103bf565b91506001016102f6565b60005b8281101561034c576103446103348286016103bf565b61033f8387016103d0565b61038d565b60010161031e565b604080518481526020810184905281517fe849f68c74be0ec2d162615e7bc539b752b8e3e7db7ccb69f93eb19c85597f7e929181900390910190a150505050565b6bffffffffffffffffffffffff19821660009081526020818152604090912082516103ba9284019061050c565b505050565b60006103ca82610490565b92915050565b60408051606480825260a082019092526060916000919060208201610c8080388339019050509150600090505b606481101561048a5760c060405190810160405280609681526020016105a860969139805160328506830190811061043157fe5b90602001015160f860020a900460f860020a02828281518110151561045257fe5b9060200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506001016103fd565b50919050565b60008115156104c057507f3030303030303030303030303030303030303030000000000000000000000000610507565b6000821115610507576101006c010000000000000000000000008204046c01000000000000000000000000029050600a820660300160f860020a0217600a820491506104c0565b919050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061054d57805160ff191683800117855561057a565b8280016001018555821561057a579182015b8281111561057a57825182559160200191906001019061055f565b5061058692915061058a565b5090565b6105a491905b808211156105865760008155600101610590565b9056006162636465666768696a6b6c6d6e6f707172737475767778792324255e262a28295f2b5b5d7b7d7c3b3a2c2e2f3c3e3f607e6162636465666768696a6b6c6d6e6f707172737475767778792324255e262a28295f2b5b5d7b7d7c3b3a2c2e2f3c3e3f607e6162636465666768696a6b6c6d6e6f707172737475767778792324255e262a28295f2b5b5d7b7d7c3b3a2c2e2f3c3e3f607ea165627a7a7230582057064680bc1d92fb3eb1693cd3f18830fe3adc6f1a95e25a2798e3e44bfadba70029`

// DeployIOHeavy deploys a new Klaytn contract, binding an instance of IOHeavy to it.
func DeployIOHeavy(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *IOHeavy, error) {
	parsed, err := abi.JSON(strings.NewReader(IOHeavyABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(IOHeavyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &IOHeavy{IOHeavyCaller: IOHeavyCaller{contract: contract}, IOHeavyTransactor: IOHeavyTransactor{contract: contract}, IOHeavyFilterer: IOHeavyFilterer{contract: contract}}, nil
}

// IOHeavy is an auto generated Go binding around a Klaytn contract.
type IOHeavy struct {
	IOHeavyCaller     // Read-only binding to the contract
	IOHeavyTransactor // Write-only binding to the contract
	IOHeavyFilterer   // Log filterer for contract events
}

// IOHeavyCaller is an auto generated read-only Go binding around a Klaytn contract.
type IOHeavyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOHeavyTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IOHeavyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOHeavyFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IOHeavyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOHeavySession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IOHeavySession struct {
	Contract     *IOHeavy          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IOHeavyCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IOHeavyCallerSession struct {
	Contract *IOHeavyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// IOHeavyTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IOHeavyTransactorSession struct {
	Contract     *IOHeavyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IOHeavyRaw is an auto generated low-level Go binding around a Klaytn contract.
type IOHeavyRaw struct {
	Contract *IOHeavy // Generic contract binding to access the raw methods on
}

// IOHeavyCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IOHeavyCallerRaw struct {
	Contract *IOHeavyCaller // Generic read-only contract binding to access the raw methods on
}

// IOHeavyTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IOHeavyTransactorRaw struct {
	Contract *IOHeavyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIOHeavy creates a new instance of IOHeavy, bound to a specific deployed contract.
func NewIOHeavy(address common.Address, backend bind.ContractBackend) (*IOHeavy, error) {
	contract, err := bindIOHeavy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IOHeavy{IOHeavyCaller: IOHeavyCaller{contract: contract}, IOHeavyTransactor: IOHeavyTransactor{contract: contract}, IOHeavyFilterer: IOHeavyFilterer{contract: contract}}, nil
}

// NewIOHeavyCaller creates a new read-only instance of IOHeavy, bound to a specific deployed contract.
func NewIOHeavyCaller(address common.Address, caller bind.ContractCaller) (*IOHeavyCaller, error) {
	contract, err := bindIOHeavy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IOHeavyCaller{contract: contract}, nil
}

// NewIOHeavyTransactor creates a new write-only instance of IOHeavy, bound to a specific deployed contract.
func NewIOHeavyTransactor(address common.Address, transactor bind.ContractTransactor) (*IOHeavyTransactor, error) {
	contract, err := bindIOHeavy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IOHeavyTransactor{contract: contract}, nil
}

// NewIOHeavyFilterer creates a new log filterer instance of IOHeavy, bound to a specific deployed contract.
func NewIOHeavyFilterer(address common.Address, filterer bind.ContractFilterer) (*IOHeavyFilterer, error) {
	contract, err := bindIOHeavy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IOHeavyFilterer{contract: contract}, nil
}

// bindIOHeavy binds a generic wrapper to an already deployed contract.
func bindIOHeavy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IOHeavyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IOHeavy *IOHeavyRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IOHeavy.Contract.IOHeavyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IOHeavy *IOHeavyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOHeavy.Contract.IOHeavyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IOHeavy *IOHeavyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IOHeavy.Contract.IOHeavyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IOHeavy *IOHeavyCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _IOHeavy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IOHeavy *IOHeavyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOHeavy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IOHeavy *IOHeavyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IOHeavy.Contract.contract.Transact(opts, method, params...)
}

// Get is a free data retrieval call binding the contract method 0x5acecc78.
//
// Solidity: function get(key bytes20) constant returns(bytes)
func (_IOHeavy *IOHeavyCaller) Get(opts *bind.CallOpts, key [20]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _IOHeavy.contract.Call(opts, out, "get", key)
	return *ret0, err
}

// Get is a free data retrieval call binding the contract method 0x5acecc78.
//
// Solidity: function get(key bytes20) constant returns(bytes)
func (_IOHeavy *IOHeavySession) Get(key [20]byte) ([]byte, error) {
	return _IOHeavy.Contract.Get(&_IOHeavy.CallOpts, key)
}

// Get is a free data retrieval call binding the contract method 0x5acecc78.
//
// Solidity: function get(key bytes20) constant returns(bytes)
func (_IOHeavy *IOHeavyCallerSession) Get(key [20]byte) ([]byte, error) {
	return _IOHeavy.Contract.Get(&_IOHeavy.CallOpts, key)
}

// RevertScan is a paid mutator transaction binding the contract method 0xc315d63e.
//
// Solidity: function revert_scan(start_key uint256, size uint256, signature uint256) returns()
func (_IOHeavy *IOHeavyTransactor) RevertScan(opts *bind.TransactOpts, start_key *big.Int, size *big.Int, signature *big.Int) (*types.Transaction, error) {
	return _IOHeavy.contract.Transact(opts, "revert_scan", start_key, size, signature)
}

// RevertScan is a paid mutator transaction binding the contract method 0xc315d63e.
//
// Solidity: function revert_scan(start_key uint256, size uint256, signature uint256) returns()
func (_IOHeavy *IOHeavySession) RevertScan(start_key *big.Int, size *big.Int, signature *big.Int) (*types.Transaction, error) {
	return _IOHeavy.Contract.RevertScan(&_IOHeavy.TransactOpts, start_key, size, signature)
}

// RevertScan is a paid mutator transaction binding the contract method 0xc315d63e.
//
// Solidity: function revert_scan(start_key uint256, size uint256, signature uint256) returns()
func (_IOHeavy *IOHeavyTransactorSession) RevertScan(start_key *big.Int, size *big.Int, signature *big.Int) (*types.Transaction, error) {
	return _IOHeavy.Contract.RevertScan(&_IOHeavy.TransactOpts, start_key, size, signature)
}

// Scan is a paid mutator transaction binding the contract method 0x6531695d.
//
// Solidity: function scan(start_key uint256, size uint256, signature uint256) returns()
func (_IOHeavy *IOHeavyTransactor) Scan(opts *bind.TransactOpts, start_key *big.Int, size *big.Int, signature *big.Int) (*types.Transaction, error) {
	return _IOHeavy.contract.Transact(opts, "scan", start_key, size, signature)
}

// Scan is a paid mutator transaction binding the contract method 0x6531695d.
//
// Solidity: function scan(start_key uint256, size uint256, signature uint256) returns()
func (_IOHeavy *IOHeavySession) Scan(start_key *big.Int, size *big.Int, signature *big.Int) (*types.Transaction, error) {
	return _IOHeavy.Contract.Scan(&_IOHeavy.TransactOpts, start_key, size, signature)
}

// Scan is a paid mutator transaction binding the contract method 0x6531695d.
//
// Solidity: function scan(start_key uint256, size uint256, signature uint256) returns()
func (_IOHeavy *IOHeavyTransactorSession) Scan(start_key *big.Int, size *big.Int, signature *big.Int) (*types.Transaction, error) {
	return _IOHeavy.Contract.Scan(&_IOHeavy.TransactOpts, start_key, size, signature)
}

// Set is a paid mutator transaction binding the contract method 0xd778e2da.
//
// Solidity: function set(key bytes20, value bytes) returns()
func (_IOHeavy *IOHeavyTransactor) Set(opts *bind.TransactOpts, key [20]byte, value []byte) (*types.Transaction, error) {
	return _IOHeavy.contract.Transact(opts, "set", key, value)
}

// Set is a paid mutator transaction binding the contract method 0xd778e2da.
//
// Solidity: function set(key bytes20, value bytes) returns()
func (_IOHeavy *IOHeavySession) Set(key [20]byte, value []byte) (*types.Transaction, error) {
	return _IOHeavy.Contract.Set(&_IOHeavy.TransactOpts, key, value)
}

// Set is a paid mutator transaction binding the contract method 0xd778e2da.
//
// Solidity: function set(key bytes20, value bytes) returns()
func (_IOHeavy *IOHeavyTransactorSession) Set(key [20]byte, value []byte) (*types.Transaction, error) {
	return _IOHeavy.Contract.Set(&_IOHeavy.TransactOpts, key, value)
}

// Write is a paid mutator transaction binding the contract method 0xd4cd8790.
//
// Solidity: function write(start_key uint256, size uint256, signature uint256) returns()
func (_IOHeavy *IOHeavyTransactor) Write(opts *bind.TransactOpts, start_key *big.Int, size *big.Int, signature *big.Int) (*types.Transaction, error) {
	return _IOHeavy.contract.Transact(opts, "write", start_key, size, signature)
}

// Write is a paid mutator transaction binding the contract method 0xd4cd8790.
//
// Solidity: function write(start_key uint256, size uint256, signature uint256) returns()
func (_IOHeavy *IOHeavySession) Write(start_key *big.Int, size *big.Int, signature *big.Int) (*types.Transaction, error) {
	return _IOHeavy.Contract.Write(&_IOHeavy.TransactOpts, start_key, size, signature)
}

// Write is a paid mutator transaction binding the contract method 0xd4cd8790.
//
// Solidity: function write(start_key uint256, size uint256, signature uint256) returns()
func (_IOHeavy *IOHeavyTransactorSession) Write(start_key *big.Int, size *big.Int, signature *big.Int) (*types.Transaction, error) {
	return _IOHeavy.Contract.Write(&_IOHeavy.TransactOpts, start_key, size, signature)
}

// IOHeavyFinishScanIterator is returned from FilterFinishScan and is used to iterate over the raw logs and unpacked data for FinishScan events raised by the IOHeavy contract.
type IOHeavyFinishScanIterator struct {
	Event *IOHeavyFinishScan // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IOHeavyFinishScanIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IOHeavyFinishScan)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IOHeavyFinishScan)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IOHeavyFinishScanIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IOHeavyFinishScanIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IOHeavyFinishScan represents a FinishScan event raised by the IOHeavy contract.
type IOHeavyFinishScan struct {
	Size      *big.Int
	Signature *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFinishScan is a free log retrieval operation binding the contract event 0x2e8128137e55a67bef5f6fa7e5c6722c5632e21b8c8bcf6df64bc32239dd6a3f.
//
// Solidity: e finishScan(size uint256, signature uint256)
func (_IOHeavy *IOHeavyFilterer) FilterFinishScan(opts *bind.FilterOpts) (*IOHeavyFinishScanIterator, error) {

	logs, sub, err := _IOHeavy.contract.FilterLogs(opts, "finishScan")
	if err != nil {
		return nil, err
	}
	return &IOHeavyFinishScanIterator{contract: _IOHeavy.contract, event: "finishScan", logs: logs, sub: sub}, nil
}

// WatchFinishScan is a free log subscription operation binding the contract event 0x2e8128137e55a67bef5f6fa7e5c6722c5632e21b8c8bcf6df64bc32239dd6a3f.
//
// Solidity: e finishScan(size uint256, signature uint256)
func (_IOHeavy *IOHeavyFilterer) WatchFinishScan(opts *bind.WatchOpts, sink chan<- *IOHeavyFinishScan) (event.Subscription, error) {

	logs, sub, err := _IOHeavy.contract.WatchLogs(opts, "finishScan")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IOHeavyFinishScan)
				if err := _IOHeavy.contract.UnpackLog(event, "finishScan", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// IOHeavyFinishWriteIterator is returned from FilterFinishWrite and is used to iterate over the raw logs and unpacked data for FinishWrite events raised by the IOHeavy contract.
type IOHeavyFinishWriteIterator struct {
	Event *IOHeavyFinishWrite // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IOHeavyFinishWriteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IOHeavyFinishWrite)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IOHeavyFinishWrite)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IOHeavyFinishWriteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IOHeavyFinishWriteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IOHeavyFinishWrite represents a FinishWrite event raised by the IOHeavy contract.
type IOHeavyFinishWrite struct {
	Size      *big.Int
	Signature *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFinishWrite is a free log retrieval operation binding the contract event 0xe849f68c74be0ec2d162615e7bc539b752b8e3e7db7ccb69f93eb19c85597f7e.
//
// Solidity: e finishWrite(size uint256, signature uint256)
func (_IOHeavy *IOHeavyFilterer) FilterFinishWrite(opts *bind.FilterOpts) (*IOHeavyFinishWriteIterator, error) {

	logs, sub, err := _IOHeavy.contract.FilterLogs(opts, "finishWrite")
	if err != nil {
		return nil, err
	}
	return &IOHeavyFinishWriteIterator{contract: _IOHeavy.contract, event: "finishWrite", logs: logs, sub: sub}, nil
}

// WatchFinishWrite is a free log subscription operation binding the contract event 0xe849f68c74be0ec2d162615e7bc539b752b8e3e7db7ccb69f93eb19c85597f7e.
//
// Solidity: e finishWrite(size uint256, signature uint256)
func (_IOHeavy *IOHeavyFilterer) WatchFinishWrite(opts *bind.WatchOpts, sink chan<- *IOHeavyFinishWrite) (event.Subscription, error) {

	logs, sub, err := _IOHeavy.contract.WatchLogs(opts, "finishWrite")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IOHeavyFinishWrite)
				if err := _IOHeavy.contract.UnpackLog(event, "finishWrite", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
