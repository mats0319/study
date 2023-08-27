<template>
  <div></div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { Transaction } from "@/ts/web3_kits";
import { ElLoadingComponent } from "element-ui/types/loading";
import { sendTransaction as sendTxByWalletconnect } from "@/ts/wallet_connect";
import { sendTransaction as sendTxByMetamaskPlugin } from "@/ts/metamask_plugin";

@Component
export default class SendTransaction extends Vue {
  // modal
  private loadingDom!: ElLoadingComponent;

  private mounted() {
    // placeholder
  }

  private sendTx(tx: Transaction, flag: number, description: string): void {
    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包发送交易", // text under loading icon
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    switch (this.$store.state.connectMethod) {
      case 1:
        sendTxByWalletconnect(tx)
          .then((txHash: any) => {
            this.$message.success("> " + description + " success.");
            console.log("> " + description + " success, tx hash: ", txHash);
            this.$emit("success", flag, txHash);
          })
          .catch((err: any) => {
            this.$message.error("> " + description + " failed, error: " + err);
            console.log("> " + description + " failed, error: ", err);
          })
          .finally(() => {
            this.loadingDom.close();
          });
        break;
      case 2:
        sendTxByMetamaskPlugin(tx)
          .then((txHash: any) => {
            this.$message.success("> " + description + " success.");
            console.log("> " + description + " success, tx hash: ", txHash);
            this.$emit("success", flag, txHash as string);
          })
          .catch((err: any) => {
            this.$message.error("> " + description + " failed, error: " + err);
            console.log("> " + description + " failed, error: ", err);
          })
          .finally(() => {
            this.loadingDom.close();
          });
        break;
      default:
        console.log("unknown connect method: ", this.$store.state.connectMethod);
        this.loadingDom.close();
    }
  }
}
</script>
