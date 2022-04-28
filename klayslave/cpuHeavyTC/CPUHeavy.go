// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package cpuHeavyTC

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

// CPUHeavyABI is the input ABI used to generate the binding from.
const CPUHeavyABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"size\",\"type\":\"uint256\"},{\"name\":\"signature\",\"type\":\"uint256\"}],\"name\":\"sort\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"checkResult\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"signature\",\"type\":\"uint256\"}],\"name\":\"sortSingle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"empty\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"size\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"signature\",\"type\":\"uint256\"}],\"name\":\"finish\",\"type\":\"event\"}]"

// CPUHeavyBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const CPUHeavyBinRuntime = `0x6080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416637b395ec28114610066578063a21d942f14610083578063e71c6c82146100ac578063f2a75fe4146100c4575b600080fd5b34801561007257600080fd5b506100816004356024356100d9565b005b34801561008f57600080fd5b5061009861018c565b604080519115158252519081900360200190f35b3480156100b857600080fd5b506100816004356101d9565b3480156100d057600080fd5b5061008161024d565b6060600083604051908082528060200260200182016040528015610107578160200160208202803883390190505b509150600090505b815181101561013b57808403828281518110151561012957fe5b6020908102909101015260010161010f565b61014b8260006001855103610258565b604080518581526020810185905281517fd596fdad182d29130ce218f4c1590c4b5ede105bee36690727baa6592bd2bfc8929181900390910190a150505050565b6001805460009182905b60148110156101ce57600181601481106101ac57fe5b01549150818311156101c157600093506101d3565b9091508190600101610196565b600193505b50505090565b601460005b601481101561020157808203600182601481106101f757fe5b01556001016101de565b61020d60006013610407565b604080518381526020810185905281517fd596fdad182d29130ce218f4c1590c4b5ede105bee36690727baa6592bd2bfc8929181900390910190a1505050565b600080546001019055565b6000806000838510156103ff5750839150829050815b81831015610377575b858181518110151561028557fe5b90602001906020020151868481518110151561029d57fe5b90602001906020020151111580156102b457508383105b156102c457600190920191610277565b85818151811015156102d257fe5b9060200190602002015186838151811015156102ea57fe5b90602001906020020151111561030657600019909101906102c4565b8183101561037257858281518110151561031c57fe5b90602001906020020151868481518110151561033457fe5b90602001906020020151878581518110151561034c57fe5b906020019060200201888581518110151561036357fe5b60209081029091010191909152525b61026e565b858281518110151561038557fe5b90602001906020020151868281518110151561039d57fe5b9060200190602002015187838151811015156103b557fe5b90602001906020020188858151811015156103cc57fe5b602090810290910101919091525260018211156103f1576103f1868660018503610258565b6103ff868360010186610258565b505050505050565b60008060008385101561054c5750839150829050815b818310156104e4575b6001816014811061043357fe5b01546001846014811061044257fe5b01541115801561045157508383105b1561046157600190920191610426565b6001816014811061046e57fe5b01546001836014811061047d57fe5b015411156104915760001990910190610461565b818310156104df57600182601481106104a657fe5b0154600184601481106104b557fe5b0154600185601481106104c457fe5b016000600186601481106104d457fe5b019290925591909155505b61041d565b600182601481106104f157fe5b01546001826014811061050057fe5b01546001836014811061050f57fe5b0160006001866014811061051f57fe5b01929092559190915550600182111561053f5761053f8560018403610407565b61054c8260010185610407565b50505050505600a165627a7a72305820da2b6c258fcb8ad97e0dc6bdf373257eb28dfc7b5ba8cfdf1860152c23429ac70029`

