<template>
    <div class="container-wrap">
        <div class="container">
        </div>
    </div>
</template>

<style lang="less" scoped>
.container-wrap {
    width: 100%;
    min-width: 1200px;
}
</style>

<script>
export default {
    name: "Reserve",
    data () {
        return {
            role : this.$store.getters['App/getCurrentRole']
        };
    },
    methods: {
    },
    created() {
        if (this.role !== 'user') {
            this.$Message.info('请先登录');
            this.$router.push({
                name: "Login"
            });
        }
    }
};
</script>