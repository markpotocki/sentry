package user

import (
	"idendity-provider/database"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	db := &database.Database{Host: "localhost", Port: 27017}
	db.Connect()
	dls := DatabaseLoginService{db}

	// we know test, test is in database
	result := dls.Login("test", "test")

	if !result {
		t.Log("failed to authenticate user")
		t.Fail()
	}
}

func TestLoginFail(t *testing.T) {
	db := &database.Database{Host: "localhost", Port: 27017}
	db.Connect()
	dls := DatabaseLoginService{db}

	// we know test, test is in database
	result := dls.Login("test", "wrong")

	if result {
		t.Log("user authenticated when should not have been")
		t.Fail()
	}
}
