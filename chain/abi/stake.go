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

// StakeContractMetaData contains all meta data concerning the StakeContract contract.
var StakeContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newUnlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMin\",\"type\":\"uint256\"}],\"name\":\"ParametersUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"}],\"name\":\"Unstaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGlobalStats\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_totalStaked\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_totalUnlocked\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"contractBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"stakedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockedAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_unlockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minStakeAmount\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minStakeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"stakers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"stakedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockedAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalStaked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalUnlocked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlockPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newUnlockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_newMinStake\",\"type\":\"uint256\"}],\"name\":\"updateParameters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// StakeContractABI is the input ABI used to generate the binding from.
// Deprecated: Use StakeContractMetaData.ABI instead.
var StakeContractABI = StakeContractMetaData.ABI

// StakeContract is an auto generated Go binding around an Ethereum contract.
type StakeContract struct {
	StakeContractCaller     // Read-only binding to the contract
	StakeContractTransactor // Write-only binding to the contract
	StakeContractFilterer   // Log filterer for contract events
}

// StakeContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakeContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakeContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakeContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakeContractSession struct {
	Contract     *StakeContract    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakeContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakeContractCallerSession struct {
	Contract *StakeContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// StakeContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakeContractTransactorSession struct {
	Contract     *StakeContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// StakeContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakeContractRaw struct {
	Contract *StakeContract // Generic contract binding to access the raw methods on
}

// StakeContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakeContractCallerRaw struct {
	Contract *StakeContractCaller // Generic read-only contract binding to access the raw methods on
}

// StakeContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakeContractTransactorRaw struct {
	Contract *StakeContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakeContract creates a new instance of StakeContract, bound to a specific deployed contract.
func NewStakeContract(address common.Address, backend bind.ContractBackend) (*StakeContract, error) {
	contract, err := bindStakeContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakeContract{StakeContractCaller: StakeContractCaller{contract: contract}, StakeContractTransactor: StakeContractTransactor{contract: contract}, StakeContractFilterer: StakeContractFilterer{contract: contract}}, nil
}

// NewStakeContractCaller creates a new read-only instance of StakeContract, bound to a specific deployed contract.
func NewStakeContractCaller(address common.Address, caller bind.ContractCaller) (*StakeContractCaller, error) {
	contract, err := bindStakeContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakeContractCaller{contract: contract}, nil
}

// NewStakeContractTransactor creates a new write-only instance of StakeContract, bound to a specific deployed contract.
func NewStakeContractTransactor(address common.Address, transactor bind.ContractTransactor) (*StakeContractTransactor, error) {
	contract, err := bindStakeContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakeContractTransactor{contract: contract}, nil
}

// NewStakeContractFilterer creates a new log filterer instance of StakeContract, bound to a specific deployed contract.
func NewStakeContractFilterer(address common.Address, filterer bind.ContractFilterer) (*StakeContractFilterer, error) {
	contract, err := bindStakeContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakeContractFilterer{contract: contract}, nil
}

// bindStakeContract binds a generic wrapper to an already deployed contract.
func bindStakeContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakeContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakeContract *StakeContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakeContract.Contract.StakeContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakeContract *StakeContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeContract.Contract.StakeContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakeContract *StakeContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakeContract.Contract.StakeContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakeContract *StakeContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakeContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakeContract *StakeContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakeContract *StakeContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakeContract.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_StakeContract *StakeContractCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StakeContract.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_StakeContract *StakeContractSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _StakeContract.Contract.UPGRADEINTERFACEVERSION(&_StakeContract.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_StakeContract *StakeContractCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _StakeContract.Contract.UPGRADEINTERFACEVERSION(&_StakeContract.CallOpts)
}

