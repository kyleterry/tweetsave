package db

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type TweetURL struct {
	gorm.Model
	URL    string
	User   User
	UserID int
}

type User struct {
	gorm.Model
	Name string
}

func New(dbURL string) *gorm.DB {
	conn, err := gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// create the tables if they don't exist
	conn.AutoMigrate(&TweetURL{}, &User{})

	return conn
}
