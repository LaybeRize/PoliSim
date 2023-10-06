package main

import (
	"PoliSim/componentHelper"
	"PoliSim/dataExtraction"
	"PoliSim/database"
	"PoliSim/htmlComposition"
	"PoliSim/htmlServer"
	"fmt"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	_, _ = fmt.Fprintf(os.Stdout, "PoliSim starting up...\n\n")

	componentHelper.ImportTranslation(os.Getenv("LANG"))

	database.ConnectDatabase()
	createAdminAccount()

	htmlServer.InstallStart()

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
	err := http.ListenAndServe(os.Getenv("ADRESS"), r)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to start the Router:\n"+err.Error()+"\n")
		os.Exit(1)
	}
}

func instigateRoutes(router *chi.Mux) {
	// sets up the 404 routing
	router.NotFound(htmlServer.GetFullPage(componentHelper.Translation["pageNotFoundTitle"]))
	htmlComposition.PageTitleMap[htmlComposition.NotFound] = componentHelper.Translation["pageNotFoundTitle"]
	router.Get("/"+htmlComposition.APIPreRoute+"*", htmlServer.GetNotFoundService)
	_, _ = fmt.Fprintf(os.Stdout, "Added htmx not found route\n")
	// sets up the standard routes
	for httpRoute, pageTitle := range htmlComposition.PageTitleMap {
		router.Get("/"+string(httpRoute), htmlServer.GetFullPage(pageTitle))
		_, _ = fmt.Fprintf(os.Stdout, "Added Route for: /"+string(httpRoute)+"\n")
	}

	for url, function := range htmlComposition.GetHTMXFunctions {
		router.Get("/"+htmlComposition.APIPreRoute+string(url), function)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Added Get Routes for htmx backend\n")
	for url, function := range htmlComposition.PostHTMXFunctions {
		router.Post("/"+htmlComposition.APIPreRoute+string(url), function)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Added Post Routes for htmx backend\n")
	for url, function := range htmlComposition.PatchHTMXFunctions {
		router.Patch("/"+htmlComposition.APIPreRoute+string(url), function)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Added Patch Routes for htmx backend\n")
	for url, function := range htmlComposition.DeleteHTMXFunctions {
		router.Delete("/"+htmlComposition.APIPreRoute+string(url), function)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Added Delete Routes for htmx backend\n")
}

func createAdminAccount() {
	ok, err := dataExtraction.RootAccountExists()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to find root account:\n"+err.Error()+"\n")
		os.Exit(1)
	} else if ok {
		return
	}

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(os.Getenv("INIT_PASSWORD")), bcrypt.DefaultCost)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to find hash root account password:\n"+err.Error()+"\n")
		os.Exit(1)
	}

	adminAccount := dataExtraction.AccountLogin{
		ID:           1,
		DisplayName:  os.Getenv("INIT_NAME"),
		Username:     os.Getenv("INIT_USERNAME"),
		Password:     string(hashedPassword),
		Suspended:    false,
		RefreshToken: uuid.New().String(),
		Role:         database.HeadAdmin,
	}
	err = adminAccount.CreateMe()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to create root account:\n"+err.Error()+"\n")
		os.Exit(1)
	}
}
