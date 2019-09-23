const path = require('path')

const host = '0.0.0.0'
const port = 9028

module.exports = {

  devServer: {
    public: process.env.VUE_APP_BASE_URL,
    port,
    host,
    hotOnly: true,
    disableHostCheck: true,
    clientLogLevel: 'warning',
    inline: true,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, PATCH, OPTIONS',
      'Access-Control-Allow-Headers': 'X-Requested-With, content-type, Authorization'
    },
    proxy: {
      '/api/': {
        target: 'http://localhost:3000',
        secure: false,
        ws: false,
      }
    },
  }
}