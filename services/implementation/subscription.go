package servicesImplementation

import (
	"context"

	"github.com/nkarakotova/lim-core/errors/servicesErrors"
	"github.com/nkarakotova/lim-core/repositories"
	"github.com/nkarakotova/lim-core/services"

	"github.com/nkarakotova/lim-core/models"

	"github.com/charmbracelet/log"
)

type SubscriptionServiceImplementation struct {
	SubscriptionRepository repositories.SubscriptionRepository
	logger                 *log.Logger
}

func NewSubscriptionServiceImplementation(
	SubscriptionRepository repositories.SubscriptionRepository,
	logger *log.Logger,
) services.SubscriptionService {

	return &SubscriptionServiceImplementation{
		SubscriptionRepository: SubscriptionRepository,
		logger:                 logger,
	}
}

func (s *SubscriptionServiceImplementation) validate(subscription *models.Subscription) error {
	if subscription.StartDate.After(subscription.EndDate) {
		return servicesErrors.SubscriptionStartDateAfterEndDate
	}

	return nil
}

func (s *SubscriptionServiceImplementation) Create(subscription *models.Subscription, clientID uint64) error {
	ctx := context.Background()

	err := s.validate(subscription)
	if err != nil {
		return err
	}

	err = s.SubscriptionRepository.Create(ctx, subscription, clientID)
	if err != nil {
		s.logger.Warn("SUBSCRIPTION! Error in repository Create", "id", subscription.ID, "error", err)
		return err
	}

	s.logger.Info("SUBSCRIPTION! Successfully create subscription", "id", subscription.ID)
	return nil
}

func (s *SubscriptionServiceImplementation) GetByID(id uint64) (*models.Subscription, error) {
	ctx := context.Background()

	subscription, err := s.SubscriptionRepository.GetByID(ctx, id)
	if err != nil {
		s.logger.Warn("SUBSCRIPTION! Error in repository method GetByID", "id", id, "error", err)
		return nil, err
	}

	s.logger.Debug("SUBSCRIPTION! Success repository method GetByID", "id", id)
	return subscription, nil
}

func (s *SubscriptionServiceImplementation) ReduceRemainingTrainingsNum(id uint64) error {
	ctx := context.Background()

	err := s.SubscriptionRepository.ReduceRemainingTrainingsNum(ctx, id)
	if err != nil {
		s.logger.Warn("SUBSCRIPTION! Error in repository method ReduceRemainingTrainingsNum", "id", id, "error", err)
		return err
	}

	s.logger.Debug("SUBSCRIPTION! Success repository method ReduceRemainingTrainingsNum", "id", id)
	return nil
}
