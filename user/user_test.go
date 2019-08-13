package user

import "testing"

func testEncodePwd(t *testing.T) {
	pwd := "thisisbreakable"
	epwd := encodePassword(pwd)

	// no errors, test if valid
	e2pwd := encodePassword(pwd)

	if epwd != e2pwd {
		t.Log("hashes of password did not match")
		t.Fail()
	}
}
