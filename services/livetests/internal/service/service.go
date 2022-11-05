package service

import (
	"github.com/digital-feather/cryptellation/services/livetests/internal/application"
)

func NewApplication() (*application.Application, error) {
	return application.New()
}

func NewMockedApplication() (*application.Application, error) {
	return application.New()
}
