package servicesImplementation

import (
	"context"
	"os"
	"testing"
	"time"
	"errors"

	repositories_mocks "github.com/nkarakotova/lim-core/repositories/mocks"
	data_builders "github.com/nkarakotova/lim-core/services/implementation/data_builders"
	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"
	"github.com/nkarakotova/lim-core/services"
	"github.com/nkarakotova/lim-core/models"

	"github.com/charmbracelet/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockCoachService struct {
	mockCoachRepository    *repositories_mocks.MockCoachRepository
	mockTrainingRepository *repositories_mocks.MockTrainingRepository
	logger                 *log.Logger
}

func createMockCoachService(controller *gomock.Controller) *mockCoachService {
	service := new(mockCoachService)

	service.mockCoachRepository = repositories_mocks.NewMockCoachRepository(controller)
	service.mockTrainingRepository = repositories_mocks.NewMockTrainingRepository(controller)
	service.logger = log.New(os.Stderr)

	return service
}

func createCoachService(service *mockCoachService) services.CoachService {
	return NewCoachServiceImplementation(service.mockCoachRepository, service.mockTrainingRepository, service.logger)
}

var testGetCoachByName = []struct {
	TestName  string
	InputData string
	Prepare   func(service *mockCoachService)
	CheckOutput func(t *testing.T, coach *models.Coach, err error)
}{
	{
		TestName:  "success get coach by name",
		InputData: "Name",
		Prepare: func(service *mockCoachService) {
			ctx := context.Background()
			service.mockCoachRepository.EXPECT().GetByName(ctx, "Name").
			Return(data_builders.NewCoachBuilder().Build(), nil)
		},
		CheckOutput: func(t *testing.T, coach *models.Coach, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, coach)
			assert.Equal(t, "Name", coach.Name)
		},
	},
	{
		TestName:  "error getting coach by name",
		InputData: "Nonexistent",
		Prepare: func(service *mockCoachService) {
			ctx := context.Background()
			service.mockCoachRepository.EXPECT().GetByName(ctx, "Nonexistent").Return(nil, errors.New("not found"))
		},
		CheckOutput: func(t *testing.T, coach *models.Coach, err error) {
			assert.Error(t, err)
			assert.Nil(t, coach)
		},
	},
}

var testCreateCoach = []struct {
	TestName  string
	InputData *models.Coach
	Prepare   func(service *mockCoachService)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName:  "success create coach",
		InputData: data_builders.NewCoachBuilder().Build(),
		Prepare: func(service *mockCoachService) {
			ctx := context.Background()
			service.mockCoachRepository.EXPECT().GetByName(ctx, "Name").Return(nil, repositoriesErrors.EntityDoesNotExists)
			service.mockCoachRepository.EXPECT().Create(ctx, data_builders.NewCoachBuilder().Build()).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
	{
		TestName:  "error creating coach",
		InputData: data_builders.NewCoachBuilder().WithName("").Build(),
		Prepare: func(service *mockCoachService) {
			ctx := context.Background()
			service.mockCoachRepository.EXPECT().GetByName(ctx, "").Return(nil, repositoriesErrors.EntityDoesNotExists)
			service.mockCoachRepository.EXPECT().Create(ctx, data_builders.NewCoachBuilder().WithName("").Build()).
			Return(errors.New("validation error"))
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.Error(t, err)
		},
	},
}

var testGetCoachByID = []struct {
	TestName  string
	InputData uint64
	Prepare   func(service *mockCoachService)
	CheckOutput func(t *testing.T, coach *models.Coach, err error)
}{
	{
		TestName:  "success get coach by ID",
		InputData: 1,
		Prepare: func(service *mockCoachService) {
			ctx := context.Background()
			service.mockCoachRepository.EXPECT().GetByID(ctx, uint64(1)).
			Return(data_builders.NewCoachBuilder().Build(), nil)
		},
		CheckOutput: func(t *testing.T, coach *models.Coach, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, coach)
			assert.Equal(t, uint64(1), coach.ID)
		},
	},
	{
		TestName:  "error getting coach by ID",
		InputData: 999,
		Prepare: func(service *mockCoachService) {
			ctx := context.Background()
			service.mockCoachRepository.EXPECT().GetByID(ctx, uint64(999)).Return(nil, errors.New("not found"))
		},
		CheckOutput: func(t *testing.T, coach *models.Coach, err error) {
			assert.Error(t, err)
			assert.Nil(t, coach)
		},
	},
}

