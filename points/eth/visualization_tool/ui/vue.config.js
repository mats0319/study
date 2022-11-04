module.exports = {
  publicPath: "./",
  assetsDir: "assets", // 打包后静态资源的位置
  productionSourceMap: false,
  css: {
    extract: false,
  },
  devServer: {
    host: "127.0.0.1",
    port: "9693",
    hot: true,
  },
}
