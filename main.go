package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kyleterry/tweetsave/api"
	"github.com/kyleterry/tweetsave/db"
	"github.com/kyleterry/tweetsave/stream"
)

func main() {
	apiBind := flag.String("api-bind", "localhost:8092", "Host and port the API http server should listen on")
	dbBind := flag.String("db-bind", "postgres://localhost:5432/tweetsave?sslmode=disable", "Postgresql URL for the database connection")
	ckey := flag.String("consumer-key", "", "Twitter consumer key")
	csecret := flag.String("consumer-secret", "", "Twitter consumer secret")
	atoken := flag.String("access-token", "", "Twitter access token")
	asecret := flag.String("access-secret", "", "Twitter access secret")

	flag.Parse()

	if *ckey == "" || *csecret == "" || *atoken == "" || *asecret == "" {
		flag.Usage()
		log.Fatal("All twitter credentials are required")
	}

	dbConn := db.New(*dbBind)

	s := stream.New(dbConn,
		&stream.Config{
			*ckey,
			*csecret,
			*atoken,
			*asecret,
		})

	go s.Start()

	apiApp := api.New(dbConn)

	http.ListenAndServe(*apiBind, apiApp)
}
