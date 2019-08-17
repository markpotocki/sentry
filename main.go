package main

import (
<<<<<<< HEAD
	session "idendity-provider/sessions"
	_ "idendity-provider/sessions/memory"
=======
	"context"
	"idendity-provider/router"
	_ "idendity-provider/router"
	"idendity-provider/view"
>>>>>>> b3489d4638612b5116e546f4b57511b1db7a24f8
	"log"
	"net/http"
)

var globalSessions *session.Manager

func main() {
	loggingMod := func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("main: request made, maybe we should be able to modify endpoints independently?")
			http.Redirect(w, r, "/modworks", 415)
			next.ServeHTTP(w, r)
		}
	}

<<<<<<< HEAD
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/register", RegisterPage)
=======
	router.GlobalRouter.WithRoute(router.MakeRoute("/register", router.Method("BYPASS"), view.RegisterPage))
	router.GlobalRouter.WithRoute(router.MakeRoute("/login", router.Method("BYPASS"), view.LoginPage))
	router.GlobalRouter.WithModifier(loggingMod)
	errchan := router.GlobalRouter.Serve(context.TODO())
>>>>>>> b3489d4638612b5116e546f4b57511b1db7a24f8

	<-errchan
}

func init() {
	globalSessions, _ = session.NewManager("memory", "SESSION_ID", 3600)
	go globalSessions.GC()
}
