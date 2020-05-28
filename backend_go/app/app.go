package app

import (
	// "database/sql"
	// "encoding/json"
	// "log"

	"net/http"

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
	"fintrack-go/socket"
)

type App struct {
	Router *mux.Router
	// Database *sqlx.DB
}

func (app *App) SetupRouter() {
	app.Router.
		Methods("GET").
		Path("/accounts").
		HandlerFunc(accounts.GetFunction())

	app.Router.
		Methods("GET").
		Path("/categories").
		HandlerFunc(categories.GetFunction())

	app.Router.
		Methods("GET").
		Path("/itemTokens").
		HandlerFunc(itemTokens.GetFunction())

	app.Router.
		Methods("GET").
		Path("/itemTokensFetchTransactions").
		HandlerFunc(itemTokens.FetchTransactionsFunction())

	app.Router.
		Methods("POST").
		Path("/plaidItemTokens").
		HandlerFunc(plaid.CreateFromPublicTokenFunction())

	app.Router.
		Methods("POST").
		Path("/plaidGeneratePublicToken").
		HandlerFunc(plaid.GeneratePublicTokenFunction())

	app.Router.
		Methods("GET").
		Path("/transactions").
		HandlerFunc(transactions.GetFunction())

		//This one will go away when we do all filtering and trees server side
	// app.Router.
	// 	Methods("GET").
	// 	Path("/transactionsRange/{range}").
	// 	HandlerFunc(transactions.GetRangeFunction())

	app.Router.
		Methods("PUT").
		Path("/transactions").
		HandlerFunc(transactions.PutFunction())

	app.Router.
		Methods("POST").
		Path("/importTransactions").
		HandlerFunc(transactions.ImportFunction())

	app.Router.
		Methods("GET").
		Path("/analysisTrees").
		HandlerFunc(analysisTrees.GetFunction())

		//This will become a generic itemToken refresh connections function for both SE and Plaid
		//Actually just going to include in fetch transactions
	// app.Router.
	// 	Methods("GET").
	// 	Path("/saltEdgeConnections").
	// 	HandlerFunc(saltedge.RefreshConnectionsFunction())

	app.Router.
		Methods("GET").
		Path("/saltEdgeRefreshInteractive/{id}").
		HandlerFunc(saltedge.RefreshConnectionInteractiveFunction())

	app.Router.
		Methods("GET").
		Path("/saltEdgeCreateInteractive").
		HandlerFunc(saltedge.CreateConnectionInteractiveFunction())

	app.Router.
		Methods("GET").
		Path("/resetDB").
		HandlerFunc(resetDB.ForceResetDBFunction())

	app.Router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWs(socket.ExportHub, w, r)
	})

	// app.Router.HandleFunc("/ws", itemTokens.FetchTransactionsFunction())
	// app.Router.
	// Methods("GET").
	// Path("/").
	// HandlerFunc()
}
