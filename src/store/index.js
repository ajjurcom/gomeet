import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

// 创建Vuex实例
const store = new Vuex.Store({
    state: {
        demoName: ''
    },
    mutations: {
        changeDemoNameFunc(state, str) {
            state.demoName = str
        }
    },
    actions: {
        changeDemoName(context, str) {
            context.commit('changeDemoNameFunc', str)
        }
    }
});

export default store;