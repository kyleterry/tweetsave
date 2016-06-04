package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var router = mux.NewRouter()

type API struct {
	dbConn *gorm.DB
}

func New(dbConn *gorm.DB) *API {
	return &API{dbConn}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
