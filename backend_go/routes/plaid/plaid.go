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
		istmt := types.PrepItemSt(txn)
		upsertItemToken(iTok, istmt)
		astmt := types.PrepAccountSt(txn)

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

func upsertItemToken(tok types.ItemToken, stmt *sqlx.NamedStmt) {

	stmt.MustExec(tok)

}

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

		if pRes.PublicToken == "" {
			errString := fmt.Sprintf("PublicToken response seems to be empty")
			log.Println(errString)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(errString))
		} else {
			resJSON := `{"public_token": "` + pRes.PublicToken + `"}`
			_, err2 := res.Write([]byte(resJSON))
			if err2 != nil {
				errString := fmt.Sprintf("Error with Plaid Generate Public Token: %v \n", err2)
				log.Println(errString)
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(errString))
			}
		}

	}
}

func RefreshConnection(iTok types.ItemToken, istmt, astmt *sqlx.NamedStmt) {

	pClient, err := newClient()
	if err != nil {
		panic(err)
	}

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
			upsertAccountWithPlaidAccount(acc, pAcc, iTok.Institution, iTok.ItemID, astmt)
		}(account)
	}
	wgAcc.Wait()
}

func FetchTransactionsForItemToken(iTok types.ItemToken, istmt *sqlx.NamedStmt, astmt *sqlx.NamedStmt, tstmt *sqlx.NamedStmt, baseCurrency string) {
	today := time.Now().Format("2006-01-02")

	pClient, err := newClient()
	if err != nil {
		panic(err)
	}

	var pTransRes plaid.GetTransactionsResponse
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
		panic(err)
	}

	for _, ptx := range pTransRes.Transactions {
		tx := types.Transaction{}

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
	}

	iTok.LastDownloadedTransactions = time.Now()
	upsertItemToken(iTok, istmt)
}
