<template>
    <div class="container-wrap">
        <div class="container">
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
                <FormItem>
                    <Button type="primary" :loading="loading" @click="handleSubmit('formValidate')">修改</Button>
                    <Button type="error" ghost  @click="handleReset('formValidate')" style="margin-left: 8px">Reset</Button>
                    <Poptip
                        confirm
                        :loading="upgradeLoading"
                        style="margin-left: 8px"
                        title="确定发起申请升级管理员?"
                        @on-ok="upgradeAdmin">
                        <Button class="info" :disabled="upgradeAdminAble" type="primary">申请升级管理员</Button>
                    </Poptip>
                </FormItem>
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
export default {
    name: "UserEdit",
    data () {
        return {
            loading: false,
            upgradeLoading: false,
            formValidate: {
                id: Number(this.$route.params.id) || -1,
                sno: '',
                phone: '',
                username: '',
                email: '',
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
            }
        };
    },
    computed: {
        upgradeAdminAble() {
            return this.$store.getters["App/getUserState"] !== "normal_user";
        }
    },
    methods: {
        handleSubmit (name) {
            this.loading = true;
            this.$refs[name].validate((valid) => {
                if (valid && this.formValidate.id !== -1) {
                    this.$service.MainAPI.putUser(this.formValidate).then(res => {
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
        },
        upgradeAdmin() {
            this.upgradeLoading = true;
            this.$service.MainAPI.applyAdmin(this.$store.getters['App/getUserID']).then(() => {
                this.$store.commit('App/setUserState', 'verify_admin');
                this.$Message.info('已发起申请, 等待管理员审核');
            }).finally(() => {
                this.upgradeLoading = false;
            })
        }
    },
    created() {
        this.$service.MainAPI.getUser(this.formValidate.id).then(res => {
            this.formValidate.sno = res.user.sno;
            this.formValidate.phone = res.user.phone;
            this.formValidate.username = res.user.username;
            this.formValidate.email = res.user.email;
        });
    }
};
</script>