/**
 * @file vuex App
 * @author chenruofeng
 */

export default {
    namespaced: true,
    state: {
        userID: -1,
        userName: 'Guest',
        isRoot: false
    },
    getters: {
        getUserID: state => state.userID,
        getUserName: state => state.userName,
        getUserIsRoot: state => state.isRoot
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
        }
    },
    actions: {},
}