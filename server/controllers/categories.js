Category = require('../models').Category;

module.exports= {
  index(req, res) {
    Category.findAll()
      .then(function (categories) {
        res.status(200).json(categories);
      })
      .catch(function (error) {
        res.status(500).json(error);
      });
  },

  show(req, res) {
    Category.findById(req.params.id)
    .then(function (category) {
      res.status(200).json(category);
    })
    .catch(function (error){
      res.status(500).json(error);
    });
  },

  create(req, res) {
    Category.create(req.body)
      .then(function (newCategory) {
        res.status(200).json(newCategory);
      })
      .catch(function (error){
        res.status(500).json(error);
      });
  },

  update(req, res) {
    Category.update(req.body, {
      where: {
        id: req.params.id
      }
    })
    .then(function (updatedRecords) {
      res.status(200).json(updatedRecords);
    })
    .catch(function (error){
      res.status(500).json(error);
    });
  },

  delete(req, res) {
    Category.destroy({
      where: {
        id: req.params.id
      }
    })
    .then(function (deletedRecords) {
      res.status(200).json(deletedRecords);
    })
    .catch(function (error){
      res.status(500).json(error);
    });
  }
};