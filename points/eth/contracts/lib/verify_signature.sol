// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

library verifySignature {
    using Strings for uint256;

    enum SignBehavior {
        Init,
        OpenNFT,
        RechargeBCash,
        AuctionLock
    }

    function calcHash(bytes memory data) pure external returns (bytes32) {
        return ECDSA.toEthSignedMessageHash(data);
    }

    function verify(
        SignBehavior behavior,
        uint256 tokenId,
        string memory timestamp,
        bytes memory signature,
        address signer
    ) view external returns (bool) {
        string memory behaviorStr;
        if (behavior == SignBehavior.OpenNFT) {
            behaviorStr = "Open NFT";
        } else if (behavior == SignBehavior.RechargeBCash) {
            behaviorStr = "Recharge BCash";
        } else if (behavior == SignBehavior.AuctionLock) {
            behaviorStr = "Auction Lock";
        } else {
            return false;
        }

        bytes memory message = generateSignMessage(behaviorStr, tokenId, timestamp);

        bytes32 hash = ECDSA.toEthSignedMessageHash(message);

        return ECDSA.recover(hash, signature) == signer;
    }

    // generateOpenNFTSignMessage return hex str
    function generateSignMessage(
        string memory behavior,
        uint256 tokenId,
        string memory timestamp
    ) view internal returns (bytes memory) {
        bytes memory params = abi.encodePacked(behavior, "(token id: ", tokenId.toHexString(), ") in ", timestamp, " on ");
        bytes memory addressBytes = hexContractAddress();

        bytes memory message = new bytes(params.length+addressBytes.length);
        uint k;
        for (uint i = 0; i < params.length; i++) {
            message[k] = params[i];
            k++;
        }

        for (uint i = 0; i < addressBytes.length; i++) {
            message[k] = addressBytes[i];
            k++;
        }

        return message;
    }

    function hexContractAddress() view internal returns (bytes memory) {
        bytes memory res = new bytes(40);
        bytes memory alphabet = "0123456789abcdef"; // 'res' generate according to this, which means letters in 'res' always in small-case
        bytes20 data = bytes20(address(this));
        for (uint i = 0; i < 20; i++) {
            uint8 char = uint8(data[i]);
//            if (41 <= char && char <= 46) {
//                char += 20;
//            }

            res[i*2 + 0] = alphabet[uint256(char) >> 4];
            res[i*2 + 1] = alphabet[uint256(char) & 15];
        }

        return res;
    }
}