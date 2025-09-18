package initializer

import (
	"log"

	"github.com/BlueNyang/theday-theplace-cron/pkg/parser/kr.go.museum"
	"github.com/BlueNyang/theday-theplace-cron/pkg/parser/registry"
)

func InitializeParsers() {
	log.Println("[INFO] (parser.init) Initializing parsers...")
	// Register museum.go.kr.go.museum parser
	registry.Register("www.museum.go.kr", kr_go_museum.GetMuseum())
	registry.Register("buyeo.museum.go.kr", kr_go_museum.GetBuyeoMuseum())
	registry.Register("gongju.museum.go.kr", kr_go_museum.GetGongjuMuseum())
	registry.Register("gimhae.museum.go.kr", kr_go_museum.GetGimhaeMuseum())
	registry.Register("iksan.museum.go.kr", kr_go_museum.GetIksanMuseum())
	registry.Register("chuncheon.museum.go.kr", kr_go_museum.GetChuncheonMuseum())
	registry.Register("cheongju.museum.go.kr", kr_go_museum.GetCheongjuMuseum())
	registry.Register("jeju.museum.go.kr", kr_go_museum.GetJejuMuseum())
	registry.Register("jeonju.museum.go.kr", kr_go_museum.GetJeonjuMuseum())
	log.Println("[INFO] (parser.init) Parsers initialized.")
}
