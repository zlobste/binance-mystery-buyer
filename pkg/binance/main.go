package binance

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	fhttp "github.com/valyala/fasthttp"
	"github.com/zlobste/binance-mystery-buyer/pkg/binance/requests"
	"net/http"
)

const (
	URLUserInfo       = "/bapi/accounts/v1/private/account/user/base-detail"
	URLMysteryBoxList = "/bapi/nft/v1/public/nft/mystery-box/list"
	URLMysteryBoxInfo = "/bapi/nft/v1/friendly/nft/mystery-box/detail"
	URLBuyMysteryBox  = "/bapi/nft/v1/private/nft/mystery-box/purchase"
	URLUserBalance    = "/bapi/nft/v1/private/nft/user-asset-board"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 100
)

const (
	BoxStatusUpcoming = -1
	BoxStatusSoldOut  = 1
)

type Client interface {
	WithSigner(csrf, cookie string) Client
	GetSignerInfo() (*requests.UserInfo, error)
	GetSignerBalance(assetList []string, fiatName string) ([]requests.Balance, error)
	GetMysteryBoxList(page, size int64) ([]requests.MysteryBoxInfo, error)
	GetUpcomingMysteryBoxList() ([]requests.MysteryBoxInfo, error)
	GetMysteryBoxInfo(id string) (*requests.MysteryBoxAdvancedInfo, error)
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

func (c *client) GetSignerInfo() (*requests.UserInfo, error) {
	res, err := c.post(fmt.Sprintf("%s%s", c.addr, URLUserInfo), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get signer info")
	}

	if res.StatusCode() != http.StatusOK {
		c.log.WithFields(logrus.Fields{
			"status": res.StatusCode(),
			"body":   string(res.Body()),
		}).Error("failed to get signer info")

		return nil, nil
	}

	result, err := requests.UnmarshalUserInfo(res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal user info")
	}

	return &result.Data, nil
}

func (c *client) GetSignerBalance(assetList []string, fiatName string) ([]requests.Balance, error) {
	req := requests.UserBalanceRequest{
		AssetList: assetList,
		FiatName:  fiatName,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal the body")
	}

	res, err := c.post(fmt.Sprintf("%s%s", c.addr, URLUserBalance), body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get signer balance")
	}

	if res.StatusCode() != http.StatusOK {
		c.log.WithFields(logrus.Fields{
			"status": res.StatusCode(),
			"body":   string(res.Body()),
		}).Error("failed to buy the boxes")
	}

	result, err := requests.UnmarshalUserBalance(res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal user balance list")
	}

	return result.Data.AssetList, nil
}

func (c *client) GetUpcomingMysteryBoxList() ([]requests.MysteryBoxInfo, error) {
	boxes, err := c.GetMysteryBoxList(DefaultPage, DefaultPageSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get a list of upcoming sales")
	}

	pendingBoxes := make([]requests.MysteryBoxInfo, 0)
	for _, val := range boxes {
		if val.MappingStatus == BoxStatusUpcoming {
			pendingBoxes = append(pendingBoxes, val)
		}
	}

	return pendingBoxes, nil
}

func (c *client) GetMysteryBoxList(page, size int64) ([]requests.MysteryBoxInfo, error) {
	res, err := c.get(fmt.Sprintf("%s%s/?page=%v&size=%v", c.addr, URLMysteryBoxList, page, size))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get a list of sales")
	}

	if res.StatusCode() != http.StatusOK {
		c.log.WithFields(logrus.Fields{
			"status": res.StatusCode(),
			"body":   string(res.Body()),
		}).Error("failed to get a list of sales")

		return nil, nil
	}

	result, err := requests.UnmarshalMysteryBoxList(res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal a list of sales")
	}

	return result.Data, nil
}

func (c *client) GetMysteryBoxInfo(id string) (*requests.MysteryBoxAdvancedInfo, error) {
	res, err := c.get(fmt.Sprintf("%s%s?productId=%s", c.addr, URLMysteryBoxInfo, id))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get mystery box info")
	}

	if res.StatusCode() != http.StatusOK {
		c.log.WithFields(logrus.Fields{
			"status": res.StatusCode(),
			"body":   string(res.Body()),
		}).Error("failed to get mystery box info")

		return nil, nil
	}

	result, err := requests.UnmarshalMysteryBoxInfo(res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal mystery box info")
	}

	return &result.Data, nil
}

func (c *client) BuyMysteryBox(id string, amount int64) error {
	req := requests.BuyMysteryBoxesRequest{
		ID:     id,
		Amount: amount,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "failed to marshal the body")
	}

	res, err := c.post(fmt.Sprintf("%s%s", c.addr, URLBuyMysteryBox), body)
	if err != nil {
		return errors.Wrap(err, "failed to buy the boxes")
	}

	if res.StatusCode() != http.StatusOK {
		c.log.WithFields(logrus.Fields{
			"status": res.StatusCode(),
			"body":   string(res.Body()),
		}).Error("failed to buy the boxes")
	}

	return nil
}
