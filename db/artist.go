package db

import (
	"database/sql"

	"github.com/nyaneet/music-streaming-backend/models"
)

func (db Database) GetAllArtists() (*models.ArtistList, error) {
	artists := &models.ArtistList{}

	rows, err := db.Conn.Query("SELECT * FROM artists ORDER BY artist_id DESC;")
	if err != nil {
		return artists, err
	}

	for rows.Next() {
		var artist models.Artist
		err := rows.Scan(&artist.Id, &artist.Name, &artist.Description)
		if err != nil {
			return artists, err
		}
		artists.Artists = append(artists.Artists, artist)
	}

	return artists, nil
}

func (db Database) GetArtistById(artistId int) (models.Artist, error) {
	artist := models.Artist{}

	query := `SELECT * FROM artists WHERE artist_id = $1;`
	row := db.Conn.QueryRow(query, artistId)
	switch err := row.Scan(&artist.Id, &artist.Name, &artist.Description); err {
	case sql.ErrNoRows:
		return artist, ErrNoMatch
	default:
		return artist, err
	}
}
