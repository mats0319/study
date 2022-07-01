import { contractABI, contractAddress } from "@/views/cashbox_controller/data";
import { sendTransaction } from "@/ts/wallet_connect";
import { Transaction } from "@/ts/web3_kits";

const Web3 = require("web3");
const web3Default = new Web3("");

const contractInsDefault = new web3Default.eth.Contract(contractABI, contractAddress);

export function setControllerContract(from: string, to: string, input: string): Promise<any> {
  const cAddr = web3Default.utils.toChecksumAddress(input);

  const tx: Transaction = {
    from: from,
    to: to,
    data: contractInsDefault.methods.setControllerContract(cAddr).encodeABI(),
  }

  return sendTransaction(tx);
}

export function addMiners(from: string, to: string, input: Array<string>): Promise<any> {
  const miners: Array<string> = [];
  for (let i = 0; i < input.length; i++) {
    miners.push(web3Default.utils.toChecksumAddress(input[i]));
  }

  const tx: Transaction = {
    from: from,
    to: to,
    data: contractInsDefault.methods.addMiners(miners).encodeABI(),
  }

  return sendTransaction(tx);
}

export function delMiners(from: string, to: string, input: Array<string>): Promise<any> {
  const miners: Array<string> = [];
  for (let i = 0; i < input.length; i++) {
    miners.push(web3Default.utils.toChecksumAddress(input[i]));
  }

  const tx: Transaction = {
    from: from,
    to: to,
    data: contractInsDefault.methods.delMiners(miners).encodeABI(),
  }

  return sendTransaction(tx);
}

export async function getMiners(chainAddr: string, contractAddr: string): Promise<any> {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);
  return await contractIns.methods.getMiners().call();
}

export function addOperators(from: string, to: string, input: Array<string>): Promise<any> {
  const operators: Array<string> = [];
  for (let i = 0; i < input.length; i++) {
    operators.push(web3Default.utils.toChecksumAddress(input[i]));
  }

  const tx: Transaction = {
    from: from,
    to: to,
    data: contractInsDefault.methods.addOperators(operators).encodeABI(),
  }

  return sendTransaction(tx);
}

export function delOperators(from: string, to: string, input: Array<string>): Promise<any> {
  const operators: Array<string> = [];
  for (let i = 0; i < input.length; i++) {
    operators.push(web3Default.utils.toChecksumAddress(input[i]));
  }

  const tx: Transaction = {
    from: from,
    to: to,
    data: contractInsDefault.methods.delOperators(operators).encodeABI(),
  }

  return sendTransaction(tx);
}

export async function getOperators(chainAddr: string, contractAddr: string): Promise<any> {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);
  return await contractIns.methods.getOperators().call();
}

export async function getBatchTokens(chainAddr: string, contractAddr: string, batchCode: string): Promise<any> {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);
  return await contractIns.methods.getBatchTokens(batchCode).call();
}

export async function getBatchTypes(chainAddr: string, contractAddr: string, batchCode: string): Promise<any> {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);
  return await contractIns.methods.getBatchTypes(batchCode).call();
}

export async function getBatchOpenStatus(chainAddr: string, contractAddr: string, batchCode: string): Promise<any> {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);
  return await contractIns.methods.getBatchOpenStatus(batchCode).call();
}

export async function getTokenStatus(chainAddr: string, contractAddr: string, tokenID: string): Promise<any> {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);

  const tokenIDBN = web3.utils.toBN(tokenID);

  return await contractIns.methods.getTokenStatus(tokenIDBN).call();
}

export async function getTokenType(chainAddr: string, contractAddr: string, tokenID: string): Promise<any> {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);

  const tokenIDBN = web3.utils.toBN(tokenID);

  return await contractIns.methods.getTokenType(tokenIDBN).call();
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
