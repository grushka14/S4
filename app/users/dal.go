package users

import (
	"strings"

	"github.com/grushka14/S4/app/db"
)

func resolveUserId(email string) string {
	for key := range db.UserData {
		if strings.EqualFold(db.UserData[key], email) {
			return key
		}
	}
	return ""
}
