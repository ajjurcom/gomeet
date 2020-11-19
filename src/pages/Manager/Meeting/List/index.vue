<template>
    <div class="container-wrap">
        <div class="container">
            <div class="select-items">
                <Select class="select-item" v-model="getMeetingsParams.campusID" placeholder="校区" @on-change="changeCampus">
                    <Option 
                        v-for="item in campusList"
                        :value="item.id"
                        :key="item.id"
                        >
                        {{item.campus_name}}
                    </Option>
                </Select>
                <Select class="select-item" v-model="getMeetingsParams.buildingID" placeholder="建筑, 先选择校区" @on-change="changeBuilding" :disabled="buildingList.length === 0 ? true : false">
                    <Option 
                        v-for="item in buildingList"
                        :value="item.id"
                        :key="item.id"
                        >
                        {{item.building_name}}
                    </Option>
                </Select>
            </div>
            <div class="list-items">
                <div
                    class="list-item"
                    v-for="item in itemList"
                    :key="item.id">
                    <div class="list-item-content">
                        <div class="list-item-content-top">
                            {{item.layer}}楼{{item.room_number}}  - {{item.meeting_name}}
                        </div>
                        <div class="list-item-content--bottom">
                            {{item.scale}} - 功能: {{item.meeting_type}}
                        </div>
                    </div>
                    <div class="list-item-buttons">
                        <Button class="list-item-button" @click="meetingEdit(item.id, getMeetingsParams.campusID)" type="info">修改</Button>
                        <Poptip
                            confirm
                            title="删除将无法恢复"
                            @on-ok="deleteBuilding(item.id)">
                            <Button class="list-item-button" :loading="loading" type="error">删除</Button>
                        </Poptip>
                    </div>
                </div>
            </div>
            <div class="list-page">
                <Page
                    :total="totalCount"
                    :page="getMeetingsParams.page"
                    :page-size="getMeetingsParams.onePageNum"
                    show-elevator
                    show-sizer
                    show-total
                    @on-change="changePage"
                    @on-page-size-change="changeSize"
                    transfer
                />
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
        .select-items {
            display: flex;
            flex-direction: row;
            margin-bottom: 20px;
            .select-item {
                width: 300px;
                margin-right: 20px;
            }
        }
        .list-items {
            width: 100%;
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
                }
                .list-item-content-top {
                    font-size: 16px;
                    color: #17233d;
                    font-weight: 550;
                    margin-bottom: 5px;
                }
                .list-item-content--bottom {
                    font-size: 10px;
                    color: #66686d;
                    font-weight: 550;
                }
                .list-item-buttons {
                    .list-item-button {
                        margin-left: 6px;
                    }
                }
            }
            .list-item:hover {
                box-shadow: 0 1px 3px rgba(0,0,0,.15);
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

<script>
export default {
    name: 'MeetingManager',
    data() {
        return {
            campusList: [],
            buildingList: [],
            totalCount: 0,
            getMeetingsParams: {
                page: 0,
                onePageNum: 10,
                campusID: Number(this.$route.query.campus_id) || -1,
                buildingID: Number(this.$route.query.building_id) || -1
            },
            itemList: [],
            loading: false
        }
    },
    methods: {
        getDataList() {
            this.$service.MainAPI.getMeetingsByPage(this.getMeetingsParams.onePageNum, this.getMeetingsParams.page, this.getMeetingsParams.buildingID).then(res => {
                this.totalCount = res.count;
                this.itemList = res.meetingList;
            })
        },
        changeCampus(value) {
            this.itemList = [];
            this.$service.MainAPI.getAllBuildingsByCampus(value).then((res) => {
                this.buildingList = res.buildings || [];
                this.$router.replace({
                    query: {'campus_id': this.getMeetingsParams.campusID}
                });
            });
        },
        changeBuilding(value) {
            if (value > 0) {
                this.$router.push({
                    query: {
                        'campus_id': this.getMeetingsParams.campusID,
                        'building_id': this.getMeetingsParams.buildingID
                    }
                });
                this.getDataList();
            }
        },
        changePage(val) {
            this.getMeetingsParams.page = val;
            if (this.getMeetingsParams.buildingID === -1) {
                this.$Message.info('选择建筑');
                return
            }
            this.getDataList();
        },
        changeSize(val) {
            this.getMeetingsParams.onePageNum = val;
            if (this.getMeetingsParams.buildingID === -1) {
                this.$Message.info('选择建筑');
                return
            }
            this.getDataList();
        },
        deleteBuilding(id) {
            this.loading = true;
            this.$service.MainAPI.deleteMeeting(id).then(res => {
                this.$Message.info('删除成功');
                this.getDataList();
            }).finally(() => {
                this.loading = false;
            });
        },
        meetingEdit(meetingID, campusID) {
            this.$router.push({
                name: 'MeetingEdit',
                params: {
                    id: meetingID
                },
                query: {
                    campus_id: campusID
                }
            })
        }
    },
    created() {
        // 获取选项
        this.$service.MainAPI.getAllCampus().then((res) => {
            this.campusList = res.campusList;
            if (this.getMeetingsParams.campusID != -1) {
                this.$service.MainAPI.getAllBuildingsByCampus(this.getMeetingsParams.campusID).then((res) => {
                    this.buildingList = res.buildings || [];
                    if (this.getMeetingsParams.buildingID != -1) {
                        this.getDataList();
                    }
                });
            }
        });
    }
}
</script>