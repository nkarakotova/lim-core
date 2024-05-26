package repositories

import (
	"context"

	"github.com/nkarakotova/lim-core/models"
)

//go:generate mockgen  -source=hall.go -destination=mocks/hall.go
type HallRepository interface {
	Create(ctx context.Context, hall *models.Hall) error
	GetByID(ctx context.Context, id uint64) (*models.Hall, error)
	GetByNumber(ctx context.Context, number uint64) (*models.Hall, error)
	GetAll(ctx context.Context) (map[uint64]models.Hall, error)
}
