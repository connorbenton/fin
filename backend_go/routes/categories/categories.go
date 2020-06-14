package categories

import (
	"encoding/json"
	"net/http"

	"fin-go/db"
	"fin-go/types"

	_ "github.com/jmoiron/sqlx"
)

func SelectAll() []types.Category {
	dbdata := []types.Category{}
	err := db.DBCon.Select(&dbdata, "SELECT * FROM `categories`")
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
