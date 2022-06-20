package models

import (
	"database/sql"
	"net/http"
)

type User struct {
	Id       int           `json: "id"`
	Password string        `json: "password"`
	Username string        `json: "username"`
	Email    string        `json: "email"`
	Role     string        `json: "role"`
	Banned   bool          `json: "banned"`
	ArtistId sql.NullInt64 `json: "artist_id"`
}

func (*User) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

type UserList struct {
	Users []User `json "users"`
}

func (*UserList) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}
