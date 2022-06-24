--- easy

--- 1

SELECT
    us.country AS country,
    count(*) AS users_amount
FROM users AS us
GROUP BY us.country
ORDER BY us.count DESC
LIMIT 5;

--- 2.1 (1й вариант, PARENT_CONTORL = FALSE) (отображаются все песни в альбоме)

SELECT
    s.title AS title,
    s.explicit AS explicit,
    s.duration AS duration
FROM songs AS s
WHERE title ILIKE '%one%'
ORDER BY title ASC;

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
    s.title ILIKE '%Alone in%'
    OR al.name ILIKE '%Alone in%'
    OR ar.name ILIKE '%Alone in%'

ORDER BY song_id DESC;

--- 2.1 (2й вариант, PARENT_CONTORL = TRUE) (песни со значением explicit = true скрыты)

SELECT
    s.title AS title,
    s.explicit AS explicit,
    s.duration AS duration
FROM songs AS s
WHERE title ILIKE '%one%'
AND CASE
        WHEN TRUE = TRUE THEN explicit <> TRUE
        ELSE TRUE
    END
ORDER BY title ASC;

--- 3 

SELECT
    al.name AS album,
    al.type AS type,
    al.date::timestamp::date AS released
FROM albums AS al
WHERE (CURRENT_DATE - al.date::timestamp::date) < 50
ORDER BY al.date DESC
LIMIT 3;

--- 4

SELECT
    count(*) as users_activity
FROM actions AS a
WHERE (CURRENT_DATE - a.date::timestamp::date) < 30;

--- medium

--- 1

SELECT
    DISTINCT a.name AS album,
    COUNT(*) OVER(PARTITION BY a.name) AS tracks_number,
    SUM(s.duration) OVER(PARTITION BY a.name) AS duration
FROM
    albums AS a
    JOIN albums_songs als ON a.album_id = als.album_id
    JOIN songs s ON s.song_id = als.song_id
    JOIN artists ar ON ar.artist_id = a.artist_id
WHERE
    ar.name = 'Normandie'
    AND a.name = 'White Flag';


--- Get track

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
    s.song_id = 1;

--- Get album

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
    al.album_id = 1;

--- Get User Tracks

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
    u.nickname = 'grumpyCat'
    AND a.type = (enum_range(null::action_type))[1]
    AND NOT EXISTS (
        SELECT 1
        FROM actions a2
        WHERE a.user_id = a2.user_id
        AND a.song_id = a2.song_id
        AND (a2.type = (enum_range(null::action_type))[3] OR a2.type = (enum_range(null::action_type))[1])
        AND a.date < a2.date )
ORDER BY a.date DESC;

--- Get all albums track

SELECT
    s.song_id as song_id,
    s.title as title,
    al.album_id AS album_id,
    al.name AS album_name,
    ar.artist_id AS artist_id,
    ar.name AS artist_name		
FROM
    songs AS s
    JOIN albums_songs als ON s.song_id = als.song_id
    JOIN albums al ON al.album_id = als.album_id
    JOIN artists ar ON al.artist_id = ar.artist_id
WHERE
    al.id = 5
ORDER BY song_id DESC;


--- 2

SELECT
    s.title AS title,
    s.explicit AS explicit,
    als.track_number AS number
FROM
    albums AS a
    JOIN albums_songs als ON a.album_id = als.album_id
    JOIN songs s ON s.song_id = als.song_id
    JOIN artists ar ON ar.artist_id = a.artist_id
WHERE ar.name = 'Beartooth'
    AND a.name = 'Aggressive'
ORDER BY als.track_number ASC;

--- 3

SELECT * FROM
    (SELECT
            s.title AS title,
            al.name AS album,
            COUNT(*) AS listening_amount,
            dense_rank() OVER (ORDER BY COUNT(*) DESC) AS rank
        FROM
            actions AS a
            JOIN songs AS s ON s.song_id = a.song_id
            JOIN albums_songs als ON als.song_id = s.song_id
            JOIN albums al ON al.album_id = als.album_id
        WHERE
            a.type = (enum_range(null::action_type))[10]
            AND a.date >= '2018-11-01'
        GROUP BY s.song_id, al.album_id
        ORDER BY listening_amount DESC) AS re
    WHERE re.rank <= 3;

--- hard

