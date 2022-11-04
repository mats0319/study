import Web3 from "web3";
import { contractABI } from "@/views/eth/data";
import { Message } from "element-ui";

export function getETHReceiver(chainAddr: string, contractAddr: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);

  contractIns.methods.getETHReceiver().call()
    .then((address: any) => {
      Message.success("当前ETH收款地址：" + address);
      console.log("> get eth receiver success, address：", address);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get eth receiver failed, error: ", err);
    });
}

export function getACashOperators(chainAddr: string, contractAddr: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);

  contractIns.methods.getACashOperators().call()
    .then((address: any) => {
      Message.success("当前ACash提现操作员地址：" + address);
      console.log("> get ACash operators success, address：", address);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get ACash operators failed, error: ", err);
    });
}

export function getBrokerageOperators(chainAddr: string, contractAddr: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);

  contractIns.methods.getBrokerageOperators().call()
    .then((address: any) => {
      Message.success("当前佣金提现操作员地址：" + address);
      console.log("> get brokerage operators success, address：", address);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get brokerage operators failed, error: ", err);
    });
}

export function getAdmin(chainAddr: string, contractAddr: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(contractABI, contractAddr);

  contractIns.methods.getAdmin().call()
    .then((address: any) => {
      Message.success("当前管理员：" + address);
      console.log("> get admin success, address：", address);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get admin failed, error: ", err);
    });
}