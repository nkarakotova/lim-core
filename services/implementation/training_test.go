package servicesImplementation

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"
	"github.com/nkarakotova/lim-core/errors/servicesErrors"
	managers_mocks "github.com/nkarakotova/lim-core/managers/mocks"
	repositories_mocks "github.com/nkarakotova/lim-core/repositories/mocks"
	"github.com/nkarakotova/lim-core/services"
	data_builders "github.com/nkarakotova/lim-core/services/implementation/data_builders"

	"github.com/nkarakotova/lim-core/models"

	"github.com/charmbracelet/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockTrainingService struct {
	mockTrainingRepository     *repositories_mocks.MockTrainingRepository
	mockClientRepository       *repositories_mocks.MockClientRepository
	mockCoachRepository        *repositories_mocks.MockCoachRepository
	mockHallRepository         *repositories_mocks.MockHallRepository
	mockTransactionManager     *managers_mocks.MockTransactionManager
	logger                     *log.Logger
}

func createMockTrainingService(controller *gomock.Controller) *mockTrainingService {
	service := new(mockTrainingService)

	service.mockTrainingRepository = repositories_mocks.NewMockTrainingRepository(controller)
	service.mockClientRepository = repositories_mocks.NewMockClientRepository(controller)
	service.mockCoachRepository = repositories_mocks.NewMockCoachRepository(controller)
	service.mockHallRepository = repositories_mocks.NewMockHallRepository(controller)
	service.mockTransactionManager = managers_mocks.NewMockTransactionManager(controller)
	service.logger = log.New(os.Stderr)

	return service
}

func createTrainingService(service *mockTrainingService) services.TrainingService {
	return NewTrainingServiceImplementation(service.mockTrainingRepository, service.mockClientRepository,
		service.mockCoachRepository, service.mockHallRepository,
		service.mockTransactionManager, service.logger)
}

var testTrainingCreate = []struct {
	TestName    string
	InputData   *models.Training
	Prepare     func(service *mockTrainingService)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName:  "success create",
		InputData: data_builders.NewTrainingBuilder().Build(),

		Prepare: func(service *mockTrainingService) {
			ctx := context.Background()
			service.mockHallRepository.EXPECT().GetByID(ctx, uint64(1)).Return(data_builders.NewHallBuilder().Build(), nil)
			service.mockTrainingRepository.EXPECT().GetAllByDateTime(ctx, time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)).Return(nil, nil)
			service.mockTrainingRepository.EXPECT().Create(ctx, data_builders.NewTrainingBuilder().Build()).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
	{
		TestName:  "error create",
		InputData: data_builders.NewTrainingBuilder().WithDateTime(time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).Build(),

		Prepare: func(service *mockTrainingService) {
		},

		CheckOutput: func(t *testing.T, err error) {
			assert.ErrorIs(t, err, servicesErrors.IncorrectTrainingTime)
		},
	},
}

var testTrainingDelete = []struct {
	TestName    string
	InputData   uint64
	Prepare     func(service *mockTrainingService)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName:  "success delete",
		InputData: 1,

		Prepare: func(service *mockTrainingService) {
			ctx := context.Background()

			service.mockTrainingRepository.EXPECT().Delete(ctx, uint64(1)).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
	{
		TestName:  "error delete",
		InputData: 1,

		Prepare: func(service *mockTrainingService) {
			ctx := context.Background()

			service.mockTrainingRepository.EXPECT().Delete(ctx, uint64(1)).Return(errors.New("not found"))
		},

		CheckOutput: func(t *testing.T, err error) {
			assert.Error(t, err)
		},
	},
}

var testGetTrainingByID = []struct {
	TestName    string
	InputData   uint64
	Prepare     func(service *mockTrainingService)
	CheckOutput func(t *testing.T, training *models.Training, err error)
}{
	{
		TestName:  "success get training by ID",
		InputData: 1,
		Prepare: func(service *mockTrainingService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(1)).
				Return(data_builders.NewTrainingBuilder().Build(), nil)
		},
		CheckOutput: func(t *testing.T, training *models.Training, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, training)
			assert.Equal(t, uint64(1), training.ID)
		},
	},
	{
		TestName:  "error getting training by ID",
		InputData: 999,
		Prepare: func(service *mockTrainingService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(999)).Return(nil, errors.New("not found"))
		},
		CheckOutput: func(t *testing.T, training *models.Training, err error) {
			assert.Error(t, err)
			assert.Nil(t, training)
		},
	},
}

