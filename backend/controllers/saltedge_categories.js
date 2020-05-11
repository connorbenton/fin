SaltEdge_Category = require('../models').SaltEdge_Category;

module.exports= {
  index(req, res) {
    SaltEdge_Category.findAll()
      .then(function (categories) {
        res.status(200).json(categories);
      })
      .catch(function (error) {
        res.status(500).json(error);
      });
  },

  show(req, res) {
    SaltEdge_Category.findById(req.params.id)
    .then(function (category) {
      res.status(200).json(category);
    })
    .catch(function (error){
      res.status(500).json(error);
    });
  },

  create(req, res) {
    SaltEdge_Category.create(req.body)
      .then(function (newSaltEdge_Category) {
        res.status(200).json(newSaltEdge_Category);
      })
      .catch(function (error){
        res.status(500).json(error);
      });
  },

  update(req, res) {
    SaltEdge_Category.update(req.body, {
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
    SaltEdge_Category.destroy({
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