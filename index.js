const express = require('express'),
app = express(),
bodyParser = require('body-parser'),
cookieParser = require('cookie-parser'),
methodOverride = require('method-override'),
router = express.Router(),
path = require('path'),
seq = require('./server/models/index')
accounts = require('./server/controllers/accounts'),
categories = require('./server/controllers/categories'),
providerTokens = require('./server/controllers/providerTokens'),
transactions = require('./server/controllers/transactions');
saltedge = require ('./server/controllers/saltedge');


// express config
app.use(bodyParser.json({ limit: '10mb'}));
app.use(bodyParser.urlencoded({ extended: true, limit: '10mb' }));
app.use(methodOverride());
app.use(cookieParser());
app.set('port', process.env.PORT || 3000);

// define accounts routes
router.get('/accounts', accounts.index);
router.get('/accounts/:id', accounts.show);
router.post('/accounts', accounts.create);
router.put('/accounts', accounts.update);
router.delete('/accounts', accounts.delete);
// define categories routes
router.get('/categories', categories.index);
router.get('/categories/:id', categories.show);
router.post('/categories', categories.create);
router.put('/categories', categories.update);
router.delete('/categories', categories.delete);
// define providerTokens routes
router.get('/providerTokens', providerTokens.index);
router.get('/providerTokens/:id', providerTokens.show);
router.post('/providerTokens', providerTokens.create);
router.put('/providerTokens', providerTokens.update);
router.delete('/providerTokens', providerTokens.delete);
// define transactions routes
router.get('/transactions', transactions.index);
router.get('/transactions/:id', transactions.show);
router.post('/transactions', transactions.create);
router.put('/transactions', transactions.update);
router.delete('/transactions', transactions.delete);
router.post('/transactionsBulk', transactions.bulkCreate);

// define SaltEdge routes
router.get('/saltEdgeConnections', saltedge.getConnections);

// register api routes
app.use('/api', router);

seq.sequelize.sync({ force: true})
  .then(() => {
    //console.log(process.env)
    console.log("Database & tables created!")
    // seed the database
    require('./server/seeders');

    // start server
    app.listen(app.get('port'), function () {
    console.log("Server started on port", app.get('port'));
    });
  })
