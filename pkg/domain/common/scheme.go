package common

import (
	"crypto/sha256"
	"fmt"
)

type Exhibition struct {
	ExhibitionId     string `json:"exhibition_id"`
	VenueVisitKor2Id string `json:"venue_visit_kor2_id"`
	DataSourceTier   int    `json:"data_source_tier"`
	Title            string `json:"title"`
	Summary          string `json:"summary"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	ImageUrl         string `json:"image_url"`
	SourceURL        string `json:"source_url"`
}

type CrawlTarget struct {
	URL      string
	Provider string
}

func GenerateExhibitionId(provider, startDate, title string) string {
	//title을 hash로 바꾸는게 좋을듯
	h := sha256.New()
	h.Write([]byte(title))
	title = fmt.Sprintf("%x", h.Sum(nil))
	//return provider + "_" + startDate + "_" + title
	exhibitionId := fmt.Sprintf("%s_%s_%s", provider, startDate, title)
	return exhibitionId
}
