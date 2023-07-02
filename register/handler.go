package register

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

func SetDatabaseRegister(databased *databases.Database) {
	db = databased.DB
}
func RegisterUser(c *gin.Context) {
	var register RegisterRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check email
	var existEmail models.User
	if err := db.Where("email = ?", register.Email).First(&existEmail).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already taken", "status": "denied"})
		return
	}
	//check username
	var existUsername models.User
	if err := db.Where("username = ?", register.Username).First(&existUsername).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already taken", "status": "denied"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password", "status": "denied"})
		return
	}
	user := models.User{
		Username: register.Username,
		Email:    register.Email,
		Password: string(hashedPassword),
	}
	// token, err := generateToken(user.ID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token "})
	// 	return
	// }
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "status": "denied"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func generateToken(UserId int) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24)
	claims := &Claims{
		UserID: uint(UserId),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	secret_key := "tokentoken"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret_key))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
// 			c.Abort()
// 			return
// 		}
// 		tokenString := authHeader[len("Bearer "):]

// 		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
// 			}
// 			return []byte(secret_key), nil
// 		})
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
// 			c.Abort()
// 			return
// 		}
// 		claims, ok := token.Claims.(*Claims)
// 		if !ok || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
// 			c.Abort()
// 			return
// 		}
// 		c.Set("user", claims)
// 		c.Next()
// 	}
// }
