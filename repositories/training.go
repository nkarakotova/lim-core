package repositories

import (
	"context"
	"time"

	"github.com/nkarakotova/lim-core/models"
)

//go:generate mockgen  -source=training.go -destination=mocks/training.go
type TrainingRepository interface {
	Create(ctx context.Context, training *models.Training) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*models.Training, error)
	GetAllByClient(ctx context.Context, id uint64) ([]models.Training, error)
	GetAllByCoachOnDate(ctx context.Context, id uint64, date time.Time) ([]models.Training, error)
	GetAllByDateTime(ctx context.Context, dateTime time.Time) ([]models.Training, error)
	GetAllBetweenDateTime(ctx context.Context, start time.Time, end time.Time) ([]models.Training, error)
	ReduceAvailablePlacesNum(ctx context.Context, id uint64) error
	IncreaseAvailablePlacesNum(ctx context.Context, id uint64) error
}
