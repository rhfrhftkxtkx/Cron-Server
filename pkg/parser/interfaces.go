package parser

import (
	"context"
	"net/url"

	"github.com/BlueNyang/theday-theplace-cron/pkg/config"
	"github.com/BlueNyang/theday-theplace-cron/pkg/domain/common"
)

type MuseumPageParser interface {
	Parsing(ctx context.Context, cfg *config.Config, job Job) (*ParseResult, error)
}

type ParseResult struct {
	FoundExhibitions []*common.Exhibition
	DiscoveredJobs   []*Job
}

type Job struct {
	Url   *url.URL
	Depth int
}
