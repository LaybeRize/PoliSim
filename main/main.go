package main

import (
	"PoliSim/handler"
	"PoliSim/handler/accounts"
	"PoliSim/handler/newspaper"
	"PoliSim/handler/notes"
	"PoliSim/handler/organisations"
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

	http.HandleFunc("GET /view/titles", titles.GetTitleView)
	http.HandleFunc("GET /single/view/title", titles.GetSingleViewTitle)

	http.HandleFunc("GET /create/organisation", organisations.GetCreateOrganisationPage)
	http.HandleFunc("POST /create/organisation", organisations.PostCreateOrganisationPage)
	http.HandleFunc("GET /edit/organisation", organisations.GetEditOrgansationPage)
	http.HandleFunc("PATCH /edit/organisation", organisations.PatchEditOrganisationPage)
	http.HandleFunc("PUT /edit/organisation/search", organisations.PutOrganisationSearchPage)

	http.HandleFunc("GET /view/organisations", organisations.GetOrganisationView)
	http.HandleFunc("GET /single/view/organisation", organisations.GetSingleOrganisationView)

	http.HandleFunc("GET /search/newspapers", newspaper.GetManageNewspaperPage)

	http.HandleFunc("GET /create/article", newspaper.GetCreateArticlePage)
	http.HandleFunc("GET /newspaper/for/account", newspaper.GetFindNewspaperForAccountPage)
	http.HandleFunc("POST /create/article", newspaper.PostCreateArticlePage)

	http.HandleFunc("GET /check/newspapers", newspaper.GetManageNewspaperPage)
	http.HandleFunc("POST /newspaper/create", newspaper.PostCreateNewspaperPage)
	http.HandleFunc("PATCH /newspaper/update", newspaper.PatchUpdateNewspaperPage)
	http.HandleFunc("PUT /newspaper/search", newspaper.PutSearchNewspaperPage)

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
