package plaid

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	// "fmt"
	"fintrack-go/db"
	"fintrack-go/types"

	"github.com/jmoiron/sqlx"
	"github.com/plaid/plaid-go/plaid"
	"github.com/shopspring/decimal"
)

func newClient() (*plaid.Client, error) {
	env := os.Getenv("PLAID_ENVIRONMENT")
	if env == "sandbox" {
		clientOptions := plaid.ClientOptions{
			os.Getenv("PLAID_CLIENT_ID"),
			os.Getenv("PLAID_SECRET_SANDBOX"),
			os.Getenv("PLAID_PUBLIC_KEY"),
			plaid.Sandbox,
			&http.Client{},
		}
		client, err := plaid.NewClient(clientOptions)
		if err != nil {
			panic(err)
		}
		return client, nil
	} else if env == "development" {
		clientOptions := plaid.ClientOptions{
			os.Getenv("PLAID_CLIENT_ID"),
			os.Getenv("PLAID_SECRET_DEVELOPMENT"),
			os.Getenv("PLAID_PUBLIC_KEY"),
			plaid.Development,
			&http.Client{},
		}
		client, err := plaid.NewClient(clientOptions)
		if err != nil {
			panic(err)
		}
		return client, nil
	} else {
		return nil, errors.New("Environment variable is not either 'sandbox' or 'development'")
	}
}

func CreateFromPublicTokenFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		var err error

		decoder := json.NewDecoder(req.Body)
		var item types.CreateTokenPost
		err = decoder.Decode(&item)
		if err != nil {
			// panic(err)
			errString := fmt.Sprintf("Error with Create Token Post Request: %v \n", err)
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		}

		pClient, err := newClient()
		if err != nil {
			// panic(err)
			errString := fmt.Sprintf("Error with Create Token Client Create: %v \n", err)
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		}

		pRes, err := pClient.ExchangePublicToken(item.Token)
		if err != nil {
			// panic(err)
			errString := fmt.Sprintf("Error with Create Token Client Exchange Token request: %v \n", err)
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		}

		txn := db.DBCon.MustBegin()

		iTok := types.ItemToken{}
		iTok.ItemID = pRes.ItemID
		iTok.AccessToken = pRes.AccessToken
		iTok.Institution = item.Name
		iTok.Provider = "Plaid"
		iTok.NeedsReLogin = false
		// upsertItemToken(iTok, txn)
		istmt := types.PrepItemSt(txn)
		// istmt := types.PrepItemSt(db.DBCon)
		upsertItemToken(iTok, istmt)
		astmt := types.PrepAccountSt(txn)
		// astmt := types.PrepAccountSt(db.DBCon)

		pAccountRes, err := pClient.GetAccounts(pRes.AccessToken)
		if err != nil {
			// panic(err)
			errString := fmt.Sprintf("Error with Create Token Client Get Accounts request: %v \n", err)
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		}

		var wg sync.WaitGroup
		for _, account := range pAccountRes.Accounts {
			wg.Add(1)
			go func(pAcc plaid.Account) {
				defer wg.Done()
				acc := types.Account{}
				// upsertAccountWithPlaidAccount(acc, pAcc, iTok.Institution, iTok.ItemID, txn)
				upsertAccountWithPlaidAccount(acc, pAcc, iTok.Institution, iTok.ItemID, astmt)
			}(account)
		}
		wg.Wait()

		errtx := txn.Commit()
		if errtx != nil {
			errString := fmt.Sprintf("Error with Create Token Txn Commit: %v \n", errtx)
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		}

		_, err2 := res.Write([]byte("Upserted " + iTok.Institution))
		if err2 != nil {
			errString := fmt.Sprintf("Error with Create Token Response Write: %v \n", err)
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		}
	}
}

// func upsertItemTokenPrep(txn *sqlx.Tx) *sqlx.NamedStmt {

// 	query := `INSERT INTO item_tokens(institution, provider, needs_re_login, access_token, item_id, last_downloaded_transactions)
// 				VALUES(:institution, :provider, :needs_re_login, :access_token, :item_id, :last_downloaded_transactions)
// 				ON CONFLICT (item_id, provider) DO UPDATE SET
// 				needs_re_login = excluded.needs_re_login
// 				last_downloaded_transactions = excluded.last_downloaded_transactions`
// 	nstmt, err := txn.PrepareNamed(query)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// nstmtx := txn.PrepareNamed(nstmt)
// 	return nstmt
// 	// _, err := txn.NamedExec(query, tok)

// }

