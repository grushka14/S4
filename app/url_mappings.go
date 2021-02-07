package app

import "github.com/grushka14/S4/handlers"

func mapUrls() {
	router.GET("/v1/auth", handlers.Login)
	router.GET("/v1/users", handlers.GetUsers)

	router.GET("/v1/files", handlers.GetFiles)
	router.PUT("/v1/files", handlers.PutFile)
	router.POST("/v1/files/share", handlers.ShareFile)
	router.DELETE("/v1/files", handlers.DeleteFile)

	router.GET("/v1/read", handlers.ReadFile)
}
