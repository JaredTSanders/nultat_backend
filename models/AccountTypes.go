package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type AccountTypes struct {
	gorm.Model
	TypeName string `json:"type"`
	Game     string `json:"game"`
	Slots    string `json:"slots"`
	Memory   string `json:"ram"`
}

func GetAllAccountTypes() []*AccountTypes {

	account := make([]*AccountTypes, 0)
	err := GetDB().Table("account_types").Find(&account).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return account
}

func GetCurrentAccountType(id uint) *AccountTypes {
	account := &AccountTypes{}
	// err := GetDB().Table("account_types").Where("name = ?", name).First(account).Error
	err := GetDB().Exec("select account_types.type_name from account_types, accounts where accounts.id = ? and account_types.type_name = accounts.acc_type", id).First(account).Error
	if err != nil {
		return nil
	}
	return account
}
