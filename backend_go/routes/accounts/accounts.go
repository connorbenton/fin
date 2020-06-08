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

func UpsertFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		p := types.Account{}

		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			panic(err)
		}

		txn := db.DBCon.MustBegin()
		astmt := types.PrepAccountUpsertSt(txn)

		astmt.MustExec(p)
		errC := txn.Commit()
		if errC != nil {
			panic(errC)
		}

		res.WriteHeader(http.StatusOK)
	}
}
