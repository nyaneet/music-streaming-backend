package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nyaneet/music-streaming-backend/db"
	jwtauth "github.com/nyaneet/music-streaming-backend/jwt-auth"
	"github.com/nyaneet/music-streaming-backend/models"
)

func auth(router chi.Router) {
	router.Post("/signin", signIn)
	router.Post("/signup", signUp)
}

func signIn(w http.ResponseWriter, req *http.Request) {
	credentials := &jwtauth.Credentials{}
	if err := render.Bind(req, credentials); err != nil {
		render.Render(w, req, ErrBadRequest)
		return
	}

	user, err := dbInstance.GetUserByName(credentials.Username)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, req, ErrUnauthorized)
			return
		}
		render.Render(w, req, ErrInternalServerError)
		return
	}

	if user.Password != credentials.Password {
		render.Render(w, req, ErrUnauthorized)
		return
	}

	token, err := jwtauth.GetToken(user)
	if err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	render.Render(w, req, &jwtauth.Payload{
		Token:    token,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	})
}

func signUp(w http.ResponseWriter, req *http.Request) {
	data := models.RegistrationData{}

	if err := render.Bind(req, &data); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}

	if err := dbInstance.AddUser(data); err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	token, err := jwtauth.GetToken(models.User{Username: data.Username, Password: data.Password})
	if err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	render.Render(w, req, &jwtauth.Payload{
		Token:    token,
		Username: data.Username,
		Email:    data.Email,
		Role:     data.Role,
	})
}
