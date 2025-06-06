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

// FileMetaContractSPPair is an auto generated low-level Go binding around an user-defined struct.
type FileMetaContractSPPair struct {
	ContractId string
	Sp         common.Address
}

// FileMetaContractMetaData contains all meta data concerning the FileMetaContract contract.
var FileMetaContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"metaData\",\"type\":\"bytes\"}],\"name\":\"FileMetaAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"contractId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"StatusUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"metaData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"size\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"contractId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"sp\",\"type\":\"address\"}],\"internalType\":\"structFileMeta.ContractSPPair[]\",\"name\":\"pairs\",\"type\":\"tuple[]\"}],\"name\":\"addFileMeta\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"fileMeta\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"cid\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"contractIds\",\"type\":\"string[]\"}],\"name\":\"getFileMeta\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"metaData\",\"type\":\"bytes\"},{\"internalType\":\"uint8[]\",\"name\":\"statuses\",\"type\":\"uint8[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractId\",\"type\":\"string\"}],\"name\":\"getSP\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sp\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractId\",\"type\":\"string\"}],\"name\":\"getStatus\",\"outputs\":[{\"internalType\":\"enumFileMeta.FileStoreStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalUsedSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractId\",\"type\":\"string\"},{\"internalType\":\"enumFileMeta.FileStoreStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"updateStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// FileMetaContractABI is the input ABI used to generate the binding from.
// Deprecated: Use FileMetaContractMetaData.ABI instead.
var FileMetaContractABI = FileMetaContractMetaData.ABI

// FileMetaContract is an auto generated Go binding around an Ethereum contract.
type FileMetaContract struct {
	FileMetaContractCaller     // Read-only binding to the contract
	FileMetaContractTransactor // Write-only binding to the contract
	FileMetaContractFilterer   // Log filterer for contract events
}

// FileMetaContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type FileMetaContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileMetaContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FileMetaContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileMetaContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FileMetaContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FileMetaContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FileMetaContractSession struct {
	Contract     *FileMetaContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FileMetaContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FileMetaContractCallerSession struct {
	Contract *FileMetaContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// FileMetaContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FileMetaContractTransactorSession struct {
	Contract     *FileMetaContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// FileMetaContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type FileMetaContractRaw struct {
	Contract *FileMetaContract // Generic contract binding to access the raw methods on
}

// FileMetaContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FileMetaContractCallerRaw struct {
	Contract *FileMetaContractCaller // Generic read-only contract binding to access the raw methods on
}

// FileMetaContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FileMetaContractTransactorRaw struct {
	Contract *FileMetaContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFileMetaContract creates a new instance of FileMetaContract, bound to a specific deployed contract.
func NewFileMetaContract(address common.Address, backend bind.ContractBackend) (*FileMetaContract, error) {
	contract, err := bindFileMetaContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FileMetaContract{FileMetaContractCaller: FileMetaContractCaller{contract: contract}, FileMetaContractTransactor: FileMetaContractTransactor{contract: contract}, FileMetaContractFilterer: FileMetaContractFilterer{contract: contract}}, nil
}

// NewFileMetaContractCaller creates a new read-only instance of FileMetaContract, bound to a specific deployed contract.
func NewFileMetaContractCaller(address common.Address, caller bind.ContractCaller) (*FileMetaContractCaller, error) {
	contract, err := bindFileMetaContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FileMetaContractCaller{contract: contract}, nil
}

// NewFileMetaContractTransactor creates a new write-only instance of FileMetaContract, bound to a specific deployed contract.
func NewFileMetaContractTransactor(address common.Address, transactor bind.ContractTransactor) (*FileMetaContractTransactor, error) {
	contract, err := bindFileMetaContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FileMetaContractTransactor{contract: contract}, nil
}

// NewFileMetaContractFilterer creates a new log filterer instance of FileMetaContract, bound to a specific deployed contract.
func NewFileMetaContractFilterer(address common.Address, filterer bind.ContractFilterer) (*FileMetaContractFilterer, error) {
	contract, err := bindFileMetaContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FileMetaContractFilterer{contract: contract}, nil
}

