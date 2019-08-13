package user

import "log"
import "golang.org/x/crypto/bcrypt"

type User interface {
	Username() string
	Pwd() string
	encodePassword(unencodedPwd string) string
}
type MyUser struct {
	username string
	password string
	role     string
	email    string
	name     string
}

/*
 * Roles
 * 1. User
 * 2. Moderator
 * 3. Administrator
 */

func (m *MyUser) Username() string {
	return m.username
}

func (m *MyUser) Pwd() string {
	return m.password
}

func MakeUser(username string, unencodedPwd string, email string, name string) *MyUser {
	return makeInternalUser(username, unencodedPwd, "USER", email, name)
}

func MakeModerator(username string, unencodedPwd string, email string, name string) *MyUser {
	return makeInternalUser(username, unencodedPwd, "MOD", email, name)
}

func MakeAdmin(username string, unencodedPwd string, email string, name string) *MyUser {
	return makeInternalUser(username, unencodedPwd, "ADMIN", email, name)
}

func makeInternalUser(username string, unencodedPwd string, role string, email string, name string) *MyUser {
	return &MyUser{username, encodePassword(unencodedPwd), role, email, name}
}

func encodePassword(unencodedPwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(unencodedPwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}
