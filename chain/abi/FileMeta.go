// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// FileMetaFileMetaData is an auto generated low-level Go binding around an user-defined struct.
type FileMetaFileMetaData struct {
	OwnerPeerId string
	From        common.Address
	FileName    string
	FileExt     string
	IsDir       bool
	FileSize    *big.Int
}

// FileMetaMetaData contains all meta data concerning the FileMeta contract.
var FileMetaMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"ownerPeerId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"fileName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"fileExt\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"isDir\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"fileSize\",\"type\":\"uint256\"}],\"internalType\":\"structFileMeta.FileMetaData\",\"name\":\"metaData\",\"type\":\"tuple\"}],\"name\":\"AddFileMeta\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"}],\"name\":\"DeleteFileMeta\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"ownerPeerId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"fileName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"fileExt\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"isDir\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"fileSize\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structFileMeta.FileMetaData\",\"name\":\"metaData\",\"type\":\"tuple\"}],\"name\":\"eventAddFileMeta\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"}],\"name\":\"eventDeleteFileMeta\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"}],\"name\":\"GetFileMeta\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"ownerPeerId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"fileName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"fileExt\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"isDir\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"fileSize\",\"type\":\"uint256\"}],\"internalType\":\"structFileMeta.FileMetaData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getImplementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// FileMetaABI is the input ABI used to generate the binding from.
// Deprecated: Use FileMetaMetaData.ABI instead.
var FileMetaABI = FileMetaMetaData.ABI

// FileMeta is an auto generated Go binding around an Ethereum contract.
type FileMeta struct {
	FileMetaCaller     // Read-only binding to the contract
	FileMetaTransactor // Write-only binding to the contract
	FileMetaFilterer   // Log filterer for contract events
}

