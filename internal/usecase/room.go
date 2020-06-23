package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/personal-work/video_server/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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

//GetAllUsers function return list of all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	for _, user := range models.UserDetailMap {

		users = append(users, *user)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	userDetail, resp, httpStatus := GetUserDetailFromUserID(w, r)
	if httpStatus != http.StatusOK {
		w.WriteHeader(httpStatus)
		json.NewEncoder(w).Encode(resp)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(user.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		//Password does not match!, update password
		pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			err := ErrorResponse{
				Err: "Password Encryption failed",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		userDetail.Password = string(pass)

	}

	//update mobile token if new token provided
	if userDetail.MobileToken != user.MobileToken {
		userDetail.MobileToken = user.MobileToken
	}

	json.NewDecoder(r.Body).Decode(userDetail)

	json.NewEncoder(w).Encode(&userDetail)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]

	ID, err := strconv.ParseUint(id, 0, 64)
	if err != nil {
		var resp = map[string]interface{}{"message": "UserId is not integer"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return

	}

	if _, ok := models.UserIDMap[ID]; !ok {
		var resp = map[string]interface{}{"message": "User with this User ID doesn't exist"}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(resp)
		return
	}

	email := models.UserIDMap[ID]

	if _, ok := models.UserDetailMap[email]; !ok {
		var resp = map[string]interface{}{"message": "User with this Email address doesn't exist"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	delete(models.UserDetailMap, email)
	delete(models.UserIDMap, ID)

	json.NewEncoder(w).Encode("User deleted")

}

func GetUser(w http.ResponseWriter, r *http.Request) {

	user, resp, status := GetUserDetailFromUserID(w, r)
	w.WriteHeader(status)
	if status != http.StatusOK {
		json.NewEncoder(w).Encode(resp)
		return
	}
	json.NewEncoder(w).Encode(&user)
	return
}

func GetUserDetailFromUserID(w http.ResponseWriter, r *http.Request) (*models.User, map[string]interface{}, int) {

	params := mux.Vars(r)
	var id = params["id"]
	var user *models.User

	ID, err := strconv.ParseUint(id, 0, 64)
	if err != nil {
		return nil, map[string]interface{}{"message": "UserId is not integer"}, http.StatusBadRequest
	}

	if _, ok := models.UserIDMap[ID]; !ok {
		return nil, map[string]interface{}{"message": "User with this User ID doesn't exist"}, http.StatusBadRequest
	}

	email := models.UserIDMap[ID]

	if _, ok := models.UserDetailMap[email]; !ok {
		return nil, map[string]interface{}{"message": "User with this Email address doesn't exist"}, http.StatusInternalServerError
	}

	user = models.UserDetailMap[email]
	return user, nil, http.StatusOK
}
