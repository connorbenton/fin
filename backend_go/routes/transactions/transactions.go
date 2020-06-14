package transactions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"fintrack-go/db"
	"fintrack-go/routes/accounts"
	"fintrack-go/routes/analysisTrees"
	"fintrack-go/types"

	_ "github.com/jmoiron/sqlx"
	"github.com/rickb777/date"
	"github.com/shopspring/decimal"
)

type account struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Institution string    `json:"institution"`
	Account_id  string    `json:"account_id"`
	Item_id     string    `json:"item_id"`
	Type        string    `json:"type"`
	Subtype     string    `json:"subtype"`
	Balance     float32   `json:"balance"`
	Limit       float32   `json:"limit"`
	Available   float32   `json:"available"`
	Currency    string    `json:"currency"`
	Provider    string    `json:"provider"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func SelectAll() []types.Transaction {
	dbdata := []types.Transaction{}
	err := db.DBCon.Select(&dbdata, "SELECT * FROM `transactions`")
	if err != nil {
		panic(err)
	}
	return dbdata
}

func GetFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		start := time.Now()
		dbdata := SelectAll()
		// log.Println("SelectAll Transactions done in:", time.Since(start))

		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(dbdata); err != nil {
			panic(err)
		}
		log.Println("SelectAll Transactions done in:", time.Since(start))
		// log.Println("Encode done:", time.Since(start))
	}
}

func PutFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		txn := db.DBCon.MustBegin()
		tstmt := types.PrepTransSt(txn)

		p := types.Transaction{}

		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			panic(err)
		}

		tstmt.MustExec(p)

		errC := txn.Commit()
		if errC != nil {
			panic(errC)
		}

		res.WriteHeader(http.StatusOK)
	}
}

func CheckFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		p := []types.ImportTransaction{}

		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			panic(err)
		}

		var resJSON struct {
			TSets struct {
				mu       sync.Mutex                 `json:"-"`
				TSingles []types.CompareTransSingle `json:"transSets"`
			} `json:"trans"`
			UCats struct {
				mu   sync.Mutex                `json:"-"`
				Cats []types.CompareCatsSingle `json:"catsToID"`
			} `json:"cats"`
		}

		txArray := make([]types.Transaction, len(p))

		for i, itx := range p {
			tx := types.Transaction{}
			if itx.Amount.IsZero() {
				continue
			}
			dt, err := date.Parse("01/02/2006", itx.Date)
			if err != nil {
				dt, err := date.Parse("1/02/2006", itx.Date)
				if err != nil {
					dt, err := date.Parse("2006-01-02", itx.Date)
					if err != nil {
						errString := fmt.Sprintf("Error with Import Date Parse: %v \n", err)
						log.Println(errString)
						res.WriteHeader(http.StatusInternalServerError)
						res.Write([]byte(errString))
					} else {
						tx.Date = dt.String()
					}
				} else {
					tx.Date = dt.String()
				}
			} else {
				tx.Date = dt.String()
			}

			tx.Description = itx.Description
			if itx.TransactionType == "debit" {
				tx.Amount = itx.Amount.Mul(decimal.NewFromInt(-1))
			} else {
				tx.Amount = itx.Amount
			}
			if itx.CurrencyCode == "" {
				tx.CurrencyCode = "USD"
			} else {
				tx.CurrencyCode = itx.CurrencyCode
			}

			if itx.Category == "" {
			} else {
				if _, ok := types.MintCatMap[itx.Category]; ok {
				} else {
					sCat := types.Category{}
					query := fmt.Sprintf(`SELECT * FROM categories WHERE sub_category = %q`, itx.Category)
					err = db.DBCon.Get(&sCat, query)
					if err != nil && err != sql.ErrNoRows {
						panic(err)
					}
					// Setting to nil so we can check with user later to ID the category
					if (types.Category{}) == sCat {
						v := false
						for _, cat := range resJSON.UCats.Cats {
							if cat.Category == itx.Category {
								v = true
								break
							}
						}
						if !v {
							compareSet := types.CompareCatsSingle{}
							compareSet.Category = itx.Category
							compareSet.AssignedCat = 106
							compareSet.AssignedCatName = "Uncategorized"
							resJSON.UCats.Cats = append(resJSON.UCats.Cats, compareSet)
						}
					} else {
					}
				}
			}

			possibleMatches := []types.Transaction{}
			query := fmt.Sprintf(`SELECT * FROM transactions WHERE amount = %q AND date = %q`, tx.Amount, tx.Date)
			err = db.DBCon.Select(&possibleMatches, query)
			if err != nil && err != sql.ErrNoRows {
				panic(err)
			}
			if len(possibleMatches) > 0 {
				for _, matchTx := range possibleMatches {

					compareSet := types.CompareTransSingle{}
					compareSet.Trans1.Date = tx.Date
					compareSet.Trans1.Description = tx.Description
					compareSet.Trans1.Amount = tx.Amount
					compareSet.Trans1.CurrencyCode = tx.CurrencyCode
					compareSet.Trans1.AccountName = itx.AccountName
					compareSet.Trans2.Date = matchTx.Date
					compareSet.Trans2.Description = matchTx.Description
					compareSet.Trans2.Amount = matchTx.Amount
					compareSet.Trans2.CurrencyCode = matchTx.CurrencyCode
					compareSet.Trans2.AccountName = matchTx.AccountName
					compareSet.Trans2.AccountID = matchTx.AccountID
					compareSet.Type = "trans"

					resJSON.TSets.TSingles = append(resJSON.TSets.TSingles, compareSet)

				}
			}

			txArray[i] = tx

		}

		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(resJSON); err != nil {
			panic(err)
		}

	}
}

func ImportFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		db.GetNewXML()
		baseCurrency := strings.ToUpper(os.Getenv("BASE_CURRENCY"))

		var cAccs struct {
			mu       sync.Mutex
			accounts []types.Account
		}

		var countInt struct {
			mu         sync.Mutex
			countDup   int
			countUncat int
			countImp   int
		}

		p := types.ImportPostData{}

		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			panic(err)
		}

		dbAccs := accounts.SelectAll()

		txn := db.DBCon.MustBegin()
		astmt := types.PrepAccountSt(txn)
		tstmt := types.PrepTransSt(txn)

		for _, itx := range p.TxSet {
			tx := types.Transaction{}
			if itx.Amount.IsZero() {
				continue
			}
			dt, err := date.Parse("01/02/2006", itx.Date)
			if err != nil {
				dt, err := date.Parse("1/02/2006", itx.Date)
				if err != nil {
					dt, err := date.Parse("2006-01-02", itx.Date)
					if err != nil {
						errString := fmt.Sprintf("Error with Import Date Parse: %v \n", err)
						log.Println(errString)
						res.WriteHeader(http.StatusInternalServerError)
						res.Write([]byte(errString))
					} else {
						tx.Date = dt.String()
					}
				} else {
					tx.Date = dt.String()
				}
			} else {
				tx.Date = dt.String()
			}
			tx.Description = itx.Description
			if itx.TransactionType == "debit" {
				tx.Amount = itx.Amount.Mul(decimal.NewFromInt(-1))
			} else {
				tx.Amount = itx.Amount
			}

			tx.Description = itx.Description
			if itx.TransactionType == "debit" {
				tx.Amount = itx.Amount.Mul(decimal.NewFromInt(-1))
			} else {
				tx.Amount = itx.Amount
			}
			if itx.CurrencyCode == "" {
				tx.CurrencyCode = "USD"
			} else {
				tx.CurrencyCode = itx.CurrencyCode
			}
			tx.NormalizedAmount = db.GetNormalizedAmount(tx.CurrencyCode, baseCurrency, tx.Date, tx.Amount)

			if itx.Category == "" {
				countInt.countUncat++
				tx.Category = 106
				tx.CategoryName = "Uncategorized"
			} else {
				if v, ok := types.MintCatMap[itx.Category]; ok {
					tx.Category = v
					var st string
					query := fmt.Sprintf(`Select sub_category FROM categories WHERE id = %d`, tx.Category)
					// log.Println(query)
					err = db.DBCon.Get(&st, query)
					if err != nil {
						panic(err)
					}
					tx.CategoryName = st
				} else {
					sCat := types.Category{}
					query := fmt.Sprintf(`SELECT * FROM categories WHERE sub_category = %q`, itx.Category)
					err = db.DBCon.Get(&sCat, query)
					if err != nil && err != sql.ErrNoRows {
						panic(err)
					}
					if (types.Category{}) == sCat {
						v := false
						for _, rCat := range p.Catres {
							if rCat.Category == itx.Category {
								tx.Category = rCat.AssignedCat
								tx.CategoryName = rCat.AssignedCatName
								v = true
								break
							}
						}
						if !v {
							countInt.countUncat++
							tx.Category = 106
							tx.CategoryName = "Uncategorized"
						}
					} else {
						tx.Category = sCat.ID
						tx.CategoryName = sCat.SubCategory
					}
				}
			}

			var newTID string
			for {
				v := true
				newTID = strconv.FormatInt(rand.Int63(), 10)
				possibleMatches := []types.Transaction{}
				query := fmt.Sprintf(`SELECT * FROM transactions WHERE transaction_id = %q`, newTID)
				err = db.DBCon.Select(&possibleMatches, query)
				if err != nil && err != sql.ErrNoRows {
					panic(err)
				}
				if len(possibleMatches) > 0 {
					v = false
				}
				if v {
					break
				}
			}
			tx.TransactionID = newTID

			v := false
			for _, acc := range cAccs.accounts {
				if acc.Name == itx.AccountName {
					tx.AccountID = acc.AccountID
					tx.AccountName = acc.Name
					v = true
					break
				}
			}
			if !v {
				v := false
				for _, acc := range p.IdentifiedAccounts {
					if acc.ImportKey == itx.AccountName {
						tx.AccountID = acc.RefAccountID
						tx.AccountName = acc.RefAccountName
						v = true
						break
					}
				}
				if !v {
					var newID string
					for {
						newID = strconv.FormatInt(rand.Int63(), 10)
						v := true
						for _, acc := range dbAccs {
							if acc.AccountID == newID {
								v = false
								break
							}
						}
						for _, acc := range cAccs.accounts {
							if acc.AccountID == newID {
								v = false
								break
							}
						}
						if v {
							break
						}
					}
					accountToCreate := types.Account{}
					accountToCreate.AccountID = newID
					tx.AccountID = newID
					tx.AccountName = itx.AccountName
					accountToCreate.Name = itx.AccountName
					accountToCreate.Institution = "Import"
					accountToCreate.Provider = "Import"
					cAccs.accounts = append(cAccs.accounts, accountToCreate)
				}
			}
			possibleMatches := []types.Transaction{}
			query := fmt.Sprintf(`SELECT * FROM transactions WHERE amount = %q AND date = %q AND account_id = %q `, tx.Amount, tx.Date, tx.AccountID)
			err = db.DBCon.Select(&possibleMatches, query)
			if err != nil {
				panic(err)
			}
			if len(possibleMatches) < 1 {
				tstmt.MustExec(tx)
				countInt.countImp++
			} else {
				countInt.countDup++
			}
		}

		for _, acc := range cAccs.accounts {
			astmt.MustExec(acc)
		}

		log.Println("duplicate number in import = " + strconv.Itoa(countInt.countDup))
		log.Println("uncategorized number in import = " + strconv.Itoa(countInt.countUncat))
		log.Println("total transactions imported = " + strconv.Itoa(countInt.countImp))

		errC := txn.Commit()
		if errC != nil {
			panic(errC)
		}

		analysisTrees.ReAnalyze()

		res.WriteHeader(http.StatusOK)

	}
}

func UpsertFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		p := []types.Transaction{}

		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			panic(err)
		}

		txn := db.DBCon.MustBegin()
		tstmt := types.PrepTransUpsertSt(txn)

		for _, tx := range p {

			tstmt.MustExec(tx)

		}

		errC := txn.Commit()
		if errC != nil {
			panic(errC)
		}

		analysisTrees.ReAnalyze()

		res.WriteHeader(http.StatusOK)
	}
}
