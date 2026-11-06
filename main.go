package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type server struct {
	addr string
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Req = Method: %s, Path: %s, request_id: %s\n", r.Method, r.URL.Path, uuid.New().String())

	switch r.Method {
	case "GET":
		if r.URL.Path == "/users" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			list := make([]User, 0, len(users))
			for _, u := range users {
				//fmt.Fprintln(w, u.ID, u.Name)
				list = append(list, u)
			}
			if err := json.NewEncoder(w).Encode(list); err != nil {
				http.Error(w, "something got wrong!", http.StatusInternalServerError)
			}
			return
		}

		// "/users/{id}"
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) == 3 && parts[1] == "users" {
			id := parts[2]
			user, ok := users[id]
			if !ok {
				http.Error(w, fmt.Sprintf("user [%s] not found", id), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(user); err != nil {
				http.Error(w, "something got wrong!", http.StatusInternalServerError)
			}
		}

	case "POST":
		if r.URL.Path == "/users" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				return
			}
			if string(body) == "" {
				http.Error(w, "name is mandatory", http.StatusBadRequest)
				return
			}

			var user User
			if err = json.Unmarshal(body, &user); err != nil {
				http.Error(w, "error unmarshalling!", http.StatusInternalServerError)
				return
			}
			id := uuid.New().String()
			user.ID = id
			users[id] = user

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(fmt.Sprintf("User with id [%s] created!", id)); err != nil {
				http.Error(w, "something got wrong!", http.StatusInternalServerError)
			}
		}
	case "DELETE":
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) == 3 && parts[1] == "users" {
			id := parts[2]
			if _, ok := users[id]; !ok {
				http.Error(w, fmt.Sprintf("user [%s] not found", id), http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			delete(users, id)
			if err := json.NewEncoder(w).Encode(fmt.Sprintf("User with id [%s] deleted!", id)); err != nil {
				http.Error(w, "something got wrong!", http.StatusInternalServerError)
			}
		}
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func main() {
	s := &server{
		addr: ":8080",
	}

	log.Println("server is running..")

	if err := http.ListenAndServe(s.addr, s); err != nil {
		log.Fatal(err)
	}
}
