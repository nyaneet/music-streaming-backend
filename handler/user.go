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

var userIdKey = "userId"

func users(router chi.Router) {
	router.Get("/", getUsers)
	router.Route("/{id}", func(router chi.Router) {
		router.Use(userCtx)
		router.Get("/", getUser)
	})
}

func userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userId := chi.URLParam(req, "id")
		if userId == "" {
			render.Render(w, req, ErrorRenderer(fmt.Errorf("User Id is required.")))
			return
		}

		id, err := strconv.Atoi(userId)
		if err != nil {
			render.Render(w, req, ErrorRenderer(fmt.Errorf("Invalid user Id.")))
			return
		}

		ctx := context.WithValue(req.Context(), userIdKey, id)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func getUser(w http.ResponseWriter, req *http.Request) {
	userId := req.Context().Value(userIdKey).(int)

	user, err := dbInstance.GetUserById(userId)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, req, ErrNotFound)
		} else {
			render.Render(w, req, ErrorRenderer(err))
		}
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
		render.Render(w, req, ErrInternalServerError)
		return
	}

	if err := render.Render(w, req, users); err != nil {
		render.Render(w, req, ErrorRenderer(err))
		return
	}
}
