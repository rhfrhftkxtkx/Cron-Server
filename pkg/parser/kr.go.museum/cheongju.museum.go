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

type CheongjuMuseum struct {
}

func (b *CheongjuMuseum) Parsing(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	switch job.Url.Path {
	case "/www/speclExbiView.do":
		return b.parseDetailPage(ctx, cfg, job)
	}

	return nil, nil
}

func (b *CheongjuMuseum) parseDetailPage(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	log.Println("[INFO] (cheongju.parseDetailPage) Parsing detail page ", job.Url)

	if job.Url.Host == "" {
		job.Url.Host = "cheongju.museum.go.kr"
		job.Url.Scheme = "https"
	}

	doc, err := crawler.DoCrawl(job.Url.String())
	if err != nil {
		log.Println("[ERROR] (cheongju.parseDetailPage) Error crawling detail page: ", err)
		return nil, err
	}

	img := doc.Find("span.photo_wrap > img")
	imageURLStr, exists := img.Attr("src")
	if !exists {
		imageURLStr = ""
	}

	imageURL, err := url.Parse(imageURLStr)
	if err != nil {
		imageURL = &url.URL{}
	}
	if imageURL.Host == "" {
		imageURL.Host = "cheongju.museum.go.kr"
		imageURL.Scheme = "https"
	}

	summary := ""
	var sDate, eDate, title string
	doc.Find("div.photo_article > ul > li").Each(func(i int, s *goquery.Selection) {
		label := s.Find("strong").Text()
		value := s.Find("span").Text()

		switch label {
		case "전시명":
			title = value
		case "전시기간":
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
		log.Println("[ERROR] (cheongju.parseExhibitions) Error parsing start date: ", err)
		return nil, err
	}
	if tempDate.Before(time.Now()) {
		log.Println("[INFO] (cheongju.parseExhibitions) Exhibition already ended, skipping")
	}

	return &parser.ParseResult{
		FoundExhibitions: []*common.Exhibition{
			{
				ExhibitionId:     common.GenerateExhibitionId("cheongjumuseum", sDate, title),
				VenueVisitKor2Id: "129779",
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

func GetCheongjuMuseum() parser.MuseumPageParser {
	return &CheongjuMuseum{}
}
