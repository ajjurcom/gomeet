<template>
    <div class="container-wrap">
        <Modal
            v-model="putModal"
            title="修改校区信息"
            ok-text="修改"
            @on-ok="changeCampus">
            <Form :label-width="80">
                <FormItem label="校区名:">
                    <Input v-model="campus.campus_name" placeholder="校区"></Input>
                </FormItem>
            </Form>
        </Modal>
        <Modal
            v-model="addModal"
            title="添加校区"
            ok-text="添加"
            @on-ok="addCampus">
            <Form :label-width="80">
                <FormItem label="校区名:">
                    <Input v-model="campus.campus_name" placeholder="校区"></Input>
                </FormItem>
            </Form>
        </Modal>
        <div class="container">
            <div class="select-items">
                <Button type="info" @click="showAddModal">新增校区</Button>
            </div>
            <div v-if="totalCount !== 0" class="list-items">
                <div
                    class="list-item"
                    v-for="item in itemList"
                    :key="item.id">
                    <div class="list-item-content">{{item.campus_name}}</div>
                    <div class="list-item-buttons">
                        <Button class="list-item-button" @click="editCampus(item)" type="info">修改</Button>
                        <Poptip
                            confirm
                            title="删除将无法恢复"
                            @on-ok="deleteCampus(item.id)">
                            <Button class="list-item-button" :loading="loading" type="error">删除</Button>
                        </Poptip>
                    </div>
                </div>
            </div>
            <div v-if="totalCount !== 0" class="list-page">
                <Page
                    :total="totalCount"
                    :current="requestObj.page"
                    :page-size="requestObj.onePageNum"
                    show-elevator
                    show-sizer
                    show-total
                    @on-change="changePage"
                    @on-page-size-change="changeSize"
                    transfer
                />
            </div>
            <no-data v-if="totalCount === 0" title="您还没有创建校区"></no-data>
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
    name: 'CampusManager',
    components: {
        NoData,
    },
    data() {
        return {
            putModal: false,
            addModal: false,
            totalCount: 0,
            requestObj: {
                page: 1,
                onePageNum: 10
            },
            itemList: [],
            loading: false,
            putLoading: false,
            addLoading: false,
            oldCampusName: "",
            campus: {
                campus_name: ""
            },
        }
    },
    methods: {
        getDataList() {
            this.$service.MainAPI.getCampusByPage(this.requestObj.onePageNum, this.requestObj.page).then(res => {
                this.totalCount = res.count;
                this.itemList = res.campusList;
            })
        },
        changePage(val) {
            this.requestObj.page = val;
            this.getDataList();
        },
        changeSize(val) {
            this.requestObj.onePageNum = val;
            this.getDataList();
        },
        deleteCampus(id) {
            this.loading = true;
            this.$service.MainAPI.deleteCampus(id).then(res => {
                this.$Message.success('删除成功');
                this.getDataList();
            }).finally(() => {
                this.loading = false;
            });
        },
        editCampus(campus) {
            this.putModal = true;
            this.oldCampusName = campus.campus_name || '';
            this.campus.campus_name = campus.campus_name || '';
            this.campus.id = campus.id;
        },
        showAddModal() {
            this.addModal = true;
            this.campus = {};
        },
        changeCampus() {
            if (this.campus.campus_name === this.oldCampusName) {
                this.$Message.success("修改成功");
                return;
            }
            this.putLoading = true;
            this.$service.MainAPI.putCampus(this.campus).then(res => {
                this.$Message.success("修改成功");
                this.getDataList();
            }).finally(() => {
                this.putLoading = false;
            });
        },
        addCampus() {
            if (this.campus.campus_name === '') {
                this.$Message.error("名字不能为空");
                return;
            }
            this.addLoading = true;
            this.$service.MainAPI.addCampus(this.campus).then(res => {
                this.$Message.success("添加成功");
                this.getDataList();
            }).finally(() => {
                this.addLoading = false;
            });
        }
    },
    created() {
        this.getDataList();
    }
}
</script>