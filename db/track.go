package db

import (
	"database/sql"

	"github.com/nyaneet/music-streaming-backend/models"
)

func (db Database) GetAllTracks() (*models.TrackList, error) {
	tracks := &models.TrackList{}
	query := `SELECT
				s.song_id as song_id,
				s.title as title,
				s.explicit AS explicit,
				s.duration AS duration,
				al.album_id AS album_id,
				al.name AS album_name,
				al.type AS album_type,
				al.year AS album_year,
				ar.artist_id AS artist_id,
				ar.name AS artist_name,
				ar.description AS artist_description
			FROM
				songs AS s
				JOIN albums_songs als ON s.song_id = als.song_id
				JOIN albums al ON al.album_id = als.album_id
				JOIN artists ar ON al.artist_id = ar.artist_id
			ORDER BY song_id DESC;`
	rows, err := db.Conn.Query(query)
	if err != nil {
		return tracks, err
	}
	defer rows.Close()

	for rows.Next() {
		var track models.Track
		err := rows.Scan(
			&track.Id,
			&track.Name,
			&track.Explicit,
			&track.Duration,
			&track.Album.Id,
			&track.Album.Name,
			&track.Album.Type,
			&track.Album.Year,
			&track.Album.Artist.Id,
			&track.Album.Artist.Name,
			&track.Album.Artist.Description,
		)
		if err != nil {
			return tracks, err
		}
		tracks.Tracks = append(tracks.Tracks, track)
	}

	return tracks, nil
}

func (db Database) GetTrackById(trackId int) (models.Track, error) {
	track := models.Track{}

	query := `SELECT
			    s.song_id as song_id,
			    s.title as title,
			    s.explicit AS explicit,
			    s.duration AS duration,
			    al.album_id AS album_id,
			    al.name AS album_name,
			    al.type AS album_type,
			    al.year AS album_year,
			    ar.artist_id AS artist_id,
			    ar.name AS artist_name,
			    ar.description AS artist_description
			FROM
			    songs AS s
			    JOIN albums_songs als ON s.song_id = als.song_id
			    JOIN albums al ON al.album_id = als.album_id
			    JOIN artists ar ON al.artist_id = ar.artist_id
			WHERE
			    s.song_id = $1;`
	row := db.Conn.QueryRow(query, trackId)

	if err := row.Scan(
		&track.Id,
		&track.Name,
		&track.Explicit,
		&track.Duration,
		&track.Album.Id,
		&track.Album.Name,
		&track.Album.Type,
		&track.Album.Year,
		&track.Album.Artist.Id,
		&track.Album.Artist.Name,
		&track.Album.Artist.Description,
	); err != nil {
		if err == sql.ErrNoRows {
			return track, ErrNoMatch
		}
		return track, err
	}

	return track, nil
}
