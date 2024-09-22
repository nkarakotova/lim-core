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
	"github.com/nkarakotova/lim-core/services"

	"github.com/charmbracelet/log"
	"github.com/golang/mock/gomock"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
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

type ClientSuite struct {
	suite.Suite
}

func (s *ClientSuite) TestGetClientByTelephoneSuccess(t provider.T) {
	t.Title("GetClientByTelephone: Success")
	t.Tags("Client")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockClientService(ctrl)
		service.mockClientRepository.EXPECT().GetByTelephone(ctx, "1234567890").Return(data_builders.NewClientBuilder().Build(), nil)

		clientService := createClientService(service)
		client, err := clientService.GetByTelephone("1234567890")

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(client)
		sCtx.Assert().Equal(data_builders.NewClientBuilder().Build(), client)
	})
}

func (s *ClientSuite) TestGetClientByTelephoneFailure(t provider.T) {
	t.Title("GetClientByTelephone: Failure")
	t.Tags("Client")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockClientService(ctrl)
		service.mockClientRepository.EXPECT().GetByTelephone(ctx, "1111111111").Return(nil, errors.New("not found"))

		clientService := createClientService(service)
		client, err := clientService.GetByTelephone("1111111111")

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(client)
	})
}

func (s *ClientSuite) TestCreateClientSuccess(t provider.T) {
	t.Title("CreateClient: Success")
	t.Tags("Client")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockClientService(ctrl)
		service.mockClientRepository.EXPECT().GetByTelephone(ctx, "1234567890").Return(nil, repositoriesErrors.EntityDoesNotExists)
		service.mockClientRepository.EXPECT().Create(ctx, data_builders.NewClientBuilder().Build()).Return(nil)

		clientService := createClientService(service)
		err := clientService.Create(data_builders.NewClientBuilder().Build())

		sCtx.Assert().NoError(err)
	})
}

func (s *ClientSuite) TestCreateClientFailure(t provider.T) {
	t.Title("CreateClient: Failure")
	t.Tags("Client")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockClientService(ctrl)
		service.mockClientRepository.EXPECT().GetByTelephone(ctx, "").Return(nil, errors.New("validation error"))

		clientService := createClientService(service)
		err := clientService.Create(data_builders.NewClientBuilder().WithTelephone("").Build())

		sCtx.Assert().Error(err)
	})
}

func (s *ClientSuite) TestLoginClientSuccess(t provider.T) {
	t.Title("LoginClient: Success")
	t.Tags("Client")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockClientService(ctrl)
		service.mockClientRepository.EXPECT().GetByTelephone(ctx, "1234567890").
				Return(data_builders.NewClientBuilder().Build(), nil)

		clientService := createClientService(service)
		client, err := clientService.Login("1234567890", "123")

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(client)
		sCtx.Assert().Equal(data_builders.NewClientBuilder().Build(), client)
	})
}

func (s *ClientSuite) TestLoginClientFailure(t provider.T) {
	t.Title("LoginClient: Failure")
	t.Tags("Client")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockClientService(ctrl)
		service.mockClientRepository.EXPECT().GetByTelephone(ctx, "1234567890").
				Return(data_builders.NewClientBuilder().Build(), nil)

		clientService := createClientService(service)
		client, err := clientService.Login("1234567890", "111")

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(client)
	})
}

func (s *ClientSuite) TestGetClientByIDSuccess(t provider.T) {
	t.Title("GetClientByID: Success")
	t.Tags("Client")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockClientService(ctrl)
		service.mockClientRepository.EXPECT().GetByID(ctx, uint64(1)).Return(data_builders.NewClientBuilder().Build(), nil)

		clientService := createClientService(service)
		client, err := clientService.GetByID(uint64(1))

		sCtx.Assert().NoError(err)
		sCtx.Assert().NotNil(client)
		sCtx.Assert().Equal(data_builders.NewClientBuilder().Build(), client)
	})
}

func (s *ClientSuite) TestGetClientByIDFailure(t provider.T) {
	t.Title("GetClientByID: Failure")
	t.Tags("Client")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockClientService(ctrl)
		service.mockClientRepository.EXPECT().GetByID(ctx, uint64(999)).Return(nil, errors.New("not found"))

		clientService := createClientService(service)
		client, err := clientService.GetByID(uint64(999))

		sCtx.Assert().Error(err)
		sCtx.Assert().Nil(client)
	})
}

func (s *ClientSuite) TestCreateAssignmentSuccess(t provider.T) {
	t.Title("CreateAssignment: Success")
	t.Tags("Client")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockClientService(ctrl)
		client := data_builders.NewClientBuilder().Build()
		training := data_builders.NewTrainingBuilder().Build()

		service.mockClientRepository.EXPECT().GetByID(ctx, uint64(1)).Return(client, nil)
		service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(1)).Return(training, nil)
		service.mockTrainingRepository.EXPECT().GetAllByClient(ctx, client.ID).Return(nil, nil)
		service.mockTransactionManager.EXPECT().WithinTransaction(ctx, gomock.Any()).Return(nil)

		clientService := createClientService(service)
		err := clientService.CreateAssignment(uint64(1), uint64(1))

		sCtx.Assert().NoError(err)
	})
}

func (s *ClientSuite) TestCreateAssignmentFailure(t provider.T) {
	t.Title("CreateAssignment: Failure")
	t.Tags("Client")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockClientService(ctrl)
		client := data_builders.NewClientBuilder().Build()
		training := data_builders.NewTrainingBuilder().WithPlacesNum(0).Build()

		service.mockClientRepository.EXPECT().GetByID(ctx, uint64(1)).Return(client, nil)
		service.mockTrainingRepository.EXPECT().GetByID(ctx, uint64(1)).Return(training, nil)

		clientService := createClientService(service)
		err := clientService.CreateAssignment(uint64(1), uint64(1))

		sCtx.Assert().Error(err)
	})
}

func (s *ClientSuite) TestDeleteAssignmentSuccess(t provider.T) {
	t.Title("DeleteAssignment: Success")
	t.Tags("Client")
	t.Parallel()
	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockClientService(ctrl)

		service.mockTransactionManager.EXPECT().WithinTransaction(ctx, gomock.Any()).Return(nil)

		clientService := createClientService(service)
		err := clientService.DeleteAssignment(uint64(1), uint64(1))

		sCtx.Assert().NoError(err)
	})
}

func (s *ClientSuite) TestDeleteAssignmentFailure(t provider.T) {
	t.Title("DeleteAssignment: Failure")
	t.Tags("Client")
	t.Parallel()
	t.WithNewStep("Failure", func(sCtx provider.StepCtx) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := createMockClientService(ctrl)

		service.mockTransactionManager.EXPECT().WithinTransaction(ctx, gomock.Any()).Return(errors.New("error"))

		clientService := createClientService(service)
		err := clientService.DeleteAssignment(uint64(1), uint64(1))

		sCtx.Assert().Error(err)
	})
}

func TestClientSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ClientSuite))
}
