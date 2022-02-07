// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package zcontract

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
)

// ZionDataSubmitBlockInfo is an auto generated low-level Go binding around an user-defined struct.
type ZionDataSubmitBlockInfo struct {
	BlockType             uint8
	BlockSize             uint16
	Data                  []byte
	Proof                 [8]*big.Int
	StoreBlockInfoOnchain bool
	AuxiliaryData         []byte
}

// ZionCMetaData contains all meta data concerning the ZionC contract.
var ZionCMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"priorityReqId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"tokenId\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"amount\",\"type\":\"uint96\"}],\"name\":\"NewDepositRequest\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_zionAddr\",\"type\":\"address\"}],\"name\":\"depositCTXC\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMerkleRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"getPendingBalance\",\"outputs\":[{\"internalType\":\"uint96\",\"name\":\"\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"initializationParameters\",\"type\":\"bytes\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"blockType\",\"type\":\"uint8\"},{\"internalType\":\"uint16\",\"name\":\"blockSize\",\"type\":\"uint16\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256[8]\",\"name\":\"proof\",\"type\":\"uint256[8]\"},{\"internalType\":\"bool\",\"name\":\"storeBlockInfoOnchain\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"auxiliaryData\",\"type\":\"bytes\"}],\"internalType\":\"structZionData.SubmitBlockInfo[]\",\"name\":\"_blocks\",\"type\":\"tuple[]\"}],\"name\":\"submitBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"upgradeParameters\",\"type\":\"bytes\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"withdrawPendingBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ZionCABI is the input ABI used to generate the binding from.
// Deprecated: Use ZionCMetaData.ABI instead.
var ZionCABI = ZionCMetaData.ABI

// ZionC is an auto generated Go binding around an Ethereum contract.
type ZionC struct {
	ZionCCaller     // Read-only binding to the contract
	ZionCTransactor // Write-only binding to the contract
	ZionCFilterer   // Log filterer for contract events
}

// ZionCCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZionCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZionCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZionCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZionCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZionCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZionCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZionCSession struct {
	Contract     *ZionC            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZionCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZionCCallerSession struct {
	Contract *ZionCCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ZionCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZionCTransactorSession struct {
	Contract     *ZionCTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZionCRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZionCRaw struct {
	Contract *ZionC // Generic contract binding to access the raw methods on
}

// ZionCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZionCCallerRaw struct {
	Contract *ZionCCaller // Generic read-only contract binding to access the raw methods on
}

// ZionCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZionCTransactorRaw struct {
	Contract *ZionCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZionC creates a new instance of ZionC, bound to a specific deployed contract.
func NewZionC(address common.Address, backend bind.ContractBackend) (*ZionC, error) {
	contract, err := bindZionC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZionC{ZionCCaller: ZionCCaller{contract: contract}, ZionCTransactor: ZionCTransactor{contract: contract}, ZionCFilterer: ZionCFilterer{contract: contract}}, nil
}

// NewZionCCaller creates a new read-only instance of ZionC, bound to a specific deployed contract.
func NewZionCCaller(address common.Address, caller bind.ContractCaller) (*ZionCCaller, error) {
	contract, err := bindZionC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZionCCaller{contract: contract}, nil
}

// NewZionCTransactor creates a new write-only instance of ZionC, bound to a specific deployed contract.
func NewZionCTransactor(address common.Address, transactor bind.ContractTransactor) (*ZionCTransactor, error) {
	contract, err := bindZionC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZionCTransactor{contract: contract}, nil
}

// NewZionCFilterer creates a new log filterer instance of ZionC, bound to a specific deployed contract.
func NewZionCFilterer(address common.Address, filterer bind.ContractFilterer) (*ZionCFilterer, error) {
	contract, err := bindZionC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZionCFilterer{contract: contract}, nil
}

// bindZionC binds a generic wrapper to an already deployed contract.
func bindZionC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ZionCABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZionC *ZionCRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZionC.Contract.ZionCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZionC *ZionCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZionC.Contract.ZionCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZionC *ZionCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZionC.Contract.ZionCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZionC *ZionCCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZionC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZionC *ZionCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZionC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZionC *ZionCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZionC.Contract.contract.Transact(opts, method, params...)
}

