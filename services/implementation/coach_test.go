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

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type CoachSuite struct {
	suite.Suite
}

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

func (s *CoachSuite) TestGetCoachByNameSuccess(t provider.T) {
	t.Title("GetCoachByName: Success")
	t.Tags("Coach")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockCoachService(ctrl)
		service.mockCoachRepository.EXPECT().GetByName(ctx, "Name").Return(data_builders.NewCoachBuilder().Build(), nil)

		coachService := createCoachService(service)
		coach, err := coachService.GetByName("Name")

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(coach)
		sCtx.Assert().Equal(data_builders.NewCoachBuilder().Build(), coach)
	})
}

func (s *CoachSuite) TestGetCoachByNameFailure(t provider.T) {
	t.Title("GetCoachByName: Failure")
	t.Tags("Coach")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockCoachService(ctrl)
		service.mockCoachRepository.EXPECT().GetByName(ctx, "Fail").Return(nil, errors.New("not found"))

		coachService := createCoachService(service)
		coach, err := coachService.GetByName("Fail")

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(coach)
	})
}

func (s *CoachSuite) TestCreateCoachSuccess(t provider.T) {
	t.Title("CreateCoach: Success")
	t.Tags("Coach")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockCoachService(ctrl)
		service.mockCoachRepository.EXPECT().GetByName(ctx, "Name").Return(nil, repositoriesErrors.EntityDoesNotExists)
		service.mockCoachRepository.EXPECT().Create(ctx, data_builders.NewCoachBuilder().Build()).Return(nil)

		coachService := createCoachService(service)
		err := coachService.Create(data_builders.NewCoachBuilder().Build())

		sCtx.Assert().NoError(err)
	})
}

func (s *CoachSuite) TestCreateCoachFailure(t provider.T) {
	t.Title("CreateCoach: Failure")
	t.Tags("Coach")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockCoachService(ctrl)
		service.mockCoachRepository.EXPECT().GetByName(ctx, "").Return(nil, repositoriesErrors.EntityDoesNotExists)
		service.mockCoachRepository.EXPECT().Create(ctx, data_builders.NewCoachBuilder().WithName("").Build()).Return(errors.New("validation error"))

		coachService := createCoachService(service)
		err := coachService.Create(data_builders.NewCoachBuilder().WithName("").Build())

		sCtx.Assert().Error(err)
	})
}

func (s *CoachSuite) TestGetCoachByIDSuccess(t provider.T) {
	t.Title("GetCoachByID: Success")
	t.Tags("Coach")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockCoachService(ctrl)
		service.mockCoachRepository.EXPECT().GetByID(ctx, uint64(1)).Return(data_builders.NewCoachBuilder().Build(), nil)

		coachService := createCoachService(service)
		coach, err := coachService.GetByID(uint64(1))

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(coach)
		sCtx.Assert().Equal(data_builders.NewCoachBuilder().Build(), coach)
	})
}

func (s *CoachSuite) TestGetCoachByIDFailure(t provider.T) {
	t.Title("GetCoachByID: Failure")
	t.Tags("Coach")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockCoachService(ctrl)
		service.mockCoachRepository.EXPECT().GetByID(ctx, uint64(999)).Return(nil, errors.New("not found"))

		coachService := createCoachService(service)
		coach, err := coachService.GetByID(uint64(999))

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(coach)
	})
}

func (s *CoachSuite) TestGetFreeTimeOnDateSuccess(t provider.T) {
	t.Title("GetFreeTimeOnDate: Success")
	t.Tags("Coach")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockCoachService(ctrl)
		service.mockTrainingRepository.EXPECT().GetAllByCoachOnDate(ctx, uint64(7), time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).
				Return([]models.Training{*data_builders.NewTrainingBuilder().Build()}, nil)

		coachService := createCoachService(service)
		slots, err := coachService.GetFreeTimeOnDate(uint64(7), time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC))

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(slots)
		sCtx.Assert().Equal([]time.Time{
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
	})
}

func (s *CoachSuite) TestGetFreeTimeOnDateFailure(t provider.T) {
	t.Title("GetFreeTimeOnDate: Failure")
	t.Tags("Coach")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockCoachService(ctrl)
		service.mockTrainingRepository.EXPECT().GetAllByCoachOnDate(ctx, uint64(1), time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).
				Return(nil, errors.New("no slots found"))

		coachService := createCoachService(service)
		slots, err := coachService.GetFreeTimeOnDate(uint64(1), time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC))

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(slots)
	})
}

func TestCoachSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(CoachSuite))
}
