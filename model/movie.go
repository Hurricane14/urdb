package model

import "errors"

type Movie struct {
	MovieInfo
	Director    *CrewMember
	Writer      *CrewMember
	Description string
}

type MovieInfo struct {
	ID
	Title  string
	Genres []string
	Year   uint16
	Brief  string
	Rating float64
}

var ErrMovieNotExist = errors.New("Movie not found")
