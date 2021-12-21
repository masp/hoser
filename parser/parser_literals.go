package parser

import (
	"strconv"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/token"
)

func (p *parser) parseIdentifier(first tokenInfo) *ast.Ident {
	next := p.peek()
	if next.tok == token.Period {
		p.eat()
		name := p.eatOnly(token.Ident)
		return &ast.Ident{V: name.lit, NamePos: name.pos, Module: first.lit, ModulePos: first.pos}
	}
	return &ast.Ident{V: first.lit, NamePos: first.pos}
}

func (p *parser) parseLiteral(tok tokenInfo) *ast.LiteralExpr {
	var parsedVal interface{}
	switch tok.tok {
	case token.Integer:
		v, err := strconv.ParseInt(tok.lit, 10, 64)
		if err != nil {
			p.error(tok.pos, err)
		}
		parsedVal = v
	case token.Float:
		v, err := strconv.ParseFloat(tok.lit, 64)
		if err != nil {
			p.error(tok.pos, err)
		}
		parsedVal = v
	case token.String:
		parsedVal = tok.lit
	}
	return &ast.LiteralExpr{Start: tok.pos, Type: tok.tok, Value: tok.lit, ParsedVal: parsedVal}
}
