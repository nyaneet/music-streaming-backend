package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nyaneet/music-streaming-backend/db"
	"github.com/nyaneet/music-streaming-backend/models"
)

func mydiscography(router chi.Router) {
	router.Put("/tracks", addArtistTrack)
	router.Delete("/tracks", removeArtistTrack)
	router.Put("/albums", addArtistAlbum)
	router.Delete("/albums", removeArtistAlbum)
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

}
