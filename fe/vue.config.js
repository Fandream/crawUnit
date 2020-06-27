module.exports = {
  publicPath: process.env.BASE_URL || '/',
  devServer: {
    proxy: {
      '/api/': {
        target: 'http://localhost:8000',
        ws: true,
        changeOrigin: true,
        pathRewrite: { '^/api': '' }
      }
    }
  },
  chainWebpack: config => {
    config.module
      .rule('eslint')
      .use('eslint-loader')
      .loader('eslint-loader')
      .tap(options => {
        options.fix = true
        return options
      })
  }
};
