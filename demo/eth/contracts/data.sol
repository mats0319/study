// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

contract Data {
    uint256 private _number = 10;
    uint256 private _numberSave = 10;

    // when use exist contract, call this method to test if instance is valid
    function getNumber() view external returns (uint256) {
        return _number;
    }

    // Q1: event is not free
    event Event1(address caller);
    event Event2(address caller, uint256 number);

    function funcWithEvent1() external {
        _number = 20;
        _number = 10;
        emit Event1(msg.sender);
    }

    function funcWithEvent2() external {
        _number = 20;
        _number = 10;
        emit Event2(msg.sender, _number);
    }

    function funcWithoutEvent() external {
        _number = 20;
        _number = 10;
    }
}
