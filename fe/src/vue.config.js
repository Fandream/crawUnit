module.exports = {
  publicPath: process.env.BASE_URL || "/",
  devServer: {
    proxy: {
      "/api/": {
        target: "http://localhost:8000",
        ws: true,
        changeOrigin: true,
        pathRewrite: {'^/api' : ''}
      }
    }
  }
};
