package app

import (
	"github.com/grushka14/S4/app/users"
)

func mapUrls() {
	router.POST("/v1/auth", users.LoginHandler)

	router.Use(users.TokenValidationMiddleware)

	router.GET("/v1/users", users.GetUsersHandler)

	router.GET("/v1/files", users.GetFilesHandler)
	router.PUT("/v1/files", users.PutFileHandler)
	router.POST("/v1/files/share", users.PostShareFile)
	router.DELETE("/v1/files", users.DeleteFileHandler)

	router.GET("/v1/read", users.GetReadFileHandler)
}
