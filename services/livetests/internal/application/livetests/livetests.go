package livetests

import (
	"context"
	"fmt"

	"github.com/digital-feather/cryptellation/services/livetests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/livetests/internal/domain/livetest"
)

type Livetests struct {
	repository vdb.Port
}

func New(repository vdb.Port) *Livetests {
	if repository == nil {
		panic("nil repository")
	}

	return &Livetests{
		repository: repository,
	}
}

func (h Livetests) Create(ctx context.Context, req livetest.NewPayload) (id uint, err error) {
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
