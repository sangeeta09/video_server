package routes

import (
	"net/http"

	"github.com/personal-work/video_server/internal/utils/auth"

	controllers "github.com/personal-work/video_server/internal/usecase"

	"github.com/gorilla/mux"
)

//InitHandlers inits all the routes
func InitHandlers() *mux.Router {

	r := mux.NewRouter()

	//healthcheck api
	r.HandleFunc("/healthcheck", controllers.HealthCheck).Methods("GET")

	//user management API
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	r.Handle("/user/{id}", auth.JwtVerify(http.HandlerFunc(controllers.GetUser))).Methods("GET")
	r.Handle("/user/{id}", auth.JwtVerify(http.HandlerFunc(controllers.UpdateUser))).Methods("PUT")
	r.Handle("/user/{id}", auth.JwtVerify(http.HandlerFunc(controllers.DeleteUser))).Methods("DELETE")

	//room mangement routes
	r.Handle("/rooms", auth.JwtVerify(http.HandlerFunc(controllers.CreateRoom))).Methods("POST")
	r.Handle("/rooms/{id}", http.HandlerFunc(controllers.GetRoomInfo)).Methods("GET")

	return r
}
