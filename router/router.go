package router

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type Router struct {
	modifiers []Modifier
	routes    []Route
	//sessionManager *session.Manager
}

type Route struct {
	path     string
	method   Method
	pathVars map[string]int
	handler  http.HandlerFunc
}

type Method string

type Modifier func(handler http.Handler) http.HandlerFunc

var GlobalRouter *Router

func init() {
	GlobalRouter = &Router{make([]Modifier, 0), make([]Route, 0)}
}

func MakeRoute(path string, method Method, handlerFunc http.HandlerFunc) Route {
	return Route{path, method, make(map[string]int), handlerFunc}
}

func (r *Router) WithRoute(route Route) {
	// split path string by delimiter "/"
	trim := strings.Trim(route.path, "/")
	split := strings.Split(trim, "/")
	// loop through split and find :var values
	for n, val := range split {
		if len(val) < 1 {
			return
		}
		// edge case nothing on path
		log.Printf("val: %v", val)
		if val[0] == byte(':') {
			route.pathVars[val[1:]] = n
		}
	}

	/* Added during serve, most likley not needed
	// modify the handler with all known modifiers
	log.Printf("is modifiers null? %v", r.modifiers)
	for _, modi := range r.modifiers {
		log.Printf("is handler null? %v", route.handler)
		log.Printf("is result null? %v", modi(route.handler))
		route.handler = modi(route.handler)
	}*/

	r.routes = append(r.routes, route)
}

func (r *Router) WithModifier(mod Modifier) {
	r.modifiers = append(r.modifiers, mod)
}

func (r *Router) Serve(ctx context.Context) chan error {
	// add file server for static resources
	// todo add config to change this directory
	log.Println("router: adding static file server at default location")
	fs := http.FileServer(http.Dir("templates/static"))

	// add handler for file server
	// todo add config
	log.Println("router: adding static handler at default prefix")
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// apply all modifiers to routes
	log.Println("router: applying modifiers to handlers")
	for _, mods := range r.modifiers {
		for _, hand := range r.routes {
			hand.handler = mods(hand.handler)
		}
	}

	// compile the routes to look for dup methods
	log.Println("router: creating routes")
	cache := make(map[string]map[Method]Route)
	for _, route := range r.routes {
		// check if route exists
		if method := cache[route.path]; method != nil {
			// weve seen this path before, save the method

			// check if method exists
			if _, ok := method[route.method]; ok != true {
				log.Fatalf("router: duplicate path configured for %s (incorrect methods?)", route.path)
			}
			// method not seen before, add to list
			method[route.method] = route
		} else {
			// route does not exist
			cache[route.path] = make(map[Method]Route)
			cache[route.path][route.method] = route
		}
	}

	for key, val := range cache {
		var gh, ph, puh, dh http.HandlerFunc
		bypass := false
		for method, route := range val {
			switch blah := strings.ToUpper(string(method)); blah {
			case "GET":
				gh = route.handler
			case "POST":
				ph = route.handler
			case "PUT":
				puh = route.handler
			case "DELETE":
				dh = route.handler
			case "BYPASS":
				log.Println("router: warning! using BYPASS method")
				http.HandleFunc(key, route.handler)
				bypass = true

			default:
				log.Fatalf("router: unsupported method %s", blah)
			}
			log.Printf("router: added route Route[path: %s, method: %s]", key, method)
		}
		if !bypass {
			hand := multiHandler(gh, ph, puh, dh)

			http.HandleFunc(key, hand)
		}
	}

	log.Println("Serving content on port 8070.")
	errchan := make(chan error)
	go func() {
		errchan <- http.ListenAndServe(":8070", nil)
	}()
	return errchan
}

func multiHandler(getHandler, postHandler, putHandler, deleteHandler http.HandlerFunc) http.HandlerFunc {
	gh, ph, puh, dh := getHandler, postHandler, putHandler, deleteHandler
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// is supported
			if gh != nil {
				gh(w, r)
			}
			// not supported
			http.Error(w, "GET is not supported", http.StatusMethodNotAllowed)
		case http.MethodPost:
			// is supported
			if ph != nil {
				ph(w, r)
			}
			// not supported
			http.Error(w, "POST is not supported", http.StatusMethodNotAllowed)
		case http.MethodPut:
			// is supported
			if puh != nil {
				puh(w, r)
			}
			// not supported
			http.Error(w, "PUT is not supported", http.StatusMethodNotAllowed)
		case http.MethodDelete:
			// is supported
			if dh != nil {
				dh(w, r)
			}
			// not supported
			http.Error(w, "DELETE is not supported", http.StatusMethodNotAllowed)
		}
	}
}

func (r *Route) PathVars() map[string]string {
	vars := make(map[string]string)
	// split path string by delimiter "/"
	split := strings.Split(r.path, "/")
	// loop through split and find :var values
	for key, val := range r.pathVars {
		vars[key] = split[val]
	}

	return vars
}
