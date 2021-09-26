package models

import (
	"encoding/json"
	fhttp "github.com/valyala/fasthttp"
)

func UnmarshalSignerInfo(res *fhttp.Response) (*SignerInfoResponse, error) {
	response := new(SignerInfoResponse)
	if err := json.Unmarshal(res.Body(), response); err != nil {
		return nil, err
	}

	return response, nil
}

func UnmarshalMysteryBoxList(res *fhttp.Response) (*MysteryBoxesListResponse, error) {
	response := new(MysteryBoxesListResponse)
	if err := json.Unmarshal(res.Body(), response); err != nil {
		return nil, err
	}

	return response, nil
}

func UnmarshalSignerBalance(res *fhttp.Response) (*UserBalanceResponse, error) {
	response := new(UserBalanceResponse)
	if err := json.Unmarshal(res.Body(), response); err != nil {
		return nil, err
	}

	return response, nil
}

func UnmarshalMysteryBoxInfo(res *fhttp.Response) (*MysteryBoxesInfoResponse, error) {
	response := new(MysteryBoxesInfoResponse)
	if err := json.Unmarshal(res.Body(), response); err != nil {
		return nil, err
	}

	return response, nil
}
