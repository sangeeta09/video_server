package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/personal-work/video_server/internal/models"
)

func generateGUID() uuid.UUID {
	return uuid.New()
}

//CreateRoom function -- creates a new room
func CreateRoom(w http.ResponseWriter, r *http.Request) {

	room := &models.Room{}
	json.NewDecoder(r.Body).Decode(room)

	// log.Println("room ", room)

	//check if host and participants are present in system
	if _, ok := models.UserIDMap[room.HostUserID]; !ok {
		var resp = map[string]interface{}{"message": "User with this User ID doesn't exist"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	for _, participant := range room.ParticipantList {
		if _, ok := models.UserIDMap[participant]; !ok {
			var resp = map[string]interface{}{"message": "User with this User ID doesn't exist"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

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

	for _, participant := range room.ParticipantList {
		models.UserRoomMap[participant] = append(models.UserRoomMap[participant], room.GUID)
	}

	log.Println(room)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&room)
}

//GetRoomInfo fetchs room info from guid
func GetRoomInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var sguid = params["guid"]
	var room *models.Room
	var ok bool
	guid, err := uuid.Parse(sguid)
	if err != nil {
		var resp = map[string]interface{}{"message": "Please sent valid Room GUID"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if room, ok = models.GuidRoomDetailMap[guid]; !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		var resp = map[string]interface{}{"message": "Room with given GUID doesn't exist"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&room)
	return
}

//JoinRoom
func JoinRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var sguid = params["guid"]
	guid, err := uuid.Parse(sguid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{"message": "Please sent valid Room GUID"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	tk, ok := getTokenFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		var resp = map[string]interface{}{"message": "Temporary error, Please try again after sometime"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	//get all rooms by userId
	if rooms, ok := models.UserRoomMap[tk.UserID]; ok {
		//check if user already in room
		for _, roomGuid := range rooms {
			if roomGuid == guid {
				w.WriteHeader(http.StatusOK)
				var resp = map[string]interface{}{"message": "User has already joined this room"}
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
	}

	//else get room details
	var room *models.Room
	if room, ok = models.GuidRoomDetailMap[guid]; !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		var resp = map[string]interface{}{"message": "Room with given GUID doesn't exist"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	//check for capacity
	if room.Capacity == len(room.ParticipantList)+1 {
		//already full
		w.WriteHeader(http.StatusTooManyRequests)
		var resp = map[string]interface{}{"message": "Room is full, try another room"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	//add user to room as participant
	room.ParticipantList = append(room.ParticipantList, tk.UserID)

	//update details in global map
	models.GuidRoomDetailMap[room.GUID] = room
	models.UserRoomMap[tk.UserID] = append(models.UserRoomMap[tk.UserID], room.GUID)
	return
}

func getTokenFromContext(ctx context.Context) (*models.Token, bool) {
	token, ok := ctx.Value("user").(*models.Token)
	return token, ok
}

//LeaveRoom
func LeaveRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var sguid = params["guid"]
	guid, err := uuid.Parse(sguid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{"message": "Please sent valid Room GUID"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	//get room details
	var room *models.Room
	var ok bool
	if room, ok = models.GuidRoomDetailMap[guid]; !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		var resp = map[string]interface{}{"message": "Room with given GUID doesn't exist"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	tk, ok := getTokenFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		var resp = map[string]interface{}{"message": "Temporary error, Please try again after sometime"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	//check if user is host
	if tk.UserID == room.HostUserID {

		//check if participant present
		length := len(room.ParticipantList)
		if length > 0 {
			//make last person in list of participant as host
			room.HostUserID = room.ParticipantList[length-1]
			room.ParticipantList = room.ParticipantList[:length-1]
		} else {
			//no one present, delete room also
			//update global maps
			if rooms, ok := models.UserRoomMap[tk.UserID]; ok {
				finalRoomList := []uuid.UUID{}
				for _, roomGuid := range rooms {
					if roomGuid != guid {
						finalRoomList = append(finalRoomList, roomGuid)
					}
				}
				if len(finalRoomList) == 0 {
					delete(models.UserRoomMap, tk.UserID)
				} else {

					models.UserRoomMap[tk.UserID] = finalRoomList
				}
			}
			delete(models.GuidRoomDetailMap, guid)

			w.WriteHeader(http.StatusOK)
			return
		}
	}

	//get all rooms by userId
	if rooms, ok := models.UserRoomMap[tk.UserID]; ok {
		//check if user present in given room
		found := false
		finalRoomList := []uuid.UUID{}
		for _, roomGuid := range rooms {

			if roomGuid == guid {
				found = true
			} else {
				finalRoomList = append(finalRoomList, roomGuid)
			}
		}

		if found {
			if len(finalRoomList) == 0 {
				delete(models.UserRoomMap, tk.UserID)
			} else {

				models.UserRoomMap[tk.UserID] = finalRoomList
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		var resp = map[string]interface{}{"message": "User has not joined this room"}
		json.NewEncoder(w).Encode(resp)
		return

	}
}

//ChangeHost changes role of particpant user to host user
func ChangeHost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var sguid = params["guid"]
	guid, err := uuid.Parse(sguid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{"message": "Please sent valid Room GUID"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	//get room details
	var room *models.Room
	var ok bool
	if room, ok = models.GuidRoomDetailMap[guid]; !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		var resp = map[string]interface{}{"message": "Room with given GUID doesn't exist"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	tk, ok := getTokenFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		var resp = map[string]interface{}{"message": "Temporary error, Please try again after sometime"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	//check if user is host
	if tk.UserID != room.HostUserID {

		found := false
		finalUserList := []uint64{}
		//check if user present in participants
		for _, participant := range room.ParticipantList {
			if participant == tk.UserID {
				found = true
			} else {
				finalUserList = append(finalUserList, participant)
			}
		}
		if !found {
			w.WriteHeader(http.StatusBadRequest)
			var resp = map[string]interface{}{"message": "User is not a participant of this room"}
			json.NewEncoder(w).Encode(resp)
			return
		}
		room.ParticipantList = finalUserList
		room.ParticipantList = append(room.ParticipantList, room.HostUserID)
		room.HostUserID = tk.UserID

	}
}

//GetRoomInfoList
func GetRoomInfoList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]

	ID, err := strconv.ParseUint(id, 0, 64)
	if err != nil {
		var resp = map[string]interface{}{"message": "UserId is not integer"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return

	}

	if _, ok := models.UserRoomMap[ID]; !ok {
		var resp = map[string]interface{}{"message": "This User doesn't have any room"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
	roomList := models.UserRoomMap[ID]

	roomInfoList := []*models.Room{}
	for _, user := range roomList {
		if room, ok := models.GuidRoomDetailMap[user]; ok {
			roomInfoList = append(roomInfoList, room)
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roomInfoList)

}

// r.Handle("/users/{id}/rooms/{guid}", auth.JwtVerify(http.HandlerFunc(controllers.JoinRoom))).Methods("POST")
// r.Handle("/users/{id}/rooms/{guid}", auth.JwtVerify(http.HandlerFunc(controllers.LeaveRoom))).Methods("DELETE")
// r.Handle("/rooms/{guid}/host", auth.JwtVerify(http.HandlerFunc(controllers.ChangeHost))).Methods("PUT")
// r.Handle("/users/{id}/rooms", http.HandlerFunc(controllers.GetRoomInfoList)).Methods("GET")
