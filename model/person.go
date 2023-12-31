package model

import "errors"

type CrewMember struct {
	ID
	Name  string
	Roles []string
}

type Person struct {
	ID
	Name       string
	BirthYear  uint16
	DeathYear  uint16
	Occupation string
	Bio        string
	Roles      []string
	Movies     []MovieInfo
}

var (
	ErrPersonNotExist = errors.New("Person not found")
)
