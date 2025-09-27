package utils

import (
	"time"
	"os"
	"github.com/golang-jwt/jwt/v5"
)

var jwtkey = []byte(os.Getenv("Security_Key"))

func GenerateJWT(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"role":   role,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtkey)
}

func ExtractSecertKey(token *jwt.Token)(interface{},error){
	return jwtkey,nil
}