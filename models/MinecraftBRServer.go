package models

import (
	"fmt"

	u "github.com/JaredTSanders/nultat_backend/utils"

	"github.com/jinzhu/gorm"
)

type MinecraftBRServer struct {
	gorm.Model
	Name   string `json:"name"`
	UserId uint
}

func (minecraftBRServer *MinecraftBRServer) Validate() (map[string]interface{}, bool) {

	if minecraftBRServer.Name == "" {
		return u.Message(false, "MinecraftBRServer name should be on the payload"), false
	}

	// if minecraftBRServer.Phone == "" {
	// 	return u.Message(false, "Phone number should be on the payload"), false
	// }

	// if minecraftBRServer.UserId <= 0 {
	// 	return u.Message(false, "User is not recognized"), false
	// }

	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (minecraftBRServer *MinecraftBRServer) Create() map[string]interface{} {

	if resp, ok := minecraftBRServer.Validate(); !ok {
		return resp
	}

	GetDB().Create(minecraftBRServer)

	resp := u.Message(true, "success")
	resp["minecraftBRServer"] = minecraftBRServer
	return resp
}

func GetMinecraftBRServer(id uint) *MinecraftBRServer {

	minecraftBRServer := &MinecraftBRServer{}
	err := GetDB().Table("minecraftBRServers").Where("id = ?", id).First(minecraftBRServer).Error
	if err != nil {
		return nil
	}
	return minecraftBRServer
}

func GetMinecraftBRServers(user uint) []*MinecraftBRServer {

	minecraftBRServers := make([]*MinecraftBRServer, 0)
	err := GetDB().Table("minecraftBRServers").Where("user_id = ?", user).Find(&minecraftBRServers).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return minecraftBRServers
}
