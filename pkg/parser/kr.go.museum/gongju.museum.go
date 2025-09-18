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

type GongjuMuseum struct {
}

func (g *GongjuMuseum) Parsing(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	switch job.Url.Path {
	case "/prog/speclDspy/kor/sub02_02_01/view.do":
		return g.parseExhibitions(ctx, cfg, job)
	default:
		return nil, nil
	}
}

func (g *GongjuMuseum) parseExhibitions(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	log.Println("[INFO] (gongju.parseDetailPage) Parsing detail page ", job.Url)

	if job.Url.Host == "" {
		job.Url.Host = "gongju.museum.go.kr"
		job.Url.Scheme = "https"
	}

	doc, err := crawler.DoCrawl(job.Url.String())
	if err != nil {
		log.Println("[ERROR] (gongju.parseDetailPage) Error crawling detail page: ", err)
		return nil, err
	}

	img := doc.Find("div.inner-box > div > div.inner > img")
	imageURLstr, exists := img.Attr("src")
	if !exists {
		imageURLstr = ""
	}

	imageURL, err := url.Parse(imageURLstr)
	if err != nil {
		imageURL = &url.URL{}
	}

	if imageURL.Host == "" {
		imageURL.Host = "gongju.museum.go.kr"
		imageURL.Scheme = "https"
	}

	title := doc.Find("div.inner-box > div > strong.title").Text()
	var sDate, eDate, summary string
	doc.Find("div.inner-box > div > ul.list-1st > li").Each(func(i int, s *goquery.Selection) {
		label := s.Find("span.tit").Text()
		value := s.Find("span.con").Text()

		if label == "기간" {
			dates := value
			if len(dates) >= 23 {
				dates = strings.TrimSpace(dates)
				dates = strings.ReplaceAll(dates, ".", "-")
				dates = strings.ReplaceAll(dates, " ~ ", "~")
				sDate = dates[0:10]
				eDate = dates[11:21]
			}
		} else if label == "설명" {
			summary = value
		}
	})

	tempDate, err := time.Parse("2006-01-02", sDate)
	if err != nil {
		log.Println("[ERROR] (gongju.parseExhibitions) Error parsing start date: ", err)
		return nil, err
	}
	if tempDate.Before(time.Now()) {
		log.Println("[INFO] (gongju.parseExhibitions) Exhibition already ended, skipping")
	}

	return &parser.ParseResult{
		FoundExhibitions: []*common.Exhibition{
			{
				ExhibitionId:     common.GenerateExhibitionId("gongjumuseum", sDate, title),
				VenueVisitKor2Id: "129787",
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

func GetGongjuMuseum() parser.MuseumPageParser {
	return &GongjuMuseum{}
}
