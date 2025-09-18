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

type IksanMuseum struct {
}

func (g *IksanMuseum) Parsing(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	switch job.Url.Path {
	case "/kor/html/sub02/0202.html":
		return g.parseExhibitions(ctx, cfg, job)
	default:
		return nil, nil
	}
}

func (g *IksanMuseum) parseExhibitions(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	log.Println("[INFO] (iksan.parseDetailPage) Parsing detail page ", job.Url)

	if job.Url.Host == "" {
		job.Url.Host = "iksan.museum.go.kr"
		job.Url.Scheme = "https"
	}

	doc, err := crawler.DoCrawl(job.Url.String())
	if err != nil {
		log.Println("[ERROR] (iksan.parseDetailPage) Error crawling detail page: ", err)
		return nil, err
	}

	img := doc.Find("div.exhibit_con_l >div > img")
	imageURLstr, exists := img.Attr("src")
	if !exists {
		imageURLstr = ""
	}

	imageURL, err := url.Parse(imageURLstr)
	if err != nil {
		imageURL = &url.URL{}
	}

	if imageURL.Host == "" {
		imageURL.Host = "iksan.museum.go.kr"
		imageURL.Scheme = "https"
	}

	title := doc.Find("p.titp").Text()
	summary := doc.Find("div.exhibit_con_r p:nth-child(2)").Text()
	dates := doc.Find("ul.sp_ul li:nth-child(2)").Text()
	var sDate, eDate string
	if len(dates) >= 23 {
		dates = strings.TrimSpace(dates)
		dates = strings.ReplaceAll(dates, ".", "-")
		dates = strings.ReplaceAll(dates, " ~ ", "~")
		sDate = dates[0:10]
		eDate = dates[11:21]
	}

	tempDate, err := time.Parse("2006-01-02", sDate)
	if err != nil {
		log.Println("[ERROR] (iksan.parseExhibitions) Error parsing start date: ", err)
		return nil, err
	}
	if tempDate.Before(time.Now()) {
		log.Println("[INFO] (iksan.parseExhibitions) Exhibition already ended, skipping")
	}

	return &parser.ParseResult{
		FoundExhibitions: []*common.Exhibition{
			{
				ExhibitionId:     common.GenerateExhibitionId("iksanmuseum", sDate, title),
				VenueVisitKor2Id: "130254",
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

func GetIksanMuseum() parser.MuseumPageParser {
	return &IksanMuseum{}
}
