package services

import "github.com/nkarakotova/lim-core/models"

type ClientService interface {
	Create(client *models.Client) error
	Login(telephone, password string) (*models.Client, error)
	GetByID(id uint64) (*models.Client, error)
	GetByTelephone(login string) (*models.Client, error)
	СreateAssignment(clientID, trainingID uint64) error
	DeleteAssignment(clientID, trainingID uint64) error
}
