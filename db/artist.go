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
	defer rows.Close()

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

	if err := row.Scan(&artist.Id, &artist.Name, &artist.Description); err != nil {
		if err == sql.ErrNoRows {
			return artist, ErrNoMatch
		}
		return artist, err
	}

	return artist, nil
}
