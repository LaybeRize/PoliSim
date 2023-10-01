package main

import (
	"PoliSim/componentHelper"
	"PoliSim/htmlComposition"
	"PoliSim/htmlServer"
	"fmt"
	"github.com/go-chi/cors"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	_, _ = fmt.Fprintf(os.Stdout, "PoliSim starting up...\n\n")

	componentHelper.ImportTranslation(os.Getenv("LANG"))

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
	for _, httpRoute := range htmlComposition.LoadingList {
		handler := htmlComposition.HandlerList[httpRoute]
		if handler == nil {
			_, _ = fmt.Fprintf(os.Stdout, "Could not find route for: /"+string(httpRoute)+"\n")
			continue
		}
		router.Get("/"+string(httpRoute), htmlServer.GetFullPage(handler.TitleText))
		if handler.GetFunction != nil {
			router.Get("/htmx/"+string(httpRoute), handler.GetFunction)
		}
		if handler.PostFunction != nil {
			router.Post("/htmx/"+string(httpRoute), handler.PostFunction)
		}
		if handler.PatchFunction != nil {
			router.Patch("/htmx/"+string(httpRoute), handler.PatchFunction)
		}
		if handler.DeleteFunction != nil {
			router.Delete("/htmx/"+string(httpRoute), handler.DeleteFunction)
		}
		_, _ = fmt.Fprintf(os.Stdout, "Added Route for: /"+string(httpRoute)+"\n")
	}
}
