package market

import (
	"context"
	"tv-bot-go/internal/binance"
)

// Service provides market data and analysis.
type Service struct {
	binanceClient *binance.Client
}

// NewService creates a new market service.
func NewService(binanceClient *binance.Client) *Service {
	return &Service{binanceClient: binanceClient}
}

// MarketData holds the raw kline data for different timeframes.
type MarketData struct {
	Klines1h  []binance.Kline
	Klines15m []binance.Kline
}

// FetchMarketData fetches kline data for a given symbol for 1h and 15m intervals.
func (s *Service) FetchMarketData(ctx context.Context, symbol string) (*MarketData, error) {
	// Fetch 1-hour klines
	klines1h, err := s.binanceClient.GetKlines(ctx, symbol, "1h", 100)
	if err != nil {
		return nil, err
	}

	// Fetch 15-minute klines
	klines15m, err := s.binanceClient.GetKlines(ctx, symbol, "15m", 100)
	if err != nil {
		return nil, err
	}

	return &MarketData{
		Klines1h:  klines1h,
		Klines15m: klines15m,
	}, nil
}
