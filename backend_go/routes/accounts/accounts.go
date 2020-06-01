package accounts

import (
	"encoding/json"
	"net/http"

	// "fmt"
	"fintrack-go/db"
	"fintrack-go/types"

	_ "github.com/jmoiron/sqlx"
)

func SelectAll() []types.Account {
	dbdata := []types.Account{}
	err := db.DBCon.Select(&dbdata, "SELECT * FROM `accounts`")
	if err != nil {
		panic(err)
	}
	return dbdata
}

func GetFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		dbdata := SelectAll()
		// dbdata := []types.Account{}
		// err := db.DBCon.Select(&dbdata, "SELECT * FROM `accounts`")
		// if err != nil {
		// log.Fatal(err)
		// }

		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(dbdata); err != nil {
			panic(err)
		}
	}
}