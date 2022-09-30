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
      <el-collapse-item title="设置数据合约持有的控制合约地址" name="setControllerContract">
        <div class="emi-item-input">
          <label>控制合约地址&#58;&nbsp;</label>
          <el-input v-model="$store.state.cashboxControllerContractAddress" placeholder="0x..." disabled />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="setCAddrInD">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="新增铸币者" name="addMiners">
        <div class="emi-item-input">
          <div v-for="(v, index) in addMinersAddrs" :key="index">
            <label>新的铸币者{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="addMinersAddrs[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="addMinersAddrs.push('')">增加铸币者</el-button>
            <el-button type="info" plain @click="addMinersAddrs.pop()">删除铸币者</el-button>
          </span>

          <el-button type="info" plain @click="addMinersList">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="删除铸币者" name="delMiners">
        <div class="emi-item-input">
          <div v-for="(v, index) in delMinersAddrs" :key="index">
            <label>铸币者{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="delMinersAddrs[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="delMinersAddrs.push('')">增加铸币者</el-button>
            <el-button type="info" plain @click="delMinersAddrs.pop()">删除铸币者</el-button>
          </span>

          <el-button type="info" plain @click="delMinersList">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询铸币者" name="getMiners">
        <div class="emi-item-button">
          <el-button type="info" plain @click="getMinersList">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="新增操作员" name="addOperator">
        <div class="emi-item-input">
          <div v-for="(v, index) in addOperatorsAddrs" :key="index">
            <label>新的操作员{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="addOperatorsAddrs[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="addOperatorsAddrs.push('')">增加操作员</el-button>
            <el-button type="info" plain @click="addOperatorsAddrs.pop()">删除操作员</el-button>
          </span>

          <el-button type="info" plain @click="addOperatorsList">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="删除操作员" name="delOperator">
        <div class="emi-item-input">
          <div v-for="(v, index) in delOperatorsAddrs" :key="index">
            <label>操作员{{ index + 1 }}&#58;&nbsp;</label>
            <el-input v-model="delOperatorsAddrs[index]" placeholder="0x..." />
          </div>
        </div>

        <div class="emi-item-button">
          <span class="left">
            <el-button type="info" plain @click="delOperatorsAddrs.push('')">增加操作员</el-button>
            <el-button type="info" plain @click="delOperatorsAddrs.pop()">删除操作员</el-button>
          </span>

          <el-button type="info" plain @click="delOperatorsList">发送交易</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询操作员" name="getOperators">
        <div class="emi-item-button">
          <el-button type="info" plain @click="getOperatorsList">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询指定批次的代币" name="getBatchTokens">
        <div class="emi-item-input">
          <label>批次码&#58;&nbsp;</label>
          <el-input v-model="batchCodeForTokens" placeholder="..." />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="getBatchTokensList">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询指定批次的代币类型" name="getBatchTypes">
        <div class="emi-item-input">
          <label>批次码&#58;&nbsp;</label>
          <el-input v-model="batchCodeForTypes" placeholder="..." />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="getBatchTypesList">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询指定批次的代币开启情况" name="getBatchOpenStatus">
        <div class="emi-item-input">
          <label>批次码&#58;&nbsp;</label>
          <el-input v-model="batchCodeForOpenStatus" placeholder="..." />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="getBatchOpenStatusList">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询指定代币的状态" name="getTokenStatus">
        <div class="emi-item-input">
          <label>token id&#58;&nbsp;</label>
          <el-input v-model="tokenIDForStatus" placeholder="0x..." />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="getTokenStatusByID">查询</el-button>
        </div>
      </el-collapse-item>

      <el-collapse-item title="查询指定代币的开启结果" name="getTokenType">
        <div class="emi-item-input">
          <label>token id&#58;&nbsp;</label>
          <el-input v-model="tokenIDForType" placeholder="0x..." />
        </div>

        <div class="emi-item-button">
          <el-button type="info" plain @click="getTokenTypeByID">查询</el-button>
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
import { isValidAddress, isValidTokenID } from "@/ts/kits";
import { ElLoadingComponent } from "element-ui/types/loading";
import {
  setControllerContract,
  addMiners,
  delMiners,
  getMiners,
  addOperators,
  delOperators,
  getOperators,
  getBatchTokens,
  getBatchTypes,
  getBatchOpenStatus,
  getTokenStatus,
  getTokenType,
  replaceAdmin,
} from "@/views/cashbox_controller/invoke";

@Component
export default class CashboxControllerManage extends Vue {
  private selectedItem = ""; // no use, for html
  private addMinersAddrs: Array<string> = [ "" ];
  private delMinersAddrs: Array<string> = [ "" ];
  private addOperatorsAddrs: Array<string> = [ "" ];
  private delOperatorsAddrs: Array<string> = [ "" ];
  private batchCodeForTokens = "";
  private batchCodeForTypes = "";
  private batchCodeForOpenStatus = "";
  private tokenIDForStatus = "";
  private tokenIDForType = "";
  private replaceAdminInput1 = "";
  private replaceAdminInput2 = "";

  private setCAddrTxHash = "";

  private loadingDom!: ElLoadingComponent;

  private mounted() {
    // placeholder
  }

  private setCAddrInD(): void {
    if (this.$store.state.cashboxDataContractAddress.length < 1) {
      this.$message.info("请先部署数据合约")
      return;
    }

    if (this.$store.state.cashboxControllerContractAddress.length < 1) {
      this.$message.info("请先部署控制合约")
      return;
    }

    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包执行操作",
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    setControllerContract(this.$store.state.address,
      this.$store.state.cashboxDataContractAddress,
      this.$store.state.cashboxControllerContractAddress,
    )
      .then((data: any) => {
        this.$message.success("调用成功");
        console.log("> set C contract address success.", data)

        this.setCAddrTxHash = data;
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> set C contract address failed, err: ", err);
      })
      .finally(() => {
        this.loadingDom.close();
      })
  }

  private addMinersList(): void {
    if (this.addMinersAddrs.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.addMinersAddrs.length; i++) {
      if (!isValidAddress(this.addMinersAddrs[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包执行操作",
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    addMiners(this.$store.state.address, this.$store.state.cashboxControllerContractAddress, this.addMinersAddrs)
      .then((data: any) => {
        this.$message.success("调用成功")
        console.log("> add miners success.", data)
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> add miners failed, err: ", err);
      })
      .finally(() => {
        this.loadingDom.close();
      })
  }

  private delMinersList(): void {
    if (this.delMinersAddrs.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.delMinersAddrs.length; i++) {
      if (!isValidAddress(this.delMinersAddrs[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包执行操作",
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    delMiners(this.$store.state.address, this.$store.state.cashboxControllerContractAddress, this.delMinersAddrs)
      .then((data: any) => {
        this.$message.success("调用成功")
        console.log("> del miners success.", data)
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> del miners failed, err: ", err);
      })
      .finally(() => {
        this.loadingDom.close();
      })
  }

  private getMinersList(): void {
    getMiners(this.$store.state.cashboxChainAddress, this.$store.state.cashboxControllerContractAddress)
      .then((address: any) => {
        this.$message.success("当前铸币者地址列表：" + address);
        console.log("> get miners success, address：", address);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get miners failed, error: ", err);
      });
  }

  private addOperatorsList(): void {
    if (this.addOperatorsAddrs.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.addOperatorsAddrs.length; i++) {
      if (!isValidAddress(this.addOperatorsAddrs[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包执行操作",
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    addOperators(this.$store.state.address, this.$store.state.cashboxControllerContractAddress, this.addOperatorsAddrs)
      .then((data: any) => {
        this.$message.success("调用成功")
        console.log("> add operators success.", data)
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> add operators failed, err: ", err);
      })
      .finally(() => {
        this.loadingDom.close();
      })
  }

  private delOperatorsList(): void {
    if (this.delOperatorsAddrs.length < 1) {
      this.$message.info("请至少输入一个地址")
      return;
    }

    for (let i = 0; i < this.delOperatorsAddrs.length; i++) {
      if (!isValidAddress(this.delOperatorsAddrs[i])) {
        this.$message.info("当前输入参数有误，请检查后重试: " + (i + 1));
        return;
      }
    }

    this.loadingDom = this.$loading({
      lock: true,
      text: "请前往钱包执行操作",
      background: "rgba(0, 0, 0, 0.7)" // background color
    });

    delOperators(this.$store.state.address, this.$store.state.cashboxControllerContractAddress, this.delOperatorsAddrs)
      .then((data: any) => {
        this.$message.success("调用成功")
        console.log("> del operators success.", data)
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> del operators failed, err: ", err);
      })
      .finally(() => {
        this.loadingDom.close();
      })
  }

  private getOperatorsList(): void {
    getOperators(this.$store.state.cashboxChainAddress, this.$store.state.cashboxControllerContractAddress)
      .then((address: any) => {
        this.$message.success("当前操作员地址列表：" + address);
        console.log("> get operators success, address：", address);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get operators failed, error: ", err);
      });
  }

  private getBatchTokensList(): void {
    getBatchTokens(this.$store.state.cashboxChainAddress,
      this.$store.state.cashboxControllerContractAddress,
      this.batchCodeForTokens,
    )
      .then((tokenIDs: any) => {
        this.$message.success("指定批次代币列表：" + tokenIDs);
        console.log("> get batch tokens success, address：", tokenIDs);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get batch tokens failed, error: ", err);
      });
  }

  private getBatchTypesList(): void {
    getBatchTypes(this.$store.state.cashboxChainAddress,
      this.$store.state.cashboxControllerContractAddress,
      this.batchCodeForTypes,
    )
      .then((tokenTypes: any) => {
        this.$message.success("指定批次代币类型：" + tokenTypes);
        console.log("> get batch types success, address：", tokenTypes);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get batch types failed, error: ", err);
      });
  }

  private getBatchOpenStatusList(): void {
    getBatchOpenStatus(this.$store.state.cashboxChainAddress,
      this.$store.state.cashboxControllerContractAddress,
      this.batchCodeForOpenStatus,
    )
      .then((tokenOpenStatus: any) => {
        this.$message.success("指定批次代币开启情况：" + tokenOpenStatus);
        console.log("> get batch open status success, address：", tokenOpenStatus);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get batch open status failed, error: ", err);
      });
  }

  private getTokenStatusByID(): void {
    if (!isValidTokenID(this.tokenIDForStatus)) {
      this.$message.info("当前输入参数有误，请检查后重试");
      return;
    }

    getTokenStatus(this.$store.state.cashboxChainAddress,
      this.$store.state.cashboxControllerContractAddress,
      this.tokenIDForStatus,
    )
      .then((tokenStatus: any) => {
        this.$message.success("指定代币状态：" + tokenStatus);
        console.log("> get token status success, address：", tokenStatus);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get token status failed, error: ", err);
      });
  }

  private getTokenTypeByID(): void {
    if (!isValidTokenID(this.tokenIDForType)) {
      this.$message.info("当前输入参数有误，请检查后重试");
      return;
    }

    getTokenType(this.$store.state.cashboxChainAddress,
      this.$store.state.cashboxControllerContractAddress,
      this.tokenIDForType,
    )
      .then((tokenType: any) => {
        this.$message.success("指定代币开启结果：" + tokenType);
        console.log("> get token type success, address：", tokenType);
      })
      .catch((err: any) => {
        this.$message.error(err);
        console.log("> get token type failed, error: ", err);
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

    replaceAdmin(this.$store.state.address,
      this.$store.state.cashboxControllerContractAddress,
      this.replaceAdminInput1,
      this.replaceAdminInput2,
    )
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
