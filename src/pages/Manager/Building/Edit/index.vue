<template>
    <div class="container-wrap">
        <custom-menu activeItem="building-manager"></custom-menu>
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
import CustomMenu from '@/components/CustomMenu';
export default {
    name: 'BuildingEdit',
    components: {
        CustomMenu
    },
    data() {
        return {
            loading: false,
            itemObj: {},
            campusList: [],
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
                    { type: "string", pattern: /^\d+$/, message: 'the layer must be a number', trigger: 'blur' },
                ]
            }
        }
    },
    methods: {
        handleSubmit (name) {
            this.$refs[name].validate((valid) => {
                if (valid) {
                    if (!this.isChange()) {
                        this.$Message.error('建筑信息未改');
                        return
                    }
                    this.putBuilding();
                } else {
                    this.$Message.error('请检查值');
                }
            })
        },
        handleReset (name) {
            this.$refs[name].resetFields();
        },
        isChange() {
            return !(this.itemObj.building_name === this.formValidate.building_name
            && this.itemObj.campus_id === Number(this.formValidate.campus_id)
            && this.itemObj.layer === Number(this.formValidate.layer));
        },
        putBuilding() {
            this.loading = true;
            this.itemObj.building_name = this.formValidate.building_name;
            this.itemObj.campus_id = Number(this.formValidate.campus_id);
            this.itemObj.layer = Number(this.formValidate.layer);
            this.$service.MainAPI.putBuilding(this.itemObj).then(res => {
                this.$Message.success('修改信息成功');
                this.$router.push({
                    name: "BuildingManager",
                    query: {
                        campus_id: this.itemObj.campus_id
                    }
                });
            }).finally(() => {
                this.loading = false;
            });
        }
    },
    created() {
        // 1. 获取全部校区选择项
        // 获取全部校区
        this.$service.MainAPI.getAllCampus().then((res) => {
            this.campusList = res.campusList;
        });
        // 2. 获取ID建筑信息
        this.itemObj.id = Number(this.$route.params.id) || -1;
        if (this.itemObj.id === -1) {
            this.$Message.error("路由中没有建筑ID");
            return;
        }
        this.$service.MainAPI.getBuildingByID(this.itemObj.id).then(res => {
            this.itemObj = res.building;
            this.formValidate.building_name = this.itemObj.building_name;
            this.formValidate.campus_id = String(this.itemObj.campus_id);
            this.formValidate.layer = String(this.itemObj.layer);
        })
    }
}
</script>