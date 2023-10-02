package utilities

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Hashing(jsonStringToHash []byte) string {
	// Create a new SHA-256 hasher
	hasher := sha256.New()

	// Write the data from the request body to the hasher
	hasher.Write([]byte(jsonStringToHash))

	// Calculate the SHA-256 hash
	hashedData := hasher.Sum(nil)

	// Convert the hashed data to a hexadecimal string
	return hex.EncodeToString(hashedData)
}