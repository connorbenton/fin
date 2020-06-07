package analysisTrees

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	// "fmt"
	"fintrack-go/db"
	"fintrack-go/types"

	_ "github.com/jmoiron/sqlx"
	"github.com/rickb777/date"

	// "github.com/shopspring/decimal"
	"github.com/shopspring/decimal"
)

func GetFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		dbdata := []types.Tree{}
		err := db.DBCon.Select(&dbdata, "SELECT * FROM `analysis_trees`")
		if err != nil {
			panic(err)
		}

		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(dbdata); err != nil {
			panic(err)
		}
	}
}
func CustomAnalyze() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var err error

		customRange := types.CustomRange{}

		err = json.NewDecoder(req.Body).Decode(&customRange)
		if err != nil {
			panic(err)
		}

		dbcatsBase := []types.Category{}

		// dbdata := []types.Transaction{}
		err = db.DBCon.Select(&dbcatsBase, "SELECT * FROM `categories`")
		if err != nil {
			panic(err)
		}

		txn := db.DBCon.MustBegin()
		tstmt := types.PrepTreeSt(txn)

		rangedata := []types.Transaction{}
		var st, end, dbEnd string
		st = customRange.Start
		end = customRange.End
		dbEndDate, _ := date.ParseISO(customRange.End)
		dbEnd = dbEndDate.AddDate(0, 0, 1).Format("2006-01-02")
		err = db.DBCon.Select(&rangedata, "SELECT * FROM `transactions` WHERE DATE between '"+st+"' and '"+dbEnd+"'")
		// log.Println("number in ", name, len(rangedata))
		if err != nil {
			log.Println(end)
			panic(err)
		}

		tree := SetupTree(dbcatsBase, rangedata, "custom", st, end)
		tstmt.MustExec(tree)

		errC := txn.Commit()
		if errC != nil {
			panic(errC)
		}

		res.WriteHeader(http.StatusOK)
		if errW := json.NewEncoder(res).Encode(tree); errW != nil {
			panic(errW)
		}
	}
}

