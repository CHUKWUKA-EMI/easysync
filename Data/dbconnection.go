package data

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// BaseModel contains the common columns for all tables
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:char(36);primary_key;" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

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

// BeforeCreate will set a UUID rather than numeric ID.
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New()
	tx.Statement.SetColumn("ID", uuid)

	return nil
}
