import { contractABI } from "@/views/eth/data";
import { Transaction } from "@/ts/web3_kits";

const Web3 = require("web3");
const web3 = new Web3("");

const contractIns = new web3.eth.Contract(contractABI, "");

export function buildTx_ReplaceETHReceiver(from: string, to: string, newETHReceiver: string): Transaction {
  newETHReceiver = web3.utils.toChecksumAddress(newETHReceiver);

  return {
    from: from,
    to: to,
    data: contractIns.methods.replaceETHReceiver(newETHReceiver).encodeABI(),
  }
}

export function buildTx_AddACashOperators(from: string, to: string, operators: Array<string>): Transaction {
  for (let i = 0; i < operators.length; i++) {
    operators[i] = web3.utils.toChecksumAddress(operators[i]);
  }

  return {
    from: from,
    to: to,
    data: contractIns.methods.addACashOperators(operators).encodeABI(),
  }
}

export function buildTx_DelACashOperators(from: string, to: string, operators: Array<string>): Transaction {
  for (let i = 0; i < operators.length; i++) {
    operators[i] = web3.utils.toChecksumAddress(operators[i]);
  }

  return {
    from: from,
    to: to,
    data: contractIns.methods.delACashOperators(operators).encodeABI(),
  }
}

export function buildTx_AddBrokerageOperators(from: string, to: string, operators: Array<string>): Transaction {
  for (let i = 0; i < operators.length; i++) {
    operators[i] = web3.utils.toChecksumAddress(operators[i]);
  }

  return {
    from: from,
    to: to,
    data: contractIns.methods.addBrokerageOperators(operators).encodeABI(),
  }
}

export function buildTx_DelBrokerageOperators(from: string, to: string, operators: Array<string>): Transaction {
  for (let i = 0; i < operators.length; i++) {
    operators[i] = web3.utils.toChecksumAddress(operators[i]);
  }

  return {
    from: from,
    to: to,
    data: contractIns.methods.delBrokerageOperators(operators).encodeABI(),
  }
}

export function buildTx_ReplaceAdmin(from: string, to: string, oldAdmin: string, newAdmin: string): Transaction {
  oldAdmin = web3.utils.toChecksumAddress(oldAdmin);
  newAdmin = web3.utils.toChecksumAddress(newAdmin);

  return {
    from: from,
    to: to,
    data: contractIns.methods.replaceAdmin(oldAdmin, newAdmin).encodeABI(),
  }
}
