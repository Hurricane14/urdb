package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenHandler struct {
	key []byte
}

func newTokenHandler(key string) *tokenHandler {
	return &tokenHandler{
		key: []byte(key),
	}
}

func (h *tokenHandler) Create(id uint64, expireAt time.Time) string {
	token, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   fmt.Sprint(id),
			ExpiresAt: jwt.NewNumericDate(expireAt),
		}).SignedString(h.key)
	if err != nil {
		panic(err)
	}

	return token

}

func (h *tokenHandler) Parse(token string) (uint64, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return h.key, nil
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
