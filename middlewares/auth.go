package middlewares

import (
	"fmt"
	"sanchit/constants"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(req *gin.Context) {
		authHeader := req.GetHeader("Authorization")
		if authHeader == "" || (len(authHeader) <= 7) {
			req.JSON(401, gin.H{"error": "Authorization Header not found"})
			return
		}
		tokenString := authHeader[7:]
		var temp JwtClaims
		token, err := jwt.ParseWithClaims(tokenString, &temp, func(token *jwt.Token) (interface{}, error) {
			return constants.SecretKey, nil
		})
		if err != nil {
			//log.Fatal(err)
			req.JSON(401, gin.H{"error": "Authorization Header not found"})
			return
		} else if claims, ok := token.Claims.(*JwtClaims); ok {
			fmt.Println(claims.Email, claims.RegisteredClaims.Issuer)
			req.Set("Email", claims.Email)
		} else {
			//log.Fatal("Unknown claims type, cannot proceed")
			req.JSON(401, gin.H{"error": "Unknown claims type, cannot proceed"})
			return
		}
		req.Next()

	}
}
