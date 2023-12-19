package main

import (
	"context"
	"database/sql"
	"errors"
	"urdb/model"

	"github.com/mattn/go-sqlite3"
)

type usersRepository struct {
	db *sql.DB
}

type moviesRepository struct {
	db *sql.DB
}

var (
	movies *moviesRepository
	users  *usersRepository
)

func initRepositories() error {
	db, err := sql.Open("sqlite3", "sqlite.db")
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	movies = &moviesRepository{db}
	users = &usersRepository{db}

	return nil
}

func (r *moviesRepository) Latest(ctx context.Context, limit, offset uint64) (movies []model.MovieInfo, err error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, year, brief
		FROM movies
		ORDER BY added DESC
		LIMIT ?
		OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, rows.Close())
		if err != nil {
			movies = nil
		}
	}()

	movies = make([]model.MovieInfo, 0, limit)
	for rows.Next() {
		movie := model.MovieInfo{}
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Brief); err != nil {
			return nil, err
		}
		movie.Genres, err = r.genres(ctx, movie.ID)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *moviesRepository) Search(ctx context.Context, query string) (movies []model.MovieInfo, err error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, year, brief
		FROM movies
		WHERE lower(title) LIKE '%'||lower(?)||'%'
		LIMIT 5`, query,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, rows.Close())
		if err != nil {
			movies = nil
		}
	}()

	movies = make([]model.MovieInfo, 0, 5)
	for rows.Next() {
		movie := model.MovieInfo{}
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Brief); err != nil {
			return nil, err
		}
		movie.Genres, err = r.genres(ctx, movie.ID)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *moviesRepository) ByID(ctx context.Context, id model.ID) (m model.Movie, err error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT
			movies.id, movies.title,
			movies.year, movies.brief,
			movies.description,
			IFNULL(directors.id, 0) as director_id, IFNULL(directors.name, '') as director_name,
			IFNULL(writers.id, 0) as writer_id, IFNULL(writers.name, '') as writer_name
		FROM movies
		LEFT JOIN
		    (SELECT id, name from people) as directors
		    on directors.id == (
			select person_id
			from movie_crew
			where movie_id = movies.id and role == 'director'
		    )
		LEFT JOIN
		    (SELECT id, name from people) as writers
		    on writers.id == (
			select person_id
			from movie_crew
			where movie_id = movies.id and role == 'writer'
		    )
		where movies.id == ?`, id,
	)
	err = row.Err()
	if errors.Is(err, sql.ErrNoRows) {
		return model.Movie{}, model.ErrMovieNotExist
	} else if err != nil {
		return model.Movie{}, err
	}

	director := &model.CrewMember{Role: model.Director}
	writer := &model.CrewMember{Role: model.Writer}
	if err := row.Scan(
		&m.ID, &m.Title,
		&m.Year, &m.Brief, &m.Description,
		&director.ID, &director.Name, &writer.ID, &writer.Name,
	); err != nil {
		return model.Movie{}, err
	}

	genres, err := r.genres(ctx, id)
	if err != nil {
		return model.Movie{}, err
	}

	m.Genres = genres
	if writer.Name != "" {
		m.Writer = writer
	}
	if director.Name != "" {
		m.Director = director
	}

	return m, nil
}

func (r *moviesRepository) genres(ctx context.Context, movie model.ID) (genres []string, err error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT genre
		FROM genres
		WHERE movie_id = ?`, movie,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer func() {
		err = errors.Join(err, rows.Close())
		if err != nil {
			genres = nil
		}
	}()

	genres = []string{}
	for rows.Next() {
		var genre string
		if err := rows.Scan(&genre); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *usersRepository) ByEmail(ctx context.Context, email string) (user model.User, err error) {
	u := model.User{}
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, password FROM users
		WHERE email = ?`, email,
	)
	if err := row.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = model.ErrUserNotExist
		}

		return model.User{}, err
	}

	if err := row.Scan(&u.ID, &u.Name, &u.Password); err != nil {
		return model.User{}, err
	}

	u.Email = email
	return u, nil
}

func (r *usersRepository) ByID(ctx context.Context, id model.ID) (user model.User, err error) {
	u := model.User{}
	row := r.db.QueryRowContext(ctx,
		`SELECT name, email, password FROM users
		WHERE id = ? `, id,
	)
	if err := row.Err(); err != nil {
		return model.User{}, err
	}

	if err := row.Scan(&u.Name, &u.Email, &u.Password); err != nil {
		return model.User{}, err
	}

	u.ID = id
	return u, nil
}

func (r *usersRepository) Create(ctx context.Context, user model.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (name, email, password)
		VALUES (?, ?, ?)`, user.Name, user.Email, user.Password,
	)
	if err == nil {
		return nil
	}

	sErr, ok := err.(sqlite3.Error)
	if !ok {
		return err
	}

	if sErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return model.ErrUserExist
	}

	return err
}
