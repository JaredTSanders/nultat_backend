package controllers

import (
	"fmt"
	"net/http"

	"github.com/JaredTSanders/nultat_backend/models"
	u "github.com/JaredTSanders/nultat_backend/utils"
)

var GetAllUsers = func(w http.ResponseWriter, r *http.Request) {
	data := models.GetAllUsers()

	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var GetCurrentUser = func(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	sessionToken := c.Value
	res, err := models.Cache.Do("GET", sessionToken)
	if err != nil {
		resp := u.Message(false, "Error retrieving user account")
		u.Respond(w, resp)
	}
	email := fmt.Sprintf("%s", res)
	data := models.GetCurrentUser(email)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var UpdateUserProfile = func(w http.ResponseWriter, r *http.Request) {
	// id := r.Context().Value("user").(uint)
	// previous := models.GetCurrentUser(id)
	// // new := models.UpdateCurrentUser(*previous)
	// resp := u.Message(true, "success")
	// resp["data"] = new
	// u.Respond(w, resp)
}
