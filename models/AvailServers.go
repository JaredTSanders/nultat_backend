package models

import (
	"fmt"

	u "github.com/JaredTSanders/nultat_backend/utils"

	"github.com/jinzhu/gorm"
)

type AvailServer struct {
	gorm.Model
	Name     string `json:"name"`
	ServerID string `json:"serverID"`
	UserId   uint
}

func (availServer *AvailServer) Validate() (map[string]interface{}, bool) {

	if availServer.Name == "" {
		return u.Message(false, "AvailServer name should be on the payload"), false
	}

	// if availServer.Phone == "" {
	// 	return u.Message(false, "Phone number should be on the payload"), false
	// }

	// if availServer.UserId <= 0 {
	// 	return u.Message(false, "User is not recognized"), false
	// }

	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (availServer *AvailServer) Create() map[string]interface{} {

	if resp, ok := availServer.Validate(); !ok {
		return resp
	}

	GetDB().Create(availServer)

	resp := u.Message(true, "success")
	resp["availServer"] = availServer
	return resp
}

func GetAvailServer(id uint) *AvailServer {

	availServer := &AvailServer{}
	err := GetDB().Table("availServers").Where("id = ?", id).First(availServer).Error
	if err != nil {
		return nil
	}
	return availServer
}

func GetAvailServers() []*AvailServer {

	availServer := make([]*AvailServer, 0)
	err := GetDB().Table("avail_servers").Find(&availServer).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return availServer
}
