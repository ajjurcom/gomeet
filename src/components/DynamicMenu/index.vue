<template>
    <div class="custom-menu">
        <div class="main">
            <Menu
                ref="menu"
                class="main-menu"
                mode="horizontal"
                theme="light"
                :active-name="currentActiveName"
                @on-select="selectMenu">
                <div class="menu-items">
                    <div class="menu-left">
                        <div class="menu-logo">GOMEET</div>
                        <template v-for="item in menu">
                            <DynamicMenuItem
                                :config="item"
                                :key="item.name"
                                v-if="!item.isRender || item.isRender()">
                            </DynamicMenuItem>
                        </template>
                    </div>
                    <div class="menu-right">
                        <div :class="{'hidden': isAdmin}" class="reserve-info" @click="showInfo">
                            <Icon type="md-volume-up" />
                            预订须知
                        </div>
                        <Submenu :class="{'personal': !isAdmin}" name="personal">
                            <template slot="title">
                                <Icon type="ios-contact" />
                                {{this.$store.getters['App/getUserName']}}
                            </template>
                            <MenuItem
                                name="useredit"
                                :to="{
                                    name: 'UserEdit',
                                    params: {
                                        id: this.$store.getters['App/getUserID']
                                    }
                                }">修改信息</MenuItem>
                                <MenuItem
                                name="usereditpwd"
                                :to="{
                                    name: 'UserEditPwd',
                                    params: {
                                        id: this.$store.getters['App/getUserID']
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
    box-shadow: 0 0 6px rgba(0,0,0,.15);
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
                .menu-right {
                    display: flex;
                    flex-direction: row;
                    .hidden {
                        display: none;
                    }
                    .reserve-info {
                        padding: 0 12px;
                        border-left: 1px solid #dcdee2;
                        color: #9498a0;
                        cursor: pointer;
                        transition: all .2s ease-in-out;
                    }
                    .reserve-info:hover {
                        color: #515a6e;
                    }
                    .personal {
                        border-left: 1px solid #dcdee2;
                        border-right: 1px solid #dcdee2;
                    }
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
    import DynamicMenuItem from './components/DynamicMenuItem';
    import {removeLocalStorage} from '@/Utils';
    export default {
        name: 'DynamicMenu',
        props: {
            /**
             * menu 菜单配置文件 [Array]
             *     item [Object]
             *         title 标题 [String]
             *         to 链接(外链则加上协议头) [String]
             *         icon 图标 [String]
             *         children 子菜单 [Array]
             */
            menu: {
                type: Array,
                default() {
                    return [];
                }
            },
            isAdmin: {
                type: Boolean,
                required: true
            }
        },
        components: {
            DynamicMenuItem
        },
        data() {
            return {
                currentActiveName: this.$route.name.toLowerCase()
            };
        },
        watch: {
            '$route.name': 'onRouteChange',
        },
        methods: {
            selectMenu(name) {
                if (name === "personal-signOut") {
                    this.$Modal.confirm({
                        title: '退出',
                        content: '<p>是否确认退出</p>',
                        onOk: () => {
                            this.$Message.info('退出成功');
                            this.$store.commit('App/setCurrentRole', 'guest');
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
            },
            showInfo() {
                const title = '通知';
                const content = '<p>通知内容</p>';
                this.$Modal.info({
                    title: title,
                    content: content
                });
            },
            onRouteChange(route) {
                this.$nextTick(() => {
                    console.log('this.$route.name -> ', this.$route.name);
                    this.currentActiveName = this.$route.name.toLowerCase();
                    // this.$refs.menu.updateOpened();
                    this.$refs.menu.updateActiveName();
                });
            }
        },
    };
</script>
