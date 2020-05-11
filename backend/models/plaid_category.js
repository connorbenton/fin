'use strict';

/** gives us access to methods such as 
- findAll()
- create()
- update() 
- destroy()
*/

module.exports = function(sequelize, DataTypes) {
  var Plaid_Category = sequelize.define('Plaid_Category', {
    hierarchy: DataTypes.TEXT,
    catID: {type:DataTypes.TEXT, unique:true},
    linkToAppCat: DataTypes.INTEGER
  }, {
    //set the timestamps to be underscored: (created_at, updated_at)
    underscored: true,
  });
  return Plaid_Category;
};