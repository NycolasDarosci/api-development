package main

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var users = map[string]User{
	"595ea76c-b130-4c3b-9afa-454a17d90d60": {
		ID:   "595ea76c-b130-4c3b-9afa-454a17d90d60",
		Name: "Jailson",
	},
	"765a42a9-8783-4808-ae5a-af57cf07d9e0": {
		ID:   "765a42a9-8783-4808-ae5a-af57cf07d9e0",
		Name: "Mendes",
	},
	"e5c4a6e0-c68d-4389-bcf3-c3f908c27e5e": {
		ID:   "e5c4a6e0-c68d-4389-bcf3-c3f908c27e5e",
		Name: "Alice",
	},
}
