<template>
  <div class="top">
    <div class="t-connect-button">
      <el-button
        v-show="!connector || !connector.connected"
        type="info"
        plain
        @click="connectWallet"
      >
        连接钱包
      </el-button>

      <el-button
        v-show="connector && connector.connected"
        type="danger"
        plain
        @click="disConnectWallet"
      >
        断开连接
      </el-button>
    </div>

    <div class="t-address">
      <span v-show="$store.state.address.length > 0">地址&#58;</span>
      &nbsp;
      <span v-show="$store.state.address.length > 0">{{ $store.state.address }}</span>
    </div>

    <div class="t-chain-id">
      <span v-show="connector && connector.connected">链ID&#58;</span>
      &nbsp;
      <span v-show="connector && connector.connected">{{ $store.state.chainID }}</span>
    </div>

    <div class="t-dev">
      <el-button
        v-if="(!connector || !connector.connected) && $store.state.address.length < 1"
        type="info"
        plain
        @click="visitorLogin"
      >
        Visitor Login
      </el-button>

      <el-button
        v-else-if="(!connector || !connector.connected) && $store.state.address.length > 0"
        type="danger"
        plain
        @click="visitorLogout"
      >
        Visitor Logout
      </el-button>

      <el-button type="info" plain @click="showGetTxReceiptDialog = true">Get</el-button>
    </div>

    <div class="t-date">
      2022.9.30
    </div>

    <el-dialog :visible.sync="showReConnectWCDialog" title="" append-to-body width="25%">
      <div>
        连接错误&#58;&nbsp;
        <el-button type="info" plain @click="establishWCConnectAgain">重试</el-button>
      </div>
    </el-dialog>

    <el-dialog :visible.sync="showGetTxReceiptDialog" append-to-body width="fit-content">
      <div>
        <el-input v-model="chainAddress" placeholder="chain address" />
        <el-input v-model="txHash" placeholder="block/tx hash" />
      </div>

      <div slot="footer">
        <el-button type="info" plain @click="getBlock">Block</el-button>
        <el-button type="info" plain @click="getTx">Transaction</el-button>
        <el-button type="info" plain @click="getTxReceipt">Tx receipt</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { connector, subscribeSessionEvents, subscribeDisconnectEvent } from "@/ts/wallet_connect";
import {
  getBlock, Block,
  getTx, Tx,
  getTxReceipt, TransactionReceipt
} from "@/ts/web3_kits";

@Component
export default class Top extends Vue {
  private showReConnectWCDialog = false;
  private showGetTxReceiptDialog = false;

  private chainAddress = ""; // for get receipt
  private txHash = ""; // for get receipt

  // const, for html
  private connector = connector;

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

  private visitorLogin(): void {
    this.$store.state.address = "Visitor";
  }

  private visitorLogout(): void {
    this.$store.state.address = "";

    this.redirectHome();
  }

  private connectWallet(): void {
    if (connector.connected) {
      return;
    }

    this.showReConnectWCDialog = true;

    connector.createSession();

    this.subscribeWCEvents();
  }

  private getBlock(): void {
    getBlock(this.chainAddress, this.txHash)
      .then((block: Block) => {
        this.$message.success("获取区块信息成功")
        console.log("> get block info success: ", block);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get block info failed, error: ", err);
      })
  }

  private getTx(): void {
    getTx(this.chainAddress, this.txHash)
      .then((tx: Tx) => {
        this.$message.success("获取交易成功");
        console.log("> get tx success: ", tx);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get tx failed, error: ", err);
      });
  }

  private getTxReceipt(): void {
    getTxReceipt(this.chainAddress, this.txHash)
      .then((receipt: TransactionReceipt) => {
        if (!receipt) {
          this.$message.info("交易尚未打包上链，请稍等片刻后重试");
          return;
        }

        this.$message.success("获取交易收据成功");
        console.log("> get tx receipt success, receipt: ", receipt);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get tx receipt failed, error: ", err);
      });
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

    this.redirectHome();
  }

  private maintainConnectStatus(): void {
    if (connector && connector.connected && this.$store.state.address.length > 0) {
      this.subscribeWCEvents();
    } else {
      this.disConnectWallet();
    }
  }

  private redirectHome(): void {
    if (window.location.href.split("#/")[1].length > 0) {
      this.$router.push({ name: "home" });
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

  .t-dev {
    width: calc(100% - 120rem);
  }

  .t-date {
    width: 20rem;
    font-size: 1.4rem;
    text-align: right;
  }
}
</style>
