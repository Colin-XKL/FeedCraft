package adapter

import (
	"FeedCraft/internal/util"
	"fmt"
	"github.com/go-resty/resty/v2"
)

const geminiBaseUrl = "https://generativelanguage.googleapis.com/v1/models/"

type GeminiConf struct {
	baseUrl      string
	secretKey    string
	defaultModel string
}

type GeminiClient struct {
	conf   *GeminiConf
	client *resty.Client
}

func NewGeminiClient() (*GeminiClient, error) {
	envClient := util.GetEnvClient()
	geminiKey := envClient.GetString("GEMINI_SECRET_KEY")
	if geminiKey == "" {
		return nil, fmt.Errorf("GEMINI_SECRET_KEY not found in env")
	}

	client := resty.New().SetBaseURL(geminiBaseUrl)
	headers := map[string]string{
		"Cache-Control": "no-cache",
		"Content-Type":  "application/json",
	}
	client.SetBaseURL(geminiBaseUrl).SetQueryParam("key", geminiKey).SetHeaders(headers)

	return &GeminiClient{
		conf:   &GeminiConf{baseUrl: geminiBaseUrl, secretKey: geminiKey, defaultModel: "gemini-pro"},
		client: client,
	}, nil
}

func (gemini *GeminiClient) GenerateContent(payload GeminiReqPayload) (string, error) {
	apiPath := fmt.Sprintf("%s:generateContent", gemini.conf.defaultModel)
	result := &GeminiResp{}

	resp, err := gemini.client.R().SetBody(payload).Post(apiPath)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("GeminiResp status code not 200: %d", resp.StatusCode())
	}
	if (len(result.Candidates) == 0) || result.Candidates[0].Content == nil || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("unexpected empty result")
	}
	return *result.Candidates[0].Content.Parts[0].Text, nil
}

type GeminiReqPayload struct {
	Contents []Content `json:"contents"`
}

type GeminiResp struct {
	Candidates []Candidate `json:"candidates,omitempty"`
}

type Candidate struct {
	Content       *Content       `json:"content,omitempty"`
	FinishReason  *string        `json:"finishReason,omitempty"`
	Index         *int64         `json:"index,omitempty"`
	SafetyRatings []SafetyRating `json:"safetyRatings,omitempty"`
}

type Content struct {
	Parts []Part  `json:"parts,omitempty"`
	Role  *string `json:"role,omitempty"`
}

type Part struct {
	Text *string `json:"text,omitempty"`
}

type SafetyRating struct {
	Category    *string `json:"category,omitempty"`
	Probability *string `json:"probability,omitempty"`
}
