package db

import "database/sql"

func (db Database) AddTrackAction(username, actionType string, trackId int) error {
	var (
		userId   int
		actionId int
	)

	query := `
	SELECT
		user_id
	FROM
		users
	WHERE
		nickname = $1;`

	row := db.Conn.QueryRow(query, username)
	if err := row.Scan(&userId); err != nil {
		if err == sql.ErrNoRows {
			return ErrNoMatch
		}
		return err
	}

	insertQuery := `
	INSERT INTO actions
		(date, user_id, type, song_id) 
	VALUES 
		(NOW(), $1, $2, $3)
	RETURNING
		action_id;`
	if err := db.Conn.QueryRow(insertQuery, userId, actionType, trackId).Scan(&actionId); err != nil {
		return err
	}

	return nil
}
