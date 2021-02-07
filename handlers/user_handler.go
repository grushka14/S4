package handlers

import (
	"asaf_project/src/app/utils/responses"
	"asaf_project/src/domain/token"
	"asaf_project/src/domain/user"
	"asaf_project/src/encryption"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

const basePah = "C:/Users/grush/go/src/asaf_project/src/files/"

// GetUsers function
func GetUsers(c *gin.Context) {
	if c.Request.Header["Authorization"] == nil {
		c.JSON(http.StatusUnauthorized, responses.New401("unauthorized"))
		return
	}
	t := token.Token(strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1])
	res, err := t.Validate()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}
	if res == "expired" {
		c.JSON(http.StatusUnauthorized, responses.New401("Token expired"))
	}

	if res == "unauthorized" {
		c.JSON(http.StatusUnauthorized, responses.New401("unauthorized"))
	}

	if res == "authorized" {
		responseData := []string{}
		for _, u := range user.Users {
			responseData = append(responseData, u.Email)
		}
		fmt.Println(responseData)
		c.JSON(http.StatusOK, responses.New200GetUsers(responseData))
		return
	}
	c.JSON(http.StatusInternalServerError, responses.New500("server error"))
}

// Login function
func Login(c *gin.Context) {
	var user user.User
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("Invalid JSON")
		c.JSON(http.StatusBadRequest, "")
		return
	}

	if user.CheckEmail() == true {
		token, err := token.CreateToken(user)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, "")
			return
		}
		c.JSON(http.StatusOK, responses.New200Login(token))
		return
	}

	c.JSON(http.StatusUnauthorized, "")
}

// PutFile function
func PutFile(c *gin.Context) {
	if c.Request.Header["Authorization"] == nil {
		c.JSON(http.StatusUnauthorized, responses.New401("unauthorized"))
		return
	}

	t := token.Token(strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1])
	res, err := t.Validate()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}
	if res == "expired" {
		c.JSON(http.StatusUnauthorized, responses.New401("Token expired"))
		return
	}

	if res == "unauthorized" {
		c.JSON(http.StatusUnauthorized, responses.New401("unauthorized"))
		return
	}

	email, ok := t.GetEmail()
	if !ok {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}

	userId, ok := user.GetUserId(email)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}

	key := user.GetUserKey(email)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "bad_request")
		return
	}

	filename := header.Filename
	generatedFileName := strings.Replace(uuid.New().String(), "-", "", -1)

	user.AddFileToUser(email, generatedFileName, filename)

	fmt.Println(generatedFileName)

	out, err := os.Create(basePah + userId + "/" + generatedFileName)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	byteFile, err := ioutil.ReadAll(file)
	encryptedFile, err := encryption.Encrypt(key, byteFile)

	defer out.Close()
	_, err = out.Write(encryptedFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusCreated, responses.New201("created"))
}

// DeleteFile function
func DeleteFile(c *gin.Context) {
	if c.Request.Header["Authorization"] == nil {
		c.JSON(http.StatusUnauthorized, responses.New401("unauthorized"))
		return
	}

	t := token.Token(strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1])
	res, err := t.Validate()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}
	if res == "expired" {
		c.JSON(http.StatusUnauthorized, responses.New401("Token expired"))
		return
	}

	if res == "unauthorized" {
		c.JSON(http.StatusUnauthorized, responses.New401("unauthorized"))
		return
	}

	email, ok := t.GetEmail()
	if !ok {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}

	userId, ok := user.GetUserId(email)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}

	fileId := c.Query("file-id")
	if fileId == "" {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}

	owner, fileName, exist := user.CheckFile(fileId)
	if !exist {
		c.JSON(http.StatusBadRequest, responses.New400("no file with that name"))
		return
	}

	if owner != email {
		c.JSON(http.StatusBadRequest, responses.New401("this file does not belong to you"))

	}

	os.Remove(basePah + userId + "/" + fileName)

	fmt.Println(fileId)
	c.JSON(http.StatusBadRequest, responses.New204("deleted"))
}

