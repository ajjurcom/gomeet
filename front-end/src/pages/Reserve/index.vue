<template>
    <div class="container-wrap">
        <div class="params">
            <div class="params-top">
                <div class="params-select">
                    <Select
                        class="select-item"
                        prefix="ios-school"
                        v-model="params.campusID"
                        placeholder="校区"
                        @on-change="changeCampus">
                        <Option 
                            v-for="item in options.campusList"
                            :value="item.id"
                            :key="item.id"
                            >
                            {{item.campus_name}}
                        </Option>
                    </Select>
                    <div class="item-division"></div>
                    <Select
                    class="select-item"
                    prefix="ios-home"
                    v-model="params.buildingID"
                    placeholder="建筑, 先选择校区"
                    @on-change="changeBuilding"
                    :disabled="options.buildingList.length === 0 ? true : false">
                        <Option 
                            v-for="item in options.buildingList"
                            :value="item.id"
                            :key="item.id"
                            >
                            {{item.building_name}}
                        </Option>
                    </Select>
                    <div class="item-division"></div>
                </div>
                <div class="params-select">
                    <router-link
                        class="params-more"
                        :to="{
                            name: 'FastReserve',
                        }">
                        快速预定
                        <Icon class="params-more-icon" size="18" type="ios-link" />
                    </router-link>
                    <div class="params-more" @click="showMoreParams">
                        更多搜索
                        <Icon class="params-more-icon" size="24" type="ios-more" />
                    </div>
                </div>
            </div>
            <div class="params-bottom" :style="{'height': divClass.moreParamsDiv.height + 'px', 'opacity': divClass.moreParamsDiv.opacity}">
                <div class="params-bottom-left">
                    <div class="params-scale">
                        <div class="param-scale-text">会议室类型</div>
                        <div class="param-scales">
                            <Checkbox size="small" v-model="control.checkAllTypes"  @on-change="checkAllTypesFunc">全选</Checkbox>
                            <CheckboxGroup v-model="params.meetingTypes" size="small">
                                <Checkbox
                                    v-for="(item, index) in options.meetingsTypesList"
                                    :key="index"
                                    :label="item">
                                </Checkbox>
                            </CheckboxGroup>
                        </div>
                    </div>
                    <div class="params-scale">
                        <div class="param-scale-text">会议室大小</div>
                        <div class="param-scales">
                            <Checkbox size="small" v-model="control.checkAllScales" @on-change="checkAllScalesFunc">全选</Checkbox>
                            <CheckboxGroup v-model="params.meetingScales" size="small">
                                <Checkbox
                                    v-for="(item, index) in options.scalesList"
                                    :key="index"
                                    :label="item">
                                </Checkbox>
                            </CheckboxGroup>
                        </div>
                    </div>
                </div>
                <div class="params-bottom-right">
                    <Button type="success" class="params-button" @click="resetParams">清空</Button>
                    <Button type="primary" :loading="control.meetingLoading" class="params-button" @click="changeParams">确定</Button>
                </div>
            </div>
        </div>
        <Modal
            class="reserve-box"
            v-model="control.reserveModal"
            :loading="control.loading"
            :title="reserveParams.meetingName"
            @on-ok="confirmReserve">
            <div class="reserve-item">
                <div class="reserve-title">会议日期：</div>
                <div class="reserve-select">
                    <Select v-model="reserveParams.day" style="width:120px">
                        <Option
                            v-for="item in options.dateList"
                            :key="item.index"
                            :value="item.date">
                            {{item.title}}
                        </Option>
                    </Select>
                </div>
            </div>
            <div class="reserve-item">
                <div class="reserve-title">会议时间：</div>
                <div class="reserve-select">
                    <Select v-model="reserveParams.startTime" style="width:70px">
                        <Option
                            v-for="item in reserveOptions.startTimeList"
                            :key="item"
                            :value="item">
                            {{item}}:00
                        </Option>
                    </Select>
                    <div class="reserve-select-space">-</div>
                    <Select v-model="reserveParams.endTime" style="width:70px">
                        <Option
                            v-for="item in reserveOptions.endTimeList"
                            :key="item"
                            :value="item">
                            {{item}}:00
                        </Option>
                    </Select>
                </div>
            </div>
            <div class="reserve-item">
                <div class="reserve-title">会议主题：</div>
                <div class="reserve-select">
                    <Input
                        v-model="reserveParams.theme"
                        maxlength="100"
                        show-word-limit
                        clearable
                        placeholder="会议主题"
                        style="width: 400px" />
                </div>
            </div>
            <div class="reserve-item">
                <div class="reserve-title">内容简介：</div>
                <div class="reserve-select">
                    <Input
                        type="textarea"
                        :rows="4"
                        v-model="reserveParams.content"
                        maxlength="255"
                        show-word-limit
                        clearable
                        placeholder="会议内容简介"
                        style="width: 400px" />
                </div>
            </div>
            <div class="reserve-item">
                <div class="reserve-title">参会成员：</div>
                <div class="reserve-select">
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
                </div>
            </div>
            <div class="reserve-item">
                <div class="reserve-title"></div>
                <RadioGroup class="reserve-select" @on-change="changeSearchWay" v-model="search.params.searchWay">
                    <Radio
                        v-for="item in search.searchWays"
                        :key="item"
                        :label="item">
                        {{search.paramsMap[item]}}
                    </Radio>
                </RadioGroup>
            </div>
            <div class="reserve-item">
                <div class="reserve-title">参会组：</div>
                <div class="reserve-select">
                    <Select
                        ref="groups"
                        v-model="reserveParams.groupsList"
                        multiple style="width:400px"
                        @on-change="changeMember">
                        <Option
                            v-for="item in options.groupsList"
                            :disabled="control.groupLoding"
                            :key="item.id"
                            :value="item.id">
                            {{item.group_name}}
                        </Option>
                    </Select>
                </div>
            </div>
        </Modal>
        <div class="main">
            <div class="dates">
                <div
                    class="date"
                    :class="{'date-active': item.date === currentTime.dateStr}"
                    v-for="item in options.dateList"
                    :key="item.index"
                    @click="changeDate(item.date, item.index)">
                    <div>{{item.title}}</div>
                    <div v-if="item.index === 0" class="date-today">（今天）</div>
                </div>
            </div>
            <div class="date-timeline">
                <div class="timeline-info">
                </div>
                <div class="timeline-reserves">
                    <div
                        class="timeline-reserve"
                        v-for="index of (options.endTime-options.startTime)"
                        :key="index">
                        {{index-1+options.startTime}}:00
                    </div>
                </div>
            </div>
            <div
                class="meeting"
                v-for="(item, index) in options.meetingList"
                :class="{'last-meeting': index === options.meetingList.length - 1}"
                :key="item.id">
                <div class="meeting-info">
                    <div class="meeting-info-top">
                        <div class="meeting-info-name">{{item.meeting_name}}</div>
                        <div class="meeting-info-layer">F{{item.layer}}-{{item.room_number}}</div>
                    </div>
                    <div class="meeting-info-bottom">
                        {{item.meeting_type}}（最多容纳{{GetNumFromScale(item.scale)}}人）
                    </div>
                </div>
                <div class="meeting-reserves">
                    <div 
                        class="meeting-reserve"
                        v-for="index of (options.endTime-options.startTime)"
                        :class="{
                            'meeting-reserve-first': index === 1,
                            'expire': (currentTime.dayIndex === 0) && (index + options.startTime <= currentTime.hour),
                            'normal': ((currentTime.dayIndex !== 0) || (index + options.startTime > currentTime.hour)) && !reserveMap[currentTime.dateStr + '-' + item.id + '-' + (index-1+options.startTime)],
                            'reserve': reserveMap[currentTime.dateStr + '-' + item.id + '-' + (index-1+options.startTime)]
                        }"
                        :key="index">
                        <div
                            class="resereve-time"
                            v-if="((currentTime.dayIndex !== 0) || (index + options.startTime > currentTime.hour)) && !reserveMap[currentTime.dateStr + '-' + item.id + '-' + (index-1+options.startTime)]">
                            {{index-1+options.startTime}}:00
                        </div>
                        <div
                            class="resereve-time"
                            v-if="reserveMap[currentTime.dateStr + '-' + item.id + '-' + (index-1+options.startTime)]">
                            {{reserveMap[currentTime.dateStr + '-' + item.id + '-' + (index-1+options.startTime)]}}
                        </div>
                        <div
                            class="reserve-text"
                            v-if="((currentTime.dayIndex !== 0) || (index + options.startTime > currentTime.hour)) && !reserveMap[currentTime.dateStr + '-' + item.id + '-' + (index-1+options.startTime)]"
                            @click="reserveMeeting(item.id, item.meeting_name, index-1+options.startTime)">
                            <div>{{index-1+options.startTime}}:00</div>
                            <div>预订</div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="main-layer">
                <div
                    :class="{'active-layer': params.layer === 0}"
                    class="layer"
                    @click="changeLayer(0)">全部</div>
                <div 
                    class="layer"
                    v-for="index of options.layer"
                    :key="index"
                    :class="{'active-layer': index === params.layer}"
                    @click="changeLayer(index)">
                    F{{index}}
                </div>
            </div>
            <div class="main-colors">
                <div class="main-color">颜色说明:</div>
                <div class="main-color">
                    <div class="main-color-demo" style="backgroundColor: #fff"></div>
                    <div class="main-color-text">可预订</div>
                </div>
                <div class="main-color">
                    <div class="main-color-demo" style="backgroundColor: #5cadff"></div>
                    <div class="main-color-text">已被预定</div>
                </div>
                <div class="main-color">
                    <div class="main-color-demo" style="backgroundColor: #f8f8f9"></div>
                    <div class="main-color-text">已过期</div>
                </div>
            </div>
        </div>
    </div>
