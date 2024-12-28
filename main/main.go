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
	http.Handle("GET /public/*", http.StripPrefix("/public/", fs))

	http.HandleFunc("GET /create/account", account.GetCreateAccount)
	http.HandleFunc("POST /create/account", account.PostCreateAccount)
	http.HandleFunc("GET /edit/account", account.GetEditAccount)
	http.HandleFunc("POST /edit/account", account.PostEditAccount)
	http.HandleFunc("POST /edit/account/search", account.PostEditSearchAccount)

	http.HandleFunc("GET /my/profile", account.GetMyProfile)
	http.HandleFunc("POST /my/profile/password", account.PostUpdateMyPassword)
	http.HandleFunc("POST /my/profile/settings", account.PostUpdateMySettings)

	http.HandleFunc("POST /login", account.PostLoginAccount)
	http.HandleFunc("POST /logout", account.PostLogOutAccount)

	http.HandleFunc("GET /", handler.GetHomePage)

	http.HandleFunc("/", handler.GetNotFoundPage)

	http.HandleFunc("POST /markdown", handler.PostMakeMarkdown)

	_, _ = fmt.Fprintf(os.Stdout, "Starting HTML Server: Use http://"+os.Getenv("ADDRESS")+"\n")
	err := http.ListenAndServe(os.Getenv("ADDRESS"), nil)
	if err != nil {
		panic(err)
	}
}
