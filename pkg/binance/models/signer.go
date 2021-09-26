package models

type SignerInfoResponse struct {
	Data SignerInfo `json:"data"`
}

type SignerInfo struct {
	Email string `json:"email"`
}

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
