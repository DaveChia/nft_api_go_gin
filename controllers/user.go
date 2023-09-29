package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"nft_api_go_gin/database"
	"nft_api_go_gin/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserRepo struct {
	Db *gorm.DB
}

type UserCreatedOutput struct {
	User    models.User `json:"user"`
	Receipt string `json:"receipt"`
}

type UserReceiptHashingBody struct {
	Wallet string `json:"wallet"`
	Nric string `json:"nric"`
}

//create

func New() *UserRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.User{})
	return &UserRepo{Db: db}
}

type ApiError struct {
    Field string `json:"field"`
    Error   string `json:"error"`
}

func msgForTag(tag string, field string, parameter string) string {
    switch tag {
		case "required":
			return "The " + field + " field is required"
		case "max":
			return "The " + field + " field may not be greater than " + parameter + " characters."
		case "min":
			return "The " + field + " field must be at least " + parameter + " characters."
		case "alphanum":
			return "The " + field + " field may only contain letters or numbers."
		}
    return "The " + field + " field is invalid."
}

//create user
func (repository *UserRepo) CreateUser(c *gin.Context) {

	var user models.User

	if err := c.ShouldBind(&user); err != nil {
		var ve validator.ValidationErrors
        if errors.As(err, &ve) {
            out := make([]ApiError, len(ve))
            for i, fe := range ve {
				fmt.Println(fe.Param())
                out[i] = ApiError{fe.Field(), msgForTag(fe.Tag(), fe.Field(), fe.Param())}
            }
            c.JSON(http.StatusBadRequest, gin.H{"errors": out})
        }
        return
	}

	var walletFoundCount int

	walletFoundResult := repository.Db.Raw("SELECT COUNT(*) FROM users WHERE wallet = ?", user.Wallet).Scan(&walletFoundCount)

	if walletFoundResult.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": walletFoundResult.Error})
		return
	} 

	if (walletFoundCount > 0) {
		c.AbortWithStatusJSON(http.StatusBadRequest,
		gin.H{
			"error": "EXISTWALLET-1",
			"message": "The wallet exists in our database, please try with another wallet."})
		return
	}

	var nricFoundCount int

	nricFoundResult := repository.Db.Raw("SELECT COUNT(*) FROM users WHERE nric = ?", user.Nric).Scan(&nricFoundCount)

	if nricFoundResult.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": nricFoundResult.Error})
		return
	} 

	if (nricFoundCount > 0) {
		c.AbortWithStatusJSON(http.StatusBadRequest,
		gin.H{
			"error": "EXISTNRIC-1",
			"message": "The nric exists in our database, please try with another nric."})
		return
	}

	err := models.CreateUser(repository.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	body := UserReceiptHashingBody{
		Wallet : user.Wallet,
		Nric : user.Nric,
	}

	jsonString, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a new SHA-256 hasher
	hasher := sha256.New()

	// Write the data from the request body to the hasher
	hasher.Write([]byte(jsonString))

	// Calculate the SHA-256 hash
	hashedData := hasher.Sum(nil)

	// Convert the hashed data to a hexadecimal string
	receipt := hex.EncodeToString(hashedData)
	
	output :=  UserCreatedOutput{
		User: user,
		Receipt: receipt,
	}

	c.JSON(http.StatusCreated, output)
}