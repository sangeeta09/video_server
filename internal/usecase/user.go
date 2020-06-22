package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/personal-work/video_server/internal/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type ErrorResponse struct {
	Err string
}

type error interface {
	Error() string
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("tested OK"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		var resp = map[string]interface{}{"message": "Invalid request"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp, status := VerifyUser(user.Email, user.Password)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func VerifyUser(email, password string) (map[string]interface{}, int) {

	if _, ok := models.UserDetailMap[email]; !ok {
		return map[string]interface{}{"message": "Email address not found"}, http.StatusUnauthorized
	}

	user := models.UserDetailMap[email]

	expiresAt := time.Now().Add(time.Second * 30).Unix()

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		//Password does not match!
		return map[string]interface{}{"message": "Invalid login credentials. Please try again"}, http.StatusUnauthorized
	}

	tk := &models.Token{
		UserID: user.UserID,
		Name:   user.Name,
		Email:  user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println(err)
		return map[string]interface{}{"message": "Temporary Server Error"}, http.StatusInternalServerError
	}

	var resp = map[string]interface{}{"message": "logged in"}
	//Store the token in the response
	resp["token"] = tokenString
	resp["email"] = user.Email
	return resp, http.StatusOK
}

//CreateUser function -- create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	// log.Println("user ", user)

	//check if user already present
	if _, ok := models.UserDetailMap[user.Email]; ok {
		var resp = map[string]interface{}{"message": "User with this Email address already exists"}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(resp)
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		err := ErrorResponse{
			Err: "Password Encryption failed",
		}
		json.NewEncoder(w).Encode(err)
		return
	}

	user.Password = string(pass)

	models.UserIDCounter++

	user.UserID = models.UserIDCounter

	//insert in global map
	models.UserDetailMap[user.Email] = user
	models.UserIDMap[user.UserID] = user.Email

	w.WriteHeader(http.StatusCreated)
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
