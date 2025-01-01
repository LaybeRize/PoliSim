package main

import (
	"PoliSim/handler"
	"PoliSim/handler/accounts"
	"PoliSim/handler/notes"
	"PoliSim/handler/titles"
	"fmt"
	"net/http"
	"os"
)

func main() {
	_, _ = fmt.Fprintf(os.Stdout, "Registering all Pages\n")
	fs := http.FileServer(http.Dir("public"))
	http.Handle("GET /public/*", http.StripPrefix("/public/", fs))

	http.HandleFunc("GET /create/account", accounts.GetCreateAccount)
	http.HandleFunc("POST /create/account", accounts.PostCreateAccount)
	http.HandleFunc("GET /edit/account", accounts.GetEditAccount)
	http.HandleFunc("PATCH /edit/account", accounts.PostEditAccount)
	http.HandleFunc("PUT /edit/account/search", accounts.PostEditSearchAccount)

	http.HandleFunc("GET /create/title", titles.GetCreateTitlePage)
	http.HandleFunc("POST /create/title", titles.PostCreateTitlePage)
	http.HandleFunc("GET /edit/title", titles.GetEditTitlePage)
	http.HandleFunc("PATCH /edit/title", titles.PatchEditTitlePage)
	http.HandleFunc("PUT /edit/title/search", titles.PutTitleSearchPage)

	http.HandleFunc("GET /my/profile", accounts.GetMyProfile)
	http.HandleFunc("PATCH /my/profile/password", accounts.PostUpdateMyPassword)
	http.HandleFunc("PATCH /my/profile/settings", accounts.PostUpdateMySettings)

	http.HandleFunc("POST /login", accounts.PostLoginAccount)
	http.HandleFunc("POST /logout", accounts.PostLogOutAccount)

	http.HandleFunc("GET /notes/request", notes.RequestNote)
	http.HandleFunc("GET /notes", notes.GetNotesViewPage)
	http.HandleFunc("GET /create/note", notes.GetNoteCreatePage)
	http.HandleFunc("POST /create/note", notes.PostNoteCreatePage)
	http.HandleFunc("GET /search/notes", notes.GetSearchNotePage)
	http.HandleFunc("PUT /search/notes", notes.PutSearchNotePage)

	http.HandleFunc("GET /", handler.GetHomePage)

	http.HandleFunc("/", handler.GetNotFoundPage)

	http.HandleFunc("PUT /markdown", handler.PostMakeMarkdown)

	_, _ = fmt.Fprintf(os.Stdout, "Starting HTML Server: Use http://"+os.Getenv("ADDRESS")+"\n")
	err := http.ListenAndServe(os.Getenv("ADDRESS"), nil)
	if err != nil {
		panic(err)
	}
}
