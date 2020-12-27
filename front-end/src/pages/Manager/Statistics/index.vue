<template>
    <div class="container-wrap">
        <div class="container">
            <div class="day-select">
                <div class="select" style="marginRight: 20px">
                    <Col span="12">
                        <DatePicker
                            type="date" 
                            :options="ableDay" 
                            placeholder="开始日期" 
                            v-model="selectDay.startDay"
                            @on-change="endDayChange"
                            style="width: 200px">
                        </DatePicker>
                    </Col>
                </div>
                <div class="select">
                    <Col span="12">
                        <DatePicker 
                            type="date"
                            :options="ableDay" 
                            placeholder="结束日期" 
                            v-model="selectDay.endDay"
                            @on-change="startDayChange"
                            style="width: 200px">
                        </DatePicker>
                    </Col>
                </div>
            </div>
            <div class="ratio">
                <div class="title">
                    <div>在 </div>
                    <div class="var">{{selectDay.startDay | startDay}}</div>
                    <div> 至 </div>
                    <div class="var">{{selectDay.endDay | endDay}}</div>
                    <div>  时间段内, 共有 </div>
                    <div class="var">{{statisticsResult.count}}</div>
                    <div class="div">个预约会议</div>
                </div>
            </div>
            <div class="statisticsItems">
                <div
                    v-for="item in statisticsResult.items"
                    :key="item.title"
                    class="item">
                    <i-circle
                        :size="200"
                        :trail-width="4"
                        :stroke-width="5"
                        :percent="item.ratio"
                        stroke-linecap="square"
                        stroke-color="#43a3fb">
                        <div class="demo-Circle-custom">
                            <h1>{{item.num}}</h1>
                            <p>{{item.title}}</p>
                            <span>
                                占
                                <i>{{item.ratio}}%</i>
                            </span>
                        </div>
                    </i-circle>
                </div>
            </div>
            <div class="statistics-select">
                <RadioGroup v-model="statisticsChoice" type="button" @on-change="changeStatistics">
                    <Radio
                        v-for="item in statisticsOptions"
                        :disabled="chioceLoading"
                        :key="item"
                        :label="item">
                    </Radio>
                </RadioGroup>
            </div>
            <div class="list-items">
                <div v-if="statisticsChoice.indexOf('冷门') !== -1" class="info">注意该选项为至少1个预约的会议室排序，有的会议室可能0预约不会参与排序</div>
                <div
                    class="list-item"
                    v-for="item in itemList"
                    :key="item.id">
                    <div class="list-item-content">
                        <div class="list-item-content-top">
                            {{item.layer}}楼{{item.room_number}}  - {{item.meeting_name}}
                        </div>
                        <div class="list-item-content-bottom">
                            {{item.scale}} - 功能: {{item.meeting_type}}
                        </div>
                    </div>
                    <div class="list-item-right">
                        <div v-if="statisticsChoice.indexOf('无') === -1">共{{item.reverse_count}}次预约</div>
                        <div class="list-item-buttons" v-if="statisticsChoice.indexOf('无') !== -1 || statisticsChoice.indexOf('冷门') !== -1">
                            <Poptip
                                confirm
                                title="删除将无法恢复"
                                @on-ok="deleteMeeting(item.id)">
                                <Button class="list-item-button" :loading="loading" type="error">删除</Button>
                            </Poptip>
                        </div>
                    </div>
                </div>
            </div>
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
        .day-select {
            display: flex;
        }
        .ratio {
            margin-top: 30px;
            margin-bottom: 20px;
            .title {
                display: flex;
                font-size: 14px;
                height: 20px;
                line-height: 20px;
                font-weight: 400;
                .var {
                    margin: 0 3px;
                    font-weight: 500;
                    font-size: 18px;
                }
            }
        }
        .statisticsItems {
            margin: 40px 0;
            display: flex;
            justify-content: flex-start;
            .item {
                margin-right: 50px;
                .demo-Circle-custom{
                    & h1{
                        color: #3f414d;
                        font-size: 28px;
                        font-weight: normal;
                    }
                    & p{
                        color: #657180;
                        font-size: 14px;
                        margin: 10px 0 15px;
                    }
                    & span{
                        display: block;
                        padding-top: 15px;
                        color: #657180;
                        font-size: 14px;
                        &:before{
                            content: '';
                            display: block;
                            width: 50px;
                            height: 1px;
                            margin: 0 auto;
                            background: #e0e3e6;
                            position: relative;
                            top: -15px;
                        };
                    }
                    & span i{
                        font-style: normal;
                        color: #3f414d;
                    }
                }
            }
        }
        .statistics-select {
            margin: 20px 0;
        }
        .list-items {
            width: 100%;
            .info {
                margin: 10px 0;
            }
            .list-item {
                display: flex;
                flex-direction: row;
                justify-content: space-between;
                align-items: center;
                width: 100%;
                // height: 50px;
                padding: 0 20px;
                margin-bottom: 10px;
                border: 1px solid #e8eaec;
                border-radius: 4px;
                cursor: pointer;
                transition: all .2s ease;
                .list-item-content {
                    padding: 10px 0;
                    .list-item-content-top {
                        font-size: 16px;
                        color: #17233d;
                        font-weight: 550;
                        margin-bottom: 5px;
                    }
                    .list-item-content-bottom {
                        font-size: 10px;
                        color: #66686d;
                        font-weight: 550;
                    }
                }
                .list-item-right {
                    display: flex;
                    align-items: center;
                    .list-item-buttons {
                        margin-left: 10px;
                        .list-item-button {
                            margin-left: 6px;
                        }
                    }
                }
            }
            .list-item:hover {
                box-shadow: 0 1px 3px rgba(0,0,0,.15);
            }
        }
    }
}
</style>

