package searcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type GoogleSearchResult struct {
	Items []struct {
		Link string `json:"link"`
	} `json:"items"`
}

func SearchGoogle(apikey, cx, query, dataRestrict string) ([]string, error) {
	baseUrl := "https://www.googleapis.com/search?"
	params := url.Values{}

	params.Add("key", apikey)
	params.Add("cx", cx)
	params.Add("q", query)
	params.Add("dataRestrict", dataRestrict)

	resp, err := http.Get(fmt.Sprintf("%s?%s", baseUrl, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("[Error] Failed to perform Google search: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[Error] Failed to read Google search response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[Error] Google search API returned status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("[Error] Google search API returned status %d", resp.StatusCode)
	}

	var searchResult GoogleSearchResult
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return nil, fmt.Errorf("[Error] Failed to parse Google search response: %v", err)
	}

	var urls []string
	for _, item := range searchResult.Items {
		urls = append(urls, item.Link)
	}

	log.Printf("[Info] Google search fount %d results", len(urls))
	return urls, nil
}
