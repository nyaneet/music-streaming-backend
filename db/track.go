package db

import (
	"database/sql"

	"github.com/nyaneet/music-streaming-backend/models"
)

func (db Database) GetAllTracks() (*models.TrackList, error) {
	tracks := &models.TrackList{}

	rows, err := db.Conn.Query("SELECT * FROM songs ORDER BY song_id DESC;")
	if err != nil {
		return tracks, err
	}

	for rows.Next() {
		var track models.Track
		err := rows.Scan(&track.Id, &track.Name, &track.Explicit, &track.Duration)
		if err != nil {
			return tracks, err
		}
		tracks.Tracks = append(tracks.Tracks, track)
	}

	return tracks, nil
}

func (db Database) GetTrackById(trackId int) (models.Track, error) {
	track := models.Track{}

	query := `SELECT * FROM songs WHERE song_id = $1;`
	row := db.Conn.QueryRow(query, trackId)
	switch err := row.Scan(&track.Id, &track.Name, &track.Explicit, &track.Duration); err {
	case sql.ErrNoRows:
		return track, ErrNoMatch
	default:
		return track, err
	}
}
