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
	r.Handle("/users/{id}", http.HandlerFunc(controllers.GetUser)).Methods("GET")
	r.Handle("/users/{id}", auth.JwtVerify(http.HandlerFunc(controllers.UpdateUser))).Methods("PUT")
	r.Handle("/users/{id}", auth.JwtVerify(http.HandlerFunc(controllers.DeleteUser))).Methods("DELETE")

	//room mangement routes
	r.Handle("/rooms", auth.JwtVerify(http.HandlerFunc(controllers.CreateRoom))).Methods("POST")
	r.Handle("/rooms/{guid}", http.HandlerFunc(controllers.GetRoomInfo)).Methods("GET")
	r.Handle("/rooms/{guid}/users", auth.JwtVerify(http.HandlerFunc(controllers.JoinRoom))).Methods("POST")
	r.Handle("/rooms/{guid}/users", auth.JwtVerify(http.HandlerFunc(controllers.LeaveRoom))).Methods("DELETE")
	r.Handle("/rooms/{guid}", auth.JwtVerify(http.HandlerFunc(controllers.ChangeHost))).Methods("PUT")
	r.Handle("/users/{id}/rooms", http.HandlerFunc(controllers.GetRoomInfoList)).Methods("GET")

	return r
}
