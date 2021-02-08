package app

import (
	"github.com/grushka14/S4/app/users"
	"github.com/grushka14/S4/handlers"
)

func mapUrls() {
	router.POST("/v1/auth", users.LoginHandler)

	router.Use(users.TokenValidationMiddleware)

	router.GET("/v1/users", users.GetUsersHandler)

	router.GET("/v1/files", users.GetFilesHandler)
	router.PUT("/v1/files", users.PutFileHandler)
	router.POST("/v1/files/share", handlers.ShareFile)
	router.DELETE("/v1/files", users.DeleteFileHandler)

	router.GET("/v1/read", handlers.ReadFile)
}