// GetFiles function
func GetFiles(c *gin.Context) {
	if c.Request.Header["Authorization"] == nil {
		c.JSON(http.StatusUnauthorized, responses.New401("unauthorized"))
		return
	}

	t := token.Token(strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1])
	res, err := t.Validate()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}
	if res == "expired" {
		c.JSON(http.StatusUnauthorized, responses.New401("Token expired"))
		return
	}

	if res == "unauthorized" {
		c.JSON(http.StatusUnauthorized, responses.New401("unauthorized"))
		return
	}

	email, ok := t.GetEmail()
	if !ok {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}

	files := user.GetFiles(email)
	c.JSON(http.StatusOK, responses.New200GetFiles(files))
}

// ReadFile function
func ReadFile(c *gin.Context) {
	if c.Request.Header["Authorization"] == nil {
		c.JSON(http.StatusUnauthorized, responses.New401("unauthorized"))
		return
	}

	t := token.Token(strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1])
	res, err := t.Validate()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}
	if res == "expired" {
		c.JSON(http.StatusUnauthorized, responses.New401("Token expired"))
		return
	}

	if res == "unauthorized" {
		c.JSON(http.StatusUnauthorized, responses.New401("unauthorized"))
		return
	}

	email, ok := t.GetEmail()
	if !ok {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}

	userId, ok := user.GetUserId(email)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}

	fileId := c.Query("file-id")
	if fileId == "" {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}

	owner, fileName, exist := user.CheckFile(fileId)
	if !exist {
		c.JSON(http.StatusBadRequest, responses.New400("no file with that id"))
		return
	}

	if owner != email {
		c.JSON(http.StatusBadRequest, responses.New401("this file does not belong to you"))
		return
	}

	file, err := ioutil.ReadFile(basePah + userId + "/" + fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}

	key := user.GetUserKey(email)
	decryptedfile, err := encryption.Decrypt(key, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}
	c.JSON(http.StatusOK, responses.New200ReadFile(decryptedfile))
}

// ShareFile function
func ShareFile(c *gin.Context) {
	if c.Request.Header["Authorization"] == nil {
		c.JSON(http.StatusUnauthorized, responses.New401("unauthorized"))
		return
	}

	t := token.Token(strings.Split(c.Request.Header["Authorization"][0], "Bearer ")[1])
	res, err := t.Validate()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}
	if res == "expired" {
		c.JSON(http.StatusUnauthorized, responses.New401("Token expired"))
		return
	}

	if res == "unauthorized" {
		c.JSON(http.StatusUnauthorized, responses.New401("unauthorized"))
		return
	}

	var data user.ShareUser
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		fmt.Println("Invalid JSON")
		c.JSON(http.StatusBadRequest, "bad_request")
		return
	}

	email, ok := t.GetEmail()
	if !ok {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}

	myId, _ := user.GetUserId(email)
	if myId == data.UserId {
		c.JSON(http.StatusBadRequest, responses.New400("can not share file with your selfe"))
		return
	}

	owner, fileName, exist := user.CheckFile(data.FileId)
	if !exist {
		c.JSON(http.StatusBadRequest, responses.New400("no file with that id"))
		return
	}

	if owner != email {
		c.JSON(http.StatusBadRequest, responses.New401("this file does not belong to you"))
		return
	}

	file, err := ioutil.ReadFile(basePah + myId + "/" + fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}

	key := user.GetUserKey(email)
	decryptedfile, err := encryption.Decrypt(key, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.New400("bad_request"))
		return
	}

	generatedFileName := strings.Replace(uuid.New().String(), "-", "", -1)
	out, err := os.Create(basePah + data.UserId + "/" + generatedFileName)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	newUserEmail, _ := user.GetUserEmail(data.UserId)
	newKey := user.GetUserKey(newUserEmail)
	encryptedFile, err := encryption.Encrypt(newKey, decryptedfile)
	defer out.Close()
	_, err = out.Write(encryptedFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	name := user.FileSystemNameToName(fileName)
	if name == "" {
		name = fileName
	}
	user.AddFileToUser(newUserEmail, generatedFileName, name)
	c.JSON(http.StatusCreated, responses.New201("created"))
}
