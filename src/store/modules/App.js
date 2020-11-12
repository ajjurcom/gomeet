/**
 * @file vuex App
 * @author chenruofeng
 */

export default {
    namespaced: true,
    state: {
        userID: -1,
        userName: 'Guest'
    },
    getters: {
        getUserID: state => state.userID,
        getUserName: state => state.userName,
    },
    mutations: {
        setUserID(state, value) {
            state.userID = value;
        },
        setUserName(state, value) {
            state.userName = value;
        }
    },
    actions: {},
}