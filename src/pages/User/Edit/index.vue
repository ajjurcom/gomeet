<template>
    <div class="container-wrap">
        <custom-menu activeItem="personal-edit"></custom-menu>
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
    name: "UserEdit",
    components: {
        CustomMenu
    },
    data () {
        return {
            loading: false,
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