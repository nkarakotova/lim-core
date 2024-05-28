package servicesImplementation

import (
	"context"
	"time"

	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"
	"github.com/nkarakotova/lim-core/errors/servicesErrors"
	"github.com/nkarakotova/lim-core/repositories"
	"github.com/nkarakotova/lim-core/services"

	"github.com/nkarakotova/lim-core/models"

	"github.com/charmbracelet/log"
)

type CoachServiceImplementation struct {
	CoachRepository    repositories.CoachRepository
	TrainingRepository repositories.TrainingRepository
	logger             *log.Logger
}

func NewCoachServiceImplementation(
	CoachRepository repositories.CoachRepository,
	TrainingRepository repositories.TrainingRepository,
	logger *log.Logger,
) services.CoachService {

	return &CoachServiceImplementation{
		CoachRepository:    CoachRepository,
		TrainingRepository: TrainingRepository,
		logger:             logger,
	}
}

func (c *CoachServiceImplementation) validate(ctx context.Context, coach *models.Coach) error {
	_, err := c.CoachRepository.GetByName(ctx, coach.Name)
	if err != nil && err != repositoriesErrors.EntityDoesNotExists {
		c.logger.Warn("COACH! Error in repository GetByName", "name", coach.Name, "error", err)
		return err
	} else if err == nil {
		c.logger.Warn("COACH! Coach already exists", "name", coach.Name)
		return servicesErrors.CoachAlreadyExists
	}

	return nil
}

func (c *CoachServiceImplementation) GetByName(name string) (*models.Coach, error) {
	ctx := context.Background()

	coach, err := c.CoachRepository.GetByName(ctx, name)
	if err != nil {
		c.logger.Warn("COACH! Error in repository GetByName", "name", name, "error", err)
		return nil, err
	}

	c.logger.Debug("COACH! Success GetByName", "name", name)
	return coach, nil
}

func (c *CoachServiceImplementation) Create(coach *models.Coach) error {
	ctx := context.Background()

	err := c.validate(ctx, coach)
	if err != nil {
		return err
	}

	err = c.CoachRepository.Create(ctx, coach)
	if err != nil {
		c.logger.Warn("COACH! Error in repository Create", "name", coach.Name, "error", err)
		return err
	}

	c.logger.Info("COACH! Success create coach", "name", coach.Name)
	return nil
}

func (c *CoachServiceImplementation) GetByID(id uint64) (*models.Coach, error) {
	ctx := context.Background()

	coach, err := c.CoachRepository.GetByID(ctx, id)
	if err != nil {
		c.logger.Warn("COACH! Error in repository method GetByID", "id", id, "error", err)
		return nil, err
	}

	c.logger.Debug("COACH! Success GetByID", "id", id)
	return coach, nil
}

func (c *CoachServiceImplementation) AddDirection(coachID, directionID uint64) error {
	ctx := context.Background()

	err := c.CoachRepository.AddDirection(ctx, coachID, directionID)
	if err != nil {
		c.logger.Warn("COACH! Error in repository method AddDirection", "id", coachID, "error", err)
		return err
	}

	c.logger.Debug("COACH! Success AddDirection", "id", coachID)
	return nil
}

func (c *CoachServiceImplementation) GetAllByDirection(id uint64) ([]models.Coach, error) {
	ctx := context.Background()

	coaches, err := c.CoachRepository.GetAllByDirection(ctx, id)
	if err != nil {
		c.logger.Warn("COACH! Error in repository method GetAllByDirection", "id", id, "err", err)
		return nil, err
	}

	c.logger.Debug("COACH! Success GetAllByDirection", "id", id)
	return coaches, nil
}

func (c *CoachServiceImplementation) getAllSlots(date time.Time) []time.Time {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	slots := make([]time.Time, services.LastTrainingTime-services.FirstTrainingTime)
	slots[0] = date.Add(time.Duration(services.FirstTrainingTime) * time.Hour)

	for idx := 1; idx < len(slots); idx++ {
		slots[idx] = slots[idx-1].Add(time.Hour)
	}

	return slots
}

func (c *CoachServiceImplementation) GetFreeTimeOnDate(id uint64, date time.Time) ([]time.Time, error) {
	ctx := context.Background()

	trainings, err := c.TrainingRepository.GetAllByCoachOnDate(ctx, id, date)
	if err != nil {
		c.logger.Warn("TRAINING! Error in repository method GetAllByCoachOnDate", "id", id, "err", err)
		return nil, err
	}

	slots := c.getAllSlots(date)

	for _, t := range trainings {
		time := t.DateTime
		for i, slot := range slots {
			if time.Equal(slot) {
				slots = append(slots[:i], slots[i+1:]...)
				break
			}
		}
		if len(slots) == 0 {
			break
		}
	}
	
	c.logger.Debug("COACH! Success GetFreeTimeOnDate", "id", id)
	return slots, nil
}
