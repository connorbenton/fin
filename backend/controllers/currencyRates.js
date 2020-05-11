CurrencyRate = require('../models').CurrencyRate;

module.exports= {
  index(req, res) {
    CurrencyRate.findAll()
      .then(function (categories) {
        res.status(200).json(categories);
      })
      .catch(function (error) {
        res.status(500).json(error);
      });
  },

  show(req, res) {
    CurrencyRate.findById(req.params.id)
    .then(function (category) {
      res.status(200).json(category);
    })
    .catch(function (error){
      res.status(500).json(error);
    });
  },

  create(req, res) {
    CurrencyRate.create(req.body)
      .then(function (newCurrencyRate) {
        res.status(200).json(newCurrencyRate);
      })
      .catch(function (error){
        res.status(500).json(error);
      });
  },

  update(req, res) {
    CurrencyRate.update(req.body, {
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
    CurrencyRate.destroy({
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