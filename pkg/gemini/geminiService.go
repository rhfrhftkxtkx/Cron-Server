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
					"title": {
						Type:        genai.TypeString,
						Description: "The title of the exhibition.",
					},
					"summary": {
						Type:        genai.TypeString,
						Description: "A summary of the exhibition's details, with a maximum of 200 characters. If a meaningful summary of 50 characters or less cannot be extracted, the value must be null.",
					},
					"startDate": {
						Type:        genai.TypeString,
						Description: "The start date of the exhibition, formatted as YYYY-MM-DD.",
					},
					"endDate": {
						Type:        genai.TypeString,
						Description: "The end date of the exhibition, formatted as YYYY-MM-DD.",
					},
					"imageUrl": {
						Type:        genai.TypeString,
						Description: "The URL of the exhibition poster. It can be found in an `<img>` tag or a `background-image` style property. The URL must be an absolute path. If it's not an absolute path or cannot be found, the value must be null.",
					},
					"venueName": {
						Type:        genai.TypeString,
						Description: "The official name of the host institution or venue.",
					},
					"venueArea": {
						Type:        genai.TypeString,
						Description: "The location where the exhibition is held, based on the city name. Include more detailed information if available. If not found, the value must be null.",
					},
					"relatedURL": {
						Type:        genai.TypeString,
						Description: "VERY IMPORTANT: Find a URL that likely leads to more detailed information about this specific exhibition. It must be an absolute path. Do not extract URLs from links with text like 'list', 'menu', or other navigational elements. If the URL consists only of a query string (e.g., '?id=123'), construct the absolute URL based on the current URL. The found URL must not be the same as or a parent of the current URL. If no suitable URL is found or the conditions are violated, the value must be null.",
					},
					"depth": {
						Type:        genai.TypeInteger,
						Description: "This is an informational field representing the current crawling depth. Return the provided input value without any changes.",
					},
				},
				Description: "A structured representation of the exhibition information, ",
			},
		},
	}

	geminiClient := &Client{
		Client:     client,
		ModelName:  modelName,
		BasePrompt: "You are an expert AI that extracts exhibition information from a given HTML document.\nYour task is to analyze the provided HTML content and extract the information based on the defined JSON schema.\nPlease adhere to the following rules for each field:\n",
		Config:     config,
	}

	log.Println("[INFO] (gemini.init) Gemini initialization complete.")
	return geminiClient, nil
}

func (geminiClient *Client) Processing(ctx context.Context, url, webPage string, depth int) (*[]Response, error) {
	log.Println("[INFO] (gemini.processing) Processing...")
	prompt := "%s\n- **title**: Extract the main title of the exhibition.\n- **summary**: Provide a concise summary of the exhibition, up to 200 characters. If you can only find a very short or meaningless summary (less than 50 characters), return null.\n- **startDate**: Find the start date and format it as YYYY-MM-DD.\n- **endDate**: Find the end date and format it as YYYY-MM-DD.\n- **imageUrl**: Find the poster image URL. It must be an absolute URL (starts with http or https). Look for it in `<img>` tags or CSS `background-image` properties. If you cannot find an absolute URL, return null.\n- **venueName**: Extract the official name of the venue or organizer.\n- **venueArea**: Identify the city where the exhibition is held. If more specific location details are available, include them. If not found, return null.\n- **relatedURL**: This is the most critical task. Find a single URL that points to a more detailed page for THIS SPECIFIC exhibition.\n  - The URL **must be an absolute path**. If you find a relative path (e.g., `/detail/123`) or just a query string (e.g., `?id=123`), you must combine it with the `currenturl` to create a full, absolute URL.\n  - The URL **must NOT** be for a list, menu, or main page. It should be a deeper, more specific page.\n  - The URL **must NOT** be the same as the `currenturl` or a parent page of the `currenturl`.\n  - If you cannot find a URL that meets all these criteria, you **must** return null.\n- **depth**: This is just an input value. Return the exact same integer value you receive.\n\nNow, process the following data.\n\ncurrenturl: %s,\ndepth: %d,\ndocs: %s"

	content := fmt.Sprintf(prompt, geminiClient.BasePrompt, url, depth, webPage)

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
