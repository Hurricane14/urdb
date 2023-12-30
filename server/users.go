package server

import (
	"errors"
	"time"
	"urdb/components"
	"urdb/components/auth"
	"urdb/components/header"
	"urdb/model"

	"github.com/labstack/echo/v4"
)

func (s *Server) signInPage(c echo.Context) error {
	page := components.Index(
		header.Header(
			getUsernameFromCtx(c),
		),
		auth.SignIn(),
	)
	return render(c, page)
}

func (s *Server) signInForm(c echo.Context) error {
	c.Response().Header().Set("HX-Push-Url", "/signIn")
	return render(c, auth.SignIn())
}

func (s *Server) signUpPage(c echo.Context) error {
	page := components.Index(
		header.Header(
			getUsernameFromCtx(c),
		),
		auth.SignUp(),
	)
	return render(c, page)
}

func (s *Server) signUpForm(c echo.Context) error {
	c.Response().Header().Set("HX-Push-Url", "/signUp")
	return render(c, auth.SignUp())
}

type signInForm struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,password"`
}

func (s *Server) userSignIn(c echo.Context) error {
	form := signInForm{}
	if err := c.Bind(&form); err != nil {
		return s.badRequest(c)
	}

	s.router.Logger.Debug(form)

	validationErrs := s.validateForm(form)
	if len(validationErrs) != 0 {
		return render(c, auth.ValidationList(validationErrs...))
	}

	user, err := s.users.ByEmail(c.Request().Context(), form.Email)
	if err != nil && !errors.Is(err, model.ErrUserNotExist) {
		return s.internalError(c, err)
	} else if errors.Is(err, model.ErrUserNotExist) || user.Password != form.Password {
		return render(c, auth.ValidationList(model.ErrUserNotExist))
	}

	token := s.auth.CreateToken(user.ID, time.Now().Add(s.cookieTTL))
	setTokenCookie(c, token, s.cookieTTL)

	s.router.Logger.Debugf("Set cookie for user %s", user.Name)
	c.Response().Header().Set("HX-Location", "/")
	return nil
}

type signUpForm struct {
	Name          string `form:"name" validate:"min=4,max=16"`
	Email         string `form:"email" validate:"required,email"`
	Password      string `form:"password" validate:"required,password"`
	PasswordAgain string `form:"passwordAgain" validate:"required_with=Password|eqfield=Password"`
}

func (s *Server) userSignUp(c echo.Context) error {
	form := signUpForm{}
	if err := c.Bind(&form); err != nil {
		return s.badRequest(c)
	}

	s.router.Logger.Debug(form)
	validationErrs := s.validateForm(form)
	if len(validationErrs) != 0 {
		return render(c, auth.ValidationList(validationErrs...))
	}

	err := s.users.Create(c.Request().Context(), model.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	})
	if errors.Is(err, model.ErrUserExist) {
		return render(c, auth.ValidationList(model.ErrUserExist))
	} else if err != nil {
		s.internalError(c, err)
	}

	c.Response().Header().Set("HX-Location", "/signIn")
	return nil
}

func (s *Server) userSignOut(c echo.Context) error {
	setTokenCookie(c, "", 0)
	return render(c, header.Header(""))
}
