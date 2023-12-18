package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type auther struct {
	key []byte
}

func newAuther(key string) *auther {
	return &auther{
		key: []byte(key),
	}
}

func (a *auther) CreateToken(id uint64, expireAt time.Time) string {
	token, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   fmt.Sprint(id),
			ExpiresAt: jwt.NewNumericDate(expireAt),
		}).SignedString(a.key)
	if err != nil {
		panic(err)
	}

	return token

}

func (a *auther) ParseToken(token string) (uint64, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return a.key, nil
	})
	if err != nil {
		return 0, err
	}

	subj, err := t.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	if v, err := strconv.ParseUint(subj, 10, 64); err != nil {
		return 0, err
	} else {
		return v, nil
	}
}

func (a *auther) HashPassword(pass string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return hash
}

func (a *auther) MatchPassword(got string, want []byte) bool {
	err := bcrypt.CompareHashAndPassword(want, []byte(got))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false
	} else if err != nil {
		panic(err)
	} else {
		return true
	}
}
