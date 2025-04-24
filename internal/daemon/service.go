package daemon

import (
	"context"

	"github.com/rs/zerolog/log"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Start(ctx context.Context) error {
	log.Debug().Msg("[Service] start")
	return nil
}
