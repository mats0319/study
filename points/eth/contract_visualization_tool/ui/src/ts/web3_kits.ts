import { TransactionReceipt, Transaction as Tx } from "web3-core";
import { Block } from "web3-eth";

export { TransactionReceipt, Tx, Block }

const Web3 = require("web3");

export interface Transaction {
  from: string;
  to?: string;
  data?: string;
  value?: string;
}

export function encodeParameters(types: any[], args: any[]): string {
  const web3 = new Web3("");

  return web3.eth.abi.encodeParameters(types, args).slice(2);
}

export async function getBlock(chainAddr: string, blockHash: string): Promise<Block> {
  const web3 = new Web3(chainAddr);
  return await web3.eth.getBlock(blockHash);
}

export async function getTx(chainAddr: string, txHash: string): Promise<Tx> {
  const web3 = new Web3(chainAddr);
  return await web3.eth.getTransaction(txHash);
}

export async function getTxReceipt(chainAddr: string, txHash: string): Promise<TransactionReceipt> {
  const web3 = new Web3(chainAddr);
  return await web3.eth.getTransactionReceipt(txHash);
}
