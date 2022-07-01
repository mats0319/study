<template>
  <div class="eth-deploy">
    <div class="ed-info-item">
      <label>链地址&#58;&nbsp;</label>
      <el-input v-model="chainAddr" placeholder="http://..." />
    </div>

    <div class="ed-constructor-params">
      <el-form
        v-model="params"
        ref="form"
        label-width="15rem"
        status-icon
      >
        <el-form-item label="管理员1">
          <el-input v-model="params.admin1" placeholder="0x..." @change="validAdmin1" />
          <i v-show="validParams.admin1 === 0" class="el-icon-loading gray" />
          <i v-show="validParams.admin1 > 0" class="el-icon-check green" />
          <i v-show="validParams.admin1 < 0" class="el-icon-close red" />
        </el-form-item>

        <el-form-item label="管理员2">
          <el-input v-model="params.admin2" placeholder="0x..." @change="validAdmin2" />
          <i v-show="validParams.admin2 === 0" class="el-icon-loading gray" />
          <i v-show="validParams.admin2 > 0" class="el-icon-check green" />
          <i v-show="validParams.admin2 < 0" class="el-icon-close red" />
        </el-form-item>

        <el-form-item label="管理员3">
          <el-input v-model="params.admin3" placeholder="0x..." @change="validAdmin3" />
          <i v-show="validParams.admin3 === 0" class="el-icon-loading gray" />
          <i v-show="validParams.admin3 > 0" class="el-icon-check green" />
          <i v-show="validParams.admin3 < 0" class="el-icon-close red" />
        </el-form-item>

        <el-form-item label="决策者">
          <el-input v-model="params.decisionMaker" placeholder="0x..." @change="validDecisionMaker" />
          <i v-show="validParams.decisionMaker === 0" class="el-icon-loading gray" />
          <i v-show="validParams.decisionMaker > 0" class="el-icon-check green" />
          <i v-show="validParams.decisionMaker < 0" class="el-icon-close red" />
        </el-form-item>

        <el-form-item label="ether收款地址">
          <el-input v-model="params.ethReceiver" placeholder="0x..." @change="validETHReceiver" />
          <i v-show="validParams.ethReceiver === 0" class="el-icon-loading gray" />
          <i v-show="validParams.ethReceiver > 0" class="el-icon-check green" />
          <i v-show="validParams.ethReceiver < 0" class="el-icon-close red" />
        </el-form-item>
      </el-form>
    </div>

    <div class="ed-deploy-button">
      <el-button type="info" plain @click="deploy">部署合约</el-button>
      <el-button type="info" plain @click="getAddress">获取地址</el-button>
    </div>

    <div class="ed-byte-code">
      <div class="edbc-label">合约字节码（包含占位符），更新时间：2022.5.31</div>
      <el-input v-model="contractByteCode" type="textarea" :rows="10" resize="none" readonly />
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { encodeParameters, getTxReceipt, Transaction } from "@/ts/web3_kits";
import { sendTransaction } from "@/ts/wallet_connect";
import { isValidAddress } from "@/ts/kits";
import { chainAddr, contractByteCode, libraries } from "./data"
import { TransactionReceipt } from "web3-core";
import { ElLoadingComponent } from "element-ui/types/loading";

@Component
export default class ETHDeploy extends Vue {
  private params = {
    admin1: "",
    admin2: "",
    admin3: "",
    decisionMaker: "",
    ethReceiver: "",
  }

  private validParams = {
    admin1: 0,
    admin2: 0,
    admin3: 0,
    decisionMaker: 0,
    ethReceiver: 0,
  }

  private deployTxHash = "";

  private chainAddr = chainAddr

  private loadingDom!: ElLoadingComponent;

  // const, for html
  private contractByteCode = contractByteCode;

  private mounted() {
    // placeholder
  }

  private deploy(): void {
    for (let i = 0; i < libraries.length; i++) {
      if (libraries[i].address.length < 1) {
        this.$message.info("存在尚未部署的依赖，请部署全部依赖后重试");
        return;
      }
    }

    if (this.validParams.admin1 <= 0 ||
      this.validParams.admin2 <= 0 ||
      this.validParams.admin3 <= 0 ||
      this.validParams.decisionMaker <= 0 ||
      this.validParams.ethReceiver <= 0) {
      this.$message.info("当前输入参数有误，请检查后重试");
      return;
    }

    const params = encodeParameters([ "address[3]", "address", "address" ], [
      [ this.params.admin1, this.params.admin2, this.params.admin3 ],
      this.params.decisionMaker,
      this.params.ethReceiver,
    ])

    let byteCode = contractByteCode;
    for (let i = 0; i < libraries.length; i++) {
      byteCode = byteCode.replaceAll(libraries[i].placeholder, libraries[i].address.slice(2));
    }

    const tx: Transaction = {
      from: this.$store.state.address,
      data: byteCode + params,
    }

    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包执行操作",
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    sendTransaction(tx)
      .then((txHash: any) => {
        this.$message.success("部署成功，hash: " + txHash);
        this.deployTxHash = txHash;
        console.log("> deploy eth contract success, tx hash: ", txHash);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> deploy eth contract failed, error: ", err);
      })
      .finally(()=>{
        this.loadingDom.close()
      });
  }

  private getAddress(): void {
    if (this.deployTxHash.length < 1) {
      this.$message.info("请先部署，再尝试获取地址");
      return;
    }

    getTxReceipt(this.chainAddr, this.deployTxHash)
      .then((receipt: TransactionReceipt) => {
        if (!receipt) {
          this.$message.info("交易尚未打包上链，请稍等片刻后重试");
          return;
        }

        this.$message.success("获取地址成功，请保存该地址："+receipt.contractAddress);
        console.log("> get tx receipt success, receipt: ", receipt);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get tx receipt failed, error: ", err);
      });
  }

  private validAdmin1(value: string): void {
    this.validParams.admin1 = isValidAddress(value) ? 1 : -1;
  }

  private validAdmin2(value: string): void {
    this.validParams.admin2 = isValidAddress(value) ? 1 : -1;
  }

  private validAdmin3(value: string): void {
    this.validParams.admin3 = isValidAddress(value) ? 1 : -1;
  }

  private validDecisionMaker(value: string): void {
    this.validParams.decisionMaker = isValidAddress(value) ? 1 : -1;
  }

  private validETHReceiver(value: string): void {
    this.validParams.ethReceiver = isValidAddress(value) ? 1 : -1;
  }
}
</script>

<style lang="less">
.eth-deploy {
  .ed-info-item {
    display: flex;
    font-size: 1.8rem;
    line-height: 5rem;

    label {
      width: 20%;
    }

    .el-input {
      width: 60%;
    }
  }

  .ed-constructor-params {
    margin: 5rem 20rem;
    padding: 1rem;

    .el-form-item__error {
      width: 100%;
      text-align: center;
    }

    .el-input {
      width: 40rem;
    }

    i {
      width: 5rem;
      font-size: 2rem;
    }

    .gray {
      color: darkgray;
    }

    .green {
      color: forestgreen;
    }

    .red {
      color: red;
    }
  }

  .ed-deploy-button {

  }

  .ed-byte-code {
    margin: auto;
    padding-top: 5rem;
    width: 80%;

    .edbc-label {
      width: inherit;
      text-align: left;
      font-size: 1.8rem;
    }
  }
}
</style>
