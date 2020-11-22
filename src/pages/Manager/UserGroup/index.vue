<template>
    <div class="container-wrap">
        <div class="container">
            <div class="select-items">
                <Button type="info" @click="showGroupModal(false)">新增分组</Button>
                <Modal
                    class="search-box"
                    v-model="modal.show"
                    :title="modal.title"
                    :loading="loading"
                    @on-ok="onOkGroupModal">
                    <div class="group-name">
                        <div class="search-text">分组名字：</div>
                        <Input :disabled="modal.isChangeMember" v-model="search.group_name" placeholder="分组名字" style="width: 300px" />
                    </div>
                    <div class="search-way">
                        <div class="search-text">搜索方式：</div>
                        <RadioGroup @on-change="changeSearchWay" v-model="search.params.searchWay">
                            <Radio
                                v-for="item in search.searchWays"
                                :key="item"
                                :label="item">
                                {{search.paramsMap[item]}}
                            </Radio>
                        </RadioGroup>
                    </div>
                    <div class="search-input">
                        <div class="search-text">分组成员：</div>
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
                                    <Option v-for="user in search.results" :value="user.id" :key="user.id">
                                        {{user.username}}({{user.val}})
                                    </Option>
                                </Select>
                            </Col>
                        </Row>
                    </div>
                </Modal>
                <Modal
                    v-model="putNameModal"
                    title="修改名字"
                    @on-ok="putGroupName">
                    新名字：<Input v-model="currentGroup.group_name" placeholder="分组名字" style="width: 300px" />
                </Modal>
            </div>
            <div class="list-items">
                <div
                    class="list-item"
                    v-for="item in itemList"
                    :key="item.id">
                    <div class="list-item-content">
                        <div class="list-item-text">
                            {{item.group_name}}
                        </div>
                    </div>
                    <div class="list-item-buttons">
                        <Button class="list-item-button" type="info" @click="showNameModal(item)">改名</Button>
                        <Button class="list-item-button" type="success" @click="showGroupModal(true, item.id, item.group_name)">成员</Button>
                        <Poptip
                            confirm
                            title="删除将无法恢复"
                            @on-ok="deleteGroup(item.id)"
                            >
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

<style scoped>
.search-way {
    display: flex;
    flex-direction: row;
    align-items: center;
    margin: 10px 0;
}
.search-input {
    display: flex;
    flex-direction: row;
    align-items: center;
    margin: 10px 0;
}
.group-name,
.search-way,
.search-input {
    display: flex;
    flex-direction: row;
    align-items: center;
}
.search-text {
    width: 80px;
    text-align: right;
}
</style>

<script>
import {intArrayToStr} from '@/Utils';
export default {
    name: 'UserGroup',
    data() {
        return {
            loading: true,
            totalCount: 0,
            getMeetingsParams: {
                page: 1,
                onePageNum: 10,
                creator: Number(this.$store.getters['App/getUserID']),
            },
            itemList: [],
            deleteLoading: false,
            putNameModal: false,
            modal: {
                title: '添加分组',
                isChangeMember: true,
                show: false,
            },
            currentGroup: {
                id: -1,
                creator: -1,
                group_name: -1,
                member_list: ""
            },
            search: {
                id: "",
                group_name: "",
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
            }
        }
    },
    watch: {
        'search.members'(value) {
            console.log('改变了，新值 -> ', value, '最后 -> ', this.search.members);
        }
    },
    methods: {
        getDataList() {
            this.$service.MainAPI.getGroupsByPage(this.getMeetingsParams).then((res) => {
                this.totalCount = res.count || 0;
                this.itemList = res.groupList || [];
            });
        },
        changePage(val) {
            this.getMeetingsParams.page = val;
            this.getDataList();
        },
        changeSize(val) {
            this.getMeetingsParams.onePageNum = val;
            this.getDataList();
        },
        deleteGroup(id) {
            this.deleteLoading = true;
            this.$service.MainAPI.deleteGroup(id).then(res => {
                this.$Message.info('删除成功');
                this.getDataList();
            }).finally(() => {
                this.deleteLoading = false;
            });
        },
        putGroupName() {
            this.$service.MainAPI.putGroupName(this.currentGroup).then(res => {
                this.$Message.info('修改成功');
                this.getDataList();
            });
        },
        showGroupModal(isChangeMember, id, grouopName) {
            this.search.params.searchWay = "username";
            this.modal.isChangeMember = isChangeMember;
            this.modal.show = true;
            this.$refs.searchInput.setQuery('');
            // 修改分组成员信息
            if (isChangeMember) {
                this.search.id = id;
                this.modal.title = "修改分组成员"
                // 1. 获取组成员ID
                // 2. 将userIDList赋值给this.search.members
                this.$service.MainAPI.getGroupMembers(id).then(res => {
                    this.search.members = res.idList || [];
                    this.search.results = res.userList || [];
                    this.replaceShowVal('sno');
                });
                this.search.group_name = grouopName;
                return
            }
            // 添加分组
            this.search.group_name = "";
            this.search.members = [];
            this.search.results = [];
            this.modal.title = "添加分组";
        },
        showNameModal(item){
            this.currentGroup.id = item.id;
            this.currentGroup.group_name = item.group_name;
            this.putNameModal = true;
        },
        onOkGroupModal() {
            if (this.modal.isChangeMember) {
                this.putGroupMember();
            }
            else {
                this.addGroup();
            }
        },
        putGroupMember() {
            this.loading = true;
            const obj = {
                'id': this.search.id,
                'member_list': intArrayToStr(this.search.members) || ""
            };
            this.$service.MainAPI.putGroupMember(obj).then(res => {
                this.modal.show = false;
                this.$Message.info('修改成功');
            }).finally(() => {
                this.loading = false;
            });
        },
        addGroup() {
            this.loading = true;
            const obj = {
                'creator': this.getMeetingsParams.creator || -1,
                'group_name': this.search.group_name || "",
                'member_list': intArrayToStr(this.search.members) || ""
            };
            // 检查输入值
            if (obj.creator === -1) {
                this.$Message.error('请重新登录...');
                this.loading = false;
                return;
            } else if (obj.group_name === "") {
                this.$Message.error('分组名不能为空');
                this.loading = false;
                return;
            }
            this.$service.MainAPI.addGrouop(obj).then(res => {
                this.modal.show = false;
                this.$Message.info('添加成功');
                this.getDataList();
            }).finally(() => {
                this.loading = false;
            });
        },
        searchUsers(query) {
            if (query.trim() !== "") {
                if (!this.search.loading) {
                    // 实现input连续输入，只发一次请求
                    clearTimeout(this.timeout);
                    this.timeout = setTimeout(() => {
                        this.search.loading = true;
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
        changeSearchWay() {
            this.$refs.searchInput.setQuery('');
        },
    },
    created() {
        this.getDataList();
    }
}
</script>