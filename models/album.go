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

func (*Album) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

type AlbumList struct {
	Albums []Album `json:"albums"`
}

func (*AlbumList) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}
