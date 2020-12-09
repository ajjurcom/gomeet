/**
 * @file vuex App
 * @author 陈铭涛
 */

export default {
    namespaced: true,
    state: {
        userID: -1,
        userName: 'guest',
        /* 
         * 1. guest: 隐藏菜单
         * 2. admin: 显示管理菜单
         * 3. user: 显示用户菜单
         */
        currentRole: 'guest',
        isRoot: false,
        userState: 'verify_user',
    },
    getters: {
        getUserID: state => state.userID,
        getUserName: state => state.userName,
        getUserIsRoot: state => state.isRoot,
        getCurrentRole: state => state.currentRole,
        getUserState: state => state.userState
    },
    mutations: {
        setUserID(state, value) {
            state.userID = value;
        },
        setUserName(state, value) {
            state.userName = value;
        },
        setUserIsRoot(state, value) {
            state.isRoot = value;
        },
        setCurrentRole(state, value) {
            state.currentRole = value;
        },
        setUserState(state, value) {
            state.userState = value;
        }
    },
    actions: {},
}