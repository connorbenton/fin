package plaid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	// "fmt"
	"fintrack-go/db"

	_ "github.com/jmoiron/sqlx"
	_ "github.com/plaid/plaid-go/plaid"
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

func CreateFromPublicTokenFunction() func(http.ResponseWriter, *http.Request) {
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

func GeneratePublicTokenFunction() func(http.ResponseWriter, *http.Request) {
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

func RefreshConnectionsFunction() func(http.ResponseWriter, *http.Request) {
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

func FetchTransactionsForItemToken(ItemID string) {

	url := "https://www.saltedge.com/api/v5/transactions?connection_id" + ItemID
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
