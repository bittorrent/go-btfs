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

// ProposalManagerProposal is an auto generated low-level Go binding around an user-defined struct.
type ProposalManagerProposal struct {
	Id           *big.Int
	Proposer     common.Address
	ProposalType uint8
	Title        string
	Description  string
	Uri          string
	Status       uint8
	CreatedAt    *big.Int
	UpdatedAt    *big.Int
}

// ProposalContractMetaData contains all meta data concerning the ProposalContract contract.
var ProposalContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reviewer\",\"type\":\"address\"}],\"name\":\"ProposalApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"}],\"name\":\"ProposalCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reviewer\",\"type\":\"address\"}],\"name\":\"ProposalRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumProposalManager.ProposalStatus\",\"name\":\"oldStatus\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumProposalManager.ProposalStatus\",\"name\":\"newStatus\",\"type\":\"uint8\"}],\"name\":\"ProposalStatusChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumProposalManager.ProposalType\",\"name\":\"proposalType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"}],\"name\":\"ProposalSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"REVIEWER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"approveProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"cancelProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveProposalCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllProposals\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"enumProposalManager.ProposalType\",\"name\":\"proposalType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"internalType\":\"enumProposalManager.ProposalStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"internalType\":\"structProposalManager.Proposal[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getProposal\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"enumProposalManager.ProposalType\",\"name\":\"proposalType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"internalType\":\"enumProposalManager.ProposalStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"internalType\":\"structProposalManager.Proposal\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProposalCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumProposalManager.ProposalType\",\"name\":\"proposalType\",\"type\":\"uint8\"}],\"name\":\"getProposalsByType\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"enumProposalManager.ProposalType\",\"name\":\"proposalType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"internalType\":\"enumProposalManager.ProposalStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"internalType\":\"structProposalManager.Proposal[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"rejectProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumProposalManager.ProposalType\",\"name\":\"proposalType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"title\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"}],\"name\":\"submitProposal\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// ProposalContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ProposalContractMetaData.ABI instead.
var ProposalContractABI = ProposalContractMetaData.ABI

// ProposalContract is an auto generated Go binding around an Ethereum contract.
type ProposalContract struct {
	ProposalContractCaller     // Read-only binding to the contract
	ProposalContractTransactor // Write-only binding to the contract
	ProposalContractFilterer   // Log filterer for contract events
}

// ProposalContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProposalContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProposalContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProposalContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProposalContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProposalContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProposalContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProposalContractSession struct {
	Contract     *ProposalContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProposalContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProposalContractCallerSession struct {
	Contract *ProposalContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ProposalContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProposalContractTransactorSession struct {
	Contract     *ProposalContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ProposalContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProposalContractRaw struct {
	Contract *ProposalContract // Generic contract binding to access the raw methods on
}

// ProposalContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProposalContractCallerRaw struct {
	Contract *ProposalContractCaller // Generic read-only contract binding to access the raw methods on
}

// ProposalContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProposalContractTransactorRaw struct {
	Contract *ProposalContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProposalContract creates a new instance of ProposalContract, bound to a specific deployed contract.
func NewProposalContract(address common.Address, backend bind.ContractBackend) (*ProposalContract, error) {
	contract, err := bindProposalContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ProposalContract{ProposalContractCaller: ProposalContractCaller{contract: contract}, ProposalContractTransactor: ProposalContractTransactor{contract: contract}, ProposalContractFilterer: ProposalContractFilterer{contract: contract}}, nil
}

// NewProposalContractCaller creates a new read-only instance of ProposalContract, bound to a specific deployed contract.
func NewProposalContractCaller(address common.Address, caller bind.ContractCaller) (*ProposalContractCaller, error) {
	contract, err := bindProposalContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProposalContractCaller{contract: contract}, nil
}

// NewProposalContractTransactor creates a new write-only instance of ProposalContract, bound to a specific deployed contract.
func NewProposalContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ProposalContractTransactor, error) {
	contract, err := bindProposalContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProposalContractTransactor{contract: contract}, nil
}

// NewProposalContractFilterer creates a new log filterer instance of ProposalContract, bound to a specific deployed contract.
func NewProposalContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ProposalContractFilterer, error) {
	contract, err := bindProposalContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProposalContractFilterer{contract: contract}, nil
}

