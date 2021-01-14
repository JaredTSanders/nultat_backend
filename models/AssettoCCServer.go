package models

import (
	"fmt"

	u "github.com/JaredTSanders/nultat_backend/utils"

	"github.com/jinzhu/gorm"
)

type AssettoCCServer struct {
	gorm.Model
	Name   string `json:"name"`
	UserId uint
	Type   string `json: "type"`
}

func (assettoCCServer *AssettoCCServer) Validate() (map[string]interface{}, bool) {

	if assettoCCServer.Name == "" {
		return u.Message(false, "Name should be in the payload"), false
	}

	// if assettoCCServer.Phone == "" {
	// 	return u.Message(false, "Phone number should be on the payload"), false
	// }

	// if assettoCCServer.UserId <= 0 {
	// 	return u.Message(false, "User is not recognized"), false
	// }

	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (assettoCCServer *AssettoCCServer) Create() map[string]interface{} {

	if resp, ok := assettoCCServer.Validate(); !ok {
		return resp
	}

	GetDB().Create(assettoCCServer)

	resp := u.Message(true, "success")
	resp["assettoCCServer"] = assettoCCServer
	return resp
}

func GetAssettoCCServer(id uint) *AssettoCCServer {

	assettoCCServer := &AssettoCCServer{}
	err := GetDB().Table("assettoCCServers").Where("id = ?", id).First(assettoCCServer).Error
	if err != nil {
		return nil
	}
	return assettoCCServer
}

func GetAssettoCCServers(user uint) []*AssettoCCServer {

	assettoCCServers := make([]*AssettoCCServer, 0)
	err := GetDB().Table("assettoCCServers").Where("user_id = ?", user).Find(&assettoCCServers).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return assettoCCServers
}
