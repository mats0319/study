import { dataContractABI, controllerContractABI } from "@/views/cashbox/data";
import { Transaction } from "@/ts/web3_kits";

const Web3 = require("web3");
const web3 = new Web3("");

const dataContractIns = new web3.eth.Contract(dataContractABI, "");
const controllerContractIns = new web3.eth.Contract(controllerContractABI, "");

export function buildTx_SetControllerContract(from: string, to: string, newControllerContract: string): Transaction {
  newControllerContract = web3.utils.toChecksumAddress(newControllerContract);

  return {
    from: from,
    to: to,
    data: dataContractIns.methods.setControllerContract(newControllerContract).encodeABI(),
  }
}

export function buildTx_ReplaceAdmin_DataContract(from: string, to: string, oldAdmin: string, newAdmin: string): Transaction {
  oldAdmin = web3.utils.toChecksumAddress(oldAdmin);
  newAdmin = web3.utils.toChecksumAddress(newAdmin);

  return {
    from: from,
    to: to,
    data: dataContractIns.methods.replaceAdmin(oldAdmin, newAdmin).encodeABI(),
  }
}

export function buildTx_AddMiners(from: string, to: string, miners: Array<string>): Transaction {
  for (let i = 0; i < miners.length; i++) {
    miners[i] = web3.utils.toChecksumAddress(miners[i]);
  }

  return {
    from: from,
    to: to,
    data: controllerContractIns.methods.addMiners(miners).encodeABI(),
  }
}

export function buildTx_DelMiners(from: string, to: string, miners: Array<string>): Transaction {
  for (let i = 0; i < miners.length; i++) {
    miners[i] = web3.utils.toChecksumAddress(miners[i]);
  }

  return {
    from: from,
    to: to,
    data: controllerContractIns.methods.delMiners(miners).encodeABI(),
  }
}

export function buildTx_AddOperators(from: string, to: string, operators: Array<string>): Transaction {
  for (let i = 0; i < operators.length; i++) {
    operators[i] = web3.utils.toChecksumAddress(operators[i]);
  }

  return {
    from: from,
    to: to,
    data: controllerContractIns.methods.addOperators(operators).encodeABI(),
  }
}

export function buildTx_DelOperators(from: string, to: string, operators: Array<string>): Transaction {
  for (let i = 0; i < operators.length; i++) {
    operators[i] = web3.utils.toChecksumAddress(operators[i]);
  }

  return {
    from: from,
    to: to,
    data: controllerContractIns.methods.delOperators(operators).encodeABI(),
  }
}

export function buildTx_ReplaceAdmin_ControllerContract(from: string, to: string, oldAdmin: string, newAdmin: string): Transaction {
  oldAdmin = web3.utils.toChecksumAddress(oldAdmin);
  newAdmin = web3.utils.toChecksumAddress(newAdmin);

  return {
    from: from,
    to: to,
    data: controllerContractIns.methods.replaceAdmin(oldAdmin, newAdmin).encodeABI(),
  }
}
