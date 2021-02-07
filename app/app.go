package app

import (
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApp() {
	mapUrls()
	router.Run(":8081")
}
