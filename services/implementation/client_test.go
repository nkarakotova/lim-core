package servicesImplementation

import (
	"context"
	"os"
	"testing"
	"errors"

	repositories_mocks "github.com/nkarakotova/lim-core/repositories/mocks"
	managers_mocks "github.com/nkarakotova/lim-core/managers/mocks"
	data_builders "github.com/nkarakotova/lim-core/services/implementation/data_builders"
	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"
	"github.com/nkarakotova/lim-core/errors/servicesErrors"
	"github.com/nkarakotova/lim-core/services"
	"github.com/nkarakotova/lim-core/models"

	"github.com/charmbracelet/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockClientService struct {
	mockClientRepository       *repositories_mocks.MockClientRepository
	mockTrainingRepository     *repositories_mocks.MockTrainingRepository
	mockTransactionManager     *managers_mocks.MockTransactionManager
	logger                     *log.Logger
}

func createMockClientService(controller *gomock.Controller) *mockClientService {
	service := new(mockClientService)

	service.mockClientRepository = repositories_mocks.NewMockClientRepository(controller)
	service.mockTrainingRepository = repositories_mocks.NewMockTrainingRepository(controller)
	service.mockTransactionManager = managers_mocks.NewMockTransactionManager(controller)
	service.logger = log.New(os.Stderr)

	return service
}

func createClientService(service *mockClientService) services.ClientService {
	return NewClientServiceImplementation(service.mockClientRepository, service.mockTrainingRepository, service.mockTransactionManager, service.logger)
}

var testGetByTelephone = []struct {
	TestName  string
	InputData string
	Prepare   func(service *mockClientService)
	CheckOutput func(t *testing.T, client *models.Client, err error)
}{
	{
		TestName:  "success get client by telephone",
		InputData: "1234567890",
		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByTelephone(ctx, "1234567890").
			Return(data_builders.NewClientBuilder().Build(), nil)
		},
		CheckOutput: func(t *testing.T, client *models.Client, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, client)
			assert.Equal(t, "1234567890", client.Telephone)
		},
	},
	{
		TestName:  "error getting client by telephone",
		InputData: "nonexistent",
		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByTelephone(ctx, "nonexistent").
			Return(nil, errors.New("client not found"))
		},
		CheckOutput: func(t *testing.T, client *models.Client, err error) {
			assert.Error(t, err)
			assert.Nil(t, client)
		},
	},
}

var testCreateClient = []struct {
	TestName  string
	InputData *models.Client
	Prepare   func(service *mockClientService)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName:  "success create client",
		InputData: data_builders.NewClientBuilder().Build(),
		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByTelephone(ctx, "1234567890").Return(nil, repositoriesErrors.EntityDoesNotExists)
			service.mockClientRepository.EXPECT().Create(ctx, data_builders.NewClientBuilder().Build()).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
	{
		TestName:  "error creating client",
		InputData: data_builders.NewClientBuilder().WithTelephone("invalid").Build(),
		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByTelephone(ctx, "invalid").Return(nil, errors.New("validation error"))
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.Error(t, err)
		},
	},
}

var testLogin = []struct {
	TestName   string
	Tel        string
	Password   string
	Prepare    func(service *mockClientService)
	CheckOutput func(t *testing.T, client *models.Client, err error)
}{
	{
		TestName:   "success login",
		Tel:        "1234567890",
		Password:   "123",
		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByTelephone(ctx, "1234567890").
				Return(data_builders.NewClientBuilder().Build(), nil)
		},
		CheckOutput: func(t *testing.T, client *models.Client, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, client)
		},
	},
	{
		TestName:   "login with incorrect password",
		Tel:        "1234567890",
		Password:   "111",
		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByTelephone(ctx, "1234567890").
				Return(data_builders.NewClientBuilder().Build(), nil)
		},
		CheckOutput: func(t *testing.T, client *models.Client, err error) {
			assert.Error(t, err)
			assert.Nil(t, client)
		},
	},
}

