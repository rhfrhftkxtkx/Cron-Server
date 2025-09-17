package instances

import (
	"context"

	"github.com/BlueNyang/theday-theplace-cron/pkg/config"
)

type MuseumPageParser interface {
	Parsing(ctx context.Context, cfg *config.Config, job Job) (*ParseResult, error)
}
