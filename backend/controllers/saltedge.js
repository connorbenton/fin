//import Vue from 'vue'
const axios = require('axios');
const fs = require('fs');
const moment = require('moment');
//import axios from 'axios'
//import rateLimit from 'axios-rate-limit'
ItemToken = require('../models').ItemToken;
Account = require('../models').Account;
Transaction = require('../models').Transaction;
CurrencyRate = require('../models').CurrencyRate;
Category = require('../models').Category;
SaltEdge_Category = require('../models').SaltEdge_Category;

const client = axios.create({
  //  maxRequests: 50,
  //  perMilliseconds: 10,
  //baseURL: "https://cors-anywhere.herokuapp.com/https://www.saltedge.com/api/v5",
  baseURL: "https://www.saltedge.com/api/v5",
  json: true
});

let dbRates;

async function execute(method, resource, data) {
  // inject the accessToken for each request
  //let accessToken = await Vue.prototype.$auth.getAccessToken()
  return client({
    method,
    url: resource,
    data,
    headers: {
      "Accept": "application/json",
      "Content-Type": "application/json",
      "App-id": process.env.VUE_APP_SALTEDGE_APP_ID,
      "Secret": process.env.VUE_APP_SALTEDGE_APP_SECRET,
      "Expires-at": Math.floor(new Date().getTime() / 1000 + 60)
    }
  }).then(req => {
    return req.data
  })
}