var testGetAllByClient = []struct {
	TestName    string
	InputData   uint64
	Prepare     func(service *mockTrainingService)
	CheckOutput func(t *testing.T, trainings []models.Training, err error)
}{
	{
		TestName:  "success get all by client",
		InputData: 1,
		Prepare: func(service *mockTrainingService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetAllByClient(ctx, uint64(1)).
				Return([]models.Training{*data_builders.NewTrainingBuilder().Build()}, nil)
		},
		CheckOutput: func(t *testing.T, trainings []models.Training, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, trainings)
			assert.Equal(t, []models.Training{*data_builders.NewTrainingBuilder().Build()}, trainings)
		},
	},
	{
		TestName:  "error get all by client",
		InputData: 999,
		Prepare: func(service *mockTrainingService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetAllByClient(ctx, uint64(999)).
				Return(nil, errors.New("not found"))
		},
		CheckOutput: func(t *testing.T, trainings []models.Training, err error) {
			assert.Error(t, err)
			assert.Nil(t, trainings)
		},
	},
}

var testGetAllByCoachOnDate = []struct {
	TestName    string
	CoachID     uint64
	Date        time.Time
	Prepare     func(service *mockTrainingService)
	CheckOutput func(t *testing.T, trainings []models.Training, err error)
}{
	{
		TestName: "success get all by coach on date",
		CoachID:  1,
		Date:     time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC),
		Prepare: func(service *mockTrainingService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetAllByCoachOnDate(ctx, uint64(1), time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).
				Return([]models.Training{*data_builders.NewTrainingBuilder().Build()}, nil)
		},
		CheckOutput: func(t *testing.T, trainings []models.Training, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, trainings)
			assert.Equal(t, []models.Training{*data_builders.NewTrainingBuilder().Build()}, trainings)
		},
	},
	{
		TestName: "error get all by coach on date",
		CoachID:  2,
		Date:     time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC),
		Prepare: func(service *mockTrainingService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetAllByCoachOnDate(ctx, uint64(2), time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).
				Return(nil, errors.New("not found"))
		},
		CheckOutput: func(t *testing.T, trainings []models.Training, err error) {
			assert.Error(t, err)
			assert.Nil(t, trainings)
		},
	},
}

var testGetAllByDateTime = []struct {
	TestName    string
	DateTime    time.Time
	Prepare     func(service *mockTrainingService)
	CheckOutput func(t *testing.T, trainings []models.Training, err error)
}{
	{
		TestName: "success get all by date time",
		DateTime: time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC),
		Prepare: func(service *mockTrainingService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetAllByDateTime(ctx, time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)).
				Return([]models.Training{*data_builders.NewTrainingBuilder().Build()}, nil)
		},
		CheckOutput: func(t *testing.T, trainings []models.Training, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, trainings)
			assert.Equal(t, []models.Training{*data_builders.NewTrainingBuilder().Build()}, trainings)
		},
	},
	{
		TestName: "error get all by date time",
		DateTime: time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC),
		Prepare: func(service *mockTrainingService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetAllByDateTime(ctx, time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).
				Return(nil, repositoriesErrors.EntityDoesNotExists)
		},
		CheckOutput: func(t *testing.T, trainings []models.Training, err error) {
			assert.Error(t, err)
			assert.Nil(t, trainings)
		},
	},
}

var testGetAllBetweenDateTime = []struct {
	TestName      string
	StartDateTime time.Time
	EndDateTime   time.Time
	Prepare       func(service *mockTrainingService)
	CheckOutput   func(t *testing.T, trainings []models.Training, err error)
}{
	{
		TestName:      "success get all between date time",
		StartDateTime: time.Date(2024, 7, 5, 12, 0, 0, 0, time.UTC),
		EndDateTime:   time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC),
		Prepare: func(service *mockTrainingService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetAllBetweenDateTime(ctx, time.Date(2024, 7, 5, 12, 0, 0, 0, time.UTC),
				time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)).
				Return([]models.Training{*data_builders.NewTrainingBuilder().Build()}, nil)
		},
		CheckOutput: func(t *testing.T, trainings []models.Training, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, trainings)
			assert.Equal(t, []models.Training{*data_builders.NewTrainingBuilder().Build()}, trainings)
		},
	},
	{
		TestName:      "error get all between date time",
		StartDateTime: time.Date(2024, 7, 13, 0, 0, 0, 0, time.UTC),
		EndDateTime:   time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC),
		Prepare: func(service *mockTrainingService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetAllBetweenDateTime(ctx, time.Date(2024, 7, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC)).
				Return(nil, errors.New("not found"))
		},
		CheckOutput: func(t *testing.T, trainings []models.Training, err error) {
			assert.Error(t, err)
			assert.Nil(t, trainings)
		},
	},
}

