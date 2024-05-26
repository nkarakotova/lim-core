package repositories

import (
	"context"

	"github.com/nkarakotova/lim-core/models"
)

//go:generate mockgen  -source=direction.go -destination=mocks/direction.go
type DirectionRepository interface {
	Create(ctx context.Context, direction *models.Direction) error
	GetByID(ctx context.Context, id uint64) (*models.Direction, error)
	GetByName(ctx context.Context, name string) (*models.Direction, error)
	GetAll(ctx context.Context) ([]models.Direction, error)
}
