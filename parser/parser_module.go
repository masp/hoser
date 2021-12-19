package parser

import (
	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/token"
)

func (p *parser) parseModule() *ast.Module {
	p.eatOnly(token.Module)
	name := p.parseIdentifier(p.eat())

	module := &ast.Module{
		Name: name,
	}
	for {
		p.eatAll(token.Semicolon)
		block := p.parseBlock()

		if block == nil {
			return module // no more blocks
		}
		module.Blocks = append(module.Blocks, block)
	}
}
