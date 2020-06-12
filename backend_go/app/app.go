package app

import (
	// "database/sql"
	// "encoding/json"
	// "log"

	"github.com/gorilla/mux"
	// "github.com/jmoiron/sqlx"

	"fintrack-go/routes/accounts"
	"fintrack-go/routes/analysisTrees"
	"fintrack-go/routes/categories"
	"fintrack-go/routes/itemTokens"
	"fintrack-go/routes/plaid"
	"fintrack-go/routes/resetDB"
	"fintrack-go/routes/saltedge"
	"fintrack-go/routes/transactions"
)

type App struct {
	Router *mux.Router
	// Database *sqlx.DB
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

		//This one will go away when we do all filtering and trees server side
	// app.Router.
	// 	Methods("GET").
	// 	Path("/transactionsRange/{range}").
	// 	HandlerFunc(transactions.GetRangeFunction())

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

		//This will become a generic itemToken refresh connections function for both SE and Plaid
		//Actually just going to include in fetch transactions
	// app.Router.
	// 	Methods("GET").
	// 	Path("/saltEdgeConnections").
	// 	HandlerFunc(saltedge.RefreshConnectionsFunction())

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

	// app.Router.
	// 	Methods("GET").
	// 	Path("/api/resetToken").
	// 	HandlerFunc(plaid.ResetToken())

	// app.Router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	socket.ServeWs(socket.ExportHub, w, r)
	// })

	// app.Router.HandleFunc("/ws", itemTokens.FetchTransactionsFunction())
	// app.Router.
	// Methods("GET").
	// Path("/").
	// HandlerFunc()
}