// bindFileMetaContract binds a generic wrapper to an already deployed contract.
func bindFileMetaContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FileMetaContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FileMetaContract *FileMetaContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FileMetaContract.Contract.FileMetaContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FileMetaContract *FileMetaContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileMetaContract.Contract.FileMetaContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FileMetaContract *FileMetaContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FileMetaContract.Contract.FileMetaContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FileMetaContract *FileMetaContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FileMetaContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FileMetaContract *FileMetaContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileMetaContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FileMetaContract *FileMetaContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FileMetaContract.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_FileMetaContract *FileMetaContractCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FileMetaContract.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_FileMetaContract *FileMetaContractSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _FileMetaContract.Contract.UPGRADEINTERFACEVERSION(&_FileMetaContract.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_FileMetaContract *FileMetaContractCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _FileMetaContract.Contract.UPGRADEINTERFACEVERSION(&_FileMetaContract.CallOpts)
}

// FileMeta is a free data retrieval call binding the contract method 0x6fa361a2.
//
// Solidity: function fileMeta(string ) view returns(bytes)
func (_FileMetaContract *FileMetaContractCaller) FileMeta(opts *bind.CallOpts, arg0 string) ([]byte, error) {
	var out []interface{}
	err := _FileMetaContract.contract.Call(opts, &out, "fileMeta", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// FileMeta is a free data retrieval call binding the contract method 0x6fa361a2.
//
// Solidity: function fileMeta(string ) view returns(bytes)
func (_FileMetaContract *FileMetaContractSession) FileMeta(arg0 string) ([]byte, error) {
	return _FileMetaContract.Contract.FileMeta(&_FileMetaContract.CallOpts, arg0)
}

// FileMeta is a free data retrieval call binding the contract method 0x6fa361a2.
//
// Solidity: function fileMeta(string ) view returns(bytes)
func (_FileMetaContract *FileMetaContractCallerSession) FileMeta(arg0 string) ([]byte, error) {
	return _FileMetaContract.Contract.FileMeta(&_FileMetaContract.CallOpts, arg0)
}

// GetFileMeta is a free data retrieval call binding the contract method 0x87f03f0d.
//
// Solidity: function getFileMeta(string cid, string[] contractIds) view returns(bytes metaData, uint8[] statuses)
func (_FileMetaContract *FileMetaContractCaller) GetFileMeta(opts *bind.CallOpts, cid string, contractIds []string) (struct {
	MetaData []byte
	Statuses []uint8
}, error) {
	var out []interface{}
	err := _FileMetaContract.contract.Call(opts, &out, "getFileMeta", cid, contractIds)

	outstruct := new(struct {
		MetaData []byte
		Statuses []uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MetaData = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.Statuses = *abi.ConvertType(out[1], new([]uint8)).(*[]uint8)

	return *outstruct, err

}

// GetFileMeta is a free data retrieval call binding the contract method 0x87f03f0d.
//
// Solidity: function getFileMeta(string cid, string[] contractIds) view returns(bytes metaData, uint8[] statuses)
func (_FileMetaContract *FileMetaContractSession) GetFileMeta(cid string, contractIds []string) (struct {
	MetaData []byte
	Statuses []uint8
}, error) {
	return _FileMetaContract.Contract.GetFileMeta(&_FileMetaContract.CallOpts, cid, contractIds)
}

// GetFileMeta is a free data retrieval call binding the contract method 0x87f03f0d.
//
// Solidity: function getFileMeta(string cid, string[] contractIds) view returns(bytes metaData, uint8[] statuses)
func (_FileMetaContract *FileMetaContractCallerSession) GetFileMeta(cid string, contractIds []string) (struct {
	MetaData []byte
	Statuses []uint8
}, error) {
	return _FileMetaContract.Contract.GetFileMeta(&_FileMetaContract.CallOpts, cid, contractIds)
}

// GetSP is a free data retrieval call binding the contract method 0x0f6e3a4f.
//
// Solidity: function getSP(string contractId) view returns(address sp)
func (_FileMetaContract *FileMetaContractCaller) GetSP(opts *bind.CallOpts, contractId string) (common.Address, error) {
	var out []interface{}
	err := _FileMetaContract.contract.Call(opts, &out, "getSP", contractId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetSP is a free data retrieval call binding the contract method 0x0f6e3a4f.
//
// Solidity: function getSP(string contractId) view returns(address sp)
func (_FileMetaContract *FileMetaContractSession) GetSP(contractId string) (common.Address, error) {
	return _FileMetaContract.Contract.GetSP(&_FileMetaContract.CallOpts, contractId)
}

// GetSP is a free data retrieval call binding the contract method 0x0f6e3a4f.
//
// Solidity: function getSP(string contractId) view returns(address sp)
func (_FileMetaContract *FileMetaContractCallerSession) GetSP(contractId string) (common.Address, error) {
	return _FileMetaContract.Contract.GetSP(&_FileMetaContract.CallOpts, contractId)
}

// GetStatus is a free data retrieval call binding the contract method 0x22b05ed2.
//
// Solidity: function getStatus(string contractId) view returns(uint8 status)
func (_FileMetaContract *FileMetaContractCaller) GetStatus(opts *bind.CallOpts, contractId string) (uint8, error) {
	var out []interface{}
	err := _FileMetaContract.contract.Call(opts, &out, "getStatus", contractId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetStatus is a free data retrieval call binding the contract method 0x22b05ed2.
//
// Solidity: function getStatus(string contractId) view returns(uint8 status)
func (_FileMetaContract *FileMetaContractSession) GetStatus(contractId string) (uint8, error) {
	return _FileMetaContract.Contract.GetStatus(&_FileMetaContract.CallOpts, contractId)
}

// GetStatus is a free data retrieval call binding the contract method 0x22b05ed2.
//
// Solidity: function getStatus(string contractId) view returns(uint8 status)
func (_FileMetaContract *FileMetaContractCallerSession) GetStatus(contractId string) (uint8, error) {
	return _FileMetaContract.Contract.GetStatus(&_FileMetaContract.CallOpts, contractId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FileMetaContract *FileMetaContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FileMetaContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FileMetaContract *FileMetaContractSession) Owner() (common.Address, error) {
	return _FileMetaContract.Contract.Owner(&_FileMetaContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FileMetaContract *FileMetaContractCallerSession) Owner() (common.Address, error) {
	return _FileMetaContract.Contract.Owner(&_FileMetaContract.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FileMetaContract *FileMetaContractCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FileMetaContract.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FileMetaContract *FileMetaContractSession) ProxiableUUID() ([32]byte, error) {
	return _FileMetaContract.Contract.ProxiableUUID(&_FileMetaContract.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FileMetaContract *FileMetaContractCallerSession) ProxiableUUID() ([32]byte, error) {
	return _FileMetaContract.Contract.ProxiableUUID(&_FileMetaContract.CallOpts)
}

// TotalUsedSize is a free data retrieval call binding the contract method 0xd48462b6.
//
// Solidity: function totalUsedSize() view returns(uint256)
func (_FileMetaContract *FileMetaContractCaller) TotalUsedSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FileMetaContract.contract.Call(opts, &out, "totalUsedSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalUsedSize is a free data retrieval call binding the contract method 0xd48462b6.
//
// Solidity: function totalUsedSize() view returns(uint256)
func (_FileMetaContract *FileMetaContractSession) TotalUsedSize() (*big.Int, error) {
	return _FileMetaContract.Contract.TotalUsedSize(&_FileMetaContract.CallOpts)
}

// TotalUsedSize is a free data retrieval call binding the contract method 0xd48462b6.
//
// Solidity: function totalUsedSize() view returns(uint256)
func (_FileMetaContract *FileMetaContractCallerSession) TotalUsedSize() (*big.Int, error) {
	return _FileMetaContract.Contract.TotalUsedSize(&_FileMetaContract.CallOpts)
}

// AddFileMeta is a paid mutator transaction binding the contract method 0xd1f94065.
//
// Solidity: function addFileMeta(string cid, bytes metaData, uint256 size, (string,address)[] pairs) returns()
func (_FileMetaContract *FileMetaContractTransactor) AddFileMeta(opts *bind.TransactOpts, cid string, metaData []byte, size *big.Int, pairs []FileMetaContractSPPair) (*types.Transaction, error) {
	return _FileMetaContract.contract.Transact(opts, "addFileMeta", cid, metaData, size, pairs)
}

// AddFileMeta is a paid mutator transaction binding the contract method 0xd1f94065.
//
// Solidity: function addFileMeta(string cid, bytes metaData, uint256 size, (string,address)[] pairs) returns()
func (_FileMetaContract *FileMetaContractSession) AddFileMeta(cid string, metaData []byte, size *big.Int, pairs []FileMetaContractSPPair) (*types.Transaction, error) {
	return _FileMetaContract.Contract.AddFileMeta(&_FileMetaContract.TransactOpts, cid, metaData, size, pairs)
}

// AddFileMeta is a paid mutator transaction binding the contract method 0xd1f94065.
//
// Solidity: function addFileMeta(string cid, bytes metaData, uint256 size, (string,address)[] pairs) returns()
func (_FileMetaContract *FileMetaContractTransactorSession) AddFileMeta(cid string, metaData []byte, size *big.Int, pairs []FileMetaContractSPPair) (*types.Transaction, error) {
	return _FileMetaContract.Contract.AddFileMeta(&_FileMetaContract.TransactOpts, cid, metaData, size, pairs)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_FileMetaContract *FileMetaContractTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileMetaContract.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_FileMetaContract *FileMetaContractSession) Initialize() (*types.Transaction, error) {
	return _FileMetaContract.Contract.Initialize(&_FileMetaContract.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_FileMetaContract *FileMetaContractTransactorSession) Initialize() (*types.Transaction, error) {
	return _FileMetaContract.Contract.Initialize(&_FileMetaContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FileMetaContract *FileMetaContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FileMetaContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FileMetaContract *FileMetaContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _FileMetaContract.Contract.RenounceOwnership(&_FileMetaContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FileMetaContract *FileMetaContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FileMetaContract.Contract.RenounceOwnership(&_FileMetaContract.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FileMetaContract *FileMetaContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FileMetaContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FileMetaContract *FileMetaContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FileMetaContract.Contract.TransferOwnership(&_FileMetaContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FileMetaContract *FileMetaContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FileMetaContract.Contract.TransferOwnership(&_FileMetaContract.TransactOpts, newOwner)
}

// UpdateStatus is a paid mutator transaction binding the contract method 0x53616982.
//
// Solidity: function updateStatus(string contractId, uint8 status) returns()
func (_FileMetaContract *FileMetaContractTransactor) UpdateStatus(opts *bind.TransactOpts, contractId string, status uint8) (*types.Transaction, error) {
	return _FileMetaContract.contract.Transact(opts, "updateStatus", contractId, status)
}

// UpdateStatus is a paid mutator transaction binding the contract method 0x53616982.
//
// Solidity: function updateStatus(string contractId, uint8 status) returns()
func (_FileMetaContract *FileMetaContractSession) UpdateStatus(contractId string, status uint8) (*types.Transaction, error) {
	return _FileMetaContract.Contract.UpdateStatus(&_FileMetaContract.TransactOpts, contractId, status)
}

// UpdateStatus is a paid mutator transaction binding the contract method 0x53616982.
//
// Solidity: function updateStatus(string contractId, uint8 status) returns()
func (_FileMetaContract *FileMetaContractTransactorSession) UpdateStatus(contractId string, status uint8) (*types.Transaction, error) {
	return _FileMetaContract.Contract.UpdateStatus(&_FileMetaContract.TransactOpts, contractId, status)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FileMetaContract *FileMetaContractTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FileMetaContract.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FileMetaContract *FileMetaContractSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FileMetaContract.Contract.UpgradeToAndCall(&_FileMetaContract.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FileMetaContract *FileMetaContractTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FileMetaContract.Contract.UpgradeToAndCall(&_FileMetaContract.TransactOpts, newImplementation, data)
}

// FileMetaContractFileMetaAddedIterator is returned from FilterFileMetaAdded and is used to iterate over the raw logs and unpacked data for FileMetaAdded events raised by the FileMetaContract contract.
type FileMetaContractFileMetaAddedIterator struct {
	Event *FileMetaContractFileMetaAdded // Event containing the contract specifics and raw log

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
func (it *FileMetaContractFileMetaAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileMetaContractFileMetaAdded)
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
		it.Event = new(FileMetaContractFileMetaAdded)
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
func (it *FileMetaContractFileMetaAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileMetaContractFileMetaAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileMetaContractFileMetaAdded represents a FileMetaAdded event raised by the FileMetaContract contract.
type FileMetaContractFileMetaAdded struct {
	Cid      string
	MetaData []byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterFileMetaAdded is a free log retrieval operation binding the contract event 0x4aa65a9b4b57eefd0bf11e716e53b8896cb8ee8d0a4ec9bb9c231b091e0045a7.
//
// Solidity: event FileMetaAdded(string cid, bytes metaData)
func (_FileMetaContract *FileMetaContractFilterer) FilterFileMetaAdded(opts *bind.FilterOpts) (*FileMetaContractFileMetaAddedIterator, error) {

	logs, sub, err := _FileMetaContract.contract.FilterLogs(opts, "FileMetaAdded")
	if err != nil {
		return nil, err
	}
	return &FileMetaContractFileMetaAddedIterator{contract: _FileMetaContract.contract, event: "FileMetaAdded", logs: logs, sub: sub}, nil
}

// WatchFileMetaAdded is a free log subscription operation binding the contract event 0x4aa65a9b4b57eefd0bf11e716e53b8896cb8ee8d0a4ec9bb9c231b091e0045a7.
//
// Solidity: event FileMetaAdded(string cid, bytes metaData)
func (_FileMetaContract *FileMetaContractFilterer) WatchFileMetaAdded(opts *bind.WatchOpts, sink chan<- *FileMetaContractFileMetaAdded) (event.Subscription, error) {

	logs, sub, err := _FileMetaContract.contract.WatchLogs(opts, "FileMetaAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileMetaContractFileMetaAdded)
				if err := _FileMetaContract.contract.UnpackLog(event, "FileMetaAdded", log); err != nil {
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

// ParseFileMetaAdded is a log parse operation binding the contract event 0x4aa65a9b4b57eefd0bf11e716e53b8896cb8ee8d0a4ec9bb9c231b091e0045a7.
//
// Solidity: event FileMetaAdded(string cid, bytes metaData)
func (_FileMetaContract *FileMetaContractFilterer) ParseFileMetaAdded(log types.Log) (*FileMetaContractFileMetaAdded, error) {
	event := new(FileMetaContractFileMetaAdded)
	if err := _FileMetaContract.contract.UnpackLog(event, "FileMetaAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileMetaContractInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the FileMetaContract contract.
type FileMetaContractInitializedIterator struct {
	Event *FileMetaContractInitialized // Event containing the contract specifics and raw log

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
func (it *FileMetaContractInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileMetaContractInitialized)
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
		it.Event = new(FileMetaContractInitialized)
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
func (it *FileMetaContractInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileMetaContractInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileMetaContractInitialized represents a Initialized event raised by the FileMetaContract contract.
type FileMetaContractInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_FileMetaContract *FileMetaContractFilterer) FilterInitialized(opts *bind.FilterOpts) (*FileMetaContractInitializedIterator, error) {

	logs, sub, err := _FileMetaContract.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FileMetaContractInitializedIterator{contract: _FileMetaContract.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_FileMetaContract *FileMetaContractFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FileMetaContractInitialized) (event.Subscription, error) {

	logs, sub, err := _FileMetaContract.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileMetaContractInitialized)
				if err := _FileMetaContract.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_FileMetaContract *FileMetaContractFilterer) ParseInitialized(log types.Log) (*FileMetaContractInitialized, error) {
	event := new(FileMetaContractInitialized)
	if err := _FileMetaContract.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileMetaContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FileMetaContract contract.
type FileMetaContractOwnershipTransferredIterator struct {
	Event *FileMetaContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *FileMetaContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileMetaContractOwnershipTransferred)
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
		it.Event = new(FileMetaContractOwnershipTransferred)
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
func (it *FileMetaContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileMetaContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileMetaContractOwnershipTransferred represents a OwnershipTransferred event raised by the FileMetaContract contract.
type FileMetaContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FileMetaContract *FileMetaContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FileMetaContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FileMetaContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FileMetaContractOwnershipTransferredIterator{contract: _FileMetaContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FileMetaContract *FileMetaContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FileMetaContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FileMetaContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileMetaContractOwnershipTransferred)
				if err := _FileMetaContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_FileMetaContract *FileMetaContractFilterer) ParseOwnershipTransferred(log types.Log) (*FileMetaContractOwnershipTransferred, error) {
	event := new(FileMetaContractOwnershipTransferred)
	if err := _FileMetaContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileMetaContractStatusUpdatedIterator is returned from FilterStatusUpdated and is used to iterate over the raw logs and unpacked data for StatusUpdated events raised by the FileMetaContract contract.
type FileMetaContractStatusUpdatedIterator struct {
	Event *FileMetaContractStatusUpdated // Event containing the contract specifics and raw log

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
func (it *FileMetaContractStatusUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileMetaContractStatusUpdated)
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
		it.Event = new(FileMetaContractStatusUpdated)
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
func (it *FileMetaContractStatusUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileMetaContractStatusUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileMetaContractStatusUpdated represents a StatusUpdated event raised by the FileMetaContract contract.
type FileMetaContractStatusUpdated struct {
	ContractId common.Hash
	Status     uint8
	Timestamp  *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterStatusUpdated is a free log retrieval operation binding the contract event 0xc1c0ab06b259f062bfd174cf912a7c24b40820e430b30e09c9e96e17737bc9d3.
//
// Solidity: event StatusUpdated(string indexed contractId, uint8 status, uint256 timestamp)
func (_FileMetaContract *FileMetaContractFilterer) FilterStatusUpdated(opts *bind.FilterOpts, contractId []string) (*FileMetaContractStatusUpdatedIterator, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _FileMetaContract.contract.FilterLogs(opts, "StatusUpdated", contractIdRule)
	if err != nil {
		return nil, err
	}
	return &FileMetaContractStatusUpdatedIterator{contract: _FileMetaContract.contract, event: "StatusUpdated", logs: logs, sub: sub}, nil
}

// WatchStatusUpdated is a free log subscription operation binding the contract event 0xc1c0ab06b259f062bfd174cf912a7c24b40820e430b30e09c9e96e17737bc9d3.
//
// Solidity: event StatusUpdated(string indexed contractId, uint8 status, uint256 timestamp)
func (_FileMetaContract *FileMetaContractFilterer) WatchStatusUpdated(opts *bind.WatchOpts, sink chan<- *FileMetaContractStatusUpdated, contractId []string) (event.Subscription, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _FileMetaContract.contract.WatchLogs(opts, "StatusUpdated", contractIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileMetaContractStatusUpdated)
				if err := _FileMetaContract.contract.UnpackLog(event, "StatusUpdated", log); err != nil {
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

// ParseStatusUpdated is a log parse operation binding the contract event 0xc1c0ab06b259f062bfd174cf912a7c24b40820e430b30e09c9e96e17737bc9d3.
//
// Solidity: event StatusUpdated(string indexed contractId, uint8 status, uint256 timestamp)
func (_FileMetaContract *FileMetaContractFilterer) ParseStatusUpdated(log types.Log) (*FileMetaContractStatusUpdated, error) {
	event := new(FileMetaContractStatusUpdated)
	if err := _FileMetaContract.contract.UnpackLog(event, "StatusUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FileMetaContractUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the FileMetaContract contract.
type FileMetaContractUpgradedIterator struct {
	Event *FileMetaContractUpgraded // Event containing the contract specifics and raw log

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
func (it *FileMetaContractUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FileMetaContractUpgraded)
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
		it.Event = new(FileMetaContractUpgraded)
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
func (it *FileMetaContractUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FileMetaContractUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FileMetaContractUpgraded represents a Upgraded event raised by the FileMetaContract contract.
type FileMetaContractUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FileMetaContract *FileMetaContractFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*FileMetaContractUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FileMetaContract.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &FileMetaContractUpgradedIterator{contract: _FileMetaContract.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FileMetaContract *FileMetaContractFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *FileMetaContractUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FileMetaContract.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FileMetaContractUpgraded)
				if err := _FileMetaContract.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_FileMetaContract *FileMetaContractFilterer) ParseUpgraded(log types.Log) (*FileMetaContractUpgraded, error) {
	event := new(FileMetaContractUpgraded)
	if err := _FileMetaContract.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
