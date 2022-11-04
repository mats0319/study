<template>
  <div class="eth-deploy">
    <div class="ed-info-item">
      <label>链地址&#58;&nbsp;</label>
      <el-input v-model="$store.state.cashboxChainAddress" placeholder="http://127.0.0.1:8545" />
    </div>

    <div class="ed-constructor-params">
      <el-form v-model="params" ref="form" label-width="15rem" status-icon>
        <el-form-item label="管理员1">
          <el-input v-model="params.admin1" placeholder="0x..." @change="validAdmin1" />
          <i v-show="validParams.admin1 > 0" class="el-icon-check green" />
          <i v-show="validParams.admin1 < 0" class="el-icon-close red" />
        </el-form-item>

        <el-form-item label="管理员2">
          <el-input v-model="params.admin2" placeholder="0x..." @change="validAdmin2" />
          <i v-show="validParams.admin2 > 0" class="el-icon-check green" />
          <i v-show="validParams.admin2 < 0" class="el-icon-close red" />
        </el-form-item>

        <el-form-item label="管理员3">
          <el-input v-model="params.admin3" placeholder="0x..." @change="validAdmin3" />
          <i v-show="validParams.admin3 > 0" class="el-icon-check green" />
          <i v-show="validParams.admin3 < 0" class="el-icon-close red" />
        </el-form-item>

        <el-form-item label="决策者">
          <el-input v-model="params.decisionMaker" placeholder="0x..." @change="validDecisionMaker" />
          <i v-show="validParams.decisionMaker > 0" class="el-icon-check green" />
          <i v-show="validParams.decisionMaker < 0" class="el-icon-close red" />
        </el-form-item>
      </el-form>
    </div>

    <div>
      <el-button type="info" plain @click="deploy">部署合约</el-button>
      <el-button type="info" plain @click="getAddress">获取地址</el-button>
    </div>

    <div class="ed-byte-code">
      <div class="edbc-label">合约字节码(hardhat bytecode)，仅供了解</div>
      <el-input v-model="contractByteCode" type="textarea" :rows="10" resize="none" readonly />
    </div>

    <send-transaction ref="sendTx" @success="onSendTxSuccess" />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { encodeParameters, getTxReceipt, Transaction } from "@/ts/web3_kits";
import { isValidAddress } from "@/ts/utils";
import { dataContractByteCode, libraries } from "@/views/cashbox/data"
import { TransactionReceipt } from "web3-core";
import SendTransaction from "@/components/send_transaction.vue";

@Component({
  components: { SendTransaction }
})
export default class CashboxControllerDeployD extends Vue {
  private params = {
    admin1: "",
    admin2: "",
    admin3: "",
    decisionMaker: "",
  }

  private validParams = {
    admin1: -1,
    admin2: -1,
    admin3: -1,
    decisionMaker: -1,
  }

  private deployTxHash = "";

  // const, for html
  private contractByteCode = dataContractByteCode;

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
      this.validParams.decisionMaker <= 0
    ) {
      this.$message.info("当前输入参数有误，请检查后重试");
      return;
    }

    let byteCode = dataContractByteCode;
    for (let i = 0; i < libraries.length; i++) {
      byteCode = byteCode.replaceAll(libraries[i].placeholder, libraries[i].address.slice(2));
    }

    const params = encodeParameters([ "address[3]", "address" ], [
      [ this.params.admin1, this.params.admin2, this.params.admin3 ],
      this.params.decisionMaker,
    ])

    const tx: Transaction = {
      from: this.$store.state.address,
      data: byteCode + params,
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(tx, 0, "deploy cashbox D contract");
  }

  private onSendTxSuccess(flag: number, txHash: any): void {
    this.deployTxHash = txHash;
  }

  private getAddress(): void {
    if (this.deployTxHash.length < 1) {
      this.$message.info("请先部署，再尝试获取地址");
      return;
    }

    getTxReceipt(this.$store.state.cashboxChainAddress, this.deployTxHash, (receipt: TransactionReceipt) => {
      this.$message.success("获取地址成功，请保存该地址：" + receipt.contractAddress);
      this.$store.state.cashboxDataContractAddress = receipt.contractAddress;
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
    border: 1px solid darkgray;

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

    .green {
      color: forestgreen;
    }

    .red {
      color: red;
    }
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
