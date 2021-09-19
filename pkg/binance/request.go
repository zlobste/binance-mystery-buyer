package binance

import (
	fhttp "github.com/valyala/fasthttp"
)

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
