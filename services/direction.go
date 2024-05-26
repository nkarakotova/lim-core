package services

import "github.com/nkarakotova/lim-core/models"

type DirectionService interface {
	Create(direction *models.Direction) error
	GetByID(id uint64) (*models.Direction, error)
	GetByName(name string) (*models.Direction, error)
	GetAll() ([]models.Direction, error)
}
