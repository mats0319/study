pragma solidity 0.8.6;

contract Support {
    mapping(string => bool) exist;
    mapping(string => uint256) number;

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
