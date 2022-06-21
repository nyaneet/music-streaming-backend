package handler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nyaneet/music-streaming-backend/db"
	jwtauth "github.com/nyaneet/music-streaming-backend/jwt-auth"
)

func auth(router chi.Router) {
	router.Post("/sign_in", signIn)
	router.Post("/sign_up", signUp)
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

	expirationTime := time.Now().Add(jwtauth.JWT_LIFETIME)
	claims := &jwtauth.Claims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(jwtauth.JWTKey)
	if err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}

	err = render.Render(w, req, &jwtauth.Payload{
		Token:    token,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	})
	if err != nil {
		render.Render(w, req, ErrInternalServerError)
		return
	}
}

// TODO
func signUp(w http.ResponseWriter, req *http.Request) {

}
