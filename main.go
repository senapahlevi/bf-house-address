package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type House struct {
	ID        uint      `gorm:"id"`
	Tipe      string    `gorm:"tipe"`
	Alamat    string    `gorm:"alamat"`
	Lat       string    `gorm:"lat"`
	Long      string    `gorm:"long"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
	DeletedAt time.Time `gorm:"deleted_at"`
}

type Calculate struct {
	OriginID        int       `gorm:"originid"`
	DestinationID   int       `gorm:"destinationid"`
	LatOrigin       string    `gorm:"lat_origin"`
	LongOrigin      string    `gorm:"long_origin"`
	LatDestination  string    `gorm:"lat_destination"`
	LongDestination string    `gorm:"long_destination"`
	CreatedAt       time.Time `gorm:"created_at"`
	UpdatedAt       time.Time `gorm:"updated_at"`
}

func calculateDistance(originLat, originLong, destinationLat, destinationLong float64) float64 {

	lat1 := originLat
	long1 := originLong
	lat2 := destinationLat
	long2 := destinationLong

	return math.Sqrt(math.Pow(lat2-lat1, 2) + math.Pow(long2-long1, 2))
}

//  radius bumi dalam km
const earthRadius = 6371.0

// konversi degrees ke radians
func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// Kalkulasi 2 titik dalam km
func haversineDistance(originLong float64, originLat float64, destinationLong float64, destinationLat float64) float64 {

	// Convert the latitude and longitude values to radians
	originLat = deg2rad(originLat)
	originLong = deg2rad(originLong)
	destinationLat = deg2rad(destinationLat)
	destinationLong = deg2rad(destinationLong)

	// Calculate the differences between the latitude and longitude values
	dLat := destinationLat - originLat
	dLon := destinationLong - originLong

	// Apply the haversine formula
	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(originLat)*math.Cos(destinationLat)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := earthRadius * c
	fmt.Println("hello ini KM haversineDistance", d)
	// Return the distance in kilometers
	return d
}

func main() {
	//cloud start

	//r
	dsn := "host=containers-us-west-98.railway.app user=postgres password=04njxElMvMSRpaqSSafl dbname=railway port=7927 TimeZone=Asia/Jakarta" //local
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		fmt.Println("gagal")
	}
	//

	fmt.Println("berhasil konek")
	// err = db.AutoMigrate(&Calculate{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "https://next-bf-home-routes.vercel.app"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type"},
		// ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
	}))

	api := router.Group("/api/v1/")
	api.GET("/data-house", func(c *gin.Context) {
		var house []House
		result := db.Find(&house)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, house)
	})

	api.POST("/house", func(c *gin.Context) {
		var house House
		err := c.ShouldBindJSON(&house)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := db.Create(&house)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, house)
	})

	api.GET("/house/:id", func(c *gin.Context) {
		var house House
		id := c.Param("id")
		result := db.Where("id = ?", id).First(&house)

		// err := c.ShouldBindJSON(&house)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "house not found"})
			return
		}

		c.JSON(http.StatusOK, house)
	})
	api.PUT("/house-update/:id", func(c *gin.Context) {
		var house House
		id := c.Param("id")
		// Pencarian data berdasarkan ID
		result := db.Model(&house).Where("id = ?", id).First(&house)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "house not found"})
			return
		}

		// Perbarui data berdasarkan permintaan HTTP
		err := c.ShouldBindJSON(&house)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Simpan perubahan ke dalam database
		result = db.Save(&house)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, house)
	})

	api.DELETE("/house/:id", func(c *gin.Context) {
		// var house House
		var house House
		id := c.Param("id")
		result := db.Where("id = ?", id).First(&house)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "house not found"})
			return
		}

		result = db.Delete(&house)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, house)
	})

	type ResponseCalculate struct {
		Haversine string `json:"haversine"`
		Euclidean string `json:"euclidean"`
	}
	fmt.Println("hello port", os.Getenv(fmt.Sprint("PORT")))

	api.POST("/calculate-route", func(c *gin.Context) {
		var dataOriginHouse House
		var dataDestinationHouse House
		var calculate Calculate

		err := c.ShouldBindJSON(&calculate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		originData := db.Where("id = ?", calculate.OriginID).First(&dataOriginHouse)

		fmt.Println("calcuate origin lat", dataOriginHouse.Lat)
		// id := c.Param("id")
		if originData.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data Origin not found"})
			return
		}

		destinationData := db.Where("id = ?", calculate.DestinationID).First(&dataDestinationHouse)

		fmt.Println("originData", originData)
		if destinationData.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error or=": destinationData.Error.Error()})
			return
		}

		originLat, _ := strconv.ParseFloat(dataOriginHouse.Lat, 64)
		originLong, _ := strconv.ParseFloat(dataOriginHouse.Long, 64)
		destinationLat, _ := strconv.ParseFloat(dataDestinationHouse.Lat, 64)
		destinationLong, _ := strconv.ParseFloat(dataDestinationHouse.Long, 64)

		resultHaversine := haversineDistance(originLong, originLat, destinationLong, destinationLat)
		resultCalculate := calculateDistance(originLat, originLong, destinationLat, destinationLong)

		// optimumRoute := append(calculate, calculateDistance)
		// var hasil []CalculateResult
		HaversineResponse := strconv.FormatFloat(resultHaversine, 'f', -1, 64) // s = "64.2345"
		EuclideanResponse := strconv.FormatFloat(resultCalculate, 'f', -1, 64) // s = "64.2345"

		response := ResponseCalculate{
			Haversine: HaversineResponse,
			Euclidean: EuclideanResponse,
		}
		c.JSON(http.StatusOK, response)
		// c.JSON(http.StatusOK, route)
	})
	router.Run(":8080")
	// router.Run(":" + os.Getenv("PORT"))
	// router.Run(":0.0.0.0")

	// router.Run(":" + os.Getenv("PORT"))
}

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
