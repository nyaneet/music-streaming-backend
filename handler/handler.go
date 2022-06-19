package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nyaneet/music-streaming-backend/db"
)

var dbInstance db.Database

func NewHandler(db db.Database) http.Handler {
	router := chi.NewRouter()
	dbInstance = db
	router.MethodNotAllowed(notAllowedHandler)
	router.NotFound(notFoundHandler)
	router.Route("/tracks", tracks)
	router.Route("/artists", artists)
	router.Route("/albums", albums)
	return router
}

func notAllowedHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(405)
	render.Render(w, req, ErrNotAllowed)
}

func notFoundHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(400)
	render.Render(w, req, ErrNotFound)
}
