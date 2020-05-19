package analysisTrees

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	// "fmt"
	"fintrack-go/db"

	_ "github.com/jmoiron/sqlx"
)

type Account struct {
	Id   		int			`json:"id"`
	Name		string		`json:"name" db:"name"`		
	Institution string		`json:"institution" db:"institution"`
	Account_id 	string		`json:"account_id" db:"account_id"`
	Item_id 	string		`json:"item_id" db:"item_id"`
	Type 		string 		`json:"type" db:"type"`
	Subtype 	string		`json:"subtype" db:"subtype"`
	Balance 	float32		`json:"balance" db:"balance"`
	Limit 		float32 	`json:"limit" db:"limit"`
	Available 	float32 	`json:"available" db:"available"`
	Currency 	string 		`json:"currency" db:"currency"`
	Provider	string 		`json:"provider" db:"provider"`
	RunningTotal	float32	`json:"running_total" db:"running_total"`
	CreatedAt time.Time 	`json:"created_at" db:"created_at"`
	UpdatedAt time.Time 	`json:"updated_at" db:"updated_at"`
}


func GetFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req * http.Request) {

	dbdata := []Account{}
	err := db.DBCon.Select(&dbdata, "SELECT * FROM `accounts`")
	if err != nil {
		log.Fatal(err)
		}

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(dbdata); err != nil {
		panic(err)
		}
	}
}