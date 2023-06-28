package address

import (
	"fmt"
	"housemap/databases"
	"housemap/models"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDatabase(database *databases.Database) {
	db = database.DB
}
func GetDataHouse(c *gin.Context) {
	var house []models.House
	result := db.Find(&house)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, house)
}

func StoreDataHouse(c *gin.Context) {

	var house models.House
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

}

func GetDataHouseID(c *gin.Context) {

	var house models.House
	id := c.Param("id")
	result := db.Where("id = ?", id).First(&house)

	// err := c.ShouldBindJSON(&house)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "house not found"})
		return
	}

	c.JSON(http.StatusOK, house)

}

func UpdateDataHouseID(c *gin.Context) {

	var house models.House
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

}

func DeleteDataHouse(c *gin.Context) {

	// var house House
	var house models.House
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

}

func CalculateHouse(c *gin.Context) {
	var dataOriginHouse models.House
	var dataDestinationHouse models.House
	var otherHouse models.House
	var calculate models.Calculate

	err := c.ShouldBindJSON(&calculate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	originData := db.Where("id = ?", calculate.OriginID).First(&dataOriginHouse)

	if originData.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data Origin not found"})
		return
	}

	destinationData := db.Where("id = ?", calculate.DestinationID).First(&dataDestinationHouse)

	if destinationData.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error or=": destinationData.Error.Error()})
		return
	}

	if calculate.OtherStatus == 1 { //if using 3 posisi
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

	if otherHouse.ID == 0 && calculate.OtherStatus == 0 {
		otherHouseLat = destinationLat
		otherHouseLong = destinationLong
	}
	resultHaversine := haversineDistance(originLong, originLat, destinationLong, destinationLat) + haversineDistance(destinationLong, destinationLat, otherHouseLong, otherHouseLat)

	HaversineResponse := strconv.FormatFloat(resultHaversine, 'f', -1, 64) // s = "64.2345"

	response := ResponseCalculate{
		Haversine: HaversineResponse,
		// Euclidean:       EuclideanResponse,
		OriginLong:      originLong,
		OriginLat:       originLat,
		DestinationLong: destinationLong,
		DestinationLat:  destinationLat,
		OtherLong:       otherHouseLong,
		OtherLat:        otherHouseLat,
	}
	c.JSON(http.StatusOK, response)

}

//  radius bumi dalam km
const earthRadius = 6371.0

// konversi degrees ke radians
func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}
func haversineDistance(originLong float64, originLat float64, destinationLong float64, destinationLat float64) float64 {

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
