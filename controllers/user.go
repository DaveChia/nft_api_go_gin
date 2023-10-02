package controllers

import (
	"encoding/json"
	"net/http"
	"nft_api_go_gin/database"
	"nft_api_go_gin/models"
	"nft_api_go_gin/utilities"
	"nft_api_go_gin/validators"

	"github.com/gin-gonic/gin"
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

//	Initiate the database connection and create the users table automatically if the table has not been created
func New() *UserRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.User{})
	return &UserRepo{Db: db}
}

//	Record the user's NRIC and wallet address when they mint an NFT
//	Each user's NRIC can only mint an NFT once
//	Each user's wallet address can only mint an NFT once
func (repository *UserRepo) CreateUser(c *gin.Context) {

	var user models.User

	//	Validate the request's form body
	//	Return an array of validation errors based on the field and error type
	//	Example:
	// [
    //     {
    //         "field": "Wallet",
    //         "error": "The Wallet field is required"
    //     },
    //     {
    //         "field": "Nric",
    //         "error": "The Nric field is required"
    //     }
    // ]
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validators.GenerateSplitValidatorErrorMessages(err)})
        return
	}

	var walletFoundCount int

	//	Check whether the wallet exists in the users, table, return error if it exists since each wallet can only mint an NFT once
	walletFoundResult := repository.Db.Raw("SELECT COUNT(*) FROM users WHERE wallet = ?", user.Wallet).Scan(&walletFoundCount)
	if walletFoundResult.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": walletFoundResult.Error})
		return
	} 
	if (walletFoundCount > 0) {
		c.AbortWithStatusJSON(http.StatusBadRequest,
		gin.H{
			"error": "The wallet exists in our database, please try with another wallet."})
		return
	}

	var nricFoundCount int

	//	Check whether the nric exists in the users, table, return error if it exists since each nric can only mint an NFT once
	nricFoundResult := repository.Db.Raw("SELECT COUNT(*) FROM users WHERE nric = ?", user.Nric).Scan(&nricFoundCount)
	if nricFoundResult.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": nricFoundResult.Error})
		return
	} 
	if (nricFoundCount > 0) {
		c.AbortWithStatusJSON(http.StatusBadRequest,
		gin.H{
			"error": "The nric exists in our database, please try with another nric."})
		return
	}

	//	Insert the row into the database
	err := models.CreateUser(repository.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	//	Generate the JSON object to be encrypted
	body := UserReceiptHashingBody{
		Wallet : user.Wallet,
		Nric : user.Nric,
	}

	//	Convert the JSON object into a json string
	jsonString, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//	Hash the json string into a sha256 hash as the minting transaction's unique receipt
	receipt := utilities.Sha256Hashing(jsonString)
	
	output :=  UserCreatedOutput{
		User: user,
		Receipt: receipt,
	}

	c.JSON(http.StatusCreated, output)
}