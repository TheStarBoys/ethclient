// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ContractsABI is the input ABI used to generate the binding from.
const ContractsABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"counter\",\"type\":\"uint256\"}],\"name\":\"CounterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"arg1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"arg2\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"arg3\",\"type\":\"bytes\"}],\"name\":\"FuncEvent1\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"arg1\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"arg2\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"arg3\",\"type\":\"bytes\"}],\"name\":\"testFunc1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testReverted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ContractsBin is the compiled bytecode used for deploying new contracts.
var ContractsBin = "0x608060405234801561001057600080fd5b506103d3806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806361bc221a146100465780636c6dd6031461006457806388655d981461006e575b600080fd5b61004e6101ca565b6040518082815260200191505060405180910390f35b61006c6101d0565b005b6101c86004803603606081101561008457600080fd5b81019080803590602001906401000000008111156100a157600080fd5b8201836020820111156100b357600080fd5b803590602001918460018302840111640100000000831117156100d557600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290803590602001909291908035906020019064010000000081111561014257600080fd5b82018360208201111561015457600080fd5b8035906020019184600183028401116401000000008311171561017657600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050919291929050505061023e565b005b60005481565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600d8152602001807f746573742072657665727465640000000000000000000000000000000000000081525060200191505060405180910390fd5b7fee7ebd5ac9177b3cfe282c440d0220335dc60bc4472338132f06af7b4b9432fc838383604051808060200184815260200180602001838103835286818151815260200191508051906020019080838360005b838110156102ac578082015181840152602081019050610291565b50505050905090810190601f1680156102d95780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b838110156103125780820151818401526020810190506102f7565b50505050905090810190601f16801561033f5780820380516001836020036101000a031916815260200191505b509550505050505060405180910390a1600160008082825401925050819055507f4785d80d2593e2cb7a3331d31eb5106408bdde2aab0db9e9b616b036a1b6039d6000546040518082815260200191505060405180910390a150505056fea26469706673582212200a605c4881a6c7323cfd00843e81261800bb5015a497ee47a7d92a78d68c252364736f6c63430007060033"

// DeployContracts deploys a new Ethereum contract, binding an instance of Contracts to it.
func DeployContracts(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contracts, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ContractsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contracts{ContractsCaller: ContractsCaller{contract: contract}, ContractsTransactor: ContractsTransactor{contract: contract}, ContractsFilterer: ContractsFilterer{contract: contract}}, nil
}

// Contracts is an auto generated Go binding around an Ethereum contract.
type Contracts struct {
	ContractsCaller     // Read-only binding to the contract
	ContractsTransactor // Write-only binding to the contract
	ContractsFilterer   // Log filterer for contract events
}

// ContractsCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractsSession struct {
	Contract     *Contracts        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractsCallerSession struct {
	Contract *ContractsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ContractsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractsTransactorSession struct {
	Contract     *ContractsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ContractsRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractsRaw struct {
	Contract *Contracts // Generic contract binding to access the raw methods on
}

// ContractsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractsCallerRaw struct {
	Contract *ContractsCaller // Generic read-only contract binding to access the raw methods on
}

// ContractsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractsTransactorRaw struct {
	Contract *ContractsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContracts creates a new instance of Contracts, bound to a specific deployed contract.
func NewContracts(address common.Address, backend bind.ContractBackend) (*Contracts, error) {
	contract, err := bindContracts(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contracts{ContractsCaller: ContractsCaller{contract: contract}, ContractsTransactor: ContractsTransactor{contract: contract}, ContractsFilterer: ContractsFilterer{contract: contract}}, nil
}

// NewContractsCaller creates a new read-only instance of Contracts, bound to a specific deployed contract.
func NewContractsCaller(address common.Address, caller bind.ContractCaller) (*ContractsCaller, error) {
	contract, err := bindContracts(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractsCaller{contract: contract}, nil
}

// NewContractsTransactor creates a new write-only instance of Contracts, bound to a specific deployed contract.
func NewContractsTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractsTransactor, error) {
	contract, err := bindContracts(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractsTransactor{contract: contract}, nil
}

// NewContractsFilterer creates a new log filterer instance of Contracts, bound to a specific deployed contract.
func NewContractsFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractsFilterer, error) {
	contract, err := bindContracts(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractsFilterer{contract: contract}, nil
}

// bindContracts binds a generic wrapper to an already deployed contract.
func bindContracts(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contracts *ContractsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contracts.Contract.ContractsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contracts *ContractsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contracts.Contract.ContractsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contracts *ContractsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contracts.Contract.ContractsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contracts *ContractsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contracts.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contracts *ContractsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contracts.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contracts *ContractsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contracts.Contract.contract.Transact(opts, method, params...)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_Contracts *ContractsCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contracts.contract.Call(opts, &out, "counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_Contracts *ContractsSession) Counter() (*big.Int, error) {
	return _Contracts.Contract.Counter(&_Contracts.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() view returns(uint256)
func (_Contracts *ContractsCallerSession) Counter() (*big.Int, error) {
	return _Contracts.Contract.Counter(&_Contracts.CallOpts)
}

// TestFunc1 is a paid mutator transaction binding the contract method 0x88655d98.
//
// Solidity: function testFunc1(string arg1, uint256 arg2, bytes arg3) returns()
func (_Contracts *ContractsTransactor) TestFunc1(opts *bind.TransactOpts, arg1 string, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _Contracts.contract.Transact(opts, "testFunc1", arg1, arg2, arg3)
}

// TestFunc1 is a paid mutator transaction binding the contract method 0x88655d98.
//
// Solidity: function testFunc1(string arg1, uint256 arg2, bytes arg3) returns()
func (_Contracts *ContractsSession) TestFunc1(arg1 string, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _Contracts.Contract.TestFunc1(&_Contracts.TransactOpts, arg1, arg2, arg3)
}

// TestFunc1 is a paid mutator transaction binding the contract method 0x88655d98.
//
// Solidity: function testFunc1(string arg1, uint256 arg2, bytes arg3) returns()
func (_Contracts *ContractsTransactorSession) TestFunc1(arg1 string, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _Contracts.Contract.TestFunc1(&_Contracts.TransactOpts, arg1, arg2, arg3)
}

// TestReverted is a paid mutator transaction binding the contract method 0x6c6dd603.
//
// Solidity: function testReverted() returns()
func (_Contracts *ContractsTransactor) TestReverted(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contracts.contract.Transact(opts, "testReverted")
}

// TestReverted is a paid mutator transaction binding the contract method 0x6c6dd603.
//
// Solidity: function testReverted() returns()
func (_Contracts *ContractsSession) TestReverted() (*types.Transaction, error) {
	return _Contracts.Contract.TestReverted(&_Contracts.TransactOpts)
}

// TestReverted is a paid mutator transaction binding the contract method 0x6c6dd603.
//
// Solidity: function testReverted() returns()
func (_Contracts *ContractsTransactorSession) TestReverted() (*types.Transaction, error) {
	return _Contracts.Contract.TestReverted(&_Contracts.TransactOpts)
}

// ContractsCounterUpdatedIterator is returned from FilterCounterUpdated and is used to iterate over the raw logs and unpacked data for CounterUpdated events raised by the Contracts contract.
type ContractsCounterUpdatedIterator struct {
	Event *ContractsCounterUpdated // Event containing the contract specifics and raw log

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
func (it *ContractsCounterUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractsCounterUpdated)
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
		it.Event = new(ContractsCounterUpdated)
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
func (it *ContractsCounterUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractsCounterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractsCounterUpdated represents a CounterUpdated event raised by the Contracts contract.
type ContractsCounterUpdated struct {
	Counter *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterCounterUpdated is a free log retrieval operation binding the contract event 0x4785d80d2593e2cb7a3331d31eb5106408bdde2aab0db9e9b616b036a1b6039d.
//
// Solidity: event CounterUpdated(uint256 counter)
func (_Contracts *ContractsFilterer) FilterCounterUpdated(opts *bind.FilterOpts) (*ContractsCounterUpdatedIterator, error) {

	logs, sub, err := _Contracts.contract.FilterLogs(opts, "CounterUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractsCounterUpdatedIterator{contract: _Contracts.contract, event: "CounterUpdated", logs: logs, sub: sub}, nil
}

// WatchCounterUpdated is a free log subscription operation binding the contract event 0x4785d80d2593e2cb7a3331d31eb5106408bdde2aab0db9e9b616b036a1b6039d.
//
// Solidity: event CounterUpdated(uint256 counter)
func (_Contracts *ContractsFilterer) WatchCounterUpdated(opts *bind.WatchOpts, sink chan<- *ContractsCounterUpdated) (event.Subscription, error) {

	logs, sub, err := _Contracts.contract.WatchLogs(opts, "CounterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractsCounterUpdated)
				if err := _Contracts.contract.UnpackLog(event, "CounterUpdated", log); err != nil {
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

// ParseCounterUpdated is a log parse operation binding the contract event 0x4785d80d2593e2cb7a3331d31eb5106408bdde2aab0db9e9b616b036a1b6039d.
//
// Solidity: event CounterUpdated(uint256 counter)
func (_Contracts *ContractsFilterer) ParseCounterUpdated(log types.Log) (*ContractsCounterUpdated, error) {
	event := new(ContractsCounterUpdated)
	if err := _Contracts.contract.UnpackLog(event, "CounterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractsFuncEvent1Iterator is returned from FilterFuncEvent1 and is used to iterate over the raw logs and unpacked data for FuncEvent1 events raised by the Contracts contract.
type ContractsFuncEvent1Iterator struct {
	Event *ContractsFuncEvent1 // Event containing the contract specifics and raw log

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
func (it *ContractsFuncEvent1Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractsFuncEvent1)
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
		it.Event = new(ContractsFuncEvent1)
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
func (it *ContractsFuncEvent1Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractsFuncEvent1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractsFuncEvent1 represents a FuncEvent1 event raised by the Contracts contract.
type ContractsFuncEvent1 struct {
	Arg1 string
	Arg2 *big.Int
	Arg3 []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterFuncEvent1 is a free log retrieval operation binding the contract event 0xee7ebd5ac9177b3cfe282c440d0220335dc60bc4472338132f06af7b4b9432fc.
//
// Solidity: event FuncEvent1(string arg1, uint256 arg2, bytes arg3)
func (_Contracts *ContractsFilterer) FilterFuncEvent1(opts *bind.FilterOpts) (*ContractsFuncEvent1Iterator, error) {

	logs, sub, err := _Contracts.contract.FilterLogs(opts, "FuncEvent1")
	if err != nil {
		return nil, err
	}
	return &ContractsFuncEvent1Iterator{contract: _Contracts.contract, event: "FuncEvent1", logs: logs, sub: sub}, nil
}

// WatchFuncEvent1 is a free log subscription operation binding the contract event 0xee7ebd5ac9177b3cfe282c440d0220335dc60bc4472338132f06af7b4b9432fc.
//
// Solidity: event FuncEvent1(string arg1, uint256 arg2, bytes arg3)
func (_Contracts *ContractsFilterer) WatchFuncEvent1(opts *bind.WatchOpts, sink chan<- *ContractsFuncEvent1) (event.Subscription, error) {

	logs, sub, err := _Contracts.contract.WatchLogs(opts, "FuncEvent1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractsFuncEvent1)
				if err := _Contracts.contract.UnpackLog(event, "FuncEvent1", log); err != nil {
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

// ParseFuncEvent1 is a log parse operation binding the contract event 0xee7ebd5ac9177b3cfe282c440d0220335dc60bc4472338132f06af7b4b9432fc.
//
// Solidity: event FuncEvent1(string arg1, uint256 arg2, bytes arg3)
func (_Contracts *ContractsFilterer) ParseFuncEvent1(log types.Log) (*ContractsFuncEvent1, error) {
	event := new(ContractsFuncEvent1)
	if err := _Contracts.contract.UnpackLog(event, "FuncEvent1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
