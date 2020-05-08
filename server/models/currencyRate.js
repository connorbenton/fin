'use strict';

/** gives us access to methods such as 
- findAll()
- create()
- update() 
- destroy()
*/

module.exports = function(sequelize, DataTypes) {
  var CurrencyRate = sequelize.define('CurrencyRate', {
    dataJSON: DataTypes.JSON,
  }, {
    //set the timestamps to be underscored: (created_at, updated_at)
    underscored: true,
  });
  return CurrencyRate;
};