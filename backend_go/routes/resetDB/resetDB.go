package resetDB

import (
	"io/ioutil"
	"log"
	"net/http"

	"fintrack-go/db"

	_ "github.com/jmoiron/sqlx"
)

func ForceResetDBFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		raw, err := ioutil.ReadFile("/usr/src/app/backend_go/db/drop.sql")
		if err != nil {
			panic(err)
		}
		query := string(raw)
		if _, err := db.DBCon.Exec(query); err != nil {
			panic(err)
		}

		raw2, err2 := ioutil.ReadFile("/usr/src/app/backend_go/db/create.sql")
		query2 := string(raw2)
		if err2 != nil {
			panic(err)
		}
		if _, err := db.DBCon.Exec(query2); err != nil {
			panic(err)
		}

		log.Println("Reset (excluding itemTokens) successful")
		res.WriteHeader(http.StatusOK)

	}
}

func ForceResetDBFullFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		raw, err := ioutil.ReadFile("/usr/src/app/backend_go/db/fulldrop.sql")
		if err != nil {
			panic(err)
		}
		query := string(raw)
		if _, err := db.DBCon.Exec(query); err != nil {
			panic(err)
		}

		raw2, err2 := ioutil.ReadFile("/usr/src/app/backend_go/db/create.sql")
		query2 := string(raw2)
		if err2 != nil {
			panic(err2)
		}
		if _, err := db.DBCon.Exec(query2); err != nil {
			panic(err)
		}

		log.Println("Reset (Full reset including itemTokens) successful")
		res.WriteHeader(http.StatusOK)
	}
}
