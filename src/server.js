//This server not used anymore (using parent directory server)


const express = require('express')
const cors = require('cors')
const bodyParser = require('body-parser')
const Sequelize = require('sequelize')
const finale = require('finale-rest')

//var corsOptions = {
//  origin: function (origin, callback) {
//    if (whitelist.indexOf(origin) !== -1) {
//      callback(null, true)
//    } else {
//      callback(new Error('Not allowed by CORS'))
//    }
//  }
//}

var app = express()
//app.use(cors(corsOptions))
app.use(cors())
app.options('*', cors())
app.use(bodyParser.json())

// For ease of this tutorial, we are going to use SQLite to limit dependencies
var database = new Sequelize({
  dialect: 'sqlite',
  storage: './test.sqlite'
})

// Define our Post model
// id, createdAt, and updatedAt are added by sequelize automatically
var Transaction = database.define('transactions', {
  date: Sequelize.TEXT,
  description: Sequelize.TEXT,
  originalDescription: Sequelize.TEXT,
  amount: Sequelize.TEXT,
  transactionType: Sequelize.TEXT,
  category: Sequelize.TEXT,
  accountName: Sequelize.TEXT,
  labels: Sequelize.TEXT,
  notes: Sequelize.TEXT
})

var Account = database.define('accounts', {
  title: Sequelize.STRING,
  institution: Sequelize.TEXT,
  body: Sequelize.TEXT
})

var Category = database.define('categories', {
  topCategory: Sequelize.TEXT,
  subCategory: Sequelize.TEXT,
  excludeFromAnalysis: Sequelize.BOOLEAN
})

var InitialCategories = [
  "Auto & Transport",
  "Bills & Utilities",
  "Education",
  "Entertainment",
  "Fees & Charges",
  "Financial",
  "Food & Dining",
  "Gifts & Donations",
  "Health & Fitness",
  "Home",
  "Income",
  "Kids",
  "Misc Expenses",
  "Personal Care",
  "Shopping",
  "Taxes",
  "Transfer",
  "Travel",
  "Uncategorized",
  "Hide from Analysis"
]

