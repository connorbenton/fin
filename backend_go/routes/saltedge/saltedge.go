package saltedge

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"sync"
	"time"

	"fintrack-go/db"
	"fintrack-go/types"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func saltEdgeReq(verb string, url string, params string) string {
	var err error
	var req *http.Request
	if verb == "GET" {
		req, err = http.NewRequest("GET", url, nil)
	} else if verb == "POST" {
		jsonStr := []byte(params)
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	}
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("App-id", os.Getenv("SALTEDGE_APP_ID"))
	req.Header.Set("Secret", os.Getenv("SALTEDGE_APP_SECRET"))
	req.Header.Set("Expires-at", strconv.FormatInt((time.Now().Unix()+60), 10))
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)

}

func RefreshConnectionsFunction(istmt, astmt *sqlx.NamedStmt) {
	var wgConnections sync.WaitGroup
	url := "https://www.saltedge.com/api/v5/connections?customer_id=" + os.Getenv("SALTEDGE_CUSTOMER_ID")

	connections := saltEdgeReq("GET", url, "")
	var data types.ConnectionResponse
	json.Unmarshal([]byte(connections), &data)
	for _, connection := range data.Data {
		wgConnections.Add(1)
		go func(conn types.SEConnection) {
			defer wgConnections.Done()
			item := types.ItemToken{}
			item.Institution = conn.ProviderName
			item.Provider = "SaltEdge"
			if conn.LastAttempt.Interactive {
				item.Interactive = true
			} else {
				item.Interactive = false
			}
			item.LastRefresh = conn.LastSuccessAt
			item.NextRefreshPossible = conn.NextRefreshPossibleAt
			item.ItemID = conn.ID

			istmt.MustExec(item)
			url2 := "https://www.saltedge.com/api/v5/accounts?connection_id=" + conn.ID
			accounts := saltEdgeReq("GET", url2, "")

			var data types.AccountResponse
			json.Unmarshal([]byte(accounts), &data)
			var wgAccounts sync.WaitGroup
			for _, account := range data.Data {
				wgAccounts.Add(1)
				go func(SEAcc types.SEAccount) {
					defer wgAccounts.Done()
					acc := types.Account{}
					if SEAcc.Extra.AccountName == "" {
						acc.Name = SEAcc.Name
					} else {
						acc.Name = SEAcc.Extra.AccountName
					}
					acc.Institution = item.Institution
					acc.Provider = "SaltEdge"
					acc.AccountID = SEAcc.ID
					acc.ItemID = SEAcc.ConnectionID
					acc.Type = SEAcc.Nature
					acc.Limit = SEAcc.Extra.CreditLimit
					acc.Available = SEAcc.Extra.AvailableAmount
					acc.Balance = SEAcc.Balance
					acc.Currency = SEAcc.CurrencyCode

					astmt.MustExec(acc)
				}(account)
			}
			wgAccounts.Wait()
		}(connection)

	}

	wgConnections.Wait()

}

func RefreshConnectionInteractiveFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		var err error

		vars := mux.Vars(req)
		connID := vars["id"]

		url := "https://www.saltedge.com/api/v5/connect_sessions/refresh"
		params := fmt.Sprintf(`{
			"data": {
				"connection_id": %q,
				"attempt": {
				"return_to": %q
				}
			}
		}`, connID, os.Getenv("BASE_URL"))
		refresh := saltEdgeReq("POST", url, params)

		data := types.CreateRefreshResponse{}

		err = json.Unmarshal([]byte(refresh), &data)
		if err != nil {
			fmt.Printf("Error with Create Connection Interactive: %v \n", err)
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte(data.Data.ConnectURL))
		if err != nil {
			fmt.Printf("Error with Create Connection Interactive: %v \n", err)
		}

	}
}

func CreateConnectionInteractiveFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		url := "https://www.saltedge.com/api/v5/connect_sessions/create"
		params := fmt.Sprintf(`{
			"data": {
				"customer_id": %q,
				"consent": {
				"scopes": [
					"account_details",
					"transactions_details"
				]
				},
				"attempt": {
				"return_to": %q
				}
			}
		}`, os.Getenv("SALTEDGE_CUSTOMER_ID"), os.Getenv("BASE_URL"))
		create := saltEdgeReq("POST", url, params)

		data := types.CreateRefreshResponse{}

		err := json.Unmarshal([]byte(create), &data)
		if err != nil {
			errString := fmt.Sprintf("Error with Create Connection Interactive: %v \n", err)
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		}

		if data.Data.ConnectURL != "" {
			_, err2 := res.Write([]byte(data.Data.ConnectURL))
			if err2 != nil {
				errString := fmt.Sprintf("Error with Create Connection Interactive: %v \n", err)
				log.Println(errString)
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(errString))
			}
		} else {
			errString := fmt.Sprintf("ConnectURL field empty")
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		}
	}
}

func FetchTransactionsForItemToken(iTok types.ItemToken, istmt *sqlx.NamedStmt, astmt *sqlx.NamedStmt, tstmt *sqlx.NamedStmt, baseCurrency string) {

	url := "https://www.saltedge.com/api/v5/transactions?connection_id=" + iTok.ItemID

	res := saltEdgeReq("GET", url, "")

	var data types.TransactionsResponse
	err := json.Unmarshal([]byte(res), &data)
	if err != nil {
		panic(err)
	}
	for _, tx := range data.Data {
		var err error
		trans := types.Transaction{}
		if tx.Extra.PostingDate == "" {
			trans.Date = tx.MadeOn
		} else {
			trans.Date = tx.Extra.PostingDate
		}
		trans.Description = tx.Description
		trans.Amount = tx.Amount
		trans.AccountID = tx.AccountID

		var name string
		querytest := "SELECT name FROM accounts WHERE account_id='" + tx.AccountID + "' AND provider='SaltEdge'"
		// log.Println(querytest)
		err = db.DBCon.Get(&name, querytest)
		if err != nil {
			panic(err)
		}
		trans.AccountName = name

		trans.TransactionID = tx.ID

		trans.CurrencyCode = tx.CurrencyCode
		// log.Println(query2)
		trans.NormalizedAmount = db.GetNormalizedAmount(trans.CurrencyCode, baseCurrency, trans.Date, trans.Amount)

		//Searching for bottom category match first
		sCat := types.CategorySE{}
		query := fmt.Sprintf(`SELECT * FROM salt_edge__categories WHERE bottom_category = %q AND top_category = 'personal'`, tx.Category)
		err = db.DBCon.Get(&sCat, query)
		if err != nil && err != sql.ErrNoRows {
			panic(err)
		}
		if (types.CategorySE{}) == sCat {
			//If nil for bottom category then look for match in sub category
			query := fmt.Sprintf(`SELECT * FROM salt_edge__categories WHERE sub_category = %q AND top_category = 'personal'`, tx.Category)
			err := db.DBCon.Get(&sCat, query)
			if err != nil && err != sql.ErrNoRows {
				panic(err)
			}
			if (types.CategorySE{}) == sCat {
				//If still nil then set category to Uncategorized
				trans.Category = 106
				trans.CategoryName = "Uncategorized"
			} else {
				trans.Category = sCat.LinkToAppCat
				trans.CategoryName = sCat.AppCatName
			}
		} else {
			trans.Category = sCat.LinkToAppCat
			trans.CategoryName = sCat.AppCatName
		}
		tstmt.MustExec(trans)
	}


	if iTok.Interactive {
		iTok.LastDownloadedTransactions = iTok.LastRefresh
	} else {
		iTok.LastDownloadedTransactions = time.Now()
	}

	istmt.MustExec(iTok)

}
