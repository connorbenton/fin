package accounts

import (
	"encoding/json"
	"log"
	"net/http"

	// "fmt"
	"fintrack-go/db"
	"fintrack-go/types"

	_ "github.com/jmoiron/sqlx"
)

func GetFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		dbdata := []types.Account{}
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
