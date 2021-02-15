package users

import (
	"github.com/google/uuid"
)

//I dont have a real db there for i will use this as a mock db.

// UserData this is a state of the data base menaged in memmory.
var UserData = map[string]User{
	"4512c320-7308-4cb4-86c9-0e9a5feb2663": {
		Email:  "asafo@gmail.com",
		Files:  []UserFile{},
		Secret: []byte{116, 149, 214, 111, 240, 238, 238, 133, 152, 252, 79, 141, 156, 12, 52, 48, 21, 13, 85, 137, 204, 78, 186, 227, 82, 121, 97, 202, 68, 189, 170, 196},
	},
	"ee897b51-e797-48a4-bbbd-d4b4dd3914b9": {
		Email:  "bengi@gmail.com",
		Files:  []UserFile{},
		Secret: []byte{238, 130, 152, 252, 79, 141, 189, 149, 214, 11, 240, 68, 189, 170, 106, 238, 156, 121, 52, 182, 11, 97, 202, 48, 21, 19, 85, 137, 204, 78, 186, 250},
	},
	"a990e698-53e9-4ebc-ae3a-d4f07627aad1": {
		Email:  "grushka@gmail.com",
		Files:  []UserFile{},
		Secret: []byte{240, 238, 252, 79, 238, 130, 152, 141, 156, 121, 52, 189, 149, 214, 11, 48, 21, 19, 85, 137, 204, 178, 18, 227, 82, 121, 96, 202, 68, 189, 170, 175},
	},
}

type dataBase struct {
	Data map[string]User
}

func (db *dataBase) resolveUserID(email string) string {
	for key := range db.Data {
		if db.Data[key].Email == email {
			return key
		}
	}
	return ""
}

func (db *dataBase) getAllUsers() []getUsersResponse {
	var allUsers = make([]getUsersResponse, len(db.Data))
	i := 0
	for key := range db.Data {
		allUsers[i] = getUsersResponse{ID: key, Email: db.Data[key].Email}
		i++
	}
	return allUsers
}

func (db *dataBase) dbGetUserSecretByID(userID string) []byte {
	return db.Data[userID].Secret
}

func (db *dataBase) dbAddFileToUser(userID string, systemFileName string, userFileName string) {
	newslice := make([]UserFile, len(db.Data[userID].Files)+1)
	newslice = append(db.Data[userID].Files, UserFile{
		FileID:         (uuid.New().String()),
		FileSystemName: systemFileName,
		FileUserName:   userFileName,
	})
	tempUser := db.Data[userID]
	tempUser.Files = newslice
	db.Data[userID] = tempUser
}

func (db *dataBase) dbcheckFileOwner(userID string, fileID string) bool {
	for i := 0; i < len(db.Data[userID].Files); i++ {
		if fileID == db.Data[userID].Files[i].FileID {
			return true
		}
	}
	return false
}

func (db *dataBase) dbGetFileSystemNameByFileID(userID string, fileID string) string {
	for i := 0; i < len(db.Data[userID].Files); i++ {
		if fileID == db.Data[userID].Files[i].FileID {
			return db.Data[userID].Files[i].FileSystemName
		}
	}
	return ""
}

func (db *dataBase) dbRemoveFileFromUser(userID string, fileID string) {
	tempUser := db.Data[userID]
	for i := 0; i < len(tempUser.Files); i++ {
		if fileID == tempUser.Files[i].FileID {
			tempUser.Files = append(tempUser.Files[:i], tempUser.Files[i+1:]...)
			db.Data[userID] = tempUser
			return
		}
	}
}

func (db *dataBase) dbGetFilesByUserID(userID string) []getFilesResponse {
	var res = make([]getFilesResponse, len(db.Data[userID].Files))
	for i := 0; i < len(db.Data[userID].Files); i++ {
		res[i] = getFilesResponse{
			ID:       db.Data[userID].Files[i].FileID,
			FileName: db.Data[userID].Files[i].FileUserName,
		}
	}
	return res
}

//DB is a struct that hold the data for us.
var DB = dataBase{
	Data: UserData,
}
