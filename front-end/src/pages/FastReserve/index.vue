<template>
    <div class="container-wrap">
        <div class="container">
            <Form class="register-form" :label-width="120">
                <FormItem label="日期" prop="date">
                    <Col span="12">
                        <DatePicker v-model="formValidate.date" type="date" :options="options.date" placeholder="Select date" style="width: 200px"></DatePicker>
                    </Col>
                </FormItem>
                <FormItem label="时间" prop="time">
                    <Col span="12">
                        <TimePicker
                            hide-disabled-options
                            v-model="formValidate.time"
                            :disabled-hours="disabledTime"
                            format="HH:00"
                            type="timerange"
                            placement="bottom-end"
                            placeholder="Select time"
                            style="width: 168px">
                        </TimePicker>
                    </Col>
                </FormItem>
                <FormItem label="地点" prop="locate">
                    <Select style="width: 210px" v-model="requestObj.campus_id" placeholder="所在校区">
                        <Option 
                            v-for="item in options.campus"
                            :value="item.id"
                            :key="item.id"
                            >
                            {{item.campus_name}}
                        </Option>
                    </Select>
                </FormItem>
                <FormItem label="会议室类型" prop="meetingType">
                     <Select style="width: 210px" v-model="requestObj.meeting_type" placeholder="会议室类型">
                        <Option 
                            v-for="(item, index) in options.meetingTypes"
                            :value="item"
                            :key="index"
                            >
                            {{item}}
                        </Option>
                    </Select>
                </FormItem>
                <FormItem label="会议室大小" prop="meetingScale">
                    <Select style="width: 210px" v-model="requestObj.meeting_scale" placeholder="会议室大小">
                        <Option 
                            v-for="(item, index) in options.meetingScales"
                            :value="item"
                            :key="index"
                            >
                            {{item}}
                        </Option>
                    </Select>
                </FormItem>
                <FormItem label="会议主题" prop="theme">
                    <Input
                        v-model="requestObj.theme"
                        maxlength="100"
                        show-word-limit
                        clearable
                        placeholder="会议主题"
                        style="width: 400px" />
                </FormItem>
                <FormItem label="会议内容" prop="content">
                    <Input
                        type="textarea"
                        :rows="4"
                        v-model="requestObj.content"
                        maxlength="255"
                        show-word-limit
                        clearable
                        placeholder="会议内容简介"
                        style="width: 400px" />
                </FormItem>
                <FormItem label="参会成员" prop="members">
                    <Row>
                        <Col span="12" style="width:400px">
                            <Select
                                ref="searchInput"
                                v-model="search.members"
                                multiple
                                filterable
                                placeholder="输入关键字搜索用户"
                                :remote-method="searchUsers"
                                :loading="search.loading">
                                <Option v-for="(user, index) in search.results" :value="user.id" :key="index">
                                    {{user.username}}({{user.val}})
                                </Option>
                            </Select>
                        </Col>
                    </Row>
                    <RadioGroup @on-change="changeSearchWay" v-model="search.params.searchWay">
                        <Radio
                            v-for="item in search.searchWays"
                            :key="item"
                            :label="item">
                            {{search.paramsMap[item]}}
                        </Radio>
                    </RadioGroup>
                </FormItem>
                <FormItem label="参会组" props="groups">
                    <Select
                        ref="groups"
                        v-model="reserveParams.groupsList"
                        multiple style="width:400px"
                        @on-change="changeMember">
                        <Option
                            v-for="item in options.groupsList"
                            :disabled="groupLoding"
                            :key="item.id"
                            :value="item.id">
                            {{item.group_name}}
                        </Option>
                    </Select>
                </FormItem>
                <FormItem>
                    <Button
                        type="primary"
                        :loading="loading"
                        :style="{'width': '100px'}"
                        @click="handleSubmit">预约</Button>
                </FormItem>
            </Form>
            <Modal
                v-model="appointmentModal"
                title="预定成功">
                <div class="appointment-item">
                    <div class="appointment-title">发起人：</div>
                    <div class="appointment-content">{{appointment.creator_name}}</div>
                </div>
                <div class="appointment-item">
                    <div class="appointment-title">时间：</div>
                    <div class="appointment-content">
                        {{appointment.day | dateFormate}} {{appointment.start_time}}-{{appointment.end_time}}
                    </div>
                </div>
                <div class="appointment-item">
                    <div class="appointment-title">地点：</div>
                    <div class="appointment-content">
                        <!-- 仙溪校区 - 博学楼 - F5-502 高雅阁会议室 -->
                        {{appointment.locate}}
                    </div>
                </div>
                <div class="appointment-item">
                    <span class="appointment-title">主题：</span>
                    <div class="appointment-content">{{appointment.theme}}</div>
                </div>
                <div class="appointment-item">
                    <div class="appointment-title">内容：</div>
                    <div class="appointment-content">{{appointment.content}}</div>
                </div>
            </Modal>
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

