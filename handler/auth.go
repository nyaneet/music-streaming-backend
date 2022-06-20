package handler

import (
	"fmt"
	"net/http"
	"strings"
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

func extractTokenWithClaims(req *http.Request) (*jwt.Token, *jwtauth.Claims, error) {
	token := &jwt.Token{}
	claims := &jwtauth.Claims{}

	bearer := req.Header.Get("Authorization")
	if len(bearer) <= 7 || strings.ToUpper(bearer[0:6]) != "BEARER" {
		return token, claims, fmt.Errorf("Invalid token.")
	}

	token, err := jwt.ParseWithClaims(bearer[7:], claims, func(token *jwt.Token) (interface{}, error) {
		return jwtauth.JWTKey, nil
	})
	return token, claims, err
}

func isAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, _, err := extractTokenWithClaims(req)
		if err != nil || !token.Valid {
			render.Render(w, req, ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func isArtist(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, claims, err := extractTokenWithClaims(req)
		if err != nil || !token.Valid {
			render.Render(w, req, ErrUnauthorized)
			return
		}

		if claims.Role != "ARTIST" {
			render.Render(w, req, ErrNotAllowed)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func isAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, claims, err := extractTokenWithClaims(req)
		if err != nil || !token.Valid {
			render.Render(w, req, ErrUnauthorized)
			return
		}

		if claims.Role != "ADMIN" {
			render.Render(w, req, ErrNotAllowed)
			return
		}

		next.ServeHTTP(w, req)
	})
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
