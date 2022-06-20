package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nyaneet/music-streaming-backend/db"
)

var dbInstance db.Database

func NewHandler(db db.Database) http.Handler {
	dbInstance = db
	router := chi.NewRouter()

	router.NotFound(notFoundHandler)

	router.Route("/auth", auth)
	router.Route("/tracks", tracks)
	router.Route("/albums", albums)
	router.Route("/artists", artists)

	router.Route("/users", func(router chi.Router) {
		router.Use(isAdmin)
		router.Route("/", users)
	})

	return router
}

func notFoundHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	render.Render(w, req, ErrNotFound)
}
