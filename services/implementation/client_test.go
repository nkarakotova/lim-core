package servicesImplementation

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"
	"github.com/nkarakotova/lim-core/errors/servicesErrors"
	"github.com/nkarakotova/lim-core/services"

	managers_mocks "github.com/nkarakotova/lim-core/managers/mocks"
	repositories_mocks "github.com/nkarakotova/lim-core/repositories/mocks"

	"github.com/nkarakotova/lim-core/models"

	"github.com/charmbracelet/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockClientService struct {
	mockClientRepository       *repositories_mocks.MockClientRepository
	mockTrainingRepository     *repositories_mocks.MockTrainingRepository
	mockDirectionRepository    *repositories_mocks.MockDirectionRepository
	mockSubscriptionRepository *repositories_mocks.MockSubscriptionRepository
	mockTransactionManager     *managers_mocks.MockTransactionManager
	logger                     *log.Logger
}

func createMockClientService(controller *gomock.Controller) *mockClientService {
	service := new(mockClientService)

	service.mockClientRepository = repositories_mocks.NewMockClientRepository(controller)
	service.mockTrainingRepository = repositories_mocks.NewMockTrainingRepository(controller)
	service.mockDirectionRepository = repositories_mocks.NewMockDirectionRepository(controller)
	service.mockSubscriptionRepository = repositories_mocks.NewMockSubscriptionRepository(controller)
	service.mockTransactionManager = managers_mocks.NewMockTransactionManager(controller)
	service.logger = log.New(os.Stderr)

	return service
}

func createClientService(service *mockClientService) services.ClientService {
	return NewClientServiceImplementation(service.mockClientRepository, service.mockTrainingRepository, service.mockDirectionRepository, service.mockSubscriptionRepository, service.mockTransactionManager, service.logger)
}

//-------------------------------------------------------------------------------------------------
// create

var testCreateSuccess = []struct {
	TestName  string
	InputData struct {
		client *models.Client
	}
	Prepare     func(service *mockClientService)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "simple create",
		InputData: struct {
			client *models.Client
		}{&models.Client{SubscriptionID: 7,
			Name:      "Natali",
			Telephone: "9262218276",
			Mail:      "nka@mail.ru",
			Password:  "111",
			Age:       20,
			Gender:    models.Female}},

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByTelephone(ctx, "9262218276").Return(nil, repositoriesErrors.EntityDoesNotExists)

			service.mockClientRepository.EXPECT().Create(ctx,
				&models.Client{SubscriptionID: 7,
					Name:      "Natali",
					Telephone: "9262218276",
					Mail:      "nka@mail.ru",
					Password:  "111",
					Age:       20,
					Gender:    models.Female}).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
}

var testCreateFailure = []struct {
	TestName  string
	InputData struct {
		client *models.Client
	}
	Prepare     func(service *mockClientService)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "create error, telephone number already exists",
		InputData: struct {
			client *models.Client
		}{&models.Client{SubscriptionID: 7,
			Name:      "Natali",
			Telephone: "9262218276",
			Mail:      "nka@mail.ru",
			Password:  "111",
			Age:       20,
			Gender:    models.Female}},

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByTelephone(ctx, "9262218276").Return(
				&models.Client{SubscriptionID: 7,
					Name:      "Natali",
					Telephone: "9262218276",
					Mail:      "nka@mail.ru",
					Password:  "111",
					Age:       20,
					Gender:    models.Female}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.ErrorIs(t, err, servicesErrors.ClientAlreadyExists)
		},
	},
	{
		TestName: "create error, incorrect telephone number length",
		InputData: struct {
			client *models.Client
		}{&models.Client{SubscriptionID: 7,
			Name:      "Natali",
			Telephone: "92622182767",
			Mail:      "nka@mail.ru",
			Password:  "111",
			Age:       20,
			Gender:    models.Female}},

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByTelephone(ctx, "92622182767").Return(nil, repositoriesErrors.EntityDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.ErrorIs(t, err, servicesErrors.ClientTelephoneIncorrect)
		},
	},
	{
		TestName: "create error, letter in telephone number",
		InputData: struct {
			client *models.Client
		}{&models.Client{SubscriptionID: 7,
			Name:      "Natali",
			Telephone: "926221827g",
			Mail:      "nka@mail.ru",
			Password:  "111",
			Age:       20,
			Gender:    models.Female}},

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByTelephone(ctx, "926221827g").Return(nil, repositoriesErrors.EntityDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.ErrorIs(t, err, servicesErrors.ClientTelephoneIncorrect)
		},
	},
	{
		TestName: "create error, incorrect mail",
		InputData: struct {
			client *models.Client
		}{&models.Client{SubscriptionID: 7,
			Name:      "Natali",
			Telephone: "9262218276",
			Mail:      "nkamail.ru",
			Password:  "111",
			Age:       20,
			Gender:    models.Female}},

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByTelephone(ctx, "9262218276").Return(nil, repositoriesErrors.EntityDoesNotExists)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.ErrorIs(t, err, servicesErrors.ClientMailIncorrect)
		},
	},
}

