<template>
    <div class="custom-menu">
        <div class="main">
            <Menu class="main-menu" mode="horizontal" theme="light" :active-name="activeItem" @on-select="selectMenu">
                <div class="menu-items">
                    <div class="menu-left">
                        <div class="menu-logo">GOMEET</div>
                        <MenuItem
                            name="user"
                            :to="{
                                name: 'UserManager',
                                query: {
                                    state: 'verify_user'
                                }
                            }">
                            <Icon type="logo-octocat" />
                            用户管理
                        </MenuItem>
                        <Submenu name="campus">
                            <template slot="title">
                                <Icon type="ios-school" />
                                校区
                            </template>
                            <MenuItem
                                name="campus-manager"
                                :to="{
                                    name: 'CampusManager'
                                }">管理校区</MenuItem>
                            <MenuItem
                                name="campus-add"
                                :to="{
                                    name: 'CampusAdd'
                                }">新增校区</MenuItem>
                        </Submenu>
                        <Submenu name="building">
                            <template slot="title">
                                <Icon type="ios-pin" />
                                建筑
                            </template>
                            <MenuItem
                                name="building-manager"
                                :to="{
                                    name: 'BuildingManager'
                                }">管理建筑</MenuItem>
                            <MenuItem
                                name="building-add"
                                :to="{
                                    name: 'BuildingAdd'
                                }">新增建筑</MenuItem>
                        </Submenu>
                        <Submenu name="meeting">
                            <template slot="title">
                                <Icon type="md-text" />
                                会议室
                            </template>
                            <MenuItem
                                name="meeting-manager"
                                :to="{
                                    name: 'MeetingManager'
                                }">管理会议室</MenuItem>
                            <MenuItem
                                name="meeting-add"
                                :to="{
                                    name: 'MeetingAdd'
                                }">新增会议室</MenuItem>
                        </Submenu>
                    </div>
                    <div class="menu-right">
                        <Submenu name="personal">
                            <template slot="title">
                                <Icon type="ios-contact" />
                                {{userName}}
                            </template>
                            <MenuItem
                                name="personal-edit"
                                :to="{
                                    name: 'UserEdit',
                                    params: {
                                        id: userID
                                    }
                                }">修改信息</MenuItem>
                                <MenuItem
                                name="personal-password"
                                :to="{
                                    name: 'UserEditPwd',
                                    params: {
                                        id: userID
                                    }
                                }">修改密码</MenuItem>
                            <MenuItem name="personal-signOut">退出</MenuItem>
                        </Submenu>
                    </div>
                </div>
            </Menu>
        </div>
    </div>
</template>

<style lang="less" scoped>
.custom-menu {
    width: 100%;
    min-width: 1200px;
    .main {
        .main-menu {
            width: 100%;
            .menu-items {
                width: 80%;
                margin: 0 auto;
                display: flex;
                flex-direction: row;
                justify-content: space-between;
                .menu-left {
                    display: flex;
                    flex-direction: row;
                }
            }
            .menu-logo {
                height: inherit;
                padding: 0 20px;
                line-height: inherit;
                font-family: serif;
                font-weight: 700;
                font-size: 24px;
                color: #333;
                cursor: pointer;
                transition: all .2s ease-in-out;
            }
            .menu-logo:hover {
                color: #2d8cf0;
            }
        }
    }
}
</style>

<script>
import {removeLocalStorage} from '@/Utils';
export default {
    name: 'CustomMenu',
    data() {
        return {
            userID: this.$store.getters['App/getUserID'],
            userName: this.$store.getters['App/getUserName']
        }
    },
    props: {
        activeItem: {
            type: String,
            default: "user"
        }
    },
    methods: {
        selectMenu(name) {
            if (name === "personal-signOut") {
                this.$Modal.confirm({
                    title: '退出',
                    content: '<p>是否确认退出</p>',
                    onOk: () => {
                        this.$Message.info('退出成功');
                        removeLocalStorage('loginToken');
                        this.$router.push({
                            name: 'Login'
                        });
                    },
                    onCancel: () => {
                        this.$Message.info('取消退出');
                    }
                });
            }
        }
    },
}
</script>