package user

import (
	"idendity-provider/database"

	"go.mongodb.org/mongo-driver/bson"
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

var db = map[string]string{
	"test": encodePassword("test"),
}

func (dls DatabaseLoginService) Login(username string, pwd string) bool {
	findUser := dls.find("user", "USER")
	result := findUser("username", username)

	switch v := result.(type) {
	case MyUser:
		err := bcrypt.CompareHashAndPassword([]byte(v.Pwd()), []byte(pwd))
		if err != nil {
			return false
		}
		return true
	default:
		return false
	}
}

func (dls DatabaseLoginService) find(db string, col string) func(string, string) interface{} {
	return func(key string, val string) interface{} {
		filter := bson.D{{key, val}}
		result := dls.Db.FindOne(db, col, filter)
		return result
	}
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
