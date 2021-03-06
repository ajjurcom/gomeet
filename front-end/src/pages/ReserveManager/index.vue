<template>
    <div class="container-wrap">
        <div class="container">
            <Modal
                class="appointment-modal"
                v-model="control.appointmentModal"
                title="会议详细信息">
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
                <div class="appointment-item">
                    <div class="appointment-title">参会成员：</div>
                    <div class="appointment-content">
                        <Row>
                            <Col span="12" style="width:400px">
                                <Select
                                    disabled
                                    v-model="search.members"
                                    multiple
                                    placeholder="输入关键字搜索用户">
                                    <Option v-for="user in search.results" :value="user.id" :key="user.id">
                                        {{user.username}}({{user.val}})
                                    </Option>
                                </Select>
                            </Col>
                        </Row>
                    </div>
                </div>
            </Modal>
            <Modal
                class="invitation-modal"
                v-model="control.editModal"
                title="修改会议信息"
                @on-ok="editAppointment">
                <div class="appointment-item">
                    <span class="appointment-title">主题：</span>
                    <div class="appointment-content">
                        <Input
                            v-model="updateAppointment.theme"
                            maxlength="100"
                            show-word-limit
                            clearable
                            placeholder="会议主题"
                            style="width: 400px" />
                    </div>
                </div>
                <div class="appointment-item">
                    <div class="appointment-title">内容：</div>
                    <div class="appointment-content">
                        <Input
                            type="textarea"
                            :rows="4"
                            v-model="updateAppointment.content"
                            maxlength="255"
                            show-word-limit
                            clearable
                            placeholder="会议内容简介"
                            style="width: 400px" />
                    </div>
                </div>
                <div class="appointment-item">
                    <div class="appointment-title">参会成员：</div>
                    <div class="appointment-content">
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
                                    <Option v-for="user in search.results" :value="user.id" :key="user.id">
                                        {{user.username}}({{user.val}})
                                    </Option>
                                </Select>
                            </Col>
                        </Row>
                    </div>
                </div>
                <div class="appointment-item">
                    <div class="appointment-title"></div>
                    <div class="appointment-content">
                        <RadioGroup @on-change="changeSearchWay" v-model="search.params.searchWay">
                            <Radio
                                v-for="item in search.searchWays"
                                :key="item"
                                :label="item">
                                {{search.paramsMap[item]}}
                            </Radio>
                        </RadioGroup>
                    </div>
                </div>
                <div class="appointment-item">
                    <div class="appointment-title">参会组：</div>
                    <div class="appointment-content">
                        <Select
                            ref="groups"
                            v-model="updateAppointment.groupsList"
                            multiple
                            style="width:400px"
                            @on-change="changeMember">
                            <Option
                                v-for="item in options.group_list"
                                :disabled="control.groupLoding"
                                :key="item.id"
                                :value="item.id">
                                {{item.group_name}}
                            </Option>
                        </Select>
                    </div>
                </div>
            </Modal>
            <div
                v-if="myReserve.length !== 0"
                class="reserves">
                <div class="title">
                    您有个<span style="color: #2d8cf0">{{myReserve.length}}</span>预定记录
                </div>
                <div
                    class="reserve"
                    v-for="(item, index) in myReserve"
                    :key="item.id"
                    @click="showAppointment(item.id)">
                    <div class="info">
                        <Icon type="md-time" color="#2d8cf0" size="20"/>
                        <div class="day-time">
                            <div class="day">{{item.day | dateFormate}}</div>
                            <div class="time">{{item.start_time}} - {{item.end_time}}</div>
                        </div>
                        <div class="item-division"></div>
                        <div class="theme">
                            {{item.theme | truncateStr}}
                        </div>
                        <div class="item-division"></div>
                        <div class="state" :class="{'state-pass': item.state === 'adopt'}">
                            {{stateMap[item.state]}}
                        </div>
                    </div>
                    <div class="buttons">
                        <Button
                            :disabled="item.state === 'adopt'"
                            class="button"
                            type="primary"
                            @click.stop="showEditModal(item)">
                            修改
                        </Button>
                        <Poptip
                            confirm
                            title="退订将无法恢复"
                            @click.native.stop=""
                            @on-ok="deleteAppointment(item.id, index)">
                            <Button class="button" type="error" ghost>退订</Button>
                        </Poptip>
                    </div>
                </div>
            </div>
            <div
                v-if="otherReserve.length !== 0"
                class="reserves">
                <div class="title">
                    其他<span style="color: #2d8cf0">{{otherReserve.length}}</span>个邀约会议
                </div>
                <div
                    class="reserve"
                    v-for="item in otherReserve"
                    :key="item.id"
                    @click="showAppointment(item.id)">
                    <div class="info">
                        <Icon type="md-time" color="#2d8cf0" size="20"/>
                        <div class="day-time">
                            <div class="day">{{item.day | dateFormate}}</div>
                            <div class="time">{{item.start_time}} - {{item.end_time}}</div>
                        </div>
                        <div class="item-division"></div>
                        <div class="theme">
                            {{item.theme}}
                        </div>
                        <div class="item-division"></div>
                        <div class="state" :class="{'state-pass': item.state === 'adopt'}">
                            {{stateMap[item.state]}}
                        </div>
                    </div>
                    <div class="buttons">
                    </div>
                </div>
            </div>
            <no-data v-if="myReserve.length === 0 && otherReserve.length === 0" title="您暂时没有预定会议"></no-data>
        </div>
    </div>
