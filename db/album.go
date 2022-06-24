package db

import (
	"database/sql"

	"github.com/nyaneet/music-streaming-backend/models"
)

func (db Database) GetAllAlbums() (*models.AlbumList, error) {
	albums := &models.AlbumList{}

	query := `
	SELECT
		al.album_id AS album_id,
		al.name AS album_name,
		al.type AS album_type,
		al.year AS album_year,
		ar.artist_id AS artist_id,
		ar.name AS artist_name,
		ar.description AS artist_description
	FROM
		albums AS al
		JOIN artists ar ON al.artist_id = ar.artist_id
	ORDER BY album_id DESC;`
	rows, err := db.Conn.Query(query)
	if err != nil {
		return albums, err
	}
	defer rows.Close()

	for rows.Next() {
		album := models.Album{}
		err := rows.Scan(
			&album.Id, &album.Name, &album.Type,
			&album.Year, &album.Artist.Id, &album.Artist.Name,
			&album.Artist.Description,
		)
		if err != nil {
			return albums, err
		}
		albums.Albums = append(albums.Albums, album)
	}

	return albums, nil
}

func (db Database) GetAlbumById(artistId int) (models.Album, error) {
	album := models.Album{}

	query := `
	SELECT
		al.album_id AS album_id,
		al.name AS album_name,
		al.type AS album_type,
		al.year AS album_year,
		ar.artist_id AS artist_id,
		ar.name AS artist_name,
		ar.description AS artist_description
	FROM
		albums AS al
		JOIN artists ar ON al.artist_id = ar.artist_id
	WHERE
		al.album_id = $1;`
	row := db.Conn.QueryRow(query, artistId)

	if err := row.Scan(
		&album.Id, &album.Name, &album.Type,
		&album.Year, &album.Artist.Id, &album.Artist.Name,
		&album.Artist.Description,
	); err != nil {
		if err == sql.ErrNoRows {
			return album, ErrNoMatch
		}
		return album, err
	}

	return album, nil
}

func (db Database) GetAllArtistAlbums(username string) (*models.AlbumList, error) {
	var (
		artistId int
		query    string
	)
	albums := &models.AlbumList{}

	query = `SELECT artist_id FROM users WHERE nickname = $1;`
	row := db.Conn.QueryRow(query, username)
	if err := row.Scan(&artistId); err != nil {
		if err == sql.ErrNoRows {
			return albums, ErrNoMatch
		}
		return albums, err
	}

	query = `
	SELECT
		al.album_id AS album_id,
		al.name AS album_name,
		al.type AS album_type,
		al.year AS album_year,
		ar.artist_id AS artist_id,
		ar.name AS artist_name,
		ar.description AS artist_description
	FROM
		albums AS al
		JOIN artists ar ON al.artist_id = ar.artist_id
	WHERE
		ar.artist_id = $1`
	rows, err := db.Conn.Query(query, artistId)
	if err != nil {
		return albums, err
	}

	defer rows.Close()
	for rows.Next() {
		album := models.Album{}
		err := rows.Scan(
			&album.Id, &album.Name, &album.Type,
			&album.Year, &album.Artist.Id, &album.Artist.Name,
			&album.Artist.Description,
		)
		if err != nil {
			return albums, err
		}
		albums.Albums = append(albums.Albums, album)
	}

	return albums, nil
}

func (db Database) AddAlbum(album models.Album, username string) error {
	var (
		query    string
		row      *sql.Row
		userId   int
		artistId int
		albumId  int
	)

	query = `
	SELECT
		user_id, artist_id
	FROM
		users
	WHERE
		nickname = $1;`

	row = db.Conn.QueryRow(query, username)
	if err := row.Scan(&userId, &artistId); err != nil {
		if err == sql.ErrNoRows {
			return ErrNoMatch
		}
		return err
	}

	query = `
	INSERT INTO albums
		(name, artist_id, type, year, date)
	VALUES 
		($1, $2, $3, $4, NOW())
	RETURNING
		album_id;`
	if err := db.Conn.QueryRow(query, album.Name, artistId, album.Type, album.Year).Scan(&albumId); err != nil {
		return err
	}

	return nil
}

func (db Database) RemoveAlbum(albumId int, username string) error {
	var (
		query           string
		row             *sql.Row
		userId          int
		artistId        int
		artistIdOfAlbum int
		tracksIds       []int
	)

	// making transaction and defer a rollback in case anything fails
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// check that username is correct
	query = `SELECT user_id, artist_id FROM users WHERE nickname = $1;`
	row = tx.QueryRow(query, username)
	if err := row.Scan(&userId, &artistId); err != nil {
		if err == sql.ErrNoRows {
			return ErrNoMatch
		}
		return err
	}

	query = `SELECT artist_id FROM albums WHERE album_id = $1;`
	row = tx.QueryRow(query, albumId)
	if err := row.Scan(&artistIdOfAlbum); err != nil {
		if err == sql.ErrNoRows {
			return ErrNoMatch
		}
		return err
	}

	if artistId != artistIdOfAlbum {
		return ErrNotAllowed
	}

	// remove all album's tracks
	query = `
	SELECT
		s.song_id as song_id	
	FROM
		songs AS s
		JOIN albums_songs als ON s.song_id = als.song_id
		JOIN albums al ON al.album_id = als.album_id
		JOIN artists ar ON al.artist_id = ar.artist_id
	WHERE
		al.album_id = $1
	ORDER BY song_id DESC;`
	rows, err := tx.Query(query, albumId)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var trackId int
		err := rows.Scan(&trackId)
		if err != nil {
			return err
		}
		tracksIds = append(tracksIds, trackId)
	}

	for _, trackId := range tracksIds {
		query = `DELETE FROM songs WHERE song_id = $1;`
		if _, err := tx.Exec(query, trackId); err != nil {
			if err == sql.ErrNoRows {
				return ErrNoMatch
			}
			return err
		}
		query = `DELETE FROM albums_songs WHERE song_id = $1;`
		if _, err := tx.Exec(query, trackId); err != nil {
			if err == sql.ErrNoRows {
				return ErrNoMatch
			}
			return err
		}
	}

	// remove album
	query = `DELETE FROM albums WHERE album_id = $1;`
	if _, err := tx.Exec(query, albumId); err != nil {
		if err == sql.ErrNoRows {
			return ErrNoMatch
		}
		return err
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
