module.exports = {
  devServer: {
    proxy: {
      '/playlist': {
        target: 'http://localhost:5000',
        changeOrigin: true
      },
    }
  }
}
