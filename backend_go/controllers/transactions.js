Transaction = require('../models').Transaction;
var Server = require('../index.js');
CurrencyRate = require('../models').CurrencyRate;
const seq = require('../models/index');
const { QueryTypes } = require('sequelize');
const axios = require('axios');
var ioSock;
var realSock;

module.exports = {

  function(io) {
    ioSock = io;
  },
  function(socket) {
    realSock = socket;
  },
  //Speed up the Transactions response by using raw SQL
  async index(req, res) {
    try {
    // let tx = await Transaction.findAll();
    // const { QueryTypes } = require('sequelize');
    const tx = await seq.sequelize.query("SELECT * FROM `transactions`", {type: QueryTypes.SELECT});
    // //testing out direct SQL queries
    // const tx = await seq.sequelize.query("SELECT * FROM `transactions` WHERE date > date('now', '-30 days')", {type: QueryTypes.SELECT});
    res.status(200).json(tx);
    }
    catch (error) {
        res.status(500).json(error);
    }
      // .then(function (transactions) {
        // res.status(200).json(transactions);
      // })
      // .catch(function (error) {
      // });
  },
  async indexRange(req, res) {
    try {
    // let tx = await Transaction.findAll();
    // const { QueryTypes } = require('sequelize');
    // const tx = await seq.sequelize.query("SELECT * FROM `transactions`", {type: QueryTypes.SELECT});
    // //testing out direct SQL queries
    const tx = await seq.sequelize.query("SELECT * FROM `transactions` WHERE date > date('now', '-30 days')", {type: QueryTypes.SELECT});
    res.status(200).json(tx);
    }
    catch (error) {
        res.status(500).json(error);
    }
      // .then(function (transactions) {
        // res.status(200).json(transactions);
      // })
      // .catch(function (error) {
      // });
  },
  show(req, res) {
    Transaction.findById(req.params.id)
      .then(function (transaction) {
        res.status(200).json(transaction);
      })
      .catch(function (error) {
        res.status(500).json(error);
      });
  },

  create(req, res) {
    Transaction.create(req.body)
      .then(function (newTransaction) {
        res.status(200).json(newTransaction);
      })
      .catch(function (error) {
        res.status(500).json(error);
      });
  },

  async import(req, res) {

    try {

      let dbRates = await CurrencyRate.findAll();
      let rateJSON = dbRates[0].dataJSON.rates;
      let latestDate = Object.keys(rateJSON).reduce((a, b) => { return new Date(a) > new Date(b) ? a : b });

      let startDate = new Date(latestDate);
      startDate.setDate(startDate.getDate() + 1);
      startDate = new Date(startDate.getTime() - (startDate.getTimezoneOffset() * 60 * 1000)).toISOString().split('T')[0];

      let endDate = new Date();
      endDate.setDate(endDate.getDate() + 1);
      endDate = new Date(endDate.getTime() - (endDate.getTimezoneOffset() * 60 * 1000)).toISOString().split('T')[0];

      let currencyAPIurl = "https://api.exchangeratesapi.io/history?start_at=" + startDate + "&end_at=" + endDate + "&base=USD";

      await axios.get(currencyAPIurl).then(response => {
        var rates = response.data.rates;
        for (let [key, value] of Object.entries(rates)) {
          rateJSON[key] = value;
        }
        // console.log(response)
      })
        .catch(error => {
          console.err(error)
        });

      let newRates = {};
      newRates.dataJSON = {};
      newRates.dataJSON.rates = rateJSON;
      newRates.id = 1;

      dbRates[0].dataJSON.rates = rateJSON;
      let ratesToCheck = rateJSON;
      await CurrencyRate.upsert(newRates);
      let importArray = req.body;

      let dbCats = {};
      dbCats = await Category.findAll();
      let dbTrans = {};
      dbTrans = await Transaction.findAll();
      let dbAccounts = {};
      dbAccounts = await Account.findAll();

      const mintComparisonIndex = [
        { "MintName": "Buy", "catID": 32 },
        { "MintName": "Investments", "catID": 99 },
        { "MintName": "Dividend & Cap Gains", "catID": 64 },
        { "MintName": "Sports", "catID": 52 },
        { "MintName": "Coffee Shops", "catID": 38 },
        { "MintName": "Gift", "catID": 43 },
        { "MintName": "Shipping", "catID": 76 },
        { "MintName": "Hide from Budgets & Trends", "catID": 109 },
        { "MintName": "Business Services", "catID": 76 },
        { "MintName": "Home Improvement", "catID": 55 },
        { "MintName": "Home Services", "catID": 55 },
        { "MintName": "Withdrawal", "catID": 107 },
        { "MintName": "Office Supplies", "catID": 58 },
        { "MintName": "Pets", "catID": 76 },
        { "MintName": "Credit Card Payment", "catID": 100 },
        { "MintName": "Brokerage", "catID": 99 },
        { "MintName": "Interest Income", "catID": 64 },
        { "MintName": "Brokerage Investment", "catID": 99 },
        { "MintName": "Printing", "catID": 76 },
      ];

      let identifiedAccounts = [];
      let transArray = [];
      // let date;
      // const offset = date.getTimezoneOffset();

      // First runthrough, generate transaction IDs and identify accounts
      for (key in importArray) {
        let importTransaction = importArray[key];
        let transactionToUpload = {};

        if (importTransaction.amount == null) continue;

        // Generate random transaction ID
        let newID;
        do {
          newID = Math.random().toString(36).substr(2, 8);
        } while (dbTrans.find(x => x.transaction_id === newID) != null |
          transArray.find(x => x.transaction_id === newID) != null);
        transactionToUpload.transaction_id = newID;

        let date = new Date(importTransaction.date);
        date = new Date(date.getTime() - date.getTimezoneOffset() * 60 * 1000);

        transactionToUpload.date = date.toISOString().slice(0, 10);
        transactionToUpload.description = importTransaction.description;
        transactionToUpload.amount = (importTransaction.transactionType == "debit") ?
          importTransaction.amount * -1 : importTransaction.amount;
        transactionToUpload.currency_code = (importTransaction.currency_code == null) ?
          "USD" : importTransaction.currency_code;

        let i = 0; 
        let checkCurrencyDate = transactionToUpload.date;
        while (!ratesToCheck.hasOwnProperty(checkCurrencyDate)) {
          if (i > 20) {
            throw Error('could not find matching currency data for ' + transactionToUpload.date + ' within 20 days');
          }
          i++;
          checkCurrencyDate = new Date(checkCurrencyDate);
          checkCurrencyDate.setDate(checkCurrencyDate.getDate() - 1);
          checkCurrencyDate = new Date(checkCurrencyDate.getTime() - (checkCurrencyDate.getTimezoneOffset() * 60 * 1000)).toISOString().split('T')[0];
        }
        let compareDate = ratesToCheck[checkCurrencyDate];
        if (!compareDate.hasOwnProperty(transactionToUpload.currency_code.toUpperCase())) {
          throw Error('Could not find matching currency data for ' + transactionToUpload.currency_code + ' on ' + checkCurrencyDate);
        }
        transactionToUpload.normalized_amount = parseFloat(transactionToUpload.amount) / parseFloat(compareDate[transactionToUpload.currency_code.toUpperCase()]);

        //Find matching category
        if (!importTransaction.category || 0 === importTransaction.category.length) {
          transactionToUpload.category = 106;
        }
        else {
          let cat = mintComparisonIndex.find(x => x.MintName === importTransaction.category);
          if (cat != null) {
            transactionToUpload.category = cat.catID;
          }
          else {
            cat = dbCats.find(x => x.subCategory === importTransaction.category);
            if (cat == null) {
              cat = dbCats.find(x => x.topCategory === importTransaction.category);
            }
            if (cat == null) {
              // transactionToUpload.category = 106;
              transactionToUpload.category = importTransaction.category;
            }
            else {
              transactionToUpload.category = cat.id;
            }
          }
        }

        //Create account matching for imported account names
        let possibleMatches = dbTrans.filter(x => x.amount == transactionToUpload.amount);
        if (possibleMatches.length > 0) {
          let matchToCheck = possibleMatches.filter(x => x.date == transactionToUpload.date);
          if (matchToCheck.length > 1) {
            // throw "More than 2 matching transactions found for " + transactionToUpload.description +
            // " on " + transactionToUpload.date + " for " + transactionToUpload.amount;
            console.log("More than 2 matching transactions found for " + transactionToUpload.description +
              " on " + transactionToUpload.date + " for " + transactionToUpload.amount);
            // console.log("more than 2 matching");
          }
          // else if (matchToCheck.length === 1) {
          // transactionToUpload.removeMe = "1";
          let matchingAccount;
          for (key in matchToCheck) {
            // let matchingAccount = identifiedAccounts.find(x => x.RefAccount === matchToCheck[key].account_id); 
            matchingAccount = identifiedAccounts.find(x => x.RefAccount === matchToCheck[key].account_id);
          }
          if (matchingAccount == null) {
            for (key in matchToCheck) {
              let compareSet = {};
              compareSet.trans1 = (({ date, description, amount, currency_code }) => ({ date, description, amount, currency_code }))(transactionToUpload);
              compareSet.trans2 = (({ date, description, amount, currency_code }) => ({ date, description, amount, currency_code }))(matchToCheck[key]);
              compareSet.type = 'trans';
              await new Promise(resolve => {
                Server.socket.emit('compare', compareSet, (answer) => {
                  console.log('answer is: ' + answer);
                  if (answer === 'yes') {
                    let matchIndex = {};
                    matchIndex.ImportKey = importTransaction.accountName;
                    matchIndex.RefAccount = matchToCheck[key].account_id;
                    identifiedAccounts.push(matchIndex);
                  }
                  resolve(answer);
                });
              });
            }
          }
        }
        // Server.io.off('compareResponse');
        transactionToUpload.account_id = importTransaction.accountName;
        transArray.push(transactionToUpload);
      }

      //Identify missing categories with user
      // let uniqueCats = transArray.map(obj => {(({category}) => ({category}))(obj)});
      let uniqueCats = [];
      for (var item of transArray) {
        let pushItem = (({ category }) => ({ category }))(item);
        uniqueCats.push(pushItem);
      }
      uniqueCats = uniqueCats.filter(x => (isNaN(x.category)));
      uniqueCats = uniqueCats.filter((elem, index, self) => uniqueCats.findIndex(
        (t) => { return (t.category === elem.category) }) === index)

      let catsToTransmit = {};
      catsToTransmit.compareCats = uniqueCats;
      catsToTransmit.type = 'cats';
      catsToTransmit.dbCats = dbCats;

      let returnedCats = [];
      await new Promise(resolve => {
        Server.socket.emit('compare', catsToTransmit, (answer) => {
          returnedCats = answer;
          resolve(answer);
        });
      });

      let accountsToCreate = [];
      // Second runthrough, update manually ID'ed categories, fix account names, actually get duplicates here, get array of 'import' accounts to create
      for (key in transArray) {
        let transaction = transArray[key];

        let matchCat = returnedCats.find(x => x.category === transaction.category);
        if (matchCat != null) {
          if (matchCat.assignedCat != null) {
            transaction.category = matchCat.assignedCat;
          }
          else {
            transaction.category = 106;
          }
        }

        let matchAccount = accountsToCreate.find(x => x.name === transaction.account_id);
        if (matchAccount != null) {
          transaction.account_id = matchAccount.account_id;
        }
        else {
          let matchAccount = identifiedAccounts.find(x => x.ImportKey === transaction.account_id);
          if (matchAccount == null) {
            let accountToCreate = {};
            let newID;
            do {
              newID = Math.random().toString(36).substr(2, 8);
            } while (dbAccounts.find(x => x.account_id === newID) != null |
              accountsToCreate.find(x => x.account_id === newID) != null);
            accountToCreate.account_id = newID;
            accountToCreate.name = transaction.account_id;
            accountToCreate.institution = "Import";
            accountToCreate.provider = "Import";
            accountsToCreate.push(accountToCreate);
            transaction.account_id = newID;
          } else {
            transaction.account_id = matchAccount.RefAccount;
          }
        }
        let possibleMatches = dbTrans.filter(x => x.amount == transaction.amount);
        if (possibleMatches.length > 0) {
          let matchToCheck = possibleMatches.filter(x => x.date == transaction.date);
          if (matchToCheck.length > 0) {
            let accountMatch = matchToCheck.filter(x => x.account_id === transaction.account_id);
            if (accountMatch.length > 0) transaction.removeMe = "1";
          }
        }

      }

      for (key in accountsToCreate) {
        await Account.create(accountsToCreate[key]);
      }

      console.log("uncategorized number = " + transArray.filter(x => x.category == 106).length);
      // Remove duplicates and upload
      transArray = transArray.filter(x => x.removeMe == null)
      await Transaction.bulkCreate(transArray)
        .then(function (newTransaction) {
          res.status(200).json(newTransaction);
        })
        .catch(function (error) {
          res.status(500).json(error);
        });
    }
    catch (err) {
      console.error(err);
      res.status(500).json(err)
    }
  },

  update(req, res) {
    Transaction.update(req.body, {
      where: {
        id: req.body.id
      }
    })
      .then(function (updatedRecords) {
        res.status(200).json(updatedRecords);
      })
      .catch(function (error) {
        res.status(500).json(error);
      });
  },

  delete(req, res) {
    Transaction.destroy({
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