<template>
  <div class="eth-manage">
    <div class="em-info">
      <div class="em-info-item">
        <label>链地址&#58;&nbsp;</label>
        <el-input v-model="$store.state.ethChainAddress" placeholder="https://ropsten.infura.io/v3/63ed578c193242289ab781596516987a" />
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
          <el-input v-model="replaceETHReceiverInput" placeholder="0x..." />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="replaceETHReceiverAddress">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询ETH收款地址" name="getETHReceiver">
        <div class="emi-item-button">
          <el-button type="info" plain @click="getETHReceiverAddress">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="新增ACash提现操作员" name="addACashOperator">
        <div class="emi-item-input">
          <div v-for="(v, index) in addACashAddrs" :key="index">
            <label>新的操作员{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="addACashAddrs[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="addACashAddrs.push('')">增加操作员</el-button>
            <el-button type="info" plain @click="addACashAddrs.pop()">删除操作员</el-button>
          </span>

          <el-button type="info" plain @click="addWithdrawACashOperators">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="删除ACash提现操作员" name="delACashOperator">
        <div class="emi-item-input">
          <div v-for="(v, index) in delACashAddrs" :key="index">
            <label>操作员{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="delACashAddrs[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="delACashAddrs.push('')">增加操作员</el-button>
            <el-button type="info" plain @click="delACashAddrs.pop()">删除操作员</el-button>
          </span>

          <el-button type="info" plain @click="delWithdrawACashOperators">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询ACash提现操作员" name="getACashOperators">
        <div class="emi-item-button">
          <el-button type="info" plain @click="getWithdrawACashOperators">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="新增佣金提现操作员" name="addBrokerageOperator">
        <div class="emi-item-input">
          <div v-for="(v, index) in addBrokerageAddrs" :key="index">
            <label>新的操作员{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="addBrokerageAddrs[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="addBrokerageAddrs.push('')">增加操作员</el-button>
            <el-button type="info" plain @click="addBrokerageAddrs.pop()">删除操作员</el-button>
          </span>

          <el-button type="info" plain @click="addWithdrawBrokerageOperators">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="删除佣金提现操作员" name="delBrokerageOperator">
        <div class="emi-item-input">
          <div v-for="(v, index) in delBrokerageAddrs" :key="index">
            <label>操作员{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="delBrokerageAddrs[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="delBrokerageAddrs.push('')">增加操作员</el-button>
            <el-button type="info" plain @click="delBrokerageAddrs.pop()">删除操作员</el-button>
          </span>

          <el-button type="info" plain @click="delWithdrawBrokerageOperators">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询佣金提现操作员" name="getBrokerageOperators">
        <div class="emi-item-button">
          <el-button type="info" plain @click="getWithdrawBrokerageOperators">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="替换管理员" name="replaceAdmin">
        <div class="emi-item-input">
          <div>
            <label>老管理员&#58;&nbsp;</label>
            <el-input v-model="replaceAdminInput1" placeholder="0x..." />
          </div>

          <div>
            <label>新管理员&#58;&nbsp;</label>
            <el-input v-model="replaceAdminInput2" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="replaceSingleAdmin">发送交易</el-button>
        </div>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import {
  addACashOperators, addBrokerageOperators,
  delACashOperators, delBrokerageOperators,
  getACashOperators, getBrokerageOperators,
  getETHReceiver, replaceAdmin,
  replaceETHReceiver
} from "@/views/eth/invoke";
import { isValidAddress } from "@/ts/kits";
import { ElLoadingComponent } from "element-ui/types/loading";

@Component
export default class ETHManage extends Vue {
  private selectedItem = ""; // no use, for html
  private replaceETHReceiverInput = "";
  private addACashAddrs: Array<string> = [ "" ];
  private delACashAddrs: Array<string> = [ "" ];
  private addBrokerageAddrs: Array<string> = [ "" ];
  private delBrokerageAddrs: Array<string> = [ "" ];
  private replaceAdminInput1 = "";
  private replaceAdminInput2 = "";

  private loadingDom!: ElLoadingComponent;

  private mounted() {
    // placeholder
  }

  private replaceETHReceiverAddress(): void {
    if (!isValidAddress(this.replaceETHReceiverInput)) {
      this.$message.info("当前输入参数有误，请检查后重试");
      return;
    }

    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包执行操作",
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    replaceETHReceiver(this.$store.state.address, this.$store.state.ethContractAddress, this.replaceETHReceiverInput)
      .then((data: any) => {
        this.$message.success("调用成功")
        console.log("> replace eth receiver success.", data)
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> replace eth receiver failed, error: ", err)
      })
      .finally(() => {
        this.loadingDom.close();
      });
  }

  private getETHReceiverAddress(): void {
    getETHReceiver(this.$store.state.ethChainAddress, this.$store.state.ethContractAddress)
      .then((address: any) => {
        this.$message.success("当前ETH收款地址为：" + address);
        console.log("> get eth receiver success, address：", address);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get eth receiver failed, error: ", err);
      });
  }

  private addWithdrawACashOperators(): void {
    if (this.addACashAddrs.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.addACashAddrs.length; i++) {
      if (!isValidAddress(this.addACashAddrs[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包执行操作",
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    addACashOperators(this.$store.state.address, this.$store.state.ethContractAddress, this.addACashAddrs)
      .then((data: any) => {
        this.$message.success("调用成功")
        console.log("> add ACash operators success.", data)
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> add ACash operators failed, err: ", err);
      })
      .finally(() => {
        this.loadingDom.close();
      })
  }

  private delWithdrawACashOperators(): void {
    if (this.delACashAddrs.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.delACashAddrs.length; i++) {
      if (!isValidAddress(this.delACashAddrs[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包执行操作",
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    delACashOperators(this.$store.state.address, this.$store.state.ethContractAddress, this.delACashAddrs)
      .then((data: any) => {
        this.$message.success("调用成功")
        console.log("> del ACash operators success.", data)
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> del ACash operators failed, err: ", err);
      })
      .finally(() => {
        this.loadingDom.close();
      })
  }

  private getWithdrawACashOperators(): void {
    getACashOperators(this.$store.state.ethChainAddress, this.$store.state.ethContractAddress)
      .then((address: any) => {
        this.$message.success("当前ACash提现操作员地址列表：" + address);
        console.log("> get ACash operators success, address：", address);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get ACash operators failed, error: ", err);
      });
  }

  private addWithdrawBrokerageOperators(): void {
    if (this.addBrokerageAddrs.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.addBrokerageAddrs.length; i++) {
      if (!isValidAddress(this.addBrokerageAddrs[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包执行操作",
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    addBrokerageOperators(this.$store.state.address, this.$store.state.ethContractAddress, this.addBrokerageAddrs)
      .then((data: any) => {
        this.$message.success("调用成功")
        console.log("> add brokerage operators success.", data)
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> add brokerage operators failed, err: ", err);
      })
      .finally(() => {
        this.loadingDom.close();
      })
  }

  private delWithdrawBrokerageOperators(): void {
    if (this.delBrokerageAddrs.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.delBrokerageAddrs.length; i++) {
      if (!isValidAddress(this.delBrokerageAddrs[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包执行操作",
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    delBrokerageOperators(this.$store.state.address, this.$store.state.ethContractAddress, this.delBrokerageAddrs)
      .then((data: any) => {
        this.$message.success("调用成功")
        console.log("> del brokerage operators success.", data)
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> del brokerage operators failed, err: ", err);
      })
      .finally(() => {
        this.loadingDom.close();
      })
  }

  private getWithdrawBrokerageOperators(): void {
    getBrokerageOperators(this.$store.state.ethChainAddress, this.$store.state.ethContractAddress)
      .then((address: any) => {
        this.$message.success("当前佣金提现操作员地址列表：" + address);
        console.log("> get brokerage operators success, address：", address);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get brokerage operators failed, error: ", err);
      });
  }

  private replaceSingleAdmin(): void {
    if (!isValidAddress(this.replaceAdminInput1) || !isValidAddress(this.replaceAdminInput2)) {
      this.$message.info("当前输入参数有误，请检查后重试");
      return;
    }

    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包执行操作",
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    replaceAdmin(this.$store.state.address, this.$store.state.ethContractAddress, this.replaceAdminInput1, this.replaceAdminInput2)
      .then((data: any) => {
        this.$message.success("调用成功")
        console.log("> replace admin success.", data)
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> replace admin failed, error: ", err)
      })
      .finally(() => {
        this.loadingDom.close();
      });
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
