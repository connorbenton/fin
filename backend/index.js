var express = require('express'), http = require('http');
var app = express();
// var io = require('socket.io')(server);
// const router = app.Router();
var server = http.createServer(app);
// var io = require('socket.io').listen(server);
var io = require('socket.io')(server);

const itemTokens = require('./controllers/itemTokens');

bodyParser = require('body-parser'),
  cookieParser = require('cookie-parser'),
  methodOverride = require('method-override'),
  cors = require('cors'),
  router = express.Router(),
  path = require('path'),
  seq = require('./models/index')
  accounts = require('./controllers/accounts'),
  categories = require('./controllers/categories'),
  currencyRates = require('./controllers/currencyRates'),
  plaid_categories = require('./controllers/plaid_categories'),
  saltedge_categories = require('./controllers/saltedge_categories'),
  transactions = require('./controllers/transactions'),
  saltedge = require('./controllers/saltedge');

server.timeout = 60 * 4 * 1000; //4 min for long calls

const port = process.env.PORT || 3000;

const corsOptions = {
  origin: [process.env.VUE_APP_BASE_URL],
  allowedHeaders: ["Content-Type", "Authorization", "Access-Control-Allow-Methods", "Access-Control-Request-Headers"],
  credentials: true,
  enablePreflight: true
};

app.use(cors(corsOptions));
app.options('*', cors(corsOptions))
// cors config
// app.use(cors());

// express config
app.use(bodyParser.json({ limit: '50mb' }));
app.use(bodyParser.urlencoded({ limit: '50mb', extended: true, parameterLimit: 50000 }));
app.use(methodOverride());
app.use(cookieParser());
// app.set('port', process.env.PORT || 3000);

// define accounts routes
router.get('/accounts', accounts.index);
router.get('/accounts/:id', accounts.show);
router.post('/accounts', accounts.create);
router.put('/accounts', accounts.update);
router.delete('/accounts', accounts.delete);
// define Plaid saltedge_categories routes
router.get('/saltedge_categories', saltedge_categories.index);
router.get('/saltedge_categories/:id', saltedge_categories.show);
router.post('/saltedge_categories', saltedge_categories.create);
router.put('/saltedge_categories', saltedge_categories.update);
// define SaltEdge categories routes
router.get('/plaid_categories', plaid_categories.index);
router.get('/plaid_categories/:id', plaid_categories.show);
router.post('/plaid_categories', plaid_categories.create);
router.put('/plaid_categories', plaid_categories.update);
// define categories routes
router.get('/categories', categories.index);
router.get('/categories/:id', categories.show);
router.post('/categories', categories.create);
router.put('/categories', categories.update);
router.delete('/categories', categories.delete);
// define itemTokens routes
router.get('/itemTokens', itemTokens.index);
router.get('/itemTokens/:id', itemTokens.show);
router.post('/plaidItemTokens', itemTokens.plaidCreate);
router.post('/plaidGeneratePublicToken', itemTokens.plaidGeneratePublicToken);
router.put('/itemTokens', itemTokens.update);
router.delete('/itemTokens', itemTokens.delete);
router.get('/itemTokensFetchTransactions', itemTokens.fetchTransactions);
// define transactions routes
router.get('/transactions', transactions.index);
router.get('/transactionsRange/:range', transactions.index);
router.get('/transactions/:id', transactions.show);
router.post('/transactions', transactions.create);
router.put('/transactions', transactions.update);
router.delete('/transactions', transactions.delete);
router.post('/importTransactions', transactions.import);

// define SaltEdge routes
router.get('/saltEdgeConnections', saltedge.getConnections);
router.get('/saltEdgeRefreshInteractive/:id', saltedge.refreshConnectionInteractive);
router.get('/saltEdgeCreateInteractive', saltedge.createConnectionInteractive);

router.get('/resetDB', async function (req, res) {
  try {
    console.log('Resetting database!');
    await seq.sequelize.sync({ force: true });

    const sequelize_fixtures = require('sequelize-fixtures');

    let models = require('./models/');
    // var saltedge = require ('../controllers/saltedge');

    sequelize_fixtures.loadFiles([
      './seeders/categories.js',
      './seeders/saltedge_categories.js',
      './seeders/plaid_categories.js',
      './seeders/currencyRates.js'
    ], models).then(function () {
      // saltedge.getCategories();
      console.log('Database reset and re-seeded!')
      res.status(200).json('Reset database!');
    });
  }
  catch (error) {
    res.status(500).json(error);
  }
});

// register api routes
app.use('/api', router);

// RestartServer(false);

// async function RestartServer (resetDB) {

// if (server.listening) await server.close();

// await seq.sequelize.sync({ force: resetDB});
// seq.sequelize.sync({ force: resetDB})
seq.sequelize.sync()
  .then(() => {
    // console.log(process.env)
    console.log("Database & tables created!")
    // seed the database
    require('./seeders');

    // start server
    server.listen(port, function () {
      console.log("Server started on port", port);
    });

    io.on('connection', (socket) => {
      exports.socket = socket;
      console.log("client connected!");
    });

  })
// }

exports.server = server;
exports.io = io;