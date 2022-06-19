package models

import (
	"net/http"
)

type Album struct {
	Id     int    `json: "id"`
	Name   string `json: "name"`
	Artist Artist `json: "artist"`
	Type   string `json: "type"`
	Year   int    `json: "year"`
}

type AlbumList struct {
	Albums []Album `json:"albums"`
}

func (*Album) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func (*AlbumList) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}
