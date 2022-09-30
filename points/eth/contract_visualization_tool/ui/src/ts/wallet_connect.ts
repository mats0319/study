import WalletConnect from "@walletconnect/client";
import QRCodeModal from "@walletconnect/qrcode-modal";
import { Transaction } from "@/ts/web3_kits";

// Create a connector
export const connector = new WalletConnect({
  bridge: "https://bridge.walletconnect.org", // Required
  qrcodeModal: QRCodeModal
});

// 订阅session相关事件，包括首次创建session(connect)和更新session(session_update)
export function subscribeSessionEvents(cb: (address: string, chainID: number) => void): void {
  connector.on("connect", (error, payload) => {
    if (error) {
      throw error;
    }

    console.log("> Node: wallet connect on connect.");
    const { accounts, chainId: chainID } = payload.params[0];
    cb(accounts[0], chainID);
  });

  connector.on("session_update", (error, payload) => {
    if (error) {
      throw error;
    }

    console.log("> Node: wallet connect on session_update.");
    const { accounts, chainId: chainID } = payload.params[0];
    cb(accounts[0], chainID);
  });
}

export function subscribeDisconnectEvent(cb: () => void): void {
  connector.on("disconnect", (error, payload) => {
    if (error) {
      throw error;
    }

    console.log("> Node: wallet connect on disconnect.", payload);
    cb();

    // unsubscribe all events
    connector.off("connect");
    connector.off("session_update");
    connector.off("disconnect");
  });
}

export function sendTransaction(tx: Transaction): Promise<any> {
  console.log("> Node: show tx before send to wallet.", tx);
  return connector.sendTransaction(tx);
}
