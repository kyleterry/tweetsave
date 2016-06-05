package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var router = mux.NewRouter()

type errorHandler func(http.ResponseWriter, *http.Request) error

func (fn errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err != nil {
		log.Printf("Got error while processing the request: %s\n", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
	}
}

type API struct {
	dbConn *gorm.DB
}

func New(dbConn *gorm.DB) *API {
	a := &API{dbConn}
	router.Handle("/", errorHandler(a.IndexHandler)).Methods("GET").Name("index")

	return a
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}

func (a *API) IndexHandler(w http.ResponseWriter, r *http.Request) error {
	var result []struct {
		URL      string `json:"url"`
		PostedOn string `json:"posted_on"`
		PostedBy string `json:"posted_by"`
	}
	a.dbConn.Table("tweet_urls").Select(
		"tweet_urls.url as url, tweet_urls.created_at as posted_on, users.name as posted_by").Joins(
		"JOIN users on tweet_urls.user_id = users.id").Scan(&result)
	b, err := json.Marshal(&result)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}