func TestClientServiceImplementationCreate(t *testing.T) {
	for _, tt := range testCreateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := createMockClientService(ctrl)
			tt.Prepare(service)

			clientService := createClientService(service)

			err := clientService.Create(tt.InputData.client)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testCreateFailure {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := createMockClientService(ctrl)
			tt.Prepare(service)

			clientService := createClientService(service)

			err := clientService.Create(tt.InputData.client)

			tt.CheckOutput(t, err)
		})
	}
}

//-------------------------------------------------------------------------------------------------
// create assignment

var testCreateAssignmentSuccess = []struct {
	TestName  string
	InputData struct {
		clientID   uint64
		trainingID uint64
	}
	Prepare     func(service *mockClientService)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "simple create",
		InputData: struct {
			clientID   uint64
			trainingID uint64
		}{clientID: 3, trainingID: 7},

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByID(ctx, uint64(3)).Return(
				&models.Client{
					ID:             3,
					SubscriptionID: 5,
					Age:            20,
					Gender:         models.Female,
				}, nil)

			service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(7)).Return(
				&models.Training{
					ID:                 7,
					DirectionID:        5,
					DateTime:           time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC),
					HallID:             111,
					AvailablePlacesNum: 30,
				}, nil)

			service.mockSubscriptionRepository.EXPECT().GetByID(ctx, uint64(5)).Return(
				&models.Subscription{
					ID:                    5,
					RemainingTrainingsNum: 10,
					StartDate:             time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:               time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
				}, nil)

			service.mockDirectionRepository.EXPECT().GetByID(ctx, uint64(5)).Return(
				&models.Direction{
					ID:               22,
					AcceptableGender: models.Female,
				}, nil)

			service.mockTrainingRepository.EXPECT().GetAllByClient(ctx, uint64(3)).Return(
				[]models.Training{
					{
						ID:                 17,
						DirectionID:        22,
						DateTime:           time.Date(2024, 3, 5, 14, 0, 0, 0, time.UTC),
						HallID:             111,
						AvailablePlacesNum: 30,
					},
					{
						ID:                 10,
						DirectionID:        22,
						DateTime:           time.Date(2024, 3, 5, 17, 0, 0, 0, time.UTC),
						HallID:             111,
						AvailablePlacesNum: 30,
					},
				}, nil)

			service.mockTransactionManager.EXPECT().WithinTransaction(ctx, gomock.Any()).Return(nil)

		},
		CheckOutput: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
}

