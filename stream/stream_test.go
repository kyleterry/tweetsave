package stream

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dghubble/go-twitter/twitter"
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

type RewriteTransport struct {
	Transport http.RoundTripper
}

func (t *RewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	if t.Transport == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return t.Transport.RoundTrip(req)
}

func testServer() *http.Client {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	transport := &RewriteTransport{&http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}}
	client := &http.Client{Transport: transport}
	return client
}

func TestStreamTweetHandlerCanStoreTweetURLs(t *testing.T) {
	httpclient := testServer()

	client := twitter.NewClient(httpclient)
	stream, err := client.Streams.User(&twitter.StreamUserParams{})
	if err != nil {
		t.Fatal(err)
	}
	s := Stream{stream: stream, dbConn: DB}
	go s.Start()

	defer s.stream.Stop()

	// make fake tweet
	tweet := twitter.Tweet{
		Text: "test tweet",
		User: &twitter.User{ScreenName: "test_user"},
		Entities: &twitter.Entities{
			Urls: []twitter.URLEntity{
				twitter.URLEntity{
					ExpandedURL: "http://example.com",
				},
			},
		},
	}

	// send it over messages channel
	s.stream.Messages <- &tweet

	// This is a hack. Figure out how to get around the need for this in this test.
	<-time.After(time.Second)

	// see if the record is in the DB
	var url db.TweetURL
	DB.First(&url)

	if url.URL != "http://example.com" {
		t.Error("Resulting URL did not contain http://example.com")
	}
}

func TestStreamTweetHandlerOnlyStoresUniqueURLs(t *testing.T) {
	s := Stream{dbConn: DB}

	// make fake tweets
	tweet := twitter.Tweet{
		Text: "test tweet",
		User: &twitter.User{ScreenName: "test_user"},
		Entities: &twitter.Entities{
			Urls: []twitter.URLEntity{
				twitter.URLEntity{
					ExpandedURL: "http://example.com",
				},
			},
		},
	}

	tweet2 := twitter.Tweet{
		Text: "test tweet",
		User: &twitter.User{ScreenName: "test_user2"},
		Entities: &twitter.Entities{
			Urls: []twitter.URLEntity{
				twitter.URLEntity{
					ExpandedURL: "http://example.com",
				},
			},
		},
	}

	s.tweetHandler(&tweet)
	s.tweetHandler(&tweet2)

	// see if the record is in the DB
	var urls []db.TweetURL
	DB.Where("url = ?", "http://example.com").Find(&urls)

	if len(urls) > 1 {
		t.Error("Stream tweet handler should only store unique urls.")
	}
}
