// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package smallBankTC

import (
	"math/big"
	"strings"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

// SmallBankABI is the input ABI used to generate the binding from.
const SmallBankABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"arg0\",\"type\":\"string\"},{\"name\":\"arg1\",\"type\":\"uint256\"}],\"name\":\"updateSaving\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"arg0\",\"type\":\"string\"},{\"name\":\"arg1\",\"type\":\"uint256\"}],\"name\":\"writeCheck\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"arg0\",\"type\":\"string\"}],\"name\":\"getBalance\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"arg0\",\"type\":\"string\"},{\"name\":\"arg1\",\"type\":\"uint256\"}],\"name\":\"updateBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"arg0\",\"type\":\"string\"},{\"name\":\"arg1\",\"type\":\"string\"}],\"name\":\"almagate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"arg0\",\"type\":\"string\"},{\"name\":\"arg1\",\"type\":\"string\"},{\"name\":\"arg2\",\"type\":\"uint256\"}],\"name\":\"sendPayment\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// SmallBankBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const SmallBankBinRuntime = `0x6080604052600436106100775763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630b488b37811461007c5780630be8374d146100d95780633a51d24614610134578063870187eb1461019f578063901d706f146101fa578063ca30543514610291575b600080fd5b34801561008857600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100d7943694929360249392840191908190840183828082843750949750509335945061032a9350505050565b005b3480156100e557600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100d794369492936024939284019190819084018382808284375094975050933594506103ff9350505050565b34801561014057600080fd5b506040805160206004803580820135601f810184900484028501840190955284845261018d9436949293602493928401919081908401838280828437509497506105b19650505050505050565b60408051918252519081900360200190f35b3480156101ab57600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100d794369492936024939284019190819084018382808284375094975050933594506106819350505050565b34801561020657600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100d794369492936024939284019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497506107209650505050505050565b34801561029d57600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100d794369492936024939284019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a99988101979196509182019450925082915084018382808284375094975050933594506108839350505050565b6000806000846040518082805190602001908083835b6020831061035f5780518252601f199092019160209182019101610340565b51815160209384036101000a600019018019909216911617905292019485525060405193849003810184205488519096508795508587019460009450899350918291908401908083835b602083106103c85780518252601f1990920191602091820191016103a9565b51815160209384036101000a6000190180199092169116179052920194855250604051938490030190922092909255505050505050565b60008060006001856040518082805190602001908083835b602083106104365780518252601f199092019160209182019101610417565b51815160209384036101000a600019018019909216911617905292019485525060405193849003810184205489519097506000948a9450925082918401908083835b602083106104975780518252601f199092019160209182019101610478565b51815160209384036101000a600019018019909216911617905292019485525060405193849003019092205493508592505050828201811015610543576001818403036001866040518082805190602001908083835b6020831061050c5780518252601f1990920191602091820191016104ed565b51815160209384036101000a6000190180199092169116179052920194855250604051938490030190922092909255506105aa9050565b8083036001866040518082805190602001908083835b602083106105785780518252601f199092019160209182019101610559565b51815160209384036101000a600019018019909216911617905292019485525060405193849003019092209290925550505b5050505050565b600080600080846040518082805190602001908083835b602083106105e75780518252601f1990920191602091820191016105c8565b51815160209384036101000a60001901801990921691161790529201948552506040519384900381018420548851909650600194899450925082918401908083835b602083106106485780518252601f199092019160209182019101610629565b51815160209384036101000a60001901801990921691161790529201948552506040519384900301909220549390930195945050505050565b6000806001846040518082805190602001908083835b602083106106b65780518252601f199092019160209182019101610697565b51815160001960209485036101000a019081169019919091161790529201948552506040519384900381018420548851909650879550858701946001945089935091829190840190808383602083106103c85780518252601f1990920191602091820191016103a9565b6000806000846040518082805190602001908083835b602083106107555780518252601f199092019160209182019101610736565b51815160209384036101000a60001901801990921691161790529201948552506040519384900381018420548751909650600194889450925082918401908083835b602083106107b65780518252601f199092019160209182019101610797565b51815160209384036101000a6000190180199092169116179052920194855250604051938490038101842054885190955060009460019450899350918291908401908083835b6020831061081b5780518252601f1990920191602091820191016107fc565b51815160001960209485036101000a01908116901991909116179052920194855250604051938490038101842094909455505084518484019260009287929091829190840190808383602083106103c85780518252601f1990920191602091820191016103a9565b60008060006001866040518082805190602001908083835b602083106108ba5780518252601f19909201916020918201910161089b565b51815160209384036101000a600019018019909216911617905292019485525060405193849003810184205489519097506001948a9450925082918401908083835b6020831061091b5780518252601f1990920191602091820191016108fc565b51815160209384036101000a60001901801990921691161790529201948552506040519384900381018420548a51978990039790890196508895508794600194508b9350918291908401908083835b602083106109895780518252601f19909201916020918201910161096a565b51815160209384036101000a6000190180199092169116179052920194855250604051938490038101842094909455505086518492600192899290918291908401908083835b602083106109ee5780518252601f1990920191602091820191016109cf565b51815160209384036101000a600019018019909216911617905292019485525060405193849003019092209290925550505050505050505600a165627a7a72305820462e7aa872935a5629c50dcc94f5e114c54a38ab3e249fedcaf7ba44c4d39a940029`

// SmallBankBin is the compiled bytecode used for deploying new contracts.
const SmallBankBin = `0x608060405234801561001057600080fd5b50610a53806100206000396000f3006080604052600436106100775763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630b488b37811461007c5780630be8374d146100d95780633a51d24614610134578063870187eb1461019f578063901d706f146101fa578063ca30543514610291575b600080fd5b34801561008857600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100d7943694929360249392840191908190840183828082843750949750509335945061032a9350505050565b005b3480156100e557600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100d794369492936024939284019190819084018382808284375094975050933594506103ff9350505050565b34801561014057600080fd5b506040805160206004803580820135601f810184900484028501840190955284845261018d9436949293602493928401919081908401838280828437509497506105b19650505050505050565b60408051918252519081900360200190f35b3480156101ab57600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100d794369492936024939284019190819084018382808284375094975050933594506106819350505050565b34801561020657600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100d794369492936024939284019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a9998810197919650918201945092508291508401838280828437509497506107209650505050505050565b34801561029d57600080fd5b506040805160206004803580820135601f81018490048402850184019095528484526100d794369492936024939284019190819084018382808284375050604080516020601f89358b018035918201839004830284018301909452808352979a99988101979196509182019450925082915084018382808284375094975050933594506108839350505050565b6000806000846040518082805190602001908083835b6020831061035f5780518252601f199092019160209182019101610340565b51815160209384036101000a600019018019909216911617905292019485525060405193849003810184205488519096508795508587019460009450899350918291908401908083835b602083106103c85780518252601f1990920191602091820191016103a9565b51815160209384036101000a6000190180199092169116179052920194855250604051938490030190922092909255505050505050565b60008060006001856040518082805190602001908083835b602083106104365780518252601f199092019160209182019101610417565b51815160209384036101000a600019018019909216911617905292019485525060405193849003810184205489519097506000948a9450925082918401908083835b602083106104975780518252601f199092019160209182019101610478565b51815160209384036101000a600019018019909216911617905292019485525060405193849003019092205493508592505050828201811015610543576001818403036001866040518082805190602001908083835b6020831061050c5780518252601f1990920191602091820191016104ed565b51815160209384036101000a6000190180199092169116179052920194855250604051938490030190922092909255506105aa9050565b8083036001866040518082805190602001908083835b602083106105785780518252601f199092019160209182019101610559565b51815160209384036101000a600019018019909216911617905292019485525060405193849003019092209290925550505b5050505050565b600080600080846040518082805190602001908083835b602083106105e75780518252601f1990920191602091820191016105c8565b51815160209384036101000a60001901801990921691161790529201948552506040519384900381018420548851909650600194899450925082918401908083835b602083106106485780518252601f199092019160209182019101610629565b51815160209384036101000a60001901801990921691161790529201948552506040519384900301909220549390930195945050505050565b6000806001846040518082805190602001908083835b602083106106b65780518252601f199092019160209182019101610697565b51815160001960209485036101000a019081169019919091161790529201948552506040519384900381018420548851909650879550858701946001945089935091829190840190808383602083106103c85780518252601f1990920191602091820191016103a9565b6000806000846040518082805190602001908083835b602083106107555780518252601f199092019160209182019101610736565b51815160209384036101000a60001901801990921691161790529201948552506040519384900381018420548751909650600194889450925082918401908083835b602083106107b65780518252601f199092019160209182019101610797565b51815160209384036101000a6000190180199092169116179052920194855250604051938490038101842054885190955060009460019450899350918291908401908083835b6020831061081b5780518252601f1990920191602091820191016107fc565b51815160001960209485036101000a01908116901991909116179052920194855250604051938490038101842094909455505084518484019260009287929091829190840190808383602083106103c85780518252601f1990920191602091820191016103a9565b60008060006001866040518082805190602001908083835b602083106108ba5780518252601f19909201916020918201910161089b565b51815160209384036101000a600019018019909216911617905292019485525060405193849003810184205489519097506001948a9450925082918401908083835b6020831061091b5780518252601f1990920191602091820191016108fc565b51815160209384036101000a60001901801990921691161790529201948552506040519384900381018420548a51978990039790890196508895508794600194508b9350918291908401908083835b602083106109895780518252601f19909201916020918201910161096a565b51815160209384036101000a6000190180199092169116179052920194855250604051938490038101842094909455505086518492600192899290918291908401908083835b602083106109ee5780518252601f1990920191602091820191016109cf565b51815160209384036101000a600019018019909216911617905292019485525060405193849003019092209290925550505050505050505600a165627a7a72305820462e7aa872935a5629c50dcc94f5e114c54a38ab3e249fedcaf7ba44c4d39a940029`

// DeploySmallBank deploys a new GXP contract, binding an instance of SmallBank to it.
func DeploySmallBank(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SmallBank, error) {
	parsed, err := abi.JSON(strings.NewReader(SmallBankABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SmallBankBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SmallBank{SmallBankCaller: SmallBankCaller{contract: contract}, SmallBankTransactor: SmallBankTransactor{contract: contract}, SmallBankFilterer: SmallBankFilterer{contract: contract}}, nil
}

// SmallBank is an auto generated Go binding around an GXP contract.
type SmallBank struct {
	SmallBankCaller     // Read-only binding to the contract
	SmallBankTransactor // Write-only binding to the contract
	SmallBankFilterer   // Log filterer for contract events
}

// SmallBankCaller is an auto generated read-only Go binding around an GXP contract.
type SmallBankCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SmallBankTransactor is an auto generated write-only Go binding around an GXP contract.
type SmallBankTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SmallBankFilterer is an auto generated log filtering Go binding around an GXP contract events.
type SmallBankFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SmallBankSession is an auto generated Go binding around an GXP contract,
// with pre-set call and transact options.
type SmallBankSession struct {
	Contract     *SmallBank        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SmallBankCallerSession is an auto generated read-only Go binding around an GXP contract,
// with pre-set call options.
type SmallBankCallerSession struct {
	Contract *SmallBankCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// SmallBankTransactorSession is an auto generated write-only Go binding around an GXP contract,
// with pre-set transact options.
type SmallBankTransactorSession struct {
	Contract     *SmallBankTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// SmallBankRaw is an auto generated low-level Go binding around an GXP contract.
type SmallBankRaw struct {
	Contract *SmallBank // Generic contract binding to access the raw methods on
}

// SmallBankCallerRaw is an auto generated low-level read-only Go binding around an GXP contract.
type SmallBankCallerRaw struct {
	Contract *SmallBankCaller // Generic read-only contract binding to access the raw methods on
}

// SmallBankTransactorRaw is an auto generated low-level write-only Go binding around an GXP contract.
type SmallBankTransactorRaw struct {
	Contract *SmallBankTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSmallBank creates a new instance of SmallBank, bound to a specific deployed contract.
func NewSmallBank(address common.Address, backend bind.ContractBackend) (*SmallBank, error) {
	contract, err := bindSmallBank(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SmallBank{SmallBankCaller: SmallBankCaller{contract: contract}, SmallBankTransactor: SmallBankTransactor{contract: contract}, SmallBankFilterer: SmallBankFilterer{contract: contract}}, nil
}

// NewSmallBankCaller creates a new read-only instance of SmallBank, bound to a specific deployed contract.
func NewSmallBankCaller(address common.Address, caller bind.ContractCaller) (*SmallBankCaller, error) {
	contract, err := bindSmallBank(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SmallBankCaller{contract: contract}, nil
}

// NewSmallBankTransactor creates a new write-only instance of SmallBank, bound to a specific deployed contract.
func NewSmallBankTransactor(address common.Address, transactor bind.ContractTransactor) (*SmallBankTransactor, error) {
	contract, err := bindSmallBank(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SmallBankTransactor{contract: contract}, nil
}

// NewSmallBankFilterer creates a new log filterer instance of SmallBank, bound to a specific deployed contract.
func NewSmallBankFilterer(address common.Address, filterer bind.ContractFilterer) (*SmallBankFilterer, error) {
	contract, err := bindSmallBank(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SmallBankFilterer{contract: contract}, nil
}

// bindSmallBank binds a generic wrapper to an already deployed contract.
func bindSmallBank(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SmallBankABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SmallBank *SmallBankRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SmallBank.Contract.SmallBankCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SmallBank *SmallBankRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SmallBank.Contract.SmallBankTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SmallBank *SmallBankRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SmallBank.Contract.SmallBankTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SmallBank *SmallBankCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SmallBank.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SmallBank *SmallBankTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SmallBank.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SmallBank *SmallBankTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SmallBank.Contract.contract.Transact(opts, method, params...)
}

// GetBalance is a free data retrieval call binding the contract method 0x3a51d246.
//
// Solidity: function getBalance(arg0 string) constant returns(balance uint256)
func (_SmallBank *SmallBankCaller) GetBalance(opts *bind.CallOpts, arg0 string) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SmallBank.contract.Call(opts, out, "getBalance", arg0)
	return *ret0, err
}

// GetBalance is a free data retrieval call binding the contract method 0x3a51d246.
//
// Solidity: function getBalance(arg0 string) constant returns(balance uint256)
func (_SmallBank *SmallBankSession) GetBalance(arg0 string) (*big.Int, error) {
	return _SmallBank.Contract.GetBalance(&_SmallBank.CallOpts, arg0)
}

// GetBalance is a free data retrieval call binding the contract method 0x3a51d246.
//
// Solidity: function getBalance(arg0 string) constant returns(balance uint256)
func (_SmallBank *SmallBankCallerSession) GetBalance(arg0 string) (*big.Int, error) {
	return _SmallBank.Contract.GetBalance(&_SmallBank.CallOpts, arg0)
}

// Almagate is a paid mutator transaction binding the contract method 0x901d706f.
//
// Solidity: function almagate(arg0 string, arg1 string) returns()
func (_SmallBank *SmallBankTransactor) Almagate(opts *bind.TransactOpts, arg0 string, arg1 string) (*types.Transaction, error) {
	return _SmallBank.contract.Transact(opts, "almagate", arg0, arg1)
}

// Almagate is a paid mutator transaction binding the contract method 0x901d706f.
//
// Solidity: function almagate(arg0 string, arg1 string) returns()
func (_SmallBank *SmallBankSession) Almagate(arg0 string, arg1 string) (*types.Transaction, error) {
	return _SmallBank.Contract.Almagate(&_SmallBank.TransactOpts, arg0, arg1)
}

// Almagate is a paid mutator transaction binding the contract method 0x901d706f.
//
// Solidity: function almagate(arg0 string, arg1 string) returns()
func (_SmallBank *SmallBankTransactorSession) Almagate(arg0 string, arg1 string) (*types.Transaction, error) {
	return _SmallBank.Contract.Almagate(&_SmallBank.TransactOpts, arg0, arg1)
}

// SendPayment is a paid mutator transaction binding the contract method 0xca305435.
//
// Solidity: function sendPayment(arg0 string, arg1 string, arg2 uint256) returns()
func (_SmallBank *SmallBankTransactor) SendPayment(opts *bind.TransactOpts, arg0 string, arg1 string, arg2 *big.Int) (*types.Transaction, error) {
	return _SmallBank.contract.Transact(opts, "sendPayment", arg0, arg1, arg2)
}

// SendPayment is a paid mutator transaction binding the contract method 0xca305435.
//
// Solidity: function sendPayment(arg0 string, arg1 string, arg2 uint256) returns()
func (_SmallBank *SmallBankSession) SendPayment(arg0 string, arg1 string, arg2 *big.Int) (*types.Transaction, error) {
	return _SmallBank.Contract.SendPayment(&_SmallBank.TransactOpts, arg0, arg1, arg2)
}

// SendPayment is a paid mutator transaction binding the contract method 0xca305435.
//
// Solidity: function sendPayment(arg0 string, arg1 string, arg2 uint256) returns()
func (_SmallBank *SmallBankTransactorSession) SendPayment(arg0 string, arg1 string, arg2 *big.Int) (*types.Transaction, error) {
	return _SmallBank.Contract.SendPayment(&_SmallBank.TransactOpts, arg0, arg1, arg2)
}

// UpdateBalance is a paid mutator transaction binding the contract method 0x870187eb.
//
// Solidity: function updateBalance(arg0 string, arg1 uint256) returns()
func (_SmallBank *SmallBankTransactor) UpdateBalance(opts *bind.TransactOpts, arg0 string, arg1 *big.Int) (*types.Transaction, error) {
	return _SmallBank.contract.Transact(opts, "updateBalance", arg0, arg1)
}

// UpdateBalance is a paid mutator transaction binding the contract method 0x870187eb.
//
// Solidity: function updateBalance(arg0 string, arg1 uint256) returns()
func (_SmallBank *SmallBankSession) UpdateBalance(arg0 string, arg1 *big.Int) (*types.Transaction, error) {
	return _SmallBank.Contract.UpdateBalance(&_SmallBank.TransactOpts, arg0, arg1)
}

// UpdateBalance is a paid mutator transaction binding the contract method 0x870187eb.
//
// Solidity: function updateBalance(arg0 string, arg1 uint256) returns()
func (_SmallBank *SmallBankTransactorSession) UpdateBalance(arg0 string, arg1 *big.Int) (*types.Transaction, error) {
	return _SmallBank.Contract.UpdateBalance(&_SmallBank.TransactOpts, arg0, arg1)
}

// UpdateSaving is a paid mutator transaction binding the contract method 0x0b488b37.
//
// Solidity: function updateSaving(arg0 string, arg1 uint256) returns()
func (_SmallBank *SmallBankTransactor) UpdateSaving(opts *bind.TransactOpts, arg0 string, arg1 *big.Int) (*types.Transaction, error) {
	return _SmallBank.contract.Transact(opts, "updateSaving", arg0, arg1)
}

// UpdateSaving is a paid mutator transaction binding the contract method 0x0b488b37.
//
// Solidity: function updateSaving(arg0 string, arg1 uint256) returns()
func (_SmallBank *SmallBankSession) UpdateSaving(arg0 string, arg1 *big.Int) (*types.Transaction, error) {
	return _SmallBank.Contract.UpdateSaving(&_SmallBank.TransactOpts, arg0, arg1)
}

// UpdateSaving is a paid mutator transaction binding the contract method 0x0b488b37.
//
// Solidity: function updateSaving(arg0 string, arg1 uint256) returns()
func (_SmallBank *SmallBankTransactorSession) UpdateSaving(arg0 string, arg1 *big.Int) (*types.Transaction, error) {
	return _SmallBank.Contract.UpdateSaving(&_SmallBank.TransactOpts, arg0, arg1)
}

// WriteCheck is a paid mutator transaction binding the contract method 0x0be8374d.
//
// Solidity: function writeCheck(arg0 string, arg1 uint256) returns()
func (_SmallBank *SmallBankTransactor) WriteCheck(opts *bind.TransactOpts, arg0 string, arg1 *big.Int) (*types.Transaction, error) {
	return _SmallBank.contract.Transact(opts, "writeCheck", arg0, arg1)
}

// WriteCheck is a paid mutator transaction binding the contract method 0x0be8374d.
//
// Solidity: function writeCheck(arg0 string, arg1 uint256) returns()
func (_SmallBank *SmallBankSession) WriteCheck(arg0 string, arg1 *big.Int) (*types.Transaction, error) {
	return _SmallBank.Contract.WriteCheck(&_SmallBank.TransactOpts, arg0, arg1)
}

// WriteCheck is a paid mutator transaction binding the contract method 0x0be8374d.
//
// Solidity: function writeCheck(arg0 string, arg1 uint256) returns()
func (_SmallBank *SmallBankTransactorSession) WriteCheck(arg0 string, arg1 *big.Int) (*types.Transaction, error) {
	return _SmallBank.Contract.WriteCheck(&_SmallBank.TransactOpts, arg0, arg1)
}
