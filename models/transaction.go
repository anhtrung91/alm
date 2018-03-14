package models

const TRANSACTION_TABLE = "Transaction"

type Transaction struct {
	TransactionId       int    `json:"TransactionId"`
	UserMerchantId      int    `json:"UserMerchantId"`
	Amount              string `json:"Amount"`
	ReceiveValue        string `json:"ReceiveValue"`
	ReceiveMerchantId   string `json:"ReceiveMerchantId"`
	ReceiveProductId    string `json:"ReceiveProductId"`
	ReceiveRedeemId     string `json:"ReceiveRedeemId"`
	Blockid             string `json:"Blockid"`
	TransactionDatetime string `json:"TransactionDatetime"`
	MerchantId          int    `json:"MerchantId"`
}
