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
	if err := components.Movies(moviesInfo).Render(c.Request().Context(), c.Response().Writer); err != nil {
		return err
	}
	return components.More(
		more(moviesInfo, q.Limit), q.Limit, q.Limit+q.Offset,
	).Render(c.Request().Context(), c.Response().Writer)
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
		return components.MoviesDiv(
			components.Movies(moviesInfo),
			components.MoviesLoadingIndicator(),
		).Render(c.Request().Context(), c.Response().Writer)
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
	return movies.Render(c.Request().Context(), c.Response().Writer)
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

	return components.Index(
		components.Header(getUsernameFromCtx(c)),
		components.Movie(movie),
	).Render(c.Request().Context(), c.Response().Writer)
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
	return components.Movie(movie).
		Render(c.Request().Context(), c.Response().Writer)
}
