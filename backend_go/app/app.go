package app

import (
	"github.com/gorilla/mux"

	"fin-go/routes/accounts"
	"fin-go/routes/analysisTrees"
	"fin-go/routes/categories"
	"fin-go/routes/itemTokens"
	"fin-go/routes/plaid"
	"fin-go/routes/resetDB"
	"fin-go/routes/saltedge"
	"fin-go/routes/transactions"
)

type App struct {
	Router *mux.Router
}

func (app *App) SetupRouter() {
	app.Router.
		Methods("GET").
		Path("/api/accounts").
		HandlerFunc(accounts.GetFunction())

	app.Router.
		Methods("POST").
		Path("/api/accountUpsertName").
		HandlerFunc(accounts.UpsertNameFunction())

	app.Router.
		Methods("POST").
		Path("/api/accountUpsertIgnore").
		HandlerFunc(accounts.UpsertIgnoreFunction())

	app.Router.
		Methods("POST").
		Path("/api/transactionUpsert").
		HandlerFunc(transactions.UpsertFunction())

	app.Router.
		Methods("GET").
		Path("/api/categories").
		HandlerFunc(categories.GetFunction())

	app.Router.
		Methods("GET").
		Path("/api/itemTokens").
		HandlerFunc(itemTokens.GetFunction())

	app.Router.
		Methods("GET").
		Path("/api/itemTokensFetchTransactions").
		HandlerFunc(itemTokens.FetchTransactionsFunction())

	app.Router.
		Methods("POST").
		Path("/api/plaidItemTokens").
		HandlerFunc(plaid.CreateFromPublicTokenFunction())

	app.Router.
		Methods("POST").
		Path("/api/plaidGeneratePublicToken").
		HandlerFunc(plaid.GeneratePublicTokenFunction())

	app.Router.
		Methods("GET").
		Path("/api/transactions").
		HandlerFunc(transactions.GetFunction())

	app.Router.
		Methods("PUT").
		Path("/api/transactions").
		HandlerFunc(transactions.PutFunction())

	//Step one of import
	app.Router.
		Methods("POST").
		Path("/api/checkTransactions").
		HandlerFunc(transactions.CheckFunction())

	//Step two of import
	app.Router.
		Methods("POST").
		Path("/api/importTransactions").
		HandlerFunc(transactions.ImportFunction())

	app.Router.
		Methods("GET").
		Path("/api/analysisTrees").
		HandlerFunc(analysisTrees.GetFunction())

	app.Router.
		Methods("GET").
		Path("/api/saltEdgeRefreshInteractive/{id}").
		HandlerFunc(saltedge.RefreshConnectionInteractiveFunction())

	app.Router.
		Methods("GET").
		Path("/api/saltEdgeCreateInteractive").
		HandlerFunc(saltedge.CreateConnectionInteractiveFunction())

	app.Router.
		Methods("GET").
		Path("/api/resetDB").
		HandlerFunc(resetDB.ForceResetDBFunction())

	app.Router.
		Methods("GET").
		Path("/api/resetDBFull").
		HandlerFunc(resetDB.ForceResetDBFullFunction())

	app.Router.
		Methods("POST").
		Path("/api/customTree").
		HandlerFunc(analysisTrees.CustomAnalyze())

}
