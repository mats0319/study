import detectEthereumProvider from "@metamask/detect-provider";
import { Transaction } from "./web3_kits";

const Web3 = require("web3");

// @ts-ignore
export let metamaskProvider;

export async function initMetamaskProvider(): Promise<void> {
  if (metamaskProvider) {
    return;
  }

  metamaskProvider = await detectEthereumProvider({
    mustBeMetaMask: false,
    silent: false,
    timeout: 1000 // sub timeout
  });
}

export async function sendTransaction(tx: Transaction): Promise<unknown> {
  tx.value = Web3.utils.toHex(tx.value as string)

  await initMetamaskProvider();

  return metamaskProvider.request({
    method: "eth_sendTransaction",
    params: [ tx ]
  })
}