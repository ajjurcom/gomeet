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
                            v-for="item in campusList"
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
                    :disabled="buildingList.length === 0 ? true : false">
                        <Option 
                            v-for="item in buildingList"
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
            <div class="params-bottom" :style="{'height': moreParamsDiv.height + 'px', 'opacity': moreParamsDiv.opacity}">
                <div class="params-bottom-left">
                    <div class="params-scale">
                        <div class="param-scale-text">会议室类型</div>
                        <div class="param-scales">
                            <Checkbox v-model="params.checkAllTypes"  @on-change="checkAllTypesFunc">全选</Checkbox>
                            <CheckboxGroup v-model="params.meetingTypes">
                                <Checkbox
                                    v-for="(item, index) in meetingTypes"
                                    :key="index"
                                    :label="item">
                                </Checkbox>
                            </CheckboxGroup>
                        </div>
                    </div>
                    <div class="params-scale">
                        <div class="param-scale-text">会议室大小</div>
                        <div class="param-scales">
                            <Checkbox v-model="params.checkAllScales" @on-change="checkAllScalesFunc">全选</Checkbox>
                            <CheckboxGroup v-model="params.meetingScales">
                                <Checkbox
                                    v-for="(item, index) in meetingScales"
                                    :key="index"
                                    :label="item">
                                </Checkbox>
                            </CheckboxGroup>
                        </div>
                    </div>
                </div>
                <div class="params-bottom-right">
                    <Button type="success" class="params-button" @click="resetParams">清空</Button>
                    <Button type="primary" class="params-button" @click="changeParams">确定</Button>
                </div>
            </div>
        </div>
        <div class="main">
            <div class="date"></div>
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
                        margin-right: 20px;
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
        width: 80%;
        margin: 0 auto;
        height: 100px;
        border-radius: 4px;
        border: 1px solid #dcdee2;
        .date {
            height: 48px;
            width: 100%;
            background-color: #ddd;
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
export default {
    name: "Reserve",
    data () {
        return {
            role : this.$store.getters['App/getCurrentRole'],
            params: {
                campusID: -1,
                buildingID: -1,
                meetingTypes: [],
                meetingScales: [],
                checkAllScales: false,
                checkAllTypes: false
            },
            meetingScales: [],
            meetingTypes: [],
            campusList: [],
            buildingList: [],
            meetingList: [],
            moreParamsDiv: {
                height: 0,
                opacity: 0
            },
            
        };
    },
    methods: {
        initOptions() {
            this.$service.MainAPI.getScheduleOptions().then(res => {
                this.meetingScales = res.meetingScales || [];
                this.meetingTypes = res.meetingTypes || [];
                this.campusList = res.campusList || [];
                if (this.campusList.length === 0) {
                    this.$Message.info('该校区中没有建筑');
                    return
                }
                this.params.campusID = this.campusList[0].id;
                this.buildingList = res.buildingList || [];
                if (this.buildingList.length === 0) {
                    this.$Message.info('该建筑中没有会议室');
                    return
                }
                this.params.buildingID = this.buildingList[0].id;
                this.meetingList = res.meetingList || [];
            });
        },
        showMoreParams() {
            this.moreParamsDiv.height = this.moreParamsDiv.height === 0 ? 100 : 0;
            this.moreParamsDiv.opacity = this.moreParamsDiv.opacity === 0 ? 1 : 0;
        },
        changeCampus(value) {
            // this.itemList = [];
            // this.$service.MainAPI.getAllBuildingsByCampus(value).then((res) => {
            //     this.buildingList = res.buildings || [];
            //     this.$router.replace({
            //         query: {'campus_id': this.params.campusID}
            //     });
            // });
        },
        changeBuilding(value) {
            // if (value > 0) {
            //     this.$router.push({
            //         query: {
            //             'campus_id': this.params.campusID,
            //             'building_id': this.params.buildingID
            //         }
            //     });
            //     this.getDataList();
            // }
        },
        changeParams() {
            console.log('当前选择会议室类型 -> ', this.params.meetingTypes);
            console.log('当前选择会议室大小 -> ', this.params.meetingScales);
        },
        resetParams() {
            this.params.meetingTypes = [];
            this.params.meetingScales = [];
            this.params.checkAllScales = false;
            this.params.checkAllTypes = false;
        },
        checkAllTypesFunc() {
            if (this.params.meetingTypes.length !== this.meetingTypes.length) {
                this.params.meetingTypes = this.meetingTypes;
                this.params.checkAllTypes = true;
            } else if (!this.params.checkAllTypes) {
                this.params.meetingTypes = [];
            }
        },
        checkAllScalesFunc() {
            if (this.params.meetingScales.length !== this.meetingScales.length) {
                this.params.meetingScales = this.meetingScales;
                this.params.checkAllScales = true;
            } else if (!this.params.checkAllScales) {
                this.params.meetingScales = [];
            }
        }
    },
    created() {
        if (this.role !== 'user') {
            this.$Message.info('请先登录');
            this.$router.push({
                name: "Login"
            });
        }
        this.initOptions();
    }
};
</script>