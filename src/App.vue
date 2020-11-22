<template>
    <div id="app">
        <DynamicMenu
          :class="{ hidden: this.$route.name === 'Login' }"
          class="menu"
          :menu="menu"
          :isAdmin="isAdmin"
        ></DynamicMenu>
        <router-view class="router-view"></router-view>
    </div>
</template>

<script>
import DynamicMenu from "@/components/DynamicMenu";
import { setLocalStorage, getLocalStorage } from "@/Utils";
import getManagerMenu from "@/managerMenu";
import getUserMenu from "@/userMenu";
export default {
  name: "App",
  components: {
    DynamicMenu,
  },
  computed: {
    /*
     * 1. guest: 隐藏菜单
     * 2. admin: 显示管理菜单
     * 3. user: 显示用户菜单
     */
    isAdmin() {
      return this.$store.getters["App/getCurrentRole"] === "admin";
    },
    menu() {
      return this.$store.getters["App/getCurrentRole"] === "admin"
        ? getManagerMenu(this)
        : getUserMenu(this);
    },
  },
  created() {
    // 初始化为游客
    this.$store.commit("App/setCurrentRole", "guest");
    /*
     * 1. 页面刷新前保存store数据到localStroge
     * 2. 页面刷新完成从localStroge恢复store数据
     */
    if (getLocalStorage("store")) {
      this.$store.replaceState(
        Object.assign(
          {},
          this.$store.state,
          JSON.parse(getLocalStorage("store"))
        )
      );
    }
    window.addEventListener("beforeunload", () => {
      setLocalStorage("store", JSON.stringify(this.$store.state));
    });
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  user-select:none;
}
.hidden {
  display: none;
}
</style>
