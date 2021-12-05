package parser

import (
	"fmt"

	"github.com/masp/hoser/lexer"
)

func (s *parserState) parseModule() (*Module, error) {
	module := &Module{
		Blocks: make(map[string]*Block),
	}

	for {
		err := s.eatAll(lexer.Semicolon)
		if err != nil {
			// no more functions
			return module, nil
		}
		fn, err := s.parseFunction()
		if err != nil {
			return nil, fmt.Errorf("%v:%v: syntax error at %v: %w", s.nextToken.Line, s.nextToken.Col, s.nextToken.Value, err)
		}

		if fn == nil {
			// no more functions
			return module, nil
		}
		module.Blocks[fn.Name.Token.Value] = fn
	}
}
