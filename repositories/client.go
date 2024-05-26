package repositories

import (
	"context"

	"github.com/nkarakotova/lim-core/models"
)

//go:generate mockgen  -source=client.go -destination=mocks/client.go
type ClientRepository interface {
	Create(ctx context.Context, client *models.Client) error
	GetByID(ctx context.Context, id uint64) (*models.Client, error)
	GetByTelephone(ctx context.Context, telephone string) (*models.Client, error)
	GetByTraining(ctx context.Context, id uint64) ([]models.Client, error)
	СreateAssignment(ctx context.Context, clientID, trainingID uint64) error
	DeleteAssignment(ctx context.Context, clientID, trainingID uint64) error
}
