package app

import "asaf_project/src/handlers"

func mapUrls() {
	router.GET("/v1/auth", handlers.Login)
	router.GET("/v1/users", handlers.GetUsers)
	router.PUT("/v1/files", handlers.PutFile)
	router.DELETE("/v1/files", handlers.DeleteFile)
	router.GET("/v1/files", handlers.GetFiles)
	router.GET("/v1/read", handlers.ReadFile)
	router.POST("/v1/files/share", handlers.ShareFile)
}