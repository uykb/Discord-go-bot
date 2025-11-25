package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"tv-bot-go/internal/analysis"
)

// Service sends technical analysis data to an AI for interpretation.
type Service struct {
	APIKey     string
	Endpoint   string
	HTTPClient *http.Client
}

// NewService creates a new AI service.
func NewService(apiKey, endpoint string) *Service {
	return &Service{
		APIKey:     apiKey,
		Endpoint:   endpoint,
		HTTPClient: &http.Client{},
	}
}

// AnalysisPayload is the data sent to the AI.
type AnalysisPayload struct {
	Symbol      string                           `json:"symbol"`
	Analysis1h  *analysis.TechnicalAnalysis `json:"analysis_1h"`
	Analysis15m *analysis.TechnicalAnalysis `json:"analysis_15m"`
}

// GenerateAnalysis sends the analysis to the AI and returns the interpretation.
func (s *Service) GenerateAnalysis(ctx context.Context, symbol string, analysis1h, analysis15m *analysis.TechnicalAnalysis) (string, error) {
	prompt := BuildPrompt(symbol, analysis1h, analysis15m)

	payload := map[string]interface{}{
		"model":    "deepseek-coder",
		"messages": []map[string]string{{"role": "user", "content": prompt}},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal AI payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.Endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create AI request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call AI API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI API returned non-200 status: %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode AI response: %w", err)
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no analysis returned from AI")
}
