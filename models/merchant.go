package models

const MERCHANT_TABLE = "Merchant"

type Merchant struct {
	MerchantId   int    `json:"MerchantId"`
	MerchantName string `json:"MerchantName"`
	TokenAmount  string `json:"TokenAmount"`
	RemainCoin   string `json:"RemainCoin"`
	SumCoin      string `json:"SumCoin"`
	RateCash     string `json:"RateCash"`
}
