package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nyaneet/music-streaming-backend/db"
)

var artistIdKey = "artistId"

var artistUrlParams = map[string]string{
	"artistIdKey": "artistId",
}

func artists(router chi.Router) {
	router.Get("/", getArtists)
	router.Route("/{id}", func(router chi.Router) {
		router.Use(extractId)
		router.Get("/", getArtist)
	})
}

func getArtist(w http.ResponseWriter, req *http.Request) {
	artistId := req.Context().Value("id").(int)

	artist, err := dbInstance.GetArtistById(artistId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, req, ErrNotFound)
		} else {
			render.Render(w, req, ErrorRenderer(err))
		}
		return
	}

	if err := render.Render(w, req, &artist); err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}
}

func getArtists(w http.ResponseWriter, req *http.Request) {
	artists, err := dbInstance.GetAllArtists()
	if err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	if err := render.Render(w, req, artists); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}
}
