package user

import (
	"idendity-provider/database"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type LoginService interface {
	Login(username string, pwd string) bool
}

type BasicLoginService struct {
	db map[string]string
}

type DatabaseLoginService struct {
	Db *database.Database
}

var LS LoginService

func (dls DatabaseLoginService) Login(username string, pwd string) bool {
	log.Printf("looking for user %s in db.\n", username)
	result := FindUserByUsername(username)

	log.Printf("result: %v", result)
	err := bcrypt.CompareHashAndPassword([]byte(result.Pwd()), []byte(pwd))
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

func (b BasicLoginService) Login(username string, pwd string) bool {
	// encode the password
	val := b.db[username]

	if val == "" {
		return false
	} else if bcrypt.CompareHashAndPassword([]byte(val), []byte(pwd)) == nil {
		return true
	} else {
		return false
	}
}
