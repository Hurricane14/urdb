package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"urdb/components"
	"urdb/model"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const defaultLimit = 1

type UsersRepository interface {
	ByEmail(ctx context.Context, email string) (model.User, error)
	ByID(ctx context.Context, id model.ID) (model.User, error)
	Create(ctx context.Context, user model.User) error
}

type MoviesRepository interface {
	Latest(ctx context.Context, limit, offset int) ([]model.MovieInfo, error)
	Search(ctx context.Context, query string) ([]model.MovieInfo, error)
	ByID(ctx context.Context, id model.ID) (model.Movie, error)
}

type AuthService interface {
	CreateToken(id model.ID, expireAt time.Time) string
	ParseToken(token string) (model.ID, error)
	HashPassword(pass string) []byte
	MatchPassword(got string, want []byte) bool
}

type Server struct {
	router    *echo.Echo
	users     UsersRepository
	movies    MoviesRepository
	auth      AuthService
	schema    *schema.Decoder
	validator *validator.Validate
	cookieTTL time.Duration
}

func New(
	users UsersRepository,
	movies MoviesRepository,
	auth AuthService,
) *Server {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("password", validatePassword)

	s := &Server{
		router:    e,
		users:     users,
		movies:    movies,
		auth:      auth,
		schema:    decoder,
		validator: validate,
		cookieTTL: 24 * time.Hour,
	}

	e.Static("/static", "static")

	e.Use(s.processAuthToken)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"localhost"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
		AllowHeaders: []string{
			"HX-Boosted",
			"HX-Current-URL",
			"HX-History-Restore-Request",
			"HX-Prompt",
			"HX-Request",
			"HX-Target",
			"HX-TriggerName",
			"HX-Trigger",
		},
	}))

	usersAPI := e.Group("/users")
	usersAPI.POST("/signIn", s.userSignIn)
	usersAPI.POST("/signUp", s.userSignUp)
	usersAPI.POST("/signOut", s.userSignOut, s.requireAuthorization)

	e.GET("/movie/:id/info", s.movieInfo)
	e.GET("/movie/:id", s.moviePage)
	e.GET("/searchMovies", s.searchMovies)
	e.GET("/latestMovies", s.latestMovies)
	e.GET("/signInForm", s.signInForm)
	e.GET("/signIn", s.signInPage)
	e.GET("/signUpForm", s.signUpForm)
	e.GET("/signUp", s.signUpPage)
	e.GET("/", s.indexPage)

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

func (s *Server) indexPage(c echo.Context) error {
	const limit, offset = defaultLimit, 0
	moviesInfo, err := s.movies.Latest(c.Request().Context(), limit, offset)
	if err != nil {
		return s.internalError(c)
	}
	header := components.Header(getUsernameFromCtx(c))
	searchBar := components.SearchBar()
	movies := components.MoviesDiv(
		components.Movies(moviesInfo),
		components.More(len(moviesInfo) == limit, limit, limit+offset),
		components.MoviesLoadingIndicator(),
	)
	page := components.Index(header, searchBar, movies)

	return page.Render(c.Request().Context(), c.Response().Writer)
}
