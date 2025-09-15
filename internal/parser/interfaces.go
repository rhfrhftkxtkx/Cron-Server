package parser

import (
	"github.com/BlueNyang/theday-theplace-cron/internal/domain/common"
	"github.com/PuerkitoBio/goquery"
)

type Parser interface {
	Parse(doc *goquery.Document) ([]common.Exhibition, error)
}
