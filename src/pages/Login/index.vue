<template>
    <h1>Hello</h1>
</template>

<style lang="less" scoped>
</style>

<script>
export default {
    name: "Login",
    created() {
        let user = {
	        "phone": "15217899926",
	        "password": "password2",
	        "is_admin": 0
        };
        this.$service.MainAPI.login(user).then((res) => {
            console.log("res -> ", res);
        })
    }
};
</script>