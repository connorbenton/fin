'use strict';

/** gives us access to methods such as 
- findAll()
- create()
- update() 
- destroy()
*/

module.exports = function(sequelize, DataTypes) {
  var Account = sequelize.define('Account', {
    title: DataTypes.STRING,
    description: DataTypes.STRING,
    institution: DataTypes.STRING,
    type: DataTypes.STRING,
    runningTotal: DataTypes.INTEGER
  }, {
    //set the timestamps to be underscored: (created_at, updated_at)
    underscored: true,
    classMethods: {
      associate: function(models) {
        //An author can have many books.
        Account.hasMany(models.Transaction, {
          onDelete: 'cascade' // when author is deleted, delete their books
        });
      }
    }
  });
  return Account;
};