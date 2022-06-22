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
