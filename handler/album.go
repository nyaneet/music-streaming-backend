package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nyaneet/music-streaming-backend/db"
)

var albumIdKey = "albumId"

func albums(router chi.Router) {
	router.Get("/", getAlbums)
	router.Route("/{id}", func(router chi.Router) {
		router.Use(extractId)
		router.Get("/", getAlbum)
	})
}

func getAlbum(w http.ResponseWriter, req *http.Request) {
	albumId := req.Context().Value("id").(int)

	album, err := dbInstance.GetAlbumById(albumId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, req, ErrNotFound)
		} else {
			render.Render(w, req, ErrorRenderer(err))
		}
		return
	}

	if err := render.Render(w, req, &album); err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}
}

func getAlbums(w http.ResponseWriter, req *http.Request) {
	albums, err := dbInstance.GetAllAlbums()
	if err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	if err := render.Render(w, req, albums); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}
}
