Account = require('../models').Account;
Transaction = require('../models').Transaction;


module.exports = {
	index(req, res){
		Account.findAll({
			// return the transactions associated with each account
			include: Transaction
		}).then(function(accounts){
			sendResult(res, accounts);
		}).catch(function(error){
			sendError(res, error);
		});
	},
	show(req, res){
		Account.findById(req.params.id, {
			// return all transactions for account
			include: Transaction
		}).then(function(account){
			sendResult(res, account);
		}).catch(function(error){
			sendError(res, error);
		});
	},
	create(req, res){
		Account.create(req.body)
			.then(function(newAccount){
				sendResult(res, newAccount);
			})
			.catch(function(error){
				sendError(res, error);
			});
	},
	update(req, res){
		Account.update(req.body, {
			where: {
				id: req.params.id
			}
		})
		.then(function(updatedAccount){
			sendResult(res, updatedAccount);
		})
		.catch(function(err){
			sendError(res, err);
		});
	},
	delete(req, res){
		Account.destroy({
			where: {
				id: req.params.id
			}
		})
		.then(function(deletedRecord){
			sendResult(res, deletedRecord);
		})
		.catch(function(err){
			sendError(res, err);
		});
	}
}

// helper functions
function sendResult(res, result){
	res.status(200).json(result);
}

function sendError(res, result){
	res.status(500).json(result);
}
