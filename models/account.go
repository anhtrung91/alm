package models

const ACCOUNT_TABLE = "Account"

// Account
type Account struct {
	ID          int    `json:"ID"`
	UserName    string `json:"UserName"`
	Password    string `json:"Password"`
	TokenAmount string    `json:"TokenAmount"`
	Name        string `json:"Name"`
	Email       string `json:"Email"`
}

type AccountArray []Account
