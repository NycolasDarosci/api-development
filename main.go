package main

import (
	"fmt"
	"log"
	"net/http"
)

type server struct {
	addr string
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from server!"))
}

func main() {
	fmt.Println("Oi")

	s := &server{
		addr: ":8080",
	}

	log.Println("server is running..")

	if err := http.ListenAndServe(s.addr, s); err != nil {
		log.Fatal(err)
	}
}