// CPUHeavyBin is the compiled bytecode used for deploying new contracts.
const CPUHeavyBin = `0x608060405234801561001057600080fd5b5061057f806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416637b395ec28114610066578063a21d942f14610083578063e71c6c82146100ac578063f2a75fe4146100c4575b600080fd5b34801561007257600080fd5b506100816004356024356100d9565b005b34801561008f57600080fd5b5061009861018c565b604080519115158252519081900360200190f35b3480156100b857600080fd5b506100816004356101d9565b3480156100d057600080fd5b5061008161024d565b6060600083604051908082528060200260200182016040528015610107578160200160208202803883390190505b509150600090505b815181101561013b57808403828281518110151561012957fe5b6020908102909101015260010161010f565b61014b8260006001855103610258565b604080518581526020810185905281517fd596fdad182d29130ce218f4c1590c4b5ede105bee36690727baa6592bd2bfc8929181900390910190a150505050565b6001805460009182905b60148110156101ce57600181601481106101ac57fe5b01549150818311156101c157600093506101d3565b9091508190600101610196565b600193505b50505090565b601460005b601481101561020157808203600182601481106101f757fe5b01556001016101de565b61020d60006013610407565b604080518381526020810185905281517fd596fdad182d29130ce218f4c1590c4b5ede105bee36690727baa6592bd2bfc8929181900390910190a1505050565b600080546001019055565b6000806000838510156103ff5750839150829050815b81831015610377575b858181518110151561028557fe5b90602001906020020151868481518110151561029d57fe5b90602001906020020151111580156102b457508383105b156102c457600190920191610277565b85818151811015156102d257fe5b9060200190602002015186838151811015156102ea57fe5b90602001906020020151111561030657600019909101906102c4565b8183101561037257858281518110151561031c57fe5b90602001906020020151868481518110151561033457fe5b90602001906020020151878581518110151561034c57fe5b906020019060200201888581518110151561036357fe5b60209081029091010191909152525b61026e565b858281518110151561038557fe5b90602001906020020151868281518110151561039d57fe5b9060200190602002015187838151811015156103b557fe5b90602001906020020188858151811015156103cc57fe5b602090810290910101919091525260018211156103f1576103f1868660018503610258565b6103ff868360010186610258565b505050505050565b60008060008385101561054c5750839150829050815b818310156104e4575b6001816014811061043357fe5b01546001846014811061044257fe5b01541115801561045157508383105b1561046157600190920191610426565b6001816014811061046e57fe5b01546001836014811061047d57fe5b015411156104915760001990910190610461565b818310156104df57600182601481106104a657fe5b0154600184601481106104b557fe5b0154600185601481106104c457fe5b016000600186601481106104d457fe5b019290925591909155505b61041d565b600182601481106104f157fe5b01546001826014811061050057fe5b01546001836014811061050f57fe5b0160006001866014811061051f57fe5b01929092559190915550600182111561053f5761053f8560018403610407565b61054c8260010185610407565b50505050505600a165627a7a72305820da2b6c258fcb8ad97e0dc6bdf373257eb28dfc7b5ba8cfdf1860152c23429ac70029`

