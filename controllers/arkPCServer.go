package controllers

import (
	"encoding/json"
	"net/http"
	// "github.com/gorilla/mux"
	"fmt"

	"github.com/JaredTSanders/nultat_backend/models"
	u "github.com/JaredTSanders/nultat_backend/utils"

)

//	router.HandleFunc("/api/me/arkpc/new/{name}/{map}/{spass}/{apass}/{backup}/{update}")
var CreateArkPCServer = func(w http.ResponseWriter, r *http.Request) {

	arkPCServer := &models.ArkPCServer{}

	err := json.NewDecoder(r.Body).Decode(arkPCServer)
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

	arkPCServer.UserId = data.ID
	resp := arkPCServer.Create()
	u.Respond(w, resp)
}
