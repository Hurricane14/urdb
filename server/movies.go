package server

import (
	"urdb/components/movies"

	"github.com/labstack/echo/v4"
)

const defaultLimit = 1

type latestMoviesQuery struct {
	Limit  uint64 `query:"limit"`
	Offset uint64 `query:"offset"`
}

func (s *Server) latestMovies(c echo.Context) error {
	q := latestMoviesQuery{}
	if err := c.Bind(&q); err != nil {
		return s.badRequest(c)
	}

	moviesInfo, err := s.movies.Latest(c.Request().Context(), q.Limit, q.Offset)
	if err != nil {
		return s.internalError(c, err)
	}
	if err := render(c, movies.Movies(moviesInfo)); err != nil {
		return err
	}

	return render(c, movies.More(
		more(moviesInfo, q.Limit), q.Limit, q.Limit+q.Offset),
	)
}

type searchQuery struct {
	Query string `query:"q"`
}

func (s *Server) searchMovies(c echo.Context) error {
	q := searchQuery{}
	if err := c.Bind(&q); err != nil {
		s.badRequest(c)
	}

	if q.Query != "" {
		moviesInfo, err := s.movies.Search(c.Request().Context(), q.Query)
		if err != nil {
			return s.internalError(c, err)
		}
		movies := movies.MoviesDiv(
			movies.Movies(moviesInfo),
			movies.MoviesLoadingIndicator(),
		)
		return render(c, movies)
	}

	var limit, offset uint64 = defaultLimit, 0
	moviesInfo, err := s.movies.Latest(c.Request().Context(), limit, offset)
	if err != nil {
		return s.internalError(c, err)
	}
	movies := movies.MoviesDiv(
		movies.Movies(moviesInfo),
		movies.More(more(moviesInfo, limit), limit, limit+offset),
		movies.MoviesLoadingIndicator(),
	)
	return render(c, movies)
}
