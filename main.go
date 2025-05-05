package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// LOGGER MIDDLEWARE

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))

		next.ServeHTTP(w, r)
	})
}

// CUSTOM HEADER ADDER MIDDLEWARE

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement Login
		w.Header().Set("X-Custom-Header", "GolangGolang")
		// End of Middleware Logic
		next.ServeHTTP(w, r)
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Home page")
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the about page")
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	// Extracting Query parameter
	query := r.URL.Query()
	name := query.Get("name")
	surname := query.Get("surname")

	if name == "" {
		name = "Guest"
	}

	fmt.Fprintf(w, "Hello %s", name)
	fmt.Fprintf(w, "Hello %s", surname)
}

// http://localhost:6969/user/123
// localhost:6969 -> 0
// user           -> 1
// 123            -> 2
func userHandler(w http.ResponseWriter, r *http.Request) {
	//Extracting Path
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) >= 3 && pathSegments[1] == "user" {
		userID := pathSegments[2]
		fmt.Fprintf(w, "User ID: %s", userID)
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", loggerMiddleware(headerMiddleware(http.HandlerFunc(homeHandler))))
	mux.Handle("/about", http.HandlerFunc(aboutHandler))
	mux.Handle("/greet", http.HandlerFunc(greetHandler))
	mux.Handle("/user/", http.HandlerFunc(userHandler))

	log.Println("Starting the SERVER at port: 6969")
	if err := http.ListenAndServe(":6969", mux); err != nil {
		log.Fatal("Server failed!!!!!!", err)
	}
}
