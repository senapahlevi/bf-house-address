package login

import (
	"housemap/databases"
	"housemap/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB
var secret_key = "tokentoken"

func SetDatabaseLogin(databased *databases.Database) {
	db = databased.DB
}
func LoginUser(c *gin.Context) {
	var request models.User
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	var user models.User
	if err := db.Where("username = ?", request.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func generateToken(UserId int) (string, error) {
	claims := jwt.MapClaims{
		"id": UserId,
		// "exp": time.Now().Add(time.Hour * 24).Unix(), //expiration for token
		"exp": time.Now().Add(time.Minute * 2).Unix(), //expiration for token
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret_key := "tokentoken"
	signedToken, err := token.SignedString([]byte(secret_key))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
