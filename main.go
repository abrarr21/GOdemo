package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strings"
// 	"time"
// )

// LOGGER MIDDLEWARE

// func loggerMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
//
// 		next.ServeHTTP(w, r)
// 	})
// }

// CUSTOM HEADER ADDER MIDDLEWARE

// func headerMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Implement Login
// 		w.Header().Set("X-Custom-Header", "GolangGolang")
// 		// End of Middleware Logic
// 		next.ServeHTTP(w, r)
// 	})
// }

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Welcome to the Home page")
// }
// func aboutHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Welcome to the about page")
// }
//
// func greetHandler(w http.ResponseWriter, r *http.Request) {
// 	// Extracting Query parameter
// 	query := r.URL.Query()
// 	name := query.Get("name")
// 	surname := query.Get("surname")
//
// 	if name == "" {
// 		name = "Guest"
// 	}
//
// 	fmt.Fprintf(w, "Hello %s", name)
// 	fmt.Fprintf(w, "Hello %s", surname)
// }

// http://localhost:6969/user/123
// localhost:6969 -> 0
// user           -> 1
// 123            -> 2
// func userHandler(w http.ResponseWriter, r *http.Request) {
// 	//Extracting Path
// 	pathSegments := strings.Split(r.URL.Path, "/")
// 	if len(pathSegments) >= 3 && pathSegments[1] == "user" {
// 		userID := pathSegments[2]
// 		fmt.Fprintf(w, "User ID: %s", userID)
// 	} else {
// 		http.NotFound(w, r)
// 	}
// }

// combining both query params + path
// http://localhost:6969/username/123?includeDetails=true
// func usernameHandler(w http.ResponseWriter, r *http.Request) {
// 	pathSegments := strings.Split(r.URL.Path, "/")
// 	query := r.URL.Query()
// 	includeDetails := query.Get("includeDetails")
//
// 	if len(pathSegments) >= 3 && pathSegments[1] == "username" {
// 		userID := pathSegments[2]
// 		response := fmt.Sprintf("UserID: %s", userID)
// 		if includeDetails == "true" {
// 			response += " (Details included)"
// 		}
// 		fmt.Fprintln(w, response)
//
// 	} else {
// 		http.NotFound(w, r)
// 	}
// }

// func main() {
// 	mux := http.NewServeMux()
//
// 	mux.Handle("/", loggerMiddleware(headerMiddleware(http.HandlerFunc(homeHandler))))
// 	mux.Handle("/about", http.HandlerFunc(aboutHandler))
// 	mux.Handle("/greet", http.HandlerFunc(greetHandler))
// 	mux.Handle("/user/", http.HandlerFunc(userHandler))
// 	mux.Handle("/username/", http.HandlerFunc(usernameHandler))
//
// 	log.Println("Starting the SERVER at port: 6969")
// 	if err := http.ListenAndServe(":6969", mux); err != nil {
// 		log.Fatal("Server failed!!!!!!", err)
// 	}
//
// }

// ----------------------------------------------------------------	WORKING WITH JSON ------------------------------------------------------
// ----------------------------------------------------------------	WORKING WITH JSON ------------------------------------------------------
// ----------------------------------------------------------------	WORKING WITH JSON ------------------------------------------------------

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	users = make(map[int]User)
	idSeq = 1
	mutex = &sync.Mutex{}
)

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cotent-Type", "application/json")

	switch r.Method {
	case "GET":
		mutex.Lock()
		defer mutex.Unlock()

		usersList := make([]User, 0, len(users))
		for _, user := range users {
			usersList = append(usersList, user)
		}

		json.NewEncoder(w).Encode(usersList)

	case "POST":
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid Json", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		user.ID = idSeq
		idSeq++
		users[user.ID] = user
		mutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func userModifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var id int
	_, err := fmt.Sscanf(r.URL.Path, "/userModify/%d", &id)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	user, exists := users[id]
	if !exists {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(user)

	case "PUT":
		var updatedUser User
		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		updatedUser.ID = id
		users[id] = updatedUser
		json.NewEncoder(w).Encode(updatedUser)

	case "DELETE":
		delete(users, id)
		w.WriteHeader(http.StatusNotFound)

	default:
		http.Error(w, "Error Not Allowed", http.StatusMethodNotAllowed)
	}
}

func main() {

	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/userModify/", userModifyHandler)

	fmt.Println("Server is running at port: 6969.....")
	log.Fatal(http.ListenAndServe(":6969", nil))
}
