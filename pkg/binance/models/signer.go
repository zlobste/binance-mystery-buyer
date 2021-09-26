package models

type SignerInfoResponse struct {
	Data SignerInfo `json:"data"`
}

type SignerInfo struct {
	Email string `json:"email" structs:"signer_email"`
}

type UserBalanceRequest struct {
	AssetList []string `json:"assetList" structs:"asset_list"`
	FiatName  string   `json:"fiatName" structs:"fiat_name"`
}

type UserBalanceResponse struct {
	Data UserBalanceData `json:"data"`
}

type UserBalanceData struct {
	AssetList []Balance `json:"assetBalanceList"`
}

type Balance struct {
	Asset          string `json:"asset" structs:"asset"`
	Free           string `json:"free" structs:"free"`
	Freeze         string `json:"freeze" structs:"freeze"`
	Total          string `json:"total" structs:"total"`
	LogoUrl        string `json:"logoUrl" structs:"logo_url"`
	TotalFiatValue string `json:"totalFiatValue" structs:"total_fiat"`
}
