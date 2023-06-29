package middleware

import (
	"fmt"
	"housemap/register"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) *register.Claims {

	//no works
	// user := c.MustGet("user").(*jwt.Token)
	// claims := user.Claims.(jwt.MapClaims)

	// userID := uint(claims["id"].(float64))

	userAuth := c.MustGet("user").(*register.Claims)
	return userAuth

}

var secret_key = "tokentoken"

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		fmt.Println("hello auth 1")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}
		tokenString := authHeader[len("Bearer "):]

		token, err := jwt.ParseWithClaims(tokenString, &register.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}
			return []byte(secret_key), nil
		})
		if err != nil {
			//if expired will showing these message
			// c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.JSON(http.StatusUnauthorized, gin.H{"status": "expired"})
			c.Abort()
			return
		}
		fmt.Println("hello auth 2")

		claims, ok := token.Claims.(*register.Claims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		fmt.Println("hello auth 3")

		c.Set("user", claims)
		fmt.Println("hello auth 4", claims)
		c.Next()

	}
}
