package controllers

import (
	"encoding/json"
	// "fmt"
	// "go/token"
	"net/http"
	"fmt"
	// "io/ioutil"

	"github.com/JaredTSanders/nultat_backend/models"
	u "github.com/JaredTSanders/nultat_backend/utils"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create() //Create account
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	
	// b, e := ioutil.ReadAll(r.Body)
	// if e != nil {
	//     panic(e)
	// }

	dbUri := fmt.Sprintf("host=%s url=%s proto=%s", r.Method, r.URL, r.Proto)
	fmt.Println(dbUri)
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Email, account.Password, w)
	u.Respond(w, resp)
}

var RefreshToken = func(w http.ResponseWriter, r *http.Request) {
	// id := r.Context().Value("user").(uint)
	// tokenData := models.RefreshToken(id)
	// resp := u.Message(true, "success")
	// resp["token"] = tokenData
	// u.Respond(w, resp)
	models.Refresh(w, r)
	resp := u.Message(true, "Token Refreshed")
	u.Respond(w, resp)
}

// var GetCurrentUser = func(w http.ResponseWriter, r *http.Request) {
// 	id := r.Context().Value("user").(uint)
// 	data := models.GetCurrentUser(id)
// 	resp := u.Message(true, "success")
// 	resp["data"] = data
// 	u.Respond(w, resp)
// }

var Logout = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

var GetUserLoginStatus = func(w http.ResponseWriter, r *http.Request) {
	u.Respond(w, u.Message(true, "logged in"))
	return
}
