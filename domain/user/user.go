package user

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

//User stract
type User struct {
	ID     string     `json:"id"`
	Email  string     `json:"email"`
	Files  []UserFile `json:"files"`
	Secret []byte     `json:"secret"`
}

//UserFile struct
type UserFile struct {
	FileID         string `json:"id"`
	FileSystemName string `json:"file"`
	FileUserName   string `json:"file_name"`
}

//ShareUser struct
type ShareUser struct {
	UserID string `json:"user"`
	FileID string `json:"file"`
}

//Users global
var Users = []User{
	{
		ID:     "4512c320-7308-4cb4-86c9-0e9a5feb2663",
		Email:  "asafo@gmail.com",
		Files:  []UserFile{},
		Secret: []byte{116, 149, 214, 111, 240, 238, 238, 133, 152, 252, 79, 141, 156, 12, 52, 48, 21, 13, 85, 137, 204, 78, 186, 227, 82, 121, 97, 202, 68, 189, 170, 196},
	},
	{
		ID:     "ee897b51-e797-48a4-bbbd-d4b4dd3914b9",
		Email:  "bengi@gmail.com",
		Files:  []UserFile{},
		Secret: []byte{238, 130, 152, 252, 79, 141, 189, 149, 214, 11, 240, 68, 189, 170, 106, 238, 156, 121, 52, 182, 11, 97, 202, 48, 21, 19, 85, 137, 204, 78, 186, 250},
	},
	{
		ID:     "a990e698-53e9-4ebc-ae3a-d4f07627aad1",
		Email:  "grushka@gmail.com",
		Files:  []UserFile{},
		Secret: []byte{240, 238, 252, 79, 238, 130, 152, 141, 156, 121, 52, 189, 149, 214, 11, 48, 21, 19, 85, 137, 204, 178, 18, 227, 82, 121, 96, 202, 68, 189, 170, 175},
	},
}

//CheckEmail method
func (u User) CheckEmail() bool {
	for i := range Users {
		if Users[i].Email == u.Email {
			return true
		}
	}
	return false
}

//GetUserID function
func GetUserID(e string) (string, bool) {
	userID := ""
	for _, user := range Users {
		if user.Email == e {
			userID = user.ID
		}
		if userID != "" {
			return userID, true
		}
	}
	return "", false
}

// AddFileToUser function
func AddFileToUser(email string, fsn string, fun string) {
	for i, user := range Users {
		if user.Email == email {
			var newFile = UserFile{
				FileID:         strings.Replace(uuid.New().String(), "-", "", -1),
				FileSystemName: fsn,
				FileUserName:   fun,
			}
			Users[i].Files = append(user.Files, newFile)
			return
		}
	}
}

//GetUserKey function
func GetUserKey(e string) []byte {

	key := []byte{}
	for _, user := range Users {
		if user.Email == e {
			key = user.Secret
		}
	}
	return key
}

// type fileRes struct {
// 	email string
// 	systemFileName string
// }

//CheckFile function
func CheckFile(f string) (string, string, bool) {
	fmt.Println(Users)
	for _, user := range Users {
		// if this is the right user
		for _, file := range user.Files {
			if file.FileID == f {
				return user.Email, file.FileSystemName, true
			}
		}
	}
	return "", "", false
}

//FileSystemNameToName function
func FileSystemNameToName(fsn string) string {
	for _, user := range Users {
		for _, file := range user.Files {
			if file.FileSystemName == fsn {
				return file.FileUserName
			}
		}
	}
	return ""
}

//GetFiles function
func GetFiles(e string) []map[string]string {
	files := []map[string]string{}
	for _, user := range Users {
		if user.Email == e {
			for _, file := range user.Files {
				newFile := map[string]string{
					"id":   file.FileID,
					"name": file.FileUserName,
				}
				files = append(files, newFile)
			}
		}
	}
	return files
}

//GetUserEmail function
func GetUserEmail(i string) (string, bool) {
	userEmail := ""
	for _, user := range Users {
		if user.ID == i {
			userEmail = user.Email
		}
		if userEmail != "" {
			return userEmail, true
		}
	}
	return "", false
}

// RemoveFileFromUserMemory func
func RemoveFileFromUserMemory(e string, f string) (bool, error) {
	for i, u := range Users {
		if u.Email == e {
			for j, file := range u.Files {
				if file.FileID == f {
					// removing index j
					Users[i].Files = append(u.Files[:j], u.Files[j+1:]...)
					return true, nil
				}
			}
		}
	}
	return false, nil
}
