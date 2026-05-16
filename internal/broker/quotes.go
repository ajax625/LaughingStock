package broker

import (
	"fmt"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

// GetLatestPrices fetches the latest trade price for each symbol using the user's Alpaca credentials.
func GetLatestPrices(alpacaKey, alpacaSecret string, symbols []string) (map[string]float64, error) {
	if len(symbols) == 0 {
		return map[string]float64{}, nil
	}

	client := marketdata.NewClient(marketdata.ClientOpts{
		APIKey:    alpacaKey,
		APISecret: alpacaSecret,
	})

	snapshots, err := client.GetSnapshots(symbols, marketdata.GetSnapshotRequest{
		Feed: marketdata.IEX,
	})
	if err != nil {
		return nil, fmt.Errorf("alpaca marketdata: get snapshots failed: %w", err)
	}

	prices := make(map[string]float64, len(snapshots))
	for sym, snap := range snapshots {
		if snap != nil && snap.LatestTrade != nil {
			prices[sym] = snap.LatestTrade.Price
		}
	}
	return prices, nil
}
