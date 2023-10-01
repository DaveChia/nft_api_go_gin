package main

import (
	"log"
	"net/http"
	"nft_api_go_gin/controllers"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
		
	router := gin.Default()
	
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Specify specific origins if needed
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}

	router.Use(cors.New(config))

	// Middleware to validate API key
	router.Use(func(c *gin.Context) {
		apiKey := c.GetHeader("Authorization")
		expectedAPIKey := os.Getenv("API_KEY")

		if apiKey != expectedAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Continue processing if the API key is valid
		c.Next()
	})

    userRepo := controllers.New()

    router.POST("/users", userRepo.CreateUser)

	router.Run(os.Getenv("API_URL"))
}

func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Replace * with your allowed origins
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Token")

        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
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