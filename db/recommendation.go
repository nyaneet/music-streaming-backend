package db

import (
	"database/sql"
	"time"

	"github.com/nyaneet/music-streaming-backend/models"
)

func subtract(tracks1, tracks2 *models.TrackList) *models.TrackList {
	sub := []models.Track{}
	seen := map[int]int{}
	for _, val := range tracks1.Tracks {
		seen[val.Id] = 1
	}
	for _, val := range tracks2.Tracks {
		_, ok := seen[val.Id]
		if !ok {
			sub = append(sub, val)
		}
	}

	return &models.TrackList{Tracks: sub}
}

const CACHE_LIFETIME = time.Second * 1

type Similarity struct {
	Value int
	Date  time.Time
}

var SimilarityCache = make(map[int](map[int]Similarity))

func (db Database) GetUserRecommendation(username string) (*models.TrackList, error) {
	recommendation := &models.TrackList{}
	var (
		query      string
		row        *sql.Row
		userId     int
		similarity int
	)

	tx, err := db.Conn.Begin()
	if err != nil {
		return recommendation, err
	}
	defer tx.Rollback()

	query = `SELECT user_id FROM users WHERE nickname = $1;`
	row = tx.QueryRow(query, username)
	if err := row.Scan(&userId); err != nil {
		if err == sql.ErrNoRows {
			return recommendation, ErrNoMatch
		}
		return recommendation, err
	}

	userIds := []int{}
	query = `SELECT user_id	FROM users WHERE user_id != $1;`
	rows, err := tx.Query(query, userId)
	if err != nil {
		return recommendation, err
	}
	defer rows.Close()

	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		if err != nil {
			return recommendation, err
		}
		userIds = append(userIds, userId)
	}

	maxSimilarity, anotherUserId := -100, 0
	for _, id := range userIds {
		query = `
		WITH
			first_user_library AS (
		SELECT
			DISTINCT ON (a.date, s.song_id)
			s.title AS title,
			s.song_id AS song_id,
			1 AS user_id,
			a.date::timestamp::date AS added
		FROM
			actions a
			RIGHT JOIN songs s ON a.song_id = s.song_id
			JOIN users u ON a.user_id = u.user_id
		WHERE
			u.user_id = $1
			AND a.type = (enum_range(null::action_type))[1]
			AND NOT EXISTS (
				SELECT 1
				FROM actions a2
				WHERE a.user_id = a2.user_id
				AND a.song_id = a2.song_id
				AND (a2.type = (enum_range(null::action_type))[3] OR a2.type = (enum_range(null::action_type))[1])
				AND a.date < a2.date )
		ORDER BY a.date DESC ),
			second_user_library AS (    
		SELECT
			DISTINCT ON (a.date, s.song_id)
			s.title AS title,
			s.song_id AS song_id,
			2 AS user_id,
			a.date::timestamp::date AS added
		FROM
			actions a
			RIGHT JOIN songs s ON a.song_id = s.song_id
			JOIN users u ON a.user_id = u.user_id
		WHERE
			u.user_id = $2
			AND a.type = (enum_range(null::action_type))[1]
			AND NOT EXISTS (
				SELECT 1
				FROM actions a2
				WHERE a.user_id = a2.user_id
				AND a.song_id = a2.song_id
				AND (a2.type = (enum_range(null::action_type))[3] OR a2.type = (enum_range(null::action_type))[1])
				AND a.date < a2.date )
		ORDER BY a.date DESC ),
			disliked AS(    
		SELECT
				DISTINCT ON (a.date, s.song_id)
				s.title AS title,
				s.song_id AS song_id,
				2 AS user_id,
				a.date::timestamp::date AS added
			FROM
				actions a
				RIGHT JOIN songs s ON a.song_id = s.song_id
				JOIN users u ON a.user_id = u.user_id
			WHERE
				u.user_id = $2
				AND a.type = (enum_range(null::action_type))[1]
				AND NOT EXISTS (
					SELECT 1
					FROM actions a2
					WHERE a.user_id = a2.user_id
					AND a.song_id = a2.song_id
					AND (a2.type = (enum_range(null::action_type))[3] OR a2.type = (enum_range(null::action_type))[1])
					AND a.date < a2.date )
				AND EXISTS (
					SELECT 1
					FROM actions a3
					WHERE a3.user_id = $1
					AND a.song_id = a3.song_id
					AND a3.type = (enum_range(null::action_type))[11])
			ORDER BY a.date DESC)
		SELECT 
			((COUNT(ful.user_id) + COUNT(sul.user_id) - COUNT(*))::real / (1 + COUNT(ful.user_id) + COUNT(d.user_id))*100)::numeric::integer AS similarity
		FROM first_user_library ful FULL OUTER JOIN second_user_library sul ON ful.song_id = sul.song_id FULL OUTER JOIN disliked d ON ful.song_id = d.song_id;`

		// check value in cache
		needToCalculate := false
		_, ok := SimilarityCache[userId]
		if !ok {
			SimilarityCache[userId] = make(map[int]Similarity)
			needToCalculate = true
		} else {
			similarity, ok := SimilarityCache[userId][id]
			if !ok {
				needToCalculate = true
			} else {
				passed := time.Now().Sub(similarity.Date)
				if passed > CACHE_LIFETIME {
					needToCalculate = true
				}
			}
		}

		if needToCalculate {
			row = tx.QueryRow(query, userId, id)
			if err := row.Scan(&similarity); err != nil {
				if err == sql.ErrNoRows {
					return recommendation, ErrNoMatch
				}
				return recommendation, err
			}
			SimilarityCache[userId][id] = Similarity{Value: similarity, Date: time.Now()}
		} else {
			similarity = SimilarityCache[userId][id].Value
		}

		if similarity > maxSimilarity {
			maxSimilarity = similarity
			anotherUserId = id
		}
	}

	if maxSimilarity == 0 {
		return recommendation, nil
	}

	query = `
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
		u.user_id = $1
		AND a.type = (enum_range(null::action_type))[1]
		AND NOT EXISTS (
			SELECT 1
			FROM actions a2
			WHERE a.user_id = a2.user_id
			AND a.song_id = a2.song_id
			AND (a2.type = (enum_range(null::action_type))[3] OR a2.type = (enum_range(null::action_type))[1])
			AND a.date < a2.date )
	ORDER BY a.date DESC;`

	tracks1 := &models.TrackList{}
	rows, err = tx.Query(query, userId)
	if err != nil {
		return recommendation, err
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
			return recommendation, err
		}

		tracks1.Tracks = append(tracks1.Tracks, track)
	}

	tracks2 := &models.TrackList{}
	rows, err = tx.Query(query, anotherUserId)
	if err != nil {
		return recommendation, err
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
			return recommendation, err
		}

		tracks2.Tracks = append(tracks2.Tracks, track)
	}
	if err := tx.Commit(); err != nil {
		return recommendation, err
	}

	recommendation = subtract(tracks1, tracks2)
	return recommendation, nil
}
