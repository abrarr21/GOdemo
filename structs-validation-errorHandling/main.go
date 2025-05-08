package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func (user *User) Normalize() {
	if user.Name == "" {
		user.Name = "Unknown"
	}
}

// func validateUser(createdUser User) is a one way of implementation
func (createdUser User) ValidateUser() error {
	// if createdUser.Name == "" {
	// 	return fmt.Errorf("missing field: Name")
	// }
	// No Longer Needed as it got normalize

	if createdUser.Email == "" {
		return fmt.Errorf("missing field: Email")
	}

	if createdUser.Age <= 0 {
		return fmt.Errorf("Invalid: Age")
	}

	return nil

}

func handler(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
	}

	// if we had used func(createdUser User) then err = validateUser(user)
	user.Normalize()
	err = user.ValidateUser()
	if err != nil {
		response := map[string]string{"error": err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode json response", http.StatusInternalServerError)
		}

		return

	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User data is valid"))
}
