package db

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"io/ioutil"

	"log"
	"net/http"
	"os"
	"strings"

	"time"

	"fintrack-go/types"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rickb777/date"
	"github.com/shopspring/decimal"
	xmlparser "github.com/tamerh/xml-stream-parser"
)

var CurrencyDBCon *sqlx.DB

// CreateCurrencyDatabase starts currency database and loads it with initial info from the XML
func CreateCurrencyDatabase() (*sqlx.DB, error) {

	start := time.Now()

	var err error
	CurrencyDBCon, err = sqlx.Open("sqlite3", "/usr/src/app/db/currencyData.sqlite")
	if err != nil {
		log.Fatalln(err)
	}

	raw2, err := os.Open("/usr/src/app/backend_go/db/currencyInitial.xml.gz")
	if err != nil {
		log.Fatalln(err)
	}
	defer raw2.Close()

	fz, err := gzip.NewReader(raw2)
	if err != nil {
		log.Fatalln(err)
	}
	defer fz.Close()

	body, err := ioutil.ReadAll(fz)
	if err != nil {
		log.Fatalln(err)
	}

	insertXMLData(body, true)

	log.Println("Currency DB setup done in:", time.Since(start))

	return DBCon, nil
}

func insertXMLData(bytesXML []byte, isInitialLoad bool) {

	txn := CurrencyDBCon.MustBegin()

	tblStr1 := `CREATE TABLE IF NOT EXISTS `
	tblStr2 := ` (fx_date DATE PRIMARY KEY, rate STRING);`
	var upsertStr1 string
	var upsertStr2 string
	if isInitialLoad {
		upsertStr1 = `INSERT OR IGNORE INTO `
		upsertStr2 = ` (fx_date, rate) VALUES($1,$2);`
	} else {
		upsertStr1 = "INSERT INTO "
		upsertStr2 = "(fx_date, rate) VALUES($1, $2) ON CONFLICT (fx_date) DO UPDATE SET rate = excluded.rate"
	}

	stream := bytes.NewReader(bytesXML)

	br := bufio.NewReaderSize(stream, 65536)

	parser := xmlparser.NewXMLParser(br, "generic:Series")

	for xml := range parser.Stream() {
		var curr string
		for _, fxKey := range xml.Childs["generic:SeriesKey"][0].Childs["generic:Value"] {
			if fxKey.Attrs["id"] == "CURRENCY" {
				curr = fxKey.Attrs["value"]
				txn.MustExec(tblStr1 + curr + tblStr2)
				break
			}
		}
		txStr := upsertStr1 + curr + upsertStr2
		fxSt := types.PrepFXTableSt(txn, txStr)

		for _, fx := range xml.Childs["generic:Obs"] {

			if isNumDot(fx.Childs["generic:ObsValue"][0].Attrs["value"]) {
				fxSt.MustExec(fx.Childs["generic:ObsDimension"][0].Attrs["value"], fx.Childs["generic:ObsValue"][0].Attrs["value"])

			}
		}
	}

	err2 := txn.Commit()
	if err2 != nil {
		panic(err2)
	}
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

func GetNewXML() {

	//Get last updated date from fx data for search (using USD)
	rows, err := CurrencyDBCon.Queryx("SELECT * FROM `USD` ORDER BY fx_date DESC LIMIT 1")
	if err != nil {
		panic(err)
	}
	var fx types.Fx
	for rows.Next() {
		err = rows.StructScan(&fx)
		if err != nil {
			panic(err)
		}
	}
	utc := time.Now().UTC()
	daysDiff := utc.Sub(fx.FxDate).Hours() / 24

	if daysDiff > 1.66 {

		start := time.Now()
		log.Println("Pulling new data from ECB API")
		// log.Println("ready to fetch")
		url1 := "https://sdw-wsrest.ecb.europa.eu/service/data/EXR/D..EUR.SP00.A?updatedAfter="
		url2 := "T16%3A30%3A00%2B00%3A00&detail=dataonly"
		urlDate := fx.FxDate.String()

		client := new(http.Client)

		request, err := http.NewRequest("GET", url1+urlDate+url2, nil)
		// request.Header.Add("Accept-Encoding", "gzip")

		// resp, err := http.Get(url1 + urlDate + url2)
		resp, err := client.Do(request)
		if err != nil {
			xmlerr := fmt.Errorf("GET error: %v", err)
			log.Println(xmlerr.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			xmlerr := fmt.Errorf("Status error: %v", resp.StatusCode)
			if resp.StatusCode == 404 {
				return
			}
			if resp.StatusCode == 500 {
				log.Println("500 Server error - ECB SDMX")
				return
			}
			if resp.StatusCode == 503 {
				log.Println("503 Server temporarily unavailable - ECB SDMX")
				return
			}
			log.Println(xmlerr.Error())
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Read error:", err)
		}
		insertXMLData(body, false)
		log.Println("Data pulled and inserted from ECB API in: ", time.Since(start))

	}

}

//GetNormalizedAmount finds and sets the correct normalized amount for transactions
func GetNormalizedAmount(code string, baseCurrency string, dt string, amt decimal.Decimal) decimal.Decimal {

	CC := strings.ToUpper(code)
	var err error
	var NormalizedAmount decimal.Decimal
	tdate, err := date.ParseISO(dt)
	if err != nil {
		panic(err)
	}

	var id int
	err = CurrencyDBCon.Get(&id, "SELECT count(*) FROM sqlite_master WHERE type='table' AND name='"+CC+"'")
	if err != nil {
		panic(err)
	}
	// want to continue if currency is EUR since there is no table in DB for it
	if CC == "EUR" {
		id = 1
	}
	// Setting to zero and skipping if table not found for currency (i.e. BTC)
	if id == 0 {
		log.Println("Currency rate table not found for " + CC)
		NormalizedAmount = decimal.Zero
	} else {
		firstDate := tdate.AddDate(0, 0, -10)
		lastDate := tdate.AddDate(0, 0, 10)
		fx := types.Fx{}
		if CC == "EUR" {
			// finding the nearest 'EUR' rate by doing a search with USD and swapping in 1.0 for rate
			CC = "USD"
			query := fmt.Sprintf(`SELECT * FROM %q WHERE fx_date BETWEEN 
					%q AND %q ORDER BY abs(%q - fx_date) LIMIT 1`, CC, firstDate, lastDate, tdate.String())
			err := CurrencyDBCon.Get(&fx, query)
			CC = "EUR"
			if err != nil && err != sql.ErrNoRows {
				panic(err)
			}
			fx.Rate = decimal.NewFromInt(1)
		} else {
			// otherwise finding the nearest rate in +/- 10 days
			query := fmt.Sprintf(`SELECT * FROM %q WHERE fx_date BETWEEN 
					%q AND %q ORDER BY abs(%q - fx_date) LIMIT 1`, CC, firstDate, lastDate, tdate.String())
			err := CurrencyDBCon.Get(&fx, query)
			if err != nil && err != sql.ErrNoRows {
				panic(err)
			}
		}
		// Setting to zero and skipping if rate not found within +/- 10 days
		if (types.Fx{}) == fx {
			log.Println("Currency rate not found for " + CC + " within +/- 10 days of " + tdate.String())
			NormalizedAmount = decimal.Zero
		} else {
			if baseCurrency == "EUR" {
				// Using no extra rate if base currency is EUR
				NormalizedAmount = amt.Div(fx.Rate)
			} else {
				// Finding second rate for base currency other than EUR
				bfx := types.Fx{}
				query := fmt.Sprintf(`SELECT * FROM %q WHERE fx_date = %q`, baseCurrency, fx.FxDate.String())
				err := CurrencyDBCon.Get(&bfx, query)
				if err != nil && err != sql.ErrNoRows {
					panic(err)
				}
				if (types.Fx{}) == bfx {
					// Setting to zero and skipping if base rate not found
					log.Println("Currency rate for base currency " + baseCurrency + " not found on date " + fx.FxDate.String())
					NormalizedAmount = decimal.Zero
				} else {
					// Decimal math to find normalized amount with rate and base rate
					NormalizedAmount = amt.Div(fx.Rate).Mul(bfx.Rate)
				}
			}
		}
	}
	return NormalizedAmount
}
