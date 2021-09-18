package account

import (
	fhttp "github.com/valyala/fasthttp"
)

const (
	URLUserInfo = "/bapi/accounts/v1/private/account/user/base-detail"
)

type Account interface {
	SignRequest(req *fhttp.Request) *fhttp.Request
	GetInfoRequest() *fhttp.Request
}

type account struct {
	Csrf   string
	Cookie string
}

func NewAccount(csrf, cookie string) Account {
	return &account{
		Csrf:   csrf,
		Cookie: cookie,
	}
}

func (a *account) GetInfoRequest() *fhttp.Request {
	return a.SignRequest(
		&fhttp.Request{},
	)
}

func (a *account) SignRequest(req *fhttp.Request) *fhttp.Request {
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("clienttype", "web")
	req.Header.Set("cookie", a.Cookie)
	req.Header.Set("csrftoken", a.Csrf)

	return req
}
