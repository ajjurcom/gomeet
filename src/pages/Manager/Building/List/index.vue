<template>
    <div class="container-wrap">
        <div class="container">
            <div class="select-building">
                <Select v-model="requestObj.campusID" style="width:200px" placeholder="选择校区"  @on-change="changeCampus">
                    <Option v-for="item in campusList" :value="item.id" :key="item.campus_name">{{ item.campus_name }}</Option>
                </Select>
            </div>
            <div class="list-items">
                <div
                    class="list-item"
                    v-for="item in itemList"
                    :key="item.id">
                    <div class="list-item-content">{{item.building_name}} - 楼高{{item.layer}}层</div>
                    <div class="list-item-buttons">
                        <Button class="list-item-button" @click="buildingEdit(item.id)" type="info">修改</Button>
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
                    :page="requestObj.page"
                    :page-size="requestObj.onePageNum"
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
        .select-building {
            margin-bottom: 20px;
        }
        .list-items {
            width: 100%;
            .list-item {
                display: flex;
                flex-direction: row;
                justify-content: space-between;
                align-items: center;
                width: 100%;
                height: 50px;
                padding: 0 20px;
                margin-bottom: 10px;
                border: 1px solid #e8eaec;
                border-radius: 4px;
                cursor: pointer;
                transition: all .2s ease;
                font-size: 14px;
                color: #17233d;
                font-weight: 550;
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
    name: 'BuildingManager',
    data() {
        return {
            campusList: [],
            totalCount: 0,
            requestObj: {
                page: 0,
                onePageNum: 10,
                campusID: Number(this.$route.query.campus_id) || -1
            },
            itemList: [],
            loading: false
        }
    },
    methods: {
        getDataList() {
            this.$service.MainAPI.getBuildingByCampusPage(this.requestObj.onePageNum, this.requestObj.page, this.requestObj.campusID).then(res => {
                this.totalCount = res.count;
                this.itemList = res.buildingList;
            })
        },
        changeCampus() {
            this.$router.replace({
                query: {'campus_id': this.requestObj.campusID}
            });
            this.getDataList();
        },
        changePage(val) {
            this.requestObj.page = val;
            if (this.requestObj.campusID === -1) {
                this.$Message.info('选择校区');
                return
            }
            this.getDataList();
        },
        changeSize(val) {
            this.requestObj.onePageNum = val;
            if (this.requestObj.campusID === -1) {
                this.$Message.info('选择校区');
                return
            }
            this.getDataList();
        },
        deleteBuilding(id) {
            this.loading = true;
            this.$service.MainAPI.deleteBuilding(id).then(res => {
                this.$Message.info('删除成功');
                this.getDataList();
            }).finally(() => {
                this.loading = false;
            });
        },
        buildingEdit(buildingID) {
            this.$router.push({
                name: 'BuildingEdit',
                params: {
                    id: buildingID
                }
            })
        }
    },
    created() {
        // 获取全部校区
        this.$service.MainAPI.getAllCampus().then((res) => {
            this.campusList = res.campusList;
            if (this.requestObj.campusID === -1 && this.campusList.length != 0) {
                this.requestObj.campusID = this.campusList[0].id;
            }
            this.getDataList();
        });
    }
}
</script>