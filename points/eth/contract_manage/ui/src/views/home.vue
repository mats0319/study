<template>
  <div class="home">
    <el-container>
      <el-header class="h-header" height="10rem">
        <top />
      </el-header>

      <el-container v-show="connector.connected" class="h-content">
        <el-aside class="h-aside" width="35rem"><left /></el-aside>

        <el-main class="h-main">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Top from "@/components/top.vue";
import Left from "@/components/left.vue";
import { connector } from '@/ts/wallet_connect';

@Component({
  components: { Top, Left },
})
export default class Home extends Vue {
  // const, for html
  private connector = connector;

  private mounted() {
    // placeholder
  }
}
</script>

<style lang="less">
.home {
  display: flex;
  height: 100%;

  .h-header {
    background-color: lightgray;
  }

  .el-header,.el-main {
    padding: 0;
  }

  .h-content {
    height: calc(100% - 10rem);

    .h-aside {
      border-right-color: darkgray;
      border-right-style: solid;
      border-right-width: 1px;
    }

    .h-main {
      padding: 8rem 10rem;
    }
  }
}
</style>
