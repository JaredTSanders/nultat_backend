package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("api/server/new", controllers.CreateAvailServer).Methods("POST")
	// router.HandleFunc("/api/server/new/server", controllers.GetServerValues).Methods("GET")
	router.HandleFunc("/api/server/new/arma2", controllers.CreateArma2Server).Methods("POST")
	router.HandleFunc("/api/server/new/arkpc", controllers.CreateArkPCServer).Methods("POST")
	router.HandleFunc("/api/server/new/assettocc", controllers.CreateAssettoCCServer).Methods("POST")
	router.HandleFunc("/api/server/new/minecraftbr", controllers.CreateMinecraftBRServer).Methods("POST")
	// router.HandleFunc("/api/server/new/minecraftbr", controllers.CreateMinecraftVanillaServer).Methods("POST")

	// router.HandleFunc("/api/me/", controllers.GetUserInfo).Methods("GET")

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
