CREATE TABLE IF NOT EXISTS users (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	email TEXT NOT NULL,
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

INSERT INTO movies(title, year, brief, added)
VALUES
	('Inglourious Basterds', 2009, 'Kristoph Waltz is amazing', date('now', '-1 day')),
	('Once Upon a Time in Holywood', 2019, 'So many actors!!', CURRENT_TIMESTAMP)
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
