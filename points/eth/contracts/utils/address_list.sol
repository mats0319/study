// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./utils.sol";

library AddressList {
    struct List {
        mapping(address => bool) isExist;
        address[] array;
    }

    function add(List storage l, address[] memory array) internal {
        for (uint8 i = 0; i < array.length; i++) {
            if (l.isExist[array[i]]) {
                continue;
            }

            l.isExist[array[i]] = true;
            l.array.push(array[i]);
        }
    }

    function del(List storage l, address[] memory array) internal {
        for (uint8 i = 0; i < array.length; i++) {
            if (!l.isExist[array[i]]) {
                continue;
            }

            l.isExist[array[i]] = false;

            uint index;
            bool isExist;
            (index, isExist) = Utils.getAddressIndex(l.array, array[i]);
            if (isExist) {
                l.array[index] = l.array[l.array.length - 1];
                l.array.pop();
            }
        }
    }

    function getList(List storage l) view internal returns (address[] memory) {
        return l.array;
    }

    function contains(List storage l, address addr) view internal returns (bool) {
        return l.isExist[addr];
    }
}
