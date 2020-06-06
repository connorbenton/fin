package itemTokens

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"

	// "fmt"
	"fintrack-go/db"
	"fintrack-go/routes/analysisTrees"
	"fintrack-go/routes/plaid"
	"fintrack-go/routes/saltedge"
	"fintrack-go/types"

	_ "github.com/jmoiron/sqlx"
)

type wsMsg struct {
	// type message struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}

// type CurrencyRate struct {
// 	Id    int            `json:"id"`
// 	Date  time.Time      `json:"date" db:"date"`
// 	Rates types.JSONText `json:"rates" db:"rates"`

// 	CreatedAt time.Time `json:"created_at" db:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
// }

func SelectAll() []types.ItemToken {
	dbdata := []types.ItemToken{}
	err := db.DBCon.Select(&dbdata, "SELECT * FROM `item_tokens`")
	if err != nil {
		panic(err)
	}
	return dbdata
}

func GetFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		dbdata := SelectAll()
		// err := db.DBCon.Select(&dbdata, "SELECT * FROM `item_tokens`")
		// if err != nil {
		// log.Fatal(err)
		// }

		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(dbdata); err != nil {
			panic(err)
		}
	}
}

//Need a 'refresh item tokens for new accounts' method to go along with a new button on Accounts tab, instead of refreshing SaltEdge and Plaid connections on each try
//Actually don't, we'll do refresh toks/accs on every trans fetch

func FetchTransactionsFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		baseCurrency := strings.ToUpper(os.Getenv("BASE_CURRENCY"))
		SE := strings.ToUpper(os.Getenv("USE_SALTEDGE"))
		Plaid := strings.ToUpper(os.Getenv("USE_PLAID"))
		var useSE, usePlaid bool

		if SE == "TRUE" {
			useSE = true
		} else {
			useSE = false
		}

		if Plaid == "TRUE" {
			usePlaid = true
		} else {
			usePlaid = false
		}

		// First get all item tokens
		// itemTokens := []types.ItemToken{}
		// err := db.DBCon.Select(&itemTokens, "SELECT * FROM `item_tokens`")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		txnPre := db.DBCon.MustBegin()

		istmtPre := types.PrepItemSt(txnPre)
		astmtPre := types.PrepAccountSt(txnPre)
		// // istmt := types.PrepItemSt(db.DBCon)
		// // astmt := types.PrepAccountSt(db.DBCon)
		// tstmt := types.PrepTransSt(txn)

		itemTokens := SelectAll()

		var wgPre sync.WaitGroup
		wgPre.Add(1)
		go func() {
			defer wgPre.Done()
			// saltedge.RefreshConnectionsFunction(txn)
			// saltedge.RefreshConnectionsFunction(istmt, astmt)
			if useSE {
				saltedge.RefreshConnectionsFunction(istmtPre, astmtPre)
			}
		}()
		wgPre.Add(1)
		go func() {
			defer wgPre.Done()
			db.GetNewXML()
		}()
		for _, itemTok := range itemTokens {
			wgPre.Add(1)
			go func(itemToken types.ItemToken) {
				defer wgPre.Done()
				if itemToken.Provider == "Plaid" && usePlaid {
					plaid.RefreshConnection(itemToken, istmtPre, astmtPre)
				}
			}(itemTok)
		}
		wgPre.Wait()

		err := txnPre.Commit()
		if err != nil {
			panic(err)
		}

		// Make sure currencies are up to date

		// Refresh Plaid and SaltEdge connections

		// Then we iterate through item tokens and process in either saltedge or plaid

		txn := db.DBCon.MustBegin()
		// istmt := types.PrepItemSt(txn)
		astmt := types.PrepAccountSt(txn)
		tstmt := types.PrepTransSt(txn)
		istmtOnlyTx := types.PrepItemStOnlyTx(txn)

		var wg sync.WaitGroup
		for _, itemTok := range itemTokens {
			wg.Add(1)
			go func(itemToken types.ItemToken) {
				defer wg.Done()
				// Think these can be done in goroutines
				// Using websocket connection to transmit which item is being currently worked on
				// message := []byte(`{ "username": "Booh", }`)
				// socket.ExportHub.Broadcast <- message

				if itemToken.Provider == "SaltEdge" && useSE {
					// saltedge.FetchTransactionsForItemToken(itemToken, txn, baseCurrency)
					// saltedge.FetchTransactionsForItemToken(itemToken, istmt, astmt, tstmt, baseCurrency)
					saltedge.FetchTransactionsForItemToken(itemToken, istmtOnlyTx, astmt, tstmt, baseCurrency)
					// if itemToken.Interactive {
					// 	// Needs to be direct DB call
					// 	itemToken.LastDownloadedTransactions = itemToken.LastRefresh
					// } else {
					// 	// Needs to be direct DB call
					// 	itemToken.LastDownloadedTransactions = time.Now()
					// }
				} else if usePlaid {
					// plaid.FetchTransactionsForItemToken(itemToken, txn, baseCurrency)
					// plaid.FetchTransactionsForItemToken(itemToken, istmt, astmt, tstmt, baseCurrency)
					plaid.FetchTransactionsForItemToken(itemToken, istmtOnlyTx, astmt, tstmt, baseCurrency)
					// Needs direct DB call here to set LastDownloadedTransactions
				}
			}(itemTok)

		}
		wg.Wait()

		err2 := txn.Commit()
		if err2 != nil {
			panic(err2)
		}

		analysisTrees.ReAnalyze()

		res.WriteHeader(http.StatusOK)
		// currencyRates := []CurrencyRate{}
		// err2 := db.DBCon.Select(&currencyRates, "SELECT * FROM `currency_rates`")
		// if err2 != nil {
		// 	log.Fatal(err)
		// }
		// rates, _ := string([]byte{json.Marshal(currencyRates[0].DataJSON)})
		// log.Println(rates)

		// res.WriteHeader(http.StatusOK)
		// if err := json.NewEncoder(res).Encode(itemTokens); err != nil {
		// 	panic(err)
		// }
	}
}