</template>

<style lang="less" scoped>
.container-wrap {
    width: 100%;
    min-width: 1200px;
    background-color: #f8f8f9;
    min-height: 800px;
    .container {
        width: 80%;
        min-width: 1024px;
        margin: 0 auto;
        padding: 20px 0;
        .reserves {
            margin: 20px 0;
            .title {
                color: #000;
                font-size: 20px;
                font-weight: 300;
            }
            .reserve {
                display: flex;
                flex-direction: row;
                align-items: center;
                justify-content: space-between;

                height: 80px;
                margin: 10px 0;
                padding: 0 15px;
                border: 1px solid #dcdee2;
                background-color: #fff;
                cursor: pointer;
                .info {
                    height: 100%;
                    display: flex;
                    flex-direction: row;
                    align-items: center;
                    .day-time {
                        margin: 0 10px;
                        width: 85px;
                    }
                    .theme {
                        font-size: 16px;
                        width: 200px;
                    }
                    .state {
                        font-size: 16px;
                        color: #2d8cf0;
                    }
                    .state-pass {
                        color: #19be6b;
                    }
                }
                .buttons {
                    display: flex;
                    flex-direction: row;
                    align-items: center;
                    .button {
                        margin-left: 10px;
                        cursor: pointer;
                    }
                }
            }
            .reserve:hover {
                box-shadow: 0 1px 3px rgba(0, 0, 0, 0.15);
            }
        }
    }
}
</style>

<style scoped>
.item-division {
    width: 1px;
    height: 86%;
    margin: 0 20px;
    background: linear-gradient(to bottom, rgba(134, 134, 134, 0), rgba(220,222,226,1), rgba(134, 134, 134, 0));
}
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
import NoData from "@/components/NoData";
import {intArrayToStr, FindDeleteIndex, NoContainEle, DeleteElements} from '@/Utils';
export default {
    name: 'ReserveManager',
    components: {
        NoData,
    },
    data() {
        return {
            groupUsers: [],
            marginLeft: 10,
            control: {
                deleteLoading: false,
                editModal: false,
                appointmentModal: false,
                appointmentClickAble: true,
                groupLoding: false
            },
            stateMap: {
                'verify': '审核中',
                'adopt': '审核通过'
            },
            appointment: {},
            updateAppointment: {},
            myReserve: [],
            otherReserve: [],
            options: {
                group_list: [],
            },
            search: {
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
                members: [],
                loading: false,
                results: []
            },
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
                    this.control.groupLoding = true;
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
                        this.control.groupLoding = false;
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
        getMyReserve() {
            this.$service.MainAPI.getMyReserve(this.$store.getters['App/getUserID']).then(res => {
                this.myReserve = res.myReserve || [];
                this.otherReserve = res.otherReserve || [];
            });
        },
        showAppointment(id) {
            if (!this.control.appointmentClickAble) {
                return
            }
            if (this.appointment.id && this.appointment.id === id) {
                this.control.appointmentModal = true;
                this.control.appointmentClickAble = true;
                return;
            }
            this.$service.MainAPI.getAppointment(id).then(res => {
                this.appointment = res.appointment;
                this.control.appointmentModal = true;
            }).finally(() => {
                this.control.appointmentClickAble = true;
            });
            this.$service.MainAPI.getUsersByID(id, 'appointment').then(res => {
                this.search.members = res.idList || [];
                this.search.results = res.userList || [];
                this.replaceShowVal('sno');
            });
        },
        showEditModal(obj) {
            this.updateAppointment.id = obj.id;
            this.updateAppointment.theme = obj.theme;
            this.updateAppointment.content = obj.content;
            this.search.params.searchWay = "username";
            this.$refs.searchInput.setQuery('');
            this.control.editModal = true;
            // 1. 获取参会人员、参会组
            // 2. 将userIDList赋值给this.search.members
            this.$service.MainAPI.getUsersByID(obj.id, 'appointment').then(res => {
                this.search.members = res.idList || [];
                this.search.results = res.userList || [];
                this.replaceShowVal('sno');
            });
            // 获取用户组选项
            if (!this.options.group_list || this.options.group_list.length === 0) {
                this.$service.MainAPI.getAllGroupsByCreator(this.$store.getters['App/getUserID']).then(res => {
                    this.options.group_list = res.groupList || [];
                });
            }
        },
        editAppointment() {
            this.updateAppointment['members'] = intArrayToStr(this.search.members);
            this.$service.MainAPI.putAppointment(this.updateAppointment).then(res => {
                this.$Message.info('修改成功');
                this.getMyReserve();
            });
        },
        deleteAppointment(id, index) {
            this.control.deleteLoading = true;
            this.$service.MainAPI.deleteAppointment(id, this.$store.getters['App/getUserID']).then(() => {
                this.$Message.info('删除成功');
                this.myReserve = this.myReserve.slice(0, index).concat(this.myReserve.slice(index+1));
            }).finally(() => {
                this.control.deleteLoading = false;
            });
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
    },
    filters: {
        truncateStr(value) {
            const maxLength = 10;
            if (value.length > maxLength) {
                return value.substr(0, maxLength) + '...';
            }
            return value;
        },
        dateFormate(value) {
            value = value || 'yyyyMMdd';
            return value.slice(0,4) + '-' + value.slice(4,6) + '-' + value.slice(6,8);
        }
    },
    created() {
        this.getMyReserve();
    }
}
</script>