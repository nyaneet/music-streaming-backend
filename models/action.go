package models

import (
	"fmt"
	"net/http"
)

var ActionTypes = map[string]string{
	"ADD_TRACK":            "ADD_SONG_TO_LIBRARY",
	"REMOVE_TRACK":         "REMOVE_SONG_FROM_LIBRARY",
	"PLAY_TRACK":           "LISTEN",
	"DISLIKE_TRACK":        "ADD_SONG_TO_DISLIKED",
	"DISLIKE_CANCEL_TRACK": "REMOVE_SONG_FROM_DISLIKED",
}

type Action struct {
	Type   string `json:"type"`
	SongId int    `json:"song_id"`
}

func validateSongId(songId int) error {
	if songId == 0 {
		return fmt.Errorf("Song id is required.")
	}
	if songId < 0 {
		return fmt.Errorf("Invalid song id.")
	}

	return nil
}

func (a *Action) Bind(req *http.Request) error {
	if a.Type == "" {
		return fmt.Errorf("Action type is required.")
	}
	actionType, ok := ActionTypes[a.Type]
	if !ok {
		return fmt.Errorf("Invalid action type.")
	}
	a.Type = actionType

	if err := validateSongId(a.SongId); err != nil {
		return err
	}

	return nil
}

func (a *Action) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

type AddTrack struct {
	Action
}

func (a *AddTrack) Bind(req *http.Request) error {
	a.Type = ActionTypes["ADD_TRACK"]

	if err := validateSongId(a.SongId); err != nil {
		return err
	}

	return nil
}

type RemoveTrack struct {
	Action
}

func (a *RemoveTrack) Bind(req *http.Request) error {
	a.Type = ActionTypes["REMOVE_TRACK"]

	if err := validateSongId(a.SongId); err != nil {
		return err
	}

	return nil
}
