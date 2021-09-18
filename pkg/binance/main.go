package binance

import (
	"encoding/json"
	"fmt"
	fhttp "github.com/valyala/fasthttp"
	"github.com/zlobste/binance-mystery-buyer/pkg/binance/account"
)

type Client interface {
	WithSigner(csrf, cookie string) Client
	GetSignerInfo() (string, error)
}

type client struct {
	signer account.Account
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
		signer: account.NewAccount(csrf, cookie),
		http:   c.http,
		addr:   c.addr,
	}
}

func (c *client) GetSignerInfo() (string, error) {
	r := &fhttp.Request{}
	req := c.signer.SignRequest(r)
	req.Header.SetRequestURI(fmt.Sprintf("%s%s", c.addr, account.URLUserInfo))
	req.Header.SetMethod(fhttp.MethodPost)
	req.Header.SetContentType("application/json")

	res := fhttp.AcquireResponse()
	if err := c.http.Do(req, res); err != nil {
		return "", err
	}

	response := struct {
		Data struct {
			Email string `json:"email"`
		} `json:"data"`
	}{}

	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return "", err
	}

	return response.Data.Email, nil
}