// GetBlockHeight is a free data retrieval call binding the contract method 0x7bb96acb.
//
// Solidity: function getBlockHeight() view returns(uint256)
func (_ZionC *ZionCCaller) GetBlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZionC.contract.Call(opts, &out, "getBlockHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBlockHeight is a free data retrieval call binding the contract method 0x7bb96acb.
//
// Solidity: function getBlockHeight() view returns(uint256)
func (_ZionC *ZionCSession) GetBlockHeight() (*big.Int, error) {
	return _ZionC.Contract.GetBlockHeight(&_ZionC.CallOpts)
}

// GetBlockHeight is a free data retrieval call binding the contract method 0x7bb96acb.
//
// Solidity: function getBlockHeight() view returns(uint256)
func (_ZionC *ZionCCallerSession) GetBlockHeight() (*big.Int, error) {
	return _ZionC.Contract.GetBlockHeight(&_ZionC.CallOpts)
}

// GetMerkleRoot is a free data retrieval call binding the contract method 0x49590657.
//
// Solidity: function getMerkleRoot() view returns(bytes32)
func (_ZionC *ZionCCaller) GetMerkleRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ZionC.contract.Call(opts, &out, "getMerkleRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetMerkleRoot is a free data retrieval call binding the contract method 0x49590657.
//
// Solidity: function getMerkleRoot() view returns(bytes32)
func (_ZionC *ZionCSession) GetMerkleRoot() ([32]byte, error) {
	return _ZionC.Contract.GetMerkleRoot(&_ZionC.CallOpts)
}

// GetMerkleRoot is a free data retrieval call binding the contract method 0x49590657.
//
// Solidity: function getMerkleRoot() view returns(bytes32)
func (_ZionC *ZionCCallerSession) GetMerkleRoot() ([32]byte, error) {
	return _ZionC.Contract.GetMerkleRoot(&_ZionC.CallOpts)
}

// GetPendingBalance is a free data retrieval call binding the contract method 0x5aca41f6.
//
// Solidity: function getPendingBalance(address _addr, address _token) view returns(uint96)
func (_ZionC *ZionCCaller) GetPendingBalance(opts *bind.CallOpts, _addr common.Address, _token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZionC.contract.Call(opts, &out, "getPendingBalance", _addr, _token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPendingBalance is a free data retrieval call binding the contract method 0x5aca41f6.
//
// Solidity: function getPendingBalance(address _addr, address _token) view returns(uint96)
func (_ZionC *ZionCSession) GetPendingBalance(_addr common.Address, _token common.Address) (*big.Int, error) {
	return _ZionC.Contract.GetPendingBalance(&_ZionC.CallOpts, _addr, _token)
}

// GetPendingBalance is a free data retrieval call binding the contract method 0x5aca41f6.
//
// Solidity: function getPendingBalance(address _addr, address _token) view returns(uint96)
func (_ZionC *ZionCCallerSession) GetPendingBalance(_addr common.Address, _token common.Address) (*big.Int, error) {
	return _ZionC.Contract.GetPendingBalance(&_ZionC.CallOpts, _addr, _token)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_ZionC *ZionCCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ZionC.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_ZionC *ZionCSession) Version() (string, error) {
	return _ZionC.Contract.Version(&_ZionC.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() pure returns(string)
func (_ZionC *ZionCCallerSession) Version() (string, error) {
	return _ZionC.Contract.Version(&_ZionC.CallOpts)
}

// DepositCTXC is a paid mutator transaction binding the contract method 0x6feec5da.
//
// Solidity: function depositCTXC(address _zionAddr) payable returns()
func (_ZionC *ZionCTransactor) DepositCTXC(opts *bind.TransactOpts, _zionAddr common.Address) (*types.Transaction, error) {
	return _ZionC.contract.Transact(opts, "depositCTXC", _zionAddr)
}

// DepositCTXC is a paid mutator transaction binding the contract method 0x6feec5da.
//
// Solidity: function depositCTXC(address _zionAddr) payable returns()
func (_ZionC *ZionCSession) DepositCTXC(_zionAddr common.Address) (*types.Transaction, error) {
	return _ZionC.Contract.DepositCTXC(&_ZionC.TransactOpts, _zionAddr)
}

// DepositCTXC is a paid mutator transaction binding the contract method 0x6feec5da.
//
// Solidity: function depositCTXC(address _zionAddr) payable returns()
func (_ZionC *ZionCTransactorSession) DepositCTXC(_zionAddr common.Address) (*types.Transaction, error) {
	return _ZionC.Contract.DepositCTXC(&_ZionC.TransactOpts, _zionAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0x439fab91.
//
// Solidity: function initialize(bytes initializationParameters) returns()
func (_ZionC *ZionCTransactor) Initialize(opts *bind.TransactOpts, initializationParameters []byte) (*types.Transaction, error) {
	return _ZionC.contract.Transact(opts, "initialize", initializationParameters)
}

// Initialize is a paid mutator transaction binding the contract method 0x439fab91.
//
// Solidity: function initialize(bytes initializationParameters) returns()
func (_ZionC *ZionCSession) Initialize(initializationParameters []byte) (*types.Transaction, error) {
	return _ZionC.Contract.Initialize(&_ZionC.TransactOpts, initializationParameters)
}

// Initialize is a paid mutator transaction binding the contract method 0x439fab91.
//
// Solidity: function initialize(bytes initializationParameters) returns()
func (_ZionC *ZionCTransactorSession) Initialize(initializationParameters []byte) (*types.Transaction, error) {
	return _ZionC.Contract.Initialize(&_ZionC.TransactOpts, initializationParameters)
}

// SubmitBlocks is a paid mutator transaction binding the contract method 0xc1e9fcc5.
//
// Solidity: function submitBlocks((uint8,uint16,bytes,uint256[8],bool,bytes)[] _blocks) returns()
func (_ZionC *ZionCTransactor) SubmitBlocks(opts *bind.TransactOpts, _blocks []ZionDataSubmitBlockInfo) (*types.Transaction, error) {
	return _ZionC.contract.Transact(opts, "submitBlocks", _blocks)
}

// SubmitBlocks is a paid mutator transaction binding the contract method 0xc1e9fcc5.
//
// Solidity: function submitBlocks((uint8,uint16,bytes,uint256[8],bool,bytes)[] _blocks) returns()
func (_ZionC *ZionCSession) SubmitBlocks(_blocks []ZionDataSubmitBlockInfo) (*types.Transaction, error) {
	return _ZionC.Contract.SubmitBlocks(&_ZionC.TransactOpts, _blocks)
}

// SubmitBlocks is a paid mutator transaction binding the contract method 0xc1e9fcc5.
//
// Solidity: function submitBlocks((uint8,uint16,bytes,uint256[8],bool,bytes)[] _blocks) returns()
func (_ZionC *ZionCTransactorSession) SubmitBlocks(_blocks []ZionDataSubmitBlockInfo) (*types.Transaction, error) {
	return _ZionC.Contract.SubmitBlocks(&_ZionC.TransactOpts, _blocks)
}

// Upgrade is a paid mutator transaction binding the contract method 0x25394645.
//
// Solidity: function upgrade(bytes upgradeParameters) returns()
func (_ZionC *ZionCTransactor) Upgrade(opts *bind.TransactOpts, upgradeParameters []byte) (*types.Transaction, error) {
	return _ZionC.contract.Transact(opts, "upgrade", upgradeParameters)
}

// Upgrade is a paid mutator transaction binding the contract method 0x25394645.
//
// Solidity: function upgrade(bytes upgradeParameters) returns()
func (_ZionC *ZionCSession) Upgrade(upgradeParameters []byte) (*types.Transaction, error) {
	return _ZionC.Contract.Upgrade(&_ZionC.TransactOpts, upgradeParameters)
}

// Upgrade is a paid mutator transaction binding the contract method 0x25394645.
//
// Solidity: function upgrade(bytes upgradeParameters) returns()
func (_ZionC *ZionCTransactorSession) Upgrade(upgradeParameters []byte) (*types.Transaction, error) {
	return _ZionC.Contract.Upgrade(&_ZionC.TransactOpts, upgradeParameters)
}

// WithdrawPendingBalance is a paid mutator transaction binding the contract method 0x8ddb4eba.
//
// Solidity: function withdrawPendingBalance(address _owner, address _token) returns()
func (_ZionC *ZionCTransactor) WithdrawPendingBalance(opts *bind.TransactOpts, _owner common.Address, _token common.Address) (*types.Transaction, error) {
	return _ZionC.contract.Transact(opts, "withdrawPendingBalance", _owner, _token)
}

// WithdrawPendingBalance is a paid mutator transaction binding the contract method 0x8ddb4eba.
//
// Solidity: function withdrawPendingBalance(address _owner, address _token) returns()
func (_ZionC *ZionCSession) WithdrawPendingBalance(_owner common.Address, _token common.Address) (*types.Transaction, error) {
	return _ZionC.Contract.WithdrawPendingBalance(&_ZionC.TransactOpts, _owner, _token)
}

// WithdrawPendingBalance is a paid mutator transaction binding the contract method 0x8ddb4eba.
//
// Solidity: function withdrawPendingBalance(address _owner, address _token) returns()
func (_ZionC *ZionCTransactorSession) WithdrawPendingBalance(_owner common.Address, _token common.Address) (*types.Transaction, error) {
	return _ZionC.Contract.WithdrawPendingBalance(&_ZionC.TransactOpts, _owner, _token)
}

// ZionCNewDepositRequestIterator is returned from FilterNewDepositRequest and is used to iterate over the raw logs and unpacked data for NewDepositRequest events raised by the ZionC contract.
type ZionCNewDepositRequestIterator struct {
	Event *ZionCNewDepositRequest // Event containing the contract specifics and raw log

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
func (it *ZionCNewDepositRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZionCNewDepositRequest)
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
		it.Event = new(ZionCNewDepositRequest)
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
func (it *ZionCNewDepositRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZionCNewDepositRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZionCNewDepositRequest represents a NewDepositRequest event raised by the ZionC contract.
type ZionCNewDepositRequest struct {
	Sender        common.Address
	Receiver      common.Address
	PriorityReqId uint64
	TokenId       uint16
	Amount        *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterNewDepositRequest is a free log retrieval operation binding the contract event 0x53a959842d0c8b0e8e28011927f401d5d1258c0db5ebb6043fd6949220fdc778.
//
// Solidity: event NewDepositRequest(address sender, address receiver, uint64 priorityReqId, uint16 tokenId, uint96 amount)
func (_ZionC *ZionCFilterer) FilterNewDepositRequest(opts *bind.FilterOpts) (*ZionCNewDepositRequestIterator, error) {

	logs, sub, err := _ZionC.contract.FilterLogs(opts, "NewDepositRequest")
	if err != nil {
		return nil, err
	}
	return &ZionCNewDepositRequestIterator{contract: _ZionC.contract, event: "NewDepositRequest", logs: logs, sub: sub}, nil
}

// WatchNewDepositRequest is a free log subscription operation binding the contract event 0x53a959842d0c8b0e8e28011927f401d5d1258c0db5ebb6043fd6949220fdc778.
//
// Solidity: event NewDepositRequest(address sender, address receiver, uint64 priorityReqId, uint16 tokenId, uint96 amount)
func (_ZionC *ZionCFilterer) WatchNewDepositRequest(opts *bind.WatchOpts, sink chan<- *ZionCNewDepositRequest) (event.Subscription, error) {

	logs, sub, err := _ZionC.contract.WatchLogs(opts, "NewDepositRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZionCNewDepositRequest)
				if err := _ZionC.contract.UnpackLog(event, "NewDepositRequest", log); err != nil {
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

// ParseNewDepositRequest is a log parse operation binding the contract event 0x53a959842d0c8b0e8e28011927f401d5d1258c0db5ebb6043fd6949220fdc778.
//
// Solidity: event NewDepositRequest(address sender, address receiver, uint64 priorityReqId, uint16 tokenId, uint96 amount)
func (_ZionC *ZionCFilterer) ParseNewDepositRequest(log types.Log) (*ZionCNewDepositRequest, error) {
	event := new(ZionCNewDepositRequest)
	if err := _ZionC.contract.UnpackLog(event, "NewDepositRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
