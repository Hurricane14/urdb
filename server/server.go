package server

import (
	"context"
	"fmt"
	"time"
	"urdb/components"
	"urdb/model"

	"github.com/labstack/echo/v4"
)

const defaultLimit = 1

type MoviesRepository interface {
	Latest(ctx context.Context, limit, offset int) ([]model.MovieInfo, error)
	Search(ctx context.Context, query string) ([]model.MovieInfo, error)
}

type Server struct {
	router *echo.Echo
	movies MoviesRepository
}

func New(repo MoviesRepository) *Server {
	e := echo.New()
	s := &Server{
		router: e,
		movies: repo,
	}

	e.Static("/static", "static")
	// users := e.Group("/users")
	// users.POST("/signIn", s.userSignIn)
	e.GET("/searchMovies", s.searchMovies)
	e.GET("/latestMovies", s.latestMovies)
	e.GET("/signInForm", s.signInForm)
	e.GET("/signIn", s.signIn)
	e.GET("/signUpForm", s.signUpForm)
	e.GET("/signUp", s.signUp)
	e.GET("/", s.index)

	return s
}

func (s *Server) Run(port int) error {
	return s.router.Start(fmt.Sprintf("localhost:%d", port))
}

func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.router.Shutdown(ctx)
}

func (s *Server) index(c echo.Context) error {
	const limit, offset = defaultLimit, 0
	moviesInfo, err := s.movies.Latest(c.Request().Context(), limit, offset)
	if err != nil {
		return s.internalError(c)
	}
	header := components.Header()
	searchBar := components.SearchBar()
	movies := components.MoviesDiv(
		components.Movies(moviesInfo),
		components.More(len(moviesInfo) == limit, limit, limit+offset),
		components.MoviesLoadingIndicator(),
	)
	page := components.Index(header, searchBar, movies)

	return page.Render(c.Request().Context(), c.Response().Writer)
}

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

func (s *Server) signIn(c echo.Context) error {
	return components.Index(
		components.Header(),
		components.SignIn(),
	).Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) signInForm(c echo.Context) error {
	time.Sleep(time.Second)
	return components.
		SignIn().
		Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) signUp(c echo.Context) error {
	return components.Index(
		components.Header(),
		components.SignUp(),
	).Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) signUpForm(c echo.Context) error {
	time.Sleep(time.Second)
	return components.
		SignUp().
		Render(c.Request().Context(), c.Response().Writer)
}
