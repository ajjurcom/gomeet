import axios from 'axios';
import {showMessage} from '@/Utils';


export const variables = {
    getApiUrl: function () {
        return process.env.VUE_APP_API_URL;
    }
};

export const WebHttp = axios.create({
    baseURL: variables.getApiUrl(),
    withCredentials: false
});

WebHttp.interceptors.response.use(res => {
    const {success, data, msg} = res.data;
    if (!success) {
        showMessage('error', msg);
        return Promise.reject(msg);
    }
    return data;
}, error => {
    if (error.response && error.response.data && error.response.data.msg) {
        msg = error.response.data.msg;
    }
    showMessage('error', msg);
    return Promise.reject(error);
})
