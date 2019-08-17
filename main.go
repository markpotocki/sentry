package main

import (
	"context"
	"idendity-provider/router"
	_ "idendity-provider/router"
	"idendity-provider/view"
	"log"
	"net/http"
)

func main() {
	loggingMod := func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("main: request made, maybe we should be able to modify endpoints independently?")
			http.Redirect(w, r, "/modworks", 415)
			next.ServeHTTP(w, r)
		}
	}

	router.GlobalRouter.WithRoute(router.MakeRoute("/register", router.Method("BYPASS"), view.RegisterPage))
	router.GlobalRouter.WithRoute(router.MakeRoute("/login", router.Method("BYPASS"), view.LoginPage))
	router.GlobalRouter.WithModifier(loggingMod)
	errchan := router.GlobalRouter.Serve(context.TODO())

	<-errchan
}
