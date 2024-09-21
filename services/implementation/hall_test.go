package servicesImplementation

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"
	repositories_mocks "github.com/nkarakotova/lim-core/repositories/mocks"
	"github.com/nkarakotova/lim-core/services"
	data_builders "github.com/nkarakotova/lim-core/services/implementation/data_builders"

	"github.com/charmbracelet/log"
	"github.com/golang/mock/gomock"
	"github.com/nkarakotova/lim-core/models"
	"github.com/stretchr/testify/assert"
)

type mockHallService struct {
	mockHallRepository     *repositories_mocks.MockHallRepository
	mockTrainingRepository *repositories_mocks.MockTrainingRepository
	logger                 *log.Logger
}

func createMockHallService(controller *gomock.Controller) *mockHallService {
	service := new(mockHallService)

	service.mockHallRepository = repositories_mocks.NewMockHallRepository(controller)
	service.mockTrainingRepository = repositories_mocks.NewMockTrainingRepository(controller)
	service.logger = log.New(os.Stderr)

	return service
}

func createHallService(service *mockHallService) services.HallService {
	return NewHallServiceImplementation(service.mockHallRepository, service.mockTrainingRepository, service.logger)
}

var testGetHallByNumber = []struct {
	TestName    string
	InputData   uint64
	Prepare     func(service *mockHallService)
	CheckOutput func(t *testing.T, hall *models.Hall, err error)
}{
	{
		TestName:  "success get hall by name",
		InputData: 1,
		Prepare: func(service *mockHallService) {
			ctx := context.Background()
			service.mockHallRepository.EXPECT().GetByNumber(ctx, uint64(1)).
				Return(data_builders.NewHallBuilder().Build(), nil)
		},
		CheckOutput: func(t *testing.T, hall *models.Hall, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, hall)
			assert.Equal(t, uint64(1), hall.Number)
		},
	},
	{
		TestName:  "error getting hall by name",
		InputData: 2,
		Prepare: func(service *mockHallService) {
			ctx := context.Background()
			service.mockHallRepository.EXPECT().GetByNumber(ctx, uint64(2)).Return(nil, errors.New("not found"))
		},
		CheckOutput: func(t *testing.T, hall *models.Hall, err error) {
			assert.Error(t, err)
			assert.Nil(t, hall)
		},
	},
}

var testCreateHall = []struct {
	TestName    string
	InputData   *models.Hall
	Prepare     func(service *mockHallService)
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName:  "success create hall",
		InputData: data_builders.NewHallBuilder().Build(),
		Prepare: func(service *mockHallService) {
			ctx := context.Background()
			service.mockHallRepository.EXPECT().GetByNumber(ctx, uint64(1)).Return(nil, repositoriesErrors.EntityDoesNotExists)
			service.mockHallRepository.EXPECT().Create(ctx, data_builders.NewHallBuilder().Build()).Return(nil)
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.NoError(t, err)
		},
	},
	{
		TestName:  "error creating hall",
		InputData: data_builders.NewHallBuilder().WithNumber(0).Build(),
		Prepare: func(service *mockHallService) {
			ctx := context.Background()
			service.mockHallRepository.EXPECT().GetByNumber(ctx, uint64(0)).Return(nil, repositoriesErrors.EntityDoesNotExists)
			service.mockHallRepository.EXPECT().Create(ctx, data_builders.NewHallBuilder().WithNumber(0).Build()).
				Return(errors.New("validation error"))
		},
		CheckOutput: func(t *testing.T, err error) {
			assert.Error(t, err)
		},
	},
}

var testGetHallByID = []struct {
	TestName    string
	InputData   uint64
	Prepare     func(service *mockHallService)
	CheckOutput func(t *testing.T, hall *models.Hall, err error)
}{
	{
		TestName:  "success get hall by ID",
		InputData: 1,
		Prepare: func(service *mockHallService) {
			ctx := context.Background()
			service.mockHallRepository.EXPECT().GetByID(ctx, uint64(1)).
				Return(data_builders.NewHallBuilder().Build(), nil)
		},
		CheckOutput: func(t *testing.T, hall *models.Hall, err error) {
			assert.NoError(t, err)
			assert.NotNil(t, hall)
			assert.Equal(t, uint64(1), hall.ID)
		},
	},
	{
		TestName:  "error getting hall by ID",
		InputData: 999,
		Prepare: func(service *mockHallService) {
			ctx := context.Background()
			service.mockHallRepository.EXPECT().GetByID(ctx, uint64(999)).Return(nil, errors.New("not found"))
		},
		CheckOutput: func(t *testing.T, hall *models.Hall, err error) {
			assert.Error(t, err)
			assert.Nil(t, hall)
		},
	},
}

var testGetFreeOnDateTime = []struct {
	TestName  string
	InputData struct {
		slot time.Time
	}
	Prepare     func(service *mockHallService)
	CheckOutput func(t *testing.T, freeHalls map[uint64]models.Hall, err error)
}{
	{
		TestName: "success get free time on day",
		InputData: struct {
			slot time.Time
		}{slot: time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)},

		Prepare: func(service *mockHallService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetAllByDateTime(ctx, time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)).Return(
				[]models.Training{*data_builders.NewTrainingBuilder().Build()}, nil)

			service.mockHallRepository.EXPECT().GetAll(ctx).Return(
				map[uint64]models.Hall{1: *data_builders.NewHallBuilder().Build(), 2: *data_builders.NewHallBuilder().WithID(2).Build()}, nil)
		},
		CheckOutput: func(t *testing.T, freeHalls map[uint64]models.Hall, err error) {
			assert.NoError(t, err)
			assert.Equal(t, map[uint64]models.Hall{2: *data_builders.NewHallBuilder().WithID(2).Build()}, freeHalls)
		},
	},
}

func TestHallServiceImplementation(t *testing.T) {
	t.Run("GetHallByNumber", func(t *testing.T) {
		for _, tt := range testGetHallByNumber {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockHallService(ctrl)
				tt.Prepare(service)

				hallService := createHallService(service)
				hall, err := hallService.GetByNumber(tt.InputData)

				tt.CheckOutput(t, hall, err)
			})
		}
	})

	t.Run("CreateHall", func(t *testing.T) {
		for _, tt := range testCreateHall {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockHallService(ctrl)
				tt.Prepare(service)

				hallService := createHallService(service)
				err := hallService.Create(tt.InputData)

				tt.CheckOutput(t, err)
			})
		}
	})

	t.Run("GetHallByID", func(t *testing.T) {
		for _, tt := range testGetHallByID {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockHallService(ctrl)
				tt.Prepare(service)

				hallService := createHallService(service)
				hall, err := hallService.GetByID(tt.InputData)

				tt.CheckOutput(t, hall, err)
			})
		}
	})

	t.Run("GetFreeOnDateTime", func(t *testing.T) {
		for _, tt := range testGetFreeOnDateTime {
			tt := tt
			t.Run(tt.TestName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				service := createMockHallService(ctrl)
				tt.Prepare(service)

				hallService := createHallService(service)

				slots, err := hallService.GetFreeOnDateTime(tt.InputData.slot)

				tt.CheckOutput(t, slots, err)
			})
		}
	})
}
