package indicators

// CalculateMFI calculates the Money Flow Index (MFI).
// It requires high, low, close prices, and volumes, along with a period.
func CalculateMFI(highs, lows, closes, volumes []float64, period int) []float64 {
	if len(highs) <= period {
		return nil
	}

	typicalPrices := make([]float64, len(highs))
	rawMoneyFlows := make([]float64, len(highs))

	for i := 0; i < len(highs); i++ {
		typicalPrices[i] = (highs[i] + lows[i] + closes[i]) / 3
		rawMoneyFlows[i] = typicalPrices[i] * volumes[i]
	}

	positiveMoneyFlows := make([]float64, len(highs))
	negativeMoneyFlows := make([]float64, len(highs))

	for i := 1; i < len(highs); i++ {
		if typicalPrices[i] > typicalPrices[i-1] {
			positiveMoneyFlows[i] = rawMoneyFlows[i]
		} else if typicalPrices[i] < typicalPrices[i-1] {
			negativeMoneyFlows[i] = rawMoneyFlows[i]
		}
	}

	mfiValues := make([]float64, len(highs)-period)

	for i := period; i < len(highs); i++ {
		sumPositiveMF := 0.0
		sumNegativeMF := 0.0

		for j := i - period + 1; j <= i; j++ {
			sumPositiveMF += positiveMoneyFlows[j]
			sumNegativeMF += negativeMoneyFlows[j]
		}

		if sumNegativeMF == 0 {
			// Avoid division by zero; if negative flow is zero, MFI is 100.
			mfiValues[i-period] = 100
			continue
		}

		moneyFlowRatio := sumPositiveMF / sumNegativeMF
		mfi := 100 - (100 / (1 + moneyFlowRatio))
		mfiValues[i-period] = mfi
	}

	return mfiValues
}
