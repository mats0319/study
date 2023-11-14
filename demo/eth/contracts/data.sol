// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

contract Data {
    uint256 private _number = 0;

    // when use exist contract, call this method to test if instance is valid
    function getNumber() view external returns (uint256) {
        return _number;
    }

    // Q1: event is not free
    event Event1(address caller);
    event Event2(address caller, uint256 number);

    function funcWithEvent1() external {
        _number = 0;
        _number = 100;
        emit Event1(msg.sender);
    }

    function funcWithEvent2() external {
        _number = 0;
        _number = 100;
        emit Event2(msg.sender, _number);
    }

    function funcWithoutEvent() external {
        _number = 0;
        _number = 100;
    }

    // Q2: cost gas(calculate and storage)
    function calcMul2() external {
        _number = 0;
        _number = _number * 2;
    }

    function calcMul2Double() external {
        _number = 0;
        _number = _number * 2 * 2;
    }

    function storageNothing() external returns (uint256) {
        _number = 0;
        uint256 v = _number * 2;
        return v;
    }

    function storageVariable() external returns (uint256) {
        _number = 0;
        _number = _number * 2;
        return _number;
    }
}
