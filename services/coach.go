package services

import (
	"time"

	"github.com/nkarakotova/lim-core/models"
)

type CoachService interface {
	Create(coach *models.Coach) error
	GetByID(id uint64) (*models.Coach, error)
	GetByName(name string) (*models.Coach, error)
	AddDirection(coachID, directionID uint64) error
	GetAllByDirection(id uint64) ([]models.Coach, error)
	GetFreeTimeOnDate(id uint64, date time.Time) ([]time.Time, error)
}
