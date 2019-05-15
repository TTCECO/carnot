// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	"github.com/TTCECO/gttc/accounts/abi"
	"github.com/TTCECO/gttc/accounts/abi/bind"
	"github.com/TTCECO/gttc/common"
	"github.com/TTCECO/gttc/core/types"
)

// ContractABI is the input ABI used to generate the binding from.
const ContractABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"delValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_num\",\"type\":\"uint256\"}],\"name\":\"setMinConfirmNum\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"setDepositFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"addValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_id\",\"type\":\"string\"},{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getConfirmStatus\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"depositFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_id\",\"type\":\"string\"},{\"name\":\"_target\",\"type\":\"address\"},{\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"confirm\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"depositRecords\",\"outputs\":[{\"name\":\"orderID\",\"type\":\"string\"},{\"name\":\"targetAddress\",\"type\":\"address\"},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\"},{\"name\":\"confirmCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minConfirmNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() constant returns(uint256)
func (_Contract *ContractCaller) DepositFee(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "depositFee")
	return *ret0, err
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() constant returns(uint256)
func (_Contract *ContractSession) DepositFee() (*big.Int, error) {
	return _Contract.Contract.DepositFee(&_Contract.CallOpts)
}

// DepositFee is a free data retrieval call binding the contract method 0x67a52793.
//
// Solidity: function depositFee() constant returns(uint256)
func (_Contract *ContractCallerSession) DepositFee() (*big.Int, error) {
	return _Contract.Contract.DepositFee(&_Contract.CallOpts)
}

// DepositRecords is a free data retrieval call binding the contract method 0xcf50cb6b.
//
// Solidity: function depositRecords( bytes32) constant returns(orderID string, targetAddress address, name string, value uint256, status uint8, confirmCount uint256)
func (_Contract *ContractCaller) DepositRecords(opts *bind.CallOpts, arg0 [32]byte) (struct {
	OrderID       string
	TargetAddress common.Address
	Name          string
	Value         *big.Int
	Status        uint8
	ConfirmCount  *big.Int
}, error) {
	ret := new(struct {
		OrderID       string
		TargetAddress common.Address
		Name          string
		Value         *big.Int
		Status        uint8
		ConfirmCount  *big.Int
	})
	out := ret
	err := _Contract.contract.Call(opts, out, "depositRecords", arg0)
	return *ret, err
}

// DepositRecords is a free data retrieval call binding the contract method 0xcf50cb6b.
//
// Solidity: function depositRecords( bytes32) constant returns(orderID string, targetAddress address, name string, value uint256, status uint8, confirmCount uint256)
func (_Contract *ContractSession) DepositRecords(arg0 [32]byte) (struct {
	OrderID       string
	TargetAddress common.Address
	Name          string
	Value         *big.Int
	Status        uint8
	ConfirmCount  *big.Int
}, error) {
	return _Contract.Contract.DepositRecords(&_Contract.CallOpts, arg0)
}

// DepositRecords is a free data retrieval call binding the contract method 0xcf50cb6b.
//
// Solidity: function depositRecords( bytes32) constant returns(orderID string, targetAddress address, name string, value uint256, status uint8, confirmCount uint256)
func (_Contract *ContractCallerSession) DepositRecords(arg0 [32]byte) (struct {
	OrderID       string
	TargetAddress common.Address
	Name          string
	Value         *big.Int
	Status        uint8
	ConfirmCount  *big.Int
}, error) {
	return _Contract.Contract.DepositRecords(&_Contract.CallOpts, arg0)
}

// GetConfirmStatus is a free data retrieval call binding the contract method 0x57ab0d0a.
//
// Solidity: function getConfirmStatus(_id string, _addr address) constant returns(bool)
func (_Contract *ContractCaller) GetConfirmStatus(opts *bind.CallOpts, _id string, _addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "getConfirmStatus", _id, _addr)
	return *ret0, err
}

// GetConfirmStatus is a free data retrieval call binding the contract method 0x57ab0d0a.
//
// Solidity: function getConfirmStatus(_id string, _addr address) constant returns(bool)
func (_Contract *ContractSession) GetConfirmStatus(_id string, _addr common.Address) (bool, error) {
	return _Contract.Contract.GetConfirmStatus(&_Contract.CallOpts, _id, _addr)
}

// GetConfirmStatus is a free data retrieval call binding the contract method 0x57ab0d0a.
//
// Solidity: function getConfirmStatus(_id string, _addr address) constant returns(bool)
func (_Contract *ContractCallerSession) GetConfirmStatus(_id string, _addr common.Address) (bool, error) {
	return _Contract.Contract.GetConfirmStatus(&_Contract.CallOpts, _id, _addr)
}

// MinConfirmNum is a free data retrieval call binding the contract method 0xffa5a563.
//
// Solidity: function minConfirmNum() constant returns(uint256)
func (_Contract *ContractCaller) MinConfirmNum(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "minConfirmNum")
	return *ret0, err
}

// MinConfirmNum is a free data retrieval call binding the contract method 0xffa5a563.
//
// Solidity: function minConfirmNum() constant returns(uint256)
func (_Contract *ContractSession) MinConfirmNum() (*big.Int, error) {
	return _Contract.Contract.MinConfirmNum(&_Contract.CallOpts)
}

