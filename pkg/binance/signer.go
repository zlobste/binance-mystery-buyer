package binance

import (
	fhttp "github.com/valyala/fasthttp"
)

type Signer interface {
	SignRequest(req *fhttp.Request) *fhttp.Request
}

type signer struct {
	Csrf   string
	Cookie string
}

func NewAccount(csrf, cookie string) Signer {
	return &signer{
		Csrf:   csrf,
		Cookie: cookie,
	}
}

func (a *signer) SignRequest(req *fhttp.Request) *fhttp.Request {
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("clienttype", "web")
	req.Header.Set("cookie", a.Cookie)
	req.Header.Set("csrftoken", a.Csrf)

	return req
}
