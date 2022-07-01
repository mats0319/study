// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./multi_admin.sol";

library multiSigAddrArray {
    struct data {
        address[][] optionalAddrArrayForAdd;
        bool[3][] optionalAddrArrayForAddConfirmed;

        address[][] optionalAddrArrayForDel;
        bool[3][] optionalAddrArrayForDelConfirmed;
    }

    // clearAddData clear all optional data for add, used when init / do replace / 'admin' replaced
    function clearAddData(data storage d) public {
        d.optionalAddrArrayForAdd = new address[][](0);
        d.optionalAddrArrayForAddConfirmed = new bool[3][](0);
    }

    // clearDelData clear all optional data for del, used when init / do replace / 'admin' replaced
    function clearDelData(data storage d) public {
        d.optionalAddrArrayForDel = new address[][](0);
        d.optionalAddrArrayForDelConfirmed = new bool[3][](0);
    }

    // addAddrArray returns if should do add
    // different order of 'new addrs' are regard as different data
    function addAddrArray(multiAdmin.data storage admin, data storage d, address[] memory newAddrs) external returns (bool) {
        uint senderIndex;
        bool isExist;
        (senderIndex, isExist) = multiAdmin.isAdminExist(admin.admin, msg.sender);
        require(isExist, "permission denied: need admin");

        if (d.optionalAddrArrayForAdd.length == 0) {
            clearAddData(d);
        }

        uint optionalDataIndex;
        (optionalDataIndex, isExist) = isAddressArrayExist(d.optionalAddrArrayForAdd, newAddrs);

        // a new add request
        if (!isExist) {
            d.optionalAddrArrayForAdd.push(newAddrs);

            bool[3] memory newConfirmedData;
            newConfirmedData[senderIndex] = true;
            d.optionalAddrArrayForAddConfirmed.push(newConfirmedData);

            return false;
        }

        bool[3] memory confirmedData = d.optionalAddrArrayForAddConfirmed[optionalDataIndex];

        // repeated confirm
        if (confirmedData[senderIndex]) {
            return false;
        }

        confirmedData[senderIndex] = true;
        d.optionalAddrArrayForAddConfirmed[optionalDataIndex] = confirmedData;

        uint confirmedAmount;
        for (uint i = 0; i < 3; i++) {// '3' admin
            if (confirmedData[i]) {
                confirmedAmount += 1;
            }
        }

        if (confirmedAmount < 2) {
            return false;
        }

        clearAddData(d);

        return true;
    }

    // delAddrArray returns if should do del
    // different order of 'new addrs' are regard as different data
    function delAddrArray(multiAdmin.data storage admin, data storage d, address[] memory newAddrs) external returns (bool) {
        uint senderIndex;
        bool isExist;
        (senderIndex, isExist) = multiAdmin.isAdminExist(admin.admin, msg.sender);
        require(isExist, "permission denied: need admin");

        if (d.optionalAddrArrayForDel.length == 0) {
            clearDelData(d);
        }

        uint optionalDataIndex;
        (optionalDataIndex, isExist) = isAddressArrayExist(d.optionalAddrArrayForDel, newAddrs);

        // a new add request
        if (!isExist) {
            d.optionalAddrArrayForDel.push(newAddrs);

            bool[3] memory newConfirmedData;
            newConfirmedData[senderIndex] = true;
            d.optionalAddrArrayForDelConfirmed.push(newConfirmedData);

            return false;
        }

        bool[3] memory confirmedData = d.optionalAddrArrayForDelConfirmed[optionalDataIndex];

        // repeated confirm
        if (confirmedData[senderIndex]) {
            return false;
        }

        confirmedData[senderIndex] = true;
        d.optionalAddrArrayForDelConfirmed[optionalDataIndex] = confirmedData;

        uint confirmedAmount;
        for (uint i = 0; i < 3; i++) {// '3' admin
            if (confirmedData[i]) {
                confirmedAmount += 1;
            }
        }

        if (confirmedAmount < 2) {
            return false;
        }

        clearDelData(d);

        return true;
    }

    function isAddressArrayExist(address[][] memory arr, address[] memory item) pure internal returns (uint index, bool isExist) {
        for (uint i = 0; i < arr.length; i++) {
            if (compareOnAddressArray(arr[i], item)) {
                index = i;
                isExist = true;
                break;
            }
        }

        return (index, isExist);
    }

    // compareOnAddressArray return if two address array are strict equal, include order
    function compareOnAddressArray(address[] memory a, address[] memory b) pure internal returns (bool) {
        if (a.length != b.length) {
            return false;
        }

        bool isEqual = true;
        for (uint i = 0; i < a.length; i++) {
            if (a[i] != b[i]) {
                isEqual = false;
                break;
            }
        }

        return isEqual;
    }
}
