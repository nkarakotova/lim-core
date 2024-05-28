package servicesImplementation

import (
	"context"
	"time"

	"github.com/nkarakotova/lim-core/errors/servicesErrors"
	"github.com/nkarakotova/lim-core/managers"
	"github.com/nkarakotova/lim-core/repositories"
	"github.com/nkarakotova/lim-core/services"

	"github.com/nkarakotova/lim-core/models"

	"github.com/charmbracelet/log"
)

type TrainingServiceImplementation struct {
	TrainingRepository     repositories.TrainingRepository
	ClientRepository       repositories.ClientRepository
	CoachRepository        repositories.CoachRepository
	SubscriptionRepository repositories.SubscriptionRepository
	HallRepository         repositories.HallRepository
	TransactionManager     managers.TransactionManager
	logger                 *log.Logger
}

func NewTrainingServiceImplementation(
	TrainingRepository repositories.TrainingRepository,
	ClientRepository repositories.ClientRepository,
	CoachRepository repositories.CoachRepository,
	SubscriptionRepository repositories.SubscriptionRepository,
	HallRepository repositories.HallRepository,
	TransactionManager managers.TransactionManager,
	logger *log.Logger,
) services.TrainingService {

	return &TrainingServiceImplementation{
		TrainingRepository:     TrainingRepository,
		ClientRepository:       ClientRepository,
		CoachRepository:        CoachRepository,
		SubscriptionRepository: SubscriptionRepository,
		HallRepository:         HallRepository,
		TransactionManager:     TransactionManager,
		logger:                 logger,
	}
}

func (t *TrainingServiceImplementation) validate(ctx context.Context, training *models.Training) error {
	h, m, s := training.DateTime.Clock()
	if h < services.FirstTrainingTime || h > services.LastTrainingTime || m != 0 || s != 0 {
		t.logger.Warn("TRAINING! Incorrect start time", "id", training.ID)
		return servicesErrors.IncorrectTrainingTime
	}

	hall, err := t.HallRepository.GetByID(ctx, training.HallID)
	if err != nil {
		t.logger.Warn("HALL! Error in repository GetByID", "id", training.HallID, "error", err)
		return err
	}
	if training.PlacesNum > hall.Capacity {
		t.logger.Warn("HALL! Places num more then capacity", "id", training.HallID, "error", err)
		return servicesErrors.PlacesNumMoreThenCapacity
	}

	coaches, err := t.CoachRepository.GetAllByDirection(ctx, training.DirectionID)
	if err != nil {
		t.logger.Warn("TRAINING! Error in repository GetAllByDirection", "id", training.ID, "error", err)
		return err
	}

	var flag = false
	for _, c := range coaches {
		if c.ID == training.CoachID {
			flag = true
		}
	}
	if !flag {
		t.logger.Warn("TRAINING! Not coach by direction", "id", training.ID, "error", err)
		return servicesErrors.NotCoachByDirection
	}

	trainings, err := t.GetAllByDateTime(training.DateTime)
	if err != nil {
		t.logger.Warn("TRAINING! Error in repository GetAllByDateTime", "id", training.ID, "error", err)
		return err
	}
	for _, tr := range trainings {
		if tr.CoachID == training.CoachID || tr.HallID == training.HallID {
			t.logger.Warn("TRAINING! Bysy date time", "id", training.ID)
			return servicesErrors.BysyDateTime
		}
	}

	return nil
}

func (t *TrainingServiceImplementation) Create(training *models.Training) error {
	ctx := context.Background()

	err := t.validate(ctx, training)
	if err != nil {
		return err
	}

	err = t.TrainingRepository.Create(ctx, training)
	if err != nil {
		t.logger.Warn("TRAINING! Error in repository Create", "id", training.ID, "error", err)
		return err
	}

	t.logger.Info("TRAINING! Success create training", "id", training.ID)
	return nil
}

