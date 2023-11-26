package data

import (
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains the common columns for all tables
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:char(36);primary_key;" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

// DB points to SQL Database connection ...
var DB *gorm.DB

// AstraDBSession points to Astra DB (Cassandra) session
var AstraDBSession *gocql.Session

// InitDatabaseConnection ...
func InitDatabaseConnection() {
	db, err := connectToSQLDB()

	if err != nil {
		log.Fatal("Error initializing database connection: ", err.Error())
	}

	println("SQL Database connected!")
	DB = db

	astraDb, err := connectToAstraDB()

	if err != nil {
		log.Fatalf("unable to connect astraDB session: %v", err)
	}
	println("AstraDB Connected!")
	AstraDBSession = astraDb

}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New()
	tx.Statement.SetColumn("ID", uuid)

	return nil
}
