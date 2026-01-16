package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// User model
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// In-memory data
var users []User
var nextID = 1

func main() {
	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/users", usersHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Home route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server is running 🚀"))
}

// CRUD handler
func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	// CREATE
	case http.MethodPost:
		var newUser User
		json.NewDecoder(r.Body).Decode(&newUser)

		newUser.ID = nextID
		nextID++
		users = append(users, newUser)

		json.NewEncoder(w).Encode(newUser)

	// READ
	case http.MethodGet:
		idParam := r.URL.Query().Get("id")

		// Single user
		if idParam != "" {
			id, _ := strconv.Atoi(idParam)

			for _, user := range users {
				if user.ID == id {
					json.NewEncoder(w).Encode(user)
					return
				}
			}

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "User not found",
			})
			return
		}

		// All users
		json.NewEncoder(w).Encode(users)

	// UPDATE
	case http.MethodPut:
		idParam := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idParam)

		var updatedUser User
		json.NewDecoder(r.Body).Decode(&updatedUser)

		for i, user := range users {
			if user.ID == id {
				users[i].Name = updatedUser.Name
				users[i].Email = updatedUser.Email
				json.NewEncoder(w).Encode(users[i])
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "User not found",
		})

	// DELETE
	case http.MethodDelete:
		idParam := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idParam)

		for i, user := range users {
			if user.ID == id {
				users = append(users[:i], users[i+1:]...)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "User deleted",
				})
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "User not found",
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
