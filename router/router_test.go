package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestWithRoute(t *testing.T) {

	router := &Router{make([]Modifier, 0), make([]Route, 0)}
	route := MakeRoute("/test/hello", Method("GET"), nil)
	router.WithRoute(route)

	// check that route exists
	if bl := router.routes; len(bl) != 1 {
		t.Log("router_test: route was unable to be found by path index")
		t.FailNow()
	}
}

func TestWith100Route(t *testing.T) {

	router := &Router{make([]Modifier, 0), make([]Route, 0)}
	for i := 0; i < 100; i++ {
		route := MakeRoute(fmt.Sprint(i), Method("GET"), nil)
		router.WithRoute(route)
	}

	// check that route exists
	if bl := router.routes; len(bl) != 100 {
		t.Logf("router_test: route length did not match (length: %d)", len(bl))
		t.FailNow()
	}
}

func TestPathVars(t *testing.T) {

	router := &Router{make([]Modifier, 0), make([]Route, 0)}
	route := MakeRoute("/test/:hello", Method("GET"), nil)
	router.WithRoute(route)

	// check pathvars are set
	for key, ind := range route.pathVars {
		if key != "hello" || ind != 1 {
			t.Logf("router_test: pathvars did not validate. key: %s, ind: %d. Wanted key: hello, ind: 1", key, ind)
			t.FailNow()
		}
	}

}

func TestServe(t *testing.T) {
	result := make(chan bool, 1)
	router := &Router{make([]Modifier, 0), make([]Route, 0)}

	modFunc := func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			result <- true
			next.ServeHTTP(w, r)
		}
	}

	testFunc := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hello")
	}
	route := MakeRoute("/test/:hello", Method("GET"), testFunc)

	router.WithRoute(route)
	router.WithModifier(Modifier(modFunc))

	router.Serve(context.TODO())

	// check if result channel can be heard

}
