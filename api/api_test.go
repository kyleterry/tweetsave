package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/kyleterry/tweetsave/db"
)

func setup() *gorm.DB {
	return db.New("postgres://localhost:5432/tweetsavetest?sslmode=disable")
}

func teardown(dbconn *gorm.DB) {
	dbconn.DropTableIfExists(&db.TweetURL{}, &db.User{})
}

func TestIndexHandlerReturnsJsonResponse(t *testing.T) {
	dbconn := setup()

	// clean up the test DB
	defer teardown(dbconn)

	dbconn.Create(
		&db.TweetURL{
			URL: "http://for.tn/1t8eq4f",
			User: db.User{
				Name: "FortuneMagazine",
			},
		}).Create(
		&db.TweetURL{
			URL: "http://slnm.us/CVeIunQ",
			User: db.User{
				Name: "Salon",
			},
		})

	server := httptest.NewServer(New(dbconn))
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatal("expected 200, got", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range []string{"http://for.tn/1t8eq4f", "FortuneMagazine", "http://slnm.us/CVeIunQ", "Salon"} {
		if !strings.Contains(string(b), item) {
			t.Error("result string did not contain", item)
		}
	}
}
