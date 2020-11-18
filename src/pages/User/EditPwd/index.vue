<template>
    <div class="container-wrap">
        <custom-menu activeItem="personal-password"></custom-menu>
        <div class="container">
            <Form class="register-form" ref="formValidate" :model="formValidate" :rules="ruleValidate" :label-width="80">
                <FormItem label="旧密码" prop="old_password">
                    <Input type="password" v-model="formValidate.old_password" placeholder="Enter your old password"></Input>
                </FormItem>
                <FormItem label="新密码" prop="password">
                    <Input type="password" v-model="formValidate.password" placeholder="Enter your new password"></Input>
                </FormItem>
                <FormItem label="验证" prop="passwdCheck">
                    <Input type="password" v-model="formValidate.passwdCheck" placeholder="Enter your new password again"></Input>
                </FormItem>
                <FormItem>
                    <Button type="primary" :loading="loading" @click="handleSubmit('formValidate')">修改</Button>
                    <Button @click="handleReset('formValidate')" style="margin-left: 8px">Reset</Button>
                </FormItem>
                <FormItem></FormItem>
            </Form>
        </div>
    </div>
</template>

<style lang="less" scoped>
.container-wrap {
    width: 100%;
    min-width: 1200px;
    .container {
        width: 80%;
        min-width: 1024px;
        margin: 0 auto;
        padding: 20px 0;
    }
}
</style>

<script>
import CustomMenu from '@/components/CustomMenu';
export default {
    name: "UserEditPwd",
    components: {
        CustomMenu
    },
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
                id: Number(this.$route.params.id) || -1,
                old_password: '',
                password: '',
                passwdCheck: ''
            },
            ruleValidate: {
                old_password: [
                    { required: true, message: 'The password cannot be empty', trigger: 'blur' },
                    { type: 'string', min: 6, message: 'length must longger than 6' },
                    { type: 'string', max: 16, message: 'length cannot longger than 50' },
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
                if (valid && this.formValidate.id !== -1) {
                    this.$service.MainAPI.putUserPwd(this.formValidate).then(res => {
                        this.$Message.success('修改成功');
                    }).finally(() => {
                        this.loading = false;
                    });
                }
                else {
                    this.$Message.error('请检查输入值');
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