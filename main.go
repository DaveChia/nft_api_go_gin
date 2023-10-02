package main

import (
	"log"
	"net/http"
	"nft_api_go_gin/controllers"
	"nft_api_go_gin/middlewares"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
		
	router := gin.Default()

	//	Middleware to setup CORS
	router.Use(middlewares.SetupCors())

	// Middleware to validate API key
	router.Use(func(c *gin.Context) {

		if (middlewares.AuthorizeApiKey(c) == false) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Continue processing if the API key is valid
		c.Next()
	})

    userRepo := controllers.New()

    router.POST("/users", userRepo.CreateUser)

	router.Run(os.Getenv("API_URL") +":" + os.Getenv("API_PORT"))
}

// func generateSecureAPIKey(keyLength int) (string, error) {
//     // Generate random bytes
//     keyBytes := make([]byte, keyLength)
//     _, err := rand.Read(keyBytes)
//     if err != nil {
//         return "", err
//     }

//     // Encode the random bytes as a base64 string
//     apiKey := base64.URLEncoding.EncodeToString(keyBytes)
//     return apiKey, nil
// }