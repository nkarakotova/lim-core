package models

type Direction struct {
	ID               uint64
	Name             string
	Description      string
	AcceptableGender Gender
}
