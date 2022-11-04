// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./address_list.sol";
import "./utils.sol";

// MultiSigAddressArray is an extend of 'multi-signature admin' lib,
// this lib implements 'add'/'del' method on 'address[]' type under multi-sig rules
library MSAddressArray {
    struct Data {
        AddressList.List array;

        address[][] optionalAddData;
        bool[3][] optionalAddDataConfirmedInfo;

        address[][] optionalDelData;
        bool[3][] optionalDelDataConfirmedInfo;
    }

    // add returns if this invoke done add and which 'admin' confirmed it(when done add)
    // note:
    // 1. different order of 'new addrs' are regard as different data, because solidity can not loop keys of map type
    // 2. duplicated data in 'new addrs' will be handle in 'add' step
    function add(Data storage d, address[3] memory admin, address[] memory newAddrs) external returns (address[] memory, bool) {
        uint senderIndex;
        bool isExist;
        (senderIndex, isExist) = Utils.getAdminIndex(admin, msg.sender);
        require(isExist, "permission denied: need admin");

        // if optional data exist
        uint optionalDataIndex;
        (optionalDataIndex, isExist) = isAddressArrayExist(d.optionalAddData, newAddrs);

        // optional data not exist, record it
        if (!isExist) {
            d.optionalAddData.push(newAddrs);

            bool[3] memory newConfirmedInfo;
            newConfirmedInfo[senderIndex] = true;
            d.optionalAddDataConfirmedInfo.push(newConfirmedInfo);

            return (new address[](0), false);
        }

        // optional data exist
        bool[3] memory confirmedInfo = d.optionalAddDataConfirmedInfo[optionalDataIndex];
        if (confirmedInfo[senderIndex]) {// repeated confirm
            return (new address[](0), false);
        }

        confirmedInfo[senderIndex] = true;
        d.optionalAddDataConfirmedInfo[optionalDataIndex] = confirmedInfo;

        // count confirmed admin and prepare to do add
        uint confirmedAmount;
        address[] memory confirmedAdmin = new address[](2);
        for (uint i = 0; i < 3 && confirmedAmount < 2; i++) {// '3' admin
            if (confirmedInfo[i]) {
                confirmedAmount += 1;
                confirmedAdmin[confirmedAmount - 1] = admin[i];
            }
        }

        if (confirmedAmount < 2) {// '2' confirmed
            return (new address[](0), false);
        }

        // do add and clear optional add data
        AddressList.add(d.array, newAddrs);
        clearOptionalAddData(d);

        return (confirmedAdmin, true);
    }

    // del returns if this invoke done del and which 'admin' confirmed it(when done del)
    // note:
    // 1. different order of 'new addrs' are regard as different data, because solidity can not loop keys of map type
    // 2. duplicated data in 'new addrs' will be handle in 'del' step
    function del(Data storage d, address[3] memory admin, address[] memory newAddrs) external returns (address[] memory, bool) {
        uint senderIndex;
        bool isExist;
        (senderIndex, isExist) = Utils.getAdminIndex(admin, msg.sender);
        require(isExist, "permission denied: need admin");

        // if optional data exist
        uint optionalDataIndex;
        (optionalDataIndex, isExist) = isAddressArrayExist(d.optionalDelData, newAddrs);

        // optional data not exist, record it
        if (!isExist) {
            d.optionalDelData.push(newAddrs);

            bool[3] memory newConfirmedInfo;
            newConfirmedInfo[senderIndex] = true;
            d.optionalDelDataConfirmedInfo.push(newConfirmedInfo);

            return (new address[](0), false);
        }

        // optional data exist
        bool[3] memory confirmedInfo = d.optionalDelDataConfirmedInfo[optionalDataIndex];
        if (confirmedInfo[senderIndex]) {// repeated confirm
            return (new address[](0), false);
        }

        confirmedInfo[senderIndex] = true;
        d.optionalDelDataConfirmedInfo[optionalDataIndex] = confirmedInfo;

        // count confirmed admin and prepare to do del
        uint confirmedAmount;
        address[] memory confirmedAdmin = new address[](2);
        for (uint i = 0; i < 3 && confirmedAmount < 2; i++) {// '3' admin
            if (confirmedInfo[i]) {
                confirmedAmount += 1;
                confirmedAdmin[confirmedAmount - 1] = admin[i];
            }
        }

        if (confirmedAmount < 2) {// '2' confirmed
            return (new address[](0), false);
        }

        // do del and clear optional add data
        AddressList.del(d.array, newAddrs);
        clearOptionalDelData(d);

        return (confirmedAdmin, true);
    }

    function getList(Data storage d) view external returns (address[] memory) {
        return AddressList.getList(d.array);
    }

    function contains(Data storage d, address addr) view external returns (bool) {
        return AddressList.contains(d.array, addr);
    }

    function clearOptionalData(Data storage d) external {
        clearOptionalAddData(d);
        clearOptionalDelData(d);
    }

    // clearOptionalAddData clear all optional add data, used when init / do add / 'admin' replaced
    function clearOptionalAddData(Data storage d) internal {
        d.optionalAddData = new address[][](0);
        d.optionalAddDataConfirmedInfo = new bool[3][](0);
    }

    // clearOptionalDelData clear all optional del data, used when init / do del / 'admin' replaced
    function clearOptionalDelData(Data storage d) internal {
        d.optionalDelData = new address[][](0);
        d.optionalDelDataConfirmedInfo = new bool[3][](0);
    }

    function isAddressArrayExist(address[][] memory array, address[] memory item) pure internal returns (uint index, bool isExist) {
        for (uint i = 0; i < array.length; i++) {
            if (compareOnAddressArray(array[i], item)) {
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
