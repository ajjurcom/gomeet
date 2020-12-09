<template>
    <div class="container">
        <div class="back-image"></div>
        <div class="main">
            <div class="main-title">
                <div class="title-name">GoMeet</div>
                <div class="title-info">佛科院会议室预定系统</div>
            </div>
            <div class="login-container">
                <div class="login-role">
                    <div class="login-user" :class="{'active-role': currentRole === 'user'}" @click="changeRole('user')">会议预定</div>
                    <div class="login-admin" :class="{'active-role': currentRole === 'admin'}" @click="changeRole('admin')">后台管理</div>
                </div>
                <div class="login-form">
                    <input class="form-input" v-model="value" :placeholder="loginWay" style="width: 100%" />
                    <input type="password" class="form-input" v-model="password" placeholder="密码" style="width: 100%" />
                </div>
                <div class="login-bottom">
                    <div class="login-way" @click="changeWay">{{valueTips}}</div>
                    <router-link class="login-register" to="Register">注册</router-link>
                </div>
                <Button class="login-button" :loading="loading" @click="Login">登录</Button>
            </div>
        </div>
    </div>
</template>

<style lang="less" scoped>
.container {
    .back-image {
        position: fixed;
        top: 0;
        left: 0;
        width:100%;
        height:100%;
        min-width: 1000px;
        z-index:-10;
        zoom: 1;
        background-repeat: no-repeat;
        background-size: cover;
        background-position: center 0;
        background: url("http://static.mittacy.com/blog/login.JPG");
        background-size: 100% 100%;
    }
    .main {
        margin: 0 auto;
        margin-top: 50px;
        width: 400px;
        .main-title {
            text-align: center;
            margin-bottom: 24px;
            .title-name {
                font-size: 36px;
                font-weight: 700;
                color: #0983ff;

            }
            .title-info {
                font-size: 10px;
                color: #47494a;
            }
        }
        .login-container {
            width: 400px;
            padding: 0 24px;
            background-color: #fff;
            opacity: 0.8;
            .login-role {
                display: flex;
                flex-direction: row;
                .active-role {
                    border-bottom: 3px solid #0084ff;
                }
                .login-user,
                .login-admin {
                    font-size: 16px;
                    font-weight: 600;
                    line-height: 60px;
                    height: 60px;
                    margin-right: 24px;
                    color: #444;
                    cursor: pointer;
                }
            }
            .login-form {
                width: 100%;
                .form-input {
                    width: 100%;
                    height: 48px;
                    margin-bottom: 10px;
                    margin-top: 24px;
                    border: none;
                    border-bottom: 1px solid #ebebeb;
                    outline:none;
                    padding-left: 3px;
                }
            }
            .login-bottom {
                width: 100%;
                display: flex;
                flex-direction: row;
                justify-content: space-between;
                .login-way {
                    cursor: pointer;
                    margin-left: 2px;
                    color: #175199;
                }
                .login-register {
                    cursor: pointer;
                    margin-right: 2px;
                    color: #8590a6;
                }
                .login-way:hover {
                    color: #76839b;
                }
                .login-register:hover {
                    color: #575a61;
                }
            }
            .login-button {
                width: 100%;
                margin: 20px 0px;
                color: #fff;
                background-color: #0084ff;
            }
            .login-button:hover {
                box-sizing: border-box;
                border-color: #0077e6;
                background-color: #0077e6;
            }
        }
    }
}
</style>

<script>
import {showMessage, setLocalStorage} from '@/Utils';
export default {
    name: "Login",
    data () {
        return {
            loading: false,
            currentRole: 'user',    // admin / user
            value: '',
            loginWay: '学号',
            password: '',
            valueTips: '使用手机登录'
        };
    },
    methods: {
        changeRole(role) {
            this.currentRole = role
        },
        changeWay() {
            this.valueTips = this.valueTips === '使用手机登录' ? '使用学号登录' : '使用手机登录';
            this.loginWay = this.loginWay === '学号' ? '手机号' : '学号';
            this.value = '';
            this.password = '';
        },
        Login() {
            this.loading = true;
            // 1. 检查输入值是否满足
            if (this.value === "" || this.password === "") {
                this.loading = false;
                this.$Message.error('请输入账号和密码');
                return 
            }
            let postUser = {
                sno: '',
                phone: '',
                password: this.password,
                is_admin: false
            };
            // 2. 判断角色
            if (this.currentRole === 'admin') {
                postUser.is_admin = true;
            }
            // 3. 判断登录方式: 学号/手机号
            if (this.loginWay === '学号') {
                postUser.sno = this.value;
            }
            else {
                postUser.phone = this.value;
            }
            // 4. 登录请求
            this.$service.MainAPI.login(postUser).then((res) => {
                // 存储token
                if (!res.loginToken) {
                    showMessage('error', '获取token失败');
                }
                else {
                    // 保存token
                    setLocalStorage('loginToken', res.loginToken);
                    // 保存用户ID、Name、isRoot到store和localStroge
                    this.$store.commit('App/setUserID', res.id || -1);
                    this.$store.commit('App/setUserName', res.username || 'Guest');
                    this.$store.commit('App/setUserIsRoot', res.isRoot);
                    this.$store.commit('App/setUserState', res['state']);
                    // 登录成功跳转
                    showMessage('info', '登录成功');
                    if (this.currentRole === 'admin') {
                        this.$store.commit('App/setCurrentRole', 'admin');
                        // 保存store到localStroge
                        setLocalStorage("store", JSON.stringify(this.$store.state));
                        this.$router.push({
                            name: "AppointmentManager"
                        });
                    } else {
                        // 跳转到会议室预定界面
                        this.$store.commit('App/setCurrentRole', 'user');
                        // 保存store到localStroge
                        setLocalStorage("store", JSON.stringify(this.$store.state));
                        this.$router.push({
                            name: "ReserveMeeting",
                        });
                    }
                }
            }).finally(() => {
                this.loading = false;
            });
        }
    }
};
</script>