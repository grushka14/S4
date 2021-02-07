package users

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