package server

import (
	"net/http"
	"strconv"
	"urdb/model"

	"github.com/labstack/echo/v4"
)

func more(movies []model.MovieInfo, limit int) bool {
	return len(movies) == limit
}

func (s *Server) internalError(c echo.Context) error {
	c.Response().WriteHeader(http.StatusInternalServerError)
	return nil
}

func (s *Server) badRequest(c echo.Context) error {
	c.Response().WriteHeader(http.StatusBadRequest)
	return nil
}

func intQueryParamWithDefault(c echo.Context, name string, dflt int) int {
	p := c.QueryParam(name)
	if v, err := strconv.Atoi(p); err != nil {
		return dflt
	} else {
		return v
	}
}
