package db

import (
	// "database/sql"

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
	// var wg sync.WaitGroup
	// wg.Add(2)
	// go func() {
	// defer wg.Done()
	var err error
	CurrencyDBCon, err = sqlx.Open("sqlite3", "/usr/src/app/db/currencyData.sqlite")
	if err != nil {
		log.Fatalln(err)
	}
	// }()

	// log.Println("extracting currencyInitial")
	// raw2, err := ioutil.ReadFile("/usr/src/app/backend_go/db/currencyCreate.sql")
	// raw2, err := ioutil.ReadFile("/usr/src/app/backend_go/db/currencyInitial.xml")
	// raw2, err := ioutil.ReadFile("/usr/src/app/backend_go/db/currencyInitial.xml.gz")
	// if err != nil {
	// 	return nil, err
	// }

	// go func() {
	// defer wg.Done()
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

	// s, err := ioutil.ReadAll(fz)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Println("uploading initial currency XML")

	// xmlString := string(s)

	// data := types.EcbFX{}

	// if err := xml.NewDecoder(fz).Decode(&data); err != nil {
	// 	return nil, err
	// }

	// insertXMLData(xmlString, true)
	// insertXMLData(data, true)
	// insertXMLData(fz, true)
	insertXMLData(body, true)
	// }()
	// wg.Wait()
	// GetNewXML()

	// raw2, err := os.Open("/usr/src/app/backend_go/db/createCurrency.sql.gz")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer raw2.Close()
	// log.Println("Open sql done:", time.Since(start))

	// fz, err := gzip.NewReader(raw2)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer fz.Close()
	// log.Println("Unzip sql done:", time.Since(start))

	// raw, err := ioutil.ReadAll(fz)
	// query := string(raw)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Println("Read sql done:", time.Since(start))
	// CurrencyDBCon.MustExec(query)

	log.Println("Currency DB setup done in:", time.Since(start))

	return DBCon, nil
}

