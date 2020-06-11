// const path = require('path')
// const WorkerPlugin = require('worker-plugin')
// const { GenerateSW } = require("workbox-webpack-plugin");
// const VuetifyLoaderPlugin = require('vuetify-loader/lib/plugin')

const host = '0.0.0.0'
const port = 9028

module.exports = {
  configureWebpack: {
      output: {
        globalObject: "this"
      },
      // plugins: [new GenerateSW()]
      // plugins: [
      //   new WorkerPlugin(),
      //   // new VuetifyLoaderPlugin()
      // ]
    },
    // chainWebpack: config => {
    //   config.plugins.delete('prefetch')
    //   // config.plugins.delete('preload')
    // },
  pwa : {
    iconPaths: {
      favicon32: 'img/icons/favicon-32x32.png',
      favicon16: 'img/icons/favicon-16x16.png',
      appleTouchIcon: 'img/icons/apple-touch-icon-152x152.png',
      maskIcon: 'img/icons/safari-pinned-tab.svg',
      msTileImage: 'img/icons/msapplication-icon-144x144.png'
    }  
  },
  devServer: {
    public: process.env.BASE_URL,
    port,
    host,
    hotOnly: true,
    disableHostCheck: true,
    clientLogLevel: 'warning',
    inline: true,
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, PATCH, OPTIONS',
      'Access-Control-Allow-Headers': 'Origin, Accept, X-Requested-With, content-type, Authorization'
    },
    proxy: {
      '^/api': {
        target: 'http://fintrack-go:6060',
        secure: false,
        ws: false,
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