package kr_go_museum

import (
	"context"
	"log"
	"regexp"

	"github.com/BlueNyang/theday-theplace-cron/pkg/config"
	"github.com/BlueNyang/theday-theplace-cron/pkg/parser"
)

type Buyeo struct {
}

func (b Buyeo) Parsing(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	urlDetailPattern := regexp.MustCompile(`https://buyeo.museum.go.kt/speclExhibi/view.do\?.*`)
	if urlDetailPattern.MatchString(*job.Url) {
		return parseDetailPage(ctx, cfg, job)
	}
	return nil, nil
}

func parseDetailPage(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	log.Println("[INFO] (buyeo.parseDetailPage) Parsing detail page ", job.Url)

	//baseURL, err := url.Parse(*job.Url)
	//if err != nil {
	//	log.Println("[ERROR] (buyeo.parseDetailPage) Error parsing URL: ", err)
	//	return nil, err
	//}
	//
	//domain = baseURL.Hostname()

	//doc, err := crawler.DoCrawl(*job.Url)
	//if err != nil {
	//	log.Println("[ERROR] (buyeo.parseDetailPage) Error crawling detail page: ", err)
	//	return nil, err
	//}
	//
	//img := doc.Find(".swiper-slide > img")

	return nil, nil
}

func GetBuyeo() parser.MuseumPageParser {
	return &Buyeo{}
}
