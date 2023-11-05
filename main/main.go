package main

import (
	"PoliSim/data/database"
	"PoliSim/data/extraction"
	"PoliSim/data/validation"
	"PoliSim/helper"
	"PoliSim/html/builder"
	"PoliSim/html/composition"
	"PoliSim/html/serving"
	"fmt"
	"github.com/go-chi/cors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	_, _ = fmt.Fprintf(os.Stdout, "PoliSim starting up...\n\n")

	builder.ImportTranslation(os.Getenv("LANG"))
	helper.UpdateAttributes()

	database.ConnectDatabase()
	_, _ = fmt.Fprintf(os.Stdout, "Updating title groups\n")
	extraction.UpdateTitleGroupMap()
	err := extraction.StartupUpdateOrganisation()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Organisationlist Update Error:\n"+err.Error())
	}
	createAdminAccount()

	_, _ = fmt.Fprintf(os.Stdout, "Installing Router paths\n")
	serving.InstallStart()
	serving.InstallAccountManagment()
	serving.InstallErrorPage()
	serving.InstallTitlePages()
	serving.InstallOrganisationPages()
	serving.InstallMarkdown()
	serving.InstallPress()
	serving.InstallLetter()
	serving.InstallDocumentText()

	_, _ = fmt.Fprintf(os.Stdout, "Creating cookie store\n")
	validation.CreateStore()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{os.Getenv("CORS_URL")},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE"},
	}))

	_, _ = fmt.Fprintf(os.Stdout, "Serving static files\n")
	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	_, _ = fmt.Fprintf(os.Stdout, "Setting up dynamic routes\n")
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/start", http.StatusMovedPermanently)
	})
	instigateRoutes(r)

	_, _ = fmt.Fprintf(os.Stdout, "PoliSim is trying to start the listener...\n")
	err = http.ListenAndServe(os.Getenv("ADDRESS"), r)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to start the Router:\n"+err.Error()+"\n")
		os.Exit(1)
	}
}

func instigateRoutes(router *chi.Mux) {
	// sets up the 404 routing
	router.NotFound(serving.GetFullPage(builder.Translation["pageNotFoundTitle"]))
	composition.PageTitleMap[composition.NotFound] = builder.Translation["pageNotFoundTitle"]
	router.Get("/"+composition.APIPreRoute+"*", serving.NotFoundService)
	router.Post("/"+composition.APIPreRoute+"*", serving.NotFoundService)
	router.Patch("/"+composition.APIPreRoute+"*", serving.NotFoundService)
	router.Delete("/"+composition.APIPreRoute+"*", serving.NotFoundService)
	_, _ = fmt.Fprintf(os.Stdout, "Added htmx not found routing\n")

	// sets up the standard routes
	for httpRoute, pageTitle := range composition.PageTitleMap {
		router.Get("/"+string(httpRoute), serving.GetFullPage(pageTitle))
		_, _ = fmt.Fprintf(os.Stdout, "Added Route for: /"+string(httpRoute)+"\n")
	}

	for url, function := range composition.GetHTMXFunctions {
		router.Get("/"+composition.APIPreRoute+string(url), function)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Added Get Routes for htmx backend\n")
	for url, function := range composition.PostHTMXFunctions {
		router.Post("/"+composition.APIPreRoute+string(url), function)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Added Post Routes for htmx backend\n")
	for url, function := range composition.PatchHTMXFunctions {
		router.Patch("/"+composition.APIPreRoute+string(url), function)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Added Patch Routes for htmx backend\n")
}

func createAdminAccount() {
	ok, err := extraction.RootAccountExists()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to find root account:\n"+err.Error()+"\n")
		os.Exit(1)
	} else if ok {
		return
	}
	_, _ = fmt.Fprintf(os.Stdout, "Creating admin account\n")

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(os.Getenv("INIT_PASSWORD")), bcrypt.DefaultCost)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to find hash root account password:\n"+err.Error()+"\n")
		os.Exit(1)
	}

	adminAccount := extraction.AccountLogin{
		ID:          1,
		DisplayName: os.Getenv("INIT_NAME"),
		Username:    os.Getenv("INIT_USERNAME"),
		Password:    string(hashedPassword),
		Suspended:   false,
		Role:        database.HeadAdmin,
	}
	err = adminAccount.CreateMe()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to create root account:\n"+err.Error()+"\n")
		os.Exit(1)
	}
}
