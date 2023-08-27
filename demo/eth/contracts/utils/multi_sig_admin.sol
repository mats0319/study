// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./utils.sol";

// Admin implement self-define multi-signature rule on chain
// our multi-sig rule:
// 1. there are 3 'admin' and 1 'decision maker'
// 2. 'admin' hold the highest permission of contract
//    1. if 'admin' want set something, need >=2 'admin' confirmed(send tx)
// 3. 'decision maker' only used to replace 'admin'
//    1. replace 'admin', need 'decision maker' and >=2 'admin' confirmed(send tx)
// 4. all optional replace data will be saved, if one of them confirmed, clear others
contract MSAdmin {
    event ReplaceAdmin(address oldAdmin, address newAdmin, address[] confirmedAdmin);

    // multi-signature rule: 3 admin, 2 confirmed
    struct Data {
        address[3] admin;
        address decisionMaker; // can not modify, only use to modify 'admin'
        bool initFlag; // if data is initialized

        // keep all optional replace data and confirmed info
        address[2][] replaceData; // 0: old admin, 1: new admin
        bool[4][] replaceDataConfirmedInfo; // last item means if 'decision maker' confirmed
    }

    Data _admin;

    function getAdmin() view external returns (address[3] memory, address) {
        return (_admin.admin, _admin.decisionMaker);
    }

    // init initialize default 'admin' and 'decision maker', can only invoke once
    // init func defines for situation: skip multi-signature rule, set data without twice invoke, e.g. constructor()
    function init(address[3] memory admin, address decisionMaker) internal {
        if (_admin.initFlag) {
            return;
        }

        _admin.admin = admin;
        _admin.decisionMaker = decisionMaker;
        clearReplaceData();

        _admin.initFlag = true;
    }

    // replace replace admin when 'decision maker' and >=2 'admin' confirmed
    // return if this invoke done replace and which 'admin' confirmed it(when done replace)
    // @param oldAdmin: address you want to del from 'admin', require exist in 'admin'
    // @param newAdmin: address you want to add to 'admin', require not exist in 'admin'
    function replace(address oldAdmin, address newAdmin) internal returns (address[] memory, bool) {
        uint senderIndex;
        uint oldAdminIndex;
        bool isExist;

        if (msg.sender == _admin.decisionMaker) {// 'decision maker' confirm
            senderIndex = 3;
        } else {// 'admin' confirm or invalid caller
            bool isValidCaller;
            (senderIndex, isValidCaller) = Utils.getAdminIndex(_admin.admin, msg.sender);
            require(isValidCaller, "permission denied: invalid caller");
        }

        (oldAdminIndex, isExist) = Utils.getAdminIndex(_admin.admin, oldAdmin);
        require(isExist, "invalid old admin");

        (, isExist) = Utils.getAdminIndex(_admin.admin, newAdmin);
        require(!isExist, "invalid new admin");

        // if replace data exist
        uint replaceDataIndex;
        (replaceDataIndex, isExist) = getReplaceDataIndex(_admin.replaceData, oldAdmin, newAdmin);

        // replace data not exist, record it
        if (!isExist) {
            address[2] memory newReplaceData = [oldAdmin, newAdmin];
            _admin.replaceData.push(newReplaceData);

            bool[4] memory newConfirmedInfo;
            newConfirmedInfo[senderIndex] = true;
            _admin.replaceDataConfirmedInfo.push(newConfirmedInfo);

            return (new address[](0), false);
        }

        // replace data exist
        bool[4] memory confirmedInfo = _admin.replaceDataConfirmedInfo[replaceDataIndex];
        if (confirmedInfo[senderIndex]) {// repeated confirm
            return (new address[](0), false);
        }

        confirmedInfo[senderIndex] = true;
        _admin.replaceDataConfirmedInfo[replaceDataIndex] = confirmedInfo;

        if (!confirmedInfo[3]) {// 'decision maker' not confirmed
            return (new address[](0), false);
        }

        // 'decision maker' confirmed, count confirmed admin and prepare to do replace
        uint confirmedAmount;
        address[] memory confirmedAdmin = new address[](2);
        for (uint i = 0; i < 3 && confirmedAmount < 2; i++) {// '3' admin
            if (confirmedInfo[i]) {
                confirmedAmount += 1;
                confirmedAdmin[confirmedAmount - 1] = _admin.admin[i];
            }
        }

        if (confirmedAmount < 2) {// '2' confirmed
            return (new address[](0), false);
        }

        // do replace and clear replace data
        _admin.admin[oldAdminIndex] = newAdmin;
        clearReplaceData();

        emit ReplaceAdmin(oldAdmin, newAdmin, confirmedAdmin);

        return (confirmedAdmin, true);
    }

    function clearReplaceData() internal {
        _admin.replaceData = new address[2][](0);
        _admin.replaceDataConfirmedInfo = new bool[4][](0);
    }

    function getReplaceDataIndex(address[2][] memory array, address oldAdmin, address newAdmin) pure internal returns (uint index, bool isExist) {
        for (uint i = 0; i < array.length; i++) {
            if (array[i][0] == oldAdmin && array[i][1] == newAdmin) {
                index = i;
                isExist = true;
                break;
            }
        }

        return (index, isExist);
    }
}
