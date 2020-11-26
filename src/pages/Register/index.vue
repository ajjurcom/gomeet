<template>
    <div class="container">
        <div class="back-image"></div>
        <div class="main">
            <div class="main-title">
                <div class="title-name">GoMeet</div>
                <div class="title-info">佛科院会议室预定系统</div>
            </div>
            <div class="register-container">
                <div class="back-login">
                    <router-link to="Login"><Icon type="ios-arrow-back" />登录</router-link>
                </div>
                <Form class="register-form" ref="formValidate" :model="formValidate" :rules="ruleValidate" :label-width="80">
                    <FormItem label="学号" prop="sno">
                        <Input v-model="formValidate.sno" placeholder="Enter your sno"></Input>
                    </FormItem>
                    <FormItem label="手机" prop="phone">
                        <Input v-model="formValidate.phone" placeholder="Enter your phone"></Input>
                    </FormItem>
                    <FormItem label="名字" prop="username">
                        <Input v-model="formValidate.username" placeholder="Enter your name"></Input>
                    </FormItem>
                    <FormItem label="邮箱" prop="email">
                        <Input v-model="formValidate.email" placeholder="Enter your e-mail"></Input>
                    </FormItem>
                    <FormItem label="密码" prop="password">
                        <Input type="password" v-model="formValidate.password" placeholder="Enter your password"></Input>
                    </FormItem>
                    <FormItem label="验证" prop="passwdCheck">
                        <Input type="password" v-model="formValidate.passwdCheck" placeholder="Enter your password again"></Input>
                    </FormItem>
                    <FormItem>
                        <Button type="primary" :loading="loading" @click="handleSubmit('formValidate')">Submit</Button>
                        <Button @click="handleReset('formValidate')" style="margin-left: 8px">Reset</Button>
                    </FormItem>
                    <FormItem></FormItem>
                </Form>
            </div>
        </div>
    </div>
</template>

<style lang="less" scoped>
.container {
    width: 100%;
    min-width: 1200px;
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
        .register-container {
            width: 400px;
            padding: 0 24px;
            background-color: #fff;
            opacity: 0.8;
            .back-login {
                font-size: 16px;
                font-weight: 600;
                line-height: 60px;
                height: 60px;
                margin-right: 24px;
                color: #444;
            }
        }
    }
}
</style>

<script>
import {showMessage} from '@/Utils';
export default {
    name: "Register",
    data () {
        const validatePass = (rule, value, callback) => {
            if (value === '') {
                callback(new Error('Please enter your password'));
            } else {
                if (this.formValidate.passwdCheck !== '') {
                    // 对第二个密码框单独验证
                    this.$refs.formValidate.validateField('passwdCheck');
                }
                callback();
            }
        };
        const validatePassCheck = (rule, value, callback) => {
            if (value === '') {
                callback(new Error('Please enter your password again'));
            } else if (value !== this.formValidate.password) {
                callback(new Error('The two input passwords do not match!'));
            } else {
                callback();
            }
        };
        return {
            loading: false,
            campusList: [],
            formValidate: {
                sno: '',
                phone: '',
                username: '',
                email: '',
                password: '',
                passwdCheck: '',
            },
            ruleValidate: {
                sno: [
                    { required: true, message: 'The sno cannot be empty', trigger: 'blur' },
                    { type: 'string', min: 11, max:11, message: 'The sno is 11 digits', trigger: 'blur' }
                ],
                phone: [
                    { required: true, message: 'The phone cannot be empty', trigger: 'blur' },
                    { type: 'string', min: 11, max:11, message: 'The phone is 11 digits', trigger: 'blur' }
                ],
                username: [
                    {required: true, message: 'The username cannot be empty', trigger: 'blur'},
                    { type: 'string', min: 1, max: 50, message: 'length cannot longger than 50', trigger: 'blur' }
                ],
                email: [
                    { required: true, message: 'Mailbox cannot be empty', trigger: 'blur' },
                    { type: 'email', message: 'Incorrect email format', trigger: 'blur' }
                ],
                password: [
                    { required: true, message: 'The password cannot be empty', trigger: 'blur' },
                    { type: 'string', min: 6, message: 'length must longger than 6' },
                    { type: 'string', max: 16, message: 'length cannot longger than 50' },
                    { validator: validatePass, trigger: 'blur' }
                ],
                passwdCheck: [
                    { required: true, message: 'Please verify the password', trigger: 'blur' },
                    { validator: validatePassCheck, trigger: 'blur' }
                ]
            }
        };
    },
    methods: {
        handleSubmit (name) {
            this.loading = true;
            this.$refs[name].validate((valid) => {
                if (valid && this.formValidate.password === this.formValidate.passwdCheck) {
                    this.loading = false;
                    this.$service.MainAPI.addUser(this.formValidate).then((res) => {
                        showMessage('success', '注册成功, 请等待管理员审核');
                        this.$router.push({
                            name: "Login",
                        });
                    }).finally(() => {
                        this.loading = false;
                    });
                }
                else {
                    this.$Message.error('Fail!');
                    this.loading = false;
                }
            })
        },
        handleReset (name) {
            this.$refs[name].resetFields();
        }
    }
};
</script>