var testGetFreeTimeOnDate = []struct {
	TestName  string
	InputData struct {
		coachID uint64
		date    time.Time
	}
	Prepare     func(service *mockCoachService)
	CheckOutput func(t *testing.T, slots []time.Time, err error)
}{
	{
		TestName: "error get free time on day",
		InputData: struct {
			coachID uint64
			date    time.Time
		}{coachID: 1, date: time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)},

		Prepare: func(service *mockCoachService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetAllByCoachOnDate(ctx, uint64(1), time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).
				Return(nil, errors.New("no slots found"))
		},
		CheckOutput: func(t *testing.T, slots []time.Time, err error) {
			assert.Error(t, err)
			assert.Nil(t, slots)
		},
	},
	{
		TestName: "success get free time on day",
		InputData: struct {
			coachID uint64
			date    time.Time
		}{coachID: 7, date: time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)},

		Prepare: func(service *mockCoachService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetAllByCoachOnDate(ctx, uint64(7), time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).
				Return([]models.Training{*data_builders.NewTrainingBuilder().Build()}, nil)
		},
		CheckOutput: func(t *testing.T, slots []time.Time, err error) {
			assert.NoError(t, err)
			assert.Equal(t, []time.Time{
				time.Date(2024, 7, 7, 10, 0, 0, 0, time.UTC),
				time.Date(2024, 7, 7, 11, 0, 0, 0, time.UTC),
				time.Date(2024, 7, 7, 13, 0, 0, 0, time.UTC),
				time.Date(2024, 7, 7, 14, 0, 0, 0, time.UTC),
				time.Date(2024, 7, 7, 15, 0, 0, 0, time.UTC),
				time.Date(2024, 7, 7, 16, 0, 0, 0, time.UTC),
				time.Date(2024, 7, 7, 17, 0, 0, 0, time.UTC),
				time.Date(2024, 7, 7, 18, 0, 0, 0, time.UTC),
				time.Date(2024, 7, 7, 19, 0, 0, 0, time.UTC),
				time.Date(2024, 7, 7, 20, 0, 0, 0, time.UTC),
				time.Date(2024, 7, 7, 21, 0, 0, 0, time.UTC),
			}, slots)
		},
	},
}

func TestCoachServiceImplementation(t *testing.T) {
	t.Run("GetCoachByName", func(t *testing.T) {
		for _, tt := range testGetCoachByName {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockCoachService(ctrl)
				tt.Prepare(service)

				coachService := createCoachService(service)
				coach, err := coachService.GetByName(tt.InputData)

				tt.CheckOutput(t, coach, err)
			})
		}
	})

	t.Run("Create", func(t *testing.T) {
		for _, tt := range testCreateCoach {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockCoachService(ctrl)
				tt.Prepare(service)

				coachService := createCoachService(service)
				err := coachService.Create(tt.InputData)

				tt.CheckOutput(t, err)
			})
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		for _, tt := range testGetCoachByID {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockCoachService(ctrl)
				tt.Prepare(service)

				coachService := createCoachService(service)
				coach, err := coachService.GetByID(tt.InputData)

				tt.CheckOutput(t, coach, err)
			})
		}
	})

	t.Run("GetFreeTimeOnDate", func(t *testing.T) {
		for _, tt := range testGetFreeTimeOnDate {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockCoachService(ctrl)
				tt.Prepare(service)
				coachService := createCoachService(service)
				
				slots, err := coachService.GetFreeTimeOnDate(tt.InputData.coachID, tt.InputData.date)
				tt.CheckOutput(t, slots, err)
			})
		}
	})
}
