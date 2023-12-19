package server

import (
	"errors"
	"fmt"
	"net/http"
	"urdb/components"
	"urdb/model"

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
	if err := render(c, components.Movies(moviesInfo)); err != nil {
		return err
	}

	return render(c, components.More(
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
		movies := components.MoviesDiv(
			components.Movies(moviesInfo),
			components.MoviesLoadingIndicator(),
		)
		return render(c, movies)
	}

	var limit, offset uint64 = defaultLimit, 0
	moviesInfo, err := s.movies.Latest(c.Request().Context(), limit, offset)
	if err != nil {
		return s.internalError(c, err)
	}
	movies := components.MoviesDiv(
		components.Movies(moviesInfo),
		components.More(more(moviesInfo, limit), limit, limit+offset),
		components.MoviesLoadingIndicator(),
	)
	return render(c, movies)
}

type moviePageQuery struct {
	ID model.ID `param:"id"`
}

func (s *Server) moviePage(c echo.Context) error {
	q := moviePageQuery{}
	if err := c.Bind(&q); err != nil {
		return s.badRequest(c)
	}

	movie, err := s.movies.ByID(c.Request().Context(), q.ID)
	if errors.Is(err, model.ErrMovieNotExist) {
		return c.String(http.StatusNotFound, "Not Found")
	} else if err != nil {
		return s.internalError(c, err)
	}

	page := components.Index(
		components.Header(getUsernameFromCtx(c)),
		components.Movie(movie),
	)
	return render(c, page)
}

func (s *Server) movieInfo(c echo.Context) error {
	q := moviePageQuery{}
	if err := c.Bind(&q); err != nil {
		return s.badRequest(c)
	}

	movie, err := s.movies.ByID(c.Request().Context(), q.ID)
	if errors.Is(err, model.ErrMovieNotExist) {
		return c.String(http.StatusNotFound, "Not Found")
	} else if err != nil {
		return s.internalError(c, err)
	}

	c.Response().Header().Set("HX-Push-Url", fmt.Sprintf("/movie/%d", q.ID))
	return render(c, components.Movie(movie))
}
