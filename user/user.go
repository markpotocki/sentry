package user

import (
	"context"
	"fmt"
	"idendity-provider/database"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type User interface {
	Username() string
	Pwd() string
	encodePassword(unencodedPwd string) string
}
type MyUser struct {
	UserID   string
	Password string
	Role     string
	Email    string
	Name     string
}

/*
 * Roles
 * 1. User
 * 2. Moderator
 * 3. Administrator
 */

func (m *MyUser) Username() string {
	return m.UserID
}

func (m *MyUser) Pwd() string {
	return m.Password
}

func MakeUser(username string, unencodedPwd string, email string, name string) *MyUser {
	u := makeInternalUser(username, unencodedPwd, "USER", email, name)
	c := database.Mongo.Database("user").Collection("USER")
	insertR, err := c.InsertMany(context.TODO(), []interface{}{u})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted documents: ", insertR.InsertedIDs)
	return u
}

func FindUserByUsername(username string) *MyUser {
	var result MyUser
	filter := bson.D{{"userid", username}}
	c := database.Mongo.Database("user").Collection("USER")
	log.Printf("searching key [%s] with value [%s]\n", "userid", username)
	err := c.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return &result
}

func MakeModerator(username string, unencodedPwd string, email string, name string) *MyUser {
	return makeInternalUser(username, unencodedPwd, "MOD", email, name)
}

func MakeAdmin(username string, unencodedPwd string, email string, name string) *MyUser {
	return makeInternalUser(username, unencodedPwd, "ADMIN", email, name)
}

func makeInternalUser(username string, unencodedPwd string, role string, email string, name string) *MyUser {
	hash, err := bcrypt.GenerateFromPassword([]byte(unencodedPwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return &MyUser{username, string(hash), role, email, name}
}

func EncodePassword(unencodedPwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(unencodedPwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}
