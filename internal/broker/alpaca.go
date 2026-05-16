package broker

import (
	"context"
	"fmt"

	alpaca "github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/shopspring/decimal"
)

type OrderResult struct {
	ID     string
	Symbol string
	Side   string
	Qty    float64
	Status string
}

// PlaceMarketOrder places a paper market order using the user's Alpaca credentials.
func PlaceMarketOrder(ctx context.Context, alpacaKey, alpacaSecret, symbol, side string, qty float64) (*OrderResult, error) {
	client := alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    alpacaKey,
		APISecret: alpacaSecret,
		BaseURL:   "https://paper-api.alpaca.markets",
	})

	alpacaSide := alpaca.Buy
	if side == "sell" {
		alpacaSide = alpaca.Sell
	}

	qtyDecimal := decimal.NewFromFloat(qty)
	order, err := client.PlaceOrder(alpaca.PlaceOrderRequest{
		Symbol:      symbol,
		Qty:         &qtyDecimal,
		Side:        alpacaSide,
		Type:        alpaca.Market,
		TimeInForce: alpaca.Day,
	})
	if err != nil {
		return nil, fmt.Errorf("alpaca: place order failed: %w", err)
	}

	return &OrderResult{
		ID:     order.ID,
		Symbol: order.Symbol,
		Side:   string(order.Side),
		Qty:    qty,
		Status: string(order.Status),
	}, nil
}
