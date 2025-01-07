package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Transactions struct {
	Result []struct {
		Type        string    `json:"type"`
		Date        time.Time `json:"date"`
		MainContent struct {
			CoinIcons  []string      `json:"coinIcons"`
			CoinAssets []interface{} `json:"coinAssets"`
		} `json:"mainContent"`
		CoinData struct {
			Count        float64 `json:"count"`
			Symbol       string  `json:"symbol"`
			CurrentValue float64 `json:"currentValue"`
		} `json:"coinData"`
		ProfitLoss struct {
			Profit        float64 `json:"profit"`
			ProfitPercent float64 `json:"profitPercent"`
			CurrentValue  float64 `json:"currentValue"`
		} `json:"profitLoss"`
		Transactions []struct {
			Action string `json:"action"`
			Items  []struct {
				Id         string  `json:"id"`
				Count      float64 `json:"count"`
				TotalWorth float64 `json:"totalWorth"`
				Coin       struct {
					Id     string `json:"id"`
					Name   string `json:"name"`
					Symbol string `json:"symbol"`
					Icon   string `json:"icon"`
				} `json:"coin"`
			} `json:"items"`
		} `json:"transactions"`
	} `json:"result"`
	Meta struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
	} `json:"meta"`
}

type Portfolio struct {
	CoinId          string  `json:"coinId"`
	Amount          float64 `json:"amount"`
	Decimals        int     `json:"decimals,omitempty"`
	ContractAddress string  `json:"contractAddress,omitempty"`
	Chain           string  `json:"chain"`
	Name            string  `json:"name"`
	Symbol          string  `json:"symbol"`
	Price           float64 `json:"price"`
	PriceBtc        float64 `json:"priceBtc"`
	ImgUrl          string  `json:"imgUrl"`
	Rank            int     `json:"rank"`
	Volume          float64 `json:"volume"`
	PCh24H          float64 `json:"pCh24h,omitempty"`
}

type AddressStatus struct {
	Status string `json:"status"`
}

func (h *Helper) coinStatsHelper(ctx context.Context, payload []byte, url string, requestMethod string) ([]byte, int, error) {

	req, _ := http.NewRequest(requestMethod, url, bytes.NewBuffer(payload))

	coinStatsSecretKey := h.cfg.CoinStatsAPIKey

	req.Header.Add("accept", "text/plain")
	req.Header.Add("content-type", "application/json")

	req.Header.Add("X-API-KEY", coinStatsSecretKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}

	return body, res.StatusCode, nil
}

func (h *Helper) GetAddressTransactions(address string) (transactions Transactions, err error) {

	from := time.Now().AddDate(0, 0, -7).Format("2006-01-02T15:04:05Z")
	to := time.Now().AddDate(0, 0, 0).Format("2006-01-02T15:04:05Z")
	url := fmt.Sprintf("%s?address=%s&from=%s&to=%s&connectionId=solana", h.cfg.CoinStatsTransactionsRoute, address, from, to)

	response, statusCode, err := h.coinStatsHelper(context.Background(), nil, url, "GET")
	if err != nil {
		return Transactions{}, err
	}

	if statusCode != 200 {
		return transactions, fmt.Errorf("Error getting : %s", response)
	}

	err = json.Unmarshal(response, &transactions)
	if err != nil {
		return Transactions{}, err
	}
	return transactions, nil
}

func (h *Helper) AddAddressTransactions(address string) (addressStatus AddressStatus, err error) {
	url := fmt.Sprintf("%s?address=%s&connectionId=solana", h.cfg.CoinStatsAddAddressRoute, address)
	response, statusCode, err := h.coinStatsHelper(context.Background(), nil, url, "PATCH")
	if err != nil {
		return addressStatus, err
	}

	if statusCode != 200 {
		return addressStatus, fmt.Errorf("Error getting portfolio: %s", response)
	}

	err = json.Unmarshal(response, &addressStatus)
	if err != nil {
		return addressStatus, err
	}
	return addressStatus, nil
}

func (h *Helper) AddressStatus(address string) (addressStatus AddressStatus, err error) {
	url := fmt.Sprintf("%s?address=%s&connectionId=solana", h.cfg.CoinStatsAddressStatusRoute, address)
	response, statusCode, err := h.coinStatsHelper(context.Background(), nil, url, "GET")
	if err != nil {
		return addressStatus, err
	}

	if statusCode != 200 {
		return addressStatus, fmt.Errorf("Error getting portfolio: %s", response)
	}

	err = json.Unmarshal(response, &addressStatus)
	if err != nil {
		return addressStatus, err
	}
	return addressStatus, nil
}

