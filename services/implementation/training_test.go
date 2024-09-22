package servicesImplementation

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/nkarakotova/lim-core/errors/repositoriesErrors"
	// "github.com/nkarakotova/lim-core/errors/servicesErrors"
	managers_mocks "github.com/nkarakotova/lim-core/managers/mocks"
	repositories_mocks "github.com/nkarakotova/lim-core/repositories/mocks"
	"github.com/nkarakotova/lim-core/services"
	data_builders "github.com/nkarakotova/lim-core/services/implementation/data_builders"

	"github.com/nkarakotova/lim-core/models"

	"github.com/charmbracelet/log"
	"github.com/golang/mock/gomock"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
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

type TrainingSuite struct {
	suite.Suite
}

func (s *TrainingSuite) TestCreateTrainingSuccess(t provider.T) {
	t.Title("CreateTraining: Success")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)
		service.mockTrainingRepository.EXPECT().GetAllByDateTime(ctx, time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)).Return(nil, nil)
		service.mockTrainingRepository.EXPECT().Create(ctx, data_builders.NewTrainingBuilder().Build()).Return(nil)

		trainingService := createTrainingService(service)
		err := trainingService.Create(data_builders.NewTrainingBuilder().Build())

		sCtx.Assert().NoError(err)
	})
}

func (s *TrainingSuite) TestCreateTrainingFailure(t provider.T) {
	t.Title("CreateTraining: Failure")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)

		trainingService := createTrainingService(service)
		err := trainingService.Create(data_builders.NewTrainingBuilder().WithDateTime(time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).Build())

		sCtx.Assert().Error(err)
	})
}

func (s *TrainingSuite) TestDeleteTrainingSuccess(t provider.T) {
	t.Title("DeleteTraining: Success")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)
		service.mockTrainingRepository.EXPECT().Delete(ctx, uint64(1)).Return(nil)

		trainingService := createTrainingService(service)
		err := trainingService.Delete(uint64(1))

		sCtx.Assert().NoError(err)
	})
}

func (s *TrainingSuite) TestDeleteTrainingFailure(t provider.T) {
	t.Title("DeleteTraining: Failure")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)
		service.mockTrainingRepository.EXPECT().Delete(ctx, uint64(1)).Return(errors.New("not found"))

		trainingService := createTrainingService(service)
		err := trainingService.Delete(uint64(1))

		sCtx.Assert().Error(err)
	})
}

func (s *TrainingSuite) TestGetTrainingByIDSuccess(t provider.T) {
	t.Title("GetTrainingByID: Success")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)
		service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(1)).Return(data_builders.NewTrainingBuilder().Build(), nil)

		trainingService := createTrainingService(service)
		training, err := trainingService.GetByID(uint64(1))

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(training)
		sCtx.Assert().Equal(data_builders.NewTrainingBuilder().Build(), training)
	})
}

func (s *TrainingSuite) TestGetTrainingByIDFailure(t provider.T) {
	t.Title("GetTrainingByID: Failure")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)
		service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(999)).Return(nil, errors.New("not found"))

		trainingService := createTrainingService(service)
		training, err := trainingService.GetByID(uint64(999))

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(training)
	})
}

func (s *TrainingSuite) TestGetAllByClientSuccess(t provider.T) {
	t.Title("GetAllByClient: Success")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)
		service.mockTrainingRepository.EXPECT().GetAllByClient(ctx, uint64(1)).
			Return([]models.Training{*data_builders.NewTrainingBuilder().Build()}, nil)

		trainingService := createTrainingService(service)
		training, err := trainingService.GetAllByClient(uint64(1))

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(training)
		sCtx.Assert().Equal([]models.Training{*data_builders.NewTrainingBuilder().Build()}, training)
	})
}

