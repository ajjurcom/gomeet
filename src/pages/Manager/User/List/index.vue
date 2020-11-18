<template>
    <div class="container-wrap">
        <custom-menu activeItem="user"></custom-menu>
        <div class="container">
            <div class="select-items">
                <Select class="select-item" v-model="getMeetingsParams.state" placeholder="状态" @on-change="changeState">
                    <Option 
                        v-for="state in stateList"
                        :value="state"
                        :key="state"
                        >
                        {{stateMap[state]}}
                    </Option>
                </Select>
            </div>
            <div class="list-items">
                <div
                    class="list-item"
                    v-for="item in itemList"
                    :key="item.id">
                    <div class="list-item-content">
                        <div class="list-item-text">
                            {{item.sno}} - {{item.username}}
                        </div>
                    </div>
                    <div class="list-item-buttons">
                        <Select class="list-item-buttons-state" v-model="item.state" placeholder="状态" @on-change="putUserState(item.id, item.state)">
                            <Option 
                                v-for="state in stateList"
                                :value="state"
                                :key="state"
                                >
                                {{state}}
                            </Option>
                        </Select>
                        <Button class="list-item-button" type="success" :loading="detailsLoading" @click="showUserInfo(item.id)">详情</Button>
                        <Poptip
                            confirm
                            title="删除将无法恢复"
                            @on-ok="deleteUser(item.id)">
                            <Button class="list-item-button" :loading="deleteLoading" type="error">删除</Button>
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
                    display: flex;
                    flex-direction: row;
                    .list-item-buttons-state {
                        width: 115px;
                        margin-right: 10px;
                    }
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
import CustomMenu from '@/components/CustomMenu';
export default {
    name: 'UserManager',
    components: {
        CustomMenu
    },
    data() {
        return {
            isRoot: this.$store.getters['App/getUserIsRoot'],
            stateMap: {
                "verify_user": "待审核用户",
                "normal_user": "正常用户",
                "refuse_user": "拒绝用户",
                "blacklist": "黑名单",
                "verify_admin": "待审核管理员",
                "normal_admin": "正常管理员"
            },
            stateList: [],
            totalCount: 0,
            getMeetingsParams: {
                page: 0,
                onePageNum: 10,
                state: this.$route.query.state || ''
            },
            userInfo: {},
            itemList: [],
            deleteLoading: false,
            detailsLoading: false,
        }
    },
    methods: {
        getDataList() {
            this.$service.MainAPI.getUsersByPage(this.getMeetingsParams).then((res) => {
                this.totalCount = res.count || 0;
                this.itemList = res.userList || [];
            });
        },
        changeState() {
            this.$service.MainAPI.getUsersByPage(this.getMeetingsParams).then((res) => {
                this.totalCount = res.count || 0;
                this.itemList = res.userList || [];
                this.$router.replace({
                    query: {'state': this.getMeetingsParams.state}
                });
            });
        },
        changePage(val) {
            this.getMeetingsParams.page = val;
            if (this.getMeetingsParams.state === "") {
                this.$Message.info('请先选择状态');
                return;
            }
            this.getDataList();
        },
        changeSize(val) {
            this.getMeetingsParams.onePageNum = val;
            if (this.getMeetingsParams.state === "") {
                this.$Message.info('请先选择状态');
                return;
            }
            this.getDataList();
        },
        deleteUser(id) {
            this.deleteLoading = true;
            this.$service.MainAPI.deleteUser(id).then(res => {
                this.$Message.info('删除成功');
                this.getDataList();
            }).finally(() => {
                this.deleteLoading = false;
            });
        },
        putUserState(id, state) {
            this.$service.MainAPI.putUserState(id, state).then(res => {
                this.$Message.success("修改成功");
            });
        },
        showUserInfo(id) {
            this.detailsLoading = true;
            this.$service.MainAPI.getUserInfo(id).then(res => {
                this.userInfo = res.user || {};
                const title = '用户详情';
                let content = `<p>学号: ${this.userInfo.sno}</p>
                    <p>名字: ${this.userInfo.username}</p>
                    <p>电话: ${this.userInfo.phone}</p>
                    <p>邮箱: ${this.userInfo.email}</p>
                    <p>状态: ${this.stateMap[this.userInfo.state]}</p>`;
                if (this.userInfo.state === 'blacklist') {
                    content += `<p>解封时间: ${this.userInfo.ban}</p>`;
                }
                this.$Modal.info({
                    title: title,
                    content: content
                });
            }).finally(() => {
                this.detailsLoading = false;
            });
        },
    },
    created() {
        // 获取选项
        this.$service.MainAPI.getUserOptions(this.isRoot ? 'root' : 'admin').then(res => {
            this.stateList = res.stateList || [];
        });
        // 获取用户列表
        if (this.getMeetingsParams.state === "") {
            this.$Message.info("请选择状态");
            return;
        }
        this.$service.MainAPI.getUsersByPage(this.getMeetingsParams).then((res) => {
            this.itemList = res.userList || [];
            this.totalCount = res.count || 0;
        });
    }
}
</script>