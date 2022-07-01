<template>
  <div class="top">
    <div class="t-connect-button">
      <el-button v-show="!connector || !connector.connected" type="info" plain @click="connectWallet">连接钱包</el-button>
      <el-button v-show="connector && connector.connected" type="danger" plain @click="disConnectWallet">断开连接</el-button>
    </div>

    <div class="t-address">
      <span v-show="connector && connector.connected">地址&#58;</span>
      &nbsp;
      <span v-show="connector && connector.connected">{{ $store.state.address }}</span>
    </div>

    <div class="t-chain-id">
      <span v-show="connector && connector.connected">链ID&#58;</span>
      &nbsp;
      <span v-show="connector && connector.connected">{{ $store.state.chainID }}</span>
    </div>

    <div class="t-empty">
      &nbsp;
      <el-button v-show="false" type="info" plain @click="TEST_func">test func</el-button>
    </div>

    <div class="t-date">
      2022.6.1
    </div>

    <el-dialog :visible.sync="showReConnectWCDialog" title="" append-to-body width="25%">
      <div>
        连接错误&#58;&nbsp;
        <el-button type="info" plain @click="establishWCConnectAgain">重试</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { connector, subscribeSessionEvents, subscribeDisconnectEvent } from "@/ts/wallet_connect";
import { libraries } from "@/views/eth/data";

@Component
export default class Top extends Vue {
  private showReConnectWCDialog = false;

  // const, for html
  private connector = connector;

  private TEST_func(): void {
    console.log("> Node: show libraries.", libraries);
  }

  private mounted() {
    this.maintainConnectStatus();

    if (this.$store.state.reConnect > 0) {
      this.$store.state.reConnect = 0;
      if (connector && connector.connected) {
        connector.killSession();
      }
      this.connectWallet();
    }
  }

  private connectWallet(): void {
    if (connector.connected) {
      return;
    }

    this.showReConnectWCDialog = true;

    connector.createSession();

    this.subscribeWCEvents();
  }

  private disConnectWallet(): void {
    if (!connector || !connector.connected) {
      return;
    }

    connector.killSession();

    this.onDisconnectWallet();
  }

  private establishWCConnectAgain(): void {
    this.$store.state.reConnect = 1;

    location.reload();
  }

  private subscribeWCEvents(): void {
    subscribeSessionEvents((address: string, chainID: number) => {
      this.onConnectWallet(address, chainID);

      this.showReConnectWCDialog = false;
    });
    subscribeDisconnectEvent(() => {
      this.onDisconnectWallet();

      location.reload();
    });
  }

  private onConnectWallet(address: string, chainID: number): void {
    this.$store.state.address = address;
    this.$store.state.chainID = chainID;
  }

  private onDisconnectWallet(): void {
    this.$store.state.address = "";
    this.$store.state.chainID = 0;
  }

  private maintainConnectStatus(): void {
    if (connector && connector.connected && this.$store.state.address.length > 0) {
      this.subscribeWCEvents();
    } else {
      this.disConnectWallet();
    }
  }
}
</script>

<style lang="less">
.top {
  display: flex;
  height: 4rem;
  width: calc(100% - 20rem);
  line-height: 4rem;
  padding: 3rem 10rem;
  text-align: left;
  font-size: 1.8rem;

  .t-connect-button {
    width: 20rem;
  }

  .t-address {
    width: 60rem;
  }

  .t-chain-id {
    width: 20rem;
  }

  .t-empty {
    width: calc(100% - 120rem);
  }

  .t-date {
    width: 20rem;
    font-size: 1.4rem;
    text-align: right;
  }
}
</style>