// DeployCPUHeavy deploys a new Klaytn contract, binding an instance of CPUHeavy to it.
func DeployCPUHeavy(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CPUHeavy, error) {
	parsed, err := abi.JSON(strings.NewReader(CPUHeavyABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CPUHeavyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CPUHeavy{CPUHeavyCaller: CPUHeavyCaller{contract: contract}, CPUHeavyTransactor: CPUHeavyTransactor{contract: contract}, CPUHeavyFilterer: CPUHeavyFilterer{contract: contract}}, nil
}

// CPUHeavy is an auto generated Go binding around a Klaytn contract.
type CPUHeavy struct {
	CPUHeavyCaller     // Read-only binding to the contract
	CPUHeavyTransactor // Write-only binding to the contract
	CPUHeavyFilterer   // Log filterer for contract events
}

// CPUHeavyCaller is an auto generated read-only Go binding around a Klaytn contract.
type CPUHeavyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CPUHeavyTransactor is an auto generated write-only Go binding around a Klaytn contract.
type CPUHeavyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CPUHeavyFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type CPUHeavyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CPUHeavySession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type CPUHeavySession struct {
	Contract     *CPUHeavy         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CPUHeavyCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type CPUHeavyCallerSession struct {
	Contract *CPUHeavyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// CPUHeavyTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type CPUHeavyTransactorSession struct {
	Contract     *CPUHeavyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// CPUHeavyRaw is an auto generated low-level Go binding around a Klaytn contract.
type CPUHeavyRaw struct {
	Contract *CPUHeavy // Generic contract binding to access the raw methods on
}

// CPUHeavyCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type CPUHeavyCallerRaw struct {
	Contract *CPUHeavyCaller // Generic read-only contract binding to access the raw methods on
}

// CPUHeavyTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type CPUHeavyTransactorRaw struct {
	Contract *CPUHeavyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCPUHeavy creates a new instance of CPUHeavy, bound to a specific deployed contract.
func NewCPUHeavy(address common.Address, backend bind.ContractBackend) (*CPUHeavy, error) {
	contract, err := bindCPUHeavy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CPUHeavy{CPUHeavyCaller: CPUHeavyCaller{contract: contract}, CPUHeavyTransactor: CPUHeavyTransactor{contract: contract}, CPUHeavyFilterer: CPUHeavyFilterer{contract: contract}}, nil
}

// NewCPUHeavyCaller creates a new read-only instance of CPUHeavy, bound to a specific deployed contract.
func NewCPUHeavyCaller(address common.Address, caller bind.ContractCaller) (*CPUHeavyCaller, error) {
	contract, err := bindCPUHeavy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CPUHeavyCaller{contract: contract}, nil
}

// NewCPUHeavyTransactor creates a new write-only instance of CPUHeavy, bound to a specific deployed contract.
func NewCPUHeavyTransactor(address common.Address, transactor bind.ContractTransactor) (*CPUHeavyTransactor, error) {
	contract, err := bindCPUHeavy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CPUHeavyTransactor{contract: contract}, nil
}

// NewCPUHeavyFilterer creates a new log filterer instance of CPUHeavy, bound to a specific deployed contract.
func NewCPUHeavyFilterer(address common.Address, filterer bind.ContractFilterer) (*CPUHeavyFilterer, error) {
	contract, err := bindCPUHeavy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CPUHeavyFilterer{contract: contract}, nil
}

// bindCPUHeavy binds a generic wrapper to an already deployed contract.
func bindCPUHeavy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CPUHeavyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CPUHeavy *CPUHeavyRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CPUHeavy.Contract.CPUHeavyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CPUHeavy *CPUHeavyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CPUHeavy.Contract.CPUHeavyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CPUHeavy *CPUHeavyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CPUHeavy.Contract.CPUHeavyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CPUHeavy *CPUHeavyCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CPUHeavy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CPUHeavy *CPUHeavyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CPUHeavy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CPUHeavy *CPUHeavyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CPUHeavy.Contract.contract.Transact(opts, method, params...)
}

// CheckResult is a free data retrieval call binding the contract method 0xa21d942f.
//
// Solidity: function checkResult() constant returns(bool)
func (_CPUHeavy *CPUHeavyCaller) CheckResult(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CPUHeavy.contract.Call(opts, out, "checkResult")
	return *ret0, err
}

// CheckResult is a free data retrieval call binding the contract method 0xa21d942f.
//
// Solidity: function checkResult() constant returns(bool)
func (_CPUHeavy *CPUHeavySession) CheckResult() (bool, error) {
	return _CPUHeavy.Contract.CheckResult(&_CPUHeavy.CallOpts)
}

// CheckResult is a free data retrieval call binding the contract method 0xa21d942f.
//
// Solidity: function checkResult() constant returns(bool)
func (_CPUHeavy *CPUHeavyCallerSession) CheckResult() (bool, error) {
	return _CPUHeavy.Contract.CheckResult(&_CPUHeavy.CallOpts)
}

// Empty is a paid mutator transaction binding the contract method 0xf2a75fe4.
//
// Solidity: function empty() returns()
func (_CPUHeavy *CPUHeavyTransactor) Empty(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CPUHeavy.contract.Transact(opts, "empty")
}

// Empty is a paid mutator transaction binding the contract method 0xf2a75fe4.
//
// Solidity: function empty() returns()
func (_CPUHeavy *CPUHeavySession) Empty() (*types.Transaction, error) {
	return _CPUHeavy.Contract.Empty(&_CPUHeavy.TransactOpts)
}

// Empty is a paid mutator transaction binding the contract method 0xf2a75fe4.
//
// Solidity: function empty() returns()
func (_CPUHeavy *CPUHeavyTransactorSession) Empty() (*types.Transaction, error) {
	return _CPUHeavy.Contract.Empty(&_CPUHeavy.TransactOpts)
}

// Sort is a paid mutator transaction binding the contract method 0x7b395ec2.
//
// Solidity: function sort(size uint256, signature uint256) returns()
func (_CPUHeavy *CPUHeavyTransactor) Sort(opts *bind.TransactOpts, size *big.Int, signature *big.Int) (*types.Transaction, error) {
	return _CPUHeavy.contract.Transact(opts, "sort", size, signature)
}

// Sort is a paid mutator transaction binding the contract method 0x7b395ec2.
//
// Solidity: function sort(size uint256, signature uint256) returns()
func (_CPUHeavy *CPUHeavySession) Sort(size *big.Int, signature *big.Int) (*types.Transaction, error) {
	return _CPUHeavy.Contract.Sort(&_CPUHeavy.TransactOpts, size, signature)
}

// Sort is a paid mutator transaction binding the contract method 0x7b395ec2.
//
// Solidity: function sort(size uint256, signature uint256) returns()
func (_CPUHeavy *CPUHeavyTransactorSession) Sort(size *big.Int, signature *big.Int) (*types.Transaction, error) {
	return _CPUHeavy.Contract.Sort(&_CPUHeavy.TransactOpts, size, signature)
}

// SortSingle is a paid mutator transaction binding the contract method 0xe71c6c82.
//
// Solidity: function sortSingle(signature uint256) returns()
func (_CPUHeavy *CPUHeavyTransactor) SortSingle(opts *bind.TransactOpts, signature *big.Int) (*types.Transaction, error) {
	return _CPUHeavy.contract.Transact(opts, "sortSingle", signature)
}

// SortSingle is a paid mutator transaction binding the contract method 0xe71c6c82.
//
// Solidity: function sortSingle(signature uint256) returns()
func (_CPUHeavy *CPUHeavySession) SortSingle(signature *big.Int) (*types.Transaction, error) {
	return _CPUHeavy.Contract.SortSingle(&_CPUHeavy.TransactOpts, signature)
}

// SortSingle is a paid mutator transaction binding the contract method 0xe71c6c82.
//
// Solidity: function sortSingle(signature uint256) returns()
func (_CPUHeavy *CPUHeavyTransactorSession) SortSingle(signature *big.Int) (*types.Transaction, error) {
	return _CPUHeavy.Contract.SortSingle(&_CPUHeavy.TransactOpts, signature)
}

// CPUHeavyFinishIterator is returned from FilterFinish and is used to iterate over the raw logs and unpacked data for Finish events raised by the CPUHeavy contract.
type CPUHeavyFinishIterator struct {
	Event *CPUHeavyFinish // Event containing the contract specifics and raw log

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
func (it *CPUHeavyFinishIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CPUHeavyFinish)
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
		it.Event = new(CPUHeavyFinish)
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
func (it *CPUHeavyFinishIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CPUHeavyFinishIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CPUHeavyFinish represents a Finish event raised by the CPUHeavy contract.
type CPUHeavyFinish struct {
	Size      *big.Int
	Signature *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFinish is a free log retrieval operation binding the contract event 0xd596fdad182d29130ce218f4c1590c4b5ede105bee36690727baa6592bd2bfc8.
//
// Solidity: e finish(size uint256, signature uint256)
func (_CPUHeavy *CPUHeavyFilterer) FilterFinish(opts *bind.FilterOpts) (*CPUHeavyFinishIterator, error) {

	logs, sub, err := _CPUHeavy.contract.FilterLogs(opts, "finish")
	if err != nil {
		return nil, err
	}
	return &CPUHeavyFinishIterator{contract: _CPUHeavy.contract, event: "finish", logs: logs, sub: sub}, nil
}

// WatchFinish is a free log subscription operation binding the contract event 0xd596fdad182d29130ce218f4c1590c4b5ede105bee36690727baa6592bd2bfc8.
//
// Solidity: e finish(size uint256, signature uint256)
func (_CPUHeavy *CPUHeavyFilterer) WatchFinish(opts *bind.WatchOpts, sink chan<- *CPUHeavyFinish) (event.Subscription, error) {

	logs, sub, err := _CPUHeavy.contract.WatchLogs(opts, "finish")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CPUHeavyFinish)
				if err := _CPUHeavy.contract.UnpackLog(event, "finish", log); err != nil {
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
