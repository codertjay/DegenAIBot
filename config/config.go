package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Environment string

// Config defines app config
type Config struct {
	CoinStatsAPIKey             string   `envconfig:"COIN_STATS_API_KEY" default:""`
	CoinStatsTransactionsRoute  string   `envconfig:"COIN_STATS_TRANSACTIONS_ROUTE" default:"https://openapiv1.coinstats.app/wallet/transactions"`
	CoinStatsAddressStatusRoute string   `envconfig:"COIN_STATS_ADDRESS_STATUS_ROUTE" default:"https://openapiv1.coinstats.app/wallet/status"`
	CoinStatsAddAddressRoute    string   `envconfig:"COIN_STATS_ADD_ADDRESS_ROUTE" default:"https://openapiv1.coinstats.app/wallet/transactions"`
	CoinStatsPortfolioRoute     string   `envconfig:"COIN_STATS_PORTFOLIO_ROUTE" default:"https://openapiv1.coinstats.app/wallet/balance"`
	TwitterBearerToken          string   `envconfig:"TWITTER_BEARER_TOKEN" default:""`
	TwitterSendTweetRoute       string   `envconfig:"TWITTER_SEND_TWEET_ROUTE" default:"https://api.twitter.com/2/tweets"`
	TwitterAPIKey               string   `envconfig:"TWITTER_API_KEY" default:""`
	TwitterAPISecretKey         string   `envconfig:"TWITTER_API_SECRET_KEY" default:""`
	TwitterAccessToken          string   `envconfig:"TWITTER_ACCESS_TOKEN" default:""`
	TwitterAccessTokenSecret    string   `envconfig:"TWITTER_ACCESS_TOKEN_SECRET" default:""`
	SolanaAddresses             []string `envconfig:"SOLANA_ADDRESSES" default:""`
}

func Load(envFile string) (Config, error) {
	var cfg Config

	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return Config{}, fmt.Errorf("error loading .env file: %w", err)
		}
	} else {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file ", err.Error())
		}
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}

	cfg.SolanaAddresses = []string{
		"HdxkiXqeN6qpK2YbG51W23QSWj3Yygc1eEk2zwmKJExp",
		"FGZtWHYMww8vCEcgbdLjv22QoEtzGgX8UqB7Efvf8ZWf",
		"2QLWGa5Dhos5zNwXL17JRYPT9S13riFpqWJPgUQhuVai",
		"4EsY8HQB4Ak65diFrSHjwWhKSGC8sKmnzyusM993gk2w",
		"45yBcpnzFTqLYQJtjxsa1DdZkgrTYponCg6yLQ6LQPu6",
		"DNfuF1L62WWyW3pNakVkyGGFzVVhj4Yr52jSmdTyeBHm",
	}
	return cfg, nil
}
