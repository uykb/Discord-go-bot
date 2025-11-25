package indicators

// MACDIndicator holds the calculated values for the MACD indicator.
type MACDIndicator struct {
	MACD      float64 `json:"macd"`
	Signal    float64 `json:"signal"`
	Histogram float64 `json:"histogram"`
}

// calculateEMA calculates the Exponential Moving Average.
func calculateEMA(data []float64, period int) []float64 {
	if len(data) < period {
		return nil
	}

	ema := make([]float64, len(data))
	k := 2.0 / float64(period+1)

	// Calculate initial SMA
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += data[i]
	}
	ema[period-1] = sum / float64(period)

	// Calculate EMA for the rest of the data
	for i := period; i < len(data); i++ {
		ema[i] = (data[i] * k) + (ema[i-1] * (1 - k))
	}

	return ema
}

// CalculateMACD calculates the MACD indicator values.
func CalculateMACD(prices []float64, fastPeriod, slowPeriod, signalPeriod int) []MACDIndicator {
	if len(prices) < slowPeriod {
		return nil
	}

	emaFast := calculateEMA(prices, fastPeriod)
	emaSlow := calculateEMA(prices, slowPeriod)

	if emaFast == nil || emaSlow == nil {
		return nil
	}

	macdLine := make([]float64, len(prices))
	for i := slowPeriod - 1; i < len(prices); i++ {
		macdLine[i] = emaFast[i] - emaSlow[i]
	}

	// We need to pass the valid part of the macdLine to calculate the signal line
	signalLine := calculateEMA(macdLine[slowPeriod-1:], signalPeriod)
	if signalLine == nil {
		return nil
	}

	results := make([]MACDIndicator, 0)
	// The starting point for the histogram is when the signal line calculation is complete
	histogramStartIndex := slowPeriod - 1 + signalPeriod - 1

	for i := histogramStartIndex; i < len(prices); i++ {
		macdVal := macdLine[i]
		// Adjust index for signalLine as it was calculated on a slice of macdLine
		signalVal := signalLine[i-(slowPeriod-1)]
		histogramVal := macdVal - signalVal

		results = append(results, MACDIndicator{
			MACD:      macdVal,
			Signal:    signalVal,
			Histogram: histogramVal,
		})
	}

	return results
}
