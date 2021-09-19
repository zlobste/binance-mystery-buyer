package requests

import (
	"encoding/json"
	fhttp "github.com/valyala/fasthttp"
)

type UserBalanceRequest struct {
	AssetList []string `json:"assetList"`
	FiatName  string   `json:"fiatName"`
}

type UserBalanceResponse struct {
	Data UserBalanceData `json:"data"`
}

type UserBalanceData struct {
	AssetList []Balance `json:"assetBalanceList"`
}

type Balance struct {
	Asset          string `json:"asset"`
	Free           string `json:"free"`
	Freeze         string `json:"freeze"`
	Total          string `json:"total"`
	LogoUrl        string `json:"logoUrl"`
	TotalFiatValue string `json:"totalFiatValue"`
}

func UnmarshalUserBalance(res *fhttp.Response) (*UserBalanceResponse, error) {
	response := new(UserBalanceResponse)
	if err := json.Unmarshal(res.Body(), response); err != nil {
		return nil, err
	}

	return response, nil
}
