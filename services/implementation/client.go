package servicesImplementation

import (
	"context"
	"net/mail"
	"strconv"

	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"
	"github.com/nkarakotova/lim-core/errors/servicesErrors"
	"github.com/nkarakotova/lim-core/managers"
	"github.com/nkarakotova/lim-core/repositories"
	"github.com/nkarakotova/lim-core/services"

	"github.com/nkarakotova/lim-core/models"

	"github.com/charmbracelet/log"
)

type ClientServiceImplementation struct {
	ClientRepository       repositories.ClientRepository
	TrainingRepository     repositories.TrainingRepository
	DirectionRepository    repositories.DirectionRepository
	SubscriptionRepository repositories.SubscriptionRepository
	TransactionManager     managers.TransactionManager
	logger                 *log.Logger
}

func NewClientServiceImplementation(
	ClientRepository repositories.ClientRepository,
	TrainingRepository repositories.TrainingRepository,
	DirectionRepository repositories.DirectionRepository,
	SubscriptionRepository repositories.SubscriptionRepository,
	TransactionManager managers.TransactionManager,
	logger *log.Logger,
) services.ClientService {

	return &ClientServiceImplementation{
		ClientRepository:       ClientRepository,
		TrainingRepository:     TrainingRepository,
		DirectionRepository:    DirectionRepository,
		SubscriptionRepository: SubscriptionRepository,
		TransactionManager:     TransactionManager,
		logger:                 logger,
	}
}

const TelephoneNumberLen = 10

func (c *ClientServiceImplementation) validate(ctx context.Context, client *models.Client) error {
	_, err := c.ClientRepository.GetByTelephone(ctx, client.Telephone)
	if err != nil && err != repositoriesErrors.EntityDoesNotExists {
		c.logger.Warn("CLIENT! Error in repository GetClientByTelephone", "telephone", client.Telephone, "error", err)
		return err
	} else if err == nil {
		c.logger.Warn("CLIENT! Client already exists", "telephone", client.Telephone)
		return servicesErrors.ClientAlreadyExists
	}

	if len(client.Telephone) != TelephoneNumberLen {
		return servicesErrors.ClientTelephoneIncorrect
	}

	_, err = strconv.Atoi(client.Telephone)
	if err != nil {
		return servicesErrors.ClientTelephoneIncorrect
	}

	_, err = mail.ParseAddress(client.Mail)
	if err != nil {
		return servicesErrors.ClientMailIncorrect
	}

	return nil
}

func (c *ClientServiceImplementation) GetByTelephone(telephone string) (*models.Client, error) {
	ctx := context.Background()

	client, err := c.ClientRepository.GetByTelephone(ctx, telephone)
	if err != nil {
		c.logger.Warn("CLIENT! Error in repository GetClientByTelephone", "telephone", telephone, "error", err)
		return nil, err
	}

	c.logger.Debug("CLIENT! Successfully GetClientByTelephone", "telephone", telephone)
	return client, nil
}

func (c *ClientServiceImplementation) Create(client *models.Client) error {
	ctx := context.Background()

	err := c.validate(ctx, client)
	if err != nil {
		return err
	}

	err = c.ClientRepository.Create(ctx, client)
	if err != nil {
		c.logger.Warn("CLIENT! Error in repository Create", "telephone", client.Telephone, "error", err)
		return err
	}

	c.logger.Info("CLIENT! Successfully create client", "login", client.Telephone, "id", client.ID)
	return nil
}

func (c *ClientServiceImplementation) Login(telepnone, password string) (*models.Client, error) {
	ctx := context.Background()

	tempClient, err := c.ClientRepository.GetByTelephone(ctx, telepnone)
	if err != nil && err == repositoriesErrors.EntityDoesNotExists {
		c.logger.Warn("CLIENT! Error, client with this telephone does not exists", "telephone", telepnone, "error", err)
		return nil, servicesErrors.ClientDoesNotExists
	} else if err != nil {
		c.logger.Warn("CLIENT! Error in repository method GetByTelephone", "telepnone", telepnone, "error", err)
		return nil, err
	}

	if password != tempClient.Password {
		c.logger.Warn("CLIENT! Error client password", "telephone", telepnone)
		return nil, servicesErrors.InvalidPassword
	}

	c.logger.Info("CLIENT! Success login with", "telepnone", telepnone, "id", tempClient.ID)
	return tempClient, nil
}

func (c *ClientServiceImplementation) GetByID(id uint64) (*models.Client, error) {
	ctx := context.Background()

	client, err := c.ClientRepository.GetByID(ctx, id)
	if err != nil {
		c.logger.Warn("CLIENT! Error in repository method GetByID", "id", id, "error", err)
		return nil, err
	}

	c.logger.Debug("CLIENT! Success repository method GetByID", "id", id)
	return client, nil
}

func (c *ClientServiceImplementation) trainingInSameDateTimeCheck(ctx context.Context, client *models.Client, training *models.Training) error {
	curStart := training.DateTime

	clientTrainings, err := c.TrainingRepository.GetAllByClient(ctx, client.ID)
	if err != nil {
		c.logger.Warn("TRAINING! Error in repository method GetAllByClient", "id", client.ID, "error", err)
		return err
	}

	for _, t := range clientTrainings {
		if t.DateTime == curStart {
			c.logger.Warn("CLIENT! There is already an assignment for this time", "id", client.ID, "error", err)
			return servicesErrors.AssignmentOnThisTimeAlreadyExists
		}
	}

	return nil
}

