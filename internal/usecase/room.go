package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/personal-work/video_server/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func generateGUID() uuid.UUID {
	return uuid.New()
}

//CreateRoom function -- creates a new room
func CreateRoom(w http.ResponseWriter, r *http.Request) {

	room := &models.Room{}
	json.NewDecoder(r.Body).Decode(room)

	// log.Println("room ", room)

	if len(room.GUID) != 0 {
		//check if room already created
		if _, ok := models.GuidRoomDetailMap[room.GUID]; ok {
			var resp = map[string]interface{}{"message": "room with this GUID already exists"}
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
	//generate guid
	//fill in guid info
	room.GUID = generateGUID()
	//if room capacity not set
	if room.Capacity == 0 {
		//default room capacity
		room.Capacity = 5
	}
	//insert in global map
	models.GuidRoomDetailMap[room.GUID] = room
	models.UserRoomMap[room.HostUserID] = append(models.UserRoomMap[room.HostUserID], room.GUID)

	log.Println(room)

	w.WriteHeader(http.StatusCreated)
}

//GetRoomInfo fetchs room info from guid
func GetRoomInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var sguid = params["id"]
	var room *models.Room
	var ok bool
	guid, err := uuid.Parse(sguid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if room, ok = models.GuidRoomDetailMap[guid]; !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&room)
	return
}
