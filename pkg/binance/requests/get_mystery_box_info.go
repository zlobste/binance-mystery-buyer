package requests

import (
	"encoding/json"
	fhttp "github.com/valyala/fasthttp"
)

type MysteryBoxesInfoResponse struct {
	Data MysteryBoxAdvancedInfo `json:"data"`
}

type MysteryBoxAdvancedInfo struct {
	ID            string `json:"productId"`
	Name          string `json:"name"`
	Price         string `json:"price"`
	Currency      string `json:"currency"`
	StartTime     uint64 `json:"startTime"`
	EndTime       uint64 `json:"endTime"`
	MappingStatus int    `json:"mappingStatus"`
	CurrentStore  uint64 `json:"currentStore"`
	TotalStore    uint64 `json:"totalStore"`
	LimitPerTime  uint64 `json:"limitPerTime"`
}

func UnmarshalMysteryBoxInfo(res *fhttp.Response) (*MysteryBoxesInfoResponse, error) {
	response := new(MysteryBoxesInfoResponse)
	if err := json.Unmarshal(res.Body(), response); err != nil {
		return nil, err
	}

	return response, nil
}
