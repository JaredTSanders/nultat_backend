package controllers

import (
	"encoding/json"
	"net/http"
	// "github.com/gorilla/mux"
	"fmt"

	"github.com/JaredTSanders/nultat_backend/models"
	u "github.com/JaredTSanders/nultat_backend/utils"

)

var GetPodLogs = func(w http.ResponseWriter, r *http.Request) {
	pod := &models.Pod{}

	err := json.NewDecoder(r.Body).Decode(pod)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
    c, err := r.Cookie("session_token")
	sessionToken := c.Value
	res, error := models.Cache.Do("GET", sessionToken)
	if error != nil {
		resp := u.Message(false, "Error retrieving user account")
		u.Respond(w, resp)
	}

	email := fmt.Sprintf("%s", res)

	data := models.GetUserByEmail(email)
	pod.UserId = data.ID
	pod.Namespace = data.Namespace

	// fmt.Println(pod.Namespace)
	// fmt.Println(arkPCServer.UserId)
	// fmt.Println(arkPCServer.Namespace)

	pod.GetPodLogs()
	// fmt.Println(logs)
	u.Respond(w, u.Message(true, "Retrieved logs"))
}

var SendCommand = func(w http.ResponseWriter, r *http.Request) {

	pod := &models.Pod{}

	err := json.NewDecoder(r.Body).Decode(pod)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}
    c, err := r.Cookie("session_token")
	sessionToken := c.Value
	res, error := models.Cache.Do("GET", sessionToken)
	if error != nil {
		resp := u.Message(false, "Error retrieving user account")
		u.Respond(w, resp)
	}

	email := fmt.Sprintf("%s", res)

	data := models.GetUserByEmail(email)
	pod.UserId = data.ID
	pod.Namespace = data.Namespace

	// fmt.Println(pod.Namespace)
	// fmt.Println(arkPCServer.UserId)
	// fmt.Println(arkPCServer.Namespace)

	pod.SendCommand()
	// fmt.Println(logs)
	u.Respond(w, u.Message(true, "Retrieved logs"))


}
