import { TransactionReceipt, Transaction as Tx } from "web3-core";
import { Block } from "web3-eth";
import { Message } from "element-ui";

export { TransactionReceipt, Tx, Block }

const Web3 = require("web3");

export interface Transaction {
  from: string;
  to?: string;
  data?: string;
  value?: string;
}

export function hexToNumber(hex: string): number {
  const web3 = new Web3("");
  return web3.utils.hexToNumber(hex);
}

export function encodeParameters(types: any[], args: any[]): string {
  const web3 = new Web3("");
  return web3.eth.abi.encodeParameters(types, args).slice(2);
}

export function getBlock(chainAddr: string, blockHash: string): void {
  const web3 = new Web3(chainAddr);
  web3.eth.getBlock(blockHash)
    .then((block: Block) => {
      Message.success("获取区块信息成功");
      console.log("> get block info success: ", block);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get block info failed, error: ", err);
    });
}

export function getTx(chainAddr: string, txHash: string): void {
  const web3 = new Web3(chainAddr);
  web3.eth.getTransaction(txHash)
    .then((tx: Tx) => {
      Message.success("获取交易成功");
      console.log("> get tx success: ", tx);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get tx failed, error: ", err);
    });
}

export function getTxReceipt(chainAddr: string, txHash: string, cb: (receipt: TransactionReceipt) => void): void {
  const web3 = new Web3(chainAddr);
  web3.eth.getTransactionReceipt(txHash)
    .then((receipt: TransactionReceipt) => {
      if (!receipt) {
        Message.info("交易尚未打包上链，请稍等片刻后重试");
        return;
      }

      console.log("> get tx receipt success, receipt: ", receipt);

      cb(receipt);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get tx receipt failed, error: ", err);
    });
}
