package model

import "errors"

type Role uint16

const (
	Unknown Role = iota
	Director
	Writer
	Actor
)

func (r Role) String() string {
	var s string
	switch r {
	case Unknown:
		s = "Unknown"
	case Director:
		s = "Director"
	case Writer:
		s = "Writer"
	case Actor:
		s = "Actor"
	default:
		panic("Unexpected role value")
	}
	return s
}

type CrewMember struct {
	ID
	Name string
	Role Role
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
