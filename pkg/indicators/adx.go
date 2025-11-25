package indicators

import "math"

// CalculateADX calculates the Average Directional Index (ADX).
// It requires high, low, and close prices, and a period for the calculation.
func CalculateADX(highs, lows, closes []float64, period int) []float64 {
	if len(highs) < period*2-1 || len(lows) < period*2-1 || len(closes) < period*2-1 {
		return nil
	}

	trueRanges := calculateTrueRange(highs, lows, closes)
	diPlus, diMinus := calculateDirectionalIndicators(highs, lows, closes, trueRanges, period)

	dx := make([]float64, len(diPlus))
	for i := range diPlus {
		if diPlus[i]+diMinus[i] != 0 {
			dx[i] = 100 * math.Abs(diPlus[i]-diMinus[i]) / (diPlus[i] + diMinus[i])
		}
	}

	// Smooth the DX to get ADX. We start from where DX is available.
	adx := smooth(dx[period-1:], period)

	return adx
}

// calculateTrueRange calculates the True Range for each period.
func calculateTrueRange(highs, lows, closes []float64) []float64 {
	tr := make([]float64, len(highs))
	for i := 1; i < len(highs); i++ {
		highLow := highs[i] - lows[i]
		highClose := math.Abs(highs[i] - closes[i-1])
		lowClose := math.Abs(lows[i] - closes[i-1])
		tr[i] = math.Max(highLow, math.Max(highClose, lowClose))
	}
	return tr
}

// calculateDirectionalIndicators calculates the +DI and -DI values.
func calculateDirectionalIndicators(highs, lows, closes, trueRanges []float64, period int) ([]float64, []float64) {
	dmPlus := make([]float64, len(highs))
	dmMinus := make([]float64, len(highs))

	for i := 1; i < len(highs); i++ {
		upMove := highs[i] - highs[i-1]
		downMove := lows[i-1] - lows[i]

		if upMove > downMove && upMove > 0 {
			dmPlus[i] = upMove
		}
		if downMove > upMove && downMove > 0 {
			dmMinus[i] = downMove
		}
	}

	smoothedDMPlus := smooth(dmPlus, period)
	smoothedDMMinus := smooth(dmMinus, period)
	smoothedTR := smooth(trueRanges, period)

	diPlus := make([]float64, len(smoothedTR))
	diMinus := make([]float64, len(smoothedTR))

	for i := range smoothedTR {
		if smoothedTR[i] != 0 {
			diPlus[i] = 100 * smoothedDMPlus[i] / smoothedTR[i]
			diMinus[i] = 100 * smoothedDMMinus[i] / smoothedTR[i]
		}
	}

	return diPlus, diMinus
}

// smooth applies a Wilder's smoothing method.
func smooth(data []float64, period int) []float64 {
	if len(data) < period {
		return nil
	}
	smoothed := make([]float64, len(data))

	// Initial sum
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += data[i]
	}
	smoothed[period-1] = sum

	// Subsequent smoothing
	for i := period; i < len(data); i++ {
		smoothed[i] = smoothed[i-1] - (smoothed[i-1] / float64(period)) + data[i]
	}

	// The first 'period-1' values are not valid, so we return the valid part
	return smoothed
}
