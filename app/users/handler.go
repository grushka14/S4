package users

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/grushka14/S4/encryption"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoginHandler function handler
func LoginHandler(c *gin.Context) {
	var request userLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Parameters")
		return
	}
	token := validateUser(request.Email)
	if !strings.EqualFold(token, "") {
		c.JSON(http.StatusOK, token)
		return
	}
	c.JSON(http.StatusBadRequest, "Unauthorized")
}

// TokenValidationMiddleware Handler function handler
func TokenValidationMiddleware(c *gin.Context) {
	if c.Request.Header["Authorization"] == nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized - Please use a valid token")
		c.Abort()
		return
	}
	token := strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1]
	claims, ok := extractClaims(token)
	if !ok {
		c.JSON(http.StatusUnauthorized, "Unauthorized - Please use a valid token")
		c.Abort()
		return
	}
	userID := claims["id"].(string)
	if strings.EqualFold(userID, "") {
		c.JSON(http.StatusUnauthorized, "Unauthorized - Please use a valid token")
		c.Abort()
		return
	}
}

// GetUsersHandler function handler
func GetUsersHandler(c *gin.Context) {
	users := getUsers()
	c.JSON(http.StatusOK, users)
	return
}

// PutFileHandler function handler
func PutFileHandler(c *gin.Context) {
	userID := extractUserIDFromToken(strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1])
	userKey := getUserKey(userID)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	filename := header.Filename
	generatedFileName := strings.Replace(uuid.New().String(), "-", "", -1)
	path, _ := filepath.Abs("../S4/files")
	path = strings.Replace(path, "\\", "/", -1)
	out, err := os.Create(path + "/" + userID + "/" + generatedFileName)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, "Internal server error")
		return
	}
	byteFile, err := ioutil.ReadAll(file)
	encryptedFile, err := encryption.Encrypt(userKey, byteFile)

	addFileToUser(userID, generatedFileName, filename)

	defer out.Close()
	_, err = out.Write(encryptedFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusCreated, "Created")
}

// DeleteFileHandler function handler
func DeleteFileHandler(c *gin.Context) {
	userID := extractUserIDFromToken(strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1])
	fileID := c.Query("file-id")
	if !checkFileOwner(userID, fileID) {
		c.JSON(http.StatusUnauthorized, "You can only delete your files")
		return
	}
	fileName := getFileSystemNameByFileID(userID, fileID)
	if strings.EqualFold(fileName, "") {
		c.JSON(http.StatusBadRequest, "No file with that id")
		return
	}
	path, _ := filepath.Abs("../S4/files")
	path = strings.Replace(path, "\\", "/", -1)
	os.Remove(path + "/" + userID + "/" + fileName)
	removeFileFromUser(userID, fileID)
	c.JSON(http.StatusOK, "Deleted")
}

// GetFilesHandler function handler
func GetFilesHandler(c *gin.Context) {
	userID := extractUserIDFromToken(strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1])
	files := getFilesByUserID(userID)
	c.JSON(http.StatusOK, files)
}

// GetReadFileHandler function handler
func GetReadFileHandler(c *gin.Context) {
	userID := extractUserIDFromToken(strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1])
	fileID := c.Query("file-id")
	userKey := getUserKey(userID)
	if !checkFileOwner(userID, fileID) {
		c.JSON(http.StatusUnauthorized, "You can only delete your files")
		return
	}
	fileName := getFileSystemNameByFileID(userID, fileID)
	if strings.EqualFold(fileName, "") {
		c.JSON(http.StatusBadRequest, "No file with that id")
		return
	}
	path, _ := filepath.Abs("../S4/files")
	path = strings.Replace(path, "\\", "/", -1)
	file, err := ioutil.ReadFile(path + "/" + userID + "/" + fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Something whent wrong while reading the file")
		return
	}
	decryptedfile, err := encryption.Decrypt(userKey, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Something whent wrong while decrypting the file")
		return
	}
	c.JSON(http.StatusOK, decryptedfile)
}

func PostShareFile(c *gin.Context) {
	userID := extractUserIDFromToken(strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1])
	userKey := getUserKey(userID)
	var request shareFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid Params")
		return
	}
	if strings.EqualFold(request.FileID, "") {
		c.JSON(http.StatusBadRequest, "Invalid file_id")
	}
	if strings.EqualFold(request.UserID, "") {
		c.JSON(http.StatusBadRequest, "Invalid user_id")
	}
	if !checkFileOwner(userID, request.FileID) {
		c.JSON(http.StatusUnauthorized, "You can only delete your files")
		return
	}
	fileName := getFileSystemNameByFileID(userID, request.FileID)
	if strings.EqualFold(fileName, "") {
		c.JSON(http.StatusBadRequest, "No file with that id")
		return
	}
	path, _ := filepath.Abs("../S4/files")
	path = strings.Replace(path, "\\", "/", -1)
	file, err := ioutil.ReadFile(path + "/" + userID + "/" + fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Something whent wrong while reading the file")
		return
	}
	decryptedfile, err := encryption.Decrypt(userKey, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Something whent wrong while decrypting the file")
		return
	}
	generatedFileName := strings.Replace(uuid.New().String(), "-", "", -1)
	out, err := os.Create(path + "/" + request.UserID + "/" + generatedFileName)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, "Failed to genetate new name for the file")
		return
	}
	key := getUserKey(request.UserID)
	encryptedFile, err := encryption.Encrypt(key, decryptedfile)
	defer out.Close()
	_, err = out.Write(encryptedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to encrypt new name for the file")
		return
	}
	addFileToUser(request.UserID, generatedFileName, fileName)
	c.JSON(http.StatusCreated, "created")
}