func upsertItemToken(tok types.ItemToken, stmt *sqlx.NamedStmt) {

	stmt.MustExec(tok)
	// _, err := txn.NamedExec(query, tok)
	// if err != nil {
	// 	panic(err)
	// }

}

// func upsertAccountPrep(txn *sqlx.Tx) *sqlx.NamedStmt {
// 	query := `INSERT INTO accounts(name, institution, provider, account_id, item_id, type, 'limit', available, balance, currency, subtype)
// 							VALUES(:name, :institution, :provider, :account_id, :item_id, :type, :limit, :available, :balance, :currency, :subtype)
// 							ON CONFLICT (account_id, provider) DO UPDATE SET
// 							'limit' = excluded.'limit',
// 							available = excluded.available,
// 							balance = excluded.balance`
// 	nstmt, err := txn.PrepareNamed(query)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// nstmtx := txn.Stmtx(nstmt)
// 	return nstmt
// }

func upsertAccountWithPlaidAccount(acc types.Account, pAcc plaid.Account, inst string, itemID string, stmt *sqlx.NamedStmt) {
	acc.Name = pAcc.Name
	acc.Institution = inst
	acc.AccountID = pAcc.AccountID
	acc.Provider = "Plaid"
	if pAcc.Type == "credit" {
		acc.Balance = decimal.NewFromFloat(pAcc.Balances.Current * -1)
	} else {
		acc.Balance = decimal.NewFromFloat(pAcc.Balances.Current)
	}
	acc.Limit = decimal.NewFromFloat(pAcc.Balances.Limit)
	acc.Available = decimal.NewFromFloat(pAcc.Balances.Available)
	acc.Currency = pAcc.Balances.ISOCurrencyCode
	acc.Type = pAcc.Type
	acc.Subtype = pAcc.Subtype
	acc.ItemID = itemID

	stmt.MustExec(acc)
	// _, err := db.DBCon.NamedExec(query, acc)
	// if err != nil {
	// 	panic(err)
	// }

}

func GeneratePublicTokenFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		var err error

		decoder := json.NewDecoder(req.Body)
		var item types.GenerateTokenPost
		err = decoder.Decode(&item)
		if err != nil {
			// panic(err)
			errString := fmt.Sprintf("Error with Generate Token Post Request: %v \n", err)
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		}

		var access string

		query := fmt.Sprintf(`SELECT access_token FROM item_tokens WHERE item_id = %q`, item.ItemID)
		err = db.DBCon.Get(&access, query)
		if err != nil {
			// panic(err)
			errString := fmt.Sprintf("Error with Generate Token Item Query: %v \n", err)
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		}

		pClient, err := newClient()
		if err != nil {
			// panic(err)
			errString := fmt.Sprintf("Error with Generate Token Client Create: %v \n", err)
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		}

		pRes, err := pClient.CreatePublicToken(access)
		if err != nil {
			// panic(err)
			errString := fmt.Sprintf("Error with Generate Token Client Create Token request: %v \n", err)
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		}

		// data := types.CreateRefreshResponse{}

		if pRes.PublicToken == "" {
			errString := fmt.Sprintf("PublicToken response seems to be empty")
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		} else {
			resJSON := `{"public_token": "` + pRes.PublicToken + `"}`
			_, err2 := res.Write([]byte(resJSON))
			// _, err2 := res.Write([]byte(pRes.PublicToken))
			if err2 != nil {
				errString := fmt.Sprintf("Error with Plaid Generate Public Token: %v \n", err2)
				log.Println(errString)
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(errString))
			}
		}

	}
}

// Don't think refreshconnections is needed for Plaid since we'll refresh on every transaction fetch

// func RefreshConnectionsFunction() func(http.ResponseWriter, *http.Request) {
// 	return func(res http.ResponseWriter, req *http.Request) {

// 		dbdata := []types.Account{}
// 		// err := app.Database.Query("SELECT * FROM `categories`", id).Scan(&dbdata.id, &dbdata.topCategory, &dbdata.subCategory)
// 		err := db.DBCon.Select(&dbdata, "SELECT * FROM `categories`")
// 		if err != nil {
// 			log.Fatal("Database SELECT failed")
// 			// fmt.Println("Database SELECT failed")
// 			// fmt.Println(err)
// 			// return
// 		}

