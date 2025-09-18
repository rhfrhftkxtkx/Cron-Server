package kr_go_museum

import (
	"context"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/BlueNyang/theday-theplace-cron/pkg/config"
	"github.com/BlueNyang/theday-theplace-cron/pkg/crawler"
	"github.com/BlueNyang/theday-theplace-cron/pkg/domain/common"
	"github.com/BlueNyang/theday-theplace-cron/pkg/parser"
	"github.com/PuerkitoBio/goquery"
)

type ChuncheonMuseum struct {
}

func (b *ChuncheonMuseum) Parsing(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	switch job.Url.Path {
	case "/prog/spclExht/kor/sub02_03/view.do":
		return b.parseDetailPage(ctx, cfg, job)
	}

	return nil, nil
}

func (b *ChuncheonMuseum) parseDetailPage(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	log.Println("[INFO] (chuncheon.parseDetailPage) Parsing detail page ", job.Url)

	if job.Url.Host == "" {
		job.Url.Host = "chuncheon.museum.go.kr"
		job.Url.Scheme = "https"
	}

	doc, err := crawler.DoCrawl(job.Url.String())
	if err != nil {
		log.Println("[ERROR] (chuncheon.parseDetailPage) Error crawling detail page: ", err)
		return nil, err
	}

	img := doc.Find("img.card-img-top")
	imageURLStr, exists := img.Attr("src")
	if !exists {
		imageURLStr = ""
	}

	imageURL, err := url.Parse(imageURLStr)
	if err != nil {
		imageURL = &url.URL{}
	}
	if imageURL.Host == "" {
		imageURL.Host = "chuncheon.museum.go.kr"
		imageURL.Scheme = "https"
	}

	title := doc.Find("strong.title > em").Text()
	summary := doc.Find("div#boxscroll > div > p:first-child").Text()
	var sDate, eDate string
	doc.Find("ul.list-1st > li").Each(func(i int, s *goquery.Selection) {
		label := s.Find("em").Text()
		value := s.Find("em").Remove().End().Text()

		if label == "기간" {
			dates := value
			if len(dates) >= 23 {
				dates = strings.TrimSpace(dates)
				dates = strings.ReplaceAll(dates, ".", "-")
				dates = strings.ReplaceAll(dates, " ~ ", "~")
				sDate = dates[0:10]
				eDate = dates[11:21]
			}
		}
	})

	tempDate, err := time.Parse("2006-01-02", sDate)
	if err != nil {
		log.Println("[ERROR] (chuncheon.parseExhibitions) Error parsing start date: ", err)
		return nil, err
	}
	if tempDate.Before(time.Now()) {
		log.Println("[INFO] (chuncheon.parseExhibitions) Exhibition already ended, skipping")
	}

	return &parser.ParseResult{
		FoundExhibitions: []*common.Exhibition{
			{
				ExhibitionId:     common.GenerateExhibitionId("chuncheonmuseum", sDate, title),
				VenueVisitKor2Id: "130654",
				DataSourceTier:   1,
				Title:            title,
				Summary:          summary,
				StartDate:        sDate,
				EndDate:          eDate,
				ImageUrl:         imageURL.String(),
				SourceURL:        job.Url.String(),
			},
		},
		DiscoveredJobs: []*parser.Job{},
	}, nil
}

func GetChuncheonMuseum() parser.MuseumPageParser {
	return &ChuncheonMuseum{}
}
