package db

import (
	"github.com/grushka14/S4/app/users"
)

//I dont have a real db there for i will use this as a mock db.

// UserData this is a state of the data base menaged in memmory.
var UserData = map[string]users.User{
	"4512c320-7308-4cb4-86c9-0e9a5feb2663": {
		Email:  "asafo@gmail.com",
		Files:  []users.UserFile{},
		Secret: []byte{116, 149, 214, 111, 240, 238, 238, 133, 152, 252, 79, 141, 156, 12, 52, 48, 21, 13, 85, 137, 204, 78, 186, 227, 82, 121, 97, 202, 68, 189, 170, 196},
	},
	"ee897b51-e797-48a4-bbbd-d4b4dd3914b9": {
		Email:  "bengi@gmail.com",
		Files:  []users.UserFile{},
		Secret: []byte{238, 130, 152, 252, 79, 141, 189, 149, 214, 11, 240, 68, 189, 170, 106, 238, 156, 121, 52, 182, 11, 97, 202, 48, 21, 19, 85, 137, 204, 78, 186, 250},
	},
	"a990e698-53e9-4ebc-ae3a-d4f07627aad1": {
		Email:  "grushka@gmail.com",
		Files:  []users.UserFile{},
		Secret: []byte{240, 238, 252, 79, 238, 130, 152, 141, 156, 121, 52, 189, 149, 214, 11, 48, 21, 19, 85, 137, 204, 178, 18, 227, 82, 121, 96, 202, 68, 189, 170, 175},
	},
}
