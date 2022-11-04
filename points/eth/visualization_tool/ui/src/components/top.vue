<template>
  <div class="top">
    <div class="t-connect-button">
      <connect-wallet />
    </div>

    <div class="t-address">
      <span v-show="$store.state.address.length > 0">地址&#58;&nbsp;</span>
      <span v-show="$store.state.address.length > 0">{{ $store.state.address }}</span>
    </div>

    <div class="t-chain-id">
      <span v-show="$store.state.isConnected">链ID&#58;&nbsp;</span>
      <span v-show="$store.state.isConnected">{{ $store.state.chainID }}</span>
    </div>

    <div class="t-dev">
      <el-button
        v-if="!$store.state.isConnected && $store.state.address.length < 1"
        type="info"
        plain
        @click="visitorLogin"
      >
        Visitor Login
      </el-button>

      <el-button
        v-else-if="!$store.state.isConnected && $store.state.address.length > 0"
        type="danger"
        plain
        @click="visitorLogout"
      >
        Visitor Logout
      </el-button>

      <el-button type="info" plain @click="showGetDevTool = true">Get</el-button>
    </div>

    <div class="t-date">2022.10.20</div>

    <el-dialog :visible.sync="showGetDevTool" append-to-body width="fit-content">
      <div>
        <el-input v-model="chainAddress" placeholder="http://127.0.0.1:8545" />
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
import { Block, getBlock, Tx, getTx, TransactionReceipt, getTxReceipt } from "@/ts/web3_kits";
import ConnectWallet from "@/components/connect_wallet.vue";

@Component({
  components: { ConnectWallet }
})
export default class Top extends Vue {
  private showGetDevTool = false; // for dev tool
  private chainAddress = ""; // for dev tool
  private txHash = ""; // for dev tool

  private mounted() {
    // placeholder
  }

  private visitorLogin(): void {
    this.$store.state.address = "Visitor";
  }

  private visitorLogout(): void {
    this.$store.state.address = "";

    this.redirectHome();
  }

  private getBlock(): void {
    getBlock(this.chainAddress, this.txHash);
  }

  private getTx(): void {
    getTx(this.chainAddress, this.txHash);
  }

  private getTxReceipt(): void {
    getTxReceipt(this.chainAddress, this.txHash, (receipt: TransactionReceipt) => {
      this.$message.success("获取交易收据成功");
    });
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
