package servicesDataBuilder

import (
	"github.com/nkarakotova/lim-core/models"
)

type HallBuilder struct {
	hall *models.Hall
}

func NewHallBuilder() *HallBuilder {
	return &HallBuilder{
		hall: &models.Hall{
			ID: 1,
			Number: 1,
		},
	}
}

func (b *HallBuilder) WithID(id uint64) *HallBuilder {
	b.hall.ID = id
	return b
}

func (b *HallBuilder) WithNumber(number uint64) *HallBuilder {
	b.hall.Number = number
	return b
}

func (b *HallBuilder) Build() *models.Hall {
	return b.hall
}
