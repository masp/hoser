package parser

import (
	"errors"
	"fmt"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/lexer"
	"github.com/masp/hoser/token"
)

type parser struct {
	file    *token.File
	scanner *lexer.Scanner
	errors  token.ErrorList

	// if peek() is called, peeked will be the cached values so that
	// calls to eat() will consume this rather than calling next() again.
	peeked tokenInfo
}

type tokenInfo struct {
	pos token.Pos
	tok token.Token
	lit string
}

var (
	ErrUnexpectedEnd      = errors.New("unexpected end of text")
	ErrExpectedExpression = errors.New("expected expression")
)

func ParseModule(file *token.File, src []byte) (module *ast.Module, err error) {
	scanner := lexer.NewScanner(file, src)
	p := parser{file: file, scanner: scanner}
	defer p.handleErrors(&err)
	module = p.parseModule()
	return
}

func ParseBlock(file *token.File, src []byte) (block *ast.BlockDecl, err error) {
	scanner := lexer.NewScanner(file, src)
	p := parser{file: file, scanner: scanner}
	defer p.handleErrors(&err)
	block = p.parseBlock()
	return
}

func ParseExpression(file *token.File, src []byte) (expr ast.Expr, err error) {
	scanner := lexer.NewScanner(file, src)
	p := parser{file: file, scanner: scanner}
	defer p.handleErrors(&err)
	expr = p.parseExpression(token.Invalid)
	return
}

// recover bailout panics where we have gotten too many errors
func (p *parser) handleErrors(err *error) {
	if e := recover(); e != nil {
		// resume same panic if it's not a bailout
		if _, ok := e.(bailout); !ok {
			panic(e)
		}
	}
	p.errors.Sort()
	*err = p.errors.Err()
}

// A bailout panic is raised to indicate early termination.
type bailout struct{}

func (p *parser) error(pos token.Pos, err error) {
	epos := p.file.Position(pos)

	// If AllErrors is not set, discard errors reported on the same line
	// as the last recorded error and stop parsing if there are more than
	// 10 errors.
	n := len(p.errors)
	if n > 0 && p.errors[n-1].Pos.Line == epos.Line {
		return // discard - likely a spurious error
	}
	if n > 10 {
		panic(bailout{})
	}

	p.errors.Add(epos, err)
}

func (p *parser) expectedError(posOrTok interface{}, msg string) {
	msg = "expected " + msg
	var pos token.Pos
	if got, ok := posOrTok.(tokenInfo); ok {
		// the error happened at the current position;
		// make the error message more specific
		pos = got.pos
		switch {
		case got.tok == token.Semicolon && got.lit == "\n":
			msg += ", found newline"
		case got.tok.IsLiteral():
			// print 123 rather than 'INT', etc.
			msg += ", found " + got.lit
		default:
			msg += ", found '" + got.tok.String() + "'"
		}
	} else if pos, ok = posOrTok.(token.Pos); !ok {
		panic("invalid arg, expected token.Pos or tokenInfo")
	}
	p.error(pos, fmt.Errorf(msg))
}

func (p *parser) next() {
	p.peeked.pos, p.peeked.tok, p.peeked.lit = p.scanner.Next()
}

func (p *parser) peek() tokenInfo {
	if p.peeked.tok == token.Invalid {
		p.next()
	}
	return p.peeked
}

func (p *parser) eat() tokenInfo {
	tok := p.peek()
	p.peeked = tokenInfo{} // token is consumed immediately
	return tok
}

func (p *parser) eatOnly(expected token.Token) tokenInfo {
	next := p.peek()
	if next.tok == expected {
		return p.eat()
	}
	p.error(next.pos, fmt.Errorf("expected token %v, got %v", expected, next.tok))
	return next
}

func (p *parser) eatAll(expected token.Token) {
	for {
		next := p.peek()
		if next.tok != expected {
			break
		}
		p.eat()
	}
}

func precedence(kind token.Token) int {
	switch kind {
	case token.Semicolon:
		return -1 // End of expression always
	case token.Invalid:
		return 0 // lexer.Invalid is a special token meaning no parent
	case token.Equals:
		return 1
	case token.Colon:
		return 4
	case token.LParen:
		return 5
	default:
		// Every other token is lower precedence than these and signal an end to an expression
		return -1
	}
}

func (p *parser) parseStmt() ast.Stmt {
	return &ast.ExprStmt{X: p.parseExpression(token.Invalid)}
}

func (p *parser) parseExpression(parent token.Token) ast.Expr {
	left := p.parsePrefix()
	for {
		next := p.peek()
		if precedence(parent) >= precedence(next.tok) {
			return left
		}

		left = p.parseInfixExpr(left)
	}
}

func (p *parser) parsePrefix() ast.Expr {
	next := p.eat()

	switch next.tok {
	case token.LParen:
		return p.parseLParen(next)
	case token.Ident:
		return p.parseIdentifier(next)
	case token.String, token.Integer, token.Float:
		return p.parseLiteral(next)
	case token.LCurlyBrack:
		return p.parseFieldList(next)
	default:
		p.error(next.pos, fmt.Errorf("expected expression got %v", next.tok))
		return nil
	}
}

func (p *parser) parseInfixExpr(left ast.Expr) ast.Expr {
	next := p.eat()
	switch next.tok {
	case token.Equals:
		return p.parseEquals(left, next)
	case token.Colon:
		return p.parseField(left, next)
	case token.LParen:
		return p.parseBlockCall(left, next)
	default:
		p.error(next.pos, fmt.Errorf("invalid token for infix expression: %v", next.tok))
		return nil
	}
}
