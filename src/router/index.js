/**
 * @file 路由配置文件
 * @author chenmingtao
 */

import Vue from 'vue';
import Router from 'vue-router';
import viewDesign from 'view-design';
import LoginComponent from 'pages/Login';
import RegisterComponent from 'pages/Register';
import CampusManager from 'pages/Manager/CampusManager';
import CampusEdit from 'pages/Manager/CampusEdit';
import CampusAdd from 'pages/Manager/CampusAdd';

Vue.use(Router);
Vue.use(viewDesign);

const routes = [
    {
        name: 'Login',
        path: '/login',
        component: LoginComponent
    },
    {
        name: 'Register',
        path: '/register',
        component: RegisterComponent
    },
    {
        name: 'CampusManager',
        path: '/back/campus/manager',
        component: CampusManager
    },
    {
        name: 'CampusEdit',
        path: '/back/campus/edit',
        component: CampusEdit
    },
    {
        name: 'CampusAdd',
        path: '/back/campus/add',
        component: CampusAdd
    }
];

const router = new Router({
    mode: 'history',
    routes
});

router.beforeEach((to, from, next) => {
    viewDesign.LoadingBar.start();
    next();
});

router.afterEach(route => {
    viewDesign.LoadingBar.finish();
})

export default router;