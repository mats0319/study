import { contractABI, contractAddress } from "@/views/eth/data";
import { sendTransaction } from "@/ts/wallet_connect";
import { Transaction } from "@/ts/web3_kits";

const Web3 = require("web3");
const web3Default = new Web3("");

const contractInsDefault = new web3Default.eth.Contract(contractABI, contractAddress);

export async function replaceETHReceiver(from: string, to: string, input: string): Promise<any> {
  const replacedAddr = web3Default.utils.toChecksumAddress(input);

  const tx: Transaction = {
    from: from,
    to: to,
    data: contractInsDefault.methods.replaceETHReceiver(replacedAddr).encodeABI(),
  }

  return sendTransaction(tx);
}

export async function getETHReceiver(chainAddr: string, contractAddr: string): Promise<any> {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);
  return await contractIns.methods.getETHReceiver().call();
}

export function addACashOperators(from: string, to: string, input: Array<string>): Promise<any> {
  const operators: Array<string> = [];
  for (let i = 0; i < input.length; i++) {
    operators.push(web3Default.utils.toChecksumAddress(input[i]));
  }

  const tx: Transaction = {
    from: from,
    to: to,
    data: contractInsDefault.methods.addACashOperators(operators).encodeABI(),
  }

  return sendTransaction(tx);
}

export function delACashOperators(from: string, to: string, input: Array<string>): Promise<any> {
  const operators: Array<string> = [];
  for (let i = 0; i < input.length; i++) {
    operators.push(web3Default.utils.toChecksumAddress(input[i]));
  }

  const tx: Transaction = {
    from: from,
    to: to,
    data: contractInsDefault.methods.delACashOperators(operators).encodeABI(),
  }

  return sendTransaction(tx);
}

export async function getACashOperators(chainAddr: string, contractAddr: string): Promise<any> {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);
  return await contractIns.methods.getACashOperators().call();
}

export function addBrokerageOperators(from: string, to: string, input: Array<string>): Promise<any> {
  const operators: Array<string> = [];
  for (let i = 0; i < input.length; i++) {
    operators.push(web3Default.utils.toChecksumAddress(input[i]));
  }

  const tx: Transaction = {
    from: from,
    to: to,
    data: contractInsDefault.methods.addBrokerageOperators(operators).encodeABI(),
  }

  return sendTransaction(tx);
}

export function delBrokerageOperators(from: string, to: string, input: Array<string>): Promise<any> {
  const operators: Array<string> = [];
  for (let i = 0; i < input.length; i++) {
    operators.push(web3Default.utils.toChecksumAddress(input[i]));
  }

  const tx: Transaction = {
    from: from,
    to: to,
    data: contractInsDefault.methods.delBrokerageOperators(operators).encodeABI(),
  }

  return sendTransaction(tx);
}

export async function getBrokerageOperators(chainAddr: string, contractAddr: string): Promise<any> {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);
  return await contractIns.methods.getBrokerageOperators().call();
}

export function replaceAdmin(from: string, to: string, input1: string, input2: string): Promise<any> {
  const oldAdmin = web3Default.utils.toChecksumAddress(input1);
  const newAdmin = web3Default.utils.toChecksumAddress(input2);

  const tx: Transaction = {
    from: from,
    to: to,
    data: contractInsDefault.methods.replaceAdmin(oldAdmin, newAdmin).encodeABI(),
  }

  return sendTransaction(tx);
}
