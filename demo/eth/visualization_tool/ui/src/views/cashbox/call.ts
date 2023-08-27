import Web3 from "web3";
import { dataContractABI, controllerContractABI } from "@/views/cashbox/data";
import { Message } from "element-ui";

export function call_GetControllerContract(chainAddr: string, contractAddr: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(dataContractABI, contractAddr);

  contractIns.methods.getControllerContract().call()
    .then((address: any) => {
      Message.success("当前C合约地址：" + address);
      console.log("> get C contract success, address：", address);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get C contract failed, error: ", err);
    });
}

export function call_GetAdmin_DataContract(chainAddr: string, contractAddr: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(dataContractABI, contractAddr);

  contractIns.methods.getAdmin().call()
    .then((address: any) => {
      Message.success("当前管理员：" + address);
      console.log("> get D contract admin success, address：", address);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get D contract admin failed, error: ", err);
    });
}

export function call_GetBatchTokens(chainAddr: string, contractAddr: string, batchCode: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(dataContractABI, contractAddr);

  contractIns.methods.getBatchTokens(batchCode).call()
    .then((tokens: any) => {
      Message.success("指定批次代币数量：" + tokens.length);
      console.log("> get batch tokens success, address：", tokens);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get batch tokens failed, error: ", err);
    });
}

export function call_GetBatchTypes(chainAddr: string, contractAddr: string, batchCode: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(dataContractABI, contractAddr);

  contractIns.methods.getBatchTypes(batchCode).call()
    .then((tokenTypes: any) => {
      Message.success("指定批次类型：" + tokenTypes);
      console.log("> get batch types success, address：", tokenTypes);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get batch types failed, error: ", err);
    });
}

export function call_GetBatchOpenStatus(chainAddr: string, contractAddr: string, batchCode: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(dataContractABI, contractAddr);

  contractIns.methods.getBatchOpenStatus(batchCode).call()
    .then((openStatus: any) => {
      Message.success("指定批次开启情况：" + openStatus);
      console.log("> get batch open status success, address：", openStatus);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get batch open status failed, error: ", err);
    });
}

export function call_GetTokenStatus(chainAddr: string, contractAddr: string, tokenID: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(dataContractABI, contractAddr);

  contractIns.methods.getTokenStatus(tokenID).call()
    .then((tokenStatus: any) => {
      Message.success("指定代币状态：" + tokenStatus);
      console.log("> get token status success, address：", tokenStatus);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get token status failed, error: ", err);
    });
}

export function call_GetTokenOpenResult(chainAddr: string, contractAddr: string, tokenID: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(dataContractABI, contractAddr);

  contractIns.methods.getTokenOpenResult(tokenID).call()
    .then((tokenStatus: any) => {
      Message.success("指定代币开启结果：" + tokenStatus);
      console.log("> get token open result success, address：", tokenStatus);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get token open result failed, error: ", err);
    });
}

export function call_GetMiners(chainAddr: string, contractAddr: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(controllerContractABI, contractAddr);

  contractIns.methods.getMiners().call()
    .then((miners: any) => {
      Message.success("当前铸币者：" + miners);
      console.log("> get miners success, address：", miners);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get miners failed, error: ", err);
    });
}

export function call_GetOperators(chainAddr: string, contractAddr: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(controllerContractABI, contractAddr);

  contractIns.methods.getOperators().call()
    .then((miners: any) => {
      Message.success("当前操作员：" + miners);
      console.log("> get operators success, address：", miners);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get operators failed, error: ", err);
    });
}

export function call_GetAdmin_ControllerContract(chainAddr: string, contractAddr: string): void {
  const web3 = new Web3(chainAddr);
  const contractIns = new web3.eth.Contract(controllerContractABI, contractAddr);

  contractIns.methods.getAdmin().call()
    .then((miners: any) => {
      Message.success("当前管理员：" + miners);
      console.log("> get C contract admin success, address：", miners);
    })
    .catch((err: any) => {
      Message.error(err);
      console.log("> get C contract admin failed, error: ", err);
    });
}
