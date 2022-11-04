// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/Strings.sol";

// VerifySignature
contract VerifySignature {
    // usedSignature record hash of used signature
    // in design, format of 'message' is '[behavior] [token id] - [timestamp]',
    // e.g. 'Open NFT (token id: [xxx]) - 123456000'
    // which means same 'message' normally incorrect
    mapping(bytes32 => bool) usedSignature;

    function calcHash(bytes memory data) pure internal returns (bytes32) {
        return ECDSA.toEthSignedMessageHash(data);
    }

    function verify(
        uint behavior,
        uint256 tokenId,
        string memory timestamp,
        bytes memory signature,
        address signer,
        address contractAddress
    ) external {
        bytes32 signatureHash = calcHash(signature);
        require(!usedSignature[signatureHash], "duplicated signature");

        bytes memory message = generateSignMessage(signBehaviorToStr(behavior), tokenId, timestamp, contractAddress);
        bytes32 msgHash = calcHash(message);
        require(ECDSA.recover(msgHash, signature) == signer, "invalid signature");

        usedSignature[signatureHash] = true;
    }

    function signBehaviorToStr(uint behavior) pure internal returns (string memory) {
        string memory res;
        if (behavior == 0) {
            res = "Open NFT";
        } else if (behavior == 1) {
            res = "Recharge BCash";
        } else if (behavior == 2) {
            res = "Auction Lock";
        } else {
            require(false, "unknown behavior");
        }

        return res;
    }

    // generateOpenNFTSignMessage return hex str
    function generateSignMessage(
        string memory behavior,
        uint256 tokenId,
        string memory timestamp,
        address contractAddress
    ) pure internal returns (bytes memory) {
        bytes memory params = abi.encodePacked(behavior, "(token id: ", Strings.toHexString(tokenId), ") in ", timestamp, " on ");
        bytes memory addressBytes = hexContractAddress(contractAddress);

        bytes memory message = new bytes(params.length + addressBytes.length);
        uint index;
        for (uint i = 0; i < params.length; i++) {
            message[index] = params[i];
            index++;
        }

        for (uint i = 0; i < addressBytes.length; i++) {
            message[index] = addressBytes[i];
            index++;
        }

        return message;
    }

    function hexContractAddress(address contractAddress) pure internal returns (bytes memory) {
        bytes memory res = new bytes(40);
        bytes memory alphabet = "0123456789abcdef";
        bytes20 data = bytes20(contractAddress);
        for (uint i = 0; i < 20; i++) {
            uint8 char = uint8(data[i]);
            res[i * 2 + 0] = alphabet[uint256(char) >> 4];
            res[i * 2 + 1] = alphabet[uint256(char) & 15];
        }

        return res;
    }
}