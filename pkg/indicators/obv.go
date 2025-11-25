package indicators

// CalculateOBV calculates the On-Balance Volume (OBV).
// It requires closing prices and volumes.
func CalculateOBV(closes, volumes []float64) []float64 {
	if len(closes) != len(volumes) || len(closes) == 0 {
		return nil
	}

	obv := make([]float64, len(closes))
	obv[0] = 0 // OBV starts at 0

	for i := 1; i < len(closes); i++ {
		if closes[i] > closes[i-1] {
			obv[i] = obv[i-1] + volumes[i] // Price went up, add volume
		} else if closes[i] < closes[i-1] {
			obv[i] = obv[i-1] - volumes[i] // Price went down, subtract volume
		} else {
			obv[i] = obv[i-1] // Price is unchanged, OBV is unchanged
		}
	}

	return obv
}
