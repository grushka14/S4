package users

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/grushka14/S4/encryption"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoginHandler function handler
func LoginHandler(c *gin.Context) {
	var request userLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "Bad Parameters")
		return
	}
	token := validateUser(request.Email)
	if token != "" {
		c.String(http.StatusOK, token)
		return
	}
	c.String(http.StatusUnauthorized, "Unauthorized")
}

// TokenValidationMiddleware Handler function handler
func TokenValidationMiddleware(c *gin.Context) {
	if c.Request.Header["Authorization"] == nil {
		c.String(http.StatusUnauthorized, "Unauthorized - Please use a valid token")
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
	if userID == "" {
		c.String(http.StatusUnauthorized, "Unauthorized - Please use a valid token")
		c.Abort()
		return
	}
}

// GetUsersHandler function handler
func GetUsersHandler(c *gin.Context) {
	c.JSON(http.StatusOK, getUsers())
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
	out, err := os.Create(service.Path + "/" + userID + "/" + generatedFileName)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Internal server error")
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
	c.String(http.StatusCreated, "Created")
}

// DeleteFileHandler function handler
func DeleteFileHandler(c *gin.Context) {
	userID := extractUserIDFromToken(strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1])
	fileID := c.Query("file-id")
	if !checkFileOwner(userID, fileID) {
		c.String(http.StatusUnauthorized, "You can only delete your files")
		return
	}
	fileName := getFileSystemNameByFileID(userID, fileID)
	if fileName == "" {
		c.String(http.StatusBadRequest, "No file with that id")
		return
	}

	os.Remove(service.Path + "/" + userID + "/" + fileName)
	removeFileFromUser(userID, fileID)
	c.String(http.StatusOK, "Deleted")
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
		c.String(http.StatusUnauthorized, "You can only delete your files")
		return
	}
	fileName := getFileSystemNameByFileID(userID, fileID)
	if fileName == "" {
		c.String(http.StatusBadRequest, "No file with that id")
		return
	}
	file, err := ioutil.ReadFile(service.Path + "/" + userID + "/" + fileName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something whent wrong while reading the file")
		return
	}
	decryptedfile, err := encryption.Decrypt(userKey, file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something whent wrong while decrypting the file")
		return
	}

	generatedFileName := strings.Replace(uuid.New().String(), "-", "", -1)
	filePath := service.Path + "/" + userID + "/" + generatedFileName + "DELETE"
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}
	defer out.Close()
	_, err = out.Write(decryptedfile)

	c.File(filePath)

	os.Remove(filePath)

	return

}

// PostShareFile function handler
func PostShareFile(c *gin.Context) {
	userID := extractUserIDFromToken(strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1])
	userKey := getUserKey(userID)
	var request shareFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "Invalid Params")
		return
	}
	if request.FileID == "" {
		c.String(http.StatusBadRequest, "Invalid file_id")
		return
	}
	if request.UserID == "" {
		c.String(http.StatusBadRequest, "Invalid user_id")
		return
	}
	if !checkFileOwner(userID, request.FileID) {
		c.String(http.StatusUnauthorized, "You can only share your files")
		return
	}
	fileName := getFileSystemNameByFileID(userID, request.FileID)
	if fileName == "" {
		c.String(http.StatusBadRequest, "No file with that id")
		return
	}

	file, err := ioutil.ReadFile(service.Path + "/" + userID + "/" + fileName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something whent wrong while reading the file")
		return
	}
	decryptedfile, err := encryption.Decrypt(userKey, file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something whent wrong while decrypting the file")
		return
	}
	generatedFileName := strings.Replace(uuid.New().String(), "-", "", -1)
	out, err := os.Create(service.Path + "/" + request.UserID + "/" + generatedFileName)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "Failed to genetate new name for the file")
		return
	}
	key := getUserKey(request.UserID)
	encryptedFile, err := encryption.Encrypt(key, decryptedfile)
	defer out.Close()
	_, err = out.Write(encryptedFile)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to encrypt new name for the file")
		return
	}
	addFileToUser(request.UserID, generatedFileName, fileName)
	c.String(http.StatusCreated, "created")

	return
}
