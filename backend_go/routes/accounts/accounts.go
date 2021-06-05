package accounts

import (
	"encoding/json"
	"net/http"

	// "fmt"
	"fin-go/db"
	"fin-go/types"

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

		res.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(res).Encode(dbdata); err != nil {
			panic(err)
		}
	}
}

func UpsertIgnoreFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		p := types.Account{}

		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			panic(err)
		}

		txn := db.DBCon.MustBegin()
		astmt := types.PrepAccountUpsertIgnoreSt(txn)

		astmt.MustExec(p)
		errC := txn.Commit()
		if errC != nil {
			panic(errC)
		}

		res.WriteHeader(http.StatusOK)
	}
}

func UpsertNameFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		p := types.Account{}

		err := json.NewDecoder(req.Body).Decode(&p)
		if err != nil {
			panic(err)
		}

		txn := db.DBCon.MustBegin()
		astmt := types.PrepAccountUpsertNameSt(txn)

		astmt.MustExec(p)

		poststmt, err := txn.Preparex(`UPDATE transactions SET account_name = $1 WHERE account_id = $2`)
		if err != nil {
			panic(err)
		}

		poststmt.MustExec(p.Name, p.AccountID)

		errC := txn.Commit()
		if errC != nil {
			panic(errC)
		}

		res.WriteHeader(http.StatusOK)
	}
}