func (t *TrainingServiceImplementation) delete(ctx context.Context, clients []models.Client, id uint64) error {
	return t.TransactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		for _, c := range clients {
			err := t.SubscriptionRepository.IncreaseRemainingTrainingsNum(ctx, c.SubscriptionID)
			if err != nil {
				t.logger.Warn("SUBSCRIPTION! Error in repository IncreaseRemainingTrainingsNum")
				return err
			}
		}

		err := t.TrainingRepository.Delete(ctx, id)
		if err != nil {
			t.logger.Warn("TRAINING! Error in repository Delete", "id", id, "error", err)
			return err
		}

		return nil
	})
}

func (t *TrainingServiceImplementation) Delete(id uint64) error {
	ctx := context.Background()

	clients, err := t.ClientRepository.GetByTraining(ctx, id)
	if err != nil {
		t.logger.Warn("CLIENT! Error in repository GetByTraining")
		return err
	}

	err = t.delete(ctx, clients, id)
	if err != nil {
		return err
	}

	t.logger.Info("TRAINING! Success delete training", "id", id)
	return nil
}

func (t *TrainingServiceImplementation) GetByID(id uint64) (*models.Training, error) {
	ctx := context.Background()

	training, err := t.TrainingRepository.GetByID(ctx, id)

	if err != nil {
		t.logger.Warn("TRAINING! Error in repository method GetByID", "id", id, "error", err)
		return nil, err
	}

	t.logger.Debug("TRAINING! Success GetByID", "id", id)
	return training, nil
}

func (t *TrainingServiceImplementation) GetAllByClient(id uint64) ([]models.Training, error) {
	ctx := context.Background()

	trainings, err := t.TrainingRepository.GetAllByClient(ctx, id)

	if err != nil {
		t.logger.Warn("TRAINING! Error in repository method GetAllByClient", "id", id, "err", err)
		return nil, err
	}

	t.logger.Debug("TRAINING! Success GetAllByClient", "id", id)
	return trainings, nil
}

func (t *TrainingServiceImplementation) GetAllByCoachOnDate(id uint64, date time.Time) ([]models.Training, error) {
	ctx := context.Background()

	trainings, err := t.TrainingRepository.GetAllByCoachOnDate(ctx, id, date)
	if err != nil {
		t.logger.Warn("TRAINING! Error in service method GetAllByCoachOnDate", "id", id, "err", err)
		return nil, err
	}

	t.logger.Debug("TRAINING! Successfully service method GetAllByCoachOnDate", "id", id)
	return trainings, nil
}

func (t *TrainingServiceImplementation) GetAllByDateTime(dateTime time.Time) ([]models.Training, error) {
	ctx := context.Background()

	trainings, err := t.TrainingRepository.GetAllByDateTime(ctx, dateTime)

	if err != nil {
		t.logger.Warn("TRAINING! Error in repository method GetAllByDateTime", "dateTime", dateTime, "err", err)
		return nil, err
	}

	t.logger.Debug("TRAINING! Success GetAllByDateTime", "dateTime", dateTime)
	return trainings, nil
}

func (t *TrainingServiceImplementation) GetAllBetweenDateTime(start time.Time, end time.Time) ([]models.Training, error) {
	if start.After(end) {
		t.logger.Warn("TRAINING! Start after end", "start", start, "end", end)
		return nil, servicesErrors.StartDateAfterEndDate
	}

	ctx := context.Background()

	trainings, err := t.TrainingRepository.GetAllBetweenDateTime(ctx, start, end)
	if err != nil {
		t.logger.Warn("TRAINING! Error in repository method GetAllBetweenDateTime", "start", start, "end", end, "err", err)
		return nil, err
	}

	t.logger.Debug("TRAINING! Success GetAllBetweenDateTime", "start", start, "end", end)
	return trainings, nil
}

func (t *TrainingServiceImplementation) ReduceAvailablePlacesNum(id uint64) error {
	ctx := context.Background()

	err := t.TrainingRepository.ReduceAvailablePlacesNum(ctx, id)
	if err != nil {
		t.logger.Warn("TRAINING! Error in repository method ReduceAvailablePlacesNum", "id", id, "error", err)
		return err
	}

	t.logger.Info("TRAINING! Success ReduceAvailablePlacesNum", "id", id)
	return nil
}
