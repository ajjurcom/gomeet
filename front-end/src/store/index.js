/**
 * @file Vuex配置文件
 * @author chenmingtao
 */

import Vue from 'vue';
import Vuex from 'vuex';
import App from './modules/App';

Vue.use(Vuex);

export default new Vuex.Store({
    modules: {
        App
    }
});
