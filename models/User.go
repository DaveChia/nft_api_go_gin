package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID    int 
	Wallet  string `gorm:"type:varchar(255);not null;unique;" json:"wallet" binding:"required,max=255,min=8,alphanum"`
	Nric string `gorm:"type:varchar(255);not null;unique;" json:"nric" binding:"required,max=10,min=8,alphanum"`
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