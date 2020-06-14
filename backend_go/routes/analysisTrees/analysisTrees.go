package analysisTrees

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"fin-go/db"
	"fin-go/types"

	_ "github.com/jmoiron/sqlx"
	"github.com/rickb777/date"

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
		stDt, _ := date.ParseISO(customRange.Start)
		endDt, _ := date.ParseISO(customRange.End)
		dbEnd = endDt.AddDate(0, 0, 1).String()
		err = db.DBCon.Select(&rangedata, "SELECT * FROM `transactions` WHERE DATE between '"+st+"' and '"+dbEnd+"'")
		// log.Println("number in ", name, len(rangedata))
		if err != nil {
			log.Println(end)
			panic(err)
		}

		diff := int64(endDt.Sub(stDt) + 1)

		tree := SetupTree(dbcatsBase, rangedata, "custom", st, end, diff)
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
	log.Println("today: ", date.Today().String(), date.Today().FormatISO(4), date.Today().String())

	start := time.Now()
	dbcatsBase := []types.Category{}

	err := db.DBCon.Select(&dbcatsBase, "SELECT * FROM `categories`")
	if err != nil {
		panic(err)
	}

	txn := db.DBCon.MustBegin()

	tstmt := types.PrepTreeSt(txn)

	var wg sync.WaitGroup
	for _, name := range types.TreeRanges {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			rangedata := []types.Transaction{}
			today := date.Today()
			var stDt, endDt date.Date
			var st, end, dbEnd string
			var err error
			switch name {
			case "last30":
				stDt = today.AddDate(0, 0, -29)
				endDt = today
			case "thisMonth":
				stDt = date.New(today.Year(), today.Month(), 1)
				endDt = today
			case "lastMonth":
				stPre := today.AddDate(0, -1, 0)
				stDt = date.New(stPre.Year(), stPre.Month(), 1)
				endDt = date.New(stPre.Year(), stPre.Month(), 1).AddDate(0, 1, -1)
			case "last6Months":
				stDt = today.AddDate(0, -6, 0)
				endDt = today
			case "thisYear":
				stDt = date.New(today.Year(), 1, 1)
				endDt = today
			case "lastYear":
				stPre := today.AddDate(-1, 0, 0)
				stDt = date.New(stPre.Year(), 1, 1)
				endDt = date.New(stPre.Year(), 1, 1).AddDate(1, 0, -1)
			case "fromBeginning":
				var dbst sql.NullString
				err2 := db.DBCon.Get(&dbst, "SELECT MIN(date) FROM `transactions`")
				if err2 != nil {
					panic(err2)
				}
				if !dbst.Valid {
					return
				}
				st = dbst.String
				stDt, _ = date.ParseISO(st)
				endDt = today
			case "custom":
				stDt = today.AddDate(0, 0, -29)
				endDt = today
			}

			st = stDt.String()
			end = endDt.String()
			dbEnd = endDt.AddDate(0, 0, 1).String()

			err = db.DBCon.Select(&rangedata, "SELECT * FROM `transactions` WHERE DATE between '"+st+"' and '"+dbEnd+"'")
			// log.Println("number in ", name, len(rangedata))
			if err != nil {
				log.Println(end)
				panic(err)
			}

			diff := int64(endDt.Sub(stDt) + 1)

			tree := SetupTree(dbcatsBase, rangedata, name, st, end, diff)

			tstmt.MustExec(tree)

		}(name)
	}

	wg.Wait()
	errC := txn.Commit()
	if errC != nil {
		panic(errC)
	}
	log.Println("Regenerate analysis trees done in:", time.Since(start))
}

func SetupTree(dbcatsBase []types.Category, rangedata []types.Transaction, name, st, end string, diff int64) types.Tree {
	dbcats := append(dbcatsBase[:0:0], dbcatsBase...)
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
	var totalcount int64
	tree.Data, totalcount = GenerateDataTree(dbcats, diff, 0, false)
	tree.DataNoInvest, _ = GenerateDataTree(dbcatsNoInvest, diff, totalcount, true)
	return tree
}

func GenerateDataTree(dbcats []types.Category, diff, trueTotal int64, useTrueTotalFlag bool) (string, int64) {

	var Zero = decimal.New(0, 1)
	treeData := types.TreeData{}

	treeData.Name = "Transactions by Category"
	treeData.Children = []types.ChildTop{}
	treeData.Value = Zero
	treeData.Count = 0
	treeData.TrueCount = 0
	treeData.IncomeTotal = Zero

	for _, cat := range dbcats {
		if cat.SubCategory == cat.TopCategory {
			topChild := types.ChildTop{}
			topChild.Name = cat.TopCategory
			topChild.Children = []types.ChildSub{}
			topChild.Value = Zero
			topChild.Count = 0
			topChild.TrueCount = 0
			topChild.DbID = 0

			for _, cat2 := range dbcats {
				if cat2.TopCategory == cat.TopCategory {
					subChild := types.ChildSub{}
					subChild.Name = ""
					subChild.Value = Zero
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
					subChild.Per30 = subChild.Value.Div(decimal.NewFromInt(diff)).Mul(decimal.NewFromInt(30))
					subChild.Count = int64(cat2.Count)
					subChild.TrueCount = int64(cat2.Count)
					if cat2.ExcludeFromAnalysis {
						subChild.Value = Zero
						subChild.Count = 0
					}
					if cat2.TopCategory == "Income" {
						treeData.IncomeTotal = treeData.IncomeTotal.Add(subChild.Value.Abs())
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
				if topChild.Children[i].Value.IsNegative() {
					topChild.Value = topChild.Value.Sub(topChild.Children[i].Value)
					topChild.Count = topChild.Count - topChild.Children[i].Count
					topChild.Children[i].Value = Zero
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
			treeData.Children[i].Per30 = treeData.Children[i].Value.Div(decimal.NewFromInt(diff)).Mul(decimal.NewFromInt(30))
		} else {
			treeData.Children[i].Percent = "0%"
		}
		if treeData.Children[i].Value.IsNegative() {
			treeData.Value = treeData.Value.Sub(treeData.Children[i].Value)
			treeData.Count = treeData.Count - treeData.Children[i].Count
			treeData.Children[i].Value = decimal.New(0, 1)
		}
	}
	treeData.Per30 = treeData.Value.Div(decimal.NewFromInt(diff)).Mul(decimal.NewFromInt(30))

	if useTrueTotalFlag {
		treeData.TrueCount = trueTotal
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(treeData)

	return string(buf.String()), treeData.TrueCount
}
