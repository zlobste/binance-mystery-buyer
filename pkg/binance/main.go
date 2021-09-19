package binance

import (
	"fmt"
	fhttp "github.com/valyala/fasthttp"
	"github.com/zlobste/binance-mystery-buyer/pkg/binance/requests"
)

const (
	URLUserInfo       = "/bapi/accounts/v1/private/account/user/base-detail"
	URLMysteryBoxList = "/bapi/nft/v1/public/nft/mystery-box/list"
	URLMysteryBoxInfo = "/bapi/nft/v1/friendly/nft/mystery-box/detail"
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
	GetSignerInfo() (requests.UserInfo, error)
	GetMysteryBoxList(page, size int64) ([]requests.MysteryBoxInfo, error)
	GetPendingMysteryBoxList() ([]requests.MysteryBoxInfo, error)
	GetMysteryBoxInfo(id string) (requests.MysteryBoxAdvancedInfo, error)
}

type client struct {
	signer Signer
	http   *fhttp.Client
	addr   string
}

func New(addr string) Client {
	return &client{
		http: &fhttp.Client{},
		addr: addr,
	}
}

func (c *client) WithSigner(csrf, cookie string) Client {
	return &client{
		signer: NewAccount(csrf, cookie),
		http:   c.http,
		addr:   c.addr,
	}
}

func (c *client) GetSignerInfo() (requests.UserInfo, error) {
	res, err := c.post(fmt.Sprintf("%s%s", c.addr, URLUserInfo), nil)
	if err != nil {
		return requests.UserInfo{}, err
	}

	result, err := requests.UnmarshalUserInfo(res)
	if err != nil {
		return requests.UserInfo{}, err
	}

	return result.Data, nil
}

func (c *client) GetPendingMysteryBoxList() ([]requests.MysteryBoxInfo, error) {
	boxes, err := c.GetMysteryBoxList(DefaultPage, DefaultPageSize)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	result, err := requests.UnmarshalMysteryBoxList(res)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *client) GetMysteryBoxInfo(id string) (requests.MysteryBoxAdvancedInfo, error) {
	res, err := c.get(fmt.Sprintf("%s%s?productId=%s", c.addr, URLMysteryBoxInfo, id))
	if err != nil {
		return requests.MysteryBoxAdvancedInfo{}, err
	}

	result, err := requests.UnmarshalMysteryBoxInfo(res)
	if err != nil {
		return requests.MysteryBoxAdvancedInfo{}, err
	}

	return result.Data, nil
}
