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
	"github.com/golang/geo/r3"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type House struct {
	ID        uint      `gorm:"id"`
	Tipe      string    `gorm:"tipe" binding:"required"`
	Alamat    string    `gorm:"alamat" binding:"required"`
	Lat       string    `gorm:"lat" binding:"required"`
	Long      string    `gorm:"long" binding:"required"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
	DeletedAt time.Time `gorm:"deleted_at"`
}

type Calculate struct {
	OriginID        int    `gorm:"originid binding:"`
	DestinationID   int    `gorm:"destinationid"`
	LatOrigin       string `gorm:"lat_origin"`
	LongOrigin      string `gorm:"long_origin"`
	LatDestination  string `gorm:"lat_destination"`
	LongDestination string `gorm:"long_destination"`
	//
	OtherStatus  int       `gorm:"other_status"`
	OtherID      int       `gorm:"otherid"`
	LatOther     string    `gorm:"lat_other"`
	LongLatOther string    `gorm:"long_other"`
	CreatedAt    time.Time `gorm:"created_at"`
	UpdatedAt    time.Time `gorm:"updated_at"`
}

func EuclideanDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Earth radius in kilometers
	const R = 6371

	// Convert latitude and longitude to Cartesian coordinates
	p1 := r3.Vector{
		X: R * math.Cos(lat1) * math.Cos(lon1),
		Y: R * math.Cos(lat1) * math.Sin(lon1),
		Z: R * math.Sin(lat1),
	}

	p2 := r3.Vector{
		X: R * math.Cos(lat2) * math.Cos(lon2),
		Y: R * math.Cos(lat2) * math.Sin(lon2),
		Z: R * math.Sin(lat2),
	}

	// Compute the Euclidean distance between the points
	return p1.Sub(p2).Norm()
}

//  radius bumi dalam km
const earthRadius = 6371.0

// konversi degrees ke radians
func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// Kalkulasi 2 titik dalam km
func haversineDistance(originLong float64, originLat float64, destinationLong float64, destinationLat, otherHouseLong, otherHouseLat float64) float64 {

	// Convert the latitude and longitude values to radians
	originLat = deg2rad(originLat)
	originLong = deg2rad(originLong)
	destinationLat = deg2rad(destinationLat)
	destinationLong = deg2rad(destinationLong)

	// kalkulasi  latitude dan longitude
	dLat := destinationLat - originLat
	dLon := destinationLong - originLong

	// rumus  haversine
	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(originLat)*math.Cos(destinationLat)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := earthRadius * c
	fmt.Println("hello ini KM haversineDistance", d)
	// satuan km
	return d
}

func main() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	dbnames := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASS")
	dbport := os.Getenv("DB_PORT")

	// dsn := "host=containers-us-west-98.railway.app user=postgres password=04njxElMvMSRpaqSSafl dbname=railway port=7927 TimeZone=Asia/Jakarta" //local
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Jakarta", host, user, pass, dbnames, dbport) //local
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
		Haversine       string  `json:"haversine"`
		Euclidean       string  `json:"euclidean"`
		OriginLong      float64 `json:"origin_long"`
		OriginLat       float64 `json:"origin_lat"`
		DestinationLat  float64 `json:"destination_lat"`
		DestinationLong float64 `json:"destination_long"`
		OtherLat        float64 `json:"other_lat"`
		OtherLong       float64 `json:"other_long"`
	}

	api.POST("/calculate-route", func(c *gin.Context) {
		var dataOriginHouse House
		var dataDestinationHouse House
		var otherHouse House
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
		fmt.Println(calculate.OtherStatus, "hello otherstatus")
		fmt.Println(calculate.OtherID, "hello other id")
		if calculate.OtherStatus == 1 {
			otherHouseData := db.Where("id = ?", calculate.OtherID).First(&otherHouse)
			if otherHouseData.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": otherHouseData.Error.Error()})
				return
			}
		}

		originLat, _ := strconv.ParseFloat(dataOriginHouse.Lat, 64)
		originLong, _ := strconv.ParseFloat(dataOriginHouse.Long, 64)
		destinationLat, _ := strconv.ParseFloat(dataDestinationHouse.Lat, 64)
		destinationLong, _ := strconv.ParseFloat(dataDestinationHouse.Long, 64)

		otherHouseLat, _ := strconv.ParseFloat(otherHouse.Lat, 64)
		otherHouseLong, _ := strconv.ParseFloat(otherHouse.Long, 64)

		resultHaversine := haversineDistance(originLong, originLat, destinationLong, destinationLat, otherHouseLong, otherHouseLat)
		// resultCalculate := calculateDistance(originLat, originLong, destinationLat, destinationLong, otherHouseLong, otherHouseLat)
		// 	lat1 := -6.21462 * math.Pi / 180 // convert degrees to radians
		lat1 := originLat * math.Pi / 180 // convert degrees to radians
		long1 := originLong * math.Pi / 180
		lat2 := destinationLat * math.Pi / 180
		long2 := destinationLong * math.Pi / 180
		lat3 := otherHouseLat * math.Pi / 180
		long3 := otherHouseLong * math.Pi / 180
		fmt.Println("long3 hello ", long3)
		fmt.Println("lat3 hello ", lat3)
		resultEuclidean := EuclideanDistance(lat1, long1, lat2, long2) + EuclideanDistance(lat2, long2, lat3, long3)
		fmt.Println(EuclideanDistance(lat1, long1, lat2, long2), "hello ab")
		fmt.Println(EuclideanDistance(lat2, long2, lat3, long3), "hello bc")
		// fmt.Println(EuclideanDistance(latA, lonA, latB, lonB) + EuclideanDistance(latB, lonB, latC, lonC))
		// }
		// optimumRoute := append(calculate, calculateDistance)
		// var hasil []CalculateResult
		HaversineResponse := strconv.FormatFloat(resultHaversine, 'f', -1, 64) // s = "64.2345"
		// EuclideanResponse := strconv.FormatFloat(resultCalculate, 'f', -1, 64) // s = "64.2345"
		EuclideanResponse := strconv.FormatFloat(resultEuclidean, 'f', -1, 64) // s = "64.2345"

		response := ResponseCalculate{
			Haversine:       HaversineResponse,
			Euclidean:       EuclideanResponse,
			OriginLong:      originLong,
			OriginLat:       originLat,
			DestinationLong: destinationLong,
			DestinationLat:  destinationLat,
			OtherLong:       otherHouseLong,
			OtherLat:        otherHouseLat,
		}
		c.JSON(http.StatusOK, response)
		// c.JSON(http.StatusOK, route)
	})
	// router.Run(":8080")
	router.Run(":" + os.Getenv("PORT"))
	// router.Run(":0.0.0.0")

	// router.Run(":" + os.Getenv("PORT"))
}

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
