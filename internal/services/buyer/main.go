package buyer

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zlobste/binance-mystery-buyer/internal/config"
)

type Service interface {
	Run(ctx context.Context) error
}

type service struct {
	config config.Config
	logger *logrus.Logger
}

func New(cfg config.Config) Service {
	return &service{
		config: cfg,
		logger: cfg.Log(),
	}
}

func (s *service) Run(ctx context.Context) error {
	return nil
}
