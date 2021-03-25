package main

import (
	"log"
	"net/http"
	"os"

    "github.com/honeycombio/beeline-go/wrappers/hnygorilla"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
	"github.com/JaredTSanders/nultat_backend/app"
	"github.com/JaredTSanders/nultat_backend/controllers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	/*
	   =================
	    Router initialization:
	   =================
	*/
	router := mux.NewRouter()

	router.Use(hnygorilla.Middleware)

	/*
	   =================
	    GET endpoints:
	   =================
	*/

	router.HandleFunc("/api/user/logout", controllers.Logout).Methods("GET")
	router.HandleFunc("/api/server/get/avail_server", controllers.GetAllAvailServers).Methods("GET")
	router.HandleFunc("/api/user/get_all", controllers.GetAllUsers).Methods("GET")
	router.HandleFunc("/api/user/me", controllers.GetCurrentUser).Methods("GET")
	router.HandleFunc("/api/user/refresh", controllers.RefreshToken).Methods("GET")
	router.HandleFunc("/api/types/me", controllers.GetCurrentAccountType).Methods("GET")
	router.HandleFunc("/api/types/all", controllers.GetAllAccountTypes).Methods("GET")
	router.HandleFunc("/api/user/me/status", controllers.GetUserLoginStatus).Methods("GET")
	router.HandleFunc("/api/server/logs", controllers.GetPodLogs).Methods("GET")
	/*
	   =================
	    POST endpoints:
	   =================
	*/

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/server/new/avail_server", controllers.CreateAvailServer).Methods("POST")
	router.HandleFunc("/api/server/new/arma2", controllers.CreateArma2Server).Methods("POST")
	router.HandleFunc("/api/server/new/arkpc", controllers.CreateArkPCServer).Methods("POST")
	router.HandleFunc("/api/server/new/assettocc", controllers.CreateAssettoCCServer).Methods("POST")
	router.HandleFunc("/api/server/new/minecraftbr", controllers.CreateMinecraftBRServer).Methods("POST")
	// router.HandleFun("/api/auth/validate", controllers.ValidateToken).Methods("POST")
	// router.HandleFunc("/api/server/new/minecraftbr", controllers.CreateMinecraftVanillaServer).Methods("POST")

	/*
	   =================
	    DELETE endpoints:
	   =================
	*/

	/*
	   =================
	    Router Logic:
	   =================
	*/

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	log.Fatal(http.ListenAndServe(":8000", c.Handler(hnynethttp.WrapHandler(router))))
	// log.Fatal(http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Headers", "Origin", "Access-Control-Request-Headers", "credentials", "Content-Type", "content-Length", "Accept", "X-CSRF-Token", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"http://localhost:4200"}))(router)))
}

