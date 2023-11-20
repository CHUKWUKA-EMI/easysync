package data

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB ...
var DB *gorm.DB

// InitDatabaseConnection ...
func InitDatabaseConnection() {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")),
		DefaultStringSize: 256,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatal("Error initializing database connection: ", err.Error())
		os.Exit(1)
	}

	println("Database connected!")
	DB = db
}
