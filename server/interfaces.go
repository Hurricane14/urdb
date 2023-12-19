package server

import (
	"context"
	"time"
	"urdb/model"
)

type Logger interface {
	Debug(any)
	Debugf(string, ...any)
	Error(any)
	Errorf(string, ...any)
}

type UsersRepository interface {
	ByEmail(ctx context.Context, email string) (model.User, error)
	ByID(ctx context.Context, id model.ID) (model.User, error)
	Create(ctx context.Context, user model.User) error
}

type MoviesRepository interface {
	Latest(ctx context.Context, limit, offset uint64) ([]model.MovieInfo, error)
	Search(ctx context.Context, query string) ([]model.MovieInfo, error)
	ByID(ctx context.Context, id model.ID) (model.Movie, error)
}

type AuthService interface {
	CreateToken(id model.ID, expireAt time.Time) string
	ParseToken(token string) (model.ID, error)
	HashPassword(pass string) []byte
	MatchPassword(got string, want []byte) bool
}
