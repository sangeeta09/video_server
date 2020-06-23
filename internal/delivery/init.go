package routes

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/personal-work/video_server/internal/models"
)

func Init() {

	//init user management maps
	models.UserDetailMap = make(map[string]*models.User)
	models.UserIDMap = make(map[uint64]string)

	//init room maps
	models.GuidRoomDetailMap = make(map[uuid.UUID]*models.Room)
	models.UserRoomMap = make(map[uint64][]uuid.UUID)

	// Handle routes
	http.Handle("/", InitHandlers())

}
