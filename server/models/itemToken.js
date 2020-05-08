'use strict';

/** gives us access to methods such as 
- findAll()
- create()
- update() 
- destroy()
*/

module.exports = function(sequelize, DataTypes) {
  var ItemToken = sequelize.define('ItemToken', {
    // providerName: DataTypes.STRING, //todo maybe later add providername (i.e. Plaid, SaltEdge, etc.)
    institution: DataTypes.STRING,
    access_token: DataTypes.STRING,
    item_id: { type: DataTypes.STRING, unique: 'compositeIndex'},
    provider: { type: DataTypes.STRING, unique: 'compositeIndex'},
    interactive: DataTypes.BOOLEAN,
    needsReLogin: DataTypes.BOOLEAN,
    lastRefresh: DataTypes.DATE,
    nextRefreshPossible: DataTypes.DATE,
    lastDownloadedTransactions: DataTypes.DATE
  }, {
    //set the timestamps to be underscored: (created_at, updated_at)
    underscored: true,
    classMethods: {
      associate: function(models) {
        //An author can have many books.
        ItemToken.hasMany(models.Account, {
          onDelete: 'cascade' // when author is deleted, delete their books
        });
      }
    }
  });
  return ItemToken;
};