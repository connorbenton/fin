package categories

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	// "fmt"
	"fintrack-go/db"

	_ "github.com/jmoiron/sqlx"
)

type Category struct {
	Id   int				`json:"id"`
	TopCategory	string		`json:"top_category" db:"top_category"`		
	SubCategory string		`json:"sub_category" db:"sub_category"`
	ExcludeFromAnalysis int	`json:"exclude_from_analysis" db:"exclude_from_analysis"`
	CreatedAt time.Time 	`json:"created_at" db:"created_at"`
	UpdatedAt time.Time 	`json:"updated_at" db:"updated_at"`
}


func GetFunction() func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req * http.Request) {

	dbdata := []Category{}
	err := db.DBCon.Select(&dbdata, "SELECT * FROM `categories`")
	if err != nil {
		log.Fatal(err)
		}

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(dbdata); err != nil {
		panic(err)
		}
	}
}