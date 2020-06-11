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

	// "fmt"
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
		// err := db.DBCon.Select(&dbdata, "SELECT * FROM `item_tokens`")
		// if err != nil {
		// log.Fatal(err)
		// }

		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(dbdata); err != nil {
			panic(err)
		}
		log.Println("SelectAll Transactions done in:", time.Since(start))
		// log.Println("Encode done:", time.Since(start))
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

		// var mAccs struct {
		// mu       sync.Mutex
		// accounts []types.MatchingAccount `json:"accounts"`
		// }

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

		// var wgFirst sync.WaitGroup
		// for i, transaction := range p {
		for i, itx := range p {
			// wgFirst.Add(1)
			// go func(i int, itx types.ImportTransaction) {
			// defer wgFirst.Done()
			tx := types.Transaction{}
			if itx.Amount.IsZero() {
				continue
			}
			dt, err := date.Parse("01/02/2006", itx.Date)
			if err != nil {
				dt, err := date.Parse("1/02/2006", itx.Date)
				if err != nil {
					errString := fmt.Sprintf("Error with Import Date Parse: %v \n", err)
					log.Println(errString)
					res.WriteHeader(http.StatusInternalServerError)
					res.Write([]byte(errString))
				} else {
					tx.Date = dt.Format("2006-01-02")
				}
			} else {
				tx.Date = dt.Format("2006-01-02")
			}

			// dt, err := date.Parse("01/02/2006", itx.Date)
			// if err != nil {
			// 	dt, err := date.Parse("1/02/2006", itx.Date)
			// 	if err != nil {
			// 		errString := fmt.Sprintf("Error with Import Date Parse: %v \n", err)
			// 		log.Println(errString)
			// 		res.WriteHeader(http.StatusInternalServerError)
			// 		res.Write([]byte(errString))
			// 	} else {
			// 		tx.Date = dt.Format("2006-01-02")
			// 	}
			// } else {
			// 	tx.Date = dt.Format("2006-01-02")
			// }
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
			// tx.NormalizedAmount = db.GetNormalizedAmount(tx.CurrencyCode, baseCurrency, tx.Date, tx.Amount)

			if itx.Category == "" {
				// tx.Category = 106
				// tx.CategoryName = "Uncategorized"
			} else {
				if _, ok := types.MintCatMap[itx.Category]; ok {
					// tx.Category = v
					// tx.CategoryName = itx.Category
				} else {
					sCat := types.Category{}
					query := fmt.Sprintf(`SELECT * FROM categories WHERE sub_category = %q`, itx.Category)
					err = db.DBCon.Get(&sCat, query)
					if err != nil && err != sql.ErrNoRows {
						panic(err)
					}
					// Setting to nil so we can check with user later to ID the category
					if (types.Category{}) == sCat {
						// if err == sql.ErrNoRows {
						// tx.Category = 0
						// tx.CategoryName = itx.Category
						// UCats.mu.Lock()
						v := false
						for _, cat := range resJSON.UCats.Cats {
							if cat.Category == itx.Category {
								v = true
								break
							}
						}
						if !v {
							// UCats.mu.Lock()
							compareSet := types.CompareCatsSingle{}
							compareSet.Category = itx.Category
							compareSet.AssignedCat = 106
							compareSet.AssignedCatName = "Uncategorized"
							resJSON.UCats.Cats = append(resJSON.UCats.Cats, compareSet)
							// UCats.mu.Unlock()
						}
					} else {
						// tx.Category = sCat.ID
						// tx.CategoryName = sCat.SubCategory
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
					// if matchTx.Date == tx.Date {
					// TSets.mu.Lock()
					// v := false
					// for _, mAcc := range mAccs.accounts {
					// 	if mAcc.RefAccountID == matchTx.AccountID {
					// 		v = true
					// 		// TSets.mu.Unlock()
					// 		break
					// 	}
					// }
					// if !v {
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
					// TSets.mu.Lock()
					resJSON.TSets.TSingles = append(resJSON.TSets.TSingles, compareSet)
					// TSets.mu.Unlock()
					// messageData, err := json.Marshal(compareSet)
					// if err != nil {
					// 	panic(err)
					// }

					// message := types.WsMsg{}
					// message.Name = "compare"
					// message.Data = messageData

					// messageJSON, err := json.Marshal(message)
					// if err != nil {
					// 	panic(err)
					// }

					// // message := []byte(`{ "username": "Booh", }`)
					// socket.ExportHub.Broadcast <- messageJSON
					// answerJSON := <-socket.ExportHub.Response
					// // answer := types.CompareTransAnswer{}
					// // answer, err := json.Unmarshal(answerJSON)
					// if string(answerJSON) == "yes" {
					// 	mAcc := types.MatchingAccount{}
					// mAcc.ImportKey = itx.AccountName
					// 	mAcc.RefAccountID = matchTx.AccountID
					// 	mAcc.RefAccountName = matchTx.AccountName
					// 	mAccs.mu.Lock()
					// 	mAccs.accounts = append(mAccs.accounts, mAcc)
					// 	mAccs.mu.Unlock()
					// }

					// TSets.mu.Unlock()

					// }
					// TSets.mu.Unlock()
					// }
				}
			}

			// tx.AccountName = itx.AccountName
			txArray[i] = tx

			// }(i, transaction)
		}
		// wgFirst.Wait()

		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(resJSON); err != nil {
			panic(err)
		}
		// if err := json.NewEncoder(res).Encode(TSets.TSingles); err != nil {
		// 	panic(err)
		// }
		// if err := json.NewEncoder(res).Encode(UCats.Cats); err != nil {
		// 	panic(err)
		// }

		// for _, set := range TSets.TSingles {
		// 	// compareSet := types.CompareTransSingle{}
		// 	// compareSet.Trans1.Date = tx.Date
		// 	// compareSet.Trans1.Description = tx.Description
		// 	// compareSet.Trans1.Amount = tx.Amount
		// 	// compareSet.Trans1.CurrencyCode = tx.CurrencyCode
		// 	// compareSet.Trans2.Date = matchTx.Date
		// 	// compareSet.Trans2.Description = matchTx.Description
		// 	// compareSet.Trans2.Amount = matchTx.Amount
		// 	// compareSet.Trans2.CurrencyCode = matchTx.CurrencyCode
		// 	// compareSet.Type = "trans"
		// 	// TSets.TSingles = append(TSets.TSingles, compareSet)
		// 	// TSets.mu.Unlock()
		// 	messageData, err := json.Marshal(set)
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	message := types.WsMsg{}
		// 	message.Name = "compare"
		// 	message.Data = messageData

		// 	messageJSON, err := json.Marshal(message)
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	// message := []byte(`{ "username": "Booh", }`)
		// 	socket.ExportHub.Broadcast <- messageJSON
		// 	answerJSON := <-socket.ExportHub.Response
		// 	// answer := types.CompareTransAnswer{}
		// 	// answer, err := json.Unmarshal(answerJSON)
		// 	if string(answerJSON) == "yes" {
		// 		// mAcc := types.MatchingAccount{}
		// 		// mAcc.ImportKey = itx.AccountName
		// 		// mAcc.RefAccountID = matchTx.AccountID
		// 		// mAcc.RefAccountName = matchTx.AccountName
		// 		// mAccs.mu.Lock()
		// 		// mAccs.accounts = append(mAccs.accounts, mAcc)
		// 		// mAccs.mu.Unlock()
		// 	}
		// }

		// // for _, mAcc := range mAccs.accounts {
		// // if mAcc.RefAccountID == matchTx.AccountID {
		// // 	v = true
		// // 	// TSets.mu.Unlock()
		// // 	break
		// // }
		// // }

		// // messageData, err := json.Marshal(compareSet)
		// // if err != nil {
		// // 	panic(err)
		// // }

		// // message := types.WsMsg{}
		// // message.Name = "compare"
		// // message.Data = messageData

		// // messageJSON, err := json.Marshal(message)
		// // if err != nil {
		// // 	panic(err)
		// // }

		// // // message := []byte(`{ "username": "Booh", }`)
		// // socket.ExportHub.Broadcast <- messageJSON
		// // answerJSON := <-socket.ExportHub.Response
		// // // answer := types.CompareTransAnswer{}
		// // // answer, err := json.Unmarshal(answerJSON)
		// // if string(answerJSON) == "yes" {
		// // 	mAcc := types.MatchingAccount{}
		// // mAcc.ImportKey = itx.AccountName
		// // 	mAcc.RefAccount = matchTx.AccountID
		// // 	mAccs.mu.Lock()
		// // 	mAccs.accounts = append(mAccs.accounts, mAcc)
		// // 	mAccs.mu.Unlock()
		// // }
		// answer := []types.CompareCatsResponse{}
		// if len(UCats.Cats) > 0 {

		// 	compareSet := types.CompareCatsSet{}
		// 	compareSet.DbCats = categories.SelectAll()
		// 	compareSet.CompareCats = UCats.Cats
		// 	compareSet.Type = "cats"

		// 	messageData, err := json.Marshal(compareSet)
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	message := types.WsMsg{}
		// 	message.Name = "compare"
		// 	message.Data = messageData

		// 	messageJSON, err := json.Marshal(message)
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	// message := []byte(`{ "username": "Booh", }`)
		// 	socket.ExportHub.Broadcast <- messageJSON
		// 	answerJSON := <-socket.ExportHub.Response

		// 	// answer := []types.CompareCatsResponse{}
		// 	if err := json.Unmarshal(answerJSON, &answer); err != nil {
		// 		panic(err)
		// 	}
		// }
		// var cAccs struct {
		// 	mu       sync.Mutex
		// 	accounts []types.Account
		// }

		// var countInt struct {
		// 	mu         sync.Mutex
		// 	countDup   int
		// 	countUncat int
		// }

		// dbAccs := accounts.SelectAll()

		// txn := db.DBCon.MustBegin()

		// astmt := types.PrepAccountSt(txn)
		// tstmt := types.PrepTransSt(txn)

		// var wgSecond sync.WaitGroup
		// for _, transaction := range txArray {
		// 	wgSecond.Add(1)
		// 	go func(tx types.Transaction) {
		// 		defer wgSecond.Done()
		// 		if tx.Category == 0 {
		// 			v := false
		// 			for _, rCat := range answer {
		// 				if rCat.Category == tx.CategoryName {
		// 					tx.Category = rCat.AssignedCat
		// 					tx.CategoryName = rCat.AssignedCatName
		// 					v = true
		// 					break
		// 				}
		// 			}
		// 			if !v {
		// 				countInt.mu.Lock()
		// 				countInt.countUncat++
		// 				countInt.mu.Unlock()
		// 				tx.Category = 106
		// 				tx.CategoryName = "Uncategorized"
		// 			}
		// 		}
		// 		cAccs.mu.Lock()
		// 		v := false
		// 		for _, acc := range cAccs.accounts {
		// 			if acc.Name == tx.AccountName {
		// 				tx.AccountID = acc.AccountID
		// 				v = true
		// 				break
		// 			}
		// 		}
		// 		if !v {
		// 			v := false
		// 			for _, acc := range mAccs.accounts {
		// 				if acc.ImportKey == tx.AccountName {
		// 					tx.AccountID = acc.RefAccountID
		// 					tx.AccountName = acc.RefAccountName
		// 					v = true
		// 					break
		// 				}
		// 			}
		// 			if !v {
		// 				var newID string
		// 				for {
		// 					newID = strconv.FormatInt(rand.Int63(), 10)
		// 					v := true
		// 					for _, acc := range dbAccs {
		// 						if acc.AccountID == newID {
		// 							v = false
		// 							break
		// 						}
		// 					}
		// 					for _, acc := range cAccs.accounts {
		// 						if acc.AccountID == newID {
		// 							v = false
		// 							break
		// 						}
		// 					}
		// 					if v {
		// 						break
		// 					}
		// 				}
		// 				accountToCreate := types.Account{}
		// 				accountToCreate.AccountID = newID
		// 				tx.AccountID = newID
		// 				accountToCreate.Name = tx.AccountName
		// 				accountToCreate.Institution = "Import"
		// 				accountToCreate.Provider = "Import"
		// 				cAccs.accounts = append(cAccs.accounts, accountToCreate)
		// 			}
		// 		}
		// 		cAccs.mu.Unlock()
		// 		possibleMatches := []types.Transaction{}
		// 		query := fmt.Sprintf(`SELECT * FROM transactions WHERE amount = %q AND 'date' = %q AND account_id = %q `, tx.Amount, tx.Date, tx.AccountID)
		// 		err = db.DBCon.Get(&possibleMatches, query)
		// 		if err != nil {
		// 			panic(err)
		// 		}
		// 		if len(possibleMatches) < 1 {
		// 			tstmt.MustExec(tx)
		// 		} else {
		// 			countInt.mu.Lock()
		// 			countInt.countDup++
		// 			countInt.mu.Unlock()
		// 		}
		// 	}(transaction)
		// }

		// for _, acc := range cAccs.accounts {
		// 	astmt.MustExec(acc)
		// }

		// wgSecond.Wait()

		// log.Println("duplicate number in import = " + strconv.Itoa(countInt.countDup))
		// log.Println("uncategorized number in import = " + strconv.Itoa(countInt.countUncat))

		// errC := txn.Commit()
		// if errC != nil {
		// 	panic(errC)
		// }

		// res.WriteHeader(http.StatusOK)

		// // dbdata := []account{}
		// // // err := app.Database.Query("SELECT * FROM `categories`", id).Scan(&dbdata.id, &dbdata.topCategory, &dbdata.subCategory)
		// // err := db.DBCon.Select(&dbdata, "SELECT * FROM `categories`")
		// // if err != nil {
		// // 	log.Fatal("Database SELECT failed")
		// // 	// fmt.Println("Database SELECT failed")
		// // 	// fmt.Println(err)
		// // 	// return
		// // }

		// // log.Println("You fetched a thing!")
		// // res.WriteHeader(http.StatusOK)
		// // if err := json.NewEncoder(res).Encode(dbdata); err != nil {
		// // 	panic(err)
		// // }
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

		// var wgSecond sync.WaitGroup
		// for _, transaction := range p.TxSet {
		for _, itx := range p.TxSet {
			// wgSecond.Add(1)
			// go func(itx types.ImportTransaction) {
			// defer wgSecond.Done()
			tx := types.Transaction{}
			if itx.Amount.IsZero() {
				// return
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
						tx.Date = dt.Format("2006-01-02")
					}
				} else {
					tx.Date = dt.Format("2006-01-02")
				}
			} else {
				tx.Date = dt.Format("2006-01-02")
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
				// countInt.mu.Lock()
				countInt.countUncat++
				// countInt.mu.Unlock()
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
							// countInt.mu.Lock()
							countInt.countUncat++
							// countInt.mu.Unlock()
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

			// cAccs.mu.Lock()
			v := false
			for _, acc := range cAccs.accounts {
				if acc.Name == itx.AccountName {
					tx.AccountID = acc.AccountID
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
			// cAccs.mu.Unlock()
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
				// countInt.mu.Lock()
				countInt.countDup++
				// countInt.mu.Unlock()
			}
			// }(transaction)
		}

		for _, acc := range cAccs.accounts {
			astmt.MustExec(acc)
		}

		// wgSecond.Wait()

		log.Println("duplicate number in import = " + strconv.Itoa(countInt.countDup))
		log.Println("uncategorized number in import = " + strconv.Itoa(countInt.countUncat))
		log.Println("total transactions imported = " + strconv.Itoa(countInt.countImp))

		errC := txn.Commit()
		if errC != nil {
			panic(errC)
		}

		analysisTrees.ReAnalyze()

		res.WriteHeader(http.StatusOK)

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
