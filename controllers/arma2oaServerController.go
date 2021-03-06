package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JaredTSanders/nultat_backend/models"
	u "github.com/JaredTSanders/nultat_backend/utils"
)

var CreateArma2Server = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	arma2Server := &models.Arma2Server{}

	err := json.NewDecoder(r.Body).Decode(arma2Server)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	arma2Server.UserId = user
	resp := arma2Server.Create()
	u.Respond(w, resp)
}
