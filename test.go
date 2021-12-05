package main

import (
	"errors"
	"fmt"
	"log"
)

type Expr interface {
	String() string
}

type TokenKind int

const (
	TokLowest TokenKind = iota
	TokAdd
	TokMult
	TokId
	TokNeg
	TokInvalid TokenKind = -999
)

func tokType(tok byte) TokenKind {
	switch tok {
	case '-':
		return TokNeg
	case 'a', 'b', 'c', 'd':
		return TokId
	case '+':
		return TokAdd
	case '*':
		return TokMult
	}
	return TokInvalid
}

func lex(text string) []TokenKind {
	toks := make([]TokenKind, len(text))
	for i := 0; i < len(text); i++ {
		toks[i] = tokType(text[i])
	}
	return toks
}

type Id struct {
}

func (i *Id) String() string {
	return "a"
}

type Plus struct {
	left  Expr
	right Expr
}

func (i *Plus) String() string {
	return fmt.Sprintf("(%v+%v)", i.left, i.right)
}

type Mult struct {
	left  Expr
	right Expr
}

func (i *Mult) String() string {
	return fmt.Sprintf("(%v*%v)", i.left, i.right)
}

type Neg struct {
	what Expr
}

func (i *Neg) String() string {
	return fmt.Sprintf("(-%v)", i.what)
}

type parser struct {
	text []TokenKind
	cur  int
}

func (p *parser) peek() (TokenKind, error) {
	if p.cur >= len(p.text) {
		return 0, nil
	}
	return p.text[p.cur], nil
}

func (p *parser) eat() (TokenKind, error) {
	next, err := p.peek()
	if err != nil {
		return 0, err
	}

	p.cur += 1
	return next, nil
}

func (p *parser) parseExpr(precedence TokenKind) (Expr, error) {
	left, err := p.parsePrefix()
	if err != nil {
		return nil, err
	}

	for {
		next, err := p.peek()
		if err != nil {
			return nil, err
		}

		if int(precedence) >= int(next) {
			return left, nil
		}

		left, err = p.parseInfix(left)
		if err != nil {
			return nil, err
		}
	}
}

var EndOfExpression = errors.New("end of expression")

func (p *parser) parsePrefix() (Expr, error) {
	tok, err := p.eat()
	if err != nil {
		return nil, err
	}
	switch tok {
	case TokNeg:
		return p.parseNeg(tok)
	case TokId:
		return p.parseId(tok)
	default:
		return nil, EndOfExpression
	}
}

func (p *parser) parseInfix(left Expr) (Expr, error) {
	tok, err := p.eat()
	if err != nil {
		return nil, err
	}

	switch tok {
	case TokAdd:
		return p.parseAdd(left, tok)
	case TokMult:
		return p.parseMult(left, tok)
	default:
		return nil, EndOfExpression
	}
}

func (p *parser) parseNeg(tok TokenKind) (Expr, error) {
	what, err := p.parseExpr(tok)
	if err != nil {
		return nil, err
	}
	return &Neg{what: what}, nil
}

func (p *parser) parseId(tok TokenKind) (Expr, error) {
	return &Id{}, nil
}

func (p *parser) parseAdd(left Expr, tok TokenKind) (Expr, error) {
	right, err := p.parseExpr(tok)
	if err != nil {
		return nil, err
	}
	return &Plus{left: left, right: right}, nil
}

func (p *parser) parseMult(left Expr, tok TokenKind) (Expr, error) {
	right, err := p.parseExpr(tok)
	if err != nil {
		return nil, err
	}
	return &Mult{left: left, right: right}, nil
}

func main() {
	p := &parser{text: lex("a+a+a")}
	expr, err := p.parseExpr(TokLowest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(expr)
}
