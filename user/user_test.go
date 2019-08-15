package user

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func testEncodePwd(t *testing.T) {
	pwd := "thisisbreakable"
	epwd := EncodePassword(pwd)

	// no errors, test if valid
	e2pwd := EncodePassword(pwd)

	if epwd != e2pwd {
		t.Log("hashes of password did not match")
		t.Fail()
	}

	if err := bcrypt.CompareHashAndPassword([]byte(epwd), []byte(pwd)); err != nil {
		t.Fail()
	}
}
