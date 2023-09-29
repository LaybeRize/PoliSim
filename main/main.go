package main

import (
	. "PoliSim/componentBuilder"
	"fmt"
	"github.com/go-chi/cors"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	_, _ = fmt.Fprintf(os.Stdout, "PoliSim starting up...\n\n")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{os.Getenv("CORS_URL")},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE"},
	}))

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		err := RenderHTMLDoc(w,
			El(HEAD, El(TITLE, Text("Test"))),
			El(BODY, Raw("<p>test</p>"), El(DIV, Attr(HXPOST, "/test"), Text("test"))))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	_, _ = fmt.Fprintf(os.Stdout, "PoliSim is trying to start listener.\n")
	err := http.ListenAndServe(os.Getenv("ADRESS"), r)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to start the Router:\n"+err.Error()+"\n")
		os.Exit(1)
	}
}