func (s *TrainingSuite) TestGetAllByClientFailure(t provider.T) {
	t.Title("GetAllByClient: Failure")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)
		service.mockTrainingRepository.EXPECT().GetAllByClient(ctx, uint64(999)).
			Return(nil, errors.New("not found"))

		trainingService := createTrainingService(service)
		training, err := trainingService.GetAllByClient(uint64(999))

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(training)
	})
}

func (s *TrainingSuite) TestGetAllByCoachOnDateSuccess(t provider.T) {
	t.Title("GetAllByCoachOnDate: Success")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)
		service.mockTrainingRepository.EXPECT().GetAllByCoachOnDate(ctx, uint64(1), time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).
				Return([]models.Training{*data_builders.NewTrainingBuilder().Build()}, nil)

		trainingService := createTrainingService(service)
		training, err := trainingService.GetAllByCoachOnDate(uint64(1), time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC))

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(training)
		sCtx.Assert().Equal([]models.Training{*data_builders.NewTrainingBuilder().Build()}, training)
	})
}

func (s *TrainingSuite) TestGetAllByCoachOnDateFailure(t provider.T) {
	t.Title("GetAllByCoachOnDate: Failure")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)
		service.mockTrainingRepository.EXPECT().GetAllByCoachOnDate(ctx, uint64(2), time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).
				Return(nil, errors.New("not found"))

		trainingService := createTrainingService(service)
		training, err := trainingService.GetAllByCoachOnDate(uint64(2), time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC))

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(training)
	})
}

func (s *TrainingSuite) TestGetAllByDateTimeSuccess(t provider.T) {
	t.Title("GetAllByDateTime: Success")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)
		service.mockTrainingRepository.EXPECT().GetAllByDateTime(ctx, time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)).
				Return([]models.Training{*data_builders.NewTrainingBuilder().Build()}, nil)

		trainingService := createTrainingService(service)
		training, err := trainingService.GetAllByDateTime(time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC))

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(training)
		sCtx.Assert().Equal([]models.Training{*data_builders.NewTrainingBuilder().Build()}, training)
	})
}

func (s *TrainingSuite) TestGetAllByDateTimeFailure(t provider.T) {
	t.Title("GetAllByDateTime: Failure")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)
		service.mockTrainingRepository.EXPECT().GetAllByDateTime(ctx, time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC)).
				Return(nil, repositoriesErrors.EntityDoesNotExists)

		trainingService := createTrainingService(service)
		training, err := trainingService.GetAllByDateTime(time.Date(2024, 7, 7, 0, 0, 0, 0, time.UTC))

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(training)
	})
}

func (s *TrainingSuite) TestGetAllBetweenDateTimeSuccess(t provider.T) {
	t.Title("GetAllBetweenDateTime: Success")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)
		service.mockTrainingRepository.EXPECT().GetAllBetweenDateTime(ctx, time.Date(2024, 7, 5, 12, 0, 0, 0, time.UTC), time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC)).
				Return([]models.Training{*data_builders.NewTrainingBuilder().Build()}, nil)

		trainingService := createTrainingService(service)
		training, err := trainingService.GetAllBetweenDateTime(time.Date(2024, 7, 5, 12, 0, 0, 0, time.UTC), time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC))

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(training)
		sCtx.Assert().Equal([]models.Training{*data_builders.NewTrainingBuilder().Build()}, training)
	})
}

func (s *TrainingSuite) TestGetAllBetweenDateTimeFailure(t provider.T) {
	t.Title("GetAllBetweenDateTime: Failure")
	t.Tags("Training")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockTrainingService(ctrl)
		service.mockTrainingRepository.EXPECT().GetAllBetweenDateTime(ctx, time.Date(2024, 7, 13, 0, 0, 0, 0, time.UTC), time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC)).
				Return(nil, repositoriesErrors.EntityDoesNotExists)

		trainingService := createTrainingService(service)
		training, err := trainingService.GetAllBetweenDateTime(time.Date(2024, 7, 13, 0, 0, 0, 0, time.UTC), time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC))

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(training)
	})
}

func TestTrainingSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(TrainingSuite))
}
