package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nyaneet/music-streaming-backend/models"
)

func me(router chi.Router) {
	router.Get("/tracks", getUserTracks)
	router.Put("/tracks", addUserTrack)
	router.Delete("/tracks", removeUserTrack)
	router.Post("/tracks", addUserTrackMetadata)

	router.Get("/top", getRecommendedTracks)
}

func getUserTracks(w http.ResponseWriter, req *http.Request) {
	auth := req.Context().Value("auth").(map[string]string)
	username, ok := auth["username"]
	if !ok {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	tracks, err := dbInstance.GetAllUserTracks(username)
	if err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	if err := render.Render(w, req, tracks); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}
}

func addUserTrack(w http.ResponseWriter, req *http.Request) {
	auth := req.Context().Value("auth").(map[string]string)
	username, ok := auth["username"]
	if !ok {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	action := models.AddTrack{}
	if err := render.Bind(req, &action); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}

	if err := dbInstance.AddTrackAction(username, action.Type, action.SongId); err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}
}

func removeUserTrack(w http.ResponseWriter, req *http.Request) {
	auth := req.Context().Value("auth").(map[string]string)
	username, ok := auth["username"]
	if !ok {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	action := models.RemoveTrack{}
	if err := render.Bind(req, &action); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}

	if err := dbInstance.AddTrackAction(username, action.Type, action.SongId); err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}
}

func addUserTrackMetadata(w http.ResponseWriter, req *http.Request) {
	auth := req.Context().Value("auth").(map[string]string)
	username, ok := auth["username"]
	if !ok {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	action := models.Action{}
	if err := render.Bind(req, &action); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}

	if err := dbInstance.AddTrackAction(username, action.Type, action.SongId); err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}
}

func getRecommendedTracks(w http.ResponseWriter, req *http.Request) {
	auth := req.Context().Value("auth").(map[string]string)
	username, ok := auth["username"]
	if !ok {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	tracks, err := dbInstance.GetUserRecommendation(username)
	if err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	if err := render.Render(w, req, tracks); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}
}
