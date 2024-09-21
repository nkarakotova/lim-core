package servicesDataBuilder

import (
	"github.com/nkarakotova/lim-core/models"
)

type CoachBuilder struct {
	coach *models.Coach
}

func NewCoachBuilder() *CoachBuilder {
	return &CoachBuilder{
		coach: &models.Coach{
			ID: 1,
			Name: "Name",
		},
	}
}

func (b *CoachBuilder) WithID(id uint64) *CoachBuilder {
	b.coach.ID = id
	return b
}

func (b *CoachBuilder) WithName(name string) *CoachBuilder {
	b.coach.Name = name
	return b
}

func (b *CoachBuilder) Build() *models.Coach {
	return b.coach
}
