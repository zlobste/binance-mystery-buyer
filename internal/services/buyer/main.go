package buyer

import (
	"github.com/jasonlvhit/gocron"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zlobste/binance-mystery-buyer/internal/config"
	"github.com/zlobste/binance-mystery-buyer/pkg/binance"
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
	Run() error
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

func (s *service) Run() error {
	s.logger.Info("Buyer has started...")

	info, err := s.client.GetSignerInfo()
	if err != nil {
		return errors.Wrap(err, "error on getting signer info")
	}
	s.logger.WithField("signer_info", info).Info("Signed in successfully")

	if err := gocron.Every(SalesCheckPeriod).Hours().Do(s.checkUpcomingSales); err != nil {
		return errors.Wrap(err, "Failed to set scheduled check upcoming sales")
	}
	<-gocron.Start()

	return nil
}

func (s *service) checkUpcomingSales() {
	s.RLock()
	defer s.RUnlock()

	sales, err := s.client.GetUpcomingMysteryBoxList()
	if err != nil {
		s.logger.WithError(err).Error("Failed to get upcoming sales")

		return
	}
	s.logger.WithField("upcoming sales", sales).Info("Upcoming sales have been checked")

	for _, sale := range sales {
		_, exists := s.pendingJobs[sale.ID]
		if !exists {
			s.pendingJobs[sale.ID] = true
			s.logger.WithField("sale_id", sale.ID).Info("Upcoming sale has been added to buyer list")

			go s.prepareToBuy(sale.ID)
		}
	}
}

func (s *service) prepareToBuy(saleID string) {
	box, err := s.client.GetMysteryBoxInfo(saleID)
	if err != nil {
		s.logger.WithError(err).WithField("sale_id", saleID).Error("Failed to get mystery box info")

		return
	}
	s.logger.WithField("box_info", box).Info("Got mystery box info")

	startNano := time.Unix(box.StartTime, 0).UTC().UnixNano()
	nowNano := time.Now().UTC().UnixNano()
	time.Sleep(time.Duration(startNano-nowNano)*time.Nanosecond - 5*time.Minute)

	balances, err := s.client.GetSignerBalance(DefaultFiat, box.Currency)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get signer balances")

		return
	}

	price, err := strconv.ParseFloat(box.Price, 64)
	if err != nil {
		s.logger.WithError(err).Error("Failed to parse box price")

		return
	}

	freeBalance, err := strconv.ParseFloat(balances[0].Free, 64)
	if err != nil {
		s.logger.WithError(err).Error("Failed to parse signer balance")

		return
	}
	if freeBalance == 0 {
		s.logger.WithField("balances", balances).Error("Zero balance, failed to buy boxes")

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
	s.logger.WithFields(logrus.Fields{
		"sale_id": box.ID,
		"to_buy":  ableToBuy,
	}).Info("Trying to buy mystery boxes")

	nowNano = time.Now().UTC().UnixNano()
	time.Sleep(time.Duration(startNano-nowNano) * time.Nanosecond)

	s.buyBox(box.ID, ableToBuy)
}

func (s *service) buyBox(saleID string, countToBuy int64) {
	s.RLock()
	defer s.RUnlock()
	defer delete(s.pendingJobs, saleID)

	if err := s.client.BuyMysteryBox(saleID, countToBuy); err != nil {
		s.logger.WithError(err).Error("Failed to buy mystery boxes")

		return
	}
}
