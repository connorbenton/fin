package app

import (
	// "database/sql"
	// "encoding/json"
	// "log"
	// "net/http"

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
	Router   *mux.Router
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

	app.Router.
		Methods("GET").
		Path("/saltEdgeConnections").
		HandlerFunc(saltedge.GetConnectionsFunction())

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

	// app.Router.
		// Methods("GET").
		// Path("/").
		// HandlerFunc()
}