package repositories

import (
	"context"

	"github.com/nkarakotova/lim-core/models"
)

//go:generate mockgen  -source=coach.go -destination=mocks/coach.go
type CoachRepository interface {
	Create(ctx context.Context, coach *models.Coach) error
	GetByID(ctx context.Context, id uint64) (*models.Coach, error)
	GetByName(ctx context.Context, name string) (*models.Coach, error)
	AddDirection(ctx context.Context, coachID, directionID uint64) error
	GetAllByDirection(ctx context.Context, id uint64) ([]models.Coach, error)
}
