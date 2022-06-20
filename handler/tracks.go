package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nyaneet/music-streaming-backend/db"
)

var trackIdKey = "trackId"

func tracks(router chi.Router) {
	router.Get("/", getTracks)
	router.Route("/{id}", func(router chi.Router) {
		router.Use(trackCtx)
		router.Get("/", getTrack)
	})
}

func trackCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		trackId := chi.URLParam(req, "id")
		if trackId == "" {
			render.Render(w, req, ErrorRenderer(fmt.Errorf("Track Id is required.")))
			return
		}

		id, err := strconv.Atoi(trackId)
		if err != nil {
			render.Render(w, req, ErrorRenderer(fmt.Errorf("Invalid track Id.")))
			return
		}

		ctx := context.WithValue(req.Context(), trackIdKey, id)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func getTrack(w http.ResponseWriter, req *http.Request) {
	trackId := req.Context().Value(trackIdKey).(int)

	track, err := dbInstance.GetTrackById(trackId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, req, ErrNotFound)
		} else {
			render.Render(w, req, ErrorRenderer(err))
		}
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
