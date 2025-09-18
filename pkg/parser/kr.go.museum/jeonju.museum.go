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
)

type JeonjuMuseum struct {
}

func (g *JeonjuMuseum) Parsing(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	switch job.Url.Path {
	case "/special.es":
		return g.parseExhibitions(ctx, cfg, job)
	default:
		return nil, nil
	}
}

func (g *JeonjuMuseum) parseExhibitions(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	log.Println("[INFO] (jeonju.parseExhibitions) Parsing detail page ", job.Url)

	if job.Url.Host == "" {
		job.Url.Host = "jeonju.museum.go.kr"
		job.Url.Scheme = "https"
	}

	doc, err := crawler.DoCrawl(job.Url.String())
	if err != nil {
		log.Println("[ERROR] (jeonju.parseExhibitions) Error crawling detail page: ", err)
		return nil, err
	}

	img := doc.Find("p.pic > img")
	imageURLstr, exists := img.Attr("src")
	if !exists {
		imageURLstr = ""
	}

	imageURL, err := url.Parse(imageURLstr)
	if err != nil {
		imageURL = &url.URL{}
	}

	if imageURL.Host == "" {
		imageURL.Host = "jeonju.museum.go.kr"
		imageURL.Scheme = "https"
	}

	title := doc.Find("p.title").Text()
	summary := ""
	dates := doc.Find("dd.color-pink").Text()
	dates = strings.TrimSpace(dates)
	dates = strings.ReplaceAll(dates, ".", "-")
	dates = strings.ReplaceAll(dates, " ~ ", "~")
	sDate := dates[0:10]
	eDate := dates[11:21]

	tempDate, err := time.Parse("2006-01-02", sDate)
	if err != nil {
		log.Println("[ERROR] (jeonju.parseExhibitions) Error parsing start date: ", err)
		return nil, err
	}
	if tempDate.Before(time.Now()) {
		log.Println("[INFO] (jeonju.parseExhibitions) Exhibition already ended, skipping")
	}

	return &parser.ParseResult{
		FoundExhibitions: []*common.Exhibition{
			{
				ExhibitionId:     common.GenerateExhibitionId("jeonjumuseum", sDate, title),
				VenueVisitKor2Id: "129786",
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

func GetJeonjuMuseum() parser.MuseumPageParser {
	return &JeonjuMuseum{}
}
