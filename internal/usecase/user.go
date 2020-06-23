package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/personal-work/video_server/internal/models"

	jwt "github.com/dgrijalva/jwt-go"
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
