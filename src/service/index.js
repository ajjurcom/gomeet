import axios from 'axios';
import {showMessage, getLocalStorage} from '@/Utils';
import Vue from 'vue';
import viewDesign from 'view-design';
import vue from "../main";

Vue.use(viewDesign);

export const variables = {
    getApiUrl: function () {
        return process.env.VUE_APP_API_URL;
    }
};

export const WebHttp = axios.create({
    baseURL: variables.getApiUrl(),
    headers: { 'loginToken': getLocalStorage('loginToken') === null ? '' : getLocalStorage('loginToken')},
    withCredentials: false
});

WebHttp.defaults.headers.post['Content-Type'] = 'application/json;charset=UTF-8';


// 请求前
WebHttp.interceptors.response.use(config => {
    viewDesign.LoadingBar.start();
    return config;
}, error => {
    viewDesign.LoadingBar.error();
    return Promise.reject(error)
})


// 相应处理
WebHttp.interceptors.response.use(res => {
    const {success, data, msg} = res.data;
    if (!success) {
        viewDesign.LoadingBar.error();
        showMessage('info', msg);
        return Promise.reject(msg);
    }
    viewDesign.LoadingBar.finish();
    return data;
}, error => {
    viewDesign.LoadingBar.error();
    let msg = '';
    if (error.response.status === 401) {
        msg = "暂无权限, 即将跳转登录页面!";
        setTimeout(() => {
            // 跳转页面
            vue.$router.push({
                name: "Login"
            });
        }, 2000)
    }
    else if (error.response && error.response.data && error.response.data.msg) {
        msg = error.response.data.msg;
    }
    showMessage('error', msg);
    return Promise.reject(error);
})
