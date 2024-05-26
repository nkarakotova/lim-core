package models

import "time"

type Training struct {
	ID                 uint64
	CoachID            uint64
	HallID             uint64
	DirectionID        uint64
	Name               string
    DateTime        time.Time
	PlacesNum          uint64
	AvailablePlacesNum uint64
	AcceptableAge      uint16
}