// 		log.Println("You fetched a thing!")
// 		res.WriteHeader(http.StatusOK)
// 		if err := json.NewEncoder(res).Encode(dbdata); err != nil {
// 			panic(err)
// 		}
// 	}
// }
// func ResetToken() func(http.ResponseWriter, *http.Request) {
// 	return func(res http.ResponseWriter, req *http.Request) {
// 		pClient, err := newClient()
// 		if err != nil {
// 			panic(err)
// 		}
// 		pClient.ResetSandboxItem("access-sandbox-84adcd32-fe27-46ac-928c-1ae0b857f3c3")
// 		res.WriteHeader(http.StatusOK)
// 	}
// }

func RefreshConnection(iTok types.ItemToken, istmt, astmt *sqlx.NamedStmt) {

	pClient, err := newClient()
	if err != nil {
		panic(err)
	}

	// istmt := upsertItemTokenPrep(txn)
	// astmt := upsertAccountPrep(txn)

	// queryIns := `INSERT INTO transactions('date', transaction_id, description, amount, normalized_amount, category,
	// 			category_name, account_name, currency_code, account_id)
	// 			VALUES(:date, :transaction_id, :description, :amount, :normalized_amount, :category,
	// 			:category_name, :account_name, :currency_code, :account_id)
	// 			ON CONFLICT (transaction_id) DO UPDATE SET
	// 			'date' = excluded.'date',
	// 			description = excluded.description,
	// 			amount = excluded.amount
	// 			normalized_amount = excluded.normalized_amount
	// 			category = excluded.category
	// 			category_name = excluded.category_name`

	// tstmt, err := txn.PrepareNamed(queryIns)
	// if err != nil {
	// 	panic(err)
	// }

	// txnPre := db.DBCon.MustBegin()

	// astmtPre := types.PrepAccountSt(txnPre)

	// var wgWrap sync.WaitGroup
	// wgWrap.Add(2)
	// go func() {
	// 	defer wgWrap.Done()
	pAccountRes, err := pClient.GetAccounts(iTok.AccessToken)
	if err != nil {
		perr, ok := err.(plaid.Error)
		if ok {
			log.Println(fmt.Sprintf("Plaid error: %v", perr))
			if perr.ErrorCode == "ITEM_LOGIN_REQUIRED" {
				iTok.NeedsReLogin = true
				upsertItemToken(iTok, istmt)
				return
			}
		}
		errString := fmt.Sprintf("Error with Create Token Client Get Accounts request: %v \n", err)
		log.Println(errString)
		panic(err)
	}

	var wgAcc sync.WaitGroup
	for _, account := range pAccountRes.Accounts {
		wgAcc.Add(1)
		go func(pAcc plaid.Account) {
			defer wgAcc.Done()
			acc := types.Account{}
			// upsertAccountWithPlaidAccount(acc, pAcc, iTok.Institution, iTok.ItemID, txn)
			upsertAccountWithPlaidAccount(acc, pAcc, iTok.Institution, iTok.ItemID, astmt)
		}(account)
	}
	wgAcc.Wait()

	// err2 := txnPre.Commit()
	// if err2 != nil {
	// 	panic(err2)
	// }
}

