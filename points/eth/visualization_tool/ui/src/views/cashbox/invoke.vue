<template>
  <div class="eth-manage">
    <div class="em-info">
      <div class="em-info-item">
        <label>链地址&#58;&nbsp;</label>
        <el-input v-model="$store.state.cashboxChainAddress" placeholder="http://127.0.0.1:8545" />
      </div>

      <div class="em-info-item">
        <label>数据合约地址&#58;&nbsp;</label>
        <el-input v-model="$store.state.cashboxDataContractAddress" placeholder="0x..." />
      </div>

      <div class="em-info-item">
        <label>控制合约地址&#58;&nbsp;</label>
        <el-input v-model="$store.state.cashboxControllerContractAddress" placeholder="0x..." />
      </div>
    </div>

    <el-collapse
      v-show="$store.state.cashboxChainAddress.length > 0 &&
      $store.state.cashboxDataContractAddress.length > 0 &&
      $store.state.cashboxControllerContractAddress.length > 0"
      class="em-invokes"
      v-model="selectedItem"
      accordion
    >
      <el-collapse-item title="设置控制合约地址（数据合约）" name="setControllerContract">
        <div class="emi-item-input">
          <label>控制合约地址&#58;&nbsp;</label>
          <el-input v-model="newControllerContractAddress" placeholder="0x..." />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="sendTx_SetControllerContract">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询控制合约地址（数据合约）" name="getControllerContract">
        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetControllerContract">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="替换管理员（数据合约）" name="replaceAdmin_DataContract">
        <div class="emi-item-input">
          <div>
            <label>老管理员&#58;&nbsp;</label>
            <el-input v-model="oldAdminInDataContract" placeholder="0x..." />
          </div>

          <div>
            <label>新管理员&#58;&nbsp;</label>
            <el-input v-model="newAdminInDataContract" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="sendTx_ReplaceAdmin_DataContract">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询管理员（数据合约）" name="getAdmin_DataContract">
        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetAdmin_DataContract">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="新增铸币者" name="addMiners">
        <div class="emi-item-input">
          <div v-for="(v, index) in addMinerList" :key="index">
            <label>新的铸币者{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="addMinerList[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="addMinerList.push('')">增加铸币者</el-button>
            <el-button type="info" plain @click="addMinerList.pop()">删除铸币者</el-button>
          </span>

          <el-button type="info" plain @click="sendTx_AddMiners">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="删除铸币者" name="delMiners">
        <div class="emi-item-input">
          <div v-for="(v, index) in delMinerList" :key="index">
            <label>铸币者{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="delMinerList[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="delMinerList.push('')">增加铸币者</el-button>
            <el-button type="info" plain @click="delMinerList.pop()">删除铸币者</el-button>
          </span>

          <el-button type="info" plain @click="sendTx_DelMiners">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询铸币者" name="getMiners">
        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetMiners">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="新增操作员" name="addOperator">
        <div class="emi-item-input">
          <div v-for="(v, index) in addOperatorList" :key="index">
            <label>新的操作员{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="addOperatorList[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="addOperatorList.push('')">增加操作员</el-button>
            <el-button type="info" plain @click="addOperatorList.pop()">删除操作员</el-button>
          </span>

          <el-button type="info" plain @click="sendTx_AddOperators">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="删除操作员" name="delOperator">
        <div class="emi-item-input">
          <div v-for="(v, index) in delOperatorList" :key="index">
            <label>操作员{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="delOperatorList[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="delOperatorList.push('')">增加操作员</el-button>
            <el-button type="info" plain @click="delOperatorList.pop()">删除操作员</el-button>
          </span>

          <el-button type="info" plain @click="sendTx_DelOperators">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询操作员" name="getOperators">
        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetOperators">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="替换管理员" name="replaceAdmin">
        <div class="emi-item-input">
          <div>
            <label>老管理员&#58;&nbsp;</label>
            <el-input v-model="oldAdminInControllerContract" placeholder="0x..." />
          </div>

          <div>
            <label>新管理员&#58;&nbsp;</label>
            <el-input v-model="newAdminInControllerContract" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="sendTx_ReplaceAdmin_ControllerData">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询管理员" name="getAdmin">
        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetAdmin_ControllerData">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询指定批次的代币（数据合约）" name="getBatchTokens">
        <div class="emi-item-input">
          <label>批次码&#58;&nbsp;</label>
          <el-input v-model="batchCodeForTokens" placeholder="..." />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetBatchTokens">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询指定批次的代币类型（数据合约）" name="getBatchTypes">
        <div class="emi-item-input">
          <label>批次码&#58;&nbsp;</label>
          <el-input v-model="batchCodeForTypes" placeholder="..." />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetBatchTypes">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询指定批次的代币的开启情况（数据合约）" name="getBatchOpenStatus">
        <div class="emi-item-input">
          <label>批次码&#58;&nbsp;</label>
          <el-input v-model="batchCodeForOpenStatus" placeholder="..." />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetBatchOpenStatus">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询指定代币的状态（数据合约）" name="getTokenStatus">
        <div class="emi-item-input">
          <label>token id&#58;&nbsp;</label>
          <el-input v-model="tokenIDForStatus" placeholder="0x..." />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetTokenStatus">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询指定代币的开启结果（数据合约）" name="getTokenOpenResult">
        <div class="emi-item-input">
          <label>token id&#58;&nbsp;</label>
          <el-input v-model="tokenIDForOpenResult" placeholder="0x..." />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetTokenOpenResult">查询</el-button>
        </div>
      </el-collapse-item>
    </el-collapse>

    <send-transaction ref="sendTx" />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { isValidAddress, isValidTokenID } from "@/ts/utils";
import {
  buildTx_SetControllerContract,
  buildTx_ReplaceAdmin_DataContract,
  buildTx_AddMiners,
  buildTx_DelMiners,
  buildTx_AddOperators,
  buildTx_DelOperators,
  buildTx_ReplaceAdmin_ControllerContract,
} from "@/views/cashbox/build_tx";
import {
  call_GetControllerContract,
  call_GetAdmin_DataContract,
  call_GetMiners,
  call_GetOperators,
  call_GetAdmin_ControllerContract,
  call_GetBatchTokens,
  call_GetBatchTypes,
  call_GetBatchOpenStatus,
  call_GetTokenStatus,
  call_GetTokenOpenResult,
} from "@/views/cashbox/call";
import SendTransaction from "@/components/send_transaction.vue";

@Component({
  components: { SendTransaction }
})
export default class CashboxControllerManage extends Vue {
  private selectedItem = ""; // no use

  private newControllerContractAddress = "";
  private oldAdminInDataContract = "";
  private newAdminInDataContract = "";
  private addMinerList: Array<string> = [ "" ];
  private delMinerList: Array<string> = [ "" ];
  private addOperatorList: Array<string> = [ "" ];
  private delOperatorList: Array<string> = [ "" ];
  private oldAdminInControllerContract = "";
  private newAdminInControllerContract = "";
  private batchCodeForTokens = "";
  private batchCodeForTypes = "";
  private batchCodeForOpenStatus = "";
  private tokenIDForStatus = "";
  private tokenIDForOpenResult = "";

  private mounted() {
    this.newControllerContractAddress = this.$store.state.cashboxControllerContractAddress;
  }

  private sendTx_SetControllerContract(): void {
    if (this.$store.state.cashboxDataContractAddress.length < 1) {
      this.$message.info("请先部署数据合约")
      return;
    }

    if (this.$store.state.cashboxControllerContractAddress.length < 1) {
      this.$message.info("请先部署控制合约")
      return;
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(buildTx_SetControllerContract(
        this.$store.state.address,
        this.$store.state.cashboxDataContractAddress,
        this.$store.state.cashboxControllerContractAddress,
      ),
      0,
      "set C contract address",
    );
  }

  private call_GetControllerContract(): void {
    call_GetControllerContract(this.$store.state.cashboxChainAddress, this.$store.state.cashboxDataContractAddress);
  }

  private sendTx_ReplaceAdmin_DataContract(): void {
    if (!isValidAddress(this.oldAdminInDataContract) || !isValidAddress(this.newAdminInDataContract)) {
      this.$message.info("当前输入参数有误，请检查后重试");
      return;
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(buildTx_ReplaceAdmin_DataContract(
        this.$store.state.address,
        this.$store.state.cashboxDataContractAddress,
        this.oldAdminInDataContract,
        this.newAdminInDataContract,
      ),
      0,
      "replace admin in data contract",
    );
  }

  private call_GetAdmin_DataContract(): void {
    call_GetAdmin_DataContract(this.$store.state.cashboxChainAddress, this.$store.state.cashboxDataContractAddress);
  }

  private sendTx_AddMiners(): void {
    if (this.addMinerList.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.addMinerList.length; i++) {
      if (!isValidAddress(this.addMinerList[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(buildTx_AddMiners(
        this.$store.state.address,
        this.$store.state.cashboxControllerContractAddress,
        this.addMinerList,
      ),
      0,
      "add miners",
    );
  }

  private sendTx_DelMiners(): void {
    if (this.delMinerList.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.delMinerList.length; i++) {
      if (!isValidAddress(this.delMinerList[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(buildTx_DelMiners(
        this.$store.state.address,
        this.$store.state.cashboxControllerContractAddress,
        this.delMinerList,
      ),
      0,
      "del miners",
    );
  }

  private call_GetMiners(): void {
    call_GetMiners(this.$store.state.cashboxChainAddress, this.$store.state.cashboxControllerContractAddress);
  }

  private sendTx_AddOperators(): void {
    if (this.addOperatorList.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.addOperatorList.length; i++) {
      if (!isValidAddress(this.addOperatorList[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(buildTx_AddOperators(
        this.$store.state.address,
        this.$store.state.cashboxControllerContractAddress,
        this.addOperatorList,
      ),
      0,
      "add operators",
    );
  }

  private sendTx_DelOperators(): void {
    if (this.delOperatorList.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.delOperatorList.length; i++) {
      if (!isValidAddress(this.delOperatorList[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(buildTx_DelOperators(
        this.$store.state.address,
        this.$store.state.cashboxControllerContractAddress,
        this.delOperatorList,
      ),
      0,
      "del operators",
    );
  }

  private call_GetOperators(): void {
    call_GetOperators(this.$store.state.cashboxChainAddress, this.$store.state.cashboxControllerContractAddress);
  }

  private sendTx_ReplaceAdmin_ControllerData(): void {
    if (!isValidAddress(this.oldAdminInControllerContract) || !isValidAddress(this.newAdminInControllerContract)) {
      this.$message.info("当前输入参数有误，请检查后重试");
      return;
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(buildTx_ReplaceAdmin_ControllerContract(
        this.$store.state.address,
        this.$store.state.cashboxControllerContractAddress,
        this.oldAdminInControllerContract,
        this.newAdminInControllerContract,
      ),
      0,
      "replace admin in controller contract",
    );
  }

  private call_GetAdmin_ControllerData(): void {
    call_GetAdmin_ControllerContract(this.$store.state.cashboxChainAddress, this.$store.state.cashboxControllerContractAddress);
  }

  private call_GetBatchTokens(): void {
    call_GetBatchTokens(this.$store.state.cashboxChainAddress,
      this.$store.state.cashboxDataContractAddress,
      this.batchCodeForTokens,
    );
  }

  private call_GetBatchTypes(): void {
    call_GetBatchTypes(this.$store.state.cashboxChainAddress,
      this.$store.state.cashboxDataContractAddress,
      this.batchCodeForTypes,
    );
  }

  private call_GetBatchOpenStatus(): void {
    call_GetBatchOpenStatus(this.$store.state.cashboxChainAddress,
      this.$store.state.cashboxDataContractAddress,
      this.batchCodeForOpenStatus,
    );
  }

  private call_GetTokenStatus(): void {
    if (!isValidTokenID(this.tokenIDForStatus)) {
      this.$message.info("当前输入参数有误，请检查后重试");
      return;
    }

    call_GetTokenStatus(this.$store.state.cashboxChainAddress,
      this.$store.state.cashboxDataContractAddress,
      this.tokenIDForStatus,
    );
  }

  private call_GetTokenOpenResult(): void {
    if (!isValidTokenID(this.tokenIDForOpenResult)) {
      this.$message.info("当前输入参数有误，请检查后重试");
      return;
    }

    call_GetTokenOpenResult(this.$store.state.cashboxChainAddress,
      this.$store.state.cashboxDataContractAddress,
      this.tokenIDForOpenResult,
    );
  }
}
</script>

<style lang="less">
.eth-manage {
  .em-info {
    text-align: left;
    font-size: 1.8rem;
    line-height: 4rem;
    padding-bottom: 10rem;

    .em-info-item {
      display: flex;
      line-height: 5rem;

      label {
        width: 20%;
      }

      .el-input {
        width: 60%;
      }
    }
  }

  .em-invokes {
    height: fit-content;
    text-align: right;

    .el-collapse-item__content {
      padding: 0;
    }

    .emi-item-input {
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

    .emi-item-button {
      padding-top: 3rem;
      padding-bottom: 3rem;
      padding-right: 10%;
      background-color: aliceblue;

      .left {
        padding-right: 50%;
      }
    }
  }
}
</style>
