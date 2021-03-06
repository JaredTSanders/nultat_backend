package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JaredTSanders/nultat_backend/models"
	u "github.com/JaredTSanders/nultat_backend/utils"
)

var CreateMinecraftBRServer = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	minecraftBRServer := &models.MinecraftBRServer{}

	err := json.NewDecoder(r.Body).Decode(minecraftBRServer)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	minecraftBRServer.UserId = user
	resp := minecraftBRServer.Create()
	u.Respond(w, resp)
}
