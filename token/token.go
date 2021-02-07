package token

import (
	"asaf_project/src/domain/user"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const secret = "supersecretkeyblabla"

//Token type
type Token (string)

//CreateToken function
func CreateToken(u user.User) (Token, error) {
	os.Setenv("ACCESS_SECRET", secret)
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["email"] = u.Email
	atClaims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return Token(token), nil
}

type tokenClaims struct {
	authorized string `json:"authorized"`
	email      string `json:"email"`
	exp        int64  `json:"exp"`
	jwt.StandardClaims
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecret := []byte(secret)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		fmt.Println("Invalid JWT Token")
		return nil, false
	}
}

// Validate method
func (t Token) Validate() (string, error) {
	tokenString := string(t)
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if !token.Valid {
		return "unauthorized", nil
	}

	// now := time.Now().UTC()
	// expiration := time.Unix(claims["exp"].(int64), 0)
	// if now.After(expiration) {
	// 	return "expired", nil
	// }

	return "authorized", nil
}
