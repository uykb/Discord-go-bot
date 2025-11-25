package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
	"tv-bot-go/internal/analysis"
)

const masterPromptTemplate = `
As a crypto market analyst, provide a brief analysis for {{.Symbol}} based on the 1-hour and 15-minute timeframes.
Focus on trend, momentum, and volume. Do not provide any financial advice, trading signals, or price predictions.

1-Hour Analysis:
{{ formatAnalysis .Analysis1h }}

15-Minute Analysis:
{{ formatAnalysis .Analysis15m }}

Synthesize these findings into a short, neutral summary.
`

var tmpl *template.Template

func init() {
	funcs := template.FuncMap{"formatAnalysis": formatAnalysis}
	tmpl = template.Must(template.New("prompt").Funcs(funcs).Parse(masterPromptTemplate))
}

// BuildPrompt creates the final prompt string sent to the AI.
func BuildPrompt(symbol string, analysis1h, analysis15m *analysis.TechnicalAnalysis) string {
	data := map[string]interface{}{
		"Symbol":      symbol,
		"Analysis1h":  analysis1h,
		"Analysis15m": analysis15m,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		// This should ideally not happen if the template is correct
		return fmt.Sprintf("Error executing template: %s", err)
	}

	return buf.String()
}

// formatAnalysis formats the technical analysis data into a human-readable string for the prompt.
func formatAnalysis(analysis *analysis.TechnicalAnalysis) string {
	if analysis == nil {
		return "No data available."
	}
	// Use json marshaling for a quick and readable format.
	b, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return "Error formatting analysis."
	}
	return string(b)
}
