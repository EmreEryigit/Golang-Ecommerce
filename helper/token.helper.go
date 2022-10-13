package helper

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type SignedDetails struct {
	Email  string
	Name   string
	UserID string
	jwt.StandardClaims
}

var SESSION_KEY = []byte(os.Getenv("SESSION_KEY"))

func GenerateJWT(userId string, name string, email string) (signedToken string, err error) {

	if err != nil {
		log.Panic(err)
		return
	}
	claims := SignedDetails{
		Email:  email,
		Name:   name,
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(SESSION_KEY)
	if err != nil {
		log.Panic(err)
		return
	}
	return token, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return SESSION_KEY, nil
		},
	)
	if err != nil {
		msg = "error while parsing jwt"
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}
	return claims, msg
}
