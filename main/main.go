package main

import (
	"PoliSim/handler"
	"PoliSim/handler/account"
	"fmt"
	"net/http"
	"os"
)

func main() {
	_, _ = fmt.Fprintf(os.Stdout, "Registering all Pages\n")
	fs := http.FileServer(http.Dir("public"))
	http.Handle("GET public/*", http.StripPrefix("/public/", fs))

	http.HandleFunc("GET /create/account", account.GetCreateAccount)
	http.HandleFunc("POST /create/account", account.PostCreateAccount)

	http.HandleFunc("POST /login", account.PostLoginAccount)

	http.HandleFunc("GET /", handler.GetHomePage)

	http.HandleFunc("/", handler.GetNotFoundPage)

	_, _ = fmt.Fprintf(os.Stdout, "Starting HTML Server\n")
	err := http.ListenAndServe(os.Getenv("ADDRESS"), nil)
	if err != nil {
		panic(err)
	}
}
