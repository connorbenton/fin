Transaction = require('../models').Transaction;

module.exports= {
  index(req, res) {
    Transaction.findAll()
      .then(function (transactions) {
        res.status(200).json(transactions);
      })
      .catch(function (error) {
        res.status(500).json(error);
      });
  },

  show(req, res) {
    Transaction.findById(req.params.id)
    .then(function (transaction) {
      res.status(200).json(transaction);
    })
    .catch(function (error){
      res.status(500).json(error);
    });
  },

  create(req, res) {
    Transaction.create(req.body)
      .then(function (newTransaction) {
        res.status(200).json(newTransaction);
      })
      .catch(function (error){
        res.status(500).json(error);
      });
  },

  bulkCreate(req, res) {
    Transaction.bulkCreate(req.body)
      .then(function (newTransaction) {
        res.status(200).json(newTransaction);
      })
      .catch(function (error){
        res.status(500).json(error);
      });
  },

  update(req, res) {
    Transaction.update(req.body, {
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
    Transaction.destroy({
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