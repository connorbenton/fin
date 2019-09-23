//import Vue from 'vue'
const axios = require('axios');
//import axios from 'axios'
//import rateLimit from 'axios-rate-limit'

const client = axios.create({
//  maxRequests: 50,
//  perMilliseconds: 10,
  //baseURL: "https://cors-anywhere.herokuapp.com/https://www.saltedge.com/api/v5",
  baseURL: "https://www.saltedge.com/api/v5",
  json: true
});

async function execute (method, resource, data) {
    // inject the accessToken for each request
    //let accessToken = await Vue.prototype.$auth.getAccessToken()
    return client({
      method,
      url: resource,
      data,
      headers: {
        "Accept":       "application/json",
        "Content-Type": "application/json",
        "App-id":       process.env.VUE_APP_SALTEDGE_APP_ID,
        "Secret":       process.env.VUE_APP_SALTEDGE_APP_SECRET,
        "Expires-at":   Math.floor(new Date().getTime() / 1000 + 60)
      }
    }).then(req => {
      return req.data
    })
  }

module.exports = {
//export default {
//  async execute (method, resource, data) {
//    // inject the accessToken for each request
//    //let accessToken = await Vue.prototype.$auth.getAccessToken()
//    return client({
//      method,
//      url: resource,
//      data,
//      headers: {
//        "Accept":       "application/json",
//        "Access-Control-Allow-Origin": '*',
//        "Content-Type": "application/json",
//        "App-id":       process.env.VUE_APP_SALTEDGE_APP_ID,
//        "Secret":       process.env.VUE_APP_SALTEDGE_APP_SECRET,
//        "Expires-at":   Math.floor(new Date().getTime() / 1000 + 60)
//      }
//    }).then(req => {
//      //console.log(req.data)
//      return req.data
//    })
//  },
  getConnections (req, res) {
    var url = '/connections?customer_id=' + process.env.VUE_APP_SALTEDGE_CUSTOMER_ID
    execute('get', url).then(data => {
      res.json({ message: 'Request received!', data })
    })
  },
//}
};





//var fs          = require("fs");
//var https       = require("https");
//var util        = require("util");
//var crypto      = require("crypto");
//
//export default {
//signedHeaders(url, method, params) {
//  expires_at = Math.floor(new Date().getTime() / 1000 + 60)
//  //payload    = expires_at + "|" + method + "|" + url + "|"
//
//  //if (method == "POST") { payload += JSON.stringify(params) }
//
//  //var privateKey = fs.readFileSync('private.pem');
//  //var signer = crypto.createSign('sha256');
//  //signer.update(payload);
//
//  return {
//    "Accept":       "application/json",
//    "Content-Type": "application/json",
//    "App-id":       process.env.VUE_APP_SALTEDGE_APP_ID,
//    "Secret":       process.env.VUE_APP_SALTEDGE_APP_SECRET,
//    "Expires-at":   expires_at,
//    //"Signature":    signer.sign(privateKey,'base64'),
//  }
//},
//
//// Use this function to verify signature in callbacks
//// https://docs.saltedge.com/account_information/v5/#callbacks-request_identification
////
//// signature - could be obtained from headers['signature']
//// callback_url - url that you add in SE dashboard
//// post_body - request body as string
//verifySignature(signature, callback_url, post_body) {
//  payload = callback_url + "|" + post_body
//
//  var publicKey = fs.readFileSync('../spectre_public.pem');
//  var verifier = crypto.createVerify('sha256');
//  verifier.update(payload);
//
//  return verifier.verify(publicKey, signature,'base64');
//},
//
//request(options) {
//  options.headers = signedHeaders(options.url, options.method, options.data);
//
//  return new Promise((resolve, reject) => {
//    var req = https.request(options.url, options, (response) => {
//      var chunks = [];
//
//      response.on('data', chunk => chunks.push(chunk))
//      response.on('end', ()=> {
//        var data = Buffer.concat(chunks).toString();
//        response.statusCode == 200 ? resolve(data) : reject(data);
//      })
//      response.on('error', ()=> {
//        var data = Buffer.concat(chunks).toString();
//        reject(data);
//      })
//    })
//
//
//    if (options.data && options.method != "GET") {
//      req.write(JSON.stringify(options.data));
//    }
//
//    req.end();
//  });
//}
//}

//url = "https://www.saltedge.com/api/v5/countries"
//
//request({
//  method:  "GET",
//  url:     url
//}).then((data)=> {
//  console.log(data)
//}).catch((data) => {
//  console.error(data)
//})
//
//
//url    = "https://www.saltedge.com/api/v5/customers"
//params = {
//  data: {
//    identifier: "my_unique_sdidentifier"
//  }
//}
//
//request({
//  method:  "POST",
//  url:     url,
//  data:    params
//}).then((data)=> {
//  console.log(data)
//}).catch((data) => {
//  console.error(data)
//})
