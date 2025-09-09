package searcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type NextPageIndex struct {
	StartIndex int
}

type GoogleSearchResult struct {
	Queries struct {
		NextPage []NextPageIndex `json:"nextPage"`
	} `json:"queries"`

	Items []struct {
		Link string `json:"link"`
	} `json:"items"`
}

type GoogleSearchConfig struct {
	APIKey       string
	CX           string
	Query        string
	Language     string
	OrTerms      string
	Safe         string
	DateRestrict string
	Num          int
	SiteSearch   string
	SiteFilter   string
	Filter       string
}

func SearchGoogle(cfg GoogleSearchConfig) (GoogleSearchResult, error) {
	baseUrl, _ := url.Parse("https://www.googleapis.com/customsearch/v1")

	params := url.Values{}

	params.Add("key", cfg.APIKey)
	params.Add("cx", cfg.CX)
	params.Add("q", cfg.Query)
	params.Add("lr", cfg.Language)
	params.Add("orTerms", cfg.OrTerms)
	params.Add("safe", cfg.Safe)
	params.Add("dateRestrict", cfg.DateRestrict)
	params.Add("num", fmt.Sprintf("%d", cfg.Num))
	params.Add("siteSearch", cfg.SiteSearch)
	params.Add("siteSearchFilter", cfg.SiteFilter)
	params.Add("filter", cfg.Filter)

	baseUrl.RawQuery = params.Encode()

	resp, err := http.Get(baseUrl.String())
	if err != nil {
		return nil, fmt.Errorf("[Error] Failed to make Google search request: %v", err)
	}

	defer TryCloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("[Error] Google search API returned status %d", resp.StatusCode)
		return nil, fmt.Errorf("[Error] Google search API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[Error] Failed to read Google search response: %v", err)
	}

	var searchResult GoogleSearchResult
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return nil, fmt.Errorf("[Error] Failed to parse Google search response: %v", err)
	}

	return searchResult, nil
}

func TryCloseBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Printf("[Error] Failed to close response body: %v", err)
	}
}
