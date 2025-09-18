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

type JejuMuseum struct {
}

func (b *JejuMuseum) Parsing(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	switch job.Url.Path {
	case "/_prog/special_exhibit/index.php":
		return b.parseDetailPage(ctx, cfg, job)
	}

	return nil, nil
}

func (b *JejuMuseum) parseDetailPage(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	log.Println("[INFO] (jeju.parseDetailPage) Parsing detail page ", job.Url)

	if job.Url.Host == "" {
		job.Url.Host = "jeju.museum.go.kr"
		job.Url.Scheme = "https"
	}

	doc, err := crawler.DoCrawl(job.Url.String())
	if err != nil {
		log.Println("[ERROR] (jeju.parseDetailPage) Error crawling detail page: ", err)
		return nil, err
	}

	img := doc.Find("figure > img")
	imageURLStr, exists := img.Attr("src")
	if !exists {
		imageURLStr = ""
	}

	imageURL, err := url.Parse(imageURLStr)
	if err != nil {
		imageURL = &url.URL{}
	}
	if imageURL.Host == "" {
		imageURL.Host = "jeju.museum.go.kr"
		imageURL.Scheme = "https"
	}

	summary := doc.Find("div.exhib_detail_txt").Find("img, br").Remove().End().Text()
	var sDate, eDate, title string
	doc.Find("div.exhib_info > ul > li").Each(func(i int, s *goquery.Selection) {
		label := s.Find("b").Text()
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
		log.Println("[ERROR] (jeju.parseExhibitions) Error parsing start date: ", err)
		return nil, err
	}
	if tempDate.Before(time.Now()) {
		log.Println("[INFO] (jeju.parseExhibitions) Exhibition already ended, skipping")
	}

	return &parser.ParseResult{
		FoundExhibitions: []*common.Exhibition{
			{
				ExhibitionId:     common.GenerateExhibitionId("jejumuseum", sDate, title),
				VenueVisitKor2Id: "130461",
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

func GetJejuMuseum() parser.MuseumPageParser {
	return &JejuMuseum{}
}
