'use strict';

/** gives us access to methods such as 
- findAll()
- create()
- update() 
- destroy()
*/

module.exports = function(sequelize, DataTypes) {
  var Account = sequelize.define('Account', {
    name: DataTypes.STRING,
    institution: DataTypes.STRING,
    account_id: { type: DataTypes.STRING, unique: 'compositeIndex'},
    item_id: DataTypes.STRING,
    type: DataTypes.STRING,
    subtype: DataTypes.STRING,
    balance: DataTypes.FLOAT,
    limit: DataTypes.FLOAT,
    available: DataTypes.FLOAT,
    currency: DataTypes.STRING,
    provider: { type: DataTypes.STRING, unique: 'compositeIndex'},
    runningTotal: DataTypes.FLOAT
  }, {
    //set the timestamps to be underscored: (created_at, updated_at)
    underscored: true,
    classMethods: {
      associate: function(models) {
        //An author can have many books.
        Account.hasMany(models.Transaction, {
          onDelete: 'cascade' // when author is deleted, delete their books
        });
        Account.belongsTo(models.ItemToken);
      }
    }
  });
  return Account;
};