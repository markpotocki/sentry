package view

import (
	"html/template"
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
