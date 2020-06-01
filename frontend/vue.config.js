const path = require('path')
const WorkerPlugin = require('worker-plugin')
// const VuetifyLoaderPlugin = require('vuetify-loader/lib/plugin')

const host = '0.0.0.0'
const port = 9028

module.exports = {
  configureWebpack: {
      output: {
        globalObject: "this"
      },
      plugins: [
        new WorkerPlugin(),
        // new VuetifyLoaderPlugin()
      ]
    },

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
      '^/api': {
        // target: 'http://localhost:3000',
        // target: 'http://localhost:6060',
        target: 'http://fintrack-go:6060',
        secure: false,
        ws: false,
      },
      // '^/socket.io': {
      '^/ws': {
        // target: 'http://localhost:3000',
        // target: 'http://localhost:6060',
        target: 'http://fintrack-go:6060',
        secure: false,
        ws: true,
      },
      '^/dbgo': {
        target: 'http://0.0.0.0:8085',
        secure: false,
        ws: false,
        changeOrigin: true,
      },
      '^/db': {
        target: 'http://0.0.0.0:8080',
        secure: false,
        ws: false,
        changeOrigin: true,
      }
    },
  }
}