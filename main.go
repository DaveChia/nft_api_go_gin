package main

import (
	"nft_api_go_gin/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

    userRepo := controllers.New()

    router.POST("/users", userRepo.CreateUser)

	router.Run("localhost:8080")
}