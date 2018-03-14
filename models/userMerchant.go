package models

const USER_MERCHANT_TABLE = "UserMerchant"

// User Merchant
type UserMerchant struct {
	UserMerchantId		int			`json:"UserMerchantId"`
	UserName 			string		`json:"UserName"`
	Password 			string 		`json:"Password"`
	MerchantId 			int 		`json:"MerchantId"`
	BalanceCoin 		string 		`json:"BalanceCoin"`
	UserId 				int 		`json:"UserId"`
}