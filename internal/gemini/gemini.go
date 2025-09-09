package gemini

import (
	"context"
	"encoding/json"
	"log"

	"google.golang.org/genai"
)

type ProcessedData struct {
	URL      string   `json:"url"`
	Title    string   `json:"title"`
	Summary  string   `json:"summary"`
	Keywords []string `json:"keywords"`
	ImageURL string   `json:"image_url"`
	Area     string   `json:"area"`
}

type GeminiClient struct {
	Client     *genai.Client
	ModelName  string
	BasePrompt string
	Config     *genai.GenerateContentConfig
}

func InitGemini(ctx context.Context, apiKey, modelName string) *GeminiClient {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"url":       {Type: genai.TypeString, Description: "The URL of the article"},
				"title":     {Type: genai.TypeString, Description: "The title of the article"},
				"summary":   {Type: genai.TypeString, Description: "A brief summary of the article"},
				"keywords":  {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "A list of relevant keywords"},
				"image_url": {Type: genai.TypeString, Description: "URL of a representative image"},
				"area":      {Type: genai.TypeString, Description: "The area or category of the article"},
			},
			PropertyOrdering: []string{"title", "area", "summary"},
		},
	}

	return &GeminiClient{
		Client:     client,
		ModelName:  modelName,
		BasePrompt: "Extract the title, summary, keywords, image URL, and area from the following content:\n",
		Config:     config,
	}
}

func (geminiClient *GeminiClient) Processing(ctx context.Context, webPage string) (*ProcessedData, error) {
	client := geminiClient.Client

	content := geminiClient.BasePrompt + webPage

	result, err := client.Models.GenerateContent(
		ctx,
		geminiClient.ModelName,
		genai.Text(content),
		geminiClient.Config,
	)
	if err != nil {
		log.Fatalf("[Error] Gemini API error: %v", err)
		return nil, err
	}

	// Parse the response to extract the structured data
	var processedData ProcessedData
	// Assume result.Content contains the JSON response
	err = json.Unmarshal([]byte(result.Text()), &processedData)
	if err != nil {
		log.Fatalf("[Error] Failed to parse Gemini response: %v", err)
		return nil, err
	}

	return &processedData, nil
}
