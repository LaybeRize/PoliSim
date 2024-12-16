package main

import (
	"PoliSim/handler/account"
	"net/http"
	"os"
)

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("GET public/*", http.StripPrefix("/public/", fs))

	http.HandleFunc("GET /create/account", account.GetCreateAccount)
	http.HandleFunc("POST /create/account", account.PostCreateAccount)

	http.HandleFunc("POST /login", account.PostLoginAccount)

	err := http.ListenAndServe(os.Getenv("ADDRESS"), nil)
	if err != nil {
		panic(err)
	}
}