// func insertXMLData(data string, isInitialLoad bool) {
// func insertXMLData(m types.EcbFX, isInitialLoad bool) {
// func insertXMLData(stream *gzip.Reader, isInitialLoad bool) {
func insertXMLData(bytesXML []byte, isInitialLoad bool) {

	// start := time.Now()

	// m := &types.EcbFX{}

	// if err := xml.Unmarshal([]byte(data), &m); err != nil {
	// 	panic(err)
	// }
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

	// d := xml.NewDecoder(stream)
	// var wg sync.WaitGroup

	stream := bytes.NewReader(bytesXML)
	// if !isInitialLoad {
	// log.Println(string(bytesXML))
	// }
	br := bufio.NewReaderSize(stream, 65536)
	// str, _ := br.ReadString('Q')
	// log.Println(str)

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
			// wg.Add(1)
			// go func(fx xmlparser.XMLElement) {
			// defer wg.Done()
			if isNumDot(fx.Childs["generic:ObsValue"][0].Attrs["value"]) {
				fxSt.MustExec(fx.Childs["generic:ObsDimension"][0].Attrs["value"], fx.Childs["generic:ObsValue"][0].Attrs["value"])
				// if !isInitialLoad {
				// 	log.Println(fx.Childs["generic:ObsDimension"][0].Attrs["value"], fx.Childs["generic:ObsValue"][0].Attrs["value"])
				// }
			}
			// }(fx)
		}
	}
	// wg.Wait()

	// for {

	// 	tok, err := d.Token()
	// 	if tok == nil || err == io.EOF {
	// 		// EOF means we're done.
	// 		break
	// 	} else if err != nil {
	// 		log.Fatalf("Error decoding token: %s", err)
	// 	}

	// 	switch ty := tok.(type) {
	// 	case xml.StartElement:
	// 		if ty.Name.Local == "Series" {

	// 			// go func(d *xml.Decoder, ty xml.StartElement) {
	// 			// start2 := time.Now()
	// 			// log.Println("Found series", time.Since(start))
	// 			var currency types.EcbFXCurrency
	// 			if err = d.DecodeElement(&currency, &ty); err != nil {
	// 				log.Fatalf("Error decoding item: %s", err)
	// 			}
	// 			// wg.Add(1)
	// 			// go func(currency types.EcbFXCurrency) {
	// 			// defer wg.Done()
	// 			// log.Println("series decoded in: ", time.Since(start2))
	// 			var curr string
	// 			for _, key := range currency.SeriesKey.Values {
	// 				if key.ID == "CURRENCY" {
	// 					curr = key.Value
	// 					txn.MustExec(tblStr1 + curr + tblStr2)
	// 				}
	// 			}
	// 			txStr := upsertStr1 + curr + upsertStr2
	// 			fxSt := types.PrepFXTableSt(txn, txStr)

	// 			for _, fx := range currency.Rates {

	// 				if isNumDot(fx.Rate.Value) {

	// 					fxSt.MustExec(fx.Date.Value, fx.Rate.Value)

	// 				}

	// 			}

	// 			// log.Println("series done in: ", time.Since(start2))
	// 			// }(currency)
	// 			// }(d, ty)

	// 		}

	// 	}
	// }

	// wg.Wait()

	err2 := txn.Commit()
	if err2 != nil {
		panic(err2)
	}
	// log.Println("insertXML done:", time.Since(start))
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

	// log.Println(utc)
	// log.Println(daysDiff)
	// hours, _, _ := utc.Clock()

	// if daysDiff > 1 && hours > 15 {
	if daysDiff > 1.66 {

		start := time.Now()
		log.Println("Pulling new data from ECB API")
		// log.Println("ready to fetch")
		url1 := "https://sdw-wsrest.ecb.europa.eu/service/data/EXR/D..EUR.SP00.A?updatedAfter="
		url2 := "T16%3A30%3A00%2B00%3A00&detail=dataonly"
		urlDate := fx.FxDate.Format("2006-01-02")

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

		// fz, err := gzip.NewReader(resp.Body)

		// if err != nil {
		// 	log.Println("Response gzip error: ", err)
		// }
		// defer fz.Close()

		// data := types.EcbFX{}

		// if err := xml.NewDecoder(resp.Body).Decode(&data); err != nil {
		// 	panic(err)
		// }

		// data, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	xmlerr := fmt.Errorf("Read body: %v", err)
		// 	log.Println(xmlerr.Error())
		// }

		// xmlString := string(data)

		// insertXMLData(xmlString, false)
		// insertXMLData(data, false)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Read error:", err)
		}
		insertXMLData(body, false)
		log.Println("Data pulled and inserted from ECB API in: ", time.Since(start))

	}

	//fx.FxDate is what we'll use to check if there's new info available

	// https://sdw-wsrest.ecb.europa.eu/service/data/EXR/D..EUR.SP00.A?updatedAfter=2020-05-15T14%3A15%3A00%2B01%3A00&detail=dataonly
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
		// firstDate := tdate.AddDate(0, 0, -10).Format("2006-01-02")
		// lastDate := tdate.AddDate(0, 0, 10).Format("2006-01-02")
		firstDate := tdate.AddDate(0, 0, -10)
		lastDate := tdate.AddDate(0, 0, 10)
		// date := date.Format("2006-01-02")
		fx := types.Fx{}
		if CC == "EUR" {
			// finding the nearest 'EUR' rate by doing a search with USD and swapping in 1.0 for rate
			CC = "USD"
			query := fmt.Sprintf(`SELECT * FROM %q WHERE fx_date BETWEEN 
					%q AND %q ORDER BY abs(%q - fx_date) LIMIT 1`, CC, firstDate, lastDate, tdate.Format("2006-01-02"))
			err := CurrencyDBCon.Get(&fx, query)
			CC = "EUR"
			if err != nil && err != sql.ErrNoRows {
				panic(err)
			}
			fx.Rate = decimal.NewFromInt(1)
		} else {
			// otherwise finding the nearest rate in +/- 10 days
			query := fmt.Sprintf(`SELECT * FROM %q WHERE fx_date BETWEEN 
					%q AND %q ORDER BY abs(%q - fx_date) LIMIT 1`, CC, firstDate, lastDate, tdate.Format("2006-01-02"))
			err := CurrencyDBCon.Get(&fx, query)
			if err != nil && err != sql.ErrNoRows {
				panic(err)
			}
		}
		// Setting to zero and skipping if rate not found within +/- 10 days
		if (types.Fx{}) == fx {
			log.Println("Currency rate not found for " + CC + " within +/- 10 days of " + tdate.Format("2006-01-02"))
			NormalizedAmount = decimal.Zero
		} else {
			if baseCurrency == "EUR" {
				// Using no extra rate if base currency is EUR
				NormalizedAmount = amt.Div(fx.Rate)
			} else {
				// Finding second rate for base currency other than EUR
				bfx := types.Fx{}
				query := fmt.Sprintf(`SELECT * FROM %q WHERE fx_date = %q`, baseCurrency, fx.FxDate.Format("2006-01-02"))
				err := CurrencyDBCon.Get(&bfx, query)
				if err != nil && err != sql.ErrNoRows {
					panic(err)
				}
				if (types.Fx{}) == bfx {
					// Setting to zero and skipping if base rate not found
					log.Println("Currency rate for base currency " + baseCurrency + " not found on date " + fx.FxDate.Format("2006-01-02"))
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
