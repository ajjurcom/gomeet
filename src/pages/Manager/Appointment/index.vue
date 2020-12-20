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
                                        {{user.username}}({{user.sno}})
                                    </Option>
                                </Select>
                            </Col>
                        </Row>
                    </div>
                </div>
            </Modal>
            <div v-if="totalCount !== 0" class="reserves">
                <div
                    class="reserve"
                    v-for="item in itemList"
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
                    </div>
                    <div class="buttons">
                        <Poptip
                            confirm
                            title="通过该会议后将无法恢复"
                            @click.native.stop=""
                            @on-ok="putState(item.id, 'adopt')">
                            <Button class="button" type="primary">通过</Button>
                        </Poptip>
                        <Poptip
                            confirm
                            title="拒绝该会议后将无法恢复"
                            @click.native.stop=""
                            @on-ok="putState(item.id, 'refuse')">
                            <Button class="button" type="error">拒绝</Button>
                        </Poptip>
                    </div>
                </div>
            </div>
            <div v-if="totalCount !== 0" class="list-page">
                <Page
                    :total="totalCount"
                    :current="requestListParams.page"
                    :page-size="requestListParams.onePageNum"
                    show-elevator
                    show-sizer
                    show-total
                    @on-change="changePage"
                    @on-page-size-change="changeSize"
                    transfer
                />
            </div>
            <no-data v-if="totalCount === 0" title="该选项暂无会议"></no-data>
        </div>
    </div>
</template>

<style lang="less" scoped>
.container-wrap {
    width: 100%;
    min-width: 1200px;
    min-height: 800px;
    .container {
        width: 80%;
        min-width: 1024px;
        margin: 0 auto;
        padding: 20px 0;
        .select-items {
            display: flex;
            flex-direction: row;
            margin-bottom: 20px;
            .select-item {
                width: 300px;
                margin-right: 20px;
            }
        }
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
        .list-page {
            display: flex;
            flex-direction: row;
            justify-content: space-around;
            margin: 20px;
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
    align-items: flex-start;
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
import {intArrayToStr} from '@/Utils';
export default {
    name: 'AppointmentManager',
    components: {
        NoData,
    },
    data() {
        return {
            control: {
                appointmentModal: false,
                appointmentClickAble: true,
                puting: false,
            },
            stateMap: {
                'verify': '等待审核',
                'adopt': '通过审核',
                'refuse': '拒绝'
            },
            appointment: {},
            options: {
                group_list: [],
                states: [],
            },
            search: {
                members: [],
                results: [],
            },
            members: [],
            itemList: [],
            totalCount: 0,
            requestListParams: {
                page: 1,
                onePageNum: 10,
                state: this.$route.query.state || ''
            },
        }
    },
    methods: {
        getDataList() {
            this.$service.MainAPI.getAppointmentByPage(this.requestListParams).then((res) => {
                this.totalCount = res.count || 0;
                this.itemList = res.appointments || [];
            });
        },
        changeState() {
            this.$service.MainAPI.getAppointmentByPage(this.requestListParams).then((res) => {
                this.totalCount = res.count || 0;
                this.itemList = res.appointments || [];
                this.$router.replace({
                    query: {'state': this.requestListParams.state}
                });
            });
        },
        changePage(val) {
            this.requestListParams.page = val;
            if (this.requestListParams.state === "") {
                this.$Message.info('请先选择状态');
                return;
            }
            this.getDataList();
        },
        changeSize(val) {
            this.requestListParams.onePageNum = val;
            if (this.requestListParams.state === "") {
                this.$Message.info('请先选择状态');
                return;
            }
            this.getDataList();
        },
        showAppointment(id) {
            if (!this.control.appointmentClickAble) {
                return
            }
            this.control.appointmentClickAble = false;
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
            });
        },
        putState(id, state) {
            if (this.control.puting) {
                return;
            }
            this.control.puting = true;
            const obj = {
                'id': id,
                'state': state
            }
            this.$service.MainAPI.putAppointmentState(obj).then(() => {
                this.$Message.info("修改成功");
                this.getDataList();
            }).finally(() => {
                this.control.puting = false;
            });
        }
    },
    filters: {
        truncateStr(value) {
            const maxLength = 20;
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
        // 获取选项
        this.$service.MainAPI.getAppointmentStates().then(res => {
            this.options.states = res.states || [];
            // 获取会议列表
            if (this.requestListParams.state === "" && this.options.states.length > 0) {
                this.requestListParams.state = this.options.states[0];
                this.$router.replace({
                    query: {'state': this.requestListParams.state}
                });
            }
            this.getDataList();
        });
    }
}
</script>