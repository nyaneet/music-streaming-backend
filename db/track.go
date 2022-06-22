package db

import (
	"database/sql"

	"github.com/nyaneet/music-streaming-backend/models"
)

func (db Database) GetAllTracks() (*models.TrackList, error) {
	tracks := &models.TrackList{}
	query := `
	SELECT
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
			&track.Id, &track.Name, &track.Explicit,
			&track.Duration, &track.Album.Id, &track.Album.Name,
			&track.Album.Type, &track.Album.Year, &track.Album.Artist.Id,
			&track.Album.Artist.Name, &track.Album.Artist.Description,
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

	query := `
	SELECT
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
		&track.Id, &track.Name, &track.Explicit,
		&track.Duration, &track.Album.Id, &track.Album.Name,
		&track.Album.Type, &track.Album.Year, &track.Album.Artist.Id,
		&track.Album.Artist.Name, &track.Album.Artist.Description,
	); err != nil {
		if err == sql.ErrNoRows {
			return track, ErrNoMatch
		}
		return track, err
	}

	return track, nil
}

func (db Database) GetAllUserTracks(username string) (*models.TrackList, error) {
	tracks := &models.TrackList{}

	query := `
	SELECT
		DISTINCT ON (a.date, s.song_id)
		s.song_id AS song_id,
		s.title AS title,
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
		actions a
		RIGHT JOIN songs s ON a.song_id = s.song_id
		JOIN users u ON a.user_id = u.user_id
		JOIN albums_songs als ON s.song_id = als.song_id
		JOIN albums al ON al.album_id = als.album_id
		JOIN artists ar ON al.artist_id = ar.artist_id
	WHERE
		u.nickname = $1
		AND a.type = (enum_range(null::action_type))[1]
		AND NOT EXISTS (
			SELECT 1
			FROM actions a2
			WHERE a.user_id = a2.user_id
			AND a.song_id = a2.song_id
			AND (a2.type = (enum_range(null::action_type))[3] OR a2.type = (enum_range(null::action_type))[1])
			AND a.date < a2.date )
	ORDER BY a.date DESC;`
	rows, err := db.Conn.Query(query, username)
	if err != nil {
		return tracks, err
	}

	defer rows.Close()
	for rows.Next() {
		var track = models.Track{}

		err := rows.Scan(
			&track.Id, &track.Name, &track.Explicit,
			&track.Duration, &track.Album.Id, &track.Album.Name,
			&track.Album.Type, &track.Album.Year, &track.Album.Artist.Id,
			&track.Album.Artist.Name, &track.Album.Artist.Description,
		)
		if err != nil {
			return tracks, err
		}

		tracks.Tracks = append(tracks.Tracks, track)
	}

	return tracks, nil
}

func (db Database) AddTrack(track models.Track, username string) error {
	var (
		query    string
		row      *sql.Row
		userId   int
		artistId int
		trackId  int
	)

	// making transaction and defer a rollback in case anything fails
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// check that albumId is correct
	query = `SELECT user_id, artist_id FROM users WHERE nickname = $1;`
	row = tx.QueryRow(query, username)
	if err := row.Scan(&userId, &track.Album.Artist.Id); err != nil {
		if err == sql.ErrNoRows {
			return ErrNoMatch
		}
		return err
	}

	query = `SELECT artist_id FROM albums WHERE album_id = $1;`
	row = tx.QueryRow(query, track.Album.Id)
	if err := row.Scan(&artistId); err != nil {
		if err == sql.ErrNoRows {
			return ErrNoMatch
		}
		return err
	}

	if artistId != track.Album.Artist.Id {
		return ErrNotAllowed
	}

	query = `
	INSERT INTO songs
		(title, explicit, duration) 
	VALUES 
		($1, $2, $3)
	RETURNING
		song_id;`
	if err := tx.QueryRow(query, track.Name, track.Explicit, track.Duration).Scan(&trackId); err != nil {
		return err
	}

	query = `
	INSERT INTO albums_songs
		(album_id, song_id) 
	VALUES 
		($1, $2)
	RETURNING
		song_id;`
	if err := tx.QueryRow(query, track.Album.Id, trackId).Scan(&trackId); err != nil {
		return err
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
