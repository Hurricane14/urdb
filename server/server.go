package server

import (
	"context"
	"fmt"
	"time"
	"urdb/components"
	"urdb/model"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const defaultLimit = 1

type UsersRepository interface {
	SignIn(ctx context.Context, email, password string) (model.User, error)
}

type MoviesRepository interface {
	Latest(ctx context.Context, limit, offset int) ([]model.MovieInfo, error)
	Search(ctx context.Context, query string) ([]model.MovieInfo, error)
}

type Server struct {
	router *echo.Echo
	users  UsersRepository
	movies MoviesRepository
}

func New(users UsersRepository, movies MoviesRepository) *Server {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	s := &Server{
		router: e,
		users:  users,
		movies: movies,
	}

	e.Static("/static", "static")
	usersAPI := e.Group("/users")
	usersAPI.POST("/signIn", s.userSignIn)
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
