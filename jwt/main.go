package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "top_secret"

// User represents login credentials
type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// Message represents a generic response message.
type Message struct {
	Status string `json:"status"`
	Msg    string `json:"message"`
}

// jsonMessageBytes returns a JSON-encoded response message.
func jsonMessageBytes(status, msg string) []byte {
	message := Message{
		Status: status,
		Msg:    msg,
	}

	b, err := json.Marshal(message)
	if err != nil {
		log.Printf("failed to marshal message: %v", err)
		return []byte(`{"status":"error", "message":"internal error"}`)
	}
	return b
}

// MyCustomClaims embeds jwt.RegisteredClaims and adds custom fields
type MyCustomClaims struct {
	UserName     string `json:"user_name"`
	LoggedInTime string `json:"logged_in_time"`
	jwt.RegisteredClaims
}

func CreateJWT() (string, error) {
	now := time.Now()

	claims := MyCustomClaims{
		UserName:     "barbarik",
		LoggedInTime: now.Format("02-01-2006 15:04:05"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Minute)),
			Issuer:    "barbarik",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

// ValidateJWT verifies and parses a JWT
func ValidateJWT(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Printf("failed to parse token: %v", err)
		return false
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		log.Printf("User: %s | LoggedIn: %s | Issuer: %s", claims.UserName, claims.LoggedInTime, claims.RegisteredClaims.Issuer)
		return true
	}

	log.Println("invalid token or claims")
	return false
}
