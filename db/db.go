package db

import (
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type TweetURL struct {
	gorm.Model
	URL    string `gorm:"not null;unique"`
	User   User
	UserID uint
}

type User struct {
	gorm.Model
	Name string `gorm:"not null;unique"`
}

func New(engine string, dbURL string) *gorm.DB {
	conn, err := gorm.Open(engine, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	conn.LogMode(true)

	// create the tables if they don't exist
	conn.AutoMigrate(&TweetURL{}, &User{})

	return conn
}

func IsUniqueConstraintErr(err error) bool {
	return strings.Contains(err.Error(), "unique constraint")
}
