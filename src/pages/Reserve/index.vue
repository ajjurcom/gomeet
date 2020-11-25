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
                <div class="params-more" @click="showMoreParams">
                    更多搜索
                    <Icon class="params-more-icon" size="24" type="ios-more" />
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
        <div class="main">
            <div class="dates">
                <div
                    class="date"
                    :class="{'date-active': item.index === params.dateID}"
                    v-for="item in options.dateList"
                    :key="item.index"
                    @click="changeDate(item.index)">
                    <div>{{item.title}}</div>
                    <div v-if="item.index === 0" class="date-today">（今天）</div>
                </div>
            </div>
            <div class="date-timeline">
                <div class="timeline-info"></div>
                <div class="timeline-reserves">
                    <div class="timeline-reserve">08:00</div>
                    <div class="timeline-reserve">09:00</div>
                    <div class="timeline-reserve">10:00</div>
                    <div class="timeline-reserve">11:00</div>
                    <div class="timeline-reserve">12:00</div>
                    <div class="timeline-reserve">13:00</div>
                    <div class="timeline-reserve">14:00</div>
                    <div class="timeline-reserve">15:00</div>
                    <div class="timeline-reserve">16:00</div>
                    <div class="timeline-reserve">17:00</div>
                    <div class="timeline-reserve">18:00</div>
                    <div class="timeline-reserve">19:00</div>
                    <div class="timeline-reserve">20:00</div>
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
                <div class="meeting-reserve"></div>
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
                border-right: 1px solid #efefef;
                font-weight: 500;
                cursor: pointer;
                .date-today {
                    font-size: 12px;
                }
            }
            .date-last {
                border-right: none;
            }
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
            .meeting-reserve {
                width: 80%;
                background-color: #f8f8f9;
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
</style>

<script>
import {GetDateObj, ShowDateFormat, GetNumFromScale} from '@/Utils';
export default {
    name: "Reserve",
    data () {
        return {
            role : this.$store.getters['App/getCurrentRole'],
            control: {
                checkAllScales: false,
                checkAllTypes: false,
                meetingLoading: false
            },
            params: {   // 已选择的参数
                campusID: -1,
                buildingID: -1,
                meetingTypes: [],
                meetingScales: [],
                dateID: 0,
                layer: 0       // 0表示全部楼层
            },
            options: {
                dateList: [],   // 日期
                campusList: [], // 校区
                buildingList: [], // 建筑
                meetingList: [],  // 会议室
                scalesList: [],   // 大小
                meetingsTypesList: [], // 类型
                layer: 0
            },
            divClass: {
                moreParamsDiv: {
                    height: 0,
                    opacity: 0
                },
            },
        };
    },
    methods: {
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
        changeDate(dateID) {
            if (this.params.dateID === dateID) {
                return;
            }
            this.params.dateID = dateID;
            this.$nextTick(() => {
                console.log('dateID -> ', this.params.dateID);
            });
        },
        initDate() {
            for (let i = 0; i <= 6; i++) {
                const obj = {
                    'index': i,
                    'title': '',
                    'date': GetDateObj(i)
                };
                obj.title = ShowDateFormat(obj.date);
                this.options.dateList[i] = obj;
            }
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
    },
    created() {
        if (this.role !== 'user') {
            this.$Message.info('请先登录');
            this.$router.push({
                name: "Login"
            });
        }
        // 1. 初始化日期
        this.initDate();
        // 2. 获取各种选项
        this.initOptions();
    }
};
</script>