func ReAnalyze() {

	start := time.Now()
	dbcatsBase := []types.Category{}

	// dbdata := []types.Transaction{}
	err := db.DBCon.Select(&dbcatsBase, "SELECT * FROM `categories`")
	if err != nil {
		panic(err)
	}

	txn := db.DBCon.MustBegin()

	tstmt := types.PrepTreeSt(txn)

	log.Println("First select of tx categories:", time.Since(start))

	var wg sync.WaitGroup
	for _, name := range types.TreeRanges {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			rangedata := []types.Transaction{}
			today := date.Today()
			var st, end, dbEnd string
			var err error
			switch name {
			case "last30":
				st = today.AddDate(0, 0, -29).Format("2006-01-02")
				end = today.Format("2006-01-02")
				dbEnd = today.AddDate(0, 0, 1).Format("2006-01-02")
			case "thisMonth":
				st = date.New(today.Year(), today.Month(), 1).Format("2006-01-02")
				end = today.Format("2006-01-02")
				dbEnd = today.AddDate(0, 0, 1).Format("2006-01-02")
			case "lastMonth":
				stdate := today.AddDate(0, -1, 0)
				st = date.New(stdate.Year(), stdate.Month(), 1).Format("2006-01-02")
				enddate := date.New(stdate.Year(), stdate.Month(), 1).AddDate(0, 1, -1)
				end = enddate.Format("2006-01-02")
				dbEnd = enddate.AddDate(0, 0, 1).Format("2006-01-02")
			case "last6Months":
				st = today.AddDate(0, -6, 0).Format("2006-01-02")
				end = today.Format("2006-01-02")
				dbEnd = today.AddDate(0, 0, 1).Format("2006-01-02")
			case "thisYear":
				st = date.New(today.Year(), 1, 1).Format("2006-01-02")
				end = today.Format("2006-01-02")
				dbEnd = today.AddDate(0, 0, 1).Format("2006-01-02")
			case "lastYear":
				stdate := today.AddDate(-1, 0, 0)
				st = date.New(stdate.Year(), 1, 1).Format("2006-01-02")
				enddate := date.New(stdate.Year(), 1, 1).AddDate(1, 0, -1)
				end = enddate.Format("2006-01-02")
				dbEnd = enddate.AddDate(0, 0, 1).Format("2006-01-02")
			case "fromBeginning":
				err2 := db.DBCon.Get(&st, "SELECT MIN(date) FROM `transactions`")
				if err2 != nil {
					panic(err2)
				}
				end = today.Format("2006-01-02")
				dbEnd = today.AddDate(0, 0, 1).Format("2006-01-02")
				// err = db.DBCon.Select(&rangedata, "SELECT * FROM `transactions`")
				// log.Println("number in ", name, len(rangedata))
			case "custom":
				st = today.AddDate(0, 0, -29).Format("2006-01-02")
				end = today.Format("2006-01-02")
				dbEnd = today.AddDate(0, 0, 1).Format("2006-01-02")
			}
			err = db.DBCon.Select(&rangedata, "SELECT * FROM `transactions` WHERE DATE between '"+st+"' and '"+dbEnd+"'")
			// log.Println("number in ", name, len(rangedata))
			if err != nil {
				log.Println(end)
				panic(err)
			}

			tree := SetupTree(dbcatsBase, rangedata, name, st, end)

			// dbcats := append(dbcatsBase[:0:0], dbcatsBase...)
			// // for _, tx := range rangedata {
			// for _, tx := range rangedata {
			// 	for i := range dbcats {
			// 		if tx.Category == dbcats[i].ID && !tx.NormalizedAmount.IsZero() {
			// 			dbcats[i].Count++
			// 			dbcats[i].Total = dbcats[i].Total.Add(tx.NormalizedAmount)
			// 			break
			// 		}
			// 	}
			// }

			// var Zero = decimal.New(0, 1)
			// dbcatsNoInvest := append(dbcats[:0:0], dbcats...)
			// for i := range dbcatsNoInvest {
			// 	if dbcatsNoInvest[i].ID == 69 {
			// 		dbcatsNoInvest[i].Count = 0
			// 		dbcatsNoInvest[i].Total = Zero
			// 		break
			// 	}
			// }

			// tree := types.Tree{}
			// tree.Name = name
			// tree.FirstDate = st
			// tree.LastDate = end
			// tree.Data = GenerateDataTree(dbcats)
			// tree.DataNoInvest = GenerateDataTree(dbcatsNoInvest)

			tstmt.MustExec(tree)

			// log.Println(dbcats[103])

			// for _, cat := range dbcats {
			// 	var err error
			// 	err = db.DBCon.Get(&cat.Count, "SELECT COUNT(id) FROM `transactions` WHERE DATE between '"+st+"' and '"+dbEnd+"' AND category = '"+strconv.Itoa(cat.ID)+"'")
			// 	err = db.DBCon.Get(&cat.Total, "SELECT SUM(normalized_amount) FROM `transactions` WHERE DATE between '"+st+"' and '"+dbEnd+"' AND category = '"+strconv.Itoa(cat.ID)+"'")
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// 	// log.Println("Count ", cat.Count, " Total ", cat.Total)
			// }

		}(name)
	}

	wg.Wait()
	errC := txn.Commit()
	if errC != nil {
		panic(errC)
	}
	log.Println("Reanalyze tx ranges done in:", time.Since(start))
}

