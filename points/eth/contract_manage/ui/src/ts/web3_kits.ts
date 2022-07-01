// https://ropsten.infura.io/v3/6bdd841c329c4b3884a737067eb8c8ff
// http://192.168.2.57:8545
import { TransactionReceipt } from "web3-core";

export const ethChainAddr = "http://192.168.2.57:8545";

const Web3 = require("web3");
const web3 = new Web3(ethChainAddr);

export interface Transaction {
  from: string;
  to?: string;
  data?: string;
  value?: string;
}

export function encodeParameters(types: any[], args: any[]): string {
  return web3.eth.abi.encodeParameters(types, args).slice(2);
}

export async function getTxReceipt(chainAddr: string, txHash: string): Promise<TransactionReceipt> {
  const web3 = new Web3(chainAddr);
  return await web3.eth.getTransactionReceipt(txHash);
}
