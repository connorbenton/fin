package db

import (
	// "database/sql"

	"io/ioutil"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DBCon *sqlx.DB

func CreateDatabase() (*sqlx.DB, error) {

	var err error
	DBCon, err = sqlx.Open("sqlite3", "/usr/src/app/db/data-go.sqlite")
	if err != nil {
		return nil, err
	}
	DBCon.Exec("PRAGMA journal_mode=WAL;")

	raw, err := ioutil.ReadFile("/usr/src/app/backend_go/db/create.sql")
	query := string(raw)
	if err != nil {
		return nil, err
	}
	if _, err := DBCon.Exec(query); err != nil {
		return nil, err
	}

	return DBCon, nil
}
