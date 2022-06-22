package db

import (
	"database/sql"

	"github.com/nyaneet/music-streaming-backend/models"
)

func (db Database) GetAllUsers() (*models.UserList, error) {
	users := &models.UserList{}
	query := `
	SELECT
		user_id,
		nickname,
		password,
		email,
		type,
		banned,
		artist_id
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
			&user.Id,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.Role,
			&user.Banned,
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
		user_id,
		nickname,
		password,
		email,
		type,
		banned,
		artist_id
	FROM
		users
	WHERE
		user_id = $1;`
	row := db.Conn.QueryRow(query, userId)

	if err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
		&user.Banned,
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
		user_id,
		nickname,
		password,
		email,
		type,
		banned,
		artist_id
	FROM
		users
	WHERE
		nickname = $1;`
	row := db.Conn.QueryRow(query, userName)

	if err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
		&user.Banned,
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

	query = `
	INSERT INTO artists
		(name, description)
	VALUES 
		($1, $2)
	RETURNING
		artist_id;`
	if data.Role == "ARTIST" {
		if err := db.Conn.QueryRow(query, data.ArtistName, data.ArtistDescription).Scan(&artistId); err != nil {
			return err
		}
	}

	query = `
	INSERT INTO users
		(nickname, password, email, country, date, type, banned)
	VALUES 
		($1, $2, $3, 'ru', NOW(), $4, FALSE)
	RETURNING
		user_id;`
	if err := db.Conn.QueryRow(query, data.Username, data.Password, data.Email, data.Role).Scan(&userId); err != nil {
		db.Conn.QueryRow(`DELETE FROM artists WHERE artist_id = $1;`, artistId)
		return err
	}

	return nil
}
