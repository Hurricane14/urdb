package server

import (
	"net/http"
	"time"
	"urdb/model"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func more(movies []model.MovieInfo, limit uint64) bool {
	return uint64(len(movies)) == limit
}

func (s *Server) internalError(c echo.Context, err error) error {
	s.router.Logger.Error(err)
	c.Response().WriteHeader(http.StatusInternalServerError)
	return nil
}

func (s *Server) badRequest(c echo.Context) error {
	c.Response().WriteHeader(http.StatusBadRequest)
	return nil
}

func setTokenCookie(c echo.Context, token string, ttl time.Duration) {
	c.SetCookie(&http.Cookie{
		Name:     "URDB-Authorization",
		Value:    token,
		Path:     "/",
		MaxAge:   int(ttl.Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
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

func getUsernameFromCtx(c echo.Context) string {
	u, ok := getUserFromCtx(c)
	if !ok {
		return ""
	}
	return u.Name
}
