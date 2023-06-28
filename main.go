package main

import (
	"fmt"
	"housemap/address"
	"housemap/databases"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// host := os.Getenv("DB_HOST")
	// user := os.Getenv("DB_USER")
	// dbnames := os.Getenv("DB_NAME")
	// pass := os.Getenv("DB_PASS")
	// dbport := os.Getenv("DB_PORT")

	// // dsn := "host=containers-us-west-98.railway.app user=postgres password=04njxElMvMSRpaqSSafl dbname=railway port=7927 TimeZone=Asia/Jakarta" //local
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Jakarta", host, user, pass, dbnames, dbport) //local
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal(err)
	// 	fmt.Println("gagal")
	// }
	fmt.Println("berhasil konek")
	// err = db.AutoMigrate(&Calculate{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	database, err := databases.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	address.SetDatabase(database)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "https://next-bf-home-routes.vercel.app"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type"},
		// ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
	}))

	api := router.Group("/api/v1/")

	api.GET("/data-house", address.GetDataHouse)
	api.POST("/house", address.StoreDataHouse)

	api.GET("/house/:id", address.GetDataHouseID)

	api.PUT("/house-update/:id", address.UpdateDataHouseID)

	api.DELETE("/house/:id", address.DeleteDataHouse)

	api.POST("/calculate-route", address.CalculateHouse)
	// router.Run(":8080")
	router.Run(":" + os.Getenv("PORT"))

}

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
