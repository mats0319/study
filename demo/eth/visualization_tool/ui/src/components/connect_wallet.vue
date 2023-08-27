<template>
  <div>
    <el-button v-if="!$store.state.isConnected" type="info" plain @click="showSelectConnectMethodDialog = true">
      连接钱包
    </el-button>

    <el-button v-else type="danger" plain @click="disconnectWallet">
      断开连接
    </el-button>

    <el-dialog
      class="select-connect-method-dialog"
      :visible.sync="showSelectConnectMethodDialog"
      append-to-body
      width="30vw"
    >
      <div class="scmd-methods">
        <div class="scmdm-item">
          <el-button type="info" plain @click="connectWallet(1)">walletconnect</el-button>
        </div>

        <div class="scmdm-item">
          <el-button type="info" plain @click="connectWallet(2)">metamask plugin</el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { connector, subscribeSessionEvents, subscribeDisconnectEvent } from "@/ts/wallet_connect";
import { initMetamaskProvider, metamaskProvider } from "@/ts/metamask_plugin";
import { hexToNumber } from "@/ts/web3_kits";

@Component
export default class ConnectWallet extends Vue {
  private showSelectConnectMethodDialog = false;

  private mounted() {
    this.maintainConnectStatus();
  }

  private async connectWallet(connectMethod: number): Promise<void> {
    switch (connectMethod) {
      case 1:
        await connector.createSession();

        this.subscribeWCEvents();
        break;
      case 2:
        await initMetamaskProvider();

        metamaskProvider.request({ method: "eth_requestAccounts" })
          .then((accounts: Array<string>) => {
            this.onConnectWallet(2, accounts[0].toLowerCase(), -1);
          });
        metamaskProvider.request({ method: "eth_chainId" })
          .then((chainID: string) => {
            this.$store.state.chainID = hexToNumber(chainID);
          });
        break;
    }
  }

  private async disconnectWallet(): Promise<void> {
    switch (this.$store.state.connectMethod) {
      case 1:
        if (connector && connector.connected) {
          await connector.killSession();
        }
        break;
      case 2:
        this.$message.info("请前往metamask浏览器插件钱包断开连接")
        break;
    }

    this.onDisconnectWallet();
  }

  private subscribeWCEvents(): void {
    subscribeSessionEvents((address: string, chainID: number) => {
      this.onConnectWallet(1, address.toLowerCase(), chainID);
    });
    subscribeDisconnectEvent(() => {
      this.onDisconnectWallet();
      location.reload();
    });
  }

  private async maintainConnectStatus(): Promise<void> {
    switch (this.$store.state.connectMethod) {
      case 1:
        if (connector && connector.connected) {
          this.$store.state.isConnected = true;
          this.subscribeWCEvents();
        } else {
          this.$store.state.isConnected = false;
          await this.disconnectWallet();
        }
        break;
      case 2:
        await initMetamaskProvider();
        this.$store.state.isConnected = !!metamaskProvider;
    }
  }

  private onConnectWallet(connectMethod: number, address: string, chainID: number): void {
    this.$store.state.isConnected = true;
    this.$store.state.connectMethod = connectMethod;
    this.$store.state.address = address;
    if (chainID > 0) {
      this.$store.state.chainID = chainID;
    }

    this.showSelectConnectMethodDialog = false;
  }

  private onDisconnectWallet(): void {
    this.$store.state.isConnected = false;
    this.$store.state.connectMethod = 0;
    this.$store.state.address = "";
    this.$store.state.chainID = 0;
  }
}
</script>

<style lang="less">
.select-connect-method-dialog {
  display: flex;
  height: fit-content;

  .scmd-methods {
    .scmdm-item {
      width: 80%;
      margin: 2vh auto;

      .el-button {
        width: 100%;
      }
    }
  }
}
</style>