</template>

<style lang="less" scoped>
.container-wrap {
    width: 100%;
    min-width: 1200px;
    .params {
        width: 80%;
        margin: 0 auto;
        margin-top: 30px;
        margin-bottom: 14px;
        .params-top {
            position: relative;
            z-index: 500;
            display: flex;
            flex-direction: row;
            align-items: center;
            justify-content: space-between;
            width: 100%;
            height: 54px;
            border: 1px solid #dcdee2;
            border-top-left-radius: 2px;
            border-top-right-radius: 2px;
            .params-select {
                height: 100%;
                display: flex;
                flex-direction: row;
                align-items: center;
                .item-division {
                    width: 1px;
                    height: 86%;
                    margin: 0 5px;
                    background: linear-gradient(to bottom, rgba(134, 134, 134, 0), rgba(220,222,226,1), rgba(134, 134, 134, 0));
                }
                .select-item {
                    width: 200px;
                    padding: 10px 5px;
                    border: none;
                    .ivu-select-selection {
                        border: none;
                    }
                }
            }
            .params-more {
                height: 100%;
                display: flex;
                flex-direction: row;
                align-items: center;
                padding: 0 10px;
                color: #2d8cf0;
                cursor: pointer;
                border-left: 1px solid #dcdee2;
                transition: all .2s ease-in-out;
            }
            .params-more:hover {
                background-color: #2d8cf0;
                color: #ffffff;
            }
        }
        .params-bottom {
            position: relative;
            z-index: 100;
            display: flex;
            flex-direction: row;
            justify-content: space-between;
            width: 100%;
            padding: 0 10px;
            border: 1px solid #dcdee2;
            border-bottom-right-radius: 2px;
            border-bottom-left-radius: 2px;
            border-top: none;
            transition: all .2s ease-in-out;
            .params-bottom-left {
                display: flex;
                flex-direction: column;
                justify-content: space-evenly;
                .params-scale {
                    display: flex;
                    flex-direction: row;
                    align-items: center;
                    .param-scale-text {
                        margin-right: 5px;
                    }
                    .param-scales {
                        display: flex;
                        flex-direction: row;
                        align-items: center;
                        color: #2d8cf0;
                        font-weight: 500;
                    }
                }
            }
            .params-bottom-right {
                display: flex;
                flex-direction: column;
                justify-content: space-evenly;
                .params-button {
                    width: 100px;
                }
            }
        }
    }
    .main {
        position: relative;
        width: 80%;
        margin: 0 auto;
        margin-bottom: 100px;
        border-radius: 4px;
        border: 1px solid #dcdee2;
        .main-layer {
            position: absolute;
            top: 0;
            left: -50px;
            width: 40px;
            border: 1px solid #dcdee2;
            border-radius: 4px;
            text-align: center;
            .layer {
                display: flex;
                flex-direction: row;
                align-items: center;
                justify-content: center;
                width: 100%;
                height: 25px;
                font-size: 12px;
                font-weight: 500;
                color: #808695;
                cursor: pointer;
            }
            .layer:hover {
                background-color: #2d8cf0;
                color: #fff;
            }
            .active-layer {
                background-color: #2d8cf0;
                color: #fff;
            }
        }
        .main-colors {
            position: absolute;
            top: 0;
            right: -110px;
            width: 100px;
            color: #afafaf;
            .main-color {
                display: flex;
                flex-direction: row;
                align-items: center;
                margin: 12px 0;
                height: 10px;
                .main-color-demo {
                    width: 12px;
                    height: 12px;
                    margin-right: 5px;
                    // background-color: #fff;
                    border: 1px solid #dcdee2;
                }
            }
        }
        .dates {
            display: flex;
            flex-direction: row;
            align-items: center;
            height: 48px;
            width: 100%;
            background-color: #f8f8f9;
            border-bottom: 1px solid #dcdee2;
            .date {
                display: flex;
                flex-direction: column;
                align-items: center;
                justify-content: center;
                height: 100%;
                flex: 1;
                text-align:center;
                color: #c5c8ce;
                font-weight: 500;
                cursor: pointer;
                .date-today {
                    font-size: 12px;
                }
            }
            // .date-last {
            //     border-right: none;
            // }
            .date-active {
                background-color: #fff;
                box-shadow: 0 0 5px #dcdee2;
                color: #2d8cf0;
            }
        }
        .date-timeline {
            display: flex;
            flex-direction: row;
            width: 100%;
            height: 38px;
            background-color: #fff;
            .timeline-info {
                width: 20%;
            }
            .timeline-reserves {
                display: flex;
                flex-direction: row;
                align-items: center;
                width: 80%;
                .timeline-reserve {
                    flex: 1;
                    text-align: left;
                    color: #808695;
                    font-weight: 500;
                }
            }
        }
        .meeting {
            display: flex;
            flex-direction: row;
            width: 100%;
            height: 58px;
            margin-bottom: 10px;
            border-top: 1px solid #dcdee2;
            border-bottom: 1px solid #dcdee2;
            .meeting-info {
                display: flex;
                flex-direction: column;
                width: 20%;
                height: 100%;
                padding: 5px;
                border-right: 1px solid #dcdee2;
                .meeting-info-top {
                    display: flex;
                    flex-direction: row;
                    align-items: center;
                    height: 60%;
                    .meeting-info-name {
                        margin-right: 5px;
                        font-size: 16px;
                        color: #17233d;
                        font-weight: 500;
                    }
                    .meeting-info-layer {
                        color: #808695;
                    }
                }
                .meeting-info-bottom {
                    height: 40%;
                    font-size: 12px;
                    color: #808695;
                }
            }
            .meeting-reserves {
                display: flex;
                flex-direction: row;
                width: 80%;
                height: 100%;
                .meeting-reserve {
                    position: relative;
                    flex: 1;
                    height: 100%;
                    border-left: 1px solid #dcdee2;
                    color: #808695;
                    font-size: 12px;
                    background-color: #f8f8f9;
                    padding: 8px 3px;
                    .resereve-time {
                        position: relative;
                        z-index: 100;
                    }
                    .reserve-text {
                        position: absolute;
                        top: -3%;
                        left: -3%;
                        width: 106%;
                        height: 106%;
                        z-index: 200;
                        background-color: #2d8cf0;
                        border-radius: 4px;
                        box-shadow: 0 0 4px #2d8cf0;
                        opacity: 0;
                        font-size: 12px;
                        padding: 8px 3px;
                    }
                    .reserve-text:hover {
                        opacity: 1;
                    }
                }
                .meeting-reserve-first {
                    border-left: none;
                }
                .expire {
                    background-color: #f8f8f9;
                    border-left: none;
                }
                .normal {
                    background-color: #fff;
                    cursor: pointer;
                }
                .normal:hover {
                    color: #fff;
                    background-color: #2d8cf0;
                }
                .reserve {
                    color: #fff;
                    background-color: rgb(92, 173, 255);
                }
            }
        }
        .last-meeting {
            margin-bottom: 0;
            border-bottom: none;
        }
    }
}
</style>

