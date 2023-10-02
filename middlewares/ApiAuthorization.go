package middlewares

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AuthorizeApiKey(c *gin.Context) bool {
	apiKey := c.GetHeader("Authorization")
	expectedAPIKey := os.Getenv("API_KEY")

	if apiKey != expectedAPIKey {
		return false
	}

	// Continue processing if the API key is valid
	return true
}