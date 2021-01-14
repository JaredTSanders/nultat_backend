package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type BasicAccount struct {
	gorm.Model
	Email       string `json:"email"`
	FName       string `json:"first_name"`
	LName       string `json:"last_name"`
	Status      string `json:"status`
	Standing    string `json:"standing"`
	MFA_Enabled string `json:"mfa_enabled"`
	AccType     string `json:"account_type"`
}

func GetAllUsers() []*BasicAccount {

	account := make([]*BasicAccount, 0)
	err := GetDB().Table("accounts").Find(&account).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return account
}

func GetCurrentUser(email string) *BasicAccount {
	account := &BasicAccount{}
	fmt.Println("EMAIL HERE", email)
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		return nil
	}
	return account
}