var testCreateAssignmentFailure = []struct {
	TestName  string
	InputData struct {
		clientID   uint64
		trainingID uint64
	}
	Prepare     func(service *mockClientService)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "create error, no avaliable places number",
		InputData: struct {
			clientID   uint64
			trainingID uint64
		}{clientID: 3, trainingID: 7},

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByID(ctx, uint64(3)).Return(
				&models.Client{
					ID:             3,
					SubscriptionID: 5,
					Age:            20,
					Gender:         models.Female,
				}, nil)

			service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(7)).Return(
				&models.Training{
					ID:                 7,
					DirectionID:        5,
					DateTime:           time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC),
					HallID:             111,
					AvailablePlacesNum: 0,
				}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.ErrorIs(t, err, servicesErrors.NoAvailablePlacesNum)
		},
	},
	{
		TestName: "create error, no subscription",
		InputData: struct {
			clientID   uint64
			trainingID uint64
		}{clientID: 3, trainingID: 7},

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByID(ctx, uint64(3)).Return(
				&models.Client{
					ID:             3,
					SubscriptionID: 0,
					Age:            20,
					Gender:         models.Female,
				}, nil)

			service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(7)).Return(
				&models.Training{
					ID:                 7,
					DirectionID:        5,
					DateTime:           time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC),
					HallID:             111,
					AvailablePlacesNum: 30,
				}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.ErrorIs(t, err, servicesErrors.ClientHasntGotSubscription)
		},
	},
	{
		TestName: "create error, no remaining trainings",
		InputData: struct {
			clientID   uint64
			trainingID uint64
		}{clientID: 3, trainingID: 7},

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByID(ctx, uint64(3)).Return(
				&models.Client{
					ID:             3,
					SubscriptionID: 5,
					Age:            20,
					Gender:         models.Female,
				}, nil)

			service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(7)).Return(
				&models.Training{
					ID:                 7,
					DirectionID:        5,
					DateTime:           time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC),
					HallID:             111,
					AvailablePlacesNum: 30,
				}, nil)

			service.mockSubscriptionRepository.EXPECT().GetByID(ctx, uint64(5)).Return(
				&models.Subscription{
					ID:                    5,
					RemainingTrainingsNum: 0,
					StartDate:             time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:               time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.ErrorIs(t, err, servicesErrors.ClientSubscriptionIsOver)
		},
	},
	{
		TestName: "create error, subscription is over",
		InputData: struct {
			clientID   uint64
			trainingID uint64
		}{clientID: 3, trainingID: 7},

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByID(ctx, uint64(3)).Return(
				&models.Client{
					ID:             3,
					SubscriptionID: 5,
					Age:            20,
					Gender:         models.Female,
				}, nil)

			service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(7)).Return(
				&models.Training{
					ID:                 7,
					DirectionID:        5,
					DateTime:           time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC),
					HallID:             111,
					AvailablePlacesNum: 30,
				}, nil)

			service.mockSubscriptionRepository.EXPECT().GetByID(ctx, uint64(5)).Return(
				&models.Subscription{
					ID:                    5,
					RemainingTrainingsNum: 10,
					StartDate:             time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:               time.Date(2024, 3, 3, 0, 0, 0, 0, time.UTC),
				}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.ErrorIs(t, err, servicesErrors.ClientSubscriptionIsOver)
		},
	},
	{
		TestName: "create error, not available age",
		InputData: struct {
			clientID   uint64
			trainingID uint64
		}{clientID: 3, trainingID: 7},

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByID(ctx, uint64(3)).Return(
				&models.Client{
					ID:             3,
					SubscriptionID: 5,
					Age:            20,
					Gender:         models.Female,
				}, nil)

			service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(7)).Return(
				&models.Training{
					ID:                 7,
					DirectionID:        5,
					DateTime:           time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC),
					HallID:             111,
					AvailablePlacesNum: 30,
					AcceptableAge:      30,
				}, nil)

			service.mockSubscriptionRepository.EXPECT().GetByID(ctx, uint64(5)).Return(
				&models.Subscription{
					ID:                    5,
					RemainingTrainingsNum: 10,
					StartDate:             time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:               time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.ErrorIs(t, err, servicesErrors.AgeNotCorrespondToAcceptableAge)
		},
	},
	{
		TestName: "create error, not available gender",
		InputData: struct {
			clientID   uint64
			trainingID uint64
		}{clientID: 3, trainingID: 7},

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByID(ctx, uint64(3)).Return(
				&models.Client{
					ID:             3,
					SubscriptionID: 5,
					Age:            20,
					Gender:         models.Female,
				}, nil)

			service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(7)).Return(
				&models.Training{
					ID:                 7,
					DirectionID:        5,
					DateTime:           time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC),
					HallID:             111,
					AvailablePlacesNum: 30,
				}, nil)

			service.mockSubscriptionRepository.EXPECT().GetByID(ctx, uint64(5)).Return(
				&models.Subscription{
					ID:                    5,
					RemainingTrainingsNum: 10,
					StartDate:             time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:               time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
				}, nil)

			service.mockDirectionRepository.EXPECT().GetByID(ctx, uint64(5)).Return(
				&models.Direction{
					ID:               22,
					AcceptableGender: models.Male,
				}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.ErrorIs(t, err, servicesErrors.GenderNotCorrespondToAcceptableGender)
		},
	},
	{
		TestName: "create error, has assigenment on this time",
		InputData: struct {
			clientID   uint64
			trainingID uint64
		}{clientID: 3, trainingID: 7},

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByID(ctx, uint64(3)).Return(
				&models.Client{
					ID:             3,
					SubscriptionID: 5,
					Age:            20,
					Gender:         models.Female,
				}, nil)

			service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(7)).Return(
				&models.Training{
					ID:                 7,
					DirectionID:        5,
					DateTime:           time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC),
					HallID:             111,
					AvailablePlacesNum: 30,
				}, nil)

			service.mockSubscriptionRepository.EXPECT().GetByID(ctx, uint64(5)).Return(
				&models.Subscription{
					ID:                    5,
					RemainingTrainingsNum: 10,
					StartDate:             time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:               time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
				}, nil)

			service.mockDirectionRepository.EXPECT().GetByID(ctx, uint64(5)).Return(
				&models.Direction{
					ID:               22,
					AcceptableGender: models.Female,
				}, nil)

			service.mockTrainingRepository.EXPECT().GetAllByClient(ctx, uint64(3)).Return(
				[]models.Training{
					{
						ID:                 17,
						DirectionID:        22,
						DateTime:           time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC),
						HallID:             111,
						AvailablePlacesNum: 30,
					},
					{
						ID:                 10,
						DirectionID:        22,
						DateTime:           time.Date(2024, 3, 5, 17, 0, 0, 0, time.UTC),
						HallID:             111,
						AvailablePlacesNum: 30,
					},
				}, nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.ErrorIs(t, err, servicesErrors.AssignmentOnThisTimeAlreadyExists)
		},
	},
}

func TestClientServiceImplementationCreateAssignment(t *testing.T) {
	for _, tt := range testCreateAssignmentSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := createMockClientService(ctrl)
			tt.Prepare(service)

			clientService := createClientService(service)

			err := clientService.СreateAssignment(tt.InputData.clientID, tt.InputData.trainingID)

			tt.CheckOutput(t, err)
		})
	}

	for _, tt := range testCreateAssignmentFailure {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := createMockClientService(ctrl)
			tt.Prepare(service)

			clientService := createClientService(service)

			err := clientService.СreateAssignment(tt.InputData.clientID, tt.InputData.trainingID)

			tt.CheckOutput(t, err)
		})
	}
}

//-------------------------------------------------------------------------------------------------
