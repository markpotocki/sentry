package main

import (
	"html/template"
	user "idendity-provider/user"
	"log"
	"net/http"
)

type LoginPageVars struct {
	PageTitle    string
	UsernameHelp string
	PwdHelp      string
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getLoginPage(w, r)
	case http.MethodPost:
		postLoginPage(w, r)
	}
}

func getLoginPage(w http.ResponseWriter, r *http.Request) {
	ses := globalSessions.SessionStart(w, r)
	// remove this check; checks for session auth and bypasses if logged in
	if ses.Get("Authenticated") != nil {
		log.Println("User is already authenticated.")
		http.Redirect(w, r, "/hello", http.StatusSeeOther)
	}
	pageVars := LoginPageVars{
		PageTitle:    "MEP -- Login",
		UsernameHelp: "Enter the username you used when registering. If your account is newer, it will be your email.",
		PwdHelp:      "Enter your password here.",
	}

	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	err := tmpl.Execute(w, pageVars)

	if err != nil {
		log.Println(err)
	}
}

func postLoginPage(w http.ResponseWriter, r *http.Request) {
	ses := globalSessions.SessionStart(w, r)
	// todo remove extra logging
	log.Println("Request to login recieved.")
	// extract username and password from request
	r.ParseForm() // todo add error check
	un := r.Form.Get("username")
	pwd := r.Form.Get("password")
	log.Printf("Request made for user: %s with pwd: %s.", un, pwd)
	isLoggedIn := user.Login(un, pwd)
	// todo save user session

	if isLoggedIn {
		log.Println("User logged in successfully.")
		ses.Set("Authenticated", isLoggedIn)
		http.Redirect(w, r, "/hello", http.StatusSeeOther)
		return
	}
	log.Println("User failed to authenticate.")
	http.Redirect(w, r, "/login?error", http.StatusUnauthorized)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getRegisterPage(w, r)
	case http.MethodPost:
		postRegisterPage(w, r)
	}
}

func getRegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/registration.html"))
	err := tmpl.Execute(w, nil)

	if err != nil {
		log.Println(err)
	}
}

func postRegisterPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved post request to register")
	log.Printf("Request: %v\n", r.Body)
}
