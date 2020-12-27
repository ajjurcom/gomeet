<template>
    <MenuItem
        class="menu-item"
        v-if="!config.children"
        :name="config.name"
        :to="config.to">
        <Icon :type="config.icon" v-if="config.icon"></Icon>
        {{config.title}}
    </MenuItem>
    <Submenu class="menu-item" :name="config.name" v-else>
        <template slot="title">
            <Icon :type="config.icon" v-if="config.icon"></Icon>
            {{config.title}}
        </template>
        <DynamicMenuItem
            v-for="children of config.children"
            :to="children.to"
            :config="children"
            :key="children.name">
        </DynamicMenuItem>
    </Submenu>
</template>

<script>
    export default {
        name: 'DynamicMenuItem',
        props: {
            config: {
                type: Object,
                default() {
                    return {};
                }
            },
        }
    };
</script>

<style lang="less" scoped>
    .menu-item.ivu-menu-submenu {
        /deep/ .ivu-menu-submenu-title-icon {
            right: 0;
        }
    }
</style>
