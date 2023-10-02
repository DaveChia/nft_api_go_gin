package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


func SetupCorsConfig() cors.Config {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Specify specific origins if needed
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}

	return config;
}

func SetupCors() gin.HandlerFunc {
	corsConnection := cors.New(SetupCorsConfig())

	return corsConnection
}