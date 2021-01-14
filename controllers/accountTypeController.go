package controllers

import (
	"net/http"

	"github.com/JaredTSanders/nultat_backend/models"
	u "github.com/JaredTSanders/nultat_backend/utils"
)

var GetAllAccountTypes = func(w http.ResponseWriter, r *http.Request) {
	data := models.GetAllAccountTypes()
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var GetCurrentAccountType = func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user").(uint)
	data := models.GetCurrentAccountType(id)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
