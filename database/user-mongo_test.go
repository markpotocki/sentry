package database

import "testing"

func ConnectionTest(t *testing.T) {
	db := Database{
		host: "localhost",
		port: 27017,
	}

	db.Connect()
}