var testGetClientByID = []struct {
	TestName  string
	ID        uint64
	Prepare   func(service *mockClientService)
	CheckOutput func(t *testing.T, client *models.Client, err error)
}{
	{
		TestName:  "success get client by ID",
		ID:        1,
		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByID(ctx, uint64(1)).
				Return(data_builders.NewClientBuilder().Build(), nil)
		},
		CheckOutput: func(t *testing.T, client *models.Client, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, client)
			assert.Equal(t, uint64(1), client.ID)
		},
	},
	{
		TestName:  "error getting client by ID",
		ID:        999,
		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			service.mockClientRepository.EXPECT().GetByID(ctx, uint64(999)).
				Return(nil, errors.New("not found"))
		},
		CheckOutput: func(t *testing.T, client *models.Client, err error) {
			assert.Error(t, err)
			assert.Nil(t, client)
		},
	},
}

var testCreateAssignment = []struct {
	TestName        string
	ClientID        uint64
	TrainingID      uint64
	DirectionID     uint64
	SubscriptionID  uint64
	Prepare      func(service *mockClientService)
	CheckOutput  func(t *testing.T, err error)
}{
	{
		TestName:    "success create assignment",
		ClientID:    1,
		TrainingID:  1,
		DirectionID: 1,
		SubscriptionID: 1,

		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			client := data_builders.NewClientBuilder().Build()
			training := data_builders.NewTrainingBuilder().Build()

			service.mockClientRepository.EXPECT().GetByID(ctx, uint64(1)).Return(client, nil)
			service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(1)).Return(training, nil)
			service.mockTrainingRepository.EXPECT().GetAllByClient(ctx, client.ID).Return(nil, nil)
			service.mockTransactionManager.EXPECT().WithinTransaction(ctx, gomock.Any()).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
	{
		TestName:    "error create assignment due to no available places",
		ClientID:    1,
		TrainingID:  1,
		Prepare: func(service *mockClientService) {
			ctx := context.Background()
			client := data_builders.NewClientBuilder().Build()
			training := data_builders.NewTrainingBuilder().WithPlacesNum(0).Build()

			service.mockClientRepository.EXPECT().GetByID(ctx, uint64(1)).Return(client, nil)
			service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(1)).Return(training, nil)

		},
		CheckOutput: func(t *testing.T, err error) {
			assert.Error(t, err)
			assert.Equal(t, servicesErrors.NoAvailablePlacesNum, err)
		},
	},
}

func TestClientServiceImplementation(t *testing.T) {
	for _, tt := range testGetByTelephone {
		t.Run(tt.TestName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := createMockClientService(ctrl)
			tt.Prepare(service)

			clientService := createClientService(service)

			client, err := clientService.GetByTelephone(tt.InputData)
			tt.CheckOutput(t, client, err)
		})
	}
	t.Run("GetByTelephone", func(t *testing.T) {
		for _, tt := range testGetByTelephone {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockClientService(ctrl)
				tt.Prepare(service)

				clientService := createClientService(service)
				client, err := clientService.GetByTelephone(tt.InputData)

				tt.CheckOutput(t, client, err)
			})
		}
	})

	t.Run("Create", func(t *testing.T) {
		for _, tt := range testCreateClient {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockClientService(ctrl)
				tt.Prepare(service)

				clientService := createClientService(service)
				err := clientService.Create(tt.InputData)

				tt.CheckOutput(t, err)
			})
		}
	})

	t.Run("Login", func(t *testing.T) {
		for _, tt := range testLogin {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockClientService(ctrl)
				tt.Prepare(service)

				clientService := createClientService(service)
				client, err := clientService.Login(tt.Tel, tt.Password)

				tt.CheckOutput(t, client, err)
			})
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		for _, tt := range testGetClientByID {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockClientService(ctrl)
				tt.Prepare(service)

				clientService := createClientService(service)
				client, err := clientService.GetByID(tt.ID)

				tt.CheckOutput(t, client, err)
			})
		}
	})

	t.Run("CreateAssignment", func(t *testing.T) {
		for _, tt := range testCreateAssignment {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockClientService(ctrl)
				tt.Prepare(service)

				clientService := createClientService(service)
				err := clientService.СreateAssignment(tt.ClientID, tt.TrainingID)

				tt.CheckOutput(t, err)
			})
		}
	})
}