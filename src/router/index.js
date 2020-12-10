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
import UserManager from 'pages/Manager/User/List';
import UserGroup from 'pages/Manager/UserGroup';
import Reserve from 'pages/Reserve';
import FastReserve from 'pages/FastReserve';
import ReserveManager from 'pages/ReserveManager';
import AppointmentManager from 'pages/Manager/Appointment';

Vue.use(Router);
Vue.use(viewDesign);

const routes = [
    {
        name: 'Login',
        path: '/login',
        meta: {
            roles: ['guest', 'user', 'admin']
        },
        component: LoginComponent
    },
    {
        name: 'Register',
        path: '/register',
        meta: {
            roles: ['guest', 'user', 'admin']
        },
        component: RegisterComponent
    },
    {
        name: 'CampusManager',
        path: '/back/campus/manager',
        meta: {
            roles: ['admin', 'root']
        },
        component: CampusManager
    },
    {
        name: 'CampusEdit',
        path: '/back/campus/edit',
        meta: {
            roles: ['admin', 'root']
        },
        component: CampusEdit
    },
    {
        name: 'CampusAdd',
        path: '/back/campus/add',
        meta: {
            roles: ['admin', 'root']
        },
        component: CampusAdd
    },
    {
        name: 'BuildingManager',
        path: '/back/building/manager',
        meta: {
            roles: ['admin', 'root']
        },
        component: BuildingManager
    },
    {
        name: 'BuildingEdit',
        path: '/back/building/edit/:id',
        meta: {
            roles: ['admin', 'root']
        },
        component: BuildingEdit
    },
    {
        name: 'BuildingAdd',
        path: '/back/building/add',
        meta: {
            roles: ['admin', 'root']
        },
        component: BuildingAdd
    },
    {
        name: 'MeetingAdd',
        path: '/back/meeting/add',
        meta: {
            roles: ['admin', 'root']
        },
        component: MeetingAdd
    },
    {
        name: 'MeetingManager',
        path: '/back/meeting/manager',
        meta: {
            roles: ['admin', 'root']
        },
        component: MeetingManager
    },
    {
        name: 'MeetingEdit',
        path: '/back/meeting/edit/:id',
        meta: {
            roles: ['admin', 'root']
        },
        component: MeetingEdit
    },
    {
        name: 'UserEdit',
        path: '/user/edit/:id',
        meta: {
            roles: ['user', 'admin', 'root']
        },
        component: UserEdit
    },
    {
        name: 'UserEditPwd',
        path: '/user/editpwd/:id',
        meta: {
            roles: ['user', 'admin', 'root']
        },
        component: UserEditPwd
    },
    {
        name: 'UserManager',
        path: '/back/user/manager',
        meta: {
            roles: ['admin', 'root']
        },
        component: UserManager,
        menuName: 'user'
    },
    {
        name: 'AppointmentManager',
        path: '/back/reserve',
        meta: {
            roles: ['admin', 'root']
        },
        component: AppointmentManager
    },
    {
        name: 'UserGroup',
        path: '/group',
        meta: {
            roles: ['user', 'admin', 'root']
        },
        component: UserGroup,
        menuName: 'usergroup'
    },
    {
        name: 'ReserveMeeting',
        path: '/',
        meta: {
            roles: ['user', 'admin', 'root']
        },
        component: Reserve
    },
    {
        name: 'FastReserve',
        path: '/fast',
        meta: {
            roles: ['user', 'admin', 'root']
        },
        component: FastReserve
    },
    {
        name: 'ReserveManager',
        path: '/my_meeting',
        meta: {
            roles: ['user', 'admin', 'root']
        },
        component: ReserveManager
    }
];

const router = new Router({
    mode: 'history',
    routes
});

router.beforeEach((to, from, next) => {
    try {
        const store = window.localStorage.getItem('store') || "{}";
        const storeJson = JSON.parse(store) || {};
        let role = storeJson.App ? storeJson.App.currentRole ? storeJson.App.currentRole : "guest" : "guest";
        if(!to.meta.roles.includes(role)) {
            next({path: "/login"});
        } else {
            viewDesign.LoadingBar.start();
            next();
        }
    } catch(e) {
        next({path: "/login"});
    }
});

router.afterEach(route => {
    viewDesign.LoadingBar.finish();
})

export default router;