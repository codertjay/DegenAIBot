package helper

import (
	"DegenAIBot/config"
)

type HelperInterface interface {
	GetAddressTransactions(address string) (Transactions, error)
	AddAddressTransactions(address string) (addressStatus AddressStatus, err error)
	AddressStatus(address string) (addressStatus AddressStatus, err error)
	CalculateTransactionPNL(transaction Transactions, address string) (pnlMessage string, err error)
	GetUserPortfolio(address string) (portfolio []Portfolio, err error)
	GetAddressPNL(address string) (pnlMessage string, err error)
	CalculatePortfolio(address string) (message string, err error)
	SendTweet(tweet string) (err error)
}

type Helper struct {
	cfg config.Config
}

func NewHelper(cfg config.Config) HelperInterface {
	return &Helper{
		cfg: cfg,
	}
}
