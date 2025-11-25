package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultBaseURL = "https://api.binance.com"
	klinesEndpoint = "/api/v3/klines"
)

// Client is a Binance API client.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Binance API client.
func NewClient() *Client {
	return &Client{
		BaseURL:    defaultBaseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Kline represents a single kline/candlestick.
type Kline struct {
	OpenTime                 int64
	Open, High, Low, Close   float64
	Volume                   float64
	CloseTime                int64
	QuoteAssetVolume         float64
	NumberOfTrades           int64
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
}

// UnmarshalJSON implements a custom unmarshaler for the Kline struct.
// Binance API returns kline data as a JSON array, not an object.
func (k *Kline) UnmarshalJSON(data []byte) error {
	var raw []interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal kline data: %w", err)
	}

	if len(raw) < 11 {
		return fmt.Errorf("invalid kline data length: expected at least 11, got %d", len(raw))
	}

	var err error
	k.OpenTime = int64(raw[0].(float64))
	k.Open, err = strconv.ParseFloat(raw[1].(string), 64)
	if err != nil { return err }
	k.High, err = strconv.ParseFloat(raw[2].(string), 64)
	if err != nil { return err }
	k.Low, err = strconv.ParseFloat(raw[3].(string), 64)
	if err != nil { return err }
	k.Close, err = strconv.ParseFloat(raw[4].(string), 64)
	if err != nil { return err }
	k.Volume, err = strconv.ParseFloat(raw[5].(string), 64)
	if err != nil { return err }
	k.CloseTime = int64(raw[6].(float64))
	k.QuoteAssetVolume, err = strconv.ParseFloat(raw[7].(string), 64)
	if err != nil { return err }
	k.NumberOfTrades = int64(raw[8].(float64))
	k.TakerBuyBaseAssetVolume, err = strconv.ParseFloat(raw[9].(string), 64)
	if err != nil { return err }
	k.TakerBuyQuoteAssetVolume, err = strconv.ParseFloat(raw[10].(string), 64)
	if err != nil { return err }

	return nil
}

// GetKlines fetches kline/candlestick data for a symbol.
func (c *Client) GetKlines(ctx context.Context, symbol, interval string, limit int) ([]Kline, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+klinesEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("interval", interval)
	q.Add("limit", strconv.Itoa(limit))
	req.URL.RawQuery = q.Encode()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var klines []Kline
	if err := json.NewDecoder(resp.Body).Decode(&klines); err != nil {
		return nil, fmt.Errorf("failed to decode klines response: %w", err)
	}

	return klines, nil
}
