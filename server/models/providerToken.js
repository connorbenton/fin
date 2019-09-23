'use strict';

/** gives us access to methods such as 
- findAll()
- create()
- update() 
- destroy()
*/

module.exports = function(sequelize, DataTypes) {
  var ProviderToken = sequelize.define('ProviderToken', {
    providerName: DataTypes.STRING,
    accountName: DataTypes.STRING,
    token: DataTypes.STRING,
    lastDownloadedTransactions: DataTypes.DATE
  }, {
    //set the timestamps to be underscored: (created_at, updated_at)
    underscored: true,
  });
  return ProviderToken;
};