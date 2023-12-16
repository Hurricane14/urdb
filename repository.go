package main

import (
	"context"
	"database/sql"
	"errors"
	"urdb/model"

	_ "github.com/mattn/go-sqlite3"
)

type repository struct {
	db *sql.DB
}

func newRepository() (*repository, error) {
	db, err := sql.Open("sqlite3", "sqlite.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &repository{db}, nil
}

func (r *repository) latestMovies(ctx context.Context, limit, offset int) (movies []model.MovieInfo, err error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, year, brief
		FROM movies
		ORDER BY added DESC
		LIMIT ?
		OFFSET ?
	`, limit, offset)
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

func (r *repository) searchMovies(ctx context.Context, query string) (movies []model.MovieInfo, err error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, year, brief
		FROM movies
		WHERE lower(title) LIKE '%'||lower(?)||'%'
		LIMIT 5
	`, query)
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

func (r *repository) genres(ctx context.Context, movieID uint64) (genres []string, err error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT genre
		FROM genres
		WHERE movie_id = ?
	`, movieID)
	if err != nil {
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
