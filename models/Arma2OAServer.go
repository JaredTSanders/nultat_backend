package models

import (
	"fmt"

	u "github.com/JaredTSanders/nultat_backend/utils"

	"github.com/jinzhu/gorm"
)

type Arma2Server struct {
	gorm.Model
	Name   string `json:"name"`
	UserId uint
	Type   string `json: "type"`
}

func (arma2Server *Arma2Server) Validate() (map[string]interface{}, bool) {

	if arma2Server.Name == "" {
		return u.Message(false, "Arma2Server name should be on the payload"), false
	}

	// if arma2Server.Phone == "" {
	// 	return u.Message(false, "Phone number should be on the payload"), false
	// }

	// if arma2Server.UserId <= 0 {
	// 	return u.Message(false, "User is not recognized"), false
	// }

	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (arma2Server *Arma2Server) Create() map[string]interface{} {

	if resp, ok := arma2Server.Validate(); !ok {
		return resp
	}

	GetDB().Create(arma2Server)

	resp := u.Message(true, "success")
	resp["arma2Server"] = arma2Server
	return resp
}

func GetArma2Server(id uint) *Arma2Server {

	arma2Server := &Arma2Server{}
	err := GetDB().Table("arma2Servers").Where("id = ?", id).First(arma2Server).Error
	if err != nil {
		return nil
	}
	return arma2Server
}

func GetArma2Servers(user uint) []*Arma2Server {

	arma2Servers := make([]*Arma2Server, 0)
	err := GetDB().Table("arma2Servers").Where("user_id = ?", user).Find(&arma2Servers).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return arma2Servers
}
