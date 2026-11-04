package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type server struct {
	addr string
}

var users = []string{"Jilmar", "Dagoberto", "Nilson"}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Path == "/users" {
			w.WriteHeader(http.StatusOK)
			for _, u := range users {
				fmt.Fprintln(w, u)
			}
		}
	case "POST":
		if r.URL.Path == "/users" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				return
			}
			if len(body) == 0 {
				http.Error(w, "name is mandatory", http.StatusBadRequest)
				return
			}
			users = append(users, string(body))

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, "User created!")
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
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
