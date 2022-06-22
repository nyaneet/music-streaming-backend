package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nyaneet/music-streaming-backend/db"
	"github.com/nyaneet/music-streaming-backend/models"
)

func mydiscography(router chi.Router) {
	router.Put("/tracks", addArtistTrack)
	router.Route("/tracks/{id}", func(router chi.Router) {
		router.Use(extractId)
		router.Delete("/", removeArtistTrack)
	})
	router.Put("/albums", addArtistAlbum)
	router.Route("/albums/{id}", func(router chi.Router) {
		router.Use(extractId)
		router.Delete("/", removeArtistAlbum)
	})
}

func addArtistTrack(w http.ResponseWriter, req *http.Request) {
	auth := req.Context().Value("auth").(map[string]string)
	username, ok := auth["username"]
	if !ok {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	track := models.Track{}
	if err := render.Bind(req, &track); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}

	if err := dbInstance.AddTrack(track, username); err != nil {
		if err == db.ErrNotAllowed {
			render.Render(w, req, ErrNotAllowed)
			return
		}
		render.Render(w, req, ErrInternalServerError)
		return
	}
}

func removeArtistTrack(w http.ResponseWriter, req *http.Request) {
	trackId := req.Context().Value("id").(int)
	auth := req.Context().Value("auth").(map[string]string)
	username, ok := auth["username"]
	if !ok {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	if err := dbInstance.RemoveTrack(trackId, username); err != nil {
		if err == db.ErrNotAllowed {
			render.Render(w, req, ErrNotAllowed)
			return
		}
		if err == db.ErrNoMatch {
			render.Render(w, req, ErrBadRequest)
			return
		}
		render.Render(w, req, ErrInternalServerError)
		return
	}
}

func addArtistAlbum(w http.ResponseWriter, req *http.Request) {
	auth := req.Context().Value("auth").(map[string]string)
	username, ok := auth["username"]
	if !ok {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	album := models.Album{}
	if err := render.Bind(req, &album); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}

	if err := dbInstance.AddAlbum(album, username); err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}
}

func removeArtistAlbum(w http.ResponseWriter, req *http.Request) {
	albumId := req.Context().Value("id").(int)
	auth := req.Context().Value("auth").(map[string]string)
	username, ok := auth["username"]
	if !ok {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	if err := dbInstance.RemoveAlbum(albumId, username); err != nil {
		log.Println(err)
		if err == db.ErrNotAllowed {
			render.Render(w, req, ErrNotAllowed)
			return
		}
		if err == db.ErrNoMatch {
			render.Render(w, req, ErrBadRequest)
			return
		}
		render.Render(w, req, ErrInternalServerError)
		return
	}
}
