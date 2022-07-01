<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class App extends Vue {
  private created() {
    if (sessionStorage.getItem("vuex")) {
      this.$store.replaceState(
        Object.assign(
          {},
          this.$store.state,
          JSON.parse(sessionStorage.getItem("vuex") as string)
        )
      );

      sessionStorage.removeItem("vuex");
    }
  }

  private mounted() {
    window.addEventListener("beforeunload", () => {
      sessionStorage.setItem("vuex", JSON.stringify(this.$store.state));
    });
  }
}
</script>

<style lang="less">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  height: 100%;
  width: 100%;
}

html {
  height: 100%;
  width: 100%;
  font-size: 62.5%;
}

body {
  height: 100%;
  width: 100%;
  font-size: 62.5%;
  margin: 0;
  padding: 0;
}
</style>
