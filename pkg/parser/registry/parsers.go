package registry

import (
	"fmt"
	"log"

	"github.com/BlueNyang/theday-theplace-cron/pkg/parser"
)

var parserRegistry = make(map[string]parser.MuseumPageParser)

func Register(name string, p parser.MuseumPageParser) {
	log.Println("[INFO] (parser.Register) Registering parser:", name)
	parserRegistry[name] = p
}

func GetParser(name string) (parser.MuseumPageParser, error) {
	log.Println("[INFO] (parser.GetParser) Getting parser:", name)
	p, exists := parserRegistry[name]
	if !exists {
		return nil, fmt.Errorf("parser %s not found", name)
	}
	return p, nil
}
