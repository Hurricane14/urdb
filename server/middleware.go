package server

import (
	"time"

	"github.com/labstack/echo/v4"
)

func (s *Server) processAuthToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		defer func() {
			err = next(c)
		}()

		token, ok := getTokenFromCookie(c)
		if !ok {
			return nil
		}

		userID, err := s.auth.ParseToken(token)
		if err != nil {
			return nil
		}

		user, err := s.users.ByID(c.Request().Context(), userID)
		if err != nil {
			return nil
		}

		setUserInCtx(c, user)
		return nil
	}
}

func (s *Server) requireAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, ok := getUserFromCtx(c)
		if !ok {
			c.Response().Header().Set("HX-Location", "/signIn")
			return nil
		}

		return next(c)
	}
}

func (s *Server) serverDelay(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if s.delay != 0 {
			time.Sleep(s.delay)
		}
		return next(c)
	}
}
