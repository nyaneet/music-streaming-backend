package db

import (
	"database/sql"

	"github.com/nyaneet/music-streaming-backend/models"
)

func (db Database) GetAllUsers() (*models.UserList, error) {
	users := &models.UserList{}
	query := `
	SELECT
		user_id, nickname, password,
		email, type, banned, artist_id
	FROM
		users
	ORDER BY user_id DESC;`
	rows, err := db.Conn.Query(query)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.Id, &user.Username, &user.Password,
			&user.Email, &user.Role, &user.Banned,
			&user.ArtistId,
		)
		if err != nil {
			return users, err
		}
		users.Users = append(users.Users, user)
	}

	return users, nil
}

func (db Database) GetUserById(userId int) (models.User, error) {
	user := models.User{}

	query := `
	SELECT
		user_id, nickname, password,
		email, type, banned, artist_id
	FROM
		users
	WHERE
		user_id = $1;`
	row := db.Conn.QueryRow(query, userId)

	if err := row.Scan(
		&user.Id, &user.Username, &user.Password,
		&user.Email, &user.Role, &user.Banned,
		&user.ArtistId,
	); err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNoMatch
		}
		return user, err
	}

	return user, nil
}

func (db Database) GetUserByName(userName string) (models.User, error) {
	user := models.User{}

	query := `
	SELECT
		user_id, nickname, password,
		email, type, banned, artist_id
	FROM
		users
	WHERE
		nickname = $1;`
	row := db.Conn.QueryRow(query, userName)

	if err := row.Scan(
		&user.Id, &user.Username, &user.Password,
		&user.Email, &user.Role, &user.Banned,
		&user.ArtistId,
	); err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNoMatch
		}
		return user, err
	}

	return user, nil
}

func (db Database) AddUser(data models.RegistrationData) error {
	var (
		query    string
		userId   int
		artistId int
	)

	// making transaction and defer a rollback in case anything fails
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if data.Role == "ARTIST" {
		query = `
		INSERT INTO artists
			(name, description)
		VALUES 
			($1, $2)
		RETURNING
			artist_id;`
		if err := tx.QueryRow(query, data.ArtistName, data.ArtistDescription).Scan(&artistId); err != nil {
			return err
		}
	}

	query = `
	INSERT INTO users
		(nickname, password, email, country, date, type, banned, artist_id)
	VALUES 
		($1, $2, $3, 'ru', NOW(), $4, FALSE, $5)
	RETURNING
		user_id;`
	validArtistId := map[bool]interface{}{true: artistId, false: nil}
	if err := tx.QueryRow(
		query, data.Username, data.Password,
		data.Email, data.Role, validArtistId[artistId != 0],
	).Scan(&userId); err != nil {
		return err
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
