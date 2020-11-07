<template>
    <div class="custom-menu">
        <div class="main">
            <Menu class="main-menu" mode="horizontal" theme="light" :active-name="activeItem">
                <div class="menu-items">
                    <div class="menu-left">
                        <div class="menu-logo">GOMEET</div>
                        <Submenu name="user">
                            <template slot="title">
                                <Icon type="logo-octocat" />
                                用户
                            </template>
                            <MenuItem name="user-verify">待审核</MenuItem>
                            <MenuItem name="user-blacklist">黑名单</MenuItem>
                            <MenuItem name="user-normal">正常状态</MenuItem>
                        </Submenu>
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
                            <MenuItem name="building-admin">管理建筑</MenuItem>
                            <MenuItem name="building-add">新增建筑</MenuItem>
                        </Submenu>
                        <Submenu name="meeting">
                            <template slot="title">
                                <Icon type="md-text" />
                                会议室
                            </template>
                            <MenuItem name="meeting-manager">管理会议室</MenuItem>
                            <MenuItem name="meeting-add">新增会议室</MenuItem>
                        </Submenu>
                    </div>
                    <div class="menu-right">
                        <Submenu name="5">
                            <template slot="title">
                                <Icon type="ios-contact" />
                                陈铭涛
                            </template>
                            <MenuItem name="5-1">修改信息</MenuItem>
                            <MenuItem name="6-2">退出</MenuItem>
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
export default {
    name: 'CustomMenu',
    props: {
        activeItem: {
            type: String,
            default: "user-normal"
        }
    },
    data() {
        return {
            
        }
    }
}
</script>