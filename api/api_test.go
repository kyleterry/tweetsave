package api

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kyleterry/tweetsave/db"
)

var DB *gorm.DB

func init() {
	// DRY this up
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 10)
	for i := 0; i < 10; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	dbname := fmt.Sprintf("tweetsave-%s.db", string(result))

	DB = db.New("sqlite3", filepath.Join(os.TempDir(), dbname))
}

func TestIndexHandlerReturnsJsonResponse(t *testing.T) {
	DB.Create(
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

	server := httptest.NewServer(New(DB))
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
