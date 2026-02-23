package service

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

type AIService struct {
	Client *genai.Client
}

func NewAIService(ctx context.Context) (*AIService, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &AIService{Client: client}, nil
}

func (s *AIService) AnalyzeFinancialHealth(ctx context.Context, transactionHistory string) (string, error) {
	prompt := fmt.Sprintf(`
		Kamu adalah asisten keuangan pribadi yang bijak.
		Berikut data transaksi user:
		%s
		
		Tugas:
		1. Analisa apakah user boros/hemat.
		2. Beri 1 saran keuangan praktis.
		Jawab singkat maksimal 3 paragraf, gaya bahasa santai tapi sopan.
	`, transactionHistory)

	resp, err := s.Client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("gagal generate content: %v", err)
	}

	if resp != nil && len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		part := resp.Candidates[0].Content.Parts[0]
		return fmt.Sprintf("%v", part.Text), nil
	}

	return "Maaf, AI tidak memberikan saran.", nil
	
}
