package main

import (
	"idendity-provider/database"
	"idendity-provider/user"
	"idendity-provider/view"
	"log"
	"net/http"
)

var gdb *database.Database

func main() {
	fs := http.FileServer(http.Dir("templates/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// init db in view
	Init()

	// save a test user
	user.MakeUser("test", "test", "test", "test")

	http.HandleFunc("/login", view.LoginPage)
	http.HandleFunc("/register", view.RegisterPage)

	log.Println("Serving content on port 8070.")
	http.ListenAndServe(":8070", nil)
}

func Init() {
	db := &database.Database{Host: "localhost", Port: 27017}
	db.Connect()
	gdb = db
	dls := user.DatabaseLoginService{db}
	user.LS = dls
}
