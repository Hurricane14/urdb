package server

import (
	"net/http"
	"strconv"
	"time"
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

func setTokenCookie(c echo.Context, token string, expireAt time.Time) {
	c.SetCookie(&http.Cookie{
		Name:     "URDB-Authorization",
		Value:    token,
		Expires:  expireAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func getTokenFromCookie(c echo.Context) (token string, ok bool) {
	cookie, err := c.Cookie("URDB-Authorization")
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}

func setUserInCtx(c echo.Context, u model.User) {
	c.Set("User", u)
}

func getUserFromCtx(c echo.Context) (model.User, bool) {
	val, ok := c.Get("User").(model.User)
	return val, ok
}
