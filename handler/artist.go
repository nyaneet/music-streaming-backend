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

var artistIdKey = "artistId"

var artistUrlParams = map[string]string{
	"artistIdKey": "artistId",
}

func artists(router chi.Router) {
	router.Get("/", getArtists)
	router.Route("/{id}", func(router chi.Router) {
		router.Use(artistCtx)
		router.Get("/", getArtist)
	})
}

func artistCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		artistId := chi.URLParam(req, "id")
		if artistId == "" {
			render.Render(w, req, ErrorRenderer(fmt.Errorf("Id is required.")))
			return
		}

		id, err := strconv.Atoi(artistId)
		if err != nil {
			render.Render(w, req, ErrorRenderer(fmt.Errorf("Invalid Id.")))
			return
		}

		ctx := context.WithValue(req.Context(), artistIdKey, id)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func getArtist(w http.ResponseWriter, req *http.Request) {
	artistId := req.Context().Value(artistIdKey).(int)

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
