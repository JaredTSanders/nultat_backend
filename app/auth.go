package app

import (
	"net/http"

	"github.com/JaredTSanders/nultat_backend/models"
	// "github.com/gorilla/securecookie"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// var hashKey = []byte(os.Getenv("token_password"))
		// var blockKey = []byte(os.Getenv("block_key"))

		notAuth := []string{"/api/user/new", "/api/user/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                               //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		c, err := r.Cookie("session_token")

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sessionToken := c.Value

		res, err := models.Cache.Do("GET", sessionToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if res == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		// fmt.Sprintf("User %", tk.UserId) //Useful for monitoring
		// ctx := context.WithValue(r.Context(), "user", tk.UserId)
		// r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
