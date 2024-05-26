package models

import "time"

type Subscription struct {
	ID                    uint64
	TrainingsNum          uint64
	RemainingTrainingsNum uint64
	Cost                  uint64
	StartDate             time.Time
	EndDate               time.Time
}
