package server

import (
	"errors"
	"time"
	"urdb/components"
	"urdb/model"

	"github.com/labstack/echo/v4"
)

func (s *Server) signInPage(c echo.Context) error {
	return components.Index(
		components.Header(
			getUsernameFromCtx(c),
		),
		components.SignIn(),
	).Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) signInForm(c echo.Context) error {
	c.Response().Header().Set("HX-Push-Url", "/signUp")
	return components.
		SignIn().
		Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) signUpPage(c echo.Context) error {
	return components.Index(
		components.Header(
			getUsernameFromCtx(c),
		),
		components.SignUp(),
	).Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) signUpForm(c echo.Context) error {
	c.Response().Header().Set("HX-Push-Url", "/signUp")
	return components.
		SignUp().
		Render(c.Request().Context(), c.Response().Writer)
}

type signInForm struct {
	Email    string `schema:"email" validate:"required,email"`
	Password string `schema:"password" validate:"required,password"`
}

func (s *Server) userSignIn(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return s.badRequest(c)
	}

	form := signInForm{}
	if err := s.schema.Decode(&form, c.Request().PostForm); err != nil {
		return s.badRequest(c)
	}

	s.router.Logger.Debug(form)

	validationErrs := s.validateForm(form)
	if len(validationErrs) != 0 {
		return components.
			ValidationList(validationErrs...).
			Render(c.Request().Context(), c.Response().Writer)
	}

	user, err := s.users.ByEmail(c.Request().Context(), form.Email)
	if err != nil && !errors.Is(err, model.ErrUserNotExist) {
		return s.internalError(c, err)
	} else if errors.Is(err, model.ErrUserNotExist) || user.Password != form.Password {
		return components.
			ValidationList(model.ErrUserNotExist).
			Render(c.Request().Context(), c.Response().Writer)
	}

	token := s.auth.CreateToken(user.ID, time.Now().Add(s.cookieTTL))
	setTokenCookie(c, token, s.cookieTTL)

	s.router.Logger.Debugf("Set cookie for user %s", user.Name)
	c.Response().Header().Set("HX-Location", "/")
	return nil
}

type signUpForm struct {
	Name          string `schema:"name" validate:"min=4,max=16"`
	Email         string `schema:"email" validate:"required,email"`
	Password      string `schema:"password" validate:"required,password"`
	PasswordAgain string `schema:"passwordAgain" validate:"required_with=Password|eqfield=Password"`
}

func (s *Server) userSignUp(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return s.badRequest(c)
	}

	form := signUpForm{}
	if err := s.schema.Decode(&form, c.Request().PostForm); err != nil {
		return s.badRequest(c)
	}

	s.router.Logger.Debug(form)

	validationErrs := s.validateForm(form)
	if len(validationErrs) != 0 {
		return components.
			ValidationList(validationErrs...).
			Render(c.Request().Context(), c.Response().Writer)
	}

	err := s.users.Create(c.Request().Context(), model.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	})
	if errors.Is(err, model.ErrUserExist) {
		return components.
			ValidationList(model.ErrUserExist).
			Render(c.Request().Context(), c.Response().Writer)
	} else if !errors.Is(err, model.ErrUserNotExist) {
		s.internalError(c, err)
	}

	c.Response().Header().Set("HX-Location", "/signIn")
	return nil
}

func (s *Server) userSignOut(c echo.Context) error {
	setTokenCookie(c, "", 0)
	return components.Header("").Render(c.Request().Context(), c.Response().Writer)
}
