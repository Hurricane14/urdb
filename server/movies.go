package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"urdb/components"
	"urdb/model"

	"github.com/labstack/echo/v4"
)

func (s *Server) latestMovies(c echo.Context) error {
	time.Sleep(1 * time.Second)
	limit := intQueryParamWithDefault(c, "limit", defaultLimit)
	offset := intQueryParamWithDefault(c, "offset", 0)
	moviesInfo, err := s.movies.Latest(c.Request().Context(), limit, offset)
	if err != nil {
		return s.internalError(c)
	}
	if err := components.Movies(moviesInfo).Render(c.Request().Context(), c.Response().Writer); err != nil {
		return err
	}
	return components.More(
		more(moviesInfo, limit), limit, limit+offset,
	).Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) searchMovies(c echo.Context) error {
	time.Sleep(1 * time.Second)
	query := c.QueryParam("q")
	if query != "" {
		moviesInfo, err := s.movies.Search(c.Request().Context(), query)
		if err != nil {
			return s.internalError(c)
		}
		return components.MoviesDiv(
			components.Movies(moviesInfo),
			components.MoviesLoadingIndicator(),
		).Render(c.Request().Context(), c.Response().Writer)
	}

	limit, offset := defaultLimit, 0
	moviesInfo, err := s.movies.Latest(c.Request().Context(), limit, offset)
	if err != nil {
		return s.internalError(c)
	}
	movies := components.MoviesDiv(
		components.Movies(moviesInfo),
		components.More(more(moviesInfo, limit), limit, limit+offset),
		components.MoviesLoadingIndicator(),
	)
	return movies.Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) moviePage(c echo.Context) error {
	params := struct {
		MovieID model.ID `param:"id"`
	}{}
	if err := c.Bind(&params); err != nil {
		return s.badRequest(c)
	}

	movie, err := s.movies.ByID(c.Request().Context(), params.MovieID)
	if errors.Is(err, model.ErrMovieNotExist) {
		return c.String(http.StatusNotFound, "Not Found")
	} else if err != nil {
		s.router.Logger.Error(err)
		return s.internalError(c)
	}

	return components.Index(
		components.Header(getUsernameFromCtx(c)),
		components.Movie(movie),
	).Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) movieInfo(c echo.Context) error {
	time.Sleep(time.Second)
	params := struct {
		MovieID model.ID `param:"id"`
	}{}
	if err := c.Bind(&params); err != nil {
		return s.badRequest(c)
	}

	movie, err := s.movies.ByID(c.Request().Context(), params.MovieID)
	if errors.Is(err, model.ErrMovieNotExist) {
		return c.String(http.StatusNotFound, "Not Found")
	} else if err != nil {
		s.router.Logger.Error(err)
		return s.internalError(c)
	}

	c.Response().Header().Set("HX-Push-Url", fmt.Sprintf("/movie/%d", params.MovieID))
	return components.Movie(movie).
		Render(c.Request().Context(), c.Response().Writer)
}
