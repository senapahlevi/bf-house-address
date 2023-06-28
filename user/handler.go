package user

import (
	"fmt"
	"housemap/databases"
	"housemap/middleware"
	"housemap/models"
	"housemap/register"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDatabaseUser(databased *databases.Database) {
	db = databased.DB
}

func UserFound(c *gin.Context) {
	user := c.MustGet("user").(*register.Claims)

	auth := middleware.Authentication

	var users models.User
	if err := db.First(&users, user.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	fmt.Println("hello 5")
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Welcome, %s!", users.Username)})

}
