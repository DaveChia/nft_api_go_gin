package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID    int
	Name  string
	Nric string
}

//create a user
func CreateUser(db *gorm.DB, User *User) (err error) {
	err = db.Create(User).Error
	if err != nil {
		return err
	}
	return nil
}