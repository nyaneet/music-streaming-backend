package models

import (
	"net/http"
)

type Track struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Album    Album  `json:"album"`
	Explicit bool   `json:"explicit"`
	Duration int    `json:"duration"`
}

func (t *Track) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func (t *Track) Bind(req *http.Request) error {
	return nil
}

type TrackList struct {
	Tracks []Track `json:"tracks"`
}

func (t *TrackList) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

type SearchQuery struct {
	Query string `json:"query"`
}

func (sq *SearchQuery) Bind(req *http.Request) error {
	return nil
}
