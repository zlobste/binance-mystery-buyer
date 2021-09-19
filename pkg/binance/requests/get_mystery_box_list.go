package requests

import (
	"encoding/json"
	fhttp "github.com/valyala/fasthttp"
)

type MysteryBoxesListResponse struct {
	Data []MysteryBoxInfo `json:"data"`
}

type MysteryBoxInfo struct {
	ID            string `json:"productId"`
	MappingStatus int    `json:"mappingStatus"`
}

func UnmarshalMysteryBoxList(res *fhttp.Response) (*MysteryBoxesListResponse, error) {
	response := new(MysteryBoxesListResponse)
	if err := json.Unmarshal(res.Body(), response); err != nil {
		return nil, err
	}

	return response, nil
}
