<template>
  <div class="eth-dependence">
    <div class="ed-info-item">
      <label>链地址&#58;&nbsp;</label>
      <el-input v-model="chainAddr" placeholder="http://..." />
    </div>

    <div class="ed-item">
      <el-collapse class="edi-invokes" v-model="selectedItem" accordion>
        <el-collapse-item v-for="item in libraries" :key="item.id" :title="item.name" :name="item.name">
          <div class="edi-item-input">
            <div>
              <label>部署交易Hash&#58;&nbsp;</label>
              <el-input v-model="item.txHash" placeholder="0x..."></el-input>
            </div>

            <div>
              <label>地址&#58;&nbsp;</label>
              <el-input v-model="item.address" placeholder="0x..."></el-input>
            </div>
          </div>

          <div class="edi-item-button">
            <el-button type="info" plain @click="deploy">部署</el-button>
            <el-button type="info" plain @click="getAddress">获取地址</el-button>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { libraries, chainAddr } from "@/views/cashbox_controller/data";
import { replacePlaceholder } from "@/ts/kits";
import { getTxReceipt, Transaction } from "@/ts/web3_kits";
import { sendTransaction } from "@/ts/wallet_connect";
import { TransactionReceipt } from "web3-core";
import { ElLoadingComponent } from "element-ui/types/loading";

@Component
export default class CashboxControllerDependence extends Vue {
  private selectedItem: string = "";
  private selectedIndex: number = -1; // calc

  private chainAddr = chainAddr

  private loadingDom!: ElLoadingComponent;

  // const, for html
  private libraries = libraries;

  private mounted() {
    // placeholder
  }

  private deploy(): void {
    for (let i = 0; i < this.libraries.length; i++) {
      if (this.selectedItem === this.libraries[i].name) {
        this.selectedIndex = i
        break;
      }
    }

    if (this.selectedIndex === -1) {
      this.$message.info("未知lib名称");
      return;
    }

    if (this.libraries[this.selectedIndex].txHash.length > 0) {
      this.$message.info("已部署，继续操作将重复部署，请确认");
    }

    // check dependencies and replace placeholder
    if (this.libraries[this.selectedIndex].dependOn.length > 0) { // have dependencies
      let errMsg = replacePlaceholder(this.libraries, this.libraries[this.selectedIndex].id);
      if (errMsg.length > 0) {
        this.$message.info(errMsg);
        return;
      }
    }

    const tx: Transaction = {
      from: this.$store.state.address,
      data: this.libraries[this.selectedIndex].byteCode,
    }

    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包执行操作",
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    sendTransaction(tx)
      .then((txHash: any) => {
        this.$message.success("部署成功，hash: " + txHash);
        this.libraries[this.selectedIndex].txHash = txHash;
        console.log("> deploy dependencies success, tx hash: ", txHash);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> deploy dependencies failed, error: ", err);
      })
      .finally(()=> {
        this.loadingDom.close();
      });
  }

  private getAddress(): void {
    for (let i = 0; i < this.libraries.length; i++) {
      if (this.selectedItem === this.libraries[i].name) {
        this.selectedIndex = i
        break;
      }
    }

    if (this.libraries[this.selectedIndex].txHash.length < 1) {
      this.$message.info("请先部署，再尝试获取地址");
      return;
    }

    getTxReceipt(this.chainAddr, this.libraries[this.selectedIndex].txHash)
      .then((receipt: TransactionReceipt) => {
        if (!receipt) {
          this.$message.info("交易尚未打包上链，请稍等片刻后重试");
          return;
        }

        this.$message.success("获取地址成功: " + receipt.contractAddress);
        this.libraries[this.selectedIndex].address = receipt.contractAddress as string;
        console.log("> get tx receipt success, receipt: ", receipt);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get tx receipt failed, error: ", err);
      });
  }
}
</script>

<style lang="less">
.eth-dependence {
  .ed-info-item {
    display: flex;
    text-align: left;
    font-size: 1.8rem;
    line-height: 5rem;
    padding-bottom: 10rem;

    label {
      width: 20%;
    }

    .el-input {
      width: 60%;
    }
  }

  .edi-invokes {
    height: fit-content;
    text-align: right;

    .el-collapse-item__content {
      padding: 0;
    }

    .edi-item-input {
      padding: 5rem 10rem;
      background-color: aliceblue;
      line-height: 5rem;

      label {
        width: 20%;
      }

      .el-input {
        width: 60%;
        padding-left: 10%;
        padding-right: 10%;
      }
    }

    .edi-item-button {
      padding-top: 3rem;
      padding-bottom: 3rem;
      padding-right: 10%;
      background-color: aliceblue;
    }
  }
}
</style>
