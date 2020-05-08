'use strict';
module.exports = function(sequelize, DataTypes) {
  var Transaction = sequelize.define('Transaction', {
    // date: {type: DataTypes.DATEONLY, get: function() {
      // return moment.utc(this.getDataValue('date')).format('YYYY-MM-DD');
    // }
    // },
    date: DataTypes.DATEONLY,
    // transaction_id: DataTypes.STRING,
    // transaction_id: { type: DataTypes.STRING, unique: 'compositeIndex'},
    transaction_id: { type: DataTypes.STRING, unique: true }, 
    description: DataTypes.TEXT,
    originalDescription: DataTypes.TEXT,
    amount: DataTypes.TEXT,
    normalized_amount: DataTypes.TEXT,
    transactionType: DataTypes.TEXT,
    category: DataTypes.INTEGER,
    accountName: DataTypes.TEXT,
    currency_code: DataTypes.STRING,
    account_id: DataTypes.STRING,
    // account_id: { type: DataTypes.STRING, unique: 'compositeIndex'},
    // account_id: { type: DataTypes.STRING, unique: true },
    labels: DataTypes.TEXT,
    notes: DataTypes.TEXT
  },
  // {
  // indexes: [
  //   {
  //     name: "unique_transaction_per_account",
  //     unique: true,
  //     fields: ["transaction_id", "account_id"]
  //   }
  // ]
  // },
   {
    underscored: true,
    classMethods: {
      associate: function(models) {
        // associations can be defined here
        Transaction.belongsTo(models.Account, {});
      }
    }
  });
  return Transaction;
};