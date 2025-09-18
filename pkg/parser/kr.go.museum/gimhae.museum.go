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

type GimhaeMuseum struct {
}

func (g *GimhaeMuseum) Parsing(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	switch job.Url.Path {
	case "/kr/html/sub02/020201.html":
		return g.parseExhibitions(ctx, cfg, job)
	default:
		return nil, nil
	}
}

func (g *GimhaeMuseum) parseExhibitions(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	log.Println("[INFO] (gimhae.parseDetailPage) Parsing detail page ", job.Url)

	if job.Url.Host == "" {
		job.Url.Host = "gimhae.museum.go.kr"
		job.Url.Scheme = "https"
	}

	doc, err := crawler.DoCrawl(job.Url.String())
	if err != nil {
		log.Println("[ERROR] (gimhae.parseDetailPage) Error crawling detail page: ", err)
		return nil, err
	}

	img := doc.Find("a.open_layer > img")
	imageURLstr, exists := img.Attr("src")
	if !exists {
		imageURLstr = ""
	}

	imageURL, err := url.Parse(imageURLstr)
	if err != nil {
		imageURL = &url.URL{}
	}

	if imageURL.Host == "" {
		imageURL.Host = "gimhae.museum.go.kr"
		imageURL.Scheme = "https"
	}

	title := doc.Find("div.tit").Find("span").Remove().End().Text()
	var sDate, eDate, summary string
	doc.Find("div.list > ul > li").Each(func(i int, s *goquery.Selection) {
		label := s.Find("em").Text()
		value := s.Find("span").Text()

		if label == "기간" {
			dates := value
			if len(dates) >= 23 {
				dates = strings.TrimSpace(dates)
				dates = strings.ReplaceAll(dates, ".", "-")
				dates = strings.ReplaceAll(dates, " ~ ", "~")
				sDate = dates[0:10]
				eDate = dates[11:21]
			}
		} else if label == "" {
			summary = s.Text()
		}
	})

	tempDate, err := time.Parse("2006-01-02", sDate)
	if err != nil {
		log.Println("[ERROR] (gimhae.parseExhibitions) Error parsing start date: ", err)
		return nil, err
	}
	if tempDate.Before(time.Now()) {
		log.Println("[INFO] (gimhae.parseExhibitions) Exhibition already ended, skipping")
	}

	return &parser.ParseResult{
		FoundExhibitions: []*common.Exhibition{
			{
				ExhibitionId:     common.GenerateExhibitionId("gimhaemuseum", sDate, title),
				VenueVisitKor2Id: "130690",
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

func GetGimhaeMuseum() parser.MuseumPageParser {
	return &GimhaeMuseum{}
}
