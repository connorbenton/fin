package main

import (
	"log"
	"net/http"

	"fin-go/app"
	"fin-go/db"

	"github.com/gorilla/mux"
)

func main() {
	// Init DB
	_, err := db.CreateDatabase()
	if err != nil {
		log.Fatal("main: cannot initialize DB: %s", err.Error())
	}

	// Init Currency DB
	_, err2 := db.CreateCurrencyDatabase()
	if err2 != nil {
		log.Fatal("main: cannot initialize Currency DB: %s", err.Error())
	}

	db.GetNewXML()

	app := &app.App{
		Router: mux.NewRouter().StrictSlash(true),
	}

	app.SetupRouter()

	log.Println("Starting HTTP server")
	log.Fatal(http.ListenAndServe(":6060", app.Router))
}
