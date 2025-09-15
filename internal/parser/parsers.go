package parser

import (
	"fmt"

	"github.com/yuin/goldmark/parser"
)

var parserRegistry = make(map[string]parser.Parser)

func Register(name string, p parser.Parser) {
	parserRegistry[name] = p
}

func GetParser(name string) (parser.Parser, error) {
	p, exists := parserRegistry[name]
	if !exists {
		return nil, fmt.Errorf("parser %s not found", name)
	}
	return p, nil
}