func SetupTree(dbcatsBase []types.Category, rangedata []types.Transaction, name, st, end string) types.Tree {
	dbcats := append(dbcatsBase[:0:0], dbcatsBase...)
	// for _, tx := range rangedata {
	for _, tx := range rangedata {
		for i := range dbcats {
			if tx.Category == dbcats[i].ID && !tx.NormalizedAmount.IsZero() {
				dbcats[i].Count++
				dbcats[i].Total = dbcats[i].Total.Add(tx.NormalizedAmount)
				break
			}
		}
	}

	var Zero = decimal.New(0, 1)
	dbcatsNoInvest := append(dbcats[:0:0], dbcats...)
	for i := range dbcatsNoInvest {
		if dbcatsNoInvest[i].ID == 69 {
			dbcatsNoInvest[i].Count = 0
			dbcatsNoInvest[i].Total = Zero
			break
		}
	}

	tree := types.Tree{}
	tree.Name = name
	tree.FirstDate = st
	tree.LastDate = end
	var totalcount int
	tree.Data, totalcount = GenerateDataTree(dbcats, 0, false)
	tree.DataNoInvest, _ = GenerateDataTree(dbcatsNoInvest, totalcount, true)
	return tree
}

func GenerateDataTree(dbcats []types.Category, trueTotal int, useTrueTotalFlag bool) (string, int) {

	var Zero = decimal.New(0, 1)
	treeData := types.TreeData{}

	// tree.Data = treeData
	treeData.Name = "Transactions by Category"
	treeData.Children = []types.ChildTop{}
	treeData.Value = Zero
	// treeData.TrueValue = Zero
	treeData.Count = 0
	treeData.TrueCount = 0

	for _, cat := range dbcats {
		if cat.SubCategory == cat.TopCategory {
			topChild := types.ChildTop{}
			topChild.Name = cat.TopCategory
			topChild.Children = []types.ChildSub{}
			topChild.Value = Zero
			// topChild.TrueValue = Zero
			topChild.Count = 0
			topChild.TrueCount = 0
			topChild.DbID = 0

			for _, cat2 := range dbcats {
				if cat2.TopCategory == cat.TopCategory {
					subChild := types.ChildSub{}
					subChild.Name = ""
					subChild.Value = Zero
					// subChild.TrueValue = Zero
					subChild.Count = 0
					subChild.TrueCount = 0
					subChild.DbID = 0

					if cat2.SubCategory == cat2.TopCategory {
						subChild.Name = cat2.SubCategory + " (General)"
						topChild.DbID = cat2.ID
					} else {
						subChild.Name = cat2.SubCategory
					}

					subChild.DbID = cat2.ID
					subChild.Value = cat2.Total.Mul(decimal.NewFromInt(-1))
					subChild.Count = cat2.Count
					subChild.TrueCount = cat2.Count
					if cat2.ExcludeFromAnalysis || cat2.TopCategory == "Income" {
						subChild.Value = Zero
						subChild.Count = 0
					}
					topChild.Value = topChild.Value.Add(subChild.Value)
					topChild.Count = topChild.Count + subChild.Count
					topChild.TrueCount = topChild.TrueCount + subChild.TrueCount
					topChild.Children = append(topChild.Children, subChild)

				}
			}
			for i := range topChild.Children {
				if !topChild.Value.IsZero() {
					topChild.Children[i].Percent = topChild.Children[i].Value.Div(topChild.Value).Mul(decimal.NewFromInt(100)).StringFixed(1) + "%"
				} else {
					topChild.Children[i].Percent = "0%"
				}
			}
			treeData.Children = append(treeData.Children, topChild)
			treeData.Value = treeData.Value.Add(topChild.Value)
			treeData.Count = treeData.Count + topChild.Count
			treeData.TrueCount = treeData.TrueCount + topChild.TrueCount
		}
	}
	for i := range treeData.Children {
		if !treeData.Value.IsZero() {
			treeData.Children[i].Percent = treeData.Children[i].Value.Div(treeData.Value).Mul(decimal.NewFromInt(100)).StringFixed(1) + "%"
		} else {
			treeData.Children[i].Percent = "0%"
		}
	}

	if useTrueTotalFlag {
		treeData.TrueCount = trueTotal
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(treeData)

	return string(buf.String()), treeData.TrueCount
}
