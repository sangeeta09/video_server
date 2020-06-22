package routes

import (
	"net/http"

	"github.com/personal-work/video_server/internal/models"
)

func Init() {

	//init maps
	models.UserDetailMap = make(map[string]*models.User)
	models.UserIDMap = make(map[uint64]string)

	// Handle routes
	http.Handle("/", InitHandlers())

}
