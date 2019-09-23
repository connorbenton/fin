'use strict';

/** gives us access to methods such as 
- findAll()
- create()
- update() 
- destroy()
*/

module.exports = function(sequelize, DataTypes) {
  var Category = sequelize.define('Category', {
    topCategory: DataTypes.TEXT,
    subCategory: DataTypes.TEXT,
    excludeFromAnalysis: DataTypes.BOOLEAN
  }, {
    //set the timestamps to be underscored: (created_at, updated_at)
    underscored: true,
  });
  return Category;
};