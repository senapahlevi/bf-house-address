package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Rumah struct {
	Tipe   string
	Alamat string
	Lat    float64
	Long   float64
}

func initialDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=123456789 dbname=rumah port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Rumah{})

	return db, nil
}

func main() {
	db, err := initialDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	// defer db.Close()
	defer db.SQLDB().Close()

	// r := gin.Default()

}
