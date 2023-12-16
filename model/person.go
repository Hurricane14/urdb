package model

type Role string

const (
	Director Role = "director"
	Acto     Role = "actor"
)

type Person struct {
	ID        uint64
	Name      string
	BirthYear uint16
	DeathYear uint16
	Role      Role
}
