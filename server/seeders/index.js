// Seed the data
const sequelize_fixtures = require('sequelize-fixtures');

var path = require('path');
var models = require('../models/');
// var saltedge = require ('../controllers/saltedge');


sequelize_fixtures.loadFiles([
	'./server/seeders/categories.js', 
	'./server/seeders/saltedge_categories.js',
	'./server/seeders/plaid_categories.js',
	'./server/seeders/currencyRates.js'
], models).then(function(){
	// saltedge.getCategories();
	console.log('Seed data loaded!');
});