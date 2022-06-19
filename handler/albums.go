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

var albumIdKey = "albumId"

func albums(router chi.Router) {
	router.Get("/", getAlbums)
	router.Route("/{id}", func(router chi.Router) {
		router.Use(AlbumCtx)
		router.Get("/", getAlbum)
	})
}

func AlbumCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		albumId := chi.URLParam(req, "id")
		if albumId == "" {
			render.Render(w, req, ErrorRenderer(fmt.Errorf("Id is required")))
			return
		}

		id, err := strconv.Atoi(albumId)
		if err != nil {
			render.Render(w, req, ErrorRenderer(fmt.Errorf("Invalid Id")))
			return
		}

		ctx := context.WithValue(req.Context(), albumIdKey, id)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func getAlbum(w http.ResponseWriter, req *http.Request) {
	albumId := req.Context().Value(albumIdKey).(int)

	album, err := dbInstance.GetAlbumById(albumId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, req, ErrNotFound)
		} else {
			render.Render(w, req, ErrorRenderer(err))
		}
		return
	}

	err = render.Render(w, req, &album)
	if err != nil {
		return
	}
}

func getAlbums(w http.ResponseWriter, req *http.Request) {
	albums, err := dbInstance.GetAllAlbums()
	if err != nil {
		return
	}

	err = render.Render(w, req, albums)
	if err != nil {
		render.Render(w, req, ErrorRenderer(err))
	}
}