// FileMetaCaller is an auto generated read-only Go binding around an Ethereum contract.
type FileMetaCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileMetaTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FileMetaTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileMetaFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FileMetaFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileMetaSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FileMetaSession struct {
	Contract     *FileMeta         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FileMetaCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FileMetaCallerSession struct {
	Contract *FileMetaCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// FileMetaTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FileMetaTransactorSession struct {
	Contract     *FileMetaTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// FileMetaRaw is an auto generated low-level Go binding around an Ethereum contract.
type FileMetaRaw struct {
	Contract *FileMeta // Generic contract binding to access the raw methods on
}

// FileMetaCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FileMetaCallerRaw struct {
	Contract *FileMetaCaller // Generic read-only contract binding to access the raw methods on
}

// FileMetaTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FileMetaTransactorRaw struct {
	Contract *FileMetaTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFileMeta creates a new instance of FileMeta, bound to a specific deployed contract.
func NewFileMeta(address common.Address, backend bind.ContractBackend) (*FileMeta, error) {
	contract, err := bindFileMeta(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FileMeta{FileMetaCaller: FileMetaCaller{contract: contract}, FileMetaTransactor: FileMetaTransactor{contract: contract}, FileMetaFilterer: FileMetaFilterer{contract: contract}}, nil
}

// NewFileMetaCaller creates a new read-only instance of FileMeta, bound to a specific deployed contract.
func NewFileMetaCaller(address common.Address, caller bind.ContractCaller) (*FileMetaCaller, error) {
	contract, err := bindFileMeta(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FileMetaCaller{contract: contract}, nil
}

// NewFileMetaTransactor creates a new write-only instance of FileMeta, bound to a specific deployed contract.
func NewFileMetaTransactor(address common.Address, transactor bind.ContractTransactor) (*FileMetaTransactor, error) {
	contract, err := bindFileMeta(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FileMetaTransactor{contract: contract}, nil
}

// NewFileMetaFilterer creates a new log filterer instance of FileMeta, bound to a specific deployed contract.
func NewFileMetaFilterer(address common.Address, filterer bind.ContractFilterer) (*FileMetaFilterer, error) {
	contract, err := bindFileMeta(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FileMetaFilterer{contract: contract}, nil
}

// bindFileMeta binds a generic wrapper to an already deployed contract.
func bindFileMeta(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FileMetaMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FileMeta *FileMetaRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FileMeta.Contract.FileMetaCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FileMeta *FileMetaRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileMeta.Contract.FileMetaTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FileMeta *FileMetaRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FileMeta.Contract.FileMetaTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FileMeta *FileMetaCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FileMeta.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FileMeta *FileMetaTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileMeta.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FileMeta *FileMetaTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FileMeta.Contract.contract.Transact(opts, method, params...)
}

// GetFileMeta is a free data retrieval call binding the contract method 0x7d21ecd0.
//
// Solidity: function GetFileMeta(string cid) view returns((string,address,string,string,bool,uint256))
func (_FileMeta *FileMetaCaller) GetFileMeta(opts *bind.CallOpts, cid string) (FileMetaFileMetaData, error) {
	var out []interface{}
	err := _FileMeta.contract.Call(opts, &out, "GetFileMeta", cid)

	if err != nil {
		return *new(FileMetaFileMetaData), err
	}

	out0 := *abi.ConvertType(out[0], new(FileMetaFileMetaData)).(*FileMetaFileMetaData)

	return out0, err

}

// GetFileMeta is a free data retrieval call binding the contract method 0x7d21ecd0.
//
// Solidity: function GetFileMeta(string cid) view returns((string,address,string,string,bool,uint256))
func (_FileMeta *FileMetaSession) GetFileMeta(cid string) (FileMetaFileMetaData, error) {
	return _FileMeta.Contract.GetFileMeta(&_FileMeta.CallOpts, cid)
}

// GetFileMeta is a free data retrieval call binding the contract method 0x7d21ecd0.
//
// Solidity: function GetFileMeta(string cid) view returns((string,address,string,string,bool,uint256))
func (_FileMeta *FileMetaCallerSession) GetFileMeta(cid string) (FileMetaFileMetaData, error) {
	return _FileMeta.Contract.GetFileMeta(&_FileMeta.CallOpts, cid)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_FileMeta *FileMetaCaller) GetImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FileMeta.contract.Call(opts, &out, "getImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_FileMeta *FileMetaSession) GetImplementation() (common.Address, error) {
	return _FileMeta.Contract.GetImplementation(&_FileMeta.CallOpts)
}

// GetImplementation is a free data retrieval call binding the contract method 0xaaf10f42.
//
// Solidity: function getImplementation() view returns(address)
func (_FileMeta *FileMetaCallerSession) GetImplementation() (common.Address, error) {
	return _FileMeta.Contract.GetImplementation(&_FileMeta.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FileMeta *FileMetaCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FileMeta.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FileMeta *FileMetaSession) Owner() (common.Address, error) {
	return _FileMeta.Contract.Owner(&_FileMeta.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FileMeta *FileMetaCallerSession) Owner() (common.Address, error) {
	return _FileMeta.Contract.Owner(&_FileMeta.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FileMeta *FileMetaCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FileMeta.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FileMeta *FileMetaSession) ProxiableUUID() ([32]byte, error) {
	return _FileMeta.Contract.ProxiableUUID(&_FileMeta.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FileMeta *FileMetaCallerSession) ProxiableUUID() ([32]byte, error) {
	return _FileMeta.Contract.ProxiableUUID(&_FileMeta.CallOpts)
}

// AddFileMeta is a paid mutator transaction binding the contract method 0x4592e045.
//
// Solidity: function AddFileMeta(string cid, (string,address,string,string,bool,uint256) metaData) returns()
func (_FileMeta *FileMetaTransactor) AddFileMeta(opts *bind.TransactOpts, cid string, metaData FileMetaFileMetaData) (*types.Transaction, error) {
	return _FileMeta.contract.Transact(opts, "AddFileMeta", cid, metaData)
}

// AddFileMeta is a paid mutator transaction binding the contract method 0x4592e045.
//
// Solidity: function AddFileMeta(string cid, (string,address,string,string,bool,uint256) metaData) returns()
func (_FileMeta *FileMetaSession) AddFileMeta(cid string, metaData FileMetaFileMetaData) (*types.Transaction, error) {
	return _FileMeta.Contract.AddFileMeta(&_FileMeta.TransactOpts, cid, metaData)
}

// AddFileMeta is a paid mutator transaction binding the contract method 0x4592e045.
//
// Solidity: function AddFileMeta(string cid, (string,address,string,string,bool,uint256) metaData) returns()
func (_FileMeta *FileMetaTransactorSession) AddFileMeta(cid string, metaData FileMetaFileMetaData) (*types.Transaction, error) {
	return _FileMeta.Contract.AddFileMeta(&_FileMeta.TransactOpts, cid, metaData)
}

// DeleteFileMeta is a paid mutator transaction binding the contract method 0x91ec67dc.
//
// Solidity: function DeleteFileMeta(string cid) returns()
func (_FileMeta *FileMetaTransactor) DeleteFileMeta(opts *bind.TransactOpts, cid string) (*types.Transaction, error) {
	return _FileMeta.contract.Transact(opts, "DeleteFileMeta", cid)
}

// DeleteFileMeta is a paid mutator transaction binding the contract method 0x91ec67dc.
//
// Solidity: function DeleteFileMeta(string cid) returns()
func (_FileMeta *FileMetaSession) DeleteFileMeta(cid string) (*types.Transaction, error) {
	return _FileMeta.Contract.DeleteFileMeta(&_FileMeta.TransactOpts, cid)
}

// DeleteFileMeta is a paid mutator transaction binding the contract method 0x91ec67dc.
//
// Solidity: function DeleteFileMeta(string cid) returns()
func (_FileMeta *FileMetaTransactorSession) DeleteFileMeta(cid string) (*types.Transaction, error) {
	return _FileMeta.Contract.DeleteFileMeta(&_FileMeta.TransactOpts, cid)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_FileMeta *FileMetaTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileMeta.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_FileMeta *FileMetaSession) Initialize() (*types.Transaction, error) {
	return _FileMeta.Contract.Initialize(&_FileMeta.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_FileMeta *FileMetaTransactorSession) Initialize() (*types.Transaction, error) {
	return _FileMeta.Contract.Initialize(&_FileMeta.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FileMeta *FileMetaTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileMeta.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FileMeta *FileMetaSession) RenounceOwnership() (*types.Transaction, error) {
	return _FileMeta.Contract.RenounceOwnership(&_FileMeta.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FileMeta *FileMetaTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FileMeta.Contract.RenounceOwnership(&_FileMeta.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FileMeta *FileMetaTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FileMeta.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FileMeta *FileMetaSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FileMeta.Contract.TransferOwnership(&_FileMeta.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FileMeta *FileMetaTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FileMeta.Contract.TransferOwnership(&_FileMeta.TransactOpts, newOwner)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FileMeta *FileMetaTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _FileMeta.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FileMeta *FileMetaSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _FileMeta.Contract.UpgradeTo(&_FileMeta.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FileMeta *FileMetaTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _FileMeta.Contract.UpgradeTo(&_FileMeta.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FileMeta *FileMetaTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FileMeta.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FileMeta *FileMetaSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FileMeta.Contract.UpgradeToAndCall(&_FileMeta.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FileMeta *FileMetaTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FileMeta.Contract.UpgradeToAndCall(&_FileMeta.TransactOpts, newImplementation, data)
}

// FileMetaAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the FileMeta contract.
type FileMetaAdminChangedIterator struct {
	Event *FileMetaAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileMetaAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileMetaAdminChanged)
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
		it.Event = new(FileMetaAdminChanged)
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
func (it *FileMetaAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileMetaAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileMetaAdminChanged represents a AdminChanged event raised by the FileMeta contract.
type FileMetaAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_FileMeta *FileMetaFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*FileMetaAdminChangedIterator, error) {

	logs, sub, err := _FileMeta.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &FileMetaAdminChangedIterator{contract: _FileMeta.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_FileMeta *FileMetaFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *FileMetaAdminChanged) (event.Subscription, error) {

	logs, sub, err := _FileMeta.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileMetaAdminChanged)
				if err := _FileMeta.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_FileMeta *FileMetaFilterer) ParseAdminChanged(log types.Log) (*FileMetaAdminChanged, error) {
	event := new(FileMetaAdminChanged)
	if err := _FileMeta.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileMetaBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the FileMeta contract.
type FileMetaBeaconUpgradedIterator struct {
	Event *FileMetaBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileMetaBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileMetaBeaconUpgraded)
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
		it.Event = new(FileMetaBeaconUpgraded)
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
func (it *FileMetaBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileMetaBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileMetaBeaconUpgraded represents a BeaconUpgraded event raised by the FileMeta contract.
type FileMetaBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_FileMeta *FileMetaFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*FileMetaBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _FileMeta.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &FileMetaBeaconUpgradedIterator{contract: _FileMeta.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_FileMeta *FileMetaFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *FileMetaBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _FileMeta.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileMetaBeaconUpgraded)
				if err := _FileMeta.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_FileMeta *FileMetaFilterer) ParseBeaconUpgraded(log types.Log) (*FileMetaBeaconUpgraded, error) {
	event := new(FileMetaBeaconUpgraded)
	if err := _FileMeta.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileMetaInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the FileMeta contract.
type FileMetaInitializedIterator struct {
	Event *FileMetaInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileMetaInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileMetaInitialized)
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
		it.Event = new(FileMetaInitialized)
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
func (it *FileMetaInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileMetaInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileMetaInitialized represents a Initialized event raised by the FileMeta contract.
type FileMetaInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_FileMeta *FileMetaFilterer) FilterInitialized(opts *bind.FilterOpts) (*FileMetaInitializedIterator, error) {

	logs, sub, err := _FileMeta.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FileMetaInitializedIterator{contract: _FileMeta.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_FileMeta *FileMetaFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FileMetaInitialized) (event.Subscription, error) {

	logs, sub, err := _FileMeta.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileMetaInitialized)
				if err := _FileMeta.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_FileMeta *FileMetaFilterer) ParseInitialized(log types.Log) (*FileMetaInitialized, error) {
	event := new(FileMetaInitialized)
	if err := _FileMeta.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileMetaOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FileMeta contract.
type FileMetaOwnershipTransferredIterator struct {
	Event *FileMetaOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileMetaOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileMetaOwnershipTransferred)
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
		it.Event = new(FileMetaOwnershipTransferred)
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
func (it *FileMetaOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileMetaOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileMetaOwnershipTransferred represents a OwnershipTransferred event raised by the FileMeta contract.
type FileMetaOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FileMeta *FileMetaFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FileMetaOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FileMeta.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FileMetaOwnershipTransferredIterator{contract: _FileMeta.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FileMeta *FileMetaFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FileMetaOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FileMeta.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileMetaOwnershipTransferred)
				if err := _FileMeta.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FileMeta *FileMetaFilterer) ParseOwnershipTransferred(log types.Log) (*FileMetaOwnershipTransferred, error) {
	event := new(FileMetaOwnershipTransferred)
	if err := _FileMeta.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileMetaUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the FileMeta contract.
type FileMetaUpgradedIterator struct {
	Event *FileMetaUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileMetaUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileMetaUpgraded)
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
		it.Event = new(FileMetaUpgraded)
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
func (it *FileMetaUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileMetaUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileMetaUpgraded represents a Upgraded event raised by the FileMeta contract.
type FileMetaUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FileMeta *FileMetaFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*FileMetaUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FileMeta.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &FileMetaUpgradedIterator{contract: _FileMeta.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FileMeta *FileMetaFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *FileMetaUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FileMeta.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileMetaUpgraded)
				if err := _FileMeta.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FileMeta *FileMetaFilterer) ParseUpgraded(log types.Log) (*FileMetaUpgraded, error) {
	event := new(FileMetaUpgraded)
	if err := _FileMeta.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileMetaEventAddFileMetaIterator is returned from FilterEventAddFileMeta and is used to iterate over the raw logs and unpacked data for EventAddFileMeta events raised by the FileMeta contract.
type FileMetaEventAddFileMetaIterator struct {
	Event *FileMetaEventAddFileMeta // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileMetaEventAddFileMetaIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileMetaEventAddFileMeta)
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
		it.Event = new(FileMetaEventAddFileMeta)
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
func (it *FileMetaEventAddFileMetaIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileMetaEventAddFileMetaIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileMetaEventAddFileMeta represents a EventAddFileMeta event raised by the FileMeta contract.
type FileMetaEventAddFileMeta struct {
	Cid      string
	MetaData FileMetaFileMetaData
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterEventAddFileMeta is a free log retrieval operation binding the contract event 0x5858f0194a72c3a15e5d5f94dbbb35a90af87476e958f990a6fc31db87a16b95.
//
// Solidity: event eventAddFileMeta(string cid, (string,address,string,string,bool,uint256) metaData)
func (_FileMeta *FileMetaFilterer) FilterEventAddFileMeta(opts *bind.FilterOpts) (*FileMetaEventAddFileMetaIterator, error) {

	logs, sub, err := _FileMeta.contract.FilterLogs(opts, "eventAddFileMeta")
	if err != nil {
		return nil, err
	}
	return &FileMetaEventAddFileMetaIterator{contract: _FileMeta.contract, event: "eventAddFileMeta", logs: logs, sub: sub}, nil
}

// WatchEventAddFileMeta is a free log subscription operation binding the contract event 0x5858f0194a72c3a15e5d5f94dbbb35a90af87476e958f990a6fc31db87a16b95.
//
// Solidity: event eventAddFileMeta(string cid, (string,address,string,string,bool,uint256) metaData)
func (_FileMeta *FileMetaFilterer) WatchEventAddFileMeta(opts *bind.WatchOpts, sink chan<- *FileMetaEventAddFileMeta) (event.Subscription, error) {

	logs, sub, err := _FileMeta.contract.WatchLogs(opts, "eventAddFileMeta")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileMetaEventAddFileMeta)
				if err := _FileMeta.contract.UnpackLog(event, "eventAddFileMeta", log); err != nil {
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

// ParseEventAddFileMeta is a log parse operation binding the contract event 0x5858f0194a72c3a15e5d5f94dbbb35a90af87476e958f990a6fc31db87a16b95.
//
// Solidity: event eventAddFileMeta(string cid, (string,address,string,string,bool,uint256) metaData)
func (_FileMeta *FileMetaFilterer) ParseEventAddFileMeta(log types.Log) (*FileMetaEventAddFileMeta, error) {
	event := new(FileMetaEventAddFileMeta)
	if err := _FileMeta.contract.UnpackLog(event, "eventAddFileMeta", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileMetaEventDeleteFileMetaIterator is returned from FilterEventDeleteFileMeta and is used to iterate over the raw logs and unpacked data for EventDeleteFileMeta events raised by the FileMeta contract.
type FileMetaEventDeleteFileMetaIterator struct {
	Event *FileMetaEventDeleteFileMeta // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FileMetaEventDeleteFileMetaIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileMetaEventDeleteFileMeta)
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
		it.Event = new(FileMetaEventDeleteFileMeta)
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
func (it *FileMetaEventDeleteFileMetaIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileMetaEventDeleteFileMetaIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileMetaEventDeleteFileMeta represents a EventDeleteFileMeta event raised by the FileMeta contract.
type FileMetaEventDeleteFileMeta struct {
	Cid string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEventDeleteFileMeta is a free log retrieval operation binding the contract event 0x75ab1296ebe483fe9506fccbfe8a4743edf9ef6cf8802c2c28c652dfa948e412.
//
// Solidity: event eventDeleteFileMeta(string cid)
func (_FileMeta *FileMetaFilterer) FilterEventDeleteFileMeta(opts *bind.FilterOpts) (*FileMetaEventDeleteFileMetaIterator, error) {

	logs, sub, err := _FileMeta.contract.FilterLogs(opts, "eventDeleteFileMeta")
	if err != nil {
		return nil, err
	}
	return &FileMetaEventDeleteFileMetaIterator{contract: _FileMeta.contract, event: "eventDeleteFileMeta", logs: logs, sub: sub}, nil
}

// WatchEventDeleteFileMeta is a free log subscription operation binding the contract event 0x75ab1296ebe483fe9506fccbfe8a4743edf9ef6cf8802c2c28c652dfa948e412.
//
// Solidity: event eventDeleteFileMeta(string cid)
func (_FileMeta *FileMetaFilterer) WatchEventDeleteFileMeta(opts *bind.WatchOpts, sink chan<- *FileMetaEventDeleteFileMeta) (event.Subscription, error) {

	logs, sub, err := _FileMeta.contract.WatchLogs(opts, "eventDeleteFileMeta")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileMetaEventDeleteFileMeta)
				if err := _FileMeta.contract.UnpackLog(event, "eventDeleteFileMeta", log); err != nil {
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

// ParseEventDeleteFileMeta is a log parse operation binding the contract event 0x75ab1296ebe483fe9506fccbfe8a4743edf9ef6cf8802c2c28c652dfa948e412.
//
// Solidity: event eventDeleteFileMeta(string cid)
func (_FileMeta *FileMetaFilterer) ParseEventDeleteFileMeta(log types.Log) (*FileMetaEventDeleteFileMeta, error) {
	event := new(FileMetaEventDeleteFileMeta)
	if err := _FileMeta.contract.UnpackLog(event, "eventDeleteFileMeta", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