// GetGlobalStats is a free data retrieval call binding the contract method 0x6b4169c3.
//
// Solidity: function getGlobalStats() view returns(uint256 _totalStaked, uint256 _totalUnlocked, uint256 contractBalance)
func (_StakeContract *StakeContractCaller) GetGlobalStats(opts *bind.CallOpts) (struct {
	TotalStaked     *big.Int
	TotalUnlocked   *big.Int
	ContractBalance *big.Int
}, error) {
	var out []interface{}
	err := _StakeContract.contract.Call(opts, &out, "getGlobalStats")

	outstruct := new(struct {
		TotalStaked     *big.Int
		TotalUnlocked   *big.Int
		ContractBalance *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalStaked = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalUnlocked = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ContractBalance = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetGlobalStats is a free data retrieval call binding the contract method 0x6b4169c3.
//
// Solidity: function getGlobalStats() view returns(uint256 _totalStaked, uint256 _totalUnlocked, uint256 contractBalance)
func (_StakeContract *StakeContractSession) GetGlobalStats() (struct {
	TotalStaked     *big.Int
	TotalUnlocked   *big.Int
	ContractBalance *big.Int
}, error) {
	return _StakeContract.Contract.GetGlobalStats(&_StakeContract.CallOpts)
}

// GetGlobalStats is a free data retrieval call binding the contract method 0x6b4169c3.
//
// Solidity: function getGlobalStats() view returns(uint256 _totalStaked, uint256 _totalUnlocked, uint256 contractBalance)
func (_StakeContract *StakeContractCallerSession) GetGlobalStats() (struct {
	TotalStaked     *big.Int
	TotalUnlocked   *big.Int
	ContractBalance *big.Int
}, error) {
	return _StakeContract.Contract.GetGlobalStats(&_StakeContract.CallOpts)
}

// GetUserStake is a free data retrieval call binding the contract method 0xbbadc93a.
//
// Solidity: function getUserStake(address user) view returns(uint256 stakedAmount, uint256 unlockTime, uint256 unlockedAmount)
func (_StakeContract *StakeContractCaller) GetUserStake(opts *bind.CallOpts, user common.Address) (struct {
	StakedAmount   *big.Int
	UnlockTime     *big.Int
	UnlockedAmount *big.Int
}, error) {
	var out []interface{}
	err := _StakeContract.contract.Call(opts, &out, "getUserStake", user)

	outstruct := new(struct {
		StakedAmount   *big.Int
		UnlockTime     *big.Int
		UnlockedAmount *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StakedAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.UnlockTime = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.UnlockedAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetUserStake is a free data retrieval call binding the contract method 0xbbadc93a.
//
// Solidity: function getUserStake(address user) view returns(uint256 stakedAmount, uint256 unlockTime, uint256 unlockedAmount)
func (_StakeContract *StakeContractSession) GetUserStake(user common.Address) (struct {
	StakedAmount   *big.Int
	UnlockTime     *big.Int
	UnlockedAmount *big.Int
}, error) {
	return _StakeContract.Contract.GetUserStake(&_StakeContract.CallOpts, user)
}

// GetUserStake is a free data retrieval call binding the contract method 0xbbadc93a.
//
// Solidity: function getUserStake(address user) view returns(uint256 stakedAmount, uint256 unlockTime, uint256 unlockedAmount)
func (_StakeContract *StakeContractCallerSession) GetUserStake(user common.Address) (struct {
	StakedAmount   *big.Int
	UnlockTime     *big.Int
	UnlockedAmount *big.Int
}, error) {
	return _StakeContract.Contract.GetUserStake(&_StakeContract.CallOpts, user)
}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_StakeContract *StakeContractCaller) MinStakeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeContract.contract.Call(opts, &out, "minStakeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_StakeContract *StakeContractSession) MinStakeAmount() (*big.Int, error) {
	return _StakeContract.Contract.MinStakeAmount(&_StakeContract.CallOpts)
}

// MinStakeAmount is a free data retrieval call binding the contract method 0xf1887684.
//
// Solidity: function minStakeAmount() view returns(uint256)
func (_StakeContract *StakeContractCallerSession) MinStakeAmount() (*big.Int, error) {
	return _StakeContract.Contract.MinStakeAmount(&_StakeContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakeContract *StakeContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakeContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakeContract *StakeContractSession) Owner() (common.Address, error) {
	return _StakeContract.Contract.Owner(&_StakeContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakeContract *StakeContractCallerSession) Owner() (common.Address, error) {
	return _StakeContract.Contract.Owner(&_StakeContract.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_StakeContract *StakeContractCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _StakeContract.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_StakeContract *StakeContractSession) ProxiableUUID() ([32]byte, error) {
	return _StakeContract.Contract.ProxiableUUID(&_StakeContract.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_StakeContract *StakeContractCallerSession) ProxiableUUID() ([32]byte, error) {
	return _StakeContract.Contract.ProxiableUUID(&_StakeContract.CallOpts)
}

// Stakers is a free data retrieval call binding the contract method 0x9168ae72.
//
// Solidity: function stakers(address ) view returns(uint256 stakedAmount, uint256 unlockTime, uint256 unlockedAmount)
func (_StakeContract *StakeContractCaller) Stakers(opts *bind.CallOpts, arg0 common.Address) (struct {
	StakedAmount   *big.Int
	UnlockTime     *big.Int
	UnlockedAmount *big.Int
}, error) {
	var out []interface{}
	err := _StakeContract.contract.Call(opts, &out, "stakers", arg0)

	outstruct := new(struct {
		StakedAmount   *big.Int
		UnlockTime     *big.Int
		UnlockedAmount *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StakedAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.UnlockTime = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.UnlockedAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Stakers is a free data retrieval call binding the contract method 0x9168ae72.
//
// Solidity: function stakers(address ) view returns(uint256 stakedAmount, uint256 unlockTime, uint256 unlockedAmount)
func (_StakeContract *StakeContractSession) Stakers(arg0 common.Address) (struct {
	StakedAmount   *big.Int
	UnlockTime     *big.Int
	UnlockedAmount *big.Int
}, error) {
	return _StakeContract.Contract.Stakers(&_StakeContract.CallOpts, arg0)
}

// Stakers is a free data retrieval call binding the contract method 0x9168ae72.
//
// Solidity: function stakers(address ) view returns(uint256 stakedAmount, uint256 unlockTime, uint256 unlockedAmount)
func (_StakeContract *StakeContractCallerSession) Stakers(arg0 common.Address) (struct {
	StakedAmount   *big.Int
	UnlockTime     *big.Int
	UnlockedAmount *big.Int
}, error) {
	return _StakeContract.Contract.Stakers(&_StakeContract.CallOpts, arg0)
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() view returns(uint256)
func (_StakeContract *StakeContractCaller) TotalStaked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeContract.contract.Call(opts, &out, "totalStaked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() view returns(uint256)
func (_StakeContract *StakeContractSession) TotalStaked() (*big.Int, error) {
	return _StakeContract.Contract.TotalStaked(&_StakeContract.CallOpts)
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() view returns(uint256)
func (_StakeContract *StakeContractCallerSession) TotalStaked() (*big.Int, error) {
	return _StakeContract.Contract.TotalStaked(&_StakeContract.CallOpts)
}

// TotalUnlocked is a free data retrieval call binding the contract method 0xa779d080.
//
// Solidity: function totalUnlocked() view returns(uint256)
func (_StakeContract *StakeContractCaller) TotalUnlocked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeContract.contract.Call(opts, &out, "totalUnlocked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalUnlocked is a free data retrieval call binding the contract method 0xa779d080.
//
// Solidity: function totalUnlocked() view returns(uint256)
func (_StakeContract *StakeContractSession) TotalUnlocked() (*big.Int, error) {
	return _StakeContract.Contract.TotalUnlocked(&_StakeContract.CallOpts)
}

// TotalUnlocked is a free data retrieval call binding the contract method 0xa779d080.
//
// Solidity: function totalUnlocked() view returns(uint256)
func (_StakeContract *StakeContractCallerSession) TotalUnlocked() (*big.Int, error) {
	return _StakeContract.Contract.TotalUnlocked(&_StakeContract.CallOpts)
}

// UnlockPeriod is a free data retrieval call binding the contract method 0x20d3a0b4.
//
// Solidity: function unlockPeriod() view returns(uint256)
func (_StakeContract *StakeContractCaller) UnlockPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeContract.contract.Call(opts, &out, "unlockPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnlockPeriod is a free data retrieval call binding the contract method 0x20d3a0b4.
//
// Solidity: function unlockPeriod() view returns(uint256)
func (_StakeContract *StakeContractSession) UnlockPeriod() (*big.Int, error) {
	return _StakeContract.Contract.UnlockPeriod(&_StakeContract.CallOpts)
}

// UnlockPeriod is a free data retrieval call binding the contract method 0x20d3a0b4.
//
// Solidity: function unlockPeriod() view returns(uint256)
func (_StakeContract *StakeContractCallerSession) UnlockPeriod() (*big.Int, error) {
	return _StakeContract.Contract.UnlockPeriod(&_StakeContract.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xe4a30116.
//
// Solidity: function initialize(uint256 _unlockPeriod, uint256 _minStakeAmount) returns()
func (_StakeContract *StakeContractTransactor) Initialize(opts *bind.TransactOpts, _unlockPeriod *big.Int, _minStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakeContract.contract.Transact(opts, "initialize", _unlockPeriod, _minStakeAmount)
}

// Initialize is a paid mutator transaction binding the contract method 0xe4a30116.
//
// Solidity: function initialize(uint256 _unlockPeriod, uint256 _minStakeAmount) returns()
func (_StakeContract *StakeContractSession) Initialize(_unlockPeriod *big.Int, _minStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakeContract.Contract.Initialize(&_StakeContract.TransactOpts, _unlockPeriod, _minStakeAmount)
}

// Initialize is a paid mutator transaction binding the contract method 0xe4a30116.
//
// Solidity: function initialize(uint256 _unlockPeriod, uint256 _minStakeAmount) returns()
func (_StakeContract *StakeContractTransactorSession) Initialize(_unlockPeriod *big.Int, _minStakeAmount *big.Int) (*types.Transaction, error) {
	return _StakeContract.Contract.Initialize(&_StakeContract.TransactOpts, _unlockPeriod, _minStakeAmount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StakeContract *StakeContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StakeContract *StakeContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _StakeContract.Contract.RenounceOwnership(&_StakeContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StakeContract *StakeContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _StakeContract.Contract.RenounceOwnership(&_StakeContract.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_StakeContract *StakeContractTransactor) Stake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeContract.contract.Transact(opts, "stake")
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_StakeContract *StakeContractSession) Stake() (*types.Transaction, error) {
	return _StakeContract.Contract.Stake(&_StakeContract.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_StakeContract *StakeContractTransactorSession) Stake() (*types.Transaction, error) {
	return _StakeContract.Contract.Stake(&_StakeContract.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StakeContract *StakeContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _StakeContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StakeContract *StakeContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _StakeContract.Contract.TransferOwnership(&_StakeContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StakeContract *StakeContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _StakeContract.Contract.TransferOwnership(&_StakeContract.TransactOpts, newOwner)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 amount) returns()
func (_StakeContract *StakeContractTransactor) Unstake(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _StakeContract.contract.Transact(opts, "unstake", amount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 amount) returns()
func (_StakeContract *StakeContractSession) Unstake(amount *big.Int) (*types.Transaction, error) {
	return _StakeContract.Contract.Unstake(&_StakeContract.TransactOpts, amount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 amount) returns()
func (_StakeContract *StakeContractTransactorSession) Unstake(amount *big.Int) (*types.Transaction, error) {
	return _StakeContract.Contract.Unstake(&_StakeContract.TransactOpts, amount)
}

// UpdateParameters is a paid mutator transaction binding the contract method 0x16128211.
//
// Solidity: function updateParameters(uint256 _newUnlockPeriod, uint256 _newMinStake) returns()
func (_StakeContract *StakeContractTransactor) UpdateParameters(opts *bind.TransactOpts, _newUnlockPeriod *big.Int, _newMinStake *big.Int) (*types.Transaction, error) {
	return _StakeContract.contract.Transact(opts, "updateParameters", _newUnlockPeriod, _newMinStake)
}

// UpdateParameters is a paid mutator transaction binding the contract method 0x16128211.
//
// Solidity: function updateParameters(uint256 _newUnlockPeriod, uint256 _newMinStake) returns()
func (_StakeContract *StakeContractSession) UpdateParameters(_newUnlockPeriod *big.Int, _newMinStake *big.Int) (*types.Transaction, error) {
	return _StakeContract.Contract.UpdateParameters(&_StakeContract.TransactOpts, _newUnlockPeriod, _newMinStake)
}

// UpdateParameters is a paid mutator transaction binding the contract method 0x16128211.
//
// Solidity: function updateParameters(uint256 _newUnlockPeriod, uint256 _newMinStake) returns()
func (_StakeContract *StakeContractTransactorSession) UpdateParameters(_newUnlockPeriod *big.Int, _newMinStake *big.Int) (*types.Transaction, error) {
	return _StakeContract.Contract.UpdateParameters(&_StakeContract.TransactOpts, _newUnlockPeriod, _newMinStake)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_StakeContract *StakeContractTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _StakeContract.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_StakeContract *StakeContractSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _StakeContract.Contract.UpgradeToAndCall(&_StakeContract.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_StakeContract *StakeContractTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _StakeContract.Contract.UpgradeToAndCall(&_StakeContract.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_StakeContract *StakeContractTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeContract.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_StakeContract *StakeContractSession) Withdraw() (*types.Transaction, error) {
	return _StakeContract.Contract.Withdraw(&_StakeContract.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_StakeContract *StakeContractTransactorSession) Withdraw() (*types.Transaction, error) {
	return _StakeContract.Contract.Withdraw(&_StakeContract.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_StakeContract *StakeContractTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeContract.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_StakeContract *StakeContractSession) Receive() (*types.Transaction, error) {
	return _StakeContract.Contract.Receive(&_StakeContract.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_StakeContract *StakeContractTransactorSession) Receive() (*types.Transaction, error) {
	return _StakeContract.Contract.Receive(&_StakeContract.TransactOpts)
}

// StakeContractInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the StakeContract contract.
type StakeContractInitializedIterator struct {
	Event *StakeContractInitialized // Event containing the contract specifics and raw log

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
func (it *StakeContractInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeContractInitialized)
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
		it.Event = new(StakeContractInitialized)
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
func (it *StakeContractInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeContractInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeContractInitialized represents a Initialized event raised by the StakeContract contract.
type StakeContractInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_StakeContract *StakeContractFilterer) FilterInitialized(opts *bind.FilterOpts) (*StakeContractInitializedIterator, error) {

	logs, sub, err := _StakeContract.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &StakeContractInitializedIterator{contract: _StakeContract.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_StakeContract *StakeContractFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *StakeContractInitialized) (event.Subscription, error) {

	logs, sub, err := _StakeContract.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeContractInitialized)
				if err := _StakeContract.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_StakeContract *StakeContractFilterer) ParseInitialized(log types.Log) (*StakeContractInitialized, error) {
	event := new(StakeContractInitialized)
	if err := _StakeContract.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the StakeContract contract.
type StakeContractOwnershipTransferredIterator struct {
	Event *StakeContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *StakeContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeContractOwnershipTransferred)
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
		it.Event = new(StakeContractOwnershipTransferred)
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
func (it *StakeContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeContractOwnershipTransferred represents a OwnershipTransferred event raised by the StakeContract contract.
type StakeContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_StakeContract *StakeContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StakeContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _StakeContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StakeContractOwnershipTransferredIterator{contract: _StakeContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_StakeContract *StakeContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StakeContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _StakeContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeContractOwnershipTransferred)
				if err := _StakeContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_StakeContract *StakeContractFilterer) ParseOwnershipTransferred(log types.Log) (*StakeContractOwnershipTransferred, error) {
	event := new(StakeContractOwnershipTransferred)
	if err := _StakeContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeContractParametersUpdatedIterator is returned from FilterParametersUpdated and is used to iterate over the raw logs and unpacked data for ParametersUpdated events raised by the StakeContract contract.
type StakeContractParametersUpdatedIterator struct {
	Event *StakeContractParametersUpdated // Event containing the contract specifics and raw log

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
func (it *StakeContractParametersUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeContractParametersUpdated)
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
		it.Event = new(StakeContractParametersUpdated)
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
func (it *StakeContractParametersUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeContractParametersUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeContractParametersUpdated represents a ParametersUpdated event raised by the StakeContract contract.
type StakeContractParametersUpdated struct {
	NewUnlock *big.Int
	NewMin    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterParametersUpdated is a free log retrieval operation binding the contract event 0xfaccb0639ff7851e0e24f3b2d9ab03cd62ffb63f5b4d90aaeff85bb078c1fa48.
//
// Solidity: event ParametersUpdated(uint256 newUnlock, uint256 newMin)
func (_StakeContract *StakeContractFilterer) FilterParametersUpdated(opts *bind.FilterOpts) (*StakeContractParametersUpdatedIterator, error) {

	logs, sub, err := _StakeContract.contract.FilterLogs(opts, "ParametersUpdated")
	if err != nil {
		return nil, err
	}
	return &StakeContractParametersUpdatedIterator{contract: _StakeContract.contract, event: "ParametersUpdated", logs: logs, sub: sub}, nil
}

// WatchParametersUpdated is a free log subscription operation binding the contract event 0xfaccb0639ff7851e0e24f3b2d9ab03cd62ffb63f5b4d90aaeff85bb078c1fa48.
//
// Solidity: event ParametersUpdated(uint256 newUnlock, uint256 newMin)
func (_StakeContract *StakeContractFilterer) WatchParametersUpdated(opts *bind.WatchOpts, sink chan<- *StakeContractParametersUpdated) (event.Subscription, error) {

	logs, sub, err := _StakeContract.contract.WatchLogs(opts, "ParametersUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeContractParametersUpdated)
				if err := _StakeContract.contract.UnpackLog(event, "ParametersUpdated", log); err != nil {
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

// ParseParametersUpdated is a log parse operation binding the contract event 0xfaccb0639ff7851e0e24f3b2d9ab03cd62ffb63f5b4d90aaeff85bb078c1fa48.
//
// Solidity: event ParametersUpdated(uint256 newUnlock, uint256 newMin)
func (_StakeContract *StakeContractFilterer) ParseParametersUpdated(log types.Log) (*StakeContractParametersUpdated, error) {
	event := new(StakeContractParametersUpdated)
	if err := _StakeContract.contract.UnpackLog(event, "ParametersUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeContractStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the StakeContract contract.
type StakeContractStakedIterator struct {
	Event *StakeContractStaked // Event containing the contract specifics and raw log

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
func (it *StakeContractStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeContractStaked)
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
		it.Event = new(StakeContractStaked)
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
func (it *StakeContractStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeContractStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeContractStaked represents a Staked event raised by the StakeContract contract.
type StakeContractStaked struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0x9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d.
//
// Solidity: event Staked(address indexed user, uint256 amount)
func (_StakeContract *StakeContractFilterer) FilterStaked(opts *bind.FilterOpts, user []common.Address) (*StakeContractStakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _StakeContract.contract.FilterLogs(opts, "Staked", userRule)
	if err != nil {
		return nil, err
	}
	return &StakeContractStakedIterator{contract: _StakeContract.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0x9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d.
//
// Solidity: event Staked(address indexed user, uint256 amount)
func (_StakeContract *StakeContractFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *StakeContractStaked, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _StakeContract.contract.WatchLogs(opts, "Staked", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeContractStaked)
				if err := _StakeContract.contract.UnpackLog(event, "Staked", log); err != nil {
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

// ParseStaked is a log parse operation binding the contract event 0x9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d.
//
// Solidity: event Staked(address indexed user, uint256 amount)
func (_StakeContract *StakeContractFilterer) ParseStaked(log types.Log) (*StakeContractStaked, error) {
	event := new(StakeContractStaked)
	if err := _StakeContract.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeContractUnstakedIterator is returned from FilterUnstaked and is used to iterate over the raw logs and unpacked data for Unstaked events raised by the StakeContract contract.
type StakeContractUnstakedIterator struct {
	Event *StakeContractUnstaked // Event containing the contract specifics and raw log

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
func (it *StakeContractUnstakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeContractUnstaked)
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
		it.Event = new(StakeContractUnstaked)
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
func (it *StakeContractUnstakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeContractUnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeContractUnstaked represents a Unstaked event raised by the StakeContract contract.
type StakeContractUnstaked struct {
	User       common.Address
	Amount     *big.Int
	UnlockTime *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUnstaked is a free log retrieval operation binding the contract event 0x7fc4727e062e336010f2c282598ef5f14facb3de68cf8195c2f23e1454b2b74e.
//
// Solidity: event Unstaked(address indexed user, uint256 amount, uint256 unlockTime)
func (_StakeContract *StakeContractFilterer) FilterUnstaked(opts *bind.FilterOpts, user []common.Address) (*StakeContractUnstakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _StakeContract.contract.FilterLogs(opts, "Unstaked", userRule)
	if err != nil {
		return nil, err
	}
	return &StakeContractUnstakedIterator{contract: _StakeContract.contract, event: "Unstaked", logs: logs, sub: sub}, nil
}

// WatchUnstaked is a free log subscription operation binding the contract event 0x7fc4727e062e336010f2c282598ef5f14facb3de68cf8195c2f23e1454b2b74e.
//
// Solidity: event Unstaked(address indexed user, uint256 amount, uint256 unlockTime)
func (_StakeContract *StakeContractFilterer) WatchUnstaked(opts *bind.WatchOpts, sink chan<- *StakeContractUnstaked, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _StakeContract.contract.WatchLogs(opts, "Unstaked", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeContractUnstaked)
				if err := _StakeContract.contract.UnpackLog(event, "Unstaked", log); err != nil {
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

// ParseUnstaked is a log parse operation binding the contract event 0x7fc4727e062e336010f2c282598ef5f14facb3de68cf8195c2f23e1454b2b74e.
//
// Solidity: event Unstaked(address indexed user, uint256 amount, uint256 unlockTime)
func (_StakeContract *StakeContractFilterer) ParseUnstaked(log types.Log) (*StakeContractUnstaked, error) {
	event := new(StakeContractUnstaked)
	if err := _StakeContract.contract.UnpackLog(event, "Unstaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeContractUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the StakeContract contract.
type StakeContractUpgradedIterator struct {
	Event *StakeContractUpgraded // Event containing the contract specifics and raw log

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
func (it *StakeContractUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeContractUpgraded)
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
		it.Event = new(StakeContractUpgraded)
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
func (it *StakeContractUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeContractUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeContractUpgraded represents a Upgraded event raised by the StakeContract contract.
type StakeContractUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_StakeContract *StakeContractFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*StakeContractUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _StakeContract.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &StakeContractUpgradedIterator{contract: _StakeContract.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_StakeContract *StakeContractFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *StakeContractUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _StakeContract.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeContractUpgraded)
				if err := _StakeContract.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_StakeContract *StakeContractFilterer) ParseUpgraded(log types.Log) (*StakeContractUpgraded, error) {
	event := new(StakeContractUpgraded)
	if err := _StakeContract.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeContractWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the StakeContract contract.
type StakeContractWithdrawnIterator struct {
	Event *StakeContractWithdrawn // Event containing the contract specifics and raw log

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
func (it *StakeContractWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeContractWithdrawn)
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
		it.Event = new(StakeContractWithdrawn)
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
func (it *StakeContractWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeContractWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeContractWithdrawn represents a Withdrawn event raised by the StakeContract contract.
type StakeContractWithdrawn struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed user, uint256 amount)
func (_StakeContract *StakeContractFilterer) FilterWithdrawn(opts *bind.FilterOpts, user []common.Address) (*StakeContractWithdrawnIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _StakeContract.contract.FilterLogs(opts, "Withdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return &StakeContractWithdrawnIterator{contract: _StakeContract.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed user, uint256 amount)
func (_StakeContract *StakeContractFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *StakeContractWithdrawn, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _StakeContract.contract.WatchLogs(opts, "Withdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeContractWithdrawn)
				if err := _StakeContract.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed user, uint256 amount)
func (_StakeContract *StakeContractFilterer) ParseWithdrawn(log types.Log) (*StakeContractWithdrawn, error) {
	event := new(StakeContractWithdrawn)
	if err := _StakeContract.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
