package binance

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	fhttp "github.com/valyala/fasthttp"
	"github.com/zlobste/binance-mystery-buyer/pkg/binance/models"
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

	BoxStatusUpcoming = -1
)

func (c *client) GetSignerInfo() (*models.SignerInfo, error) {
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

	result, err := models.UnmarshalSignerInfo(res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal user info")
	}

	return &result.Data, nil
}

func (c *client) GetSignerBalance(fiatName string, assetList ...string) ([]models.Balance, error) {
	req := models.UserBalanceRequest{
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

	result, err := models.UnmarshalSignerBalance(res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal user balance list")
	}

	return result.Data.AssetList, nil
}

func (c *client) GetUpcomingMysteryBoxList() ([]models.MysteryBoxInfo, error) {
	boxes, err := c.GetMysteryBoxList(DefaultPage, DefaultPageSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get a list of upcoming sales")
	}

	pendingBoxes := make([]models.MysteryBoxInfo, 0)
	for _, val := range boxes {
		if val.MappingStatus == BoxStatusUpcoming {
			pendingBoxes = append(pendingBoxes, val)
		}
	}

	return pendingBoxes, nil
}

func (c *client) GetMysteryBoxList(page, size int64) ([]models.MysteryBoxInfo, error) {
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

	result, err := models.UnmarshalMysteryBoxList(res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal a list of sales")
	}

	return result.Data, nil
}

func (c *client) GetMysteryBoxInfo(id string) (*models.MysteryBoxAdvancedInfo, error) {
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

	result, err := models.UnmarshalMysteryBoxInfo(res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal mystery box info")
	}

	return &result.Data, nil
}

func (c *client) BuyMysteryBox(id string, amount int64) error {
	req := models.BuyMysteryBoxesRequest{
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

func (c *client) get(url string) (*fhttp.Response, error) {
	res := fhttp.AcquireResponse()
	req := c.signer.SignRequest(fhttp.AcquireRequest())
	req.Header.SetRequestURI(url)
	req.Header.SetMethod(fhttp.MethodGet)

	if err := c.http.Do(req, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *client) post(url string, body []byte) (*fhttp.Response, error) {
	res := fhttp.AcquireResponse()
	req := c.signer.SignRequest(fhttp.AcquireRequest())
	req.Header.SetRequestURI(url)
	req.Header.SetMethod(fhttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBody(body)

	if err := c.http.Do(req, res); err != nil {
		return nil, err
	}

	return res, nil
}