// MinConfirmNum is a free data retrieval call binding the contract method 0xffa5a563.
//
// Solidity: function minConfirmNum() constant returns(uint256)
func (_Contract *ContractCallerSession) MinConfirmNum() (*big.Int, error) {
	return _Contract.Contract.MinConfirmNum(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Contract *ContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Contract *ContractSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Contract *ContractCallerSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators( address) constant returns(bool)
func (_Contract *ContractCaller) Validators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "validators", arg0)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators( address) constant returns(bool)
func (_Contract *ContractSession) Validators(arg0 common.Address) (bool, error) {
	return _Contract.Contract.Validators(&_Contract.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators( address) constant returns(bool)
func (_Contract *ContractCallerSession) Validators(arg0 common.Address) (bool, error) {
	return _Contract.Contract.Validators(&_Contract.CallOpts, arg0)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(_addr address) returns()
func (_Contract *ContractTransactor) AddValidator(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "addValidator", _addr)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(_addr address) returns()
func (_Contract *ContractSession) AddValidator(_addr common.Address) (*types.Transaction, error) {
	return _Contract.Contract.AddValidator(&_Contract.TransactOpts, _addr)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(_addr address) returns()
func (_Contract *ContractTransactorSession) AddValidator(_addr common.Address) (*types.Transaction, error) {
	return _Contract.Contract.AddValidator(&_Contract.TransactOpts, _addr)
}

// Confirm is a paid mutator transaction binding the contract method 0xc41c79c3.
//
// Solidity: function confirm(_id string, _target address, _name string, _value uint256) returns()
func (_Contract *ContractTransactor) Confirm(opts *bind.TransactOpts, _id string, _target common.Address, _name string, _value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "confirm", _id, _target, _name, _value)
}

// Confirm is a paid mutator transaction binding the contract method 0xc41c79c3.
//
// Solidity: function confirm(_id string, _target address, _name string, _value uint256) returns()
func (_Contract *ContractSession) Confirm(_id string, _target common.Address, _name string, _value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Confirm(&_Contract.TransactOpts, _id, _target, _name, _value)
}

// Confirm is a paid mutator transaction binding the contract method 0xc41c79c3.
//
// Solidity: function confirm(_id string, _target address, _name string, _value uint256) returns()
func (_Contract *ContractTransactorSession) Confirm(_id string, _target common.Address, _name string, _value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Confirm(&_Contract.TransactOpts, _id, _target, _name, _value)
}

// DelValidator is a paid mutator transaction binding the contract method 0x12ae2c65.
//
// Solidity: function delValidator(_addr address) returns()
func (_Contract *ContractTransactor) DelValidator(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "delValidator", _addr)
}

// DelValidator is a paid mutator transaction binding the contract method 0x12ae2c65.
//
// Solidity: function delValidator(_addr address) returns()
func (_Contract *ContractSession) DelValidator(_addr common.Address) (*types.Transaction, error) {
	return _Contract.Contract.DelValidator(&_Contract.TransactOpts, _addr)
}

// DelValidator is a paid mutator transaction binding the contract method 0x12ae2c65.
//
// Solidity: function delValidator(_addr address) returns()
func (_Contract *ContractTransactorSession) DelValidator(_addr common.Address) (*types.Transaction, error) {
	return _Contract.Contract.DelValidator(&_Contract.TransactOpts, _addr)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns()
func (_Contract *ContractTransactor) Finalize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "finalize")
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns()
func (_Contract *ContractSession) Finalize() (*types.Transaction, error) {
	return _Contract.Contract.Finalize(&_Contract.TransactOpts)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns()
func (_Contract *ContractTransactorSession) Finalize() (*types.Transaction, error) {
	return _Contract.Contract.Finalize(&_Contract.TransactOpts)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(_fee uint256) returns()
func (_Contract *ContractTransactor) SetDepositFee(opts *bind.TransactOpts, _fee *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setDepositFee", _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(_fee uint256) returns()
func (_Contract *ContractSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetDepositFee(&_Contract.TransactOpts, _fee)
}

// SetDepositFee is a paid mutator transaction binding the contract method 0x490ae210.
//
// Solidity: function setDepositFee(_fee uint256) returns()
func (_Contract *ContractTransactorSession) SetDepositFee(_fee *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetDepositFee(&_Contract.TransactOpts, _fee)
}

// SetMinConfirmNum is a paid mutator transaction binding the contract method 0x1b783cb0.
//
// Solidity: function setMinConfirmNum(_num uint256) returns()
func (_Contract *ContractTransactor) SetMinConfirmNum(opts *bind.TransactOpts, _num *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setMinConfirmNum", _num)
}

// SetMinConfirmNum is a paid mutator transaction binding the contract method 0x1b783cb0.
//
// Solidity: function setMinConfirmNum(_num uint256) returns()
func (_Contract *ContractSession) SetMinConfirmNum(_num *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetMinConfirmNum(&_Contract.TransactOpts, _num)
}

// SetMinConfirmNum is a paid mutator transaction binding the contract method 0x1b783cb0.
//
// Solidity: function setMinConfirmNum(_num uint256) returns()
func (_Contract *ContractTransactorSession) SetMinConfirmNum(_num *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetMinConfirmNum(&_Contract.TransactOpts, _num)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Contract *ContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Contract *ContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Contract *ContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}
