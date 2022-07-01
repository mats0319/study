// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

// multiAdmin implement self-define multi-sig rule on chain
// our multi-sig rule:
// 1. there are 3 'admin' and 1 'decision maker'
// 2. 'admin' hold the highest permission of contract
//    1. if 'admin' want set something, need >=2 'admin' confirmed(send tx)
// 3. 'decision maker' only used to replace 'admin'
//    1. replace 'admin', need 'decision maker' and >=2 'admin' confirmed(send tx)
library multiAdmin {
    // multi-sig rule: 3 admin, 2 confirmed
    struct data {
        address[3] admin;
        address decisionMaker; // can not modify, only use to modify 'admin'

        // keep all optional replace data and confirmed info,
        // clear all after do replace
        address[2][] replaceAdmin; // 0: old admin, 1: new admin
        bool[4][] replaceAdminConfirmed; // last item means if 'decision maker' confirmed
    }

    // replaceAdmin replace admin when 'decision maker' and 'work amount' amount of 'admins' confirmed
    // return if do replace
    function replaceAdmin(data storage d, address oldAdmin, address newAdmin) external returns (bool) {
        // params check
        uint senderIndex;
        uint oldAdminIndex;
        bool isExist;
        (senderIndex, isExist) = isAdminExist(d.admin, msg.sender);
        require(isExist || msg.sender == d.decisionMaker, "permission denied: invalid caller");
        // limit 'old admin' and 'new admin' are different incidentally
        (oldAdminIndex, isExist) = isAdminExist(d.admin, oldAdmin);
        require(isExist, "invalid old admin");
        (, isExist) = isAdminExist(d.admin, newAdmin);
        require(!isExist, "invalid new admin");

        if (msg.sender == d.decisionMaker) {
            senderIndex = 3;
        }

        if (d.replaceAdmin.length == 0) {
            d.replaceAdmin = new address[2][](0);
            d.replaceAdminConfirmed = new bool[4][](0);
        }

        uint replaceDataIndex;
        (replaceDataIndex, isExist) = isReplaceDataExist(d.replaceAdmin, oldAdmin, newAdmin);

        // a new replace request
        if (!isExist) {
            address[2] memory newReplaceData = [oldAdmin, newAdmin];
            d.replaceAdmin.push(newReplaceData);

            bool[4] memory newReplaceConfirmData;
            newReplaceConfirmData[senderIndex] = true;
            d.replaceAdminConfirmed.push(newReplaceConfirmData);

            return false;
        }

        // replace request is exist
        bool[4] memory replaceConfirmData = d.replaceAdminConfirmed[replaceDataIndex];
        if (replaceConfirmData[senderIndex]) { // repeated confirm, return directly
            return false;
        }

        replaceConfirmData[senderIndex] = true;
        d.replaceAdminConfirmed[replaceDataIndex] = replaceConfirmData;

        // 'decision maker' not confirmed
        if (!replaceConfirmData[3]) {
            return false;
        }

        // 'decision maker' confirmed, count confirmed admin and prepare to do replace
        uint confirmedAmount;
        for (uint i = 0; i < 3; i++) { // '3' admin
            if (replaceConfirmData[i]) {
                confirmedAmount += 1;
            }
        }

        if (confirmedAmount < 2) { // '2' confirmed
            return false;
        }

        // do replace and clear all optional data
        d.admin[oldAdminIndex] = newAdmin;
        d.replaceAdmin = new address[2][](0);
        d.replaceAdminConfirmed = new bool[4][](0);

        return true;
    }

    function isReplaceDataExist(address[2][] memory arr, address oldAdmin, address newAdmin) pure internal returns (uint index, bool isExist) {
        for (uint i = 0; i < arr.length; i++) {
            if (arr[i][0] == oldAdmin && arr[i][1] == newAdmin) {
                index = i;
                isExist = true;
                break;
            }
        }

        return (index, isExist);
    }

    function isAdminExist(address[3] memory arr, address item) pure public returns (uint index, bool isExist) {
        for (uint i = 0; i < 3; i++) {
            if (arr[i] == item) {
                index = i;
                isExist = true;
                break;
            }
        }

        return (index, isExist);
    }
}
