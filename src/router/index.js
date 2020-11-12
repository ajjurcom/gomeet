/**
 * @file 路由配置文件
 * @author chenmingtao
 */

import Vue from 'vue';
import Router from 'vue-router';
import viewDesign from 'view-design';
import LoginComponent from 'pages/Login';
import RegisterComponent from 'pages/Register';
import CampusManager from 'pages/Manager/Campus/List';
import CampusEdit from 'pages/Manager/Campus/Edit';
import CampusAdd from 'pages/Manager/Campus/Add';
import BuildingManager from 'pages/Manager/Building/List';
import BuildingEdit from 'pages/Manager/Building/Edit';
import BuildingAdd from 'pages/Manager/Building/Add';
import MeetingAdd from 'pages/Manager/Meeting/Add';
import MeetingManager from 'pages/Manager/Meeting/List';
import MeetingEdit from 'pages/Manager/Meeting/Edit';
import UserEdit from 'pages/User/Edit';
import UserEditPwd from 'pages/User/EditPwd';

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
    },
    {
        name: 'BuildingManager',
        path: '/back/building/manager',
        component: BuildingManager
    },
    {
        name: 'BuildingEdit',
        path: '/back/building/edit/:id',
        component: BuildingEdit
    },
    {
        name: 'BuildingAdd',
        path: '/back/building/add',
        component: BuildingAdd
    },
    {
        name: 'MeetingAdd',
        path: '/back/meeting/add',
        component: MeetingAdd
    },
    {
        name: 'MeetingManager',
        path: '/back/meeting/manager',
        component: MeetingManager
    },
    {
        name: 'MeetingEdit',
        path: '/back/meeting/edit/:id',
        component: MeetingEdit
    },
    {
        name: 'UserEdit',
        path: '/user/edit/:id',
        component: UserEdit
    },
    {
        name: 'UserEditPwd',
        path: '/user/editpwd/:id',
        component: UserEditPwd
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