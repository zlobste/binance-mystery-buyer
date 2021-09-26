package binance

import (
	"github.com/sirupsen/logrus"
	fhttp "github.com/valyala/fasthttp"
	"github.com/zlobste/binance-mystery-buyer/pkg/binance/models"
)

type Client interface {
	WithSigner(csrf, cookie string) Client
	GetSignerInfo() (*models.SignerInfo, error)
	GetSignerBalance(fiatName string, assetList ...string) ([]models.Balance, error)
	GetMysteryBoxList(page, size int64) ([]models.MysteryBoxInfo, error)
	GetUpcomingMysteryBoxList() ([]models.MysteryBoxInfo, error)
	GetMysteryBoxInfo(id string) (*models.MysteryBoxAdvancedInfo, error)
	BuyMysteryBox(id string, amount int64) error
}

type client struct {
	signer Signer
	addr   string
	http   *fhttp.Client
	log    *logrus.Logger
}

func New(addr string, log *logrus.Logger) Client {
	return &client{
		http: &fhttp.Client{},
		addr: addr,
		log:  log,
	}
}

func (c *client) WithSigner(csrf, cookie string) Client {
	return &client{
		signer: NewAccount(csrf, cookie),
		http:   c.http,
		addr:   c.addr,
		log:    c.log,
	}
}