function AccountCreate(institution, data, accounts) {
  // console.log(institution);
  var data = data.data;
  for (key in data) {
    var saltAcc = data[key];
    var acc = {};
    // console.log(url);
    acc.name = saltAcc.extra.account_name || saltAcc.name;
    acc.institution = institution;
    acc.account_id = saltAcc.id;
    acc.item_id = saltAcc.connection_id;
    acc.id = accounts.find(x => (x.provider === "SaltEdge" && x.item_id === acc.item_id)).id;
    acc.type = saltAcc.nature;
    acc.limit = saltAcc.extra.credit_limit;
    acc.available = saltAcc.extra.available_amount;
    acc.balance = saltAcc.balance;
    acc.provider = "SaltEdge";
    acc.currency = saltAcc.currency_code;
    Account.upsert(acc);
  }
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
  async fetchTransactions(connection) {

    let ratesToCheck = dbRates[0].dataJSON.rates;

    var url = '/transactions?connection_id=' + connection;
    // console.log("checking ", url)
    // let cats = [];
    // let s_cats = [];
    // let cats = await Category.findAll();
    let s_cats = await SaltEdge_Category.findAll();
    // await Category.findAll()
    // .then(function (categories) {
    //   cats = json(categories);
    // });

    try {
      let r = await execute('get', url);
      var data = r.data;
      let transArray = [];
      for (key in data) {
        var transToUpload = {};
        var transFromSaltEdge = data[key];
        transToUpload.date = transFromSaltEdge.extra.posting_date || transFromSaltEdge.made_on;
        transToUpload.description = transFromSaltEdge.description;
        transToUpload.amount = transFromSaltEdge.amount;
        transToUpload.currency_code = transFromSaltEdge.currency_code;

        let i = 0;
        let checkCurrencyDate = transToUpload.date;
        normalizeBlock: {
          while (!ratesToCheck.hasOwnProperty(checkCurrencyDate)) {
            if (i > 20) {
              // throw Error('could not find matching currency data for ' + transToUpload.date + ' within 20 days');
              console.error('could not find matching currency data for ' + transToUpload.date + ' within 20 days');
              transToUpload.normalized_amount = 0;
              break normalizeBlock;
            }
            i++;
            checkCurrencyDate = new Date(checkCurrencyDate);
            checkCurrencyDate.setDate(checkCurrencyDate.getDate() - 1);
            checkCurrencyDate = new Date(checkCurrencyDate.getTime() - (checkCurrencyDate.getTimezoneOffset() * 60 * 1000)).toISOString().split('T')[0];
          }
          let compareDate = ratesToCheck[checkCurrencyDate];
          if (!compareDate.hasOwnProperty(transToUpload.currency_code.toUpperCase())) {
            // throw Error('Could not find matching currency data for ' + transToUpload.currency_code + ' on ' + checkCurrencyDate);
            console.error('Could not find matching currency data for ' + transToUpload.currency_code + ' on ' + checkCurrencyDate);
            transToUpload.normalized_amount = 0;
            break normalizeBlock;
          }
          transToUpload.normalized_amount = parseFloat(transToUpload.amount) / parseFloat(compareDate[transToUpload.currency_code.toUpperCase()]);
        }

        // transToUpload.transactionType = transFromSaltEdge.description; //Should deprecate if amount can be made positive/negative
        // var salt_cat = await SaltEdge_Category.findOne({ where: {bottomCategory: transFromSaltEdge.category}});
        let salt_cat = s_cats.find(x => x.bottomCategory === transFromSaltEdge.category);
        if (salt_cat == null) {
          salt_cat = s_cats.find(x => x.subCategory === transFromSaltEdge.category);
          // salt_cat = await SaltEdge_Category.findOne({ where: {subCategory: transFromSaltEdge.category}});
        }
        if (salt_cat == null) {
          transToUpload.category = 106;
        }
        else {
          // let real_cat = cats.find(x => x.id === salt_cat.linkToAppCat)
          // var real_cat = await Category.findOne({ where: {id: salt_cat.dataValues.linkToAppCat}});
          // transToUpload.category = real_cat.id;
          transToUpload.category = salt_cat.linkToAppCat;
        }
        // transToUpload.accountName = transFromSaltEdge.description; //Should deprecate since it's covered by account_id
        transToUpload.account_id = transFromSaltEdge.account_id;
        transToUpload.transaction_id = transFromSaltEdge.id;
        transArray.push(transToUpload);

        // await Transaction.upsert(transToUpload);
        // await Transaction.findCreateFind({
        //   where: {
        //     transaction_id: transToUpload.transaction_id,
        //     account_id:     transToUpload.account_id
        //   },
        //   // transaction: t
        // })
        // .spread(function(transResult, created){
        //   if (created) {
        //     Transaction.upsert(transToUpload);
        //   }
        //   else {
        //     Transaction.update({transToUpload}, {fields: ['date', 'description', 'amount']});
        //   }
        // })
        // await Transaction.upsert({transToUpload}, {fields: [date, description, amount]});
      }
      await Transaction.bulkCreate(transArray, { updateOnDuplicate: ["date", "description", "amount"] })
        .catch(function (err) {
          console.error(err);
        });
    }
    catch (err) {
      console.error(err);
    }
    // }).catch(data => {
    // console.error(data)
    // });
  },

  createConnectionInteractive(req, res) {
    var url = '/connect_sessions/create'
    var params = {
      "data": {
        "customer_id": process.env.VUE_APP_SALTEDGE_CUSTOMER_ID,
        "consent": {
          // "from_date": moment().format('YYYY-MM-DD'),
          "scopes": [
            "account_details",
            "transactions_details"
          ]
        },
        "attempt": {
          "return_to": process.env.VUE_APP_BASE_URL,
        }
      }
    }
    execute('post', url, params)
      .then(data => {
        // console.log(data);
        var data = data.data;
        res.status(200).json(data.connect_url);
      }).catch(err => {
        // let y = err;
        console.error(err);
        // console.error(err.error_message);
        // console.log(err.request);
      });
  },

  refreshConnectionInteractive(req, res) {
    var url = '/connect_sessions/refresh'
    var params = {
      "data": {
        "connection_id": req.params.id,
        "attempt": {
          "return_to": process.env.VUE_APP_BASE_URL,
        }
      }
    }
    execute('post', url, params)
      .then(data => {
        // console.log(data);
        var data = data.data;
        res.status(200).json(data.connect_url);
      }).catch(data => {
        console.error(data)
      });
  },

  async getConnections(req, res2) {
    // return new Promise(async (resolve, reject) => {
    try {
      dbRates = await CurrencyRate.findAll();
      let rateJSON = dbRates[0].dataJSON.rates;
      let latestDate = Object.keys(rateJSON).reduce((a, b) => { return new Date(a) > new Date(b) ? a : b });

      let startDate = new Date(latestDate);
      startDate.setDate(startDate.getDate() + 1);
      startDate = new Date(startDate.getTime() - (startDate.getTimezoneOffset() * 60 * 1000)).toISOString().split('T')[0];

      let endDate = new Date();
      endDate.setDate(endDate.getDate() + 1);
      endDate = new Date(endDate.getTime() - (endDate.getTimezoneOffset() * 60 * 1000)).toISOString().split('T')[0];

      let currencyAPIurl = "https://api.exchangeratesapi.io/history?start_at=" + startDate + "&end_at=" + endDate + "&base=USD";

      let response = await axios.get(currencyAPIurl);
      var rates = response.data.rates;
      for (let [key, value] of Object.entries(rates)) {
        rateJSON[key] = value;
      }
      // console.log(response)
      // })
      // .catch(error => {
      // console.err(error)
      // });

      let newRates = {};
      newRates.dataJSON = {};
      newRates.dataJSON.rates = rateJSON;
      newRates.id = 1;

      dbRates[0].dataJSON.rates = rateJSON;
      await CurrencyRate.upsert(newRates);

      var url = '/connections?customer_id=' + process.env.VUE_APP_SALTEDGE_CUSTOMER_ID
      // var url = '/transactions?connection_id=7453032'
      // var url = '/transactions?connection_id=7453032'
      // var url = '/connections/7453032'
      // var url = '/providers/ubs_ch'
      // var url = '/connections/5364436/refresh'
      // var url = '/connect_sessions/create'
      // var url = '/customers'
      var params = {
        "data": {
          "customer_id": "4413793",
          "consent": {
            "scopes": [
              "account_details",
              "transactions_details"
            ],
            "from_date": "2018-01-01"
          }
          // "attempt": {
          // "return_to": "http://example.com/"
          // }
        }
      }
      // console.log(params)
      // execute('post', url, params).then(data => {
      let data = await execute('get', url);
      // .then(data => {
      // var data = data.data;
      data = data.data;
      // var itemId;
      var url;
      var promiseArr = [];
      // var institutionArray = [];
      for (key in data) {
        var conn = data[key];
        // console.log(conn);
        // itemId = conn.id;
        var item = {};
        item.institution = conn.provider_name;
        // institutionArray.push(item.institution);
        item.provider = "SaltEdge";
        if (conn.last_attempt.interactive === true) {
          item.interactive = true;
        }
        else item.interactive = false;
        item.lastRefresh = moment(conn.last_success_at);
        let next = conn.next_refresh_possible_at;
        item.nextRefreshPossible = (next === undefined) ? moment() : moment(conn.next_refresh_possible_at);
        item.item_id = conn.id;
        // console.log(acc);
        // console.log('item');
        // console.log(item);
        await ItemToken.upsert(item);
        url = '/accounts?connection_id=' + conn.id;
        // console.log(url);
        // var p = execute('get', url).then(data => AccountCreate(item.institution, data))
        var p = execute('get', url);
        // .catch(data => {
        // console.error(data)
        // });
        promiseArr.push(p);
      }
      // let accounts = await Account.findAll({});
      let values = await Promise.all(promiseArr)
      // .then(values => {
      for (key in values) {
        let account = values[key];
        // AccountCreate(institutionArray[key], account, accounts);
        let data = account.data;
        // function AccountCreate(institution, data, accounts) {
        for (key in data) {
          let saltAcc = data[key];
          let acc = {};
          // console.log(url);
          acc.name = saltAcc.extra.account_name || saltAcc.name;
          // acc.institution = institutionArray[key];
          acc.institution = item.institution;
          acc.account_id = saltAcc.id;
          acc.item_id = saltAcc.connection_id;
          // let match = accounts.find(x => (x.provider === "SaltEdge" && x.item_id === acc.item_id));
          // acc.id = (match === undefined) ? null : match.id; 
          acc.type = saltAcc.nature;
          acc.limit = saltAcc.extra.credit_limit;
          acc.available = saltAcc.extra.available_amount;
          acc.balance = saltAcc.balance;
          acc.provider = "SaltEdge";
          acc.currency = saltAcc.currency_code;
          // console.log(acc);
          // console.log(saltAcc);
          Account.upsert(acc);
        }
      }
      // })
      // .catch(error => {
      // console.error(error)
      // });
      // }).catch((data, res) => {
      //   console.error(data)
      // });

      // resolve(res2.status(200).json("Updated SaltEdge connections"));
      let resTokens = await ItemToken.findAll({ attributes: { exclude: ['access_token'] } });
      let resAccounts = await Account.findAll({});
      let resJSON = { resTokens: resTokens, resAccounts: resAccounts };
      if (res2 !== undefined) res2.status(200).json(resJSON);
      // }
    }
    catch (error) {
      // reject(res2.status(500).json(error));
      if (res2 !== undefined) res2.status(500).json('An error occured');
    }
    // if (success) res.status(200);
    // });
  }
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
