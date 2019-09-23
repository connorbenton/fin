ProviderToken = require('../models').ProviderToken;

module.exports= {
  index(req, res) {
    ProviderToken.findAll()
      .then(function (providerTokens) {
        res.status(200).json(providerTokens);
      })
      .catch(function (error) {
        res.status(500).json(error);
      });
  },

  show(req, res) {
    ProviderToken.findById(req.params.id)
    .then(function (providerToken) {
      res.status(200).json(providerToken);
    })
    .catch(function (error){
      res.status(500).json(error);
    });
  },

  create(req, res) {
    ProviderToken.create(req.body)
      .then(function (newProviderToken) {
        res.status(200).json(newProviderToken);
      })
      .catch(function (error){
        res.status(500).json(error);
      });
  },

  update(req, res) {
    ProviderToken.update(req.body, {
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
    ProviderToken.destroy({
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