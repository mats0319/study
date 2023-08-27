// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

library Utils {
    function getAdminIndex(address[3] memory array, address item) pure internal returns (uint index, bool isExist) {
        for (uint i = 0; i < 3; i++) {
            if (array[i] == item) {
                index = i;
                isExist = true;
                break;
            }
        }

        return (index, isExist);
    }

    function getAddressIndex(address[] memory array, address item) pure internal returns (uint index, bool isExist) {
        for (uint i = 0; i < array.length; i++) {
            if (array[i] == item) {
                index = i;
                isExist = true;
                break;
            }
        }

        return (index, isExist);
    }
}
