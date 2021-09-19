package requests

import (
	"encoding/json"
	fhttp "github.com/valyala/fasthttp"
)

type UserInfoResponse struct {
	Data UserInfo `json:"data"`
}

type UserInfo struct {
	Email string `json:"email"`
}

func UnmarshalUserInfo(res *fhttp.Response) (*UserInfoResponse, error) {
	response := new(UserInfoResponse)
	if err := json.Unmarshal(res.Body(), response); err != nil {
		return nil, err
	}

	return response, nil
}
