// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./multi_admin.sol";

library multiSigAddress {
    struct data {
        address[] optionalAddress;
        bool[3][] optionalAddressConfirmed;
    }

    // clearData clear all optional data, used when init / do replace / 'admin' replaced
    function clearData(data storage d) public {
        d.optionalAddress = new address[](0);
        d.optionalAddressConfirmed = new bool[3][](0);
    }

    // replaceAddress returns if should do replace
    // check param 'new addr' before invoke is suggested
    function replaceAddress(multiAdmin.data storage admin, data storage d, address newAddr) external returns (bool) {
        uint senderIndex;
        bool isExist;
        (senderIndex, isExist) = multiAdmin.isAdminExist(admin.admin, msg.sender);
        require(isExist, "permission denied: need admin");

        if (d.optionalAddress.length == 0) {
            clearData(d);
        }

        uint optionalDataIndex;
        (optionalDataIndex, isExist) = isAddressExist(d.optionalAddress, newAddr);

        // a new replace request
        if (!isExist) {
            d.optionalAddress.push(newAddr);

            bool[3] memory newConfirmedData;
            newConfirmedData[senderIndex] = true;
            d.optionalAddressConfirmed.push(newConfirmedData);

            return false;
        }

        bool[3] memory confirmedData = d.optionalAddressConfirmed[optionalDataIndex];

        // repeated confirm
        if (confirmedData[senderIndex]) {
            return false;
        }

        confirmedData[senderIndex] = true;
        d.optionalAddressConfirmed[optionalDataIndex] = confirmedData;

        uint confirmedAmount;
        for (uint i = 0; i < 3; i++) {// '3' admin
            if (confirmedData[i]) {
                confirmedAmount += 1;
            }
        }

        if (confirmedAmount < 2) {
            return false;
        }

        clearData(d);

        return true;
    }

    function isAddressExist(address[] memory arr, address item) pure internal returns (uint index, bool isExist) {
        for (uint i = 0; i < arr.length; i++) {
            if (arr[i] == item) {
                index = i;
                isExist = true;
                break;
            }
        }

        return (index, isExist);
    }
}
