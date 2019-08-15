package main

import (
	"idendity-provider/view"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("templates/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// init db in view
	view.Init()

	http.HandleFunc("/login", view.LoginPage)
	http.HandleFunc("/register", view.RegisterPage)

	log.Println("Serving content on port 8070.")
	http.ListenAndServe(":8070", nil)
}
