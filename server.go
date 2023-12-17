package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"urdb/components"
	"urdb/model"

	"github.com/labstack/echo/v4"
)

const defaultLimit = 1

type server struct {
	router *echo.Echo
	repo   *repository
}

func newServer(repo *repository) *server {
	e := echo.New()
	s := &server{
		router: e,
		repo:   repo,
	}

	e.Static("/static", "static")
	e.GET("/searchMovies", s.searchMovies)
	e.GET("/latestMovies", s.latestMovies)
	e.GET("/signInForm", s.signInForm)
	e.GET("/signIn", s.signIn)
	e.GET("/signUpForm", s.signUpForm)
	e.GET("/signUp", s.signUp)
	e.GET("/", s.index)

	return s
}

func (s *server) run(port int) error {
	return s.router.Start(fmt.Sprintf("localhost:%d", port))
}

func (s *server) shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.router.Shutdown(ctx)
}

func (s *server) index(c echo.Context) error {
	const limit, offset = defaultLimit, 0
	moviesInfo, err := s.repo.latestMovies(c.Request().Context(), limit, offset)
	if err != nil {
		return s.internalError(c)
	}
	header := components.Header()
	searchBar := components.SearchBar()
	movies := components.MoviesDiv(
		components.Movies(moviesInfo),
		components.More(len(moviesInfo) == limit, limit, limit+offset),
		components.LoadingIndicator(),
	)
	page := components.Index(header, searchBar, movies)

	return page.Render(c.Request().Context(), c.Response().Writer)
}

func (s *server) latestMovies(c echo.Context) error {
	time.Sleep(1 * time.Second)
	limit := intQueryParamWithDefault(c, "limit", defaultLimit)
	offset := intQueryParamWithDefault(c, "offset", 0)
	moviesInfo, err := s.repo.latestMovies(c.Request().Context(), limit, offset)
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

func (s *server) searchMovies(c echo.Context) error {
	time.Sleep(1 * time.Second)
	query := c.QueryParam("q")
	if query != "" {
		moviesInfo, err := s.repo.searchMovies(c.Request().Context(), query)
		if err != nil {
			return s.internalError(c)
		}
		return components.MoviesDiv(
			components.Movies(moviesInfo),
			components.LoadingIndicator(),
		).Render(c.Request().Context(), c.Response().Writer)
	}

	limit, offset := defaultLimit, 0
	moviesInfo, err := s.repo.latestMovies(c.Request().Context(), limit, offset)
	if err != nil {
		return s.internalError(c)
	}
	movies := components.MoviesDiv(
		components.Movies(moviesInfo),
		components.More(more(moviesInfo, limit), limit, limit+offset),
		components.LoadingIndicator(),
	)
	return movies.Render(c.Request().Context(), c.Response().Writer)
}

func (s *server) signIn(c echo.Context) error {
	return components.Index(
		components.Header(),
		components.SignIn(),
	).Render(c.Request().Context(), c.Response().Writer)
}

func (s *server) signInForm(c echo.Context) error {
	return components.
		SignIn().
		Render(c.Request().Context(), c.Response().Writer)
}

func (s *server) signUp(c echo.Context) error {
	return components.Index(
		components.Header(),
		components.SignUp(),
	).Render(c.Request().Context(), c.Response().Writer)
}

func (s *server) signUpForm(c echo.Context) error {
	return components.
		SignUp().
		Render(c.Request().Context(), c.Response().Writer)
}

func more(movies []model.MovieInfo, limit int) bool {
	return len(movies) == limit
}

func (s *server) internalError(c echo.Context) error {
	c.Response().WriteHeader(http.StatusInternalServerError)
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
