package models

import (
	"fmt"
	"net/http"
)

type Album struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Artist Artist `json:"artist"`
	Type   string `json:"type"`
	Year   int    `json:"year"`
}

func validateAlbumType(albumType string) error {
	switch albumType {
	case
		"ALBUM",
		"EP",
		"SINGLE":
		return nil
	default:
		return fmt.Errorf("Invalid album type.")
	}
}

func (a *Album) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func (a *Album) Bind(req *http.Request) error {
	if err := validateAlbumType(a.Type); err != nil {
		return err
	}

	return nil
}

type AlbumList struct {
	Albums []Album `json:"albums"`
}

func (a *AlbumList) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}
