package saltedge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	// "fmt"
	"fintrack-go/db"
	"fintrack-go/types"

	_ "github.com/jmoiron/sqlx"
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
	// appID := os.Getenv("SALTEDGE_APP_ID")
	// fmt.Println(appID)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("App-id", os.Getenv("SALTEDGE_APP_ID"))
	req.Header.Set("Secret", os.Getenv("SALTEDGE_APP_SECRET"))
	req.Header.Set("Expires-at", strconv.FormatInt((time.Now().Unix()+60), 10))
	// fmt.Println("request Headers:", req.Header)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println("response Body:", string(body))
	return string(body)

}

func RefreshConnectionsFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		var wgConnections sync.WaitGroup
		url := "https://www.saltedge.com/api/v5/connections?customer_id=" + os.Getenv("SALTEDGE_CUSTOMER_ID")

		// ch := make(chan string)
		// var responses []string
		// var user string
		// var wg sync.WaitGroup
		// wg.Add(1)
		// go saltEdgeReq("GET", url, "", ch, &wg)

		// go func() {
		// wg.Wait()
		// close(ch)
		// }()

		// for res := range ch {
		// responses = append(responses, res)
		// }

		// responses[0]
		connections := saltEdgeReq("GET", url, "")
		// log.Println(connections)

		// log.Println(responses[0])
		var data types.ConnectionResponse
		// json.Unmarshal([]byte(responses[0]), &data)
		json.Unmarshal([]byte(connections), &data)
		// fmt.Printf("Results: %v\n", data)

		txn := db.DBCon.MustBegin()

		// chAccounts := make(chan string)
		// var responsesAccounts []string
		// wgAccounts.Add(11)

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

				query := `INSERT INTO item_tokens(institution, provider, interactive, last_refresh, next_refresh_possible, item_id)
							VALUES(:institution, :provider, :interactive, :last_refresh, :next_refresh_possible, :item_id) 
							ON CONFLICT (item_id, provider) DO UPDATE SET
							interactive = excluded.interactive,
							last_refresh = excluded.last_refresh,
							next_refresh_possible = excluded.next_refresh_possible`
				_, err := txn.NamedExec(query, item)
				if err != nil {
					panic(err)
				}
				// wgAccounts.Add(1)
				url2 := "https://www.saltedge.com/api/v5/accounts?connection_id=" + conn.ID
				// go saltEdgeReq("GET", url, "", chAccounts, &wgAccounts)
				accounts := saltEdgeReq("GET", url2, "")
				// log.Println(accounts)

				// res.WriteHeader(http.StatusOK)
				// res.Write([]byte(accounts))
				// if err := json.NewEncoder(res).Encode(dbdata); err != nil {
				// if err := json.NewEncoder(res).Encode(accounts); err != nil {
				// panic(err)
				// }
				// log.Println(accounts)
				var data types.AccountResponse
				json.Unmarshal([]byte(accounts), &data)
				// var wgAccounts sync.WaitGroup
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
						query := `INSERT INTO accounts(name, institution, provider, account_id, item_id, type, 'limit', available, balance, currency)
							VALUES(:name, :institution, :provider, :account_id, :item_id, :type, :limit, :available, :balance, :currency) 
							ON CONFLICT (account_id, provider) DO UPDATE SET
							'limit' = excluded.'limit',
							available = excluded.available,
							balance = excluded.balance`
						_, err := txn.NamedExec(query, acc)
						if err != nil {
							panic(err)
						}
						// fmt.Printf("%v\n", acc.AccountID)
						// fmt.Printf("%v\n", account.Extra.AccountName)
						// fmt.Printf("%v\n", acc.Name)
					}(account)
				}
				wgAccounts.Wait()
			}(connection)

		}

		// wgAccounts.Wait()
		wgConnections.Wait()
		err := txn.Commit()
		if err != nil {
			panic(err)
		}

		fmt.Println("finished sync group")

		// close the channel in the background
		// go func() {
		// wgAccounts.Wait()
		// close(chAccounts)
		// }()
		// read from channel as they come in until its closed
		// for res := range chAccounts {
		// responsesAccounts = append(responsesAccounts, res)
		// }
		// var params = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
		// req, err := http.NewRequest("GET", url)

		// req := saltEdgeGet(url)

		// client := &http.Client{}
		// resp, err := client.Do(req)
		// if err != nil {
		// panic(err)
		// }
		// defer resp.Body.Close()

		// dbdata := []types.Account{}
		// // err := app.Database.Query("SELECT * FROM `categories`", id).Scan(&dbdata.id, &dbdata.topCategory, &dbdata.subCategory)
		// err := db.DBCon.Select(&dbdata, "SELECT * FROM `accounts`")
		// if err != nil {
		// 	log.Println("Database SELECT failed")
		// 	panic(err)
		// 	// fmt.Println("Database SELECT failed")
		// 	// fmt.Println(err)
		// 	// return
		// }

		// log.Println("You fetched a thing!")
		// res.WriteHeader(http.StatusOK)
		// // if err := json.NewEncoder(res).Encode(dbdata); err != nil {
		// if err := json.NewEncoder(res).Encode(data); err != nil {
		// 	panic(err)
		// }
	}
}

func RefreshConnectionInteractiveFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		dbdata := []types.Account{}
		// err := app.Database.Query("SELECT * FROM `categories`", id).Scan(&dbdata.id, &dbdata.topCategory, &dbdata.subCategory)
		err := db.DBCon.Select(&dbdata, "SELECT * FROM `categories`")
		if err != nil {
			log.Fatal("Database SELECT failed")
			// fmt.Println("Database SELECT failed")
			// fmt.Println(err)
			// return
		}

		log.Println("You fetched a thing!")
		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(dbdata); err != nil {
			panic(err)
		}
	}
}

func CreateConnectionInteractiveFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		dbdata := []types.Account{}
		// err := app.Database.Query("SELECT * FROM `categories`", id).Scan(&dbdata.id, &dbdata.topCategory, &dbdata.subCategory)
		err := db.DBCon.Select(&dbdata, "SELECT * FROM `categories`")
		if err != nil {
			log.Fatal("Database SELECT failed")
			// fmt.Println("Database SELECT failed")
			// fmt.Println(err)
			// return
		}

		log.Println("You fetched a thing!")
		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(dbdata); err != nil {
			panic(err)
		}
	}
}

func FetchTransactionsForItemToken(ItemID string) {

	url := "https://www.saltedge.com/api/v5/transactions?connection_id" + ItemID
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
