package main

import (
	"PoliSim/componentHelper"
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
	r.Get("/test", htmlServer.ServeTestGet)

	_, _ = fmt.Fprintf(os.Stdout, "PoliSim is trying to start the listener...\n")
	err := http.ListenAndServe(os.Getenv("ADRESS"), r)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to start the Router:\n"+err.Error()+"\n")
		os.Exit(1)
	}
}
