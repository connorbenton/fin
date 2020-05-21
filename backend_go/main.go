// // main.go
package main

import (
	_ "fmt"
	//  "flag"
	_ "fintrack-go/config"
	"log"
	"net/http"

	//  "fintrack-go/app/server"
	"fintrack-go/app"
	"fintrack-go/db"
	"fintrack-go/socket"

	"github.com/gorilla/mux"

	//  "github.com/joho/godotenv"
	//  "github.com/jmoiron/sqlx"
	//  "github.com/sirupsen/logrus"
	_ "golang.org/x/sync/errgroup"
)

func main() {
	// Load Configurations
	//  var envfile string
	//  flag.StringVar(&envfile, "env-file", "../.env", "Read in a file of environment variables")
	//  flag.Parse()
	// godotenv.Load(envfile)
	//  config, err := config.Environ()
	//  if err != nil {
	//  log.Fatal("main: invalid config: %s", err.Error())
	//   logger := logrus.WithError(err)
	// logger.Fatalln("main: invalid configuration")
	//  }
	//Init logging
	//  initLogging(config)
	// if trace level logging is enabled, output the
	//  configuration parameters.
	//  if logrus.IsLevelEnabled(logrus.TraceLevel) {
	//   fmt.Println(config.String())
	//  }
	// Init DB
	_, err := db.CreateDatabase()
	if err != nil {
		log.Fatal("main: cannot initialize DB: %s", err.Error())
		//   logger := logrus.WithError(err)
		//   logger.Fatalln("main: cannot initialize DB")
	}

	// Init Currency DB
	_, err2 := db.CreateCurrencyDatabase()
	if err2 != nil {
		log.Fatal("main: cannot initialize Currency DB: %s", err.Error())
		//   logger := logrus.WithError(err)
		//   logger.Fatalln("main: cannot initialize DB")
	}

	db.GetNewXML()

	// Start the server
	//  g := errgroup.Group{}
	//  g.Go(func() error {
	//   logrus.WithFields(
	//    logrus.Fields{
	// "Host": config.Server.Host,
	//    },
	//   ).Info("starting the http server")

	socket.StartHub()

	app := &app.App{
		Router: mux.NewRouter().StrictSlash(true),
		// Database: DB,
	}
	// app.DB = DB

	app.SetupRouter()

	log.Fatal(http.ListenAndServe(":6060", app.Router))
	//   return app.Server.ListenAndServe()
	//  })

	// Wait the gorouitine
	//  if err := g.Wait(); err != nil {
	//   logrus.WithError(err).Fatalln("program terminated")
	//  }
}

// helper function configures the logging.
// func initLogging(c config.Config) {
//  if c.Logging.Debug {
//   logrus.SetLevel(logrus.DebugLevel)
//  }
//  if c.Logging.Trace {
//   logrus.SetLevel(logrus.TraceLevel)
//  }
//  if c.Logging.Text {
//   logrus.SetFormatter(&logrus.TextFormatter{
//    ForceColors:   c.Logging.Color,
//    DisableColors: !c.Logging.Color,
//   })
//  } else {
//   logrus.SetFormatter(&logrus.JSONFormatter{
//    PrettyPrint: c.Logging.Pretty,
//   })
//  }
// }

// package main

// import (
//    "io"
//    "log"
//    "net"
//    "github.com/google/wire/cmd/wire"
//    "github.com/valyala/fasthttp"
// )

// func handleConnection(conn net.Conn) {
//    _, err := io.Copy(conn, conn)
//    if err != nil {
//       log.Fatalf("Connection error %v", err)
//    }
// }

// func main() {
//    log.Println("Starting up repeater")
//    ln, err := net.Listen("tcp", ":7080")
//    if err != nil {
//       log.Fatalf("Cannot listen on port 7080 %v", err)
//    }

//    for {
//       conn, err := ln.Accept()
//       if err != nil {
//          log.Fatalf("Accepting connection error %v", err)
//       }
//       go handleConnection(conn)
//    }

// }
