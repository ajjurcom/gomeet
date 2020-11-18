<template>
    <div class="container-wrap">
        <custom-menu activeItem="meeting-manager"></custom-menu>
        <div class="container">
            <Form class="meeting-form" ref="formValidate" :model="formValidate" :rules="ruleValidate" :label-width="100">
                <FormItem label="会议室名字" prop="meeting_name">
                    <Input v-model="formValidate.meeting_name" placeholder="Enter the meeting name"></Input>
                </FormItem>
                <FormItem label="会议室类型" prop="meeting_type">
                    <Select v-model="formValidate.meeting_type" placeholder="会议室类型">
                        <Option 
                            v-for="(item, index) in options.typeOptions"
                            :value="item"
                            :key="index"
                            >
                            {{item}}
                        </Option>
                    </Select>
                </FormItem>
                <FormItem label="会议室大小" prop="scale">
                    <Select v-model="formValidate.scale" placeholder="会议室大小">
                        <Option 
                            v-for="(item, index) in options.scaleOptions"
                            :value="item"
                            :key="index"
                            >
                            {{item}}
                        </Option>
                    </Select>
                </FormItem>
                <FormItem label="所在校区" prop="campus_id">
                    <Select v-model="formValidate.campus_id" placeholder="所在校区" @on-change="changeCampus">
                        <Option 
                            v-for="item in options.campusOptions"
                            :value="item.id"
                            :key="item.id"
                            >
                            {{item.campus_name}}
                        </Option>
                    </Select>
                </FormItem>
                <FormItem label="所在建筑" prop="building_id">
                    <Select v-model="formValidate.building_id" placeholder="所在建筑, 先选择校区" @on-change="changeBuilding"  :disabled="options.buildingOptions.length === 0 ? true : false">
                        <Option 
                            v-for="item in options.buildingOptions"
                            :value="item.id"
                            :key="item.id"
                            >
                            {{item.building_name}}
                        </Option>
                    </Select>
                </FormItem>
                <FormItem label="所在楼层" prop="layer">
                    <Select v-model="formValidate.layer" placeholder="所在楼层, 先选择建筑" :disabled="options.meetingLayer.length === 0 ? true : false">
                        <Option 
                            v-for="item in options.meetingLayer"
                            :value="item"
                            :key="item"
                            >
                            {{item}}F
                        </Option>
                    </Select>
                </FormItem>
                <FormItem label="门牌号" prop="room_number">
                    <Input v-model="formValidate.room_number" placeholder="会议室门牌号"></Input>
                </FormItem>
                <FormItem>
                    <Button type="primary" :loading="loading" @click="handleSubmit('formValidate')">修改</Button>
                    <Button @click="handleReset('formValidate')" style="margin-left: 8px">Reset</Button>
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
import CustomMenu from '@/components/CustomMenu';
export default {
    name: 'MeetingEdit',
    components: {
        CustomMenu
    },
    data() {
        return {
            loading: false,
            options: {
                "campusOptions": [],
                "scaleOptions": [],
                "typeOptions": [],
                "buildingOptions": [],
                "meetingLayer": []
            },
            formValidate: {
                meeting_name: "",
                meeting_type: "",
                scale: 0,
                campus_id: 0,
                building_id: 0,
                layer: 0,
            },
            ruleValidate: {
                meeting_name: [
                    { required: true, message: 'the meeting name cannot be empty', trigger: 'blur' },
                    { type: 'string', max:100, message: 'length cannot longger than 100', trigger: 'blur' }
                ],
                meeting_type: [
                    { required: true, message: 'Please select the meeting_type', trigger: 'blur' },
                ],
                scale: [
                    { required: true, message: 'Please select the meeting_scale', trigger: 'blur' },
                ],
                room_number: [
                    { required: true, message: 'the meeting room number cannot be empty', trigger: 'blur' },
                    { type: 'string', max:4, message: 'length cannot longger than 4', trigger: 'blur' }
                ],
            }
        }
    },
    methods: {
        handleSubmit (name) {
            this.$refs[name].validate((valid) => {
                if (valid) {
                    this.loading = true;
                    this.$service.MainAPI.putMeeting(this.formValidate).then(() => {
                        this.$Message.info("修改成功");
                        this.$router.push({
                            name: "MeetingManager",
                            query: {
                                campus_id: this.formValidate.campus_id,
                                building_id: this.formValidate.building_id
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
        changeCampus(value) {
            if (value > 0) {
                this.$service.MainAPI.getAllBuildingsByCampus(value).then((res) => {
                    this.options.buildingOptions = res.buildings || [];
                    this.formValidate.building_id = 0;
                    this.options.meetingLayer = [];
                });
            }
        },
        changeBuilding(value) {
            if (value > 0) {
                this.$service.MainAPI.getBuildingLayer(value).then((res) => {
                    this.options.meetingLayer = Array.from({length: res.building_layer}).map((v, k) => k+1) || [];
                    this.formValidate.layer = 0;
                });
            }
        }
    },
    created() {
        // 1. 获取会议室信息
        this.formValidate.id = Number(this.$route.params.id) || -1;
        this.formValidate.campus_id = Number(this.$route.query.campus_id) || -1;
        if (this.formValidate.id === -1 || this.formValidate.campus_id === -1) {
            this.$Message.error("必要要有会议室ID和校区ID");
            return;
        }
        this.$service.MainAPI.getMeetingByID(this.formValidate.id).then(res => {
            this.formValidate = res.meeting;
            this.formValidate.campus_id = Number(this.$route.query.campus_id);
            this.$service.MainAPI.getBuildingLayer(this.formValidate.building_id).then((res) => {
                this.options.meetingLayer = Array.from({length: res.building_layer}).map((v, k) => k+1) || [];
            });
            this.$service.MainAPI.getAllBuildingsByCampus(this.formValidate.campus_id).then((res) => {
                this.options.buildingOptions = res.buildings || [];
            });
        });
        // 2. 获取选项
        /* 获取全部校区
         * 获取会议室类型选项
         * 获取会议室容量选项
         */
        this.$service.MainAPI.getMeetingOptions().then((res) => {
            this.options.campusOptions = res.campusList;
            this.options.scaleOptions = res.meetingScales;
            this.options.typeOptions = res.meetingTypes;
        });
    }
}
</script>