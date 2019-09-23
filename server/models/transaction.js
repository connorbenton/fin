'use strict';
module.exports = function(sequelize, DataTypes) {
  var Transaction = sequelize.define('Transaction', {
    date: DataTypes.TEXT,
    description: DataTypes.TEXT,
    originalDescription: DataTypes.TEXT,
    amount: DataTypes.TEXT,
    transactionType: DataTypes.TEXT,
    category: DataTypes.TEXT,
    accountName: DataTypes.TEXT,
    account_id: DataTypes.INTEGER,
    labels: DataTypes.TEXT,
    notes: DataTypes.TEXT
  }, {
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