package parser

import (
	"errors"
	"fmt"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/token"
)

var (
	ErrInvalidMapValue   = errors.New("invalid value in map, must be 'key: value' pair")
	ErrInvalidEntryKey   = errors.New("key of an entry must be an identifier")
	ErrInvalidEntryValue = errors.New("value of an entry must be an identifier")
	ErrExpectedEndOfMap  = errors.New("expected ] to end map")
)

func flip(tok token.Token) token.Token {
	switch tok {
	case token.LCurlyBrack:
		return token.RCurlyBrack
	case token.RCurlyBrack:
		return token.LCurlyBrack
	case token.LParen:
		return token.RParen
	case token.RParen:
		return token.LParen
	default:
		panic(fmt.Errorf("%v is not a valid opener", tok))
	}
}

func (p *parser) parseFieldList(opener tokenInfo) (result *ast.FieldList) {
	closerTok := flip(opener.tok)
	next := p.peek()
	result = &ast.FieldList{Opener: opener.pos, Fields: nil, Closer: next.pos}
	for next.tok != closerTok {
		arg := p.parseExpression(token.Invalid)
		if ent, ok := arg.(*ast.Field); ok {
			result.Fields = append(result.Fields, ent)
		} else {
			p.expectedError(next.pos+1, "'key: value' pair")
		}
		next = p.peek()
	}

	result.Opener = opener.pos
	result.Closer = p.eatOnly(closerTok).pos
	return
}

func (p *parser) parseField(left ast.Expr, colon tokenInfo) *ast.Field {
	right := p.parseExpression(token.Colon)

	if left, ok := left.(*ast.Ident); ok {
		return &ast.Field{Key: left, Colon: colon.pos, Value: right}
	} else {
		p.expectedError(left.Pos(), "name")
	}
	return nil
}
