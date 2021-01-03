<template>
    <div class="container-wrap">
        <div class="container">
            <Form class="register-form" ref="formValidate" :model="formValidate" :rules="ruleValidate" :label-width="80">
                <FormItem label="建筑名字" prop="building_name">
                    <Input v-model="formValidate.building_name" placeholder="Enter the building name"></Input>
                </FormItem>
                <FormItem label="校区" prop="campus_id">
                    <Select v-model="formValidate.campus_id" placeholder="Select your campus">
                        <Option 
                            v-for="(item, index) in campusList"
                            :value="item.id+''"
                            :key="index"
                            >
                            {{item.campus_name}}
                        </Option>
                    </Select>
                </FormItem>
                <FormItem label="建筑楼层" prop="layer">
                    <Input v-model="formValidate.layer" placeholder="Enter number of the building floors"></Input>
                </FormItem>
                <FormItem>
                    <Button type="primary" :loading="loading" @click="handleSubmit('formValidate')">添加</Button>
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
export default {
    name: 'BuildingAdd',
    data() {
        return {
            loading: false,
            campusList: [],
            requestObj: {},
            formValidate: {
                building_name: '',
                campus_id: -1,
                layer: 0
            },
            ruleValidate: {
                building_name: [
                    { required: true, message: 'the building name cannot be empty', trigger: 'blur' },
                    { type: 'string', max:100, message: 'length cannot longger than 100', trigger: 'blur' }
                ],
                campus_id: [
                    { required: true, message: 'Please select the campus', trigger: 'blur' },
                ],
                layer: [
                    { required: true, message: 'the layer name cannot be empty', trigger: 'blur' },
                    { type: "string", pattern: /[1-9][0-9]*/, message: 'the layer must be a number', trigger: 'blur' },
                ]
            }
        }
    },
    methods: {
        handleSubmit (name) {
            this.$refs[name].validate((valid) => {
                if (valid) {
                    this.loading = true;
                    this.requestObj.building_name = this.formValidate.building_name;
                    this.requestObj.campus_id = Number(this.formValidate.campus_id);
                    this.requestObj.layer = Number(this.formValidate.layer);
                    this.$Message.info("添加成功");
                    this.$service.MainAPI.addBuilding(this.requestObj).then(() => {
                        this.$Message.info("添加成功");
                        this.$router.push({
                            name: "BuildingManager",
                            query: {
                                campus_id: this.formValidate.campus_id
                            }
                        });
                    }).finally(() => {
                        this.loading = false;
                    })
                } else {
                    this.$Message.error('请检查值');
                }
            })
        },
        handleReset (name) {
            this.$refs[name].resetFields();
        },
    },
    created() {
        // 获取全部校区
        this.$service.MainAPI.getAllCampus().then((res) => {
            this.campusList = res.campusList;
        });
    }
}
</script>