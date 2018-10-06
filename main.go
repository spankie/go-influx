package main

import (
	"net/http"

	"github.com/spankie/go-web/modules/shops"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fileServer := http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/assets")))
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		fileServer.ServeHTTP(w, r)
	})

	r.Get("/", shops.GetAllProducts)
	r.Get("/product/{ID}", shops.GetProduct)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/404.html")
	})
	log.Print("👉  Server started at 127.0.0.1:3333")
	http.ListenAndServe(":3333", r)
}
