<template>
    <div class="container-wrap">
        <div class="container">
            <Form :model="campusData" :label-width="80">
                <FormItem label="新名字">
                    <Input v-model="campusData.campus_name" placeholder="校区"></Input>
                </FormItem>
                <FormItem>
                    <Button type="primary" :loading="loading" @click="changeCampus">修改</Button>
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
    name: 'CampusEdit',
    data() {
        return {
            loading: false,
            oldCampusName: this.$route.query.campus_name,
            campusData: {
                id: Number(this.$route.query.id),
                campus_name: this.$route.query.campus_name
            }
        }
    },
    methods: {
        changeCampus() {
            if (this.campusData.campus_name === this.oldCampusName) {
                this.$Message.error("值未修改");
                return;
            }
            this.loading = true;
            this.$service.MainAPI.putCampus(this.campusData).then(res => {
                this.$Message.info("修改成功");
                this.$router.push({name:'CampusManager'})
            }).finally(() => {
                this.loading = false;
            });
        }
    }
}
</script>