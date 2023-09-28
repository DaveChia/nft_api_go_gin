package controllers

import (
	"net/http"
	"nft_api_go_gin/database"
	"nft_api_go_gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepo struct {
	Db *gorm.DB
}

func New() *UserRepo {
	db := database.InitDb()
	db.AutoMigrate(&models.User{})
	return &UserRepo{Db: db}
}

//create user
func (repository *UserRepo) CreateUser(c *gin.Context) {

	var user models.User

	if err:=c.ShouldBindJSON(&user);err!=nil{
        c.AbortWithStatusJSON(http.StatusBadRequest,
        gin.H{
            "error": "VALIDATEERR-1",
            "message": "Invalid inputs. Please check your inputs"})
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
	c.JSON(http.StatusCreated, user)
}