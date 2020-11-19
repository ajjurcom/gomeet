<template>
    <div class="container-wrap">
        <div class="container">
            <Form :model="campusData" :label-width="80">
                <FormItem label="校区名字">
                    <Input v-model="campusData.campus_name" placeholder="校区"></Input>
                </FormItem>
                <FormItem>
                    <Button type="primary" :loading="loading" @click="changeCampus">添加</Button>
                </FormItem>
            </Form>
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
    }
}
</style>

<script>
export default {
    name: 'CampusAdd',
    data() {
        return {
            loading: false,
            campusData: {
                campus_name: ''
            }
        }
    },
    methods: {
        changeCampus() {
            if (this.campusData.campus_name === '') {
                this.$Message.error("名字不能为空");
                return;
            }
            this.loading = true;
            this.$service.MainAPI.addCampus(this.campusData).then(res => {
                this.$Message.info("添加成功");
                this.$router.push({name:'CampusManager'})
            }).finally(() => {
                this.loading = false;
            });
        }
    }
}
</script>