<style scoped>
.select-item >>> .ivu-select-selection {
    border: none;
    outline: 0px;
    box-shadow: none;
}
.select-item select {
    border: red;
}
/* .reserve-box {
} */
.reserve-item {
    display: flex;
    flex-direction: row;
    align-items: center;
    margin: 10px 0;
}
.reserve-select {
    display: flex;
    flex-direction: row;
}
.reserve-select-space {
    font-size: 20px;
    font-weight: 500;
    margin: 0 10px;
}
.reserve-title {
    width: 80px;
    text-align: right;
}
</style>

<script>
import {intArrayToStr, GetDateObj, ReserveFormat, GetNumFromScale, GetNumberArr, DateFormat, FindDeleteIndex, NoContainEle, DeleteElements} from '@/Utils';
export default {
    name: "Reserve",
    data () {
        return {
            reserveMap: {},
            currentTime: {
                dayIndex: 0,
                dateStr: '',
                hour: new Date().getHours(),
                // minute: new Date().getMinutes(),
                // second: new Date().getSeconds(),
            },
            control: {
                checkAllScales: false,
                checkAllTypes: false,
                meetingLoading: false,
                reserveModal: false,
                loading: true,
                groupLoding: false
            },
            params: {   // 已选择的参数
                campusID: -1,
                buildingID: -1,
                meetingTypes: [],
                meetingScales: [],
                layer: 0       // 0表示全部楼层
            },
            reserveParams: {
                meetingID: 0,
                meetingName: '会议室名字',
                day: '',
                startTime: 8,
                endTime: 22,
                groupsList: [],
                theme: '',
                content: ''
            },
            reserveOptions: {
                startTimeList: [],
                endTimeList: []
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
            options: {
                dateList: [],   // 日期
                campusList: [], // 校区
                buildingList: [], // 建筑
                meetingList: [],  // 会议室
                scalesList: [],   // 大小
                meetingsTypesList: [], // 类型
                layer: 0,
                startTime: 8,      // 可以预定最早几点
                endTime: 22,
                groupsList: [],
            },
            divClass: {
                moreParamsDiv: {
                    height: 0,
                    opacity: 0
                },
            },
            groupUsers: [],  // 存储组成员 id列表和各个成员信息
        };
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
        confirmReserve() {
            if (this.reserveParams.theme === "") {
                this.$Message.error("会议主题不能为空");
                this.control.loading = false;
                return;
            }
            if (this.search.members.length === 0) {
                this.$Message.error("参会人员不能为空");
                this.control.loading = false;
                return;
            }
            this.control.loading = true;
            // 1. 检查参数
            const obj = {
                creator_id: this.$store.getters['App/getUserID'],
                creator_name: this.$store.getters['App/getUserName'],
                meeting_id: this.reserveParams.meetingID,
                day: this.reserveParams.day,
                start_time: this.reserveParams.startTime < 10 ? '0' + this.reserveParams.startTime + ':00' : this.reserveParams.startTime + ':00',
                end_time: this.reserveParams.endTime < 10 ? '0' + this.reserveParams.endTime + ':00' : this.reserveParams.endTime + ':00',
                theme: this.reserveParams.theme,
                content: this.reserveParams.content,
                members: intArrayToStr(this.search.members),
            };
            // 2. 发起请求
            this.$service.MainAPI.addAppointment(obj).then(res => {
                this.$Message.info('预定成功');
                this.control.reserveModal = false;
                this.initReserve();
            }).finally(() => {
                this.control.loading = false;
            });
        },
        initOptions() {
            this.$service.MainAPI.getScheduleOptions().then(res => {
                this.options.scalesList = res.meetingScales || [];
                this.options.meetingsTypesList = res.meetingTypes || [];
                this.options.campusList = res.campusList || [];
                if (this.options.campusList.length === 0) {
                    this.$Message.info('该校区中没有建筑');
                    return
                }
                this.params.campusID = this.options.campusList[0].id;
                this.options.buildingList = res.buildingList || [];
                if (this.options.buildingList.length === 0) {
                    this.$Message.info('该建筑中没有会议室');
                    return
                }
                this.options.layer = this.options.buildingList[0].layer;
                this.params.buildingID = this.options.buildingList[0].id;
                this.options.meetingList = res.meetingList || [];
                this.initReserve();
            });
        },
        updateOptions(way) {
            if (way === 'campus' || way === 'building') {
                this.params.layer = 0;
                this.params.meetingTypes = [];
                this.params.meetingScales = [];
                this.control.checkAllScales = false;
                this.control.checkAllTypes = false;
            }
            // 参数: [校区、建筑、楼层、会议室类型、会议室大小]
            this.control.meetingLoading = true;
            this.$service.MainAPI.updateOptions(this.params, way).then((res) => {
                switch(way) {
                    case 'campus':
                        this.options.buildingList = res.buildingList || [];
                        this.params.buildingID = this.options.buildingList.length > 0 ? this.options.buildingList[0].id : 0 || 0;
                    case 'campus':
                    case 'building':
                        this.options.layer = this.options.buildingList.length > 0 ? this.options.buildingList[0].layer : 0 || 0;
                    default:
                        this.options.meetingList = res.meetingList || [];
                }
            }).finally(() => {
                if (this.options.meetingList.length === 0) {
                    this.$Message.info('无会议室');
                } else {
                    this.$Message.info('会议室已更新');
                }
                this.control.meetingLoading = false;
                if (way === 'campus' || way === 'building') {
                    this.initReserve();
                }
            });
        },
        showMoreParams() {
            this.divClass.moreParamsDiv.height = this.divClass.moreParamsDiv.height === 0 ? 100 : 0;
            this.divClass.moreParamsDiv.opacity = this.divClass.moreParamsDiv.opacity === 0 ? 1 : 0;
        },
        changeLayer(value) {
            if (value === this.params.layer) {
                return;
            }
            this.params.layer = value;
            this.updateOptions('layer');
        },
        changeCampus() {
            this.updateOptions('campus');
        },
        changeBuilding(value) {
            this.updateOptions('building');
        },
        changeParams() {
            this.updateOptions('meetingType');
        },
        changeDate(dateStr, dayIndex) {
            if (this.currentTime.dateStr === dateStr) {
                return;
            }
            this.currentTime.dayIndex = dayIndex;
            this.currentTime.dateStr = dateStr;
            this.initReserve();
        },
        initDate() {
            for (let i = 0; i <= 6; i++) {
                let date = GetDateObj(i);
                const obj = {
                    'index': i,
                    'title': ReserveFormat(date),
                    'date': DateFormat(date)
                };
                this.options.dateList[i] = obj;
            }
            this.currentTime.dateStr = this.options.dateList[0].date;
        },
        GetNumFromScale(str) {
            return GetNumFromScale(str);
        },
        resetParams() {
            this.params.meetingTypes = [];
            this.params.meetingScales = [];
            this.control.checkAllScales = false;
            this.control.checkAllTypes = false;
        },
        checkAllTypesFunc() {
            if (this.params.meetingTypes.length !== this.options.meetingsTypesList.length) {
                this.params.meetingTypes = this.options.meetingsTypesList;
                this.control.checkAllTypes = true;
            } else if (!this.control.checkAllTypes) {
                this.params.meetingTypes = [];
            }
        },
        checkAllScalesFunc() {
            if (this.params.meetingScales.length !== this.options.scalesList.length) {
                this.params.meetingScales = this.options.scalesList;
                this.control.checkAllScales = true;
            } else if (!this.control.checkAllScales) {
                this.params.meetingScales = [];
            }
        },
        onTimeTask() {
            const date = new Date();
            this.currentTime.hour = date.getHours() < 10 ? '0'+date.getHours() : date.getHours();;
            // this.currentTime.minute = date.getMinutes() < 10 ? '0'+date.getMinutes() : date.getMinutes();
            // this.currentTime.second = date.getSeconds() < 10 ? '0'+date.getSeconds() : date.getSeconds();;
        },
        clearModalValue() {
            this.search.members = [];
            this.reserveParams.groupsList = [];
            this.reserveParams.theme = '';
            this.reserveParams.content = '';
            this.$refs.searchInput.setQuery('');
        },
        reserveMeeting(meetingID, meetingName, startTimeHour) {
            this.clearModalValue();
            // 会议信息赋值
            this.reserveParams.meetingID = meetingID;
            this.reserveParams.meetingName = meetingName;
            this.reserveParams.startTime = startTimeHour;
            this.reserveParams.endTime = startTimeHour+1;
            this.reserveParams.day = this.currentTime.dateStr;
            this.reserveOptions.startTimeList = GetNumberArr(startTimeHour, this.options.endTime-1);
            this.reserveOptions.endTimeList = GetNumberArr(startTimeHour+1, this.options.endTime);
            this.control.reserveModal = true;
            // 获取分组选项
            if (this.options.groupsList.length === 0) {
                this.$service.MainAPI.getAllGroupsByCreator(this.$store.getters['App/getUserID']).then(res => {
                    this.options.groupsList = res.groupList;
                })
            }
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
        changeSearchWay() {
            this.$refs.searchInput.setQuery('');
        },
        initReserve() {
            let meetingList = this.options.meetingList.map(item => {
                return item.id;
            });
            const obj = {
                'day': this.currentTime.dateStr,
                'start_time': this.currentTime.dayIndex !== 0 ? '00:00' : this.currentTime.hour < 10 ? '0'+this.currentTime.hour+':00' : this.currentTime.hour+':00',
                'meeting_id': meetingList,
            }
            this.$service.MainAPI.getReserve(obj).then(res => {
                // 存入reserveMap字典中
                this.reserveMap = {};
                if (res.appointments) {
                    for (let item of res.appointments) {
                        const start = item.start_time.substr(0,1) === '0' ? Number(item.start_time.substr(1,1)) : Number(item.start_time.substr(0,2));
                        const end = item.end_time.substr(0,1) === '0' ? Number(item.end_time.substr(1,1)) : Number(item.end_time.substr(0,2));
                        for (let i = start; i < end; i++) {
                            this.reserveMap[item.day+'-'+item.meeting_id+'-'+i] = item.creator_name;
                        }
                    }
                }
            });
            
        }
    },
    created() {
        if (this.$store.getters['App/getCurrentRole'] !== 'user') {
            this.$Message.info('请先登录');
            this.$router.push({
                name: "Login"
            });
        }
        // 1. 初始化日期
        this.initDate();
        // 2. 获取各种选项, then -> 获取显示的会议室预约情况 this.initReserve();
        this.initOptions();
        // 4. 每秒钟执行更新时间任务
        setInterval(() => {
            this.onTimeTask();
        }, 60000);
    }
};
</script>