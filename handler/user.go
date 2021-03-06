package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nyaneet/music-streaming-backend/db"
)

func users(router chi.Router) {
	router.Get("/", getUsers)
	router.Post("/ban", banUser)
	router.Route("/{id}", func(router chi.Router) {
		router.Use(extractId)
		router.Get("/", getUser)
	})
}

func getUser(w http.ResponseWriter, req *http.Request) {
	userId := req.Context().Value("id").(int)

	user, err := dbInstance.GetUserById(userId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, req, ErrNotFound)
			return
		}
		render.Render(w, req, ErrInternalServerError)
		return
	}

	if err := render.Render(w, req, &user); err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}
}

func getUsers(w http.ResponseWriter, req *http.Request) {
	users, err := dbInstance.GetAllUsers()
	if err != nil {
		log.Fatal(err.Error())
		render.Render(w, req, ErrInternalServerError)
		return
	}

	if err := render.Render(w, req, users); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}
}

type Ban struct {
	Username string `json:"username"`
}

func (b *Ban) Bind(req *http.Request) error {
	return nil
}

func banUser(w http.ResponseWriter, req *http.Request) {
	ban := Ban{}
	if err := render.Bind(req, &ban); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}

	if err := dbInstance.BanUser(ban.Username); err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}
}
