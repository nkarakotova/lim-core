package servicesImplementation

import (
	"context"

	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"
	"github.com/nkarakotova/lim-core/errors/servicesErrors"
	"github.com/nkarakotova/lim-core/repositories"
	"github.com/nkarakotova/lim-core/services"

	"github.com/nkarakotova/lim-core/models"

	"github.com/charmbracelet/log"
)

type DirectionServiceImplementation struct {
	DirectionRepository repositories.DirectionRepository
	logger              *log.Logger
}

func NewDirectionServiceImplementation(
	DirectionRepository repositories.DirectionRepository,
	logger *log.Logger,
) services.DirectionService {

	return &DirectionServiceImplementation{
		DirectionRepository: DirectionRepository,
		logger:              logger,
	}
}

func (d *DirectionServiceImplementation) validate(ctx context.Context, direction *models.Direction) error {
	_, err := d.DirectionRepository.GetByName(ctx, direction.Name)
	if err != nil && err != repositoriesErrors.EntityDoesNotExists {
		d.logger.Warn("DIRECTION! Error in repository GetByName", "name", direction.Name, "error", err)
		return err
	} else if err == nil {
		d.logger.Warn("DIRECTION! Direction already exists", "name", direction.Name)
		return servicesErrors.DirectionAlreadyExists
	}

	return nil
}

func (d *DirectionServiceImplementation) GetByName(name string) (*models.Direction, error) {
	ctx := context.Background()

	direction, err := d.DirectionRepository.GetByName(ctx, name)
	if err != nil {
		d.logger.Warn("DIRECTION! Error in repository GetByName", "name", name, "error", err)
		return nil, err
	}

	d.logger.Debug("DIRECTION! Success GetByName", "name", name)
	return direction, nil
}

func (d *DirectionServiceImplementation) Create(direction *models.Direction) error {
	ctx := context.Background()

	err := d.validate(ctx, direction)
	if err != nil {
		return err
	}

	err = d.DirectionRepository.Create(ctx, direction)
	if err != nil {
		d.logger.Warn("DIRECTION! Error in repository Create", "name", direction.Name, "error", err)
		return err
	}

	d.logger.Info("DIRECTION! Success create direction", "name", direction.Name)
	return nil
}

func (d *DirectionServiceImplementation) GetByID(id uint64) (*models.Direction, error) {
	ctx := context.Background()

	direction, err := d.DirectionRepository.GetByID(ctx, id)
	if err != nil {
		d.logger.Warn("DIRECTION! Error in repository method GetByID", "id", id, "error", err)
		return nil, err
	}

	d.logger.Debug("DIRECTION! Success GetByID", "id", id)
	return direction, nil
}

func (d *DirectionServiceImplementation) GetAll() ([]models.Direction, error) {
	ctx := context.Background()

	directions, err := d.DirectionRepository.GetAll(ctx)
	if err != nil {
		d.logger.Warn("DIRECTION! Error in repository method GetAll", "err", err)
		return nil, err
	}

	d.logger.Debug("DIRECTION! Success GetAll")
	return directions, nil
}