func (h *Helper) GetUserPortfolio(address string) (portfolio []Portfolio, err error) {
	url := fmt.Sprintf("%s?address=%s&connectionId=solana", h.cfg.CoinStatsPortfolioRoute, address)
	response, statusCode, err := h.coinStatsHelper(context.Background(), nil, url, "GET")
	if err != nil {
		return portfolio, err
	}

	if statusCode != 200 {
		return portfolio, fmt.Errorf("Error getting portfolio: %s", response)
	}

	err = json.Unmarshal(response, &portfolio)
	if err != nil {
		return portfolio, err
	}
	return portfolio, nil
}

func (h *Helper) CalculateTransactionPNL(transaction Transactions) (pnlMessage string, err error) {
	// Calculate total portfolio value and profit/loss
	var totalValue float64
	var totalProfit float64

	// Create a map to store performance by symbol
	type TokenPerformance struct {
		symbol        string
		profitPercent float64
	}
	var performances []TokenPerformance

	// Process all transactions
	for _, trans := range transaction.Result {
		totalValue += trans.CoinData.CurrentValue
		totalProfit += trans.ProfitLoss.Profit

		// Store performance data for sorting
		performances = append(performances, TokenPerformance{
			symbol:        trans.CoinData.Symbol,
			profitPercent: trans.ProfitLoss.ProfitPercent,
		})
	}

	// Calculate total profit percentage
	totalProfitPercent := 0.0
	if totalValue-totalProfit != 0 {
		totalProfitPercent = (totalProfit / (totalValue - totalProfit)) * 100
	}

	// Sort performances by profit percentage (descending)
	sort.Slice(performances, func(i, j int) bool {
		return performances[i].profitPercent > performances[j].profitPercent
	})

	// Format the timestamp
	timestamp := transaction.Result[0].Date.UTC().Format("2006-01-02 15:04 UTC")

	// Build the message
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("📊 Portfolio Update (%s)\n\n", timestamp))
	builder.WriteString(fmt.Sprintf("Portfolio Value: $%.2f\n", totalValue))

	// Add profit/loss line with color indicators
	if totalProfit >= 0 {
		builder.WriteString(fmt.Sprintf("24h P&L: +$%.2f (+%.2f%%)\n", totalProfit, totalProfitPercent))
	} else {
		builder.WriteString(fmt.Sprintf("24h P&L: -$%.2f (%.2f%%)\n", -totalProfit, totalProfitPercent))
	}

	builder.WriteString("\nTop Performers:\n")

	// Add top 3 performers
	for i, perf := range performances {
		if i >= 3 {
			break
		}
		if perf.profitPercent >= 0 {
			builder.WriteString(fmt.Sprintf("$%s: +%.2f%%\n", perf.symbol, perf.profitPercent))
		} else {
			builder.WriteString(fmt.Sprintf("$%s: %.2f%%\n", perf.symbol, perf.profitPercent))
		}
	}

	// Add worst performers section if there are any losses
	hasLosses := false
	for _, perf := range performances {
		if perf.profitPercent < 0 {
			if !hasLosses {
				builder.WriteString("\nWorst Performers:\n")
				hasLosses = true
			}
			builder.WriteString(fmt.Sprintf("$%s: %.2f%%\n", perf.symbol, perf.profitPercent))
			break // Just show the worst performer
		}
	}

	return builder.String(), nil
}

func (h *Helper) CalculatePortfolio(address string) (message string, err error) {
	portfolio, err := h.GetUserPortfolio(address)
	if err != nil {
		log.Println("Error getting user portfolio: ", err.Error())
		return "", err
	}

	var totalValue float64
	message = fmt.Sprintf("📊 **Portfolio Summary for Address:** `%s`\n\n", address)

	for _, asset := range portfolio {
		assetValue := asset.Amount * asset.Price
		totalValue += assetValue

		message += fmt.Sprintf(
			"🔹 **%s (%s)**\n   - Amount: %.6f\n   - Price: $%.6f\n   - Value: $%.2f\n   - 24h Change: %.2f%%\n\n",
			asset.Name,
			asset.Symbol,
			asset.Amount,
			asset.Price,
			assetValue,
			asset.PCh24H,
		)
	}

	message += fmt.Sprintf("💰 **Total Portfolio Value:** $%.2f\n", totalValue)
	return message, nil
}

func (h *Helper) GetAddressPNL(address string) (pnlMessage string, err error) {

	addressStatus, err := h.AddressStatus(address)
	if err != nil {
		log.Println("Error getting address status: ", err.Error())
		return "", err
	}

	if addressStatus.Status != "syncing" {
		addressTransactions, err := h.GetAddressTransactions(address)
		if err != nil {
			log.Println("Error adding address transactions: ", err.Error())
			return "", err
		}

		return h.CalculateTransactionPNL(addressTransactions)
	}

	// Wait for address to sync
	time.Sleep(time.Duration(20) * time.Second)

	// sleep and recall this function
	return h.GetAddressPNL(address)
}
