package buyer

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zlobste/binance-mystery-buyer/internal/config"
	"github.com/zlobste/binance-mystery-buyer/pkg/binance"
)

type Service interface {
	Run(ctx context.Context) error
}

type service struct {
	client binance.Client
	config config.Config
	logger *logrus.Logger
}

func New(cfg config.Config) Service {
	auth := cfg.GetAuth()
	return &service{
		client: binance.New(auth.Proxy).WithSigner(auth.CSRFToken, auth.Cookie),
		config: cfg,
		logger: cfg.Log(),
	}
}

func (s *service) Run(ctx context.Context) error {
	s.logger.Info("Buyer has started...")

	info, err := s.client.GetSignerInfo()
	if err != nil {
		return errors.Wrap(err, "error on getting user info")
	}

	fmt.Println(fmt.Sprintf("User info: %s", info))

	return nil
}
