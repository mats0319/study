pragma solidity 0.8.6;

import "./support.sol";

contract Data {
    Support c;

    mapping(string => bool) exist;
    mapping(string => uint256) number;

    constructor(address _support) {
        c = Support(_support);
    }

    function callCrossContract() external {
        if (!c.getExist("a")) {
            c.setExist("a", true);
        }

        if (c.getNumber("a") == 0) {
            c.setNumber("a", 10);
        }
    }

    function callInContract() external {
        if (!getExist("b")) {
            setExist("b", true);
        }

        if (getNumber("b") == 0) {
            setNumber("b", 10);
        }
    }

    function setExist(string memory str, bool isExist) public {
        exist[str] = isExist;
    }

    function getExist(string memory str) view public returns (bool) {
        return exist[str];
    }

    function setNumber(string memory str, uint256 _desp) public {
        number[str] = _desp;
    }

    function getNumber(string memory str) view public returns (uint256) {
        return number[str];
    }
}
