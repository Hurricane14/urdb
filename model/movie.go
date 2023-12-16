package model

type Movie struct {
	MovieInfo
	Director    Person
	Description string
}

type MovieInfo struct {
	ID     uint64
	Title  string
	Genres []string
	Year   uint16
	Brief  string
	Rating float64
}
