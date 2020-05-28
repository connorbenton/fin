package transactions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	// "fmt"
	"fintrack-go/db"
	"fintrack-go/routes/categories"
	"fintrack-go/socket"
	"fintrack-go/types"

	_ "github.com/jmoiron/sqlx"
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

// func GetRangeFunction() func(http.ResponseWriter, *http.Request) {
// 	return func(res http.ResponseWriter, req *http.Request) {

// 		dbdata := []types.Transaction{}
// 		err := db.DBCon.Select(&dbdata, "SELECT * FROM `transactions`")
// 		if err != nil {
// 			panic(err)
// 		}

// 		res.WriteHeader(http.StatusOK)
// 		if err := json.NewEncoder(res).Encode(dbdata); err != nil {
// 			panic(err)
// 		}
// 	}
// }

func PutFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		dbdata := []account{}
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

func ImportFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		baseCurrency := strings.ToUpper(os.Getenv("BASE_CURRENCY"))

		db.GetNewXML()

		p := []types.ImportTransaction{}

		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			panic(err)
		}

		var mAccs struct {
			mu       sync.Mutex
			accounts []types.MatchingAccount
		}

		var uCats struct {
			mu   sync.Mutex
			cats []string
		}

		var tSets struct {
			mu       sync.Mutex
			tSingles []types.CompareTransSingle
		}

		txArray := make([]types.Transaction, len(p))

		var wgFirst sync.WaitGroup
		for i, transaction := range p {
			wgFirst.Add(1)
			go func(i int, itx types.ImportTransaction) {
				defer wgFirst.Done()
				tx := types.Transaction{}
				if itx.Amount.IsZero() {
					return
				}
				tx.Date, err = time.Parse("01/02/2006", itx.Date)
				if err != nil {
					errString := fmt.Sprintf("Error with Import Date Parse: %v \n", err)
					log.Println(errString)
					res.WriteHeader(http.StatusInternalServerError)
					res.Write([]byte(errString))
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
					tx.Category = 106
					tx.CategoryName = "Uncategorized"
				} else {
					if v, ok := types.MintCatMap[itx.Category]; ok {
						tx.Category = v
						tx.CategoryName = itx.Category
					} else {
						sCat := types.Category{}
						query := fmt.Sprintf(`SELECT * FROM categories WHERE sub_category = %q`, itx.Category)
						err = db.DBCon.Get(&sCat, query)
						if err != nil {
							panic(err)
						}
						//Setting to nil so we can check with user later to ID the category
						if (types.Category{}) == sCat {
							tx.Category = 0
							tx.CategoryName = itx.Category
							uCats.mu.Lock()
							uCats.cats = append(uCats.cats, itx.Category)
							uCats.mu.Unlock()
						} else {
							tx.Category = sCat.ID
							tx.CategoryName = sCat.SubCategory
						}
					}
				}

				possibleMatches := []types.Transaction{}
				query := fmt.Sprintf(`SELECT * FROM transactions WHERE amount = %q`, tx.Amount)
				err = db.DBCon.Get(&possibleMatches, query)
				if err != nil {
					panic(err)
				}
				if len(possibleMatches) > 0 {
					for _, matchTx := range possibleMatches {
						if matchTx.Date == tx.Date {
							v := false
							for _, mAcc := range mAccs.accounts {
								if mAcc.RefAccount == matchTx.AccountID {
									v = true
									break
								}
							}
							if !v {
								compareSet := types.CompareTransSingle{}
								compareSet.Trans1.Date = tx.Date
								compareSet.Trans1.Description = tx.Description
								compareSet.Trans1.Amount = tx.Amount
								compareSet.Trans1.CurrencyCode = tx.CurrencyCode
								compareSet.Trans2.Date = matchTx.Date
								compareSet.Trans2.Description = matchTx.Description
								compareSet.Trans2.Amount = matchTx.Amount
								compareSet.Trans2.CurrencyCode = matchTx.CurrencyCode
								compareSet.Type = "trans"
								tSets.mu.Lock()
								tSets.tSingles = append(tSets.tSingles, compareSet)
								tSets.mu.Unlock()
								messageData, err := json.Marshal(compareSet)
								if err != nil {
									panic(err)
								}

								message := types.WsMsg{}
								message.Name = "compare"
								message.Data = messageData

								messageJSON, err := json.Marshal(message)
								if err != nil {
									panic(err)
								}

								// message := []byte(`{ "username": "Booh", }`)
								socket.ExportHub.Broadcast <- messageJSON
								answerJSON := <-socket.ExportHub.Response
								// answer := types.CompareTransAnswer{}
								// answer, err := json.Unmarshal(answerJSON)
								if string(answerJSON) == "yes" {
									mAcc := types.MatchingAccount{}
									mAcc.ImportKey = itx.AccountName
									mAcc.RefAccount = matchTx.AccountID
									mAccs.mu.Lock()
									mAccs.accounts = append(mAccs.accounts, mAcc)
									mAccs.mu.Unlock()
								}

							}
						}
					}
				}

				tx.AccountName = itx.AccountName
				txArray[i] = tx

			}(i, transaction)
		}
		wgFirst.Wait()

		messageData, err := json.Marshal(compareSet)
		if err != nil {
			panic(err)
		}

		message := types.WsMsg{}
		message.Name = "compare"
		message.Data = messageData

		messageJSON, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}

		// message := []byte(`{ "username": "Booh", }`)
		socket.ExportHub.Broadcast <- messageJSON
		answerJSON := <-socket.ExportHub.Response
		// answer := types.CompareTransAnswer{}
		// answer, err := json.Unmarshal(answerJSON)
		if string(answerJSON) == "yes" {
			mAcc := types.MatchingAccount{}
			mAcc.ImportKey = itx.AccountName
			mAcc.RefAccount = matchTx.AccountID
			mAccs.mu.Lock()
			mAccs.accounts = append(mAccs.accounts, mAcc)
			mAccs.mu.Unlock()
		}

		compareSet := types.CompareCatsSet{}
		compareSet.DbCats = categories.SelectAll()
		compareSet.CompareCats = uCats.cats
		compareSet.Type = "cats"

		messageData, err := json.Marshal(compareSet)
		if err != nil {
			panic(err)
		}

		message := types.WsMsg{}
		message.Name = "compare"
		message.Data = messageData

		messageJSON, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}

		// message := []byte(`{ "username": "Booh", }`)
		socket.ExportHub.Broadcast <- messageJSON
		answerJSON := <-socket.ExportHub.Response

		answer := []types.CompareCatsResponse{}
		if err := json.Unmarshal(answerJSON, &answer); err != nil {
			panic(err)
		}

		var wgSecond sync.WaitGroup
		for _, transaction := range txArray {
			wgSecond.Add(1)
			go func(tx types.Transaction) {
				defer wgSecond.Done()
				if tx.Category == 0 {
					for _, rCat := range answer {
						if rCat.Category == tx.CategoryName {
							tx.Category = rCat.AssignedCat
							tx.CategoryName = rCat.AssignedCatName
							break
						}
					}
				}

			}(transaction)
		}

		// dbdata := []account{}
		// // err := app.Database.Query("SELECT * FROM `categories`", id).Scan(&dbdata.id, &dbdata.topCategory, &dbdata.subCategory)
		// err := db.DBCon.Select(&dbdata, "SELECT * FROM `categories`")
		// if err != nil {
		// 	log.Fatal("Database SELECT failed")
		// 	// fmt.Println("Database SELECT failed")
		// 	// fmt.Println(err)
		// 	// return
		// }

		// log.Println("You fetched a thing!")
		// res.WriteHeader(http.StatusOK)
		// if err := json.NewEncoder(res).Encode(dbdata); err != nil {
		// 	panic(err)
		// }
	}
}
