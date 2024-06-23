package controllers

import (
	"net/http"
	"sanchit/constants"
	"sanchit/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func UserLogin(req *gin.Context) {
	var user models.User
	err := req.BindJSON(&user)
	if (err != nil) || (!models.VerifyUser(&user)) {
		req.JSON(http.StatusInternalServerError, gin.H{"error": "Login Error"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"Email": user.Email,
			"exp":   time.Now().Add(time.Hour).Unix(),
		})
	tokenString, err := token.SignedString(constants.SecretKey)
	if err != nil {
		req.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	req.JSON(http.StatusOK, gin.H{"token": tokenString})

}
