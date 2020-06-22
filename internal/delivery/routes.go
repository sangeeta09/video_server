package routes

import (
	"net/http"

	"github.com/personal-work/video_server/internal/utils/auth"

	controllers "github.com/personal-work/video_server/internal/usecase"

	"github.com/gorilla/mux"
)

//InitHandlers inits all the routes
func InitHandlers() *mux.Router {

	// r := mux.NewRouter().StrictSlash(true)
	r := mux.NewRouter()
	// r.Use(CommonMiddleware)

	//healthcheck api
	r.HandleFunc("/healthcheck", controllers.HealthCheck).Methods("GET")

	//user management API
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	r.Handle("/user/{id}", auth.JwtVerify(http.HandlerFunc(controllers.GetUser))).Methods("GET")
	r.Handle("/user/{id}", auth.JwtVerify(http.HandlerFunc(controllers.UpdateUser))).Methods("PUT")
	r.Handle("/user/{id}", auth.JwtVerify(http.HandlerFunc(controllers.DeleteUser))).Methods("DELETE")

	return r
}

// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
