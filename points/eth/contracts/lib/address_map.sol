// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

library addressMap {
    struct map {
        mapping(address => bool) data;
        address[] keys;
    }

    function add(map storage m, address key) external {
        if (m.data[key]) {
            return;
        }

        m.data[key] = true;
        m.keys.push(key);
    }

    function del(map storage m, address key) external {
        if (!m.data[key]) {
            return;
        }

        m.data[key] = false;

        for (uint8 i = 0; i < m.keys.length; i++) {
            if (m.keys[i] == key) {
                m.keys[i] = m.keys[m.keys.length-1];
                m.keys.pop();
            }
        }
    }

    function getKeys(map storage m) view external returns (address[] memory) {
        return m.keys;
    }

    function getAmount(map storage m) view external returns (uint) {
        return m.keys.length;
    }

    function isExist(map storage m, address key) view external returns (bool) {
        return m.data[key];
    }
}
