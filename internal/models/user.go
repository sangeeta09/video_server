package models

//User struct declaration
type User struct {
	UserID      uint64
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	MobileToken string `json:"mobile_token"`
}

var UserIDMap map[uint64]string

var UserDetailMap map[string]*User

var UserIDCounter uint64
