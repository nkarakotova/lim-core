package servicesImplementation

import (
	"context"
	"os"
	"testing"
	"time"

	repositories_mocks "github.com/nkarakotova/lim-core/repositories/mocks"
	"github.com/nkarakotova/lim-core/services"

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

//-------------------------------------------------------------------------------------------------
// get free on date time

var testGetFreeOnDateTimeSuccess = []struct {
	TestName  string
	InputData struct {
		slot time.Time
	}
	Prepare     func(service *mockHallService)
	CheckOutput func(t *testing.T, freeHalls map[uint64]models.Hall, err error)
}{
	{
		TestName: "simple get free time on day",
		InputData: struct {
			slot time.Time
		}{slot: time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC)},

		Prepare: func(service *mockHallService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetAllByDateTime(ctx, time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC)).Return(
				[]models.Training{
					{
						ID:       830,
						DateTime: time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC),
						HallID:   111,
					},
					{
						ID:       828,
						DateTime: time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC),
						HallID:   113,
					},
				}, nil)

			service.mockHallRepository.EXPECT().GetAll(ctx).Return(
				map[uint64]models.Hall{
					123: {
						ID:       123,
						Number:   1,
						Capacity: 20,
					},
					111: {
						ID:       111,
						Number:   2,
						Capacity: 20,
					},
				}, nil)
		},
		CheckOutput: func(t *testing.T, freeHalls map[uint64]models.Hall, err error) {
			assert.NoError(t, err)
			assert.Equal(t, map[uint64]models.Hall{
				123: {
					ID:       123,
					Number:   1,
					Capacity: 20,
				},
			}, freeHalls)
		},
	},
}

func TestHallGetFreeOnDateTime(t *testing.T) {
	for _, tt := range testGetFreeOnDateTimeSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := createMockHallService(ctrl)
			tt.Prepare(service)

			coachService := createHallService(service)

			slots, err := coachService.GetFreeOnDateTime(tt.InputData.slot)

			tt.CheckOutput(t, slots, err)
		})
	}
}

//-------------------------------------------------------------------------------------------------
