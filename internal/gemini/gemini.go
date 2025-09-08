package gemini

type ProcessedData struct {
	URL      string   `json:"url"`
	Title    string   `json:"title"`
	Summary  string   `json:"summary"`
	Keywords []string `json:"keywords"`
	ImageURL string   `json:"image_url"`
	Area     string   `json:"area"`
}

func Processing(apiKey, originalURL, text string) (*ProcessedData, error) {
	return nil, nil
}
