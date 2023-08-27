// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

contract Data {
  mapping(string => bool) exist;
  mapping(string => uint256) number;

  function setData(string[] memory strs, uint256[] memory nums) external {
    require(strs.length == nums.length, "not-matched data amount");

    for (uint8 i = 0; i < strs.length; i++) {
      string memory s = strs[i];

      require(!exist[s], "str already exist");

      exist[s] = true;
      number[s] = nums[i];
    }
  }

  function setSingleData(string memory str, uint256 num) external {
    require(!exist[str], "str already exist");

    exist[str] = true;
    number[str] = num;
  }
}
