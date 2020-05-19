package itemTokens

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	// "fmt"
	"fintrack-go/db"

	_ "github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

type ItemToken struct {
	Id   			int			`json:"id"`
	Institution		string		`json:"institution" db:"institution"`		
	Access_token	string		`json:"access_token" db:"access_token"`
	Item_id		 	string		`json:"item_id" db:"item_id"`
	Provider		string		`json:"provider" db:"provider"`
	Interactive		bool		`json:"interactive" db:"interactive"`
	NeedsReLogin	bool		`json:"needs_re_login" db:"needs_re_login"`
	LastRefresh 				time.Time 	`json:"last_refresh" db:"last_refresh"`
	NextRefreshPossible 		time.Time 	`json:"next_refresh_possible" db:"next_refresh_possible"`
	LastDownloadedTransactions 	time.Time 	`json:"last_downloaded_transactions" db:"last_downloaded_transactions"`
	CreatedAt time.Time 	`json:"created_at" db:"created_at"`
	UpdatedAt time.Time 	`json:"updated_at" db:"updated_at"`
}

type CurrencyRate struct {
	Id			int				`json:"id"`
	Date		time.Time		`json:"date" db:"date"`	
	Rates		types.JSONText  `json:"rates" db:"rates"`	


	CreatedAt 	time.Time 		`json:"created_at" db:"created_at"`
	UpdatedAt 	time.Time 		`json:"updated_at" db:"updated_at"`
}

func GetFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req * http.Request) {

	dbdata := []ItemToken{}
	err := db.DBCon.Select(&dbdata, "SELECT * FROM `item_tokens`")
	if err != nil {
		log.Fatal(err)
		}

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(dbdata); err != nil {
		panic(err)
		}
	}
}

func FetchTransactionsFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req * http.Request) {

	itemTokens := []ItemToken{}
	err := db.DBCon.Select(&itemTokens, "SELECT * FROM `item_tokens`")
	if err != nil {
		log.Fatal(err)
		}

	currencyRates := []CurrencyRate{}
	err2 := db.DBCon.Select(&currencyRates, "SELECT * FROM `currency_rates`")
	if err2 != nil {
		log.Fatal(err)
		}
	// rates, _ := string([]byte{json.Marshal(currencyRates[0].DataJSON)})
	// log.Println(rates)









	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(itemTokens); err != nil {
		panic(err)
		}
	}
}