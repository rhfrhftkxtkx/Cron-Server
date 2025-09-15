package crawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/BlueNyang/theday-theplace-cron/internal/domain/common"
	"github.com/BlueNyang/theday-theplace-cron/internal/parser"
	"github.com/PuerkitoBio/goquery"
)

type CrawlResult struct {
	URL         string
	Exhibitions []common.Exhibition
	Error       error
}

func CrawlAndParse(url string, parser parser.Parser, wg *sync.WaitGroup, resultChan chan<- CrawlResult) {
	// Signal that this goroutine is done when the function exits
	defer wg.Done()
	log.Printf("[Info] Crawling page %s", url)

	res, err := http.Get(url)
	if err != nil {
		// Send error to channel and return
		resultChan <- CrawlResult{URL: url, Error: fmt.Errorf("[Error] Crawling page %s - %s", url, err.Error())}
		return
	}

	// Ensure the response body is closed after reading
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("[Warning] Failed to close response body for %s: %s", url, err.Error())
		}
	}(res.Body)

	if res.StatusCode != 200 {
		resultChan <- CrawlResult{URL: url, Error: fmt.Errorf("[Error] Crawling page %s - status code %d", url, res.StatusCode)}
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		resultChan <- CrawlResult{URL: url, Error: fmt.Errorf("[Error] Crawling page %s - %s", url, err.Error())}
		return
	}

	exhibitions, err := parser.Parse(doc)
	if err != nil {
		resultChan <- CrawlResult{URL: url, Error: fmt.Errorf("[Error] Parsing page %s - %s", url, err.Error())}
		return
	}

	resultChan <- CrawlResult{URL: url, Exhibitions: exhibitions, Error: nil}
}
