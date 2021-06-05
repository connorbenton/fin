package itemTokens

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"

	"fin-go/db"
	"fin-go/routes/analysisTrees"
	"fin-go/routes/plaid"
	"fin-go/routes/saltedge"
	"fin-go/types"

	_ "github.com/jmoiron/sqlx"
)

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

		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(dbdata); err != nil {
			panic(err)
		}
	}
}

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

		txnPre := db.DBCon.MustBegin()

		istmtPre := types.PrepItemSt(txnPre)
		astmtPre := types.PrepAccountSt(txnPre)

		itemTokens := SelectAll()

		var wgPre sync.WaitGroup
		wgPre.Add(1)
		go func() {
			defer wgPre.Done()
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

		// Refresh connections

		// Then we iterate through item tokens and process in either saltedge or plaid

		txn := db.DBCon.MustBegin()
		astmt := types.PrepAccountSt(txn)
		tstmt := types.PrepTransSt(txn)
		istmtOnlyTx := types.PrepItemStOnlyTx(txn)

		var wg sync.WaitGroup
		for _, itemTok := range itemTokens {
			wg.Add(1)
			go func(itemToken types.ItemToken) {
				defer wg.Done()

				if itemToken.Provider == "SaltEdge" && useSE {
					saltedge.FetchTransactionsForItemToken(itemToken, istmtOnlyTx, astmt, tstmt, baseCurrency)
				} else if usePlaid && itemToken.Provider == "Plaid" {
					plaid.FetchTransactionsForItemToken(itemToken, istmtOnlyTx, astmt, tstmt, baseCurrency)
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

	}
}
