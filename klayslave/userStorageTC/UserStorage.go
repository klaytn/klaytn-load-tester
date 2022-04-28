// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package userStorageTC

import (
	"math/big"
	"strings"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

// UserStorageABI is the input ABI used to generate the binding from.
const UserStorageABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"set\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserData\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// UserStorageBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const UserStorageBinRuntime = `0x6080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b1811461005b5780636d4ce63c14610075578063ffc9896b1461009c575b600080fd5b34801561006757600080fd5b506100736004356100ca565b005b34801561008157600080fd5b5061008a6100dc565b60408051918252519081900360200190f35b3480156100a857600080fd5b5061008a73ffffffffffffffffffffffffffffffffffffffff600435166100ef565b33600090815260208190526040902055565b3360009081526020819052604090205490565b73ffffffffffffffffffffffffffffffffffffffff16600090815260208190526040902054905600a165627a7a72305820d95af2eed2aa7dc5e38559eb77cc9bd509fc3dc8bab26cc725e6f0beccac956c0029`

// UserStorageBin is the compiled bytecode used for deploying new contracts.
const UserStorageBin = `0x608060405234801561001057600080fd5b50610143806100206000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b1811461005b5780636d4ce63c14610075578063ffc9896b1461009c575b600080fd5b34801561006757600080fd5b506100736004356100ca565b005b34801561008157600080fd5b5061008a6100dc565b60408051918252519081900360200190f35b3480156100a857600080fd5b5061008a73ffffffffffffffffffffffffffffffffffffffff600435166100ef565b33600090815260208190526040902055565b3360009081526020819052604090205490565b73ffffffffffffffffffffffffffffffffffffffff16600090815260208190526040902054905600a165627a7a72305820d95af2eed2aa7dc5e38559eb77cc9bd509fc3dc8bab26cc725e6f0beccac956c0029`

// DeployUserStorage deploys a new GXP contract, binding an instance of UserStorage to it.
func DeployUserStorage(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UserStorage, error) {
	parsed, err := abi.JSON(strings.NewReader(UserStorageABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(UserStorageBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UserStorage{UserStorageCaller: UserStorageCaller{contract: contract}, UserStorageTransactor: UserStorageTransactor{contract: contract}, UserStorageFilterer: UserStorageFilterer{contract: contract}}, nil
}

// UserStorage is an auto generated Go binding around an GXP contract.
type UserStorage struct {
	UserStorageCaller     // Read-only binding to the contract
	UserStorageTransactor // Write-only binding to the contract
	UserStorageFilterer   // Log filterer for contract events
}

// UserStorageCaller is an auto generated read-only Go binding around an GXP contract.
type UserStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UserStorageTransactor is an auto generated write-only Go binding around an GXP contract.
type UserStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UserStorageFilterer is an auto generated log filtering Go binding around an GXP contract events.
type UserStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UserStorageSession is an auto generated Go binding around an GXP contract,
// with pre-set call and transact options.
type UserStorageSession struct {
	Contract     *UserStorage      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UserStorageCallerSession is an auto generated read-only Go binding around an GXP contract,
// with pre-set call options.
type UserStorageCallerSession struct {
	Contract *UserStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// UserStorageTransactorSession is an auto generated write-only Go binding around an GXP contract,
// with pre-set transact options.
type UserStorageTransactorSession struct {
	Contract     *UserStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// UserStorageRaw is an auto generated low-level Go binding around an GXP contract.
type UserStorageRaw struct {
	Contract *UserStorage // Generic contract binding to access the raw methods on
}

// UserStorageCallerRaw is an auto generated low-level read-only Go binding around an GXP contract.
type UserStorageCallerRaw struct {
	Contract *UserStorageCaller // Generic read-only contract binding to access the raw methods on
}

// UserStorageTransactorRaw is an auto generated low-level write-only Go binding around an GXP contract.
type UserStorageTransactorRaw struct {
	Contract *UserStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUserStorage creates a new instance of UserStorage, bound to a specific deployed contract.
func NewUserStorage(address common.Address, backend bind.ContractBackend) (*UserStorage, error) {
	contract, err := bindUserStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UserStorage{UserStorageCaller: UserStorageCaller{contract: contract}, UserStorageTransactor: UserStorageTransactor{contract: contract}, UserStorageFilterer: UserStorageFilterer{contract: contract}}, nil
}

// NewUserStorageCaller creates a new read-only instance of UserStorage, bound to a specific deployed contract.
func NewUserStorageCaller(address common.Address, caller bind.ContractCaller) (*UserStorageCaller, error) {
	contract, err := bindUserStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UserStorageCaller{contract: contract}, nil
}

// NewUserStorageTransactor creates a new write-only instance of UserStorage, bound to a specific deployed contract.
func NewUserStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*UserStorageTransactor, error) {
	contract, err := bindUserStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UserStorageTransactor{contract: contract}, nil
}

// NewUserStorageFilterer creates a new log filterer instance of UserStorage, bound to a specific deployed contract.
func NewUserStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*UserStorageFilterer, error) {
	contract, err := bindUserStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UserStorageFilterer{contract: contract}, nil
}

// bindUserStorage binds a generic wrapper to an already deployed contract.
func bindUserStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UserStorageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UserStorage *UserStorageRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _UserStorage.Contract.UserStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UserStorage *UserStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UserStorage.Contract.UserStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UserStorage *UserStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UserStorage.Contract.UserStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UserStorage *UserStorageCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _UserStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UserStorage *UserStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UserStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UserStorage *UserStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UserStorage.Contract.contract.Transact(opts, method, params...)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() constant returns(uint256)
func (_UserStorage *UserStorageCaller) Get(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _UserStorage.contract.Call(opts, out, "get")
	return *ret0, err
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() constant returns(uint256)
func (_UserStorage *UserStorageSession) Get() (*big.Int, error) {
	return _UserStorage.Contract.Get(&_UserStorage.CallOpts)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() constant returns(uint256)
func (_UserStorage *UserStorageCallerSession) Get() (*big.Int, error) {
	return _UserStorage.Contract.Get(&_UserStorage.CallOpts)
}

// GetUserData is a free data retrieval call binding the contract method 0xffc9896b.
//
// Solidity: function getUserData(user address) constant returns(uint256)
func (_UserStorage *UserStorageCaller) GetUserData(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _UserStorage.contract.Call(opts, out, "getUserData", user)
	return *ret0, err
}

// GetUserData is a free data retrieval call binding the contract method 0xffc9896b.
//
// Solidity: function getUserData(user address) constant returns(uint256)
func (_UserStorage *UserStorageSession) GetUserData(user common.Address) (*big.Int, error) {
	return _UserStorage.Contract.GetUserData(&_UserStorage.CallOpts, user)
}

// GetUserData is a free data retrieval call binding the contract method 0xffc9896b.
//
// Solidity: function getUserData(user address) constant returns(uint256)
func (_UserStorage *UserStorageCallerSession) GetUserData(user common.Address) (*big.Int, error) {
	return _UserStorage.Contract.GetUserData(&_UserStorage.CallOpts, user)
}

// Set is a paid mutator transaction binding the contract method 0x60fe47b1.
//
// Solidity: function set(x uint256) returns()
func (_UserStorage *UserStorageTransactor) Set(opts *bind.TransactOpts, x *big.Int) (*types.Transaction, error) {
	return _UserStorage.contract.Transact(opts, "set", x)
}

// Set is a paid mutator transaction binding the contract method 0x60fe47b1.
//
// Solidity: function set(x uint256) returns()
func (_UserStorage *UserStorageSession) Set(x *big.Int) (*types.Transaction, error) {
	return _UserStorage.Contract.Set(&_UserStorage.TransactOpts, x)
}

// Set is a paid mutator transaction binding the contract method 0x60fe47b1.
//
// Solidity: function set(x uint256) returns()
func (_UserStorage *UserStorageTransactorSession) Set(x *big.Int) (*types.Transaction, error) {
	return _UserStorage.Contract.Set(&_UserStorage.TransactOpts, x)
}
