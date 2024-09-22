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
	"github.com/nkarakotova/lim-core/models"

	"github.com/charmbracelet/log"
	"github.com/golang/mock/gomock"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type HallSuite struct {
	suite.Suite
}

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

func (s *HallSuite) TestGetHallByNumberSuccess(t provider.T) {
	t.Title("GetHallByNumber: Success")
	t.Tags("Hall")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockHallService(ctrl)
		service.mockHallRepository.EXPECT().GetByNumber(ctx, uint64(1)).Return(data_builders.NewHallBuilder().Build(), nil)

		hallService := createHallService(service)
		hall, err := hallService.GetByNumber(uint64(1))

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(hall)
		sCtx.Assert().Equal(data_builders.NewHallBuilder().Build(), hall)
	})
}

func (s *HallSuite) TestGetHallByNumberFailure(t provider.T) {
	t.Title("GetHallByNumber: Failure")
	t.Tags("Hall")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockHallService(ctrl)
		service.mockHallRepository.EXPECT().GetByNumber(ctx, uint64(1)).Return(nil, errors.New("not found"))

		hallService := createHallService(service)
		hall, err := hallService.GetByNumber(uint64(1))

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(hall)
	})
}

func (s *HallSuite) TestCreateHallSuccess(t provider.T) {
	t.Title("CreateHall: Success")
	t.Tags("Hall")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockHallService(ctrl)
		service.mockHallRepository.EXPECT().GetByNumber(ctx, uint64(1)).Return(nil, repositoriesErrors.EntityDoesNotExists)
		service.mockHallRepository.EXPECT().Create(ctx, data_builders.NewHallBuilder().Build()).Return(nil)

		hallService := createHallService(service)
		err := hallService.Create(data_builders.NewHallBuilder().Build())

		sCtx.Assert().NoError(err)
	})
}

func (s *HallSuite) TestCreateHallFailure(t provider.T) {
	t.Title("CreateHall: Failure")
	t.Tags("Hall")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockHallService(ctrl)
		service.mockHallRepository.EXPECT().GetByNumber(ctx, uint64(0)).Return(nil, repositoriesErrors.EntityDoesNotExists)
		service.mockHallRepository.EXPECT().Create(ctx, data_builders.NewHallBuilder().WithNumber(0).Build()).Return(errors.New("validation error"))

		hallService := createHallService(service)
		err := hallService.Create(data_builders.NewHallBuilder().WithNumber(0).Build())

		sCtx.Assert().Error(err)
	})
}

func (s *HallSuite) TestGetHallByIDSuccess(t provider.T) {
	t.Title("GetHallByID: Success")
	t.Tags("Hall")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockHallService(ctrl)
		service.mockHallRepository.EXPECT().GetByID(ctx, uint64(1)).Return(data_builders.NewHallBuilder().Build(), nil)

		hallService := createHallService(service)
		hall, err := hallService.GetByID(uint64(1))

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(hall)
		sCtx.Assert().Equal(data_builders.NewHallBuilder().Build(), hall)
	})
}

func (s *HallSuite) TestGetHallByIDFailure(t provider.T) {
	t.Title("GetHallByID: Failure")
	t.Tags("Hall")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockHallService(ctrl)
		service.mockHallRepository.EXPECT().GetByID(ctx, uint64(999)).Return(nil, errors.New("not found"))

		hallService := createHallService(service)
		hall, err := hallService.GetByID(uint64(999))

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(hall)
	})
}

func (s *HallSuite) TestGetFreeOnDateTimeSuccess(t provider.T) {
	t.Title("GetFreeOnDateTime: Success")
	t.Tags("Hall")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockHallService(ctrl)
		service.mockTrainingRepository.EXPECT().GetAllByDateTime(ctx, time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)).Return(
			[]models.Training{*data_builders.NewTrainingBuilder().Build()}, nil)
		service.mockHallRepository.EXPECT().GetAll(ctx).Return(
			map[uint64]models.Hall{1: *data_builders.NewHallBuilder().Build(), 2: *data_builders.NewHallBuilder().WithID(2).WithNumber(2).Build()}, nil)

		hallService := createHallService(service)
		hall, err := hallService.GetFreeOnDateTime(time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC))

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(hall)
		sCtx.Assert().Equal(map[uint64]models.Hall{2: *data_builders.NewHallBuilder().WithID(2).WithNumber(2).Build()}, hall)
	})
}

func (s *HallSuite) TestGetFreeOnDateTimeFailure(t provider.T) {
	t.Title("GetFreeOnDateTime: Failure")
	t.Tags("Hall")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockHallService(ctrl)
		service.mockTrainingRepository.EXPECT().GetAllByDateTime(ctx, time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)).Return(
			nil, errors.New("not found"))
		
		hallService := createHallService(service)
		hall, err := hallService.GetFreeOnDateTime(time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC))

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(hall)
	})
}

func TestHallSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(HallSuite))
}
