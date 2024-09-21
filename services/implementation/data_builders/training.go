package servicesDataBuilder

import (
	"time"
	"github.com/nkarakotova/lim-core/models"
)

type TrainingBuilder struct {
	training *models.Training
}

func NewTrainingBuilder() *TrainingBuilder {
	return &TrainingBuilder{
		training: &models.Training{
			ID: 1,
			CoachID: 1,
			HallID: 1,
			Name: "Name",
			DateTime: time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC),
			PlacesNum: 10,
		},
	}
}

func (b *TrainingBuilder) WithID(id uint64) *TrainingBuilder {
	b.training.ID = id
	return b
}

func (b *TrainingBuilder) WithCoachID(coachID uint64) *TrainingBuilder {
	b.training.CoachID = coachID
	return b
}

func (b *TrainingBuilder) WithHallID(hallID uint64) *TrainingBuilder {
	b.training.HallID = hallID
	return b
}

func (b *TrainingBuilder) WithName(name string) *TrainingBuilder {
	b.training.Name = name
	return b
}

func (b *TrainingBuilder) WithDateTime(dateTime time.Time) *TrainingBuilder {
	b.training.DateTime = dateTime
	return b
}

func (b *TrainingBuilder) WithPlacesNum(placesNum uint64) *TrainingBuilder {
	b.training.PlacesNum = placesNum
	return b
}

func (b *TrainingBuilder) Build() *models.Training {
	return b.training
}
