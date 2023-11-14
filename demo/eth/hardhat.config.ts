// register plugins, can import from "hardhat" directly in ts code
import "@nomicfoundation/hardhat-ethers"
import "@nomiclabs/hardhat-web3"

import { HardhatUserConfig } from "hardhat/config"

const config: HardhatUserConfig = {
  solidity: "0.8.6",
  mocha: {
    timeout: 600_000,
  },
};

export default config;
