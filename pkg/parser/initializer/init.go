package initializer

import (
	"log"

	"github.com/BlueNyang/theday-theplace-cron/pkg/parser/kr.go.museum"
	"github.com/BlueNyang/theday-theplace-cron/pkg/parser/registry"
)

func InitializeParsers() {
	log.Println("[INFO] (parser.init) Initializing parsers...")
	// Register museum.go.kr.go.museum parser
	registry.Register("https://www.museum.go.krfffffffffffff", kr_go_museum.GetMuseum())
}