<script>
export default {
    name: 'Statistics',
    data() {
        return {
            chioceLoading: false,
            loading: false,
            ableDay: {
                disabledDate (date) {
                    return date && date.valueOf() > Date.now();
                }
            },
            selectDay: {
                startDay: '',
                endDay: '',
            },
            statisticsResult: {
                count: 0,
                items: []
            },
            statisticsChoice: '',
            statisticsOptions: [],
            itemList: [],
        }
    },
    filters: {
        startDay(value) {
            return value ? value.format('yyyy年MM月dd日') : '请选择开始日期';
        },
        endDay(value) {
            return value ? value.format('yyyy年MM月dd日') : '请选择结束日期';
        }
    },
    methods: {
        startDayChange(value) {
            this.getStatisticsAppointment('true');
        },
        endDayChange(value) {
            this.getStatisticsAppointment('true');
        },
        getStatisticsAppointment(isUpdateDay) {
            const startDay = this.selectDay.startDay.format('yyyyMMdd');
            const endDay = this.selectDay.endDay.format('yyyyMMdd');
            if (startDay > endDay) {
                this.$Message.error('结束日期不能早于开始日期');
                return;
            }
            this.chioceLoading = true;
            this.$service.MainAPI.statisticsAppointment(startDay, endDay, isUpdateDay, this.statisticsChoice).then(res => {
                if (isUpdateDay) {
                    this.statisticsResult.items = res.items || [];
                    this.statisticsResult.count = res.count || 0;
                }
                this.itemList = res.statisticsList || [];
                this.$Message.success('更新成功');
            }).finally(() => {
                this.chioceLoading = false;
            });
        },
        changeStatistics() {
            this.getStatisticsAppointment(false);
        },
        deleteMeeting(id) {
            this.loading = true;
            this.$service.MainAPI.deleteMeeting(id).then(res => {
                this.$Message.success('删除成功');
                this.getStatisticsAppointment(false);
            }).finally(() => {
                this.loading = false;
            });
        }
    },
    created() {
        // 获取统计选项
        this.$service.MainAPI.statisticsOptions().then(res => {
            this.statisticsOptions = res.options || [];
            this.statisticsChoice = this.statisticsOptions ? this.statisticsOptions[0] : '';
            // 1. 获取统计数据
            this.selectDay.endDay = new Date();
            const startDay = new Date();
            startDay.setDate(startDay.getDate() - 30);
            this.selectDay.startDay = startDay;     // 设置开始日期为一个月前
            // 1. 使用当前时间段发起请求
            this.getStatisticsAppointment('true');
        });
    }
}
</script>