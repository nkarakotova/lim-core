package services

import (
	"time"

	"github.com/nkarakotova/lim-core/models"
)

type HallService interface {
	Create(hall *models.Hall) error
	GetByID(id uint64) (*models.Hall, error)
	GetByNumber(number uint64) (*models.Hall, error)
	GetFreeOnDateTime(dateTime time.Time) (map[uint64]models.Hall, error)
}
