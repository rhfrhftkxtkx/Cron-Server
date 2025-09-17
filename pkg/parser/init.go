package parser

import (
	"log"

	"github.com/BlueNyang/theday-theplace-cron/pkg/parser/instances"
	"github.com/BlueNyang/theday-theplace-cron/pkg/parser/kr.go.museum"
)

func InitializeParsers() {
	log.Println("[INFO] (parser.init) Initializing parsers...")
	// Register museum.go.kr.go.museum parser
	instances.Register("https://www.museum.go.kr/MUSEUM/contents/M0207000000.do", kr_go_museum.GetMuseum("https://www.museum.go.kr/MUSEUM/contents/M0207000000.do"))
	instances.Register("https://www.museum.go.kr/MUSEUM/contents/M0202010000.do", kr_go_museum.GetMuseum("https://www.museum.go.kr/MUSEUM/contents/M0202010000.do"))
}
