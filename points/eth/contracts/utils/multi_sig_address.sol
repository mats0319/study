// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./utils.sol";

// MultiSigAddress is an extend of 'multi-signature admin' lib,
// this lib implements 'replace' method on 'address' type under multi-sig rules
library MSAddress {
    struct Data {
        address addr;
        bool initFlag;

        address[] replaceData;
        bool[3][] replaceDataConfirmedInfo;
    }

    // init initialize default 'addr', can only invoke once
    function init(Data storage d, address addr) external {
        if (!d.initFlag) {
            d.addr = addr;
            d.initFlag = true;
        }
    }

    // replace return if this invoke done replace and which 'admin' confirmed it(when done replace)
    function replace(Data storage d, address[3] memory admin, address newAddr) external returns (address[] memory, bool) {
        require(newAddr != d.addr, "invalid new address"); // forbid self-replace

        uint senderIndex;
        bool isExist;
        (senderIndex, isExist) = Utils.getAdminIndex(admin, msg.sender);
        require(isExist, "permission denied: need admin");

        // if replace data exist
        uint replaceDataIndex;
        (replaceDataIndex, isExist) = Utils.getAddressIndex(d.replaceData, newAddr);

        // replace data not exist, record it
        if (!isExist) {
            d.replaceData.push(newAddr);

            bool[3] memory newConfirmedInfo;
            newConfirmedInfo[senderIndex] = true;
            d.replaceDataConfirmedInfo.push(newConfirmedInfo);

            return (new address[](0), false);
        }

        // replace data exist
        bool[3] memory confirmedInfo = d.replaceDataConfirmedInfo[replaceDataIndex];
        if (confirmedInfo[senderIndex]) {// repeated confirm
            return (new address[](0), false);
        }

        confirmedInfo[senderIndex] = true;
        d.replaceDataConfirmedInfo[replaceDataIndex] = confirmedInfo;

        // count confirmed admin and prepare to do replace
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

        // do replace and clear replace data
        d.addr = newAddr;
        clearReplaceData(d);

        return (confirmedAdmin, true);
    }

    // clearData clear all optional data, used when init / do replace / 'admin' replaced
    function clearReplaceData(Data storage d) public {
        d.replaceData = new address[](0);
        d.replaceDataConfirmedInfo = new bool[3][](0);
    }
}
