package server

import (
	"net/http"
	"time"
	"urdb/components"

	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
)

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

var decoder = schema.NewDecoder()

type SignInForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func (s *Server) userSignIn(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return s.badRequest(c)
	}

	form := SignInForm{}
	if err := decoder.Decode(&form, c.Request().PostForm); err != nil {
		return s.badRequest(c)
	}

	s.router.Logger.Debugf("Received Sign In request: email: %q, password: %q", form.Email, form.Password)
	c.Redirect(http.StatusMovedPermanently, "/")

	return nil
}
