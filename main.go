package main

import (
	"nft_api_go_gin/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/register", controllers.RegisterUser)
	router.Run("localhost:8080")
}