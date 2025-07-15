module.exports = {
  '/api': {
    target: 'http://127.0.0.1:8022',
    changeOrigin: true,
    pathRewrite: { '^/api': '' }
  }
}
