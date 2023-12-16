package model

type Review struct {
	UserID  uint64
	MovieID uint64
	Text    string
	Rating  float64
}
