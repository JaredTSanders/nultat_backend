package controllers

import (
	"encoding/json"
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
)

var CreateArkPCServer = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	arkPCServer := &models.ArkPCServer{}

	err := json.NewDecoder(r.Body).Decode(arkPCServer)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	arkPCServer.UserId = user
	resp := arkPCServer.Create()
	u.Respond(w, resp)
}
