package buyer

import (
	"context"
	"github.com/jasonlvhit/gocron"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zlobste/binance-mystery-buyer/internal/config"
	"github.com/zlobste/binance-mystery-buyer/pkg/binance"
	"github.com/zlobste/binance-mystery-buyer/pkg/binance/requests"
	"math"
	"strconv"
	"sync"
	"time"
)

const (
	DefaultFiat      = "USD"
	SalesCheckPeriod = 12
)

type Service interface {
	Run(ctx context.Context) error
}

type service struct {
	client binance.Client
	config config.Config
	logger *logrus.Logger

	pendingJobs map[string]bool

	sync.RWMutex
}

func New(cfg config.Config) Service {
	auth := cfg.GetAuth()
	log := cfg.Log()
	return &service{
		client:      binance.New(auth.Proxy, log).WithSigner(auth.CSRFToken, auth.Cookie),
		config:      cfg,
		logger:      log,
		pendingJobs: make(map[string]bool, 1),
	}
}

func (s *service) Run(ctx context.Context) error {
	s.logger.Info("Buyer has started...")

	info, err := s.client.GetSignerInfo()
	if err != nil {
		return errors.Wrap(err, "error on getting user info")
	}
	s.logger.WithField("signer_info", info).Info("Signed in successfully")

	if err := gocron.Every(SalesCheckPeriod).Hours().Do(s.checkUpcomingSales); err != nil {
		s.logger.WithError(err).Error("failed to check upcoming sales")
	}
	<-gocron.Start()

	return nil
}

func (s *service) checkUpcomingSales() {
	s.RLock()
	defer s.RUnlock()

	sales, err := s.client.GetUpcomingMysteryBoxList()
	if err != nil {
		s.logger.WithError(err).Error("failed to get upcoming sales")

		return
	}

	for _, sale := range sales {
		_, exists := s.pendingJobs[sale.ID]
		if !exists {
			s.pendingJobs[sale.ID] = true

			go s.prepareToBuy(sale.ID)
		}
	}
}

func (s *service) prepareToBuy(id string) {
	info, err := s.client.GetMysteryBoxInfo(id)
	if err != nil {
		s.logger.WithError(err).Error("failed to get mystery box info")

		return
	}

	s.buyBox(*info)
}

func (s *service) buyBox(box requests.MysteryBoxAdvancedInfo) {
	s.RLock()
	defer s.RUnlock()
	defer delete(s.pendingJobs, box.ID)

	startNano := time.Unix(box.StartTime, 0).UnixNano()
	nowNano := time.Now().UnixNano()
	time.Sleep(time.Duration(startNano-nowNano)*time.Nanosecond - 5*time.Minute)

	balances, err := s.client.GetSignerBalance(DefaultFiat, box.Currency)
	if err != nil {
		s.logger.WithError(err).Error("failed to get signer balances")

		return
	}

	price, err := strconv.ParseFloat(box.Price, 64)
	if err != nil {
		s.logger.WithError(err).Error("failed to parse box price")

		return
	}

	freeBalance, err := strconv.ParseFloat(balances[0].Free, 64)
	if err != nil {
		s.logger.WithError(err).Error("failed to parse signer balance")

		return
	}
	if freeBalance == 0 {
		s.logger.Info("Zero balance")

		return
	}

	ableToBuy := int64(math.Round(freeBalance / price))
	if ableToBuy == 0 {
		s.logger.WithError(err).Error("There is not enough free balance to buy boxes")

		return
	}

	if box.LimitPerTime < ableToBuy {
		ableToBuy = box.LimitPerTime
	}

	nowNano = time.Now().UnixNano()
	time.Sleep(time.Duration(startNano-nowNano) * time.Nanosecond)

	if err := s.client.BuyMysteryBox(box.ID, ableToBuy); err != nil {
		s.logger.WithError(err).Error("failed to buy mystery boxes")

		return
	}
}
