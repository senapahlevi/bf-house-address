package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type House struct {
	ID        int       `gorm:"id"`
	Tipe      string    `gorm:"tipe"`
	Alamat    string    `gorm:"alamat"`
	Lat       string    `gorm:"lat"`
	Long      string    `gorm:"long"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
	DeletedAt time.Time `gorm:"deleted_at"`
}

// func (house *House) BeforeSave(tx *gorm.DB) (err error) {
// 	long, err := strconv.ParseFloat(house.Long, 64)
// 	if err != nil {
// 		return err
// 	}
// 	lat, err := strconv.ParseFloat(house.Lat, 64)
// 	if err != nil {
// 		return err
// 	}
// 	house.Long = strconv.FormatFloat(long, 'f', -1, 64)
// 	house.Lat = strconv.FormatFloat(lat, 'f', -1, 64)
// 	return nil
// }

// func createRumah(db *gorm.DB) error {
// 	rumah := House{
// 		Tipe:   "Tipe Rumah A",
// 		Alamat: "Alamat Rumah",
// 		Lat:    123.456,
// 		Long:   789.012,
// 	}

// 	err := db.Create(&rumah).Error
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func main() {
	dsn := "host=localhost user=postgres password=123456789 dbname=address-house-map-bf port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		fmt.Println("gagal")
	}

	fmt.Println("berhasil konek")
	err = db.AutoMigrate(&House{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = createRumah(db)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	router := gin.Default()

	// config := cors.DefaultConfig()
	// config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	// config.AllowHeaders = []string{"Content-Type"}
	// router.Use(cors.New(config))

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		// AllowHeaders:     []string{"Origin"},
		// ExposeHeaders:    []string{"Content-Type"},
		AllowHeaders: []string{"Content-Type"},
		// ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
	}))
	// Access-Control-Allow-Origin

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
	api.PUT("/house/:id", func(c *gin.Context) {
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

	api.POST("/calculate-route", func(c *gin.Context) {
		// var house House
		var house House
		id := c.Param("id")
		result := db.Model(&house).Where("id = ?", id).First(&house)
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
	router.Run(":8080")
}
