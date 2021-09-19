package requests

type BuyMysteryBoxesRequest struct {
	ID     string `json:"productId"`
	Amount int64  `json:"number"`
}
