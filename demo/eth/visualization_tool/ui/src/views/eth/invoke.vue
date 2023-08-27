<template>
  <div class="eth-manage">
    <div class="em-info">
      <div class="em-info-item">
        <label>链地址&#58;&nbsp;</label>
        <el-input v-model="$store.state.ethChainAddress"
                  placeholder="https://ropsten.infura.io/v3/63ed578c193242289ab781596516987a" />
      </div>

      <div class="em-info-item">
        <label>合约地址&#58;&nbsp;</label>
        <el-input v-model="$store.state.ethContractAddress" placeholder="0x..." />
      </div>
    </div>

    <el-collapse
      v-show="$store.state.ethChainAddress.length > 0 && $store.state.ethContractAddress.length > 0"
      class="em-invokes"
      v-model="selectedItem"
      accordion
    >
      <el-collapse-item title="替换ETH收款地址" name="replaceETHReceiver">
        <div class="emi-item-input">
          <label>新的ETH收款地址&#58;&nbsp;</label>
          <el-input v-model="newETHReceiver" placeholder="0x..." />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="sendTx_ReplaceETHReceiver">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询ETH收款地址" name="getETHReceiver">
        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetETHReceiver">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="新增ACash提现操作员" name="addACashOperator">
        <div class="emi-item-input">
          <div v-for="(v, index) in addACashOperatorList" :key="index">
            <label>新的操作员{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="addACashOperatorList[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="addACashOperatorList.push('')">增加操作员</el-button>
            <el-button type="info" plain @click="addACashOperatorList.pop()">删除操作员</el-button>
          </span>

          <el-button type="info" plain @click="sendTx_AddACashOperators">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="删除ACash提现操作员" name="delACashOperator">
        <div class="emi-item-input">
          <div v-for="(v, index) in delACashOperatorList" :key="index">
            <label>操作员{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="delACashOperatorList[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="delACashOperatorList.push('')">增加操作员</el-button>
            <el-button type="info" plain @click="delACashOperatorList.pop()">删除操作员</el-button>
          </span>

          <el-button type="info" plain @click="sendTx_DelACashOperators">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询ACash提现操作员" name="getACashOperators">
        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetACashOperators">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="新增佣金提现操作员" name="addBrokerageOperator">
        <div class="emi-item-input">
          <div v-for="(v, index) in addBrokerageOperatorList" :key="index">
            <label>新的操作员{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="addBrokerageOperatorList[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="addBrokerageOperatorList.push('')">增加操作员</el-button>
            <el-button type="info" plain @click="addBrokerageOperatorList.pop()">删除操作员</el-button>
          </span>

          <el-button type="info" plain @click="sendTx_AddBrokerageOperators">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="删除佣金提现操作员" name="delBrokerageOperator">
        <div class="emi-item-input">
          <div v-for="(v, index) in delBrokerageOperatorList" :key="index">
            <label>操作员{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="delBrokerageOperatorList[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="delBrokerageOperatorList.push('')">增加操作员</el-button>
            <el-button type="info" plain @click="delBrokerageOperatorList.pop()">删除操作员</el-button>
          </span>

          <el-button type="info" plain @click="sendTx_DelBrokerageOperators">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询佣金提现操作员" name="getBrokerageOperators">
        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetBrokerageOperators">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="替换管理员" name="replaceAdmin">
        <div class="emi-item-input">
          <div>
            <label>老管理员&#58;&nbsp;</label>
            <el-input v-model="oldAdmin" placeholder="0x..." />
          </div>

          <div>
            <label>新管理员&#58;&nbsp;</label>
            <el-input v-model="newAdmin" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="sendTx_ReplaceAdmin">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询管理员" name="getAdmin">
        <div class="emi-item-button">
          <el-button type="info" plain @click="call_GetAdmin">查询</el-button>
        </div>
      </el-collapse-item>
    </el-collapse>

    <send-transaction ref="sendTx" />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import {
  buildTx_ReplaceETHReceiver,
  buildTx_AddACashOperators,
  buildTx_DelACashOperators,
  buildTx_AddBrokerageOperators,
  buildTx_DelBrokerageOperators,
  buildTx_ReplaceAdmin,
} from "@/views/eth/build_tx";
import { getETHReceiver, getACashOperators, getBrokerageOperators, getAdmin } from "@/views/eth/call"
import { isValidAddress } from "@/ts/utils";
import SendTransaction from "@/components/send_transaction.vue";

@Component({
  components: { SendTransaction }
})
export default class ETHInvoke extends Vue {
  private selectedItem = ""; // no use

  private newETHReceiver = "";
  private addACashOperatorList: Array<string> = [ "" ];
  private delACashOperatorList: Array<string> = [ "" ];
  private addBrokerageOperatorList: Array<string> = [ "" ];
  private delBrokerageOperatorList: Array<string> = [ "" ];
  private oldAdmin = "";
  private newAdmin = "";

  private mounted() {
    // placeholder
  }

  // Note: idea will mixture this func and same name method in 'contractIns.methods.[method]', so add meaningless prefix to distinguish

  private sendTx_ReplaceETHReceiver(): void {
    if (!isValidAddress(this.newETHReceiver)) {
      this.$message.info("当前输入参数有误，请检查后重试");
      return;
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(buildTx_ReplaceETHReceiver(
        this.$store.state.address,
        this.$store.state.ethContractAddress,
        this.newETHReceiver,
      ),
      0,
      "replace eth receiver",
    );
  }

  private call_GetETHReceiver(): void {
    getETHReceiver(this.$store.state.ethChainAddress, this.$store.state.ethContractAddress);
  }

  private sendTx_AddACashOperators(): void {
    if (this.addACashOperatorList.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.addACashOperatorList.length; i++) {
      if (!isValidAddress(this.addACashOperatorList[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(buildTx_AddACashOperators(
        this.$store.state.address,
        this.$store.state.ethContractAddress,
        this.addACashOperatorList,
      ),
      0,
      "add ACash operators",
    );
  }

  private sendTx_DelACashOperators(): void {
    if (this.delACashOperatorList.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.delACashOperatorList.length; i++) {
      if (!isValidAddress(this.delACashOperatorList[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(buildTx_DelACashOperators(
        this.$store.state.address,
        this.$store.state.ethContractAddress,
        this.delACashOperatorList,
      ),
      0,
      "del ACash operators",
    );
  }

  private call_GetACashOperators(): void {
    getACashOperators(this.$store.state.ethChainAddress, this.$store.state.ethContractAddress);
  }

  private sendTx_AddBrokerageOperators(): void {
    if (this.addBrokerageOperatorList.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.addBrokerageOperatorList.length; i++) {
      if (!isValidAddress(this.addBrokerageOperatorList[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(buildTx_AddBrokerageOperators(
        this.$store.state.address,
        this.$store.state.ethContractAddress,
        this.addBrokerageOperatorList,
      ),
      0,
      "add brokerage operators",
    );
  }

  private sendTx_DelBrokerageOperators(): void {
    if (this.delBrokerageOperatorList.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.delBrokerageOperatorList.length; i++) {
      if (!isValidAddress(this.delBrokerageOperatorList[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(buildTx_DelBrokerageOperators(
        this.$store.state.address,
        this.$store.state.ethContractAddress,
        this.delBrokerageOperatorList,
      ),
      0,
      "del brokerage operators",
    );
  }

  private call_GetBrokerageOperators(): void {
    getBrokerageOperators(this.$store.state.ethChainAddress, this.$store.state.ethContractAddress);
  }

  private sendTx_ReplaceAdmin(): void {
    if (!isValidAddress(this.oldAdmin) || !isValidAddress(this.newAdmin)) {
      this.$message.info("当前输入参数有误，请检查后重试");
      return;
    }

    //@ts-ignore
    this.$refs.sendTx.sendTx(buildTx_ReplaceAdmin(
        this.$store.state.address,
        this.$store.state.ethContractAddress,
        this.oldAdmin,
        this.newAdmin,
      ),
      0,
      "replace admin",
    );
  }

  private call_GetAdmin(): void {
    getAdmin(this.$store.state.ethChainAddress, this.$store.state.ethContractAddress);
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
