package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"urdb/components"
	"urdb/components/header"
	"urdb/components/movie"
	"urdb/model"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

var errUnknownInfo = errors.New("unknown info requested")

type moviePageQuery struct {
	ID   model.ID `param:"id"`
	Info string   `query:"info"`
}

func (s *Server) moviePage(c echo.Context) error {
	q := moviePageQuery{}
	if err := c.Bind(&q); err != nil {
		return s.badRequest(c)
	}

	m, err := s.movies.ByID(c.Request().Context(), q.ID)
	if errors.Is(err, model.ErrMovieNotExist) {
		return c.String(http.StatusNotFound, "Not Found")
	} else if err != nil {
		return s.internalError(c, err)
	}

	info, component, err := s.movieInfo(c.Request().Context(), q.Info, m.ID)
	if err != nil {
		return s.internalError(c, err)
	}

	c.Response().Header().Set("HX-Push-Url", fmt.Sprintf("/movie/%d?info=%s", q.ID, info))
	page := components.Index(
		header.Header(getUsernameFromCtx(c)),
		movie.Movie(m, component),
	)
	return render(c, page)
}

func (s *Server) movie(c echo.Context) error {
	q := moviePageQuery{}
	if err := c.Bind(&q); err != nil {
		return s.badRequest(c)
	}

	m, err := s.movies.ByID(c.Request().Context(), q.ID)
	if errors.Is(err, model.ErrMovieNotExist) {
		return c.String(http.StatusNotFound, "Not Found")
	} else if err != nil {
		return s.internalError(c, err)
	}

	info, infoComponent, err := s.movieInfo(c.Request().Context(), q.Info, m.ID)
	if errors.Is(err, errUnknownInfo) {
		return s.badRequest(c)
	} else if err != nil {
		return s.internalError(c, err)
	}

	c.Response().Header().Set("HX-Push-Url", fmt.Sprintf("/movie/%d?info=%s", q.ID, info))
	return render(c, movie.Movie(m, infoComponent))
}

func (s *Server) movieInfo(ctx context.Context, info string, movieID model.ID) (string, templ.Component, error) {
	if info == "" {
		info = "crew"
	}

	switch info {
	case "crew":
		crew, err := s.movies.Crew(ctx, movieID)
		if err != nil {
			return "", nil, err
		}
		return "crew", movie.Crew(crew), nil
	}

	return "", nil, errUnknownInfo
}

func (s *Server) crew(c echo.Context) error {
	q := moviePageQuery{}
	if err := c.Bind(&q); err != nil {
		return s.badRequest(c)
	}

	crew, err := s.movies.Crew(c.Request().Context(), q.ID)
	if err != nil {
		return s.internalError(c, err)
	}

	return movie.Crew(crew).Render(c.Request().Context(), c.Response())
}
