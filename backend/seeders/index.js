// Seed the data
const sequelize_fixtures = require('sequelize-fixtures');

var path = require('path');
var models = require('../models/');
// var saltedge = require ('../controllers/saltedge');


sequelize_fixtures.loadFiles([
	'./seeders/categories.js', 
	'./seeders/saltedge_categories.js',
	'./seeders/plaid_categories.js',
	'./seeders/currencyRates.js'
], models).then(function(){
	// saltedge.getCategories();
	console.log('Seed data loaded!');
});