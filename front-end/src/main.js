import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store/index';

import global from './Globals';

Vue.use(global);

Vue.config.productionTip = false;

export default new Vue({
  router,
  store,
  render: function (h) { return h(App) },
}).$mount('#app');


