Plaid_Category = require('../models').Plaid_Category;

module.exports= {
  index(req, res) {
    Plaid_Category.findAll()
      .then(function (categories) {
        res.status(200).json(categories);
      })
      .catch(function (error) {
        res.status(500).json(error);
      });
  },

  show(req, res) {
    Plaid_Category.findById(req.params.id)
    .then(function (category) {
      res.status(200).json(category);
    })
    .catch(function (error){
      res.status(500).json(error);
    });
  },

  create(req, res) {
    Plaid_Category.create(req.body)
      .then(function (newPlaid_Category) {
        res.status(200).json(newPlaid_Category);
      })
      .catch(function (error){
        res.status(500).json(error);
      });
  },

  update(req, res) {
    Plaid_Category.update(req.body, {
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
    Plaid_Category.destroy({
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