// bindProposalContract binds a generic wrapper to an already deployed contract.
func bindProposalContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ProposalContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProposalContract *ProposalContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProposalContract.Contract.ProposalContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProposalContract *ProposalContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProposalContract.Contract.ProposalContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProposalContract *ProposalContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProposalContract.Contract.ProposalContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProposalContract *ProposalContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProposalContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProposalContract *ProposalContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProposalContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProposalContract *ProposalContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProposalContract.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ProposalContract *ProposalContractCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ProposalContract *ProposalContractSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ProposalContract.Contract.DEFAULTADMINROLE(&_ProposalContract.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ProposalContract *ProposalContractCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ProposalContract.Contract.DEFAULTADMINROLE(&_ProposalContract.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_ProposalContract *ProposalContractCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_ProposalContract *ProposalContractSession) PAUSERROLE() ([32]byte, error) {
	return _ProposalContract.Contract.PAUSERROLE(&_ProposalContract.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_ProposalContract *ProposalContractCallerSession) PAUSERROLE() ([32]byte, error) {
	return _ProposalContract.Contract.PAUSERROLE(&_ProposalContract.CallOpts)
}

// REVIEWERROLE is a free data retrieval call binding the contract method 0x3b129a56.
//
// Solidity: function REVIEWER_ROLE() view returns(bytes32)
func (_ProposalContract *ProposalContractCaller) REVIEWERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "REVIEWER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// REVIEWERROLE is a free data retrieval call binding the contract method 0x3b129a56.
//
// Solidity: function REVIEWER_ROLE() view returns(bytes32)
func (_ProposalContract *ProposalContractSession) REVIEWERROLE() ([32]byte, error) {
	return _ProposalContract.Contract.REVIEWERROLE(&_ProposalContract.CallOpts)
}

// REVIEWERROLE is a free data retrieval call binding the contract method 0x3b129a56.
//
// Solidity: function REVIEWER_ROLE() view returns(bytes32)
func (_ProposalContract *ProposalContractCallerSession) REVIEWERROLE() ([32]byte, error) {
	return _ProposalContract.Contract.REVIEWERROLE(&_ProposalContract.CallOpts)
}

// UPGRADERROLE is a free data retrieval call binding the contract method 0xf72c0d8b.
//
// Solidity: function UPGRADER_ROLE() view returns(bytes32)
func (_ProposalContract *ProposalContractCaller) UPGRADERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "UPGRADER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// UPGRADERROLE is a free data retrieval call binding the contract method 0xf72c0d8b.
//
// Solidity: function UPGRADER_ROLE() view returns(bytes32)
func (_ProposalContract *ProposalContractSession) UPGRADERROLE() ([32]byte, error) {
	return _ProposalContract.Contract.UPGRADERROLE(&_ProposalContract.CallOpts)
}

// UPGRADERROLE is a free data retrieval call binding the contract method 0xf72c0d8b.
//
// Solidity: function UPGRADER_ROLE() view returns(bytes32)
func (_ProposalContract *ProposalContractCallerSession) UPGRADERROLE() ([32]byte, error) {
	return _ProposalContract.Contract.UPGRADERROLE(&_ProposalContract.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ProposalContract *ProposalContractCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ProposalContract *ProposalContractSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ProposalContract.Contract.UPGRADEINTERFACEVERSION(&_ProposalContract.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ProposalContract *ProposalContractCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ProposalContract.Contract.UPGRADEINTERFACEVERSION(&_ProposalContract.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_ProposalContract *ProposalContractCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_ProposalContract *ProposalContractSession) VERSION() (string, error) {
	return _ProposalContract.Contract.VERSION(&_ProposalContract.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_ProposalContract *ProposalContractCallerSession) VERSION() (string, error) {
	return _ProposalContract.Contract.VERSION(&_ProposalContract.CallOpts)
}

// GetActiveProposalCount is a free data retrieval call binding the contract method 0x5e6a61ce.
//
// Solidity: function getActiveProposalCount() view returns(uint256)
func (_ProposalContract *ProposalContractCaller) GetActiveProposalCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "getActiveProposalCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActiveProposalCount is a free data retrieval call binding the contract method 0x5e6a61ce.
//
// Solidity: function getActiveProposalCount() view returns(uint256)
func (_ProposalContract *ProposalContractSession) GetActiveProposalCount() (*big.Int, error) {
	return _ProposalContract.Contract.GetActiveProposalCount(&_ProposalContract.CallOpts)
}

// GetActiveProposalCount is a free data retrieval call binding the contract method 0x5e6a61ce.
//
// Solidity: function getActiveProposalCount() view returns(uint256)
func (_ProposalContract *ProposalContractCallerSession) GetActiveProposalCount() (*big.Int, error) {
	return _ProposalContract.Contract.GetActiveProposalCount(&_ProposalContract.CallOpts)
}

// GetAllProposals is a free data retrieval call binding the contract method 0xcceb68f5.
//
// Solidity: function getAllProposals() view returns((uint256,address,uint8,string,string,string,uint8,uint256,uint256)[])
func (_ProposalContract *ProposalContractCaller) GetAllProposals(opts *bind.CallOpts) ([]ProposalManagerProposal, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "getAllProposals")

	if err != nil {
		return *new([]ProposalManagerProposal), err
	}

	out0 := *abi.ConvertType(out[0], new([]ProposalManagerProposal)).(*[]ProposalManagerProposal)

	return out0, err

}

// GetAllProposals is a free data retrieval call binding the contract method 0xcceb68f5.
//
// Solidity: function getAllProposals() view returns((uint256,address,uint8,string,string,string,uint8,uint256,uint256)[])
func (_ProposalContract *ProposalContractSession) GetAllProposals() ([]ProposalManagerProposal, error) {
	return _ProposalContract.Contract.GetAllProposals(&_ProposalContract.CallOpts)
}

// GetAllProposals is a free data retrieval call binding the contract method 0xcceb68f5.
//
// Solidity: function getAllProposals() view returns((uint256,address,uint8,string,string,string,uint8,uint256,uint256)[])
func (_ProposalContract *ProposalContractCallerSession) GetAllProposals() ([]ProposalManagerProposal, error) {
	return _ProposalContract.Contract.GetAllProposals(&_ProposalContract.CallOpts)
}

// GetProposal is a free data retrieval call binding the contract method 0xc7f758a8.
//
// Solidity: function getProposal(uint256 id) view returns((uint256,address,uint8,string,string,string,uint8,uint256,uint256))
func (_ProposalContract *ProposalContractCaller) GetProposal(opts *bind.CallOpts, id *big.Int) (ProposalManagerProposal, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "getProposal", id)

	if err != nil {
		return *new(ProposalManagerProposal), err
	}

	out0 := *abi.ConvertType(out[0], new(ProposalManagerProposal)).(*ProposalManagerProposal)

	return out0, err

}

// GetProposal is a free data retrieval call binding the contract method 0xc7f758a8.
//
// Solidity: function getProposal(uint256 id) view returns((uint256,address,uint8,string,string,string,uint8,uint256,uint256))
func (_ProposalContract *ProposalContractSession) GetProposal(id *big.Int) (ProposalManagerProposal, error) {
	return _ProposalContract.Contract.GetProposal(&_ProposalContract.CallOpts, id)
}

// GetProposal is a free data retrieval call binding the contract method 0xc7f758a8.
//
// Solidity: function getProposal(uint256 id) view returns((uint256,address,uint8,string,string,string,uint8,uint256,uint256))
func (_ProposalContract *ProposalContractCallerSession) GetProposal(id *big.Int) (ProposalManagerProposal, error) {
	return _ProposalContract.Contract.GetProposal(&_ProposalContract.CallOpts, id)
}

// GetProposalCount is a free data retrieval call binding the contract method 0xc08cc02d.
//
// Solidity: function getProposalCount() view returns(uint256)
func (_ProposalContract *ProposalContractCaller) GetProposalCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "getProposalCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetProposalCount is a free data retrieval call binding the contract method 0xc08cc02d.
//
// Solidity: function getProposalCount() view returns(uint256)
func (_ProposalContract *ProposalContractSession) GetProposalCount() (*big.Int, error) {
	return _ProposalContract.Contract.GetProposalCount(&_ProposalContract.CallOpts)
}

// GetProposalCount is a free data retrieval call binding the contract method 0xc08cc02d.
//
// Solidity: function getProposalCount() view returns(uint256)
func (_ProposalContract *ProposalContractCallerSession) GetProposalCount() (*big.Int, error) {
	return _ProposalContract.Contract.GetProposalCount(&_ProposalContract.CallOpts)
}

// GetProposalsByType is a free data retrieval call binding the contract method 0x53b50594.
//
// Solidity: function getProposalsByType(uint8 proposalType) view returns((uint256,address,uint8,string,string,string,uint8,uint256,uint256)[])
func (_ProposalContract *ProposalContractCaller) GetProposalsByType(opts *bind.CallOpts, proposalType uint8) ([]ProposalManagerProposal, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "getProposalsByType", proposalType)

	if err != nil {
		return *new([]ProposalManagerProposal), err
	}

	out0 := *abi.ConvertType(out[0], new([]ProposalManagerProposal)).(*[]ProposalManagerProposal)

	return out0, err

}

// GetProposalsByType is a free data retrieval call binding the contract method 0x53b50594.
//
// Solidity: function getProposalsByType(uint8 proposalType) view returns((uint256,address,uint8,string,string,string,uint8,uint256,uint256)[])
func (_ProposalContract *ProposalContractSession) GetProposalsByType(proposalType uint8) ([]ProposalManagerProposal, error) {
	return _ProposalContract.Contract.GetProposalsByType(&_ProposalContract.CallOpts, proposalType)
}

// GetProposalsByType is a free data retrieval call binding the contract method 0x53b50594.
//
// Solidity: function getProposalsByType(uint8 proposalType) view returns((uint256,address,uint8,string,string,string,uint8,uint256,uint256)[])
func (_ProposalContract *ProposalContractCallerSession) GetProposalsByType(proposalType uint8) ([]ProposalManagerProposal, error) {
	return _ProposalContract.Contract.GetProposalsByType(&_ProposalContract.CallOpts, proposalType)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ProposalContract *ProposalContractCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ProposalContract *ProposalContractSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ProposalContract.Contract.GetRoleAdmin(&_ProposalContract.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ProposalContract *ProposalContractCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ProposalContract.Contract.GetRoleAdmin(&_ProposalContract.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ProposalContract *ProposalContractCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ProposalContract *ProposalContractSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ProposalContract.Contract.HasRole(&_ProposalContract.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ProposalContract *ProposalContractCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ProposalContract.Contract.HasRole(&_ProposalContract.CallOpts, role, account)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ProposalContract *ProposalContractCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ProposalContract *ProposalContractSession) Paused() (bool, error) {
	return _ProposalContract.Contract.Paused(&_ProposalContract.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ProposalContract *ProposalContractCallerSession) Paused() (bool, error) {
	return _ProposalContract.Contract.Paused(&_ProposalContract.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ProposalContract *ProposalContractCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ProposalContract *ProposalContractSession) ProxiableUUID() ([32]byte, error) {
	return _ProposalContract.Contract.ProxiableUUID(&_ProposalContract.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ProposalContract *ProposalContractCallerSession) ProxiableUUID() ([32]byte, error) {
	return _ProposalContract.Contract.ProxiableUUID(&_ProposalContract.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ProposalContract *ProposalContractCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ProposalContract.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ProposalContract *ProposalContractSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ProposalContract.Contract.SupportsInterface(&_ProposalContract.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ProposalContract *ProposalContractCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ProposalContract.Contract.SupportsInterface(&_ProposalContract.CallOpts, interfaceId)
}

// ApproveProposal is a paid mutator transaction binding the contract method 0x98951b56.
//
// Solidity: function approveProposal(uint256 id) returns()
func (_ProposalContract *ProposalContractTransactor) ApproveProposal(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _ProposalContract.contract.Transact(opts, "approveProposal", id)
}

// ApproveProposal is a paid mutator transaction binding the contract method 0x98951b56.
//
// Solidity: function approveProposal(uint256 id) returns()
func (_ProposalContract *ProposalContractSession) ApproveProposal(id *big.Int) (*types.Transaction, error) {
	return _ProposalContract.Contract.ApproveProposal(&_ProposalContract.TransactOpts, id)
}

// ApproveProposal is a paid mutator transaction binding the contract method 0x98951b56.
//
// Solidity: function approveProposal(uint256 id) returns()
func (_ProposalContract *ProposalContractTransactorSession) ApproveProposal(id *big.Int) (*types.Transaction, error) {
	return _ProposalContract.Contract.ApproveProposal(&_ProposalContract.TransactOpts, id)
}

// CancelProposal is a paid mutator transaction binding the contract method 0xe0a8f6f5.
//
// Solidity: function cancelProposal(uint256 id) returns()
func (_ProposalContract *ProposalContractTransactor) CancelProposal(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _ProposalContract.contract.Transact(opts, "cancelProposal", id)
}

// CancelProposal is a paid mutator transaction binding the contract method 0xe0a8f6f5.
//
// Solidity: function cancelProposal(uint256 id) returns()
func (_ProposalContract *ProposalContractSession) CancelProposal(id *big.Int) (*types.Transaction, error) {
	return _ProposalContract.Contract.CancelProposal(&_ProposalContract.TransactOpts, id)
}

// CancelProposal is a paid mutator transaction binding the contract method 0xe0a8f6f5.
//
// Solidity: function cancelProposal(uint256 id) returns()
func (_ProposalContract *ProposalContractTransactorSession) CancelProposal(id *big.Int) (*types.Transaction, error) {
	return _ProposalContract.Contract.CancelProposal(&_ProposalContract.TransactOpts, id)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ProposalContract *ProposalContractTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ProposalContract.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ProposalContract *ProposalContractSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ProposalContract.Contract.GrantRole(&_ProposalContract.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ProposalContract *ProposalContractTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ProposalContract.Contract.GrantRole(&_ProposalContract.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ProposalContract *ProposalContractTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProposalContract.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ProposalContract *ProposalContractSession) Initialize() (*types.Transaction, error) {
	return _ProposalContract.Contract.Initialize(&_ProposalContract.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ProposalContract *ProposalContractTransactorSession) Initialize() (*types.Transaction, error) {
	return _ProposalContract.Contract.Initialize(&_ProposalContract.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ProposalContract *ProposalContractTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProposalContract.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ProposalContract *ProposalContractSession) Pause() (*types.Transaction, error) {
	return _ProposalContract.Contract.Pause(&_ProposalContract.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ProposalContract *ProposalContractTransactorSession) Pause() (*types.Transaction, error) {
	return _ProposalContract.Contract.Pause(&_ProposalContract.TransactOpts)
}

// RejectProposal is a paid mutator transaction binding the contract method 0xbc28d878.
//
// Solidity: function rejectProposal(uint256 id) returns()
func (_ProposalContract *ProposalContractTransactor) RejectProposal(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _ProposalContract.contract.Transact(opts, "rejectProposal", id)
}

// RejectProposal is a paid mutator transaction binding the contract method 0xbc28d878.
//
// Solidity: function rejectProposal(uint256 id) returns()
func (_ProposalContract *ProposalContractSession) RejectProposal(id *big.Int) (*types.Transaction, error) {
	return _ProposalContract.Contract.RejectProposal(&_ProposalContract.TransactOpts, id)
}

// RejectProposal is a paid mutator transaction binding the contract method 0xbc28d878.
//
// Solidity: function rejectProposal(uint256 id) returns()
func (_ProposalContract *ProposalContractTransactorSession) RejectProposal(id *big.Int) (*types.Transaction, error) {
	return _ProposalContract.Contract.RejectProposal(&_ProposalContract.TransactOpts, id)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ProposalContract *ProposalContractTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ProposalContract.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ProposalContract *ProposalContractSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ProposalContract.Contract.RenounceRole(&_ProposalContract.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_ProposalContract *ProposalContractTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _ProposalContract.Contract.RenounceRole(&_ProposalContract.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ProposalContract *ProposalContractTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ProposalContract.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ProposalContract *ProposalContractSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ProposalContract.Contract.RevokeRole(&_ProposalContract.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ProposalContract *ProposalContractTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ProposalContract.Contract.RevokeRole(&_ProposalContract.TransactOpts, role, account)
}

// SubmitProposal is a paid mutator transaction binding the contract method 0x60b31552.
//
// Solidity: function submitProposal(uint8 proposalType, string title, string description, string uri) returns(uint256)
func (_ProposalContract *ProposalContractTransactor) SubmitProposal(opts *bind.TransactOpts, proposalType uint8, title string, description string, uri string) (*types.Transaction, error) {
	return _ProposalContract.contract.Transact(opts, "submitProposal", proposalType, title, description, uri)
}

// SubmitProposal is a paid mutator transaction binding the contract method 0x60b31552.
//
// Solidity: function submitProposal(uint8 proposalType, string title, string description, string uri) returns(uint256)
func (_ProposalContract *ProposalContractSession) SubmitProposal(proposalType uint8, title string, description string, uri string) (*types.Transaction, error) {
	return _ProposalContract.Contract.SubmitProposal(&_ProposalContract.TransactOpts, proposalType, title, description, uri)
}

// SubmitProposal is a paid mutator transaction binding the contract method 0x60b31552.
//
// Solidity: function submitProposal(uint8 proposalType, string title, string description, string uri) returns(uint256)
func (_ProposalContract *ProposalContractTransactorSession) SubmitProposal(proposalType uint8, title string, description string, uri string) (*types.Transaction, error) {
	return _ProposalContract.Contract.SubmitProposal(&_ProposalContract.TransactOpts, proposalType, title, description, uri)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ProposalContract *ProposalContractTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProposalContract.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ProposalContract *ProposalContractSession) Unpause() (*types.Transaction, error) {
	return _ProposalContract.Contract.Unpause(&_ProposalContract.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ProposalContract *ProposalContractTransactorSession) Unpause() (*types.Transaction, error) {
	return _ProposalContract.Contract.Unpause(&_ProposalContract.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ProposalContract *ProposalContractTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ProposalContract.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ProposalContract *ProposalContractSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ProposalContract.Contract.UpgradeToAndCall(&_ProposalContract.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ProposalContract *ProposalContractTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ProposalContract.Contract.UpgradeToAndCall(&_ProposalContract.TransactOpts, newImplementation, data)
}

// ProposalContractInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ProposalContract contract.
type ProposalContractInitializedIterator struct {
	Event *ProposalContractInitialized // Event containing the contract specifics and raw log

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
func (it *ProposalContractInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProposalContractInitialized)
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
		it.Event = new(ProposalContractInitialized)
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
func (it *ProposalContractInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProposalContractInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProposalContractInitialized represents a Initialized event raised by the ProposalContract contract.
type ProposalContractInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ProposalContract *ProposalContractFilterer) FilterInitialized(opts *bind.FilterOpts) (*ProposalContractInitializedIterator, error) {

	logs, sub, err := _ProposalContract.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ProposalContractInitializedIterator{contract: _ProposalContract.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ProposalContract *ProposalContractFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ProposalContractInitialized) (event.Subscription, error) {

	logs, sub, err := _ProposalContract.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProposalContractInitialized)
				if err := _ProposalContract.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ProposalContract *ProposalContractFilterer) ParseInitialized(log types.Log) (*ProposalContractInitialized, error) {
	event := new(ProposalContractInitialized)
	if err := _ProposalContract.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProposalContractPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the ProposalContract contract.
type ProposalContractPausedIterator struct {
	Event *ProposalContractPaused // Event containing the contract specifics and raw log

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
func (it *ProposalContractPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProposalContractPaused)
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
		it.Event = new(ProposalContractPaused)
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
func (it *ProposalContractPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProposalContractPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProposalContractPaused represents a Paused event raised by the ProposalContract contract.
type ProposalContractPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ProposalContract *ProposalContractFilterer) FilterPaused(opts *bind.FilterOpts) (*ProposalContractPausedIterator, error) {

	logs, sub, err := _ProposalContract.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &ProposalContractPausedIterator{contract: _ProposalContract.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ProposalContract *ProposalContractFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ProposalContractPaused) (event.Subscription, error) {

	logs, sub, err := _ProposalContract.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProposalContractPaused)
				if err := _ProposalContract.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ProposalContract *ProposalContractFilterer) ParsePaused(log types.Log) (*ProposalContractPaused, error) {
	event := new(ProposalContractPaused)
	if err := _ProposalContract.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProposalContractProposalApprovedIterator is returned from FilterProposalApproved and is used to iterate over the raw logs and unpacked data for ProposalApproved events raised by the ProposalContract contract.
type ProposalContractProposalApprovedIterator struct {
	Event *ProposalContractProposalApproved // Event containing the contract specifics and raw log

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
func (it *ProposalContractProposalApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProposalContractProposalApproved)
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
		it.Event = new(ProposalContractProposalApproved)
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
func (it *ProposalContractProposalApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProposalContractProposalApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProposalContractProposalApproved represents a ProposalApproved event raised by the ProposalContract contract.
type ProposalContractProposalApproved struct {
	Id       *big.Int
	Reviewer common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterProposalApproved is a free log retrieval operation binding the contract event 0x049c28adfe50bcf1b76fd95273b6a24566b9f377e52fddc653c3355248dad07a.
//
// Solidity: event ProposalApproved(uint256 indexed id, address indexed reviewer)
func (_ProposalContract *ProposalContractFilterer) FilterProposalApproved(opts *bind.FilterOpts, id []*big.Int, reviewer []common.Address) (*ProposalContractProposalApprovedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var reviewerRule []interface{}
	for _, reviewerItem := range reviewer {
		reviewerRule = append(reviewerRule, reviewerItem)
	}

	logs, sub, err := _ProposalContract.contract.FilterLogs(opts, "ProposalApproved", idRule, reviewerRule)
	if err != nil {
		return nil, err
	}
	return &ProposalContractProposalApprovedIterator{contract: _ProposalContract.contract, event: "ProposalApproved", logs: logs, sub: sub}, nil
}

// WatchProposalApproved is a free log subscription operation binding the contract event 0x049c28adfe50bcf1b76fd95273b6a24566b9f377e52fddc653c3355248dad07a.
//
// Solidity: event ProposalApproved(uint256 indexed id, address indexed reviewer)
func (_ProposalContract *ProposalContractFilterer) WatchProposalApproved(opts *bind.WatchOpts, sink chan<- *ProposalContractProposalApproved, id []*big.Int, reviewer []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var reviewerRule []interface{}
	for _, reviewerItem := range reviewer {
		reviewerRule = append(reviewerRule, reviewerItem)
	}

	logs, sub, err := _ProposalContract.contract.WatchLogs(opts, "ProposalApproved", idRule, reviewerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProposalContractProposalApproved)
				if err := _ProposalContract.contract.UnpackLog(event, "ProposalApproved", log); err != nil {
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

// ParseProposalApproved is a log parse operation binding the contract event 0x049c28adfe50bcf1b76fd95273b6a24566b9f377e52fddc653c3355248dad07a.
//
// Solidity: event ProposalApproved(uint256 indexed id, address indexed reviewer)
func (_ProposalContract *ProposalContractFilterer) ParseProposalApproved(log types.Log) (*ProposalContractProposalApproved, error) {
	event := new(ProposalContractProposalApproved)
	if err := _ProposalContract.contract.UnpackLog(event, "ProposalApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProposalContractProposalCancelledIterator is returned from FilterProposalCancelled and is used to iterate over the raw logs and unpacked data for ProposalCancelled events raised by the ProposalContract contract.
type ProposalContractProposalCancelledIterator struct {
	Event *ProposalContractProposalCancelled // Event containing the contract specifics and raw log

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
func (it *ProposalContractProposalCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProposalContractProposalCancelled)
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
		it.Event = new(ProposalContractProposalCancelled)
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
func (it *ProposalContractProposalCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProposalContractProposalCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProposalContractProposalCancelled represents a ProposalCancelled event raised by the ProposalContract contract.
type ProposalContractProposalCancelled struct {
	Id       *big.Int
	Proposer common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterProposalCancelled is a free log retrieval operation binding the contract event 0x74c34a008ce735d9fcf0bd03a9b238d212ad4c441c020661f4ffbb6442645b85.
//
// Solidity: event ProposalCancelled(uint256 indexed id, address indexed proposer)
func (_ProposalContract *ProposalContractFilterer) FilterProposalCancelled(opts *bind.FilterOpts, id []*big.Int, proposer []common.Address) (*ProposalContractProposalCancelledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _ProposalContract.contract.FilterLogs(opts, "ProposalCancelled", idRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return &ProposalContractProposalCancelledIterator{contract: _ProposalContract.contract, event: "ProposalCancelled", logs: logs, sub: sub}, nil
}

// WatchProposalCancelled is a free log subscription operation binding the contract event 0x74c34a008ce735d9fcf0bd03a9b238d212ad4c441c020661f4ffbb6442645b85.
//
// Solidity: event ProposalCancelled(uint256 indexed id, address indexed proposer)
func (_ProposalContract *ProposalContractFilterer) WatchProposalCancelled(opts *bind.WatchOpts, sink chan<- *ProposalContractProposalCancelled, id []*big.Int, proposer []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _ProposalContract.contract.WatchLogs(opts, "ProposalCancelled", idRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProposalContractProposalCancelled)
				if err := _ProposalContract.contract.UnpackLog(event, "ProposalCancelled", log); err != nil {
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

// ParseProposalCancelled is a log parse operation binding the contract event 0x74c34a008ce735d9fcf0bd03a9b238d212ad4c441c020661f4ffbb6442645b85.
//
// Solidity: event ProposalCancelled(uint256 indexed id, address indexed proposer)
func (_ProposalContract *ProposalContractFilterer) ParseProposalCancelled(log types.Log) (*ProposalContractProposalCancelled, error) {
	event := new(ProposalContractProposalCancelled)
	if err := _ProposalContract.contract.UnpackLog(event, "ProposalCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProposalContractProposalRejectedIterator is returned from FilterProposalRejected and is used to iterate over the raw logs and unpacked data for ProposalRejected events raised by the ProposalContract contract.
type ProposalContractProposalRejectedIterator struct {
	Event *ProposalContractProposalRejected // Event containing the contract specifics and raw log

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
func (it *ProposalContractProposalRejectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProposalContractProposalRejected)
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
		it.Event = new(ProposalContractProposalRejected)
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
func (it *ProposalContractProposalRejectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProposalContractProposalRejectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProposalContractProposalRejected represents a ProposalRejected event raised by the ProposalContract contract.
type ProposalContractProposalRejected struct {
	Id       *big.Int
	Reviewer common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterProposalRejected is a free log retrieval operation binding the contract event 0xff556cafc8033c441c6fea0e40d12f0ec0c8c9168f6bac576e84800331b1a52f.
//
// Solidity: event ProposalRejected(uint256 indexed id, address indexed reviewer)
func (_ProposalContract *ProposalContractFilterer) FilterProposalRejected(opts *bind.FilterOpts, id []*big.Int, reviewer []common.Address) (*ProposalContractProposalRejectedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var reviewerRule []interface{}
	for _, reviewerItem := range reviewer {
		reviewerRule = append(reviewerRule, reviewerItem)
	}

	logs, sub, err := _ProposalContract.contract.FilterLogs(opts, "ProposalRejected", idRule, reviewerRule)
	if err != nil {
		return nil, err
	}
	return &ProposalContractProposalRejectedIterator{contract: _ProposalContract.contract, event: "ProposalRejected", logs: logs, sub: sub}, nil
}

// WatchProposalRejected is a free log subscription operation binding the contract event 0xff556cafc8033c441c6fea0e40d12f0ec0c8c9168f6bac576e84800331b1a52f.
//
// Solidity: event ProposalRejected(uint256 indexed id, address indexed reviewer)
func (_ProposalContract *ProposalContractFilterer) WatchProposalRejected(opts *bind.WatchOpts, sink chan<- *ProposalContractProposalRejected, id []*big.Int, reviewer []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var reviewerRule []interface{}
	for _, reviewerItem := range reviewer {
		reviewerRule = append(reviewerRule, reviewerItem)
	}

	logs, sub, err := _ProposalContract.contract.WatchLogs(opts, "ProposalRejected", idRule, reviewerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProposalContractProposalRejected)
				if err := _ProposalContract.contract.UnpackLog(event, "ProposalRejected", log); err != nil {
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

// ParseProposalRejected is a log parse operation binding the contract event 0xff556cafc8033c441c6fea0e40d12f0ec0c8c9168f6bac576e84800331b1a52f.
//
// Solidity: event ProposalRejected(uint256 indexed id, address indexed reviewer)
func (_ProposalContract *ProposalContractFilterer) ParseProposalRejected(log types.Log) (*ProposalContractProposalRejected, error) {
	event := new(ProposalContractProposalRejected)
	if err := _ProposalContract.contract.UnpackLog(event, "ProposalRejected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProposalContractProposalStatusChangedIterator is returned from FilterProposalStatusChanged and is used to iterate over the raw logs and unpacked data for ProposalStatusChanged events raised by the ProposalContract contract.
type ProposalContractProposalStatusChangedIterator struct {
	Event *ProposalContractProposalStatusChanged // Event containing the contract specifics and raw log

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
func (it *ProposalContractProposalStatusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProposalContractProposalStatusChanged)
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
		it.Event = new(ProposalContractProposalStatusChanged)
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
func (it *ProposalContractProposalStatusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProposalContractProposalStatusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProposalContractProposalStatusChanged represents a ProposalStatusChanged event raised by the ProposalContract contract.
type ProposalContractProposalStatusChanged struct {
	Id        *big.Int
	OldStatus uint8
	NewStatus uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterProposalStatusChanged is a free log retrieval operation binding the contract event 0xb0978968e606225988b6464f99ce7e86b0e8770d80eaec1369bb3780e05c6438.
//
// Solidity: event ProposalStatusChanged(uint256 indexed id, uint8 oldStatus, uint8 newStatus)
func (_ProposalContract *ProposalContractFilterer) FilterProposalStatusChanged(opts *bind.FilterOpts, id []*big.Int) (*ProposalContractProposalStatusChangedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ProposalContract.contract.FilterLogs(opts, "ProposalStatusChanged", idRule)
	if err != nil {
		return nil, err
	}
	return &ProposalContractProposalStatusChangedIterator{contract: _ProposalContract.contract, event: "ProposalStatusChanged", logs: logs, sub: sub}, nil
}

// WatchProposalStatusChanged is a free log subscription operation binding the contract event 0xb0978968e606225988b6464f99ce7e86b0e8770d80eaec1369bb3780e05c6438.
//
// Solidity: event ProposalStatusChanged(uint256 indexed id, uint8 oldStatus, uint8 newStatus)
func (_ProposalContract *ProposalContractFilterer) WatchProposalStatusChanged(opts *bind.WatchOpts, sink chan<- *ProposalContractProposalStatusChanged, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ProposalContract.contract.WatchLogs(opts, "ProposalStatusChanged", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProposalContractProposalStatusChanged)
				if err := _ProposalContract.contract.UnpackLog(event, "ProposalStatusChanged", log); err != nil {
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

// ParseProposalStatusChanged is a log parse operation binding the contract event 0xb0978968e606225988b6464f99ce7e86b0e8770d80eaec1369bb3780e05c6438.
//
// Solidity: event ProposalStatusChanged(uint256 indexed id, uint8 oldStatus, uint8 newStatus)
func (_ProposalContract *ProposalContractFilterer) ParseProposalStatusChanged(log types.Log) (*ProposalContractProposalStatusChanged, error) {
	event := new(ProposalContractProposalStatusChanged)
	if err := _ProposalContract.contract.UnpackLog(event, "ProposalStatusChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProposalContractProposalSubmittedIterator is returned from FilterProposalSubmitted and is used to iterate over the raw logs and unpacked data for ProposalSubmitted events raised by the ProposalContract contract.
type ProposalContractProposalSubmittedIterator struct {
	Event *ProposalContractProposalSubmitted // Event containing the contract specifics and raw log

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
func (it *ProposalContractProposalSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProposalContractProposalSubmitted)
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
		it.Event = new(ProposalContractProposalSubmitted)
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
func (it *ProposalContractProposalSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProposalContractProposalSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProposalContractProposalSubmitted represents a ProposalSubmitted event raised by the ProposalContract contract.
type ProposalContractProposalSubmitted struct {
	Id           *big.Int
	Proposer     common.Address
	ProposalType uint8
	Title        string
	Uri          string
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterProposalSubmitted is a free log retrieval operation binding the contract event 0x5973e60d95ce805a469be386207eb73f0b25108c2436513caf143dd125161687.
//
// Solidity: event ProposalSubmitted(uint256 indexed id, address indexed proposer, uint8 proposalType, string title, string uri)
func (_ProposalContract *ProposalContractFilterer) FilterProposalSubmitted(opts *bind.FilterOpts, id []*big.Int, proposer []common.Address) (*ProposalContractProposalSubmittedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _ProposalContract.contract.FilterLogs(opts, "ProposalSubmitted", idRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return &ProposalContractProposalSubmittedIterator{contract: _ProposalContract.contract, event: "ProposalSubmitted", logs: logs, sub: sub}, nil
}

// WatchProposalSubmitted is a free log subscription operation binding the contract event 0x5973e60d95ce805a469be386207eb73f0b25108c2436513caf143dd125161687.
//
// Solidity: event ProposalSubmitted(uint256 indexed id, address indexed proposer, uint8 proposalType, string title, string uri)
func (_ProposalContract *ProposalContractFilterer) WatchProposalSubmitted(opts *bind.WatchOpts, sink chan<- *ProposalContractProposalSubmitted, id []*big.Int, proposer []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _ProposalContract.contract.WatchLogs(opts, "ProposalSubmitted", idRule, proposerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProposalContractProposalSubmitted)
				if err := _ProposalContract.contract.UnpackLog(event, "ProposalSubmitted", log); err != nil {
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

// ParseProposalSubmitted is a log parse operation binding the contract event 0x5973e60d95ce805a469be386207eb73f0b25108c2436513caf143dd125161687.
//
// Solidity: event ProposalSubmitted(uint256 indexed id, address indexed proposer, uint8 proposalType, string title, string uri)
func (_ProposalContract *ProposalContractFilterer) ParseProposalSubmitted(log types.Log) (*ProposalContractProposalSubmitted, error) {
	event := new(ProposalContractProposalSubmitted)
	if err := _ProposalContract.contract.UnpackLog(event, "ProposalSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProposalContractRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the ProposalContract contract.
type ProposalContractRoleAdminChangedIterator struct {
	Event *ProposalContractRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ProposalContractRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProposalContractRoleAdminChanged)
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
		it.Event = new(ProposalContractRoleAdminChanged)
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
func (it *ProposalContractRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProposalContractRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProposalContractRoleAdminChanged represents a RoleAdminChanged event raised by the ProposalContract contract.
type ProposalContractRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ProposalContract *ProposalContractFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ProposalContractRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _ProposalContract.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ProposalContractRoleAdminChangedIterator{contract: _ProposalContract.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ProposalContract *ProposalContractFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ProposalContractRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _ProposalContract.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProposalContractRoleAdminChanged)
				if err := _ProposalContract.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ProposalContract *ProposalContractFilterer) ParseRoleAdminChanged(log types.Log) (*ProposalContractRoleAdminChanged, error) {
	event := new(ProposalContractRoleAdminChanged)
	if err := _ProposalContract.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProposalContractRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the ProposalContract contract.
type ProposalContractRoleGrantedIterator struct {
	Event *ProposalContractRoleGranted // Event containing the contract specifics and raw log

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
func (it *ProposalContractRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProposalContractRoleGranted)
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
		it.Event = new(ProposalContractRoleGranted)
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
func (it *ProposalContractRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProposalContractRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProposalContractRoleGranted represents a RoleGranted event raised by the ProposalContract contract.
type ProposalContractRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ProposalContract *ProposalContractFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ProposalContractRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ProposalContract.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ProposalContractRoleGrantedIterator{contract: _ProposalContract.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ProposalContract *ProposalContractFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ProposalContractRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ProposalContract.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProposalContractRoleGranted)
				if err := _ProposalContract.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ProposalContract *ProposalContractFilterer) ParseRoleGranted(log types.Log) (*ProposalContractRoleGranted, error) {
	event := new(ProposalContractRoleGranted)
	if err := _ProposalContract.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProposalContractRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the ProposalContract contract.
type ProposalContractRoleRevokedIterator struct {
	Event *ProposalContractRoleRevoked // Event containing the contract specifics and raw log

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
func (it *ProposalContractRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProposalContractRoleRevoked)
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
		it.Event = new(ProposalContractRoleRevoked)
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
func (it *ProposalContractRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProposalContractRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProposalContractRoleRevoked represents a RoleRevoked event raised by the ProposalContract contract.
type ProposalContractRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ProposalContract *ProposalContractFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ProposalContractRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ProposalContract.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ProposalContractRoleRevokedIterator{contract: _ProposalContract.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ProposalContract *ProposalContractFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ProposalContractRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ProposalContract.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProposalContractRoleRevoked)
				if err := _ProposalContract.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ProposalContract *ProposalContractFilterer) ParseRoleRevoked(log types.Log) (*ProposalContractRoleRevoked, error) {
	event := new(ProposalContractRoleRevoked)
	if err := _ProposalContract.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProposalContractUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the ProposalContract contract.
type ProposalContractUnpausedIterator struct {
	Event *ProposalContractUnpaused // Event containing the contract specifics and raw log

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
func (it *ProposalContractUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProposalContractUnpaused)
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
		it.Event = new(ProposalContractUnpaused)
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
func (it *ProposalContractUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProposalContractUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProposalContractUnpaused represents a Unpaused event raised by the ProposalContract contract.
type ProposalContractUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ProposalContract *ProposalContractFilterer) FilterUnpaused(opts *bind.FilterOpts) (*ProposalContractUnpausedIterator, error) {

	logs, sub, err := _ProposalContract.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &ProposalContractUnpausedIterator{contract: _ProposalContract.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ProposalContract *ProposalContractFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ProposalContractUnpaused) (event.Subscription, error) {

	logs, sub, err := _ProposalContract.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProposalContractUnpaused)
				if err := _ProposalContract.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ProposalContract *ProposalContractFilterer) ParseUnpaused(log types.Log) (*ProposalContractUnpaused, error) {
	event := new(ProposalContractUnpaused)
	if err := _ProposalContract.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProposalContractUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ProposalContract contract.
type ProposalContractUpgradedIterator struct {
	Event *ProposalContractUpgraded // Event containing the contract specifics and raw log

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
func (it *ProposalContractUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProposalContractUpgraded)
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
		it.Event = new(ProposalContractUpgraded)
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
func (it *ProposalContractUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProposalContractUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProposalContractUpgraded represents a Upgraded event raised by the ProposalContract contract.
type ProposalContractUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ProposalContract *ProposalContractFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ProposalContractUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ProposalContract.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ProposalContractUpgradedIterator{contract: _ProposalContract.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ProposalContract *ProposalContractFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ProposalContractUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ProposalContract.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProposalContractUpgraded)
				if err := _ProposalContract.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_ProposalContract *ProposalContractFilterer) ParseUpgraded(log types.Log) (*ProposalContractUpgraded, error) {
	event := new(ProposalContractUpgraded)
	if err := _ProposalContract.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
