package users

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var tokenSecret = "522ac318-c050-4ed1-a360-10aa45211f58"

type serviceData struct {
	Path string
}

var service = serviceData{
	Path: initPath(),
}

func initPath() string {
	path, _ := filepath.Abs("../S4/files")
	return strings.Replace(path, "\\", "/", -1)
}

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
	}
	return nil, false

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
	userID := DB.resolveUserID(email)
	if !strings.EqualFold(userID, "") {
		token := createToken(userID)
		if !strings.EqualFold(token, "") {
			return token
		}
	}
	return ""
}

func getUsers() []getUsersResponse {
	return DB.getAllUsers()
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
	return DB.dbGetUserSecretByID(userID)
}

func addFileToUser(userID string, systemFileName string, userFileName string) {
	DB.dbAddFileToUser(userID, systemFileName, userFileName)
}

func checkFileOwner(userID string, fileID string) bool {
	return DB.dbcheckFileOwner(userID, fileID)
}

func getFileSystemNameByFileID(userID string, fileID string) string {
	return DB.dbGetFileSystemNameByFileID(userID, fileID)
}

func removeFileFromUser(userID string, fileID string) {
	DB.dbRemoveFileFromUser(userID, fileID)
}

func getFilesByUserID(userID string) []getFilesResponse {
	return DB.dbGetFilesByUserID(userID)
}
