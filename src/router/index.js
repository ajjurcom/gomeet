/**
 * @file 路由配置文件
 * @author chenmingtao
 */

import Vue from 'vue';
import Router from 'vue-router';
import LoginComponent from 'pages/Login'

Vue.use(Router);

const routes = [
    {
        name: 'Login',
        path: '/login',
        component: LoginComponent
    }
];

const router = new Router({
    mode: 'history',
    routes
});

export default router;