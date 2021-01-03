<template>
    <div class="no-data">
        <img
            class="no-data-img"
            :src="imgSrc"
            :alt="alt">
        <div
            :style="{'margin-left': marginLeft + 'px'}"
            class="no-data-text">
            {{title}}
        </div>
    </div>
</template>

<style lang="less" scoped>
    .no-data {
        margin: 100px 0;
        display: flex;
        flex-direction: column;
        align-items: center;
        .no-data-img {
            width: 200px;
            height: 200px;
        }
        .no-data-text {
            width: 200px;
            margin-top: 20px;
            color: #808695;
            font-size: 16px;
            font-weight: 500;
            text-align: center;
        }
    }
</style>

<script>
export default {
    name: 'NoData',
    props: {
        title: {
            type: String,
            required: true
        },
        imgSrc: {
            type: String,
            default: 'http://static.mittacy.com/blog/20201205224934.png'
        },
        alt: {
            type: String,
            default: '暂无数据'
        },
        marginLeft: {
            type: Number,
            default: 20
        }
    }
}
</script>