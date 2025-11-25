package analysis

import (
	"tv-bot-go/internal/binance"
	"tv-bot-go/pkg/indicators"
)

// TechnicalAnalysis holds the calculated indicators for a specific timeframe.
type TechnicalAnalysis struct {
	Timeframe   string                    `json:"timeframe"`
	MACD        *indicators.MACDIndicator `json:"macd,omitempty"`
	ADX         float64                   `json:"adx,omitempty"`
	OBV         float64                   `json:"obv,omitempty"`
	MFI         float64                   `json:"mfi,omitempty"`
}

// Service performs technical analysis on market data.
type Service struct{}

// NewService creates a new analysis service.
func NewService() *Service {
	return &Service{}
}

// AnalyzeKlines performs a full technical analysis on a slice of klines.
func (s *Service) AnalyzeKlines(klines []binance.Kline, timeframe string) *TechnicalAnalysis {
	if len(klines) == 0 {
		return nil
	}

	closes := getSlice(klines, "close")
	highs := getSlice(klines, "high")
	lows := getSlice(klines, "low")
	volumes := getSlice(klines, "volume")

	analysis := &TechnicalAnalysis{Timeframe: timeframe}

	// Calculate MACD
	if macdResult := indicators.CalculateMACD(closes, 12, 26, 9); len(macdResult) > 0 {
		analysis.MACD = &macdResult[len(macdResult)-1]
	}

	// Calculate ADX
	if adxResult := indicators.CalculateADX(highs, lows, closes, 14); len(adxResult) > 0 {
		analysis.ADX = adxResult[len(adxResult)-1]
	}

	// Calculate OBV
	if obvResult := indicators.CalculateOBV(closes, volumes); len(obvResult) > 0 {
		analysis.OBV = obvResult[len(obvResult)-1]
	}

	// Calculate MFI
	if mfiResult := indicators.CalculateMFI(highs, lows, closes, volumes, 14); len(mfiResult) > 0 {
		analysis.MFI = mfiResult[len(mfiResult)-1]
	}

	return analysis
}

// getSlice extracts a slice of float64 values from klines for a given field.
func getSlice(klines []binance.Kline, field string) []float64 {
	slice := make([]float64, len(klines))
	for i, k := range klines {
		switch field {
		case "close":
			slice[i] = k.Close
		case "high":
			slice[i] = k.High
		case "low":
			slice[i] = k.Low
		case "volume":
			slice[i] = k.Volume
		}
	}
	return slice
}
