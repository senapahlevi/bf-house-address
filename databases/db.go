package databases

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase() (*Database, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	dbnames := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASS")
	dbport := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Jakarta", host, user, pass, dbnames, dbport)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Database{
		DB: db,
	}, nil
}
