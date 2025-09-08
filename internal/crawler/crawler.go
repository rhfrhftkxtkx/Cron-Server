package crawler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type PageContent struct {
	URL     string
	Content string
	Error   error
}

func CrawlPage(url string, wg *sync.WaitGroup, contentChan chan<- PageContent) {
	defer wg.Done()
	log.Printf("[Info] Crawling page %s", url)

	res, err := http.Get(url)
	if err != nil {
		contentChan <- PageContent{URL: url, Error: fmt.Errorf("[Error] Crawling page %s - %s", url, err.Error())}
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		contentChan <- PageContent{URL: url, Error: fmt.Errorf("[Error] Crawling page %s - status code %d", url, res.StatusCode)}
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		contentChan <- PageContent{URL: url, Error: fmt.Errorf("[Error] Crawling page %s - %s", url, err.Error())}
		return
	}

	content := doc.Find("body").Text()
	cleanContent := strings.Join(strings.Fields(content), " ")

	contentChan <- PageContent{URL: url, Content: cleanContent, Error: nil}
}
