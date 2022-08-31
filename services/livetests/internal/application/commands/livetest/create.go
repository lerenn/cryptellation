package cmdLivetest

import (
	"context"
	"fmt"

	"github.com/digital-feather/cryptellation/services/livetests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/livetests/internal/domain/livetest"
	ticksClient "github.com/digital-feather/cryptellation/services/ticks/pkg/client"
)

type CreateHandler struct {
	repository  vdb.Port
	ticksClient ticksClient.Client
}

func NewCreateHandler(repository vdb.Port, tClient ticksClient.Client) CreateHandler {
	if repository == nil {
		panic("nil repository")
	}

	if tClient == nil {
		panic("nil tClient")
	}

	return CreateHandler{
		repository:  repository,
		ticksClient: tClient,
	}
}

func (h CreateHandler) Handle(ctx context.Context, req livetest.NewPayload) (id uint, err error) {
	bt, err := livetest.New(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("creating a new livetest from request: %w", err)
	}

	err = h.repository.CreateLivetest(ctx, &bt)
	if err != nil {
		return 0, fmt.Errorf("adding livetest to vdb: %w", err)
	}

	return bt.ID, nil
}
