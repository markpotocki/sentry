package user

import (
	"golang.org/x/crypto/bcrypt"
)

type LoginService interface {
	Login(username string, pwd string) bool
}

type BasicLoginService struct {
	db map[string]string
}

var db = map[string]string{
	"test": encodePassword("test"),
}

func Login(username string, pwd string) bool {
	// encode the password
	val := db[username]

	if val == "" {
		return false
	} else if bcrypt.CompareHashAndPassword([]byte(val), []byte(pwd)) == nil {
		return true
	} else {
		return false
	}
}
