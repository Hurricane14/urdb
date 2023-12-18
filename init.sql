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

INSERT INTO movies(id, title, year, brief, added)
VALUES
	(1, 'Inglourious Basterds', 2009, 'Kristoph Waltz is amazing', date('now', '-1 day')),
	(2, 'Once Upon a Time in Holywood', 2019, 'So many actors!!', CURRENT_TIMESTAMP)
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
	( 1, 'Quentin Tarantino', 'Director' )
;

CREATE TABLE IF NOT EXISTS movie_crew (
	movie_id INTEGER NOT NULL,
	person_id INTEGER NOT NULL,
	role TEXT NOT NULL
);

INSERT INTO movie_crew (movie_id, person_id, role)
VALUES
	( 2, 1, 'director' ),
	( 2, 1, 'writer' )
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

