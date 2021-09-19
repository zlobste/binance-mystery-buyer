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
	log := cfg.Log()
	return &service{
		client: binance.New(auth.Proxy, log).WithSigner(auth.CSRFToken, auth.Cookie),
		config: cfg,
		logger: log,
	}
}

func (s *service) Run(ctx context.Context) error {
	s.logger.Info("Buyer has started...")

	info, err := s.client.GetSignerInfo()
	if err != nil {
		return errors.Wrap(err, "error on getting user info")
	}
	if info != nil {
		fmt.Println(fmt.Sprintf("User info: %s", info))
	}

	fiatName := "USD"
	assets := []string{"BNB", "BUSD", "ETH"}
	balance, err := s.client.GetSignerBalance(assets, fiatName)
	if err != nil {
		return errors.Wrap(err, "error on getting user balance")
	}
	if info != nil {
		fmt.Println(fmt.Sprintf("User balance: %s", balance))
	}

	//boxes, err := s.client.GetUpcomingMysteryBoxList()
	//if err != nil {
	//	return errors.Wrap(err, "error on getting list of mystery boxes")
	//}
	//if boxes != nil {
	//	if err := s.client.BuyMysteryBox(boxes[0].ID, 20); err != nil {
	//		return errors.Wrap(err, "error on buying box")
	//	}
	//}

	//fmt.Println(fmt.Sprintf("Mystery boxes: %v", boxes))
	//
	//box, err := s.client.GetMysteryBoxInfo(boxes[0].ID)
	//if err != nil {
	//	return errors.Wrap(err, "error on getting box info")
	//}
	//fmt.Println(fmt.Sprintf("Box: %v", box))

	return nil
}
