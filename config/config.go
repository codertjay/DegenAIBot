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
		"D33wobbUPo9ChS8aR1RABHRYcaYZDmyajdvEM6RyLQ6a",
		"ZtNoGMyvqHsnv1iYEccMRUsWoQc62SMabFLxoVATEWd",
		"3HF5GGumrr8YAuMsC4RDycVAU9JbFE6STDZpFtkXDdYZ",
		"GM8Fvy48LSuwZ1rd1KzNgWP5xh3jDDfbCMsZZDrFck2q",
		"9osj5ksZcrYVeft5pjN1KUNpuM9xiZpjLPxtiMHG4AT5",
		"9QpiMJVDiCnwgFFTmZ3xbMznGgbdE1DA4QCwnos3D6sx",
		"5b3fdi6ZMfRjwN1DobL8qLMWF9BzU2jsJidojkmNj764",
		"4XM2PQQxoVVi7jZnGw39BGERK9icmi7zo4CcqzAqbj7y",
		"F2cGAP6qJDA451KwsdTYzs22vkhsPHJkVrZZ34y3sHFa",
		"9KLp5H95hZVEipcF9dpD8xbuZmSztebpEX3EF5KQkas9",
		"GM8Fvy48LSuwZ1rd1KzNgWP5xh3jDDfbCMsZZDrFck2q",
		"6WH2umAtRxxBvDxMeo9EE92EnMjQ3ucqg2VmkYrcL1gZ",
		"5iCopJuJxQnSUvBFHvtVRgda8XpUPQwQrsRFpMDjbB7h",
		"HCA5PfZumwMAeAmNATqNejuAhzGXWKKa1Kma9Q4MfMfo",
		"HCA5PfZumwMAeAmNATqNejuAhzGXWKKa1Kma9Q4MfMfo",
		"26rogKZkjAjtHTZnZPibpvakyCRp8FMuGu8uKXaJEdpq",
		"26rogKZkjAjtHTZnZPibpvakyCRp8FMuGu8uKXaJEdpq",
		"EByTfm7qLtGEqx4fRRt1TppZGJuxEYVNNY8oABVx1f66",
		"EByTfm7qLtGEqx4fRRt1TppZGJuxEYVNNY8oABVx1f66",
		"AhRywJxQmwDXgCdQWLo5AB5Y7y6gaahi4ax7AcskmSJ8",
		"AhRywJxQmwDXgCdQWLo5AB5Y7y6gaahi4ax7AcskmSJ8",
		"9inmKn7KywpkugeXSqdv2y84bgjx6jefxkiBCE6F4id8",
		"9inmKn7KywpkugeXSqdv2y84bgjx6jefxkiBCE6F4id8",
		"Fpbo3kPgFq9YBEXo6TYpTE1tekb7U6ZBjymsS1TMyAAJ",
		"Fpbo3kPgFq9YBEXo6TYpTE1tekb7U6ZBjymsS1TMyAAJ",
	}
	return cfg, nil
}
