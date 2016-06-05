package stream

import (
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/jinzhu/gorm"
	"github.com/kyleterry/tweetsave/db"
)

type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

type Stream struct {
	config *Config
	stream *twitter.Stream
	dbConn *gorm.DB
}

func New(dbConn *gorm.DB, C *Config) *Stream {
	c := oauth1.NewConfig(C.ConsumerKey, C.ConsumerSecret)
	token := oauth1.NewToken(C.AccessToken, C.AccessSecret)
	httpClient := c.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	params := &twitter.StreamUserParams{
		StallWarnings: twitter.Bool(true),
		With:          "followings",
	}

	stream, err := client.Streams.User(params)
	if err != nil {
		log.Fatal(err)
	}

	return &Stream{C, stream, dbConn}
}

func (s *Stream) Start() {
	// We only care about tweets, so demux those from interface{} and handle
	demux := twitter.NewSwitchDemux()
	demux.Tweet = s.tweetHandler
	demux.HandleChan(s.stream.Messages)
}

func (s *Stream) tweetHandler(tweet *twitter.Tweet) {
	user := db.User{}
	s.dbConn.FirstOrInit(&user, db.User{Name: tweet.User.ScreenName})
	for _, url := range tweet.Entities.Urls {
		log.Println("saving url from tweet")
		s.dbConn.Create(&db.TweetURL{
			URL:  url.ExpandedURL,
			User: user,
		})
	}
}