func (c *ClientServiceImplementation) subscriptionCheck(ctx context.Context, client *models.Client, training *models.Training) error {
	if client.SubscriptionID == 0 {
		c.logger.Warn("CLIENT! Client hasn't got subscription", "id", client.ID)
		return servicesErrors.ClientHasntGotSubscription
	}

	subscription, err := c.SubscriptionRepository.GetByID(ctx, client.SubscriptionID)
	if err != nil {
		c.logger.Warn("SUBSCRIPTION! Error in repository method GetByID", "id", client.SubscriptionID, "error", err)
		return err
	}

	if subscription.RemainingTrainingsNum == 0 {
		c.logger.Warn("SUBSCRIPTION! Subscription is over", "id", client.ID, "error", err)
		return servicesErrors.ClientSubscriptionIsOver
	}

	curStart := training.DateTime
	if subscription.StartDate.After(curStart) || subscription.EndDate.Before(curStart) {
		c.logger.Warn("CLIENT! Subscription is over", "id", client.ID, "error", err)
		return servicesErrors.ClientSubscriptionIsOver
	}

	return nil
}

func (c *ClientServiceImplementation) genderCheck(ctx context.Context, client *models.Client, training *models.Training) error {
	direction, err := c.DirectionRepository.GetByID(ctx, training.DirectionID)
	if err != nil {
		c.logger.Warn("DIRECTION! Error in repository GetByID", "DirectionID", training.DirectionID, "error", err)
		return err
	}

	if (direction.AcceptableGender != models.Unknown) && (client.Gender != direction.AcceptableGender) {
		c.logger.Warn("CLIENT! Gender does not correspond to the acceptable gender", "id", client.ID, "error", err)
		return servicesErrors.GenderNotCorrespondToAcceptableGender
	}

	return nil
}

func (c *ClientServiceImplementation) createAssignmentChecks(ctx context.Context, client *models.Client, training *models.Training) error {
	if training.AvailablePlacesNum == 0 {
		c.logger.Warn("TRAINING! There is no available places num", "id", training.ID)
		return servicesErrors.NoAvailablePlacesNum
	}

	err := c.subscriptionCheck(ctx, client, training)
	if err != nil {
		return err
	}

	if client.Age < training.AcceptableAge {
		c.logger.Warn("CLIENT! Age does not correspond to the acceptable age", "id", client.ID, "error", err)
		return servicesErrors.AgeNotCorrespondToAcceptableAge
	}

	err = c.genderCheck(ctx, client, training)
	if err != nil {
		return err
	}

	err = c.trainingInSameDateTimeCheck(ctx, client, training)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientServiceImplementation) createAssignment(ctx context.Context, client *models.Client, training *models.Training) error {
	return c.TransactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		err := c.ClientRepository.СreateAssignment(ctx, client.ID, training.ID)
		if err != nil {
			c.logger.Warn("CLIENT! Error in repository СreateAssignment", "clientID", client.ID, "trainingID", training.ID, "error", err)
			return err
		}

		err = c.SubscriptionRepository.ReduceRemainingTrainingsNum(ctx, client.SubscriptionID)
		if err != nil {
			c.logger.Warn("SUBSCRIPTION! Error in repository ReduceRemainingTrainingsNum", "SubscriptionID", client.SubscriptionID, "error", err)
			return err
		}

		err = c.TrainingRepository.ReduceAvailablePlacesNum(ctx, training.ID)
		if err != nil {
			c.logger.Warn("TRAINING! Error in repository ReduceAvailablePlacesNum", "trainingID", training.ID, "error", err)
			return err
		}

		return nil
	})
}

func (c *ClientServiceImplementation) СreateAssignment(clientID, trainingID uint64) error {
	ctx := context.Background()

	client, err := c.ClientRepository.GetByID(ctx, clientID)
	if err != nil {
		c.logger.Warn("CLIENT! Error in repository method GetByID", "id", clientID, "error", err)
		return err
	}

	training, err := c.TrainingRepository.GetByID(ctx, trainingID)
	if err != nil {
		c.logger.Warn("TRAINING! Error in repository method GetByID", "id", trainingID, "error", err)
		return err
	}

	err = c.createAssignmentChecks(ctx, client, training)
	if err != nil {
		return err
	}

	err = c.createAssignment(ctx, client, training)
	if err != nil {
		return err
	}

	c.logger.Info("CLIENT! Successfully create assignment", "clientID", clientID, "trainingID", clientID)
	return nil
}

func (c *ClientServiceImplementation) DeleteAssignment(clientID, trainingID uint64) error {
	ctx := context.Background()

	client, err := c.ClientRepository.GetByID(ctx, clientID)
	if err != nil {
		c.logger.Warn("CLIENT! Error in repository method GetByID", "id", clientID, "error", err)
		return err
	}

	return c.TransactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		err = c.ClientRepository.DeleteAssignment(ctx, clientID, trainingID)
		if err != nil {
			c.logger.Warn("CLIENT! Error in repository DeleteAssignment", "clientID", clientID, "trainingID", clientID, "error", err)
			return err
		}

		err = c.SubscriptionRepository.IncreaseRemainingTrainingsNum(ctx, client.SubscriptionID)
		if err != nil {
			c.logger.Warn("SUBSCRIPTION! Error in repository IncreaseRemainingTrainingsNum", "SubscriptionID", client.SubscriptionID, "error", err)
			return err
		}

		err = c.TrainingRepository.IncreaseAvailablePlacesNum(ctx, trainingID)
		if err != nil {
			c.logger.Warn("TRAINING! Error in repository IncreaseAvailablePlacesNum", "trainingID", trainingID, "error", err)
			return err
		}

		c.logger.Info("CLIENT! Successfully delete assignment", "clientID", clientID, "trainingID", clientID)
		return nil
	})
}