--- 1

SELECT
    DISTINCT ON (a.date, s.song_id)
    s.title AS title,
    s.explicit AS explicit,
    ar.name AS artist,
    al.name AS album,
    a.date::timestamp::date AS added,
    s.duration AS duration
FROM
    actions a
    RIGHT JOIN songs s ON a.song_id = s.song_id
    JOIN albums_songs als ON als.song_id = s.song_id
    JOIN albums al ON al.album_id = als.album_id
    JOIN users u ON a.user_id = u.user_id
    JOIN artists ar ON ar.artist_id = al.artist_id
WHERE
    u.nickname = 'hexsixzeros'
    AND a.type = (enum_range(null::action_type))[1]
    AND NOT EXISTS (
        SELECT 1
        FROM actions a2
        WHERE a.user_id = a2.user_id
        AND a.song_id = a2.song_id
        AND (a2.type = (enum_range(null::action_type))[3] OR a2.type = (enum_range(null::action_type))[1])
        AND a.date < a2.date )
ORDER BY a.date DESC;

--- 2

WITH
    public_playlists AS (
    SELECT
        p.playlist_id
    FROM
        actions a
        RIGHT JOIN playlists p ON a.playlist_id = p.playlist_id
    WHERE
        a.type = (enum_range(null::action_type))[6]
        AND NOT EXISTS (
            SELECT 1
            FROM actions a2
            WHERE a.playlist_id = a2.playlist_id
            AND (a2.type = (enum_range(null::action_type))[7] OR a2.type = (enum_range(null::action_type))[6])
            AND a.date < a2.date) )        
SELECT
    p.name AS playlist,
    COUNT(*) AS following
FROM
    actions a
    RIGHT JOIN playlists p ON a.playlist_id = p.playlist_id
WHERE
    a.playlist_id IN (SELECT playlist_id FROM public_playlists)
    AND a.type = (enum_range(null::action_type))[8]
    AND NOT EXISTS (
            SELECT 1
            FROM actions a2
            WHERE a.playlist_id = a2.playlist_id
            AND a.user_id = a2.user_id
            AND (a2.type = (enum_range(null::action_type))[8] OR a2.type = (enum_range(null::action_type))[9])
            AND a.date < a2.date )
GROUP BY p.name, p.playlist_id
ORDER BY following DESC
LIMIT 3;

--- 3

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
        u.nickname = 'grumpyCat'
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
        u.nickname = 'hexsixzeros'
        AND a.type = (enum_range(null::action_type))[1]
        AND NOT EXISTS (
            SELECT 1
            FROM actions a2
            WHERE a.user_id = a2.user_id
            AND a.song_id = a2.song_id
            AND (a2.type = (enum_range(null::action_type))[3] OR a2.type = (enum_range(null::action_type))[1])
            AND a.date < a2.date )
    ORDER BY a.date DESC )
    SELECT 
        ((COUNT(ful.user_id) + COUNT(sul.user_id) - COUNT(*))::real / (1 + COUNT(ful.user_id))*100)::numeric::integer AS first,
        ((COUNT(ful.user_id) + COUNT(sul.user_id) - COUNT(*))::real / (1 + COUNT(sul.user_id))*100)::numeric::integer AS second
    FROM first_user_library ful FULL OUTER JOIN second_user_library sul ON ful.song_id = sul.song_id;



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
    u.nickname = 'grumpyCat'
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
    u.nickname = 'hexsixzeros'
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
        u.nickname = 'hexsixzeros'
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
            WHERE a3.user_id = 8
            AND a.song_id = a3.song_id
            AND a3.type = (enum_range(null::action_type))[11])
    ORDER BY a.date DESC)
SELECT 
    ((COUNT(ful.user_id) + COUNT(sul.user_id) - COUNT(*))::real / (1 + COUNT(ful.user_id) + COUNT(d.user_id))*100)::numeric::integer AS similarity
FROM first_user_library ful FULL OUTER JOIN second_user_library sul ON ful.song_id = sul.song_id FULL OUTER JOIN disliked d ON ful.song_id = d.song_id;



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
        u.user_id = 63
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
            WHERE a3.user_id = 8
            AND a.song_id = a3.song_id
            AND a3.type = (enum_range(null::action_type))[11])
    ORDER BY a.date DESC;