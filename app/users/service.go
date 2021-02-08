package users

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var tokenSecret = "522ac318-c050-4ed1-a360-10aa45211f58"

func createToken(userID string) string {
	os.Setenv("ACCESS_SECRET", tokenSecret)
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = userID
	atClaims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return token
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecret := []byte(tokenSecret)
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

func validateToken(tokenString string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if !token.Valid {
		return "unauthorized", nil
	}

	return "authorized", nil
}

// Check if the email is in the database, if so generate a token and return it.
func validateUser(email string) string {
	userID := resolveUserID(email)
	if !strings.EqualFold(userID, "") {
		token := createToken(userID)
		if !strings.EqualFold(token, "") {
			return token
		}
	}
	return ""
}

func getUsers() []getUsersResponse {
	return getAllUsers()
}

func extractUserIDFromToken(tokenString string) string {
	claims := jwt.MapClaims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if token.Valid {
		return claims["id"].(string)
	}
	return ""
}

func getUserKey(userID string) []byte {
	return dbGetUserSecretByID(userID)
}

func addFileToUser(userID string, systemFileName string, userFileName string) {
	dbAddFileToUser(userID, systemFileName, userFileName)
}

func checkFileOwner(userID string, fileID string) bool {
	if dbcheckFileOwner(userID, fileID) {
		return true
	}
	return false
}

func getFileSystemNameByFileID(userID string, fileID string) string {
	return dbGetFileSystemNameByFileID(userID, fileID)
}

func removeFileFromUser(userID string, fileID string) {
	dbRemoveFileFromUser(userID, fileID)
}

func getFilesByUserID(userID string) []getFilesResponse {
	return dbGetFilesByUserID(userID)
}
