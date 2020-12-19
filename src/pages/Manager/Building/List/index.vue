<template>
    <div class="container-wrap">
        <div class="container">
            <div class="options">
                <Select v-if="!search.showInput" class="item" v-model="requestObj.campusID" style="width:200px" placeholder="选择校区"  @on-change="changeCampus">
                    <Option v-for="item in campusList" :value="item.id" :key="item.campus_name">{{ item.campus_name }}</Option>
                </Select>
                <Button v-if="!search.showInput" class="item" type="info" @click="initSearch">搜索建筑</Button>
                <Button
                    v-if="!search.showInput"
                    class="item"
                    type="info"
                    :to="{
                        name: 'BuildingAdd',
                    }">新增建筑</Button>
                <Input v-if="search.showInput" v-model="search.value" placeholder="输入建筑名" style="width: 200px" />
                <Button v-if="search.showInput" class="item" type="info" :loading="search.loading" @click="searchBuilding">搜索</Button>
                <Button v-if="search.showInput" class="item" type="error" @click="cancelSearch">取消</Button>
            </div>
            <div v-if="totalCount !== 0" class="list-items">
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
            <div v-if="totalCount !== 0" class="list-page">
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
            <no-data v-if="totalCount === 0" title="暂无建筑"></no-data>
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
        .options {
            margin-bottom: 20px;
            .item {
                margin-left: 10px;
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
import NoData from "@/components/NoData";
export default {
    name: 'BuildingManager',
    components: {
        NoData,
    },
    data() {
        return {
            search: {
                showInput: false,
                loading: false,
                value: '',
            },
            campusList: [],
            totalCount: 0,
            requestObj: {
                page: 1,
                onePageNum: 10,
                campusID: Number(this.$route.query.campus_id) || -1
            },
            itemList: [],
            loading: false
        }
    },
    methods: {
        initSearch() {
            this.totalCount = 0;
            this.itemList = [];
            this.search.showInput=true;
        },
        cancelSearch() {
            this.search.showInput=false;
            this.getDataList();
        },
        searchBuilding() {
            this.requestObj.page = 1;
            this.getSearchDataList();
        },
        getSearchDataList() {
            if (this.search.value.trim() == '') {
                this.$Message.error('搜索值不能为空');
                return;
            }
            this.search.loading = true;
            this.$service.MainAPI.searchBuildings(this.requestObj.onePageNum, this.requestObj.page, this.search.value).then(res => {
                this.totalCount = res.count;
                this.itemList = res.buildingList;
                const msg = this.totalCount === 0 ? '搜索完成, 无建筑' : '搜索完成';
                this.$Message.success(msg);
            }).finally(() => {
                this.search.loading = false;
            });
        },
        getDataList() {
            this.$service.MainAPI.getBuildingByCampusPage(this.requestObj.onePageNum, this.requestObj.page, this.requestObj.campusID).then(res => {
                this.totalCount = res.count;
                this.itemList = res.buildingList;
            });
        },
        changeCampus() {
            this.$router.replace({
                query: {'campus_id': this.requestObj.campusID}
            });
            this.requestObj.page = 1;
            this.getDataList();
        },
        changePage(val) {
            this.requestObj.page = val;
            if (this.requestObj.campusID === -1) {
                this.$Message.info('选择校区');
                return
            }
            if (this.search.showInput) {
                this.getSearchDataList();
            } else {
                this.getDataList();
            }
        },
        changeSize(val) {
            this.requestObj.onePageNum = val;
            if (this.requestObj.campusID === -1) {
                this.$Message.info('选择校区');
                return
            }
            if (this.search.showInput) {
                this.getSearchDataList();
            } else {
                this.getDataList();
            }
        },
        deleteBuilding(id) {
            this.loading = true;
            this.$service.MainAPI.deleteBuilding(id).then(res => {
                this.$Message.info('删除成功');
                if (this.search.showInput) {
                    this.getSearchDataList();
                } else {
                    this.getDataList();
                }
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