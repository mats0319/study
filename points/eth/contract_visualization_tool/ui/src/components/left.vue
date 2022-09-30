<template>
  <div class="left">
    <el-menu unique-opened>
      <el-submenu index="ETH">
        <template slot="title">ETH合约</template>

        <el-menu-item inde="dependence" @click="linkTo('/eth/dependence')">部署依赖</el-menu-item>
        <el-menu-item inde="deploy" @click="linkTo('/eth/deploy')">部署合约</el-menu-item>
        <el-menu-item inde="manage" @click="linkTo('/eth/invoke')">调用合约函数</el-menu-item>
      </el-submenu>

      <el-submenu index="cashbox_controller">
        <template slot="title">Cashbox控制合约</template>

        <el-menu-item inde="dependence" @click="linkTo('/cashbox-controller/dependence')">部署依赖</el-menu-item>
        <el-menu-item inde="deploy-d" @click="linkTo('/cashbox-controller/deploy-d')">部署数据合约</el-menu-item>
        <el-menu-item inde="deploy-c" @click="linkTo('/cashbox-controller/deploy-c')">部署控制合约</el-menu-item>
        <el-menu-item inde="manage" @click="linkTo('/cashbox-controller/invoke')">调用合约函数</el-menu-item>
      </el-submenu>
    </el-menu>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Left extends Vue {
  private mounted() {
    // placeholder
  }

  private linkTo(path: string): void {
    const url = location.href.split(":9693/#");
    if (url.length < 2 || path !== url[1]) { // avoid navigate to current page
      this.$router.push({ path: path });
    }
  }
}
</script>

<style lang="less">
.left {
  padding-top: 8rem;
  text-align: left;
}
</style>
