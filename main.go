package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type server struct {
	addr string
}

func (s *server) handleUsers(w http.ResponseWriter, r *http.Request) {
	log.Printf("new Req = Method: %s, Path: %s, request_id: %s\n", r.Method, r.URL.Path, uuid.New().String())
	list := []User{}
	for _, u := range users {
		list = append(list, u)
	}
	writeJSON(w, http.StatusOK, list)
}

func (s *server) handleUsersById(w http.ResponseWriter, r *http.Request) {
	log.Printf("Req = Method: %s, Path: %s, request_id: %s\n", r.Method, r.URL.Path, uuid.New().String())

	id := r.PathValue("id")
	user, ok := users[id]
	if !ok {
		http.Error(w, fmt.Sprintf("user [%s] not found", id), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (s *server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("Req = Method: %s, Path: %s, request_id: %s\n", r.Method, r.URL.Path, uuid.New().String())

	var reqUser CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reqUser.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	user := User{
		ID:       id,
		Name:     reqUser.Name,
		Password: reqUser.Password,
	}
	users[id] = user
	writeJSON(w, http.StatusCreated, user)
}

func (s *server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("Req = Method: %s, Path: %s, request_id: %s\n", r.Method, r.URL.Path, uuid.New().String())

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if _, ok := users[id]; !ok {
		http.Error(w, fmt.Sprintf("user [%s] not found", id), http.StatusNotFound)
		return
	}
	delete(users, id)
	writeJSON(w, http.StatusOK, fmt.Sprintf("user [%s] was deleted", id))
}

func main() {
	s := &server{
		addr: ":8080",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", s.handleUsers)
	mux.HandleFunc("GET /users/{id}", s.handleUsersById)
	mux.HandleFunc("POST /users", s.handleCreateUser)
	mux.HandleFunc("DELETE /users/{id}", s.handleDeleteUser)

	log.Println("server is running..")

	if err := http.ListenAndServe(s.addr, mux); err != nil {
		log.Fatal(err)
	}
}
