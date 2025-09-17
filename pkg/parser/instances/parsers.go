package instances

import (
	"fmt"
	"log"

	"github.com/BlueNyang/theday-theplace-cron/pkg/domain/common"
)

var parserRegistry = make(map[string]MuseumPageParser)

func Register(name string, p MuseumPageParser) {
	log.Println("[INFO] (parser.Register) Registering parser:", name)
	parserRegistry[name] = p
}

func GetParser(name string) (MuseumPageParser, error) {
	log.Println("[INFO] (parser.GetParser) Getting parser:", name)
	p, exists := parserRegistry[name]
	if !exists {
		return nil, fmt.Errorf("parser %s not found", name)
	}
	return p, nil
}

type Job struct {
	Url   *string
	Depth int
}

type ParseResult struct {
	FoundExhibitions []*common.Exhibition
	DiscoveredJobs   []*Job
}
