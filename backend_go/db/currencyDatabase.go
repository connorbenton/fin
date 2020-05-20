package db

import (
	// "database/sql"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var CurrencyDBCon *sqlx.DB

type EcbFX struct {
	Currencies []struct {
		SeriesKey struct {
			Text   string `xml:",chardata"`
			Values []struct {
				ID    string `xml:"id,attr"`
				Value string `xml:"value,attr"`
			} `xml:"Value"`
		} `xml:"SeriesKey"`
		Rates []struct {
			Date struct {
				Value string `xml:"value,attr"`
			} `xml:"ObsDimension"`
			Rate struct {
				Value string `xml:"value,attr"`
			} `xml:"ObsValue"`
		} `xml:"Obs"`
	} `xml:"DataSet>Series"`
}

//CreateCurrencyDatabase starts currency database and loads it with initial info from the XML
func CreateCurrencyDatabase() (*sqlx.DB, error) {

	var err error
	CurrencyDBCon, err = sqlx.Open("sqlite3", "/usr/src/app/db/currencyData.sqlite")
	if err != nil {
		return nil, err
	}

	// raw2, err := ioutil.ReadFile("/usr/src/app/backend_go/db/currencyCreate.sql")
	raw2, err := ioutil.ReadFile("/usr/src/app/backend_go/db/currencyInitial.xml")
	if err != nil {
		return nil, err
	}

	insertXMLData(raw2, true)

	return DBCon, nil
}

func insertXMLData(data []byte, ignoreBool bool) {
	m := &EcbFX{}

	if err := xml.Unmarshal(data, &m); err != nil {
		log.Fatal(err)
	}

	txn := CurrencyDBCon.MustBegin()
	// txn, err := CurrencyDBCon.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//Create tables if does not exist
	// tblSt, err := txn.Prepare(`CREATE TABLE IF NOT EXISTS $1 ("fx_date" DATE PRIMARY KEY, "rate" TEXT)`)
	tblStr1 := `CREATE TABLE IF NOT EXISTS `
	tblStr2 := ` (fx_date DATE PRIMARY KEY, rate STRING);`
	// Don't need these, SQLite creates on its own
	// indexStr1 := `CREATE UNIQUE INDEX IF NOT EXISTS fx ON `
	// indexStr2 := ` ("fx_date");`
	//Upsert rows to table
	var upsertStr1 string
	var upsertStr2 string
	// var upsertStr3 string
	if ignoreBool {
		upsertStr1 = `INSERT OR IGNORE INTO `
		upsertStr2 = ` (fx_date, rate) VALUES($1,$2);`
	} else {
		// upsertStr1 = `INSERT INTO ARS(fx_date, rate) VALUES($1,$2) ON CONFLICT (fx_date) DO UDPATE SET counter = counter + 1`
		upsertStr1 = "INSERT INTO "
		upsertStr2 = "(fx_date, rate) VALUES($1, $2) ON CONFLICT (fx_date) DO UPDATE SET rate = excluded.rate"
		// upsertStr2 = `(fx_date, rate) VALUES($1,$2) ON CONFLICT (fx_date) DO NOTHING;`
	}
	// upsertSt, err := txn.Prepare(upsertStr)
	for _, currency := range m.Currencies {
		var curr string
		for _, key := range currency.SeriesKey.Values {
			if key.ID == "CURRENCY" {
				curr = key.Value
				txn.MustExec(tblStr1 + curr + tblStr2)
				// txn.Exec(tblStr1 + curr + tblStr2)
				// txn.MustExec(indexStr1 + curr + tblStr2)
				// if _, err = tblSt.Exec(key.Value); err != nil {
				// 	log.Fatal(err)
				// }
			}
		}
		for _, fx := range currency.Rates {
			if isNumDot(fx.Rate.Value) {
				// if ignoreBool {
				txStr := upsertStr1 + curr + upsertStr2
				// log.Println(txStr)
				txn.MustExec(txStr, fx.Date.Value, fx.Rate.Value)
				// if _, err = txn.Exec(upsertStr1+curr+upsertStr2, fx.Date.Value, fx.Rate.Value); err != nil {
				// 	panic(err)
				// }
				// } else {
				// txStr := []byte(upsertStr1 + curr + upsertStr2)
				// str1 := []byte("INSERT INTO ARS(fx_date, rate) VALUES($1, $2) ON CONFLICT (fx_date) DO UPDATE SET rate = excluded.rate")
				// log.Println(txStr)
				// log.Println(str1)
				// txn.MustExec(str1, fx.Date.Value, fx.Rate.Value)
				// 	str2 := upsertStr1 + curr + upsertStr2
				// 	log.Println(str1)
				// 	log.Println(str2)
				// 	txn.MustExec("INSERT INTO ARS(fx_date, rate) VALUES($1, $2) ON CONFLICT (fx_date) DO UPDATE SET rate = excluded.rate", fx.Date.Value, fx.Rate.Value)
				// }
				// txn.MustExec(upsertStr1+curr+upsertStr2, fx.Date.Value, fx.Rate.Value)
			}
			// if _, err = upsertSt.Exec(curr, fx.Date.Value, fx.Rate.Value); err != nil {
			// 	log.Fatal(err)
			// }
		}
	}

	txn.Commit()
	// getNewXML()
	// if err = txn.Commit(); err != nil {
	// 	log.Fatal(err)
	// }

}

func isNumDot(s string) bool {
	dotFound := false
	for _, v := range s {
		if v == '.' {
			if dotFound {
				return false
			}
			dotFound = true
		} else if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

// Fx type exported for currency rate lookups
type Fx struct {
	FxDate time.Time `db:"fx_date"`
	Rate   float32   `db:"rate"`
}

func GetNewXML() {

	//Get last updated date from fx data for search (using USD)
	rows, err := CurrencyDBCon.Queryx("SELECT * FROM `USD` ORDER BY fx_date DESC LIMIT 1")
	if err != nil {
		panic(err)
	}
	var fx Fx
	for rows.Next() {
		err = rows.StructScan(&fx)
		if err != nil {
			panic(err)
		}
	}
	utc := time.Now().UTC()
	daysDiff := utc.Sub(fx.FxDate).Hours() / 24
	// log.Println(utc)
	// log.Println(daysDiff)
	// hours, _, _ := utc.Clock()

	// if daysDiff > 1 && hours > 15 {
	if daysDiff > 1.66 {
		// log.Println("ready to fetch")
		url1 := "https://sdw-wsrest.ecb.europa.eu/service/data/EXR/D..EUR.SP00.A?updatedAfter="
		url2 := "T16%3A30%3A00%2B00%3A00&detail=dataonly"
		urlDate := fx.FxDate.Format("2006-01-02")
		resp, err := http.Get(url1 + urlDate + url2)
		if err != nil {
			xmlerr := fmt.Errorf("GET error: %v", err)
			log.Println(xmlerr.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			xmlerr := fmt.Errorf("Status error: %v", resp.StatusCode)
			log.Println(xmlerr.Error())
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			xmlerr := fmt.Errorf("Read body: %v", err)
			log.Println(xmlerr.Error())
		}

		insertXMLData(data, false)

	}

	//fx.FxDate is what we'll use to check if there's new info available

	// https://sdw-wsrest.ecb.europa.eu/service/data/EXR/D..EUR.SP00.A?updatedAfter=2020-05-15T14%3A15%3A00%2B01%3A00&detail=dataonly
}
