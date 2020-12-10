package controllers

import (
	"encoding/json"
	u "go-contacts/utils"
	"net/http"
)

var CreateAvailServer = func(w http.ResponseWriter, r *http.Request) {

	serverID := r.Context().Value("serverID").(string) //Grab the id of the serverID that send the request
	availServer := &models.AvailServer{}

	err := json.NewDecoder(r.Body).Decode(availServer)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	availServer.ServerID = serverID
	resp := availServer.Create()
	u.Respond(w, resp)
}
