'use strict';

/** gives us access to methods such as 
- findAll()
- create()
- update() 
- destroy()
*/

module.exports = function(sequelize, DataTypes) {
  var SaltEdge_Category = sequelize.define('SaltEdge_Category', {
    topCategory: DataTypes.TEXT,
    subCategory: DataTypes.TEXT,
    bottomCategory: DataTypes.TEXT,
    linkToAppCat: DataTypes.INTEGER
  }, {
    //set the timestamps to be underscored: (created_at, updated_at)
    underscored: true,
  });
  return SaltEdge_Category;
};