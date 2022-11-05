package application

import (
	"github.com/digital-feather/cryptellation/services/livetests/internal/adapters/vdb/redis"
	"github.com/digital-feather/cryptellation/services/livetests/internal/application/livetests"
)

type Application struct {
	Livetests livetests.Operator
}

func New() (*Application, error) {
	repository, err := redis.New()
	if err != nil {
		return nil, err
	}

	return &Application{
		Livetests: livetests.New(repository),
	}, nil
}
