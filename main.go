package main

import (
	"log"
	"net/http"
	"os"

	"github.com/JaredTSanders/nultat_backend/app"
	"github.com/JaredTSanders/nultat_backend/controllers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/user/logout", controllers.Logout).Methods("GET")
	router.HandleFunc("/api/server/new/avail_server", controllers.CreateAvailServer).Methods("POST")
	router.HandleFunc("/api/server/get/avail_server", controllers.GetAllAvailServers).Methods("GET")
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

	// fmt.Println(port)

	// c := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"*"},

	// 	AllowCredentials: true,
	// })

	// handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
	// http.ListenAndServe(":"+port, handlers.CORS(handlers.Allowed)) //Launch the app, visit localhost:8000/api
	// if err != nil {
	// 	fmt.Print(err)
	// }
	//    log.Fatal(http.ListenAndServe(":3000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
