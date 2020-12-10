package models

import (
	"fmt"
	u "go-contacts/utils"

	"github.com/jinzhu/gorm"
)

type AssettoCCServer struct {
	gorm.Model
	UserId    uint
	UserEmail string
	UID       string
}

func (assettoCCServer *AssettoCCServer) Validate() (map[string]interface{}, bool) {

	if assettoCCServer.UserId == 0 {
		return u.Message(false, "User ID should be in the payload"), false
	}

	if assettoCCServer.UserEmail == "" {
		return u.Message(false, "User Email name should be in the payload"), false
	}

	if assettoCCServer.UID == "" {
		return u.Message(false, "Server UID must be generated and sent in the payload"), false
	}

	if len(assettoCCServer.UID) < 36 || len(assettoCCServer.UID) > 36 {
		return u.Message(false, "Invalid UID. Must be exactly 36 characters in length and formatted correctly"), false
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
