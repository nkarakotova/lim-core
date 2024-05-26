package repositories

import (
	"context"

	"github.com/nkarakotova/lim-core/models"
)

//go:generate mockgen  -source=subscription.go -destination=mocks/subscription.go
type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *models.Subscription, clientID uint64) error
	GetByID(ctx context.Context, id uint64) (*models.Subscription, error)
	ReduceRemainingTrainingsNum(ctx context.Context, id uint64) error
	IncreaseRemainingTrainingsNum(ctx context.Context, id uint64) error
}
