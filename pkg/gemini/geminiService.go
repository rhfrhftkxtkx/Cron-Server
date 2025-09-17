package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/genai"
)

type Response struct {
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
	ImageURL   string `json:"imageUrl"`
	VenueName  string `json:"venueName"`
	VenueArea  string `json:"venueArea"`
	RelatedURL string `json:"relatedURL"`
	Depth      int    `json:"depth"`
}

type Client struct {
	Client     *genai.Client
	ModelName  string
	BasePrompt string
	Config     *genai.GenerateContentConfig
}

func InitGemini(ctx context.Context, apiKey, modelName string) (*Client, error) {
	log.Println("[INFO] (gemini.init) Creating Gemini client...")
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("InitGemini failed: %v", err)
	}

	log.Println("[INFO] (gemini.init) Done. Initializing Gemini...")
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"title":      {Type: genai.TypeString, Description: "An officially registered title of the exhibition information"},
					"summary":    {Type: genai.TypeString, Description: "A brief summary of the overview max 500 characters"},
					"startDate":  {Type: genai.TypeString, Description: "Start date of the exhibition in YYYY-MM-DD format"},
					"endDate":    {Type: genai.TypeString, Description: "End date of the exhibition in YYYY-MM-DD format"},
					"imageUrl":   {Type: genai.TypeString, Description: "URL of a representative image"},
					"venueName":  {Type: genai.TypeString, Description: "An officially registered name of the event organizer"},
					"venueArea":  {Type: genai.TypeString, Description: "Area of the organization holding the exhibition"},
					"relatedURL": {Type: genai.TypeString, Description: "**An absolute URL** for the detailed exhibition page. If the extracted href is a relative path (e.g., '/path', './path'), it **must be combined with the base of 'currentUrl'** to form a complete, absolute URL. Null if **no separate detailed page exists** or if the found URL is functionally identical to the 'currenturl'."},
					"depth":      {Type: genai.TypeInteger, Description: "return the current depth"},
				},
				Description: "A structured representation of the exhibition information, ",
			},
		},
	}

	geminiClient := &Client{
		Client:     client,
		ModelName:  modelName,
		BasePrompt: "Extract the title, summary, period (startDate, endDate), imageUrl, venueName, venueArea, relatedURL, and depth from the following web page content about exhibition information. Provide the output in JSON array format with the specified fields. \ndata: \n",
		Config:     config,
	}

	log.Println("[INFO] (gemini.init) Gemini initialization complete.")
	return geminiClient, nil
}

func (geminiClient *Client) Processing(ctx context.Context, url, webPage string, depth int) (*[]Response, error) {
	log.Println("[INFO] (gemini.processing) Processing...")
	prompt := `%s
---
## **Overall Data Extraction Strategy**
1.  **If the current page contains detailed information,** extract data for all fields.
2.  **If the current page is a list or summary page with insufficient data** (e.g., the summary or dates are missing), you must follow these steps:
- Set fields with missing information to an empty string.
- You **MUST** find and include the 'relatedURL' that leads to the detail page. In this situation, 'relatedURL' cannot be an empty string or "null".
---
## **'relatedURL' Formatting Rules**
1. The 'relatedURL' must be a full, absolute URL starting with 'http://' or 'https://'.
2. If a relative path is found (e.g., starts with '/'), prepend the protocol and domain from the 'currenturl'.
3. If only a query string is found (e.g., starts with '?'), prepend the 'currenturl' up to its path.
4. If a URL is missing its protocol (e.g., 'www.example.com/page'), add 'https://'.
5. If the page contains sufficient data (as per the Overall Strategy), and no link to a more detailed page exists (or the link is for the current page), the value for 'relatedURL' **MUST** be an empty string.
---
## **Examples**
- **On a detail page:** If 'currenturl' is 'https://example.com/events/detail?id=123' and all data is present, fill all fields and set 'relatedURL' to an empty string.
- **On a list page:** If 'currenturl' is 'https://example.com/events/list' and it only contains titles and thumbnails for each item, you must set missing fields like 'summary' and 'startDate' to an empty string and provide a valid, absolute 'relatedURL' to its detail page.
---
data:
  current depth: %d
  current URL: %s
  document: %s
`
	content := fmt.Sprintf(prompt, geminiClient.BasePrompt, depth, url, webPage)

	result, err := geminiClient.Client.Models.GenerateContent(
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
	var resp []Response
	// Assume result.Content contains the JSON response
	err = json.Unmarshal([]byte(result.Text()), &resp)
	if err != nil {
		log.Fatalf("[Error] Failed to parse Gemini response: %v", err)
		return nil, err
	}

	log.Println("[INFO] (gemini.processing) Processing complete.")
	return &resp, nil
}
