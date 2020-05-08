ItemToken = require('../models').ItemToken;
const plaid = require('plaid');
Account = require('../models').Account;
Transaction = require('../models').Transaction;
SaltEdge = require('../controllers/saltedge.js');
const moment = require('moment');
var Server = require('../../index.js');
var ioSock;
var realSock;

// module.exports = function(io) {
// io = io;
// }

const plaidClient = new plaid.Client(
  process.env.VUE_APP_PLAID_CLIENT_ID,
  process.env['VUE_APP_PLAID_SECRET_' + process.env.VUE_APP_ENVIRONMENT.toUpperCase()],
  process.env.VUE_APP_PLAID_PUBLIC_KEY,
  // plaid.environments.sandbox,
  plaid.environments[process.env.VUE_APP_ENVIRONMENT],
  { version: '2018-05-22' }
);

function sleep(ms) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}

module.exports = {

  function(io) {
    ioSock = io;
  },
  function(socket) {
    realSock = socket;
  },
  index(req, res) {
    ItemToken.findAll({ attributes: { exclude: ['access_token'] } })
      .then(function (itemTokens) {
        res.status(200).json(itemTokens);
      })
      .catch(function (error) {
        res.status(500).json(error);
      });
  },

  show(req, res) {
    ItemToken.findById(req.params.id)
      .then(function (itemToken) {
        res.status(200).json(itemToken);
      })
      .catch(function (error) {
        res.status(500).json(error);
      });
  },

  async plaidGeneratePublicToken(req, res) {
    try {
      let toks = await ItemToken.findAll();
      let access = toks.find(x => x.item_id === req.body.item_id).access_token;
      let tok = await plaidClient.createPublicToken(access);
      res.status(200).json(tok);
    }
    catch(err) {
      console.error(err);
      res.status(500).json(err);
    }
  },

  async plaidCreate(req, res2) {
    try {
      // let toks = await ItemToken.findAll();

      // console.log(req.body.token);
      // console.log(req.body.name);

      // plaidClient.exchangePublicToken(req.body.token).then(res => {
      let res = await plaidClient.exchangePublicToken(req.body.token);
      const access_token = res.access_token;
      let item = res;
      // console.log(res);
      let accs = await plaidClient.getAccounts(access_token);
      // console.log("Res 2 ",  res);
      for (key in accs.accounts) {
        var acc = accs.accounts[key];
        var accU = {};
        // console.log(res);
        accU.name = acc.name;
        accU.institution = req.body.name;
        accU.account_id = acc.account_id;
        accU.provider = "Plaid";
        accU.balance = (acc.type == 'credit') ? acc.balances.current * -1 : acc.balances.current;
        accU.limit = acc.balances.limit;
        accU.available = acc.balances.available;
        accU.currency = acc.balances.iso_currency_code;
        accU.type = acc.type;
        accU.subtype = acc.subtype;
        // console.log(req.body);
        accU.item_id = res.item_id;
        // console.log(acc);
        await Account.upsert(accU);
      }
      // });
      item.institution = req.body.name;
      item.provider = "Plaid";
      item.needsReLogin = false;
      // console.log("res item", item);
      await ItemToken.upsert(item);
      // .then(function (newItemToken) {
      res2.status(200).json('Upserted ' + item.institution);
      // })
      // .catch(function (error) {
      // res2.status(500).json(error);
      // });
    }
    catch(err) {

      // }).catch(err => {
      // Indicates a network or runtime error.
      if (!(err instanceof plaid.PlaidError)) {
        res2.sendStatus(500);
        return;
      }

      // Indicates plaid API error
      console.log('/exchange token returned an error', {
        error_type: err.error_type,
        error_code: res.statusCode,
        error_message: err.error_message,
        display_message: err.display_message,
        request_id: err.request_id,
        status_code: err.status_code,
      });

      // Inspect error_type to handle the error in your application
      switch (err.error_type) {
        case 'INVALID_REQUEST':
          // ...
          break;
        case 'INVALID_INPUT':
          // ...
          break;
        case 'RATE_LIMIT_EXCEEDED':
          // ...
          break;
        case 'API_ERROR':
          // ...
          break;
        case 'ITEM_ERROR':
          // ...
          break;
        default:
        // fallthrough
      }

      res2.sendStatus(500);
      // });
    }
    //plaidClient.exchangePublicToken(req.body.token, function(err, res) {
    //   console.log(res);
    //     });
    //   });
  },

  update(req, res) {
    ItemToken.update(req.body, {
      where: {
        id: req.params.id
      }
    })
      .then(function (updatedRecords) {
        res.status(200).json(updatedRecords);
      })
      .catch(function (error) {
        res.status(500).json(error);
      });
  },
  async fetchTransactions(req, res2) {
    try {
      await SaltEdge.getConnections();

      let dbRates = await CurrencyRate.findAll();
      let ratesToCheck = dbRates[0].dataJSON.rates;

      let itemTokens = await ItemToken.findAll();
      // .then(async function (itemTokens) {
      // 
      let promiseArr = [];
      let itemArr = [];
      let sourceArr = [];
      let i = 0;
      let resVal = {
        "len": Object.keys(itemTokens).length,
        "curr": i,
        "name": "SaltEdge"
      };

      // res2.write(JSON.stringify(resVal));
      // Server.io.emit('check', resVal);

      for (key in itemTokens) {
        var item = itemTokens[key];
        resVal.name = "Fetching transactions from " + item.provider + ' - ' + item.institution;
        await new Promise(resolve => {
          Server.socket.emit('check', resVal, (answer) => {
            resolve(answer);
          });
        });
        // Server.io.emit('check', resVal);
        if (item.provider === "SaltEdge") {
          SaltEdge.fetchTransactions(item.item_id);
          if (item.interactive) item.lastDownloadedTransactions = item.lastRefresh;
          else item.lastDownloadedTransactions = moment();
          // console.log(item.dataValues);
          await ItemToken.upsert(item.dataValues);
          // ItemToken.update(
          // { lastDownloadedTransactions: item.lastDownloadedTransactions },
          // { where: { id: item.id } }
          // ).catch(err => { console.error(err) });
        }
        else {
          let now = moment();
          let today = now.format('YYYY-MM-DD');
          let source = item.provider + ' - ' + item.institution;

          try {
            // plaidClient.getAccounts(item.access_token).then(async res => {
            let res = await plaidClient.getAccounts(item.access_token);
            // console.log("Res 2 ",  res);
            let institution = itemTokens.find(x => x.item_id === res.item.item_id).institution;
            for (key in res.accounts) {
              var acc = res.accounts[key];
              var accU = {};
              // console.log(res);
              accU.name = acc.name;
              accU.institution = institution;
              accU.account_id = acc.account_id;
              accU.provider = "Plaid";
              accU.balance = (acc.type == 'credit') ? acc.balances.current * -1 : acc.balances.current;
              accU.limit = acc.balances.limit;
              accU.available = acc.balances.available;
              accU.currency = acc.balances.iso_currency_code;
              accU.type = acc.type;
              accU.subtype = acc.subtype;
              // console.log(req.body);
              accU.item_id = res.item.item_id;
              // console.log(acc);
              await Account.upsert(accU);
            }
            // });
            if (item.lastDownloadedTransactions === null) {
              // let transactionPlaid = plaidClient.getAllTransactions(item.access_token, '2010-01-01', today)
              // transactionPlaid.source = item.provider + ' - ' + item.institution;
              let transactionPlaid = plaidClient.getAllTransactions(item.access_token, '2010-01-01', today)
              // .then(function (result) { console.log(item.institution); return { source: source, value: result}});
              promiseArr.push(transactionPlaid)
              sourceArr.push(source)
              // .then(function (result) { return { source: item.provider + ' - ' + item.institution, value: result}});
            } else {
              // let startDate = now.subtract(30, 'days').format('YYYY-MM-DD');
              let startDate = moment(item.lastDownloadedTransactions).subtract(30, 'days').format('YYYY-MM-DD');
              // let transactionPlaid = plaidClient.getAllTransactions(item.access_token, startDate, today)
              // transactionPlaid.source = item.provider + ' - ' + item.institution;
              let transactionPlaid = plaidClient.getAllTransactions(item.access_token, startDate, today)
              // .then(function (result) { console.log(item.institution); return { source: source, value: result}});
              promiseArr.push(transactionPlaid)
              sourceArr.push(source)
              // .then(function (result) { return { source: item.provider + ' - ' + item.institution, value: result}});
              // promiseArr.push(transactionPlaid).then(function(reqResult) {
              // return {
              // source: item.provider + ' - ' + item.institution,
              // value: reqResult
              // };
              // });
            }
            itemArr.push(item);
          }
          catch (error) {
            // let y = error;
            if (error.error_code === 'ITEM_LOGIN_REQUIRED') {
              item.needsReLogin = true;
              await ItemToken.upsert(item.dataValues);
            }
            else console.error(error)
          }
        }
        // await sleep(2000);
        i++;
        resVal.curr = i;
        // res2.write(JSON.stringify(resVal));
        // Server.io.emit('check', resVal);
        // console.log(JSON.stringify(resVal));
      }
      // Promise.all(promiseArr)
      let results = await Promise.all(promiseArr)
      // .then(async function (results) {

      let cats = {};
      let p_cats = {};
      cats = await Category.findAll();
      p_cats = await Plaid_Category.findAll();

      let transArray = [];
      let data = results;
      // var i = 0;
      for (key in data) {
        var size = Object.keys(data).length;
        // let index = size / 100;
        // var curr = 0;
        // resVal.len = 100;
        resVal.len = size;
        // let limitCheck = 0;
        resVal.curr = key;
        resVal.name = "Loading transactions from " + sourceArr[key];
        await new Promise(resolve => {
          Server.socket.emit('check', resVal, (answer) => {
            resolve(answer);
          });
        });
        // Server.io.emit('check', resVal);
        let plaidTransactions = data[key].transactions;
        for (key in plaidTransactions) {
          let plaidTransaction = plaidTransactions[key];
          let transactionToUpload = {};
          transactionToUpload.date = plaidTransaction.date;
          transactionToUpload.transaction_id = plaidTransaction.transaction_id;
          transactionToUpload.description = plaidTransaction.name;
          transactionToUpload.amount = -1 * plaidTransaction.amount;
          transactionToUpload.currency_code = plaidTransaction.iso_currency_code;

          let i = 0;
          let checkCurrencyDate = transactionToUpload.date;
          normalizeBlock: {
            while (!ratesToCheck.hasOwnProperty(checkCurrencyDate)) {
              if (i > 20) {
                // throw Error('could not find matching currency data for ' + transactionToUpload.date + ' within 20 days');
                console.error('could not find matching currency data for ' + transactionToUpload.date + ' within 20 days');
                transactionToUpload.normalized_amount = 0;
                break normalizeBlock;
              }
              i++;
              checkCurrencyDate = new Date(checkCurrencyDate);
              checkCurrencyDate.setDate(checkCurrencyDate.getDate() - 1);
              checkCurrencyDate = new Date(checkCurrencyDate.getTime() - (checkCurrencyDate.getTimezoneOffset() * 60 * 1000)).toISOString().split('T')[0];
            }
            let compareDate = ratesToCheck[checkCurrencyDate];
            if (!compareDate.hasOwnProperty(transactionToUpload.currency_code.toUpperCase())) {
              // throw Error('Could not find matching currency data for ' + transactionToUpload.currency_code + ' on ' + checkCurrencyDate);
              console.error('Could not find matching currency data for ' + transactionToUpload.currency_code + ' on ' + checkCurrencyDate);
              transactionToUpload.normalized_amount = 0;
              break normalizeBlock;
            }
            transactionToUpload.normalized_amount = parseFloat(transactionToUpload.amount) / parseFloat(compareDate[transactionToUpload.currency_code.toUpperCase()]);
          }

          // let p_cat = await Plaid_Category.findOne({ where: {catID: plaidTransaction.category_id.toString()}});
          let plaid_cat = p_cats.find(x => x.catID === plaidTransaction.category_id.toString());
          if (plaid_cat == null) {
            transactionToUpload.category = 106;
          }
          else {
            // let real_cat = await Category.findOne({ where: {id: p_cat.dataValues.linkToAppCat}});
            // transactionToUpload.category = real_cat.id;
            transactionToUpload.category = plaid_cat.linkToAppCat;
          }
          // transactionToUpload.category = plaidTransaction.category.toString();
          transactionToUpload.account_id = plaidTransaction.account_id;
          // await Transaction.findCreateFind({
          //   where: {
          //     transaction_id: transactionToUpload.transaction_id,
          //     account_id:     transactionToUpload.account_id
          //   },
          //   // transaction: t
          // })
          // .spread(function(transResult, created){
          //   if (created) {
          //     Transaction.upsert(transactionToUpload);
          //   }
          //   else {
          //     Transaction.update({transactionToUpload}, {fields: ['date', 'description', 'amount']});
          //   }
          // })
          // await Transaction.upsert({transToUpload}, {fields: [date, description, amount]});
          // await Transaction.upsert(transactionToUpload);
          // await Transaction.upsert({transactionToUpload}, {fields: [date, description, amount]});
          i++;
          transArray.push(transactionToUpload);
          // if (i >= index) {
          //   curr = (index < 1) ? (curr + (1 / index)) : curr + 1; 
          //   limitCheck = (index < 1) ? (limitCheck + (1 / index)) : limitCheck + 1; 
          //   // curr++;
          //   resVal.curr = curr;
          //   if (limitCheck / 10 < 1) {
          //   Server.io.emit('check', resVal);
          //   limitCheck = 0;
          //   }
          //   i = 0;
          // }
        }
        await Transaction.bulkCreate(transArray, { updateOnDuplicate: ['date', 'description', 'amount'] })
          .catch(function (err) {
            console.error(err);
          });
      }
      for (key in itemArr) {
        let updateItem = itemArr[key];
        updateItem.lastDownloadedTransactions = moment();
        await ItemToken.upsert(updateItem.dataValues);
        // await ItemToken.update(
        // { lastDownloadedTransactions: updateItem.lastDownloadedTransactions },
        // { where: { id: updateItem.id } }
        // ).catch(err => { console.log(err) });
        // );
      }
      res2.status(200).json("Fetched transactions");
      // res2.write("Fetched transactions");
      // res2.end();
      // }).catch(err => {
      // Indicates a network or runtime error.
      // if (!(err instanceof plaid.PlaidError)) {
      // console.log(err);
      // res2.sendStatus(500);
      // return;
      // }
      // })

      // })
      // .catch(function (error) {
      // console.log(error);
      // res2.status(500).json(error);
      // });
    }
    catch (error) {
      console.error(error);
      res2.status(500).json(error);
    }
  },

  delete(req, res) {
    ItemToken.destroy({
      where: {
        id: req.params.id
      }
    })
      .then(function (deletedRecords) {
        res.status(200).json(deletedRecords);
      })
      .catch(function (error) {
        res.status(500).json(error);
      });
  }

};