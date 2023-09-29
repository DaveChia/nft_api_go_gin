package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID    int 
	Wallet  string `json:"Wallet" binding:"required,max=255,alphanum"`
	Nric string `json:"Nric" binding:"required,max=10,alphanum"`
}

//create a user
func CreateUser(db *gorm.DB, User *User) (err error) {

	err = db.Create(User).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserByWallet(db *gorm.DB, User *User, wallet string) (err error) {
	
	err = db.Where("wallet = ?", wallet).Find(User).Error

	if err != nil {
		return err
	}
	return nil
}