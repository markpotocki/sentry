package main

import (
	session "idendity-provider/sessions"
	_ "idendity-provider/sessions/memory"
	"log"
	"net/http"
)

var globalSessions *session.Manager

func main() {
	fs := http.FileServer(http.Dir("templates/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/register", RegisterPage)

	log.Println("Serving content on port 8070.")
	http.ListenAndServe(":8070", nil)
}

func init() {
	globalSessions, _ = session.NewManager("memory", "SESSION_ID", 3600)
	go globalSessions.GC()
}