func FetchTransactionsForItemToken(iTok types.ItemToken, istmt *sqlx.NamedStmt, astmt *sqlx.NamedStmt, tstmt *sqlx.NamedStmt, baseCurrency string) {
	// func FetchTransactionsForItemToken(iTok types.ItemToken, baseCurrency string) {
	today := time.Now().Format("2006-01-02")

	pClient, err := newClient()
	if err != nil {
		panic(err)
	}


	// istmt := upsertItemTokenPrep(txn)
	// astmt := upsertAccountPrep(txn)

	// queryIns := `INSERT INTO transactions('date', transaction_id, description, amount, normalized_amount, category,
	// 			category_name, account_name, currency_code, account_id)
	// 			VALUES(:date, :transaction_id, :description, :amount, :normalized_amount, :category,
	// 			:category_name, :account_name, :currency_code, :account_id)
	// 			ON CONFLICT (transaction_id) DO UPDATE SET
	// 			'date' = excluded.'date',
	// 			description = excluded.description,
	// 			amount = excluded.amount
	// 			normalized_amount = excluded.normalized_amount
	// 			category = excluded.category
	// 			category_name = excluded.category_name`

	// tstmt, err := txn.PrepareNamed(queryIns)
	// if err != nil {
	// 	panic(err)
	// }

	// txnPre := db.DBCon.MustBegin()

	// astmtPre := types.PrepAccountSt(txnPre)

	// // var wgWrap sync.WaitGroup
	// // wgWrap.Add(2)
	// // go func() {
	// // 	defer wgWrap.Done()
	// pAccountRes, err := pClient.GetAccounts(iTok.AccessToken)
	// if err != nil {
	// 	//also need login required error handling here
	// 	errString := fmt.Sprintf("Error with Create Token Client Get Accounts request: %v \n", err)
	// 	log.Println(errString)
	// 	panic(err)
	// }

	// var wgAcc sync.WaitGroup
	// for _, account := range pAccountRes.Accounts {
	// 	wgAcc.Add(1)
	// 	go func(pAcc plaid.Account) {
	// 		defer wgAcc.Done()
	// 		acc := types.Account{}
	// 		// upsertAccountWithPlaidAccount(acc, pAcc, iTok.Institution, iTok.ItemID, txn)
	// 		upsertAccountWithPlaidAccount(acc, pAcc, iTok.Institution, iTok.ItemID, astmtPre)
	// 	}(account)
	// }
	// wgAcc.Wait()

	// err2 := txnPre.Commit()
	// if err2 != nil {
	// 	panic(err2)
	// }

	// }()
	// go func() {
	// defer wgWrap.Done()
	var pTransRes plaid.GetTransactionsResponse
	// pErr := plaid.Error{}
	if iTok.LastDownloadedTransactions.IsZero() {
		pTransRes, err = pClient.GetTransactions(iTok.AccessToken, "2000-01-01", today)
	} else {
		pTransRes, err = pClient.GetTransactions(iTok.AccessToken, iTok.LastDownloadedTransactions.AddDate(0, 0, -40).Format("2006-01-02"), today)
	}
	if err != nil {
		perr, ok := err.(plaid.Error)
		if ok {
			if perr.ErrorCode == "ITEM_LOGIN_REQUIRED" {
				return
			}
		}
		// p := err.Error()
		// pErr := plaid.Error(err)
		// if p.ErrorCode == "ITEM_LOGIN_REQUIRED" {
		// iTok.NeedsReLogin = true
		// upsertItemToken(iTok, istmt)
		// return
		// }
		// Figure out how to handle login required error
		// if err.Error.ErrorCode == "ITEM_LOGIN_REQUIRED" {
		// iTok.NeedsReLogin = true
		// upsertItemToken(iTok, txn)
		// } else {
		panic(err)
		// }
	}

	// var wgTrans sync.WaitGroup
	// for _, transaction := range pTransRes.Transactions {
	for _, ptx := range pTransRes.Transactions {
		// wgTrans.Add(1)
		// go func(ptx plaid.Transaction) {
		// defer wgTrans.Done()
		tx := types.Transaction{}

		// tx.Date, err = date.Parse("2006-01-02", ptx.Date)
		tx.Date = ptx.Date
		tx.TransactionID = ptx.ID
		tx.Description = ptx.Name
		tx.Amount = decimal.NewFromFloat(ptx.Amount * -1)
		tx.CurrencyCode = ptx.ISOCurrencyCode
		tx.NormalizedAmount = db.GetNormalizedAmount(tx.CurrencyCode, baseCurrency, tx.Date, tx.Amount)

		//Searching for category ID match first
		pCat := types.CategoryPlaid{}
		query := fmt.Sprintf(`SELECT * FROM plaid__categories WHERE cat_i_d = %q`, ptx.CategoryID)
		err = db.DBCon.Get(&pCat, query)
		if err != nil && err != sql.ErrNoRows {
			panic(err)
		}
		// if (types.CategoryPlaid{}) == pCat {
		if err == sql.ErrNoRows {
			//If still nil then set category to Uncategorized
			tx.Category = 106
			tx.CategoryName = "Uncategorized"
		} else {
			tx.Category = pCat.LinkToAppCat
			tx.CategoryName = pCat.AppCatName
		}

		var name string
		err = db.DBCon.Get(&name, "SELECT name FROM accounts WHERE account_id='"+ptx.AccountID+"' AND provider='Plaid' LIMIT 1")
		if err != nil {
			panic(err)
		}
		tx.AccountName = name
		tx.AccountID = ptx.AccountID

		tstmt.MustExec(tx)
		// _, err = txn.NamedExec(queryIns, tx)
		// if err != nil {
		// 	panic(err)
		// }

		// }(transaction)
	}
	// wgTrans.Wait()

	// }()

	iTok.LastDownloadedTransactions = time.Now()
	// upsertItemToken(iTok, txn)
	upsertItemToken(iTok, istmt)

	// wgWrap.Wait()

	// url := "https://www.saltedge.com/api/v5/transactions?connection_id" + ItemID
	// fmt.Println("URL:>", url)

	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	// req.Header.Set("Content-Type", "application/json")

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
}
