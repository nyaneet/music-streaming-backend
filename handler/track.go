package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nyaneet/music-streaming-backend/db"
	"github.com/nyaneet/music-streaming-backend/models"
)

var trackIdKey = "trackId"

func tracks(router chi.Router) {
	router.Get("/", getTracks)
	router.Get("/find", findTracks)
	router.Route("/{id}", func(router chi.Router) {
		router.Use(extractId)
		router.Get("/", getTrack)
	})
}

func getTrack(w http.ResponseWriter, req *http.Request) {
	trackId := req.Context().Value("id").(int)

	track, err := dbInstance.GetTrackById(trackId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, req, ErrNotFound)
			return
		}
		render.Render(w, req, ErrInternalServerError)
		return
	}

	if err := render.Render(w, req, &track); err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}
}

func getTracks(w http.ResponseWriter, req *http.Request) {
	tracks, err := dbInstance.GetAllTracks()
	if err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	if err := render.Render(w, req, tracks); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}
}

func findTracks(w http.ResponseWriter, req *http.Request) {
	query := models.SearchQuery{}
	if err := render.Bind(req, &query); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}

	tracks, err := dbInstance.FindTracks(query.Query)
	if err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	if err := render.Render(w, req, tracks); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}
}
