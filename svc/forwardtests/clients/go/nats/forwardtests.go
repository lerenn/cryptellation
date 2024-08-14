package nats

import (
	"context"

	helpers "cryptellation/internal/asyncapi"
	"cryptellation/internal/config"
	common "cryptellation/pkg/client"
	"cryptellation/pkg/models/account"

	asyncapi "cryptellation/svc/forwardtests/api/asyncapi"
	client "cryptellation/svc/forwardtests/clients/go"
	"cryptellation/svc/forwardtests/pkg/forwardtest"

	"github.com/google/uuid"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	natsextension "github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
)

type nats struct {
	broker *natsextension.Controller
	ctrl   *asyncapi.UserController
	logger extensions.Logger
	name   string
}

func New(c config.NATS, options ...option) (client.Client, error) {
	var cl nats
	var err error

	// Execute options
	for _, option := range options {
		option(&cl)
	}

	// Create a NATS Controller
	cl.broker, err = natsextension.NewController(c.URL())
	if err != nil {
		return nats{}, err
	}

	// Create a logger if asked
	ctrlOpts := make([]asyncapi.ControllerOption, 0)
	if cl.logger != nil {
		ctrlOpts = append(ctrlOpts, asyncapi.WithLogger(cl.logger))
	} else {
		cl.logger = extensions.DummyLogger{}
	}

	// Create a new user controller
	ctrl, err := asyncapi.NewUserController(cl.broker, ctrlOpts...)
	if err != nil {
		return nats{}, err
	}
	cl.ctrl = ctrl

	return cl, nil
}

func (cl nats) CreateForwardTest(ctx context.Context, payload forwardtest.NewPayload) (uuid.UUID, error) {
	// Set message
	reqMsg := asyncapi.NewCreateRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.CreateRequestChannelPath, cl.name)
	reqMsg.Set(payload)

	// Send request
	respMsg, err := cl.ctrl.RequestToCreateOperation(ctx, reqMsg)
	if err != nil {
		return uuid.Nil, err
	}

	// Unwrap error from message
	if err := helpers.UnwrapError(respMsg.Payload.Error); err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(respMsg.Payload.Id)

}

func (cl nats) ListForwardTests(ctx context.Context) ([]uuid.UUID, error) {
	// Set message
	reqMsg := asyncapi.NewListRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ListRequestChannelPath, cl.name)

	// Send request
	respMsg, err := cl.ctrl.RequestToListOperation(ctx, reqMsg)
	if err != nil {
		return nil, err
	}

	// Unwrap error from message
	if err := helpers.UnwrapError(respMsg.Payload.Error); err != nil {
		return nil, err
	}

	// Convert response to model
	return respMsg.ToModel()
}

func (cl nats) CreateOrder(ctx context.Context, payload common.OrderCreationPayload) error {
	// Set message
	reqMsg := asyncapi.NewOrdersCreateRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.OrdersCreateRequestChannelPath, cl.name)
	reqMsg.Set(payload)

	// Send request
	respMsg, err := cl.ctrl.RequestToOrdersCreateOperation(ctx, reqMsg)
	if err != nil {
		return err
	}

	// Unwrap error from message
	return helpers.UnwrapError(respMsg.Payload.Error)
}

func (cl nats) GetAccounts(ctx context.Context, forwardTestID uuid.UUID) (map[string]account.Account, error) {
	// Set message
	reqMsg := asyncapi.NewAccountsListRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.AccountsListRequestChannelPath, cl.name)
	reqMsg.Payload.Id = asyncapi.ForwardTestIDSchema(forwardTestID.String())

	// Send request
	respMsg, err := cl.ctrl.RequestToAccountsListOperation(ctx, reqMsg)
	if err != nil {
		return nil, err
	}

	// Unwrap error from message
	if err := helpers.UnwrapError(respMsg.Payload.Error); err != nil {
		return nil, err
	}

	// Convert response to model
	return respMsg.ToModel(), nil
}

func (cl nats) GetStatus(ctx context.Context, forwardTestID uuid.UUID) (forwardtest.Status, error) {
	// Set message
	reqMsg := asyncapi.NewStatusRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.StatusRequestChannelPath, cl.name)
	reqMsg.Payload.Id = asyncapi.ForwardTestIDSchema(forwardTestID.String())

	// Send request
	respMsg, err := cl.ctrl.RequestToGetStatusOperation(ctx, reqMsg)
	if err != nil {
		return forwardtest.Status{}, err
	}

	// Unwrap error from message
	if err := helpers.UnwrapError(respMsg.Payload.Error); err != nil {
		return forwardtest.Status{}, err
	}

	return respMsg.ToModel(), nil
}

func (cl nats) ServiceInfo(ctx context.Context) (common.ServiceInfo, error) {
	// Set message
	reqMsg := asyncapi.NewServiceInfoRequestMessage()
	reqMsg.Headers.ReplyTo = helpers.AddReplyToSuffix(asyncapi.ServiceInfoRequestChannelPath, cl.name)

	// Send request
	respMsg, err := cl.ctrl.RequestToServiceInfoOperation(ctx, reqMsg)
	if err != nil {
		return common.ServiceInfo{}, err
	}

	return respMsg.ToModel(), nil
}

func (cl nats) Close(ctx context.Context) {
	cl.ctrl.Close(ctx)
	cl.broker.Close()
}
