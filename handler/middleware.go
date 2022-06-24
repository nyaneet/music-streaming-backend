package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	jwtauth "github.com/nyaneet/music-streaming-backend/jwt-auth"
)

func extractId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		if id == "" {
			render.Render(w, req, ErrorRenderer(fmt.Errorf("Id is required.")))
			return
		}

		idValue, err := strconv.Atoi(id)
		if err != nil || idValue <= 0 {
			render.Render(w, req, ErrorRenderer(fmt.Errorf("Invalid id.")))
			return
		}

		ctx := context.WithValue(req.Context(), "id", idValue)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
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

func packAuthInfo(req *http.Request, claims *jwtauth.Claims) context.Context {
	ctxValue := map[string]string{
		"username": claims.Username,
		"role":     claims.Role,
	}
	return context.WithValue(req.Context(), "auth", ctxValue)
}

func isAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, claims, err := extractTokenWithClaims(req)
		if err != nil || !token.Valid {
			render.Render(w, req, ErrUnauthorized)
			return
		}

		ctx := packAuthInfo(req, claims)
		next.ServeHTTP(w, req.WithContext(ctx))
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

		ctx := packAuthInfo(req, claims)
		next.ServeHTTP(w, req.WithContext(ctx))
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

		ctx := packAuthInfo(req, claims)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func notBanned(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, claims, err := extractTokenWithClaims(req)
		if err != nil || !token.Valid {
			render.Render(w, req, ErrUnauthorized)
			return
		}

		banned, err := dbInstance.CheckBan(claims.Username)
		if err != nil {
			render.Render(w, req, ErrInternalServerError)
			return
		}
		if banned {
			render.Render(w, req, ErrNotAllowed)
			return
		}

		ctx := packAuthInfo(req, claims)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
