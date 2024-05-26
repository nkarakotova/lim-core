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

//-------------------------------------------------------------------------------------------------
// get free time on date

var testGetFreeTimeOnDateSuccess = []struct {
	TestName  string
	InputData struct {
		coachID uint64
		date    time.Time
	}
	Prepare     func(service *mockCoachService)
	CheckOutput func(t *testing.T, slots []time.Time, err error)
}{
	{
		TestName: "simple get free time on day",
		InputData: struct {
			coachID uint64
			date    time.Time
		}{coachID: 7, date: time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC)},

		Prepare: func(service *mockCoachService) {
			ctx := context.Background()
			service.mockTrainingRepository.EXPECT().GetAllByCoachOnDate(ctx, uint64(7), time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC)).Return(
				[]models.Training{
					{
						ID:       830,
						DateTime: time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC),
						CoachID:  7,
					},
					{
						ID:       828,
						DateTime: time.Date(2024, 3, 5, 14, 0, 0, 0, time.UTC),
						CoachID:  7,
					},
				}, nil)
		},
		CheckOutput: func(t *testing.T, slots []time.Time, err error) {
			assert.NoError(t, err)
			assert.Equal(t, []time.Time{
				time.Date(2024, 3, 5, 10, 0, 0, 0, time.UTC),
				time.Date(2024, 3, 5, 11, 0, 0, 0, time.UTC),
				time.Date(2024, 3, 5, 13, 0, 0, 0, time.UTC),
				time.Date(2024, 3, 5, 15, 0, 0, 0, time.UTC),
				time.Date(2024, 3, 5, 16, 0, 0, 0, time.UTC),
				time.Date(2024, 3, 5, 17, 0, 0, 0, time.UTC),
				time.Date(2024, 3, 5, 18, 0, 0, 0, time.UTC),
				time.Date(2024, 3, 5, 19, 0, 0, 0, time.UTC),
				time.Date(2024, 3, 5, 20, 0, 0, 0, time.UTC),
				time.Date(2024, 3, 5, 21, 0, 0, 0, time.UTC),
			}, slots)
		},
	},
}

func TestCoachServiceImplementationGetFreeTimeOnDate(t *testing.T) {
	for _, tt := range testGetFreeTimeOnDateSuccess {
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
}

//-------------------------------------------------------------------------------------------------
