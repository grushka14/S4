package users

import (
	"strings"

	"github.com/google/uuid"
)

func resolveUserID(email string) string {
	for key := range UserData {
		if strings.EqualFold(UserData[key].Email, email) {
			return key
		}
	}
	return ""
}

func getAllUsers() []getUsersResponse {
	var allUsers = make([]getUsersResponse, len(UserData))
	i := 0
	for key := range UserData {
		allUsers[i] = getUsersResponse{ID: key, Email: UserData[key].Email}
		i++
	}
	return allUsers
}

func dbGetUserSecretByID(userID string) []byte {
	return UserData[userID].Secret
}

func dbAddFileToUser(userID string, systemFileName string, userFileName string) {
	newslice := make([]UserFile, len(UserData[userID].Files)+1)
	newslice = append(UserData[userID].Files, UserFile{
		FileID:         (uuid.New().String()),
		FileSystemName: systemFileName,
		FileUserName:   userFileName,
	})
	tempUser := UserData[userID]
	tempUser.Files = newslice
	UserData[userID] = tempUser
}

func dbcheckFileOwner(userID string, fileID string) bool {
	for i := 0; i < len(UserData[userID].Files); i++ {
		if strings.EqualFold(fileID, UserData[userID].Files[i].FileID) {
			return true
		}
	}
	return false
}

func dbGetFileSystemNameByFileID(userID string, fileID string) string {
	for i := 0; i < len(UserData[userID].Files); i++ {
		if strings.EqualFold(fileID, UserData[userID].Files[i].FileID) {
			return UserData[userID].Files[i].FileSystemName
		}
	}
	return ""
}

func dbRemoveFileFromUser(userID string, fileID string) {
	tempUser := UserData[userID]
	for i := 0; i < len(tempUser.Files); i++ {
		if strings.EqualFold(fileID, tempUser.Files[i].FileID) {
			tempUser.Files = append(tempUser.Files[:i], tempUser.Files[i+1:]...)
			UserData[userID] = tempUser
			return
		}
	}
}

func dbGetFilesByUserID(userID string) []getFilesResponse {
	var res = make([]getFilesResponse, len(UserData[userID].Files))
	for i := 0; i < len(UserData[userID].Files); i++ {
		res[i] = getFilesResponse{
			ID:       UserData[userID].Files[i].FileID,
			FileName: UserData[userID].Files[i].FileUserName,
		}
	}
	return res
}
