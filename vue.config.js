module.exports = {
    runtimeCompiler: true,
    outputDir: (function () {
        return {
            development: undefined,
            production: undefined,
            preview: 'dist_test'
        }[process.env.NODE_ENV] || undefined;
    })(),
    assetsDir: undefined,
    productionSourceMap: false,
    parallel: undefined,
    css: undefined,
    configureWebpack: {
        resolve: {
            alias: {
                pages: '@/pages' // 页面组件目录
            }
        },
        output: {
            filename: 'js/[name].[hash].bundle.js'
        },
    }
};