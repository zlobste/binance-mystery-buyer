package models

type BuyMysteryBoxesRequest struct {
	ID     string `json:"productId"`
	Amount int64  `json:"number"`
}

type MysteryBoxesListResponse struct {
	Data []MysteryBoxInfo `json:"data"`
}

type MysteryBoxInfo struct {
	ID            string `json:"productId"`
	MappingStatus int    `json:"mappingStatus"`
}

type MysteryBoxesInfoResponse struct {
	Data MysteryBoxAdvancedInfo `json:"data"`
}

type MysteryBoxAdvancedInfo struct {
	MysteryBoxInfo
	Name         string `json:"name"`
	Price        string `json:"price"`
	Currency     string `json:"currency"`
	StartTime    int64  `json:"startTime"`
	EndTime      int64  `json:"endTime"`
	CurrentStore int64  `json:"currentStore"`
	TotalStore   int64  `json:"totalStore"`
	LimitPerTime int64  `json:"limitPerTime"`
}
