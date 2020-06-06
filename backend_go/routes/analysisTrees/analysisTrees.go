package analysisTrees

import (
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
)

func GetFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		dbdata := []types.Tree{}
		err := db.DBCon.Select(&dbdata, "SELECT * FROM `trees`")
		if err != nil {
			panic(err)
		}

		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(dbdata); err != nil {
			panic(err)
		}
	}
}

func ReAnalyze() {

	start := time.Now()
	// dbdata := []types.Transaction{}
	// err := db.DBCon.Select(&dbdata, "SELECT * FROM `transactions`")
	// if err != nil {
	// 	panic(err)
	// }

	

	log.Println("First select of tx set:", time.Since(start))

	var wg sync.WaitGroup
	for _, name := range types.TreeRanges {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			rangedata := []types.Transaction{}
			tree := []types.Tree{}
			today := date.Today()
			var st, end, dbEnd string
			var err error
			switch name {
			case "last30":
				st = today.AddDate(0, 0, -30).Format("2006-01-02")
				end = today.Format("2006-01-02")
				dbEnd = today.AddDate(0, 0, 1).Format("2006-01-02")
				err = db.DBCon.Select(&rangedata, "SELECT * FROM `transactions` WHERE DATE between '"+st+"' and '"+dbEnd+"'")
				log.Println("number in ", name, len(rangedata))
			case "thisMonth":
				st = date.New(today.Year(), today.Month(), 1).Format("2006-01-02")
				end = today.Format("2006-01-02")
				dbEnd = today.AddDate(0, 0, 1).Format("2006-01-02")
				err = db.DBCon.Select(&rangedata, "SELECT * FROM `transactions` WHERE DATE between '"+st+"' and '"+dbEnd+"'")
				log.Println("number in ", name, len(rangedata))
			case "lastMonth":
				stdate := today.AddDate(0, -1, 0)
				st = date.New(stdate.Year(), stdate.Month(), 1).Format("2006-01-02")
				end = stdate.AddDate(0, 1, -1).Format("2006-01-02")
				dbEnd = today.AddDate(0, 0, 1).Format("2006-01-02")
				err = db.DBCon.Select(&rangedata, "SELECT * FROM `transactions` WHERE DATE between '"+st+"' and '"+dbEnd+"'")
				log.Println("number in ", name, len(rangedata))
			case "last6Months":
				st = today.AddDate(0, -6, 0).Format("2006-01-02")
				end = today.Format("2006-01-02")
				dbEnd = today.AddDate(0, 0, 1).Format("2006-01-02")
				err = db.DBCon.Select(&rangedata, "SELECT * FROM `transactions` WHERE DATE between '"+st+"' and '"+dbEnd+"'")
				log.Println("number in ", name, len(rangedata))
			case "thisYear":
				st = date.New(today.Year(), 1, 1).Format("2006-01-02")
				end = today.Format("2006-01-02")
				dbEnd = today.AddDate(0, 0, 1).Format("2006-01-02")
				err = db.DBCon.Select(&rangedata, "SELECT * FROM `transactions` WHERE DATE between '"+st+"' and '"+dbEnd+"'")
				log.Println("number in ", name, len(rangedata))
			case "lastYear":
				stdate := today.AddDate(-1, 0, 0)
				st = date.New(stdate.Year(), 1, 1).Format("2006-01-02")
				end = stdate.AddDate(1, 0, -1).Format("2006-01-02")
				dbEnd = today.AddDate(0, 0, 1).Format("2006-01-02")
				err = db.DBCon.Select(&rangedata, "SELECT * FROM `transactions` WHERE DATE between '"+st+"' and '"+dbEnd+"'")
				log.Println("number in ", name, len(rangedata))
			case "fromBeginning":
				// rangedata = dbdata
				err = db.DBCon.Select(&rangedata, "SELECT * FROM `transactions`")
				log.Println("number in ", name, len(rangedata))
			case "custom":
				st = today.AddDate(0, 0, -30).Format("2006-01-02")
				end = today.Format("2006-01-02")
				dbEnd = today.AddDate(0, 0, 1).Format("2006-01-02")
				err = db.DBCon.Select(&rangedata, "SELECT * FROM `transactions` WHERE DATE between '"+st+"' and '"+dbEnd+"'")
				log.Println("number in ", name, len(rangedata))
			}
			if err != nil {
				log.Println(end)
				panic(err)
			}
		}(name)
	}

	wg.Wait()
	log.Println("Reanalyze tx ranges done in:", time.Since(start))
}
