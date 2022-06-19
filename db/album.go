package db

import (
	"database/sql"

	"github.com/nyaneet/music-streaming-backend/models"
)

func (db Database) GetAllAlbums() (*models.AlbumList, error) {
	albums := &models.AlbumList{}
	empty := ""

	rows, err := db.Conn.Query("SELECT * FROM albums ORDER BY album_id DESC;")
	if err != nil {
		return albums, err
	}
	defer rows.Close()

	for rows.Next() {
		var album models.Album
		err := rows.Scan(&album.Id, &album.Name, &empty, &album.Type, &album.Year, &empty)
		if err != nil {
			return albums, err
		}
		albums.Albums = append(albums.Albums, album)
	}

	return albums, nil
}

func (db Database) GetAlbumById(artistId int) (models.Album, error) {
	album := models.Album{}
	empty := ""

	query := `SELECT * FROM albums WHERE album_id = $1;`
	row := db.Conn.QueryRow(query, artistId)
	switch err := row.Scan(&album.Id, &album.Name, &empty, &album.Type, &album.Year, &empty); err {
	case sql.ErrNoRows:
		return album, ErrNoMatch
	default:
		return album, err
	}
}
