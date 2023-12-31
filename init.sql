CREATE TABLE IF NOT EXISTS users (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	email TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
);

INSERT INTO users (
	id, name, email, password
) VALUES ( 1, 'Ivan', 'example@mail.com', 'password' );

CREATE TABLE IF NOT EXISTS movies (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	year INTEGER NOT NULL,
	brief TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	added TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO movies(id, title, year, brief, description, added)
VALUES
	(
		1, 'Inglourious Basterds', 2009,
		'Kristoph Waltz is amazing',
		"In Nazi-occupied France during World War II, a plan to assassinate Nazi leaders by a group of Jewish U.S. soldiers coincides with a theatre owner's vengeful plans for the same.",
		date('now', '-1 day')
	),
	(
		2, 'Once Upon a Time in Holywood', 2019,
		'So many actors!!',
		"A faded television actor and his stunt double strive to achieve fame and success in the final years of Hollywood's Golden Age in 1969 Los Angeles.",
		CURRENT_TIMESTAMP
	)
;

CREATE TABLE IF NOT EXISTS genres (
	movie_id INTEGER NOT NULL,
	genre TEXT NOT NULL
);


INSERT INTO genres (movie_id, genre)
VALUES
	( 1, 'Action' ),
	( 1, 'Drama' );


INSERT INTO genres (movie_id, genre)
VALUES
	( 2, 'Drama' );

CREATE TABLE IF NOT EXISTS people (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	birth_year INTEGER,
	death_year INTEGER,
	occupation TEXT NOT NULL
);

INSERT INTO people (id, name, occupation)
VALUES
	( 1, 'Quentin Tarantino', 'Director' ),
	( 2, 'Margot Robbie', 'Actor' )
;

CREATE TABLE IF NOT EXISTS crew (
	movie_id INTEGER NOT NULL,
	person_id INTEGER NOT NULL,
	role TEXT NOT NULL
);

INSERT INTO crew (movie_id, person_id, role)
VALUES
	( 2, 1, 'director' ),
	( 2, 1, 'writer' ),
	( 2, 2, 'actor' )
;

/*
select
	movies.id, movies.title,
	directors.id as director_id, directors.name as director_name
from movies
left join
    (select id, name from people) as directors
    on directors.id == (
        select person_id
        from movie_crew
        where movie_id = movies.id and role == 'director'
    )
where movies.id == 2
;
 */