<style scoped>
.appointment-item {
    display: flex;
    flex-direction: row;
    align-items: center;
    margin: 10px 0;
}
.appointment-title {
    width: 80px;
    font-weight: 500;
    text-align: right;
}
.appointment-content {
    width: 400px;
}
</style>

<script>
import {intArrayToStr, GetDateObj, ReserveFormat, GetNumFromScale, GetNumberArr, DateFormat, FindDeleteIndex, NoContainEle, DeleteElements} from '@/Utils';
export default {
    name: "UserEdit",
    data () {
        return {
            loading: false,
            groupLoding: false,
            appointment: {},
            appointmentModal: false,
            currentHour: new Date().getHours(),
            options: {
                date: {
                    disabledDate (date) {
                        if (!date) {
                            return true;
                        }
                        return (date.valueOf() < Date.now() - 86400000) || (date.valueOf() > Date.now() + 86400000*6);
                    },
                },
                campus: [],
                meetingTypes: [],
                meetingScales: [],
                groupsList: []
            },
            reserveParams: {
                groupsList: [],
            },
            requestObj: {
                creator_id: this.$store.getters['App/getUserID'],
                creator_name: this.$store.getters['App/getUserName'],
                day: '',
                start_time: '',
                end_time: '',
                theme: '',
                content: '',
                meeting_scale: '',
                meeting_type: '',
                campus_id: 0,
            },
            formValidate: {
                date: new Date(),
                time: [],
            },
            search: {
                members: [],
                searchWays: ["username", "sno", "phone"],
                paramsMap: {
                    'sno': '学号',
                    'phone': '手机号',
                    'username': '姓名'
                },
                params: {
                    searchWay: 'username',
                    keyword: ''
                },
                loading: false,
                results: []
            },
            groupUsers: [],  // 存储组成员 id列表和各个成员信息
        };
    },
    computed: {
        disabledTime() {
            if (this.formValidate.date.format('yyyy/MM/dd') == new Date().format('yyyy/MM/dd')) {
                let l = (Array.from({length: this.currentHour})).map((v,k) => k);
                l.push(23);
                return l;
            }
            return [0,1,2,3,4,5,6,7,23];
        }
    },
    watch: {
        appointmentModal(value) {
            if (!value) {
                this.$nextTick(() => {
                    this.$router.push({
                        name: 'ReserveManager'
                    });
                });
            }
        }
    },
    methods: {
        changeMember(value) {
            // 1. 计算是新增组还是删除组
            const oldGroups = this.$refs.groups.value;
            let changeGroup = -1;
            let isAdd = true;
            if (oldGroups.length < value.length) { // 新增
                changeGroup = value[value.length-1];
            } else {    // 删除某个数组
                isAdd = false;
                changeGroup = oldGroups[FindDeleteIndex(oldGroups, value)];
            }
            // 2. 请求组中成员数据
            if (isAdd) {
                if (!this.groupUsers[changeGroup]) {
                    this.groupLoding = true;
                    this.$service.MainAPI.getUsersByID(changeGroup, 'user_group').then(res => {
                        this.groupUsers[changeGroup] = {};
                        this.groupUsers[changeGroup].idList = res.idList || [];
                        this.groupUsers[changeGroup].userList = res.userList || [];
                        // 3. 添加到参会成员中/从参会成员中删除
                        let addMembersID = NoContainEle(this.search.members, this.groupUsers[changeGroup].idList);
                        this.search.members = this.search.members.concat(addMembersID);
                        this.search.results = this.search.results.concat(
                            this.groupUsers[changeGroup].userList.filter(item => addMembersID.indexOf(item.id) !== -1)
                        );
                        this.replaceShowVal('username');
                    }).finally(() => {
                        this.groupLoding = false;
                    });
                } else {
                    let addMembersID = NoContainEle(this.search.members, this.groupUsers[changeGroup].idList);
                    this.search.members = this.search.members.concat(addMembersID);
                    this.search.results = this.search.results.concat(
                        this.groupUsers[changeGroup].userList.filter(item => addMembersID.indexOf(item.id) !== -1)
                    );
                    this.replaceShowVal('username');
                }
            }
            if (this.groupUsers[changeGroup] && !isAdd) {
                this.search.members = DeleteElements(this.search.members, this.groupUsers[changeGroup].idList);
            }
        },
        handleSubmit() {
            if (this.requestObj.theme === "") {
                this.$Message.error('会议主题不能为空');
                return;
            }
            if (!this.search.members || this.search.members.length === 0) {
                this.$Message.error('成员不能为空');
                return;
            }
            this.requestObj.day = DateFormat(this.formValidate.date);
            this.requestObj.start_time = this.formValidate.time[0];
            this.requestObj.end_time = this.formValidate.time[1];
            this.requestObj.members = intArrayToStr(this.search.members); // todo
            this.loading = true;
            this.$service.MainAPI.addAppointmentFast(this.requestObj).then(res => {
                this.appointment = res.appointment;
                this.appointmentModal = true;
            }).finally(() => {
                this.loading = false;
            });
        },
        searchUsers(query) {
            if (query.trim() !== "") {
                if (!this.search.loading) {
                    // 实现input连续输入，只发一次请求
                    this.search.loading = true;
                    clearTimeout(this.timeout);
                    this.timeout = setTimeout(() => {
                        this.search.params.keyword = query;
                        this.$service.MainAPI.searchUsers(this.search.params).then(res => {
                            this.search.results = res.userList || [];
                            // 根据查询方式将手机或者学号赋值给val
                            this.replaceShowVal(this.search.params.searchWay);
                        }).finally(() => {
                            this.search.loading = false;
                        });
                    }, 300);
                }
            } else {
                this.search.results = [];
            }
        },
        changeSearchWay() {
            this.$refs.searchInput.setQuery('');
        },
        replaceShowVal(way) {
            if (way === 'phone') {
                for (let user of this.search.results) {
                    user.val = user.phone;
                }
                return
            }
            for (let user of this.search.results) {
                user.val = user.sno;
            }
        },
    },
    filters: {
        dateFormate(value) {
            value = value || 'yyyyMMdd';
            return value.slice(0,4) + '-' + value.slice(4,6) + '-' + value.slice(6,8);
        }
    },
    created() {
        const d = new Date();
        this.formValidate.time[0] = d.getHours()+':00';
        this.formValidate.time[1] = d.getHours()+1+':00';
        this.$service.MainAPI.getMeetingOptions().then((res) => {
            this.options.campus = res.campusList || [];
            this.options.meetingScales = res.meetingScales || [];
            this.options.meetingTypes = res.meetingTypes || [];
            if (this.options.campus.length > 0) {
                this.requestObj.campus_id = this.options.campus[0].id;
            }
            if (this.options.meetingScales.length > 0) {
                this.requestObj.meeting_scale = this.options.meetingScales[0];
            }
            if (this.options.meetingTypes.length > 0) {
                this.requestObj.meeting_type = this.options.meetingTypes[0];
            }
        });
        if (this.options.groupsList.length === 0) {
            this.$service.MainAPI.getAllGroupsByCreator(this.$store.getters['App/getUserID']).then(res => {
                this.options.groupsList = res.groupList;
            })
        }
    }
};
</script>