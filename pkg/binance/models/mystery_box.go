package models

type BuyMysteryBoxesRequest struct {
	ID     string `json:"productId" structs:"sale_id"`
	Amount int64  `json:"number" structs:"count"`
}

type MysteryBoxesListResponse struct {
	Data []MysteryBoxInfo `json:"data"`
}

type MysteryBoxInfo struct {
	ID            string `json:"productId" structs:"sale_id"`
	MappingStatus int    `json:"mappingStatus" structs:"mapping_status"`
}

type MysteryBoxesInfoResponse struct {
	Data MysteryBoxAdvancedInfo `json:"data"`
}

type MysteryBoxAdvancedInfo struct {
	MysteryBoxInfo
	Name         string `json:"name" structs:"name"`
	Price        string `json:"price" structs:"price"`
	Currency     string `json:"currency" structs:"currency"`
	StartTime    int64  `json:"startTime" structs:"start_time"`
	EndTime      int64  `json:"endTime" structs:"end_time"`
	CurrentStore int64  `json:"currentStore" structs:"current_store"`
	TotalStore   int64  `json:"totalStore" structs:"total_store"`
	LimitPerTime int64  `json:"limitPerTime" structs:"limit_per_time"`
}