const InitialCategoryData = [
  ["Auto & Transport", "Auto & Transport",false],
  ["Auto Insurance", "Auto & Transport",false],
  ["Auto Payment", "Auto & Transport",false],
  ["Gas & Fuel", "Auto & Transport",false],
  ["Parking", "Auto & Transport",false],
  ["Public Transportation", "Auto & Transport",false],
  ["Service & Parts", "Auto & Transport",false],
  ["Bills & Utilities", "Bills & Utilities",false],
  ["Home Phone", "Bills & Utilities",false],
  ["Internet", "Bills & Utilities",false],
  ["Mobile Phone", "Bills & Utilities",false],
  ["Television", "Bills & Utilities",false],
  ["Utilities", "Bills & Utilities",false],
  ["Education", "Education",false],
  ["Books & Supplies", "Education",false],
  ["Student Loan", "Education",false],
  ["Tuition", "Education",false],
  ["Entertainment", "Entertainment",false],
  ["Amusement", "Entertainment",false],
  ["Arts", "Entertainment",false],
  ["Movies & DVDs", "Entertainment",false],
  ["Music", "Entertainment",false],
  ["Newspapers & Magazines", "Entertainment",false],
  ["Fees & Charges", "Fees & Charges",false],
  ["ATM Fee", "Fees & Charges",false],
  ["Bank Fee", "Fees & Charges",false],
  ["Finance Charge", "Fees & Charges",false],
  ["Late Fee", "Fees & Charges",false],
  ["Service Fee", "Fees & Charges",false],
  ["Trade Commissions", "Fees & Charges",false],
  ["Financial", "Financial",false],
  ["Brokerage", "Financial",false],
  ["Financial Advisor", "Financial",false],
  ["Life Insurance", "Financial",false],
  ["Roth IRA", "Financial",false],
  ["Food & Dining", "Food & Dining",false],
  ["Alcohol & Bars", "Food & Dining",false],
  ["Fast Food", "Food & Dining",false],
  ["Groceries", "Food & Dining",false],
  ["Restaurants", "Food & Dining",false],
  ["Gifts & Donations", "Gifts & Donations",false],
  ["Charity", "Gifts & Donations",false],
  ["Gifts", "Gifts & Donations",false],
  ["Health & Fitness", "Health & Fitness",false],
  ["Cycling", "Health & Fitness",false],
  ["Dentist", "Health & Fitness",false],
  ["Doctor", "Health & Fitness",false],
  ["Eyecare", "Health & Fitness",false],
  ["Gym", "Health & Fitness",false],
  ["Health Insurance", "Health & Fitness",false],
  ["Pharmacy", "Health & Fitness",false],
  ["Sports", "Health & Fitness",false],
  ["Home", "Home",false],
  ["Furnishings", "Home",false],
  ["Home Improvement", "Home",false],
  ["Home Insurance", "Home",false],
  ["Home Services", "Home",false],
  ["Home Supplies", "Home",false],
  ["Lawn & Garden", "Home",false],
  ["Mortgage & Rent", "Home",false],
  ["Renter's insurance", "Home",false],
  ["Income", "Income",false],
  ["Bonus", "Income",false],
  ["Interest Income", "Income",false],
  ["Paycheck", "Income",false],
  ["Paypal Income", "Income",false],
  ["Reimbursement", "Income",false],
  ["Rental Income", "Income",false],
  ["Returned Purchase", "Income",false],
  ["Kids", "Kids",false],
  ["Allowance", "Kids",false],
  ["Baby Supplies", "Kids",false],
  ["Babysitter & Daycare", "Kids",false],
  ["Kids Activities", "Kids",false],
  ["Toys", "Kids",false],
  ["Misc Expenses", "Misc Expenses",false],
  ["Venmo expense", "Misc Expenses",false],
  ["Wedding", "Misc Expenses",false],
  ["Personal Care", "Personal Care",false],
  ["Hair", "Personal Care",false],
  ["Laundry", "Personal Care",false],
  ["Spa & Massage", "Personal Care",false],
  ["Shopping", "Shopping",false],
  ["Amazon", "Shopping",false],
  ["Amazon Prime Member", "Shopping",false],
  ["Books", "Shopping",false],
  ["Clothing", "Shopping",false],
  ["Coffee", "Shopping",false],
  ["Electronics & Software", "Shopping",false],
  ["Hobbies", "Shopping",false],
  ["Sporting Goods", "Shopping",false],
  ["Taxes", "Taxes",false],
  ["Federal Tax", "Taxes",false],
  ["Local Tax", "Taxes",false],
  ["Property Tax", "Taxes",false],
  ["Sales Tax", "Taxes",false],
  ["State Tax", "Taxes",false],
  ["Transfer", "Transfer",true],
  ["Brokerage Investment", "Transfer",true],
  ["Credit Card Payment", "Transfer",true],
  ["Travel", "Travel",false],
  ["Air Travel", "Travel",false],
  ["Hotel", "Travel",false],
  ["Rental Car & Taxi", "Travel",false],
  ["Vacation", "Travel",false],
  ["Uncategorized", "Uncategorized",false],
  ["Cash & ATM", "Uncategorized",false],
  ["Check", "Uncategorized",false],
  ["Hide from Analysis","Hide from Analysis",true]
]

// Initialize finale
finale.initialize({
  app: app,
  sequelize: database
})

// Create the dynamic REST resource for our Post model
var userResource = finale.resource({
  model: Transaction,
  endpoints: ['/transactions', '/transactions/:id']
})

var userResource2 = finale.resource({
  model: Account,
  endpoints: ['/accounts', '/accounts/:id']
})

var userResource3 = finale.resource({
  model: Category,
  endpoints: ['/categories', '/categories/:id'],
  pagination: false
})

// Resets the database and launches the express app on :8081
database
  .sync({ force: true })
  .then(() => {
    for (arr of InitialCategoryData) {
      var a = arr[1]
      var b = arr[0]
      var c = arr[2]
      Category.create({
        topCategory: a,
        subCategory: b,
        excludeFromAnalysis: c
      }).then()
    }
    app.listen(9029, () => {
      console.log('listening to port localhost:9029')
    })
  })