func TestTrainingServiceImplementationCreate(t *testing.T) {
	t.Run("CreateTraining", func(t *testing.T) {
		for _, tt := range testTrainingCreate {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockTrainingService(ctrl)
				tt.Prepare(service)

				trainingService := createTrainingService(service)

				err := trainingService.Create(tt.InputData)

				tt.CheckOutput(t, err)
			})
		}
	})
	// t.Run("DeleteTraining", func(t *testing.T) {
	// 	for _, tt := range testTrainingDelete {
	// 		tt := tt
	// 		t.Run(tt.TestName, func(t *testing.T) {
	// 			ctrl := gomock.NewController(t)
	// 			defer ctrl.Finish()

	// 			service := createMockTrainingService(ctrl)
	// 			tt.Prepare(service)

	// 			trainingService := createTrainingService(service)

	// 			err := trainingService.Delete(tt.InputData)

	// 			tt.CheckOutput(t, err)
	// 		})
	// 	}
	// })
	// t.Run("GetTrainingByID", func(t *testing.T) {
	// 	for _, tt := range testGetTrainingByID {
	// 		tt := tt
	// 		t.Run(tt.TestName, func(t *testing.T) {
	// 			ctrl := gomock.NewController(t)
	// 			defer ctrl.Finish()

	// 			service := createMockTrainingService(ctrl)
	// 			tt.Prepare(service)

	// 			trainingService := createTrainingService(service)
	// 			training, err := trainingService.GetByID(tt.InputData)

	// 			tt.CheckOutput(t, training, err)
	// 		})
	// 	}
	// })

	// t.Run("GetAllByClient", func(t *testing.T) {
	// 	for _, tt := range testGetAllByClient {
	// 		tt := tt
	// 		t.Run(tt.TestName, func(t *testing.T) {
	// 			ctrl := gomock.NewController(t)
	// 			defer ctrl.Finish()

	// 			service := createMockTrainingService(ctrl)
	// 			tt.Prepare(service)

	// 			trainingService := createTrainingService(service)
	// 			trainings, err := trainingService.GetAllByClient(tt.InputData)

	// 			tt.CheckOutput(t, trainings, err)
	// 		})
	// 	}
	// })

	// t.Run("GetAllByCoachOnDate", func(t *testing.T) {
	// 	for _, tt := range testGetAllByCoachOnDate {
	// 		tt := tt
	// 		t.Run(tt.TestName, func(t *testing.T) {
	// 			ctrl := gomock.NewController(t)
	// 			defer ctrl.Finish()

	// 			service := createMockTrainingService(ctrl)
	// 			tt.Prepare(service)

	// 			trainingService := createTrainingService(service)
	// 			trainings, err := trainingService.GetAllByCoachOnDate(tt.CoachID, tt.Date)

	// 			tt.CheckOutput(t, trainings, err)
	// 		})
	// 	}
	// })

	// t.Run("GetAllByDateTime", func(t *testing.T) {
	// 	for _, tt := range testGetAllByDateTime {
	// 		tt := tt
	// 		t.Run(tt.TestName, func(t *testing.T) {
	// 			ctrl := gomock.NewController(t)
	// 			defer ctrl.Finish()

	// 			service := createMockTrainingService(ctrl)
	// 			tt.Prepare(service)

	// 			trainingService := createTrainingService(service)
	// 			trainings, err := trainingService.GetAllByDateTime(tt.DateTime)

	// 			tt.CheckOutput(t, trainings, err)
	// 		})
	// 	}
	// })

	// t.Run("GetAllBetweenDateTime", func(t *testing.T) {
	// 	for _, tt := range testGetAllBetweenDateTime {
	// 		tt := tt
	// 		t.Run(tt.TestName, func(t *testing.T) {
	// 			ctrl := gomock.NewController(t)
	// 			defer ctrl.Finish()

	// 			service := createMockTrainingService(ctrl)
	// 			tt.Prepare(service)

	// 			trainingService := createTrainingService(service)
	// 			trainings, err := trainingService.GetAllBetweenDateTime(tt.StartDateTime, tt.EndDateTime)

	// 			tt.CheckOutput(t, trainings, err)
	// 		})
	// 	}
	// })
}
