package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JaredTSanders/nultat_backend/models"
	u "github.com/JaredTSanders/nultat_backend/utils"
)

var CreateAvailServer = func(w http.ResponseWriter, r *http.Request) {

	// user := r.Context().Value("email").(string) //Grab the id of the user that send the request
	arma2Server := &models.AvailServer{}

	err := json.NewDecoder(r.Body).Decode(arma2Server)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
	c, err := r.Cookie("session_token")
	sessionToken := c.Value
	res, err := models.Cache.Do("GET", sessionToken)
	if err != nil {
		resp := u.Message(false, "Error retrieving user account")
		u.Respond(w, resp)
	}
	email := fmt.Sprintf("%s", res)
	data := models.GetUserByEmail(email)

	arma2Server.UserId = data.ID
	resp := arma2Server.Create()
	u.Respond(w, resp)
}

var GetAllAvailServers = func(w http.ResponseWriter, r *http.Request) {
	data := models.GetAvailServers()
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
