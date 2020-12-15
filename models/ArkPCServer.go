package models

import (
	"fmt"

	u "github.com/JaredTSanders/nultat_backend/utils"

	"github.com/jinzhu/gorm"
)

type ArkPCServer struct {
	gorm.Model
	Name   string `json:"name"`
	UserId uint
}

func (arkPCServer *ArkPCServer) Validate() (map[string]interface{}, bool) {

	if arkPCServer.Name == "" {
		return u.Message(false, "ArkPCServer name should be on the payload"), false
	}

	// if arkPCServer.Phone == "" {
	// 	return u.Message(false, "Phone number should be on the payload"), false
	// }

	// if arkPCServer.UserId <= 0 {
	// 	return u.Message(false, "User is not recognized"), false
	// }

	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (arkPCServer *ArkPCServer) Create() map[string]interface{} {

	if resp, ok := arkPCServer.Validate(); !ok {
		return resp
	}

	GetDB().Create(arkPCServer)

	resp := u.Message(true, "success")
	resp["arkPCServer"] = arkPCServer
	return resp
}

func GetArkPCServer(id uint) *ArkPCServer {

	arkPCServer := &ArkPCServer{}
	err := GetDB().Table("arkPCServers").Where("id = ?", id).First(arkPCServer).Error
	if err != nil {
		return nil
	}
	return arkPCServer
}

func GetArkPCServers(user uint) []*ArkPCServer {

	arkPCServers := make([]*ArkPCServer, 0)
	err := GetDB().Table("arkPCServers").Where("user_id = ?", user).Find(&arkPCServers).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return arkPCServers
}
