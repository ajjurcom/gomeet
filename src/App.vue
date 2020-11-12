<template>
  <div id="app">
    <router-view class="router-view"></router-view>
  </div>
</template>

<script>
import {setLocalStorage, getLocalStorage} from '@/Utils';
export default {
  name: 'App',
  created() {
    /* 
     * 1. 页面刷新前保存store数据到localStroge
     * 2. 页面刷新完成从localStroge恢复store数据
     */
    if (getLocalStorage('store')) {
      this.$store.replaceState(Object.assign({}, this.$store.state, JSON.parse(getLocalStorage('store'))));
    }
    window.addEventListener("beforeunload", () => {
      setLocalStorage("store", JSON.stringify(this.$store.state));
    });
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
</style>
