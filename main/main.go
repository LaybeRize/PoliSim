package main

import (
	"PoliSim/database"
	"PoliSim/handler"
	"PoliSim/handler/accounts"
	"PoliSim/handler/documents"
	"PoliSim/handler/letter"
	"PoliSim/handler/newspaper"
	"PoliSim/handler/notes"
	"PoliSim/handler/organisations"
	"PoliSim/handler/titles"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Registering all Pages")
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("GET /public/", http.StripPrefix("/public/", fs))

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

	http.HandleFunc("GET /create/document", documents.GetCreateDocumentPage)
	http.HandleFunc("POST /create/document", documents.PostCreateDocumentPage)
	http.HandleFunc("GET /create/discussion", documents.GetCreateDiscussionPage)
	http.HandleFunc("POST /create/discussion", documents.PostCreateDiscussionPage)
	http.HandleFunc("POST /create/discussion/comment/{id}", documents.PostCreateComment)
	http.HandleFunc("GET /create/vote", documents.GetCreateVotePage)
	http.HandleFunc("POST /create/vote", documents.PostCreateVotePage)
	http.HandleFunc("GET /create/vote/element", documents.GetCreateVoteElementPage)
	http.HandleFunc("POST /create/vote/element", documents.PostCreateVoteElementPage)
	http.HandleFunc("PATCH /retrieve/vote/element", documents.PatchGetVoteElementPage)
	http.HandleFunc("GET /organisations/for/account", documents.GetFindOrganisationForAccountPage)
	http.HandleFunc("PATCH /check/reader/and/participants", documents.PatchFixUserList)

	http.HandleFunc("GET /view/document/{id}", documents.GetDocumentViewPage)
	http.HandleFunc("GET /search/documents", documents.GetSearchDocumentsPage)
	http.HandleFunc("PUT /search/documents", documents.PutSearchDocumentsPage)
	http.HandleFunc("GET /view/vote/{id}", documents.GetVoteView)
	http.HandleFunc("POST /vote/on/{id}", documents.PostVote)

	http.HandleFunc("GET /manage/tag-colors", documents.GetColorPage)
	http.HandleFunc("POST /create/tag-color", documents.PostCreateColor)
	http.HandleFunc("DELETE /delete/tag-color", documents.DeleteColor)

	http.HandleFunc("GET /view/organisations", organisations.GetOrganisationView)
	http.HandleFunc("GET /single/view/organisation", organisations.GetSingleOrganisationView)

	http.HandleFunc("GET /search/publications", newspaper.GetSearchPublicationsPage)
	http.HandleFunc("PUT /search/publications", newspaper.PutSearchPublicationPage)
	http.HandleFunc("GET /publication/view/{id}", newspaper.GetSpecificPublicationPage)
	http.HandleFunc("PATCH /publicate/{id}", newspaper.PatchPublishPublication)
	http.HandleFunc("DELETE /article/{id}", newspaper.DeleteArticle)

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
	http.HandleFunc("PATCH /change/blocked/note/{id}", notes.UnBlockNote)

	http.HandleFunc("GET /my/letter", letter.GetPagePersonalLetter)
	http.HandleFunc("PUT /my/letter", letter.PutPagePersonalLetter)
	http.HandleFunc("GET /create/letter", letter.GetCreateLetterPage)
	http.HandleFunc("POST /create/letter", letter.PostCreateLetterPage)
	http.HandleFunc("PATCH /check/create/letter", letter.PatchCheckCreateLetterPage)
	http.HandleFunc("GET /admin/letter/search", letter.GetAdminLetterSearchPage)
	http.HandleFunc("GET /view/letter/{id}", letter.GetLetterViewPage)
	http.HandleFunc("PATCH /view/letter/{id}", letter.PatchLetterViewPage)

	http.HandleFunc("GET /", handler.GetHomePage)

	http.HandleFunc("/", handler.GetNotFoundPage)

	http.HandleFunc("PUT /markdown", handler.PostMakeMarkdown)

	log.Println("Starting HTML Server: Use http://" + os.Getenv("ADDRESS"))
	serverHandling()
}

func serverHandling() {
	server := &http.Server{
		Addr: os.Getenv("ADDRESS"),
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	database.Shutdown()
	log.Println("Graceful shutdown complete.")
}
