package crawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func DoCrawl(url string) (*goquery.Document, error) {
	log.Printf("[INFO] (crawler.DoCrawl) Crawling page %s", url)

	// replace &amp; with & using string replacement
	url = strings.ReplaceAll(url, "&amp;", "&")

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] (crawler.DoCrawl) Crawling page %s - %s", url, err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("[Warning] (crawler.DoCrawl) Failed to close response body for %s: %s", url, err.Error())
		}
	}(res.Body)

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] (crawler.DoCrawl) Crawling page %s - status code %d", url, res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] (crawler.DoCrawl) Crawling page %s - %s", url, err.Error())
	}

	return doc, nil
}
