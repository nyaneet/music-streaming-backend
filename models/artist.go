package models

import (
	"database/sql"
	"net/http"
)

type Artist struct {
	Id          int            `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}

func (*Artist) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

type ArtistList struct {
	Artists []Artist `json:"artists"`
}

func (*ArtistList) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}
