package main

import (
	"nft_api_go_gin/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

    userRepo := controllers.New()

    router.POST("/users", userRepo.CreateUser)

	router.Run("localhost:8080")
}