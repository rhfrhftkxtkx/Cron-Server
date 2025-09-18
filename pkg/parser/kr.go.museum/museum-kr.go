package kr_go_museum

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"regexp"

	"github.com/BlueNyang/theday-theplace-cron/pkg/config"
	"github.com/BlueNyang/theday-theplace-cron/pkg/crawler"
	"github.com/BlueNyang/theday-theplace-cron/pkg/domain/common"
	"github.com/BlueNyang/theday-theplace-cron/pkg/gemini"
	"github.com/BlueNyang/theday-theplace-cron/pkg/parser"
	"github.com/PuerkitoBio/goquery"
)

type Museum struct {
}

func optimizeHtml(content *goquery.Selection) *string {
	log.Println("[INFO] (parser...museum.processExhibition) optimizing html...")
	content.Find("script, style").Remove()

	optRegs := []*regexp.Regexp{
		regexp.MustCompile(`(?i)<br\s*/?>`),
		regexp.MustCompile(`(?i)<script\s*.*?>.*?</script>`),
		regexp.MustCompile(`(?i)<style\s*.*?>.*?</style>`),
		regexp.MustCompile(`(?i)<!--.*?-->`),
		regexp.MustCompile(`style="font.*?"`),
		regexp.MustCompile(`class=".*?"`),
		regexp.MustCompile(`&nbsp;`),
		regexp.MustCompile(`\n+\s*`),
	}

	html, err := content.Html()
	if err != nil {
		log.Println("[ERROR] (kr.go.museum.go.museum.museum-kr.go) OptimizeHtml error: ", err)
		return nil
	}

	for _, reg := range optRegs {
		html = reg.ReplaceAllString(html, "")
	}

	//log.Println(html)
	return &html
}

func convertToExhibition(url *url.URL, data *gemini.Response) (*common.Exhibition, *gemini.Response, error) {
	log.Println("[INFO] (parser...museum.convertToExhibition) Converting Gemini response to Exhibition struct")
	if data == nil {
		return nil, nil, fmt.Errorf("data is nil")
	}

	log.Printf("[INFO] (parser...museum.convertToExhibition) Received data: %+v\n", *data)

	reg := regexp.MustCompile(`https://www.museum.go.kr.*`)

	parsedImageURL, err := url.Parse(data.ImageURL)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid image URL: %v", err)
	}

	if parsedImageURL.Hostname() == "" {
		data.ImageURL = "https://www.museum.go.kr" + data.ImageURL
	}

	if data.Title == "" || data.StartDate == "" || data.EndDate == "" {
		return nil, nil, fmt.Errorf("missing required fields in data: %+v", *data)
	}

	if data.Depth < 2 && data.RelatedURL != "" && data.RelatedURL != url.String() {
		if !reg.MatchString(data.RelatedURL) {
			data.RelatedURL = url.Scheme + "://" + url.Hostname() + data.RelatedURL
		}
		data.Depth = data.Depth + 1
		return nil, data, nil
	}

	exchibitionId := common.GenerateExhibitionId("museumkr", data.StartDate, data.Title)
	exhibition := common.Exhibition{
		ExhibitionId:     exchibitionId,
		VenueVisitKor2Id: "129703",
		DataSourceTier:   1,
		Title:            data.Title,
		Summary:          data.Summary,
		StartDate:        data.StartDate,
		EndDate:          data.EndDate,
		ImageUrl:         data.ImageURL,
		SourceURL:        url.String(),
	}

	return &exhibition, nil, nil
}

func (m Museum) Parsing(ctx context.Context, cfg *config.Config, job parser.Job) (*parser.ParseResult, error) {
	client, err := gemini.InitGemini(ctx, cfg.GoogleAPIKey, "gemini-2.5-flash-lite")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] (parser...museum.Parsing) InitGemini error: %v", err)
	}

	log.Println("[INFO] (parser...museum.Parsing) Processing URL: ", job.Url.String(), " at depth: ", job.Depth)
	doc, err := crawler.DoCrawl(job.Url.String())
	if err != nil {
		return nil, fmt.Errorf("[ERROR] (parser...museum.Parsing) DoCrawl error: %v", err)
	}

	var classToFind = ".page-content-type2"

	content := doc.Find(classToFind)
	contentHtml := optimizeHtml(content)
	if contentHtml == nil || *contentHtml == "" {
		return nil, fmt.Errorf("[ERROR] (parser...museum.processExhibition) contentHtml is nil or empty: %v", err)
	}

	dataList, err := client.Processing(ctx, job.Url, *contentHtml, job.Depth)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] (parser...museum.processExhibition) Gemini Processing error: %v", err)
	}
	log.Printf("[INFO] (parser...museum.processExhibition) Gemini returned %d items\n", len(*dataList))

	var foundedExhibitions []*common.Exhibition
	var discoveredJobs []*parser.Job

	for _, data := range *dataList {
		exhibition, subJob, err := convertToExhibition(job.Url, &data)
		if err != nil {
			log.Println("[ERROR] (parser...museum.processExhibition) convertToExhibition error: ", err)
		} else if subJob != nil {
			log.Println("[INFO] (parser...museum.processExhibition) Enqueuing sub-job for URL: ", subJob.RelatedURL)

			dataUrl, err := url.Parse(subJob.RelatedURL)
			if err != nil {
				log.Println("[ERROR] (parser...museum.processExhibition) Invalid sub-job URL: ", subJob.RelatedURL, " error: ", err)
				continue
			}
			if dataUrl.Hostname() == "" {
				subJob.RelatedURL = job.Url.Scheme + "://" + job.Url.Hostname() + subJob.RelatedURL
			}

			discoveredJobs = append(discoveredJobs, &parser.Job{
				Url:   dataUrl,
				Depth: subJob.Depth,
			})
		} else {
			log.Println("[INFO] (parser...museum.processExhibition) Sending exhibition to result channel: ", exhibition.ExhibitionId)
			foundedExhibitions = append(foundedExhibitions, exhibition)
		}
	}

	return &parser.ParseResult{
		FoundExhibitions: foundedExhibitions,
		DiscoveredJobs:   discoveredJobs,
	}, nil
}

func GetMuseum() Museum {
	return Museum{}
}
