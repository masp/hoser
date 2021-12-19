package parser

import (
	"errors"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/token"
)

var (
	ErrExpectedFnBody    = errors.New("expected start of function body with '{'")
	ErrInputOutputNotMap = errors.New("input & output of function must be maps")
	ErrInvalidBlockName  = errors.New("invalid block name, expected identifier")
	// An empty function body is different than no function body
	// `main() {}` -> empty
	// vs.
	// `main()` -> nil
	EmptyFnBody = []ast.Expr{}
)

func (p *parser) parseBlock() *ast.BlockDecl {
	ident := p.eat()
	if ident.tok == token.Eof {
		return nil // No more blocks, so we finish
	} else if ident.tok != token.Ident {
		p.eatOnly(token.Ident) // use error handling in eatOnly
		return nil             // TODO: return badblock
	}
	name := p.parseIdentifier(ident)

	inputs := p.parseArgs()

	// output is optional
	// e.g. `main() () {}`` is equivalent to `main() {}`
	next := p.peek()
	var outputs *ast.FieldList
	if next.tok == token.LParen {
		// parse output definition
		outputs = p.parseArgs()
	}

	var body []ast.Stmt
	var lbrack, rbrack token.Pos
	next = p.peek()
	if next.tok == token.LCurlyBrack {
		lbrack = p.eatOnly(token.LCurlyBrack).pos
		body = p.parseFnBody()
		rbrack = p.eatOnly(token.RCurlyBrack).pos
	} else if next.tok != token.Eof && next.tok != token.Semicolon {
		p.expectedError(next, "'{' or end of line")
	}

	return &ast.BlockDecl{
		Name:      name,
		Inputs:    inputs,
		Outputs:   outputs,
		BegLBrack: lbrack,
		Body:      body,
		EndRBrack: rbrack,
	}
}

// parseArgs takes either the input or output arguments specification and converts it to a Map
// example:
// ([name: string, value: int]) -> Map{{Key: name, Val: string}, {Key: value, Val: int}}
func (p *parser) parseArgs() *ast.FieldList {
	lparen := p.eatOnly(token.LParen).pos
	var rparen token.Pos

	next := p.peek()

	var fields []*ast.Field
	if next.tok != token.RParen {
		args := p.parseExpression(token.Invalid)

		rparen = p.eatOnly(token.RParen).pos

		switch ent := args.(type) {
		case *ast.FieldList:
			fields = ent.Fields
		case *ast.Field:
			fields = []*ast.Field{ent}
		default:
			p.error(next.pos, ErrInputOutputNotMap)
		}
	} else {
		rparen = p.eat().pos
	}
	return &ast.FieldList{
		Opener: lparen,
		Fields: fields,
		Closer: rparen,
	}
}

func (p *parser) parseFnBody() []ast.Stmt {
	var exprs []ast.Stmt
	for {
		exprTok := p.peek()

		if exprTok.tok == token.RCurlyBrack {
			break // no more expressions in body
		}
		exprs = append(exprs, p.parseStmt())

		// expressions always end in either ; or }. If there's one or more ;, we should eat them.
		p.eatAll(token.Semicolon)
	}
	return exprs
}

func (p *parser) parseBlockCall(left ast.Expr, lparen tokenInfo) *ast.CallExpr {
	if name, ok := left.(*ast.Ident); ok {
		next := p.peek()

		result := &ast.CallExpr{Name: name, Lparen: lparen.pos}
		for next.tok != token.RParen {
			// Let's parse the args (non-empty), and we continue parsing until we match our end paren and reset precedence
			arg := p.parseExpression(token.Invalid)
			result.Args = append(result.Args, arg)
			next = p.peek()
			if next.tok == token.Comma {
				p.eatOnly(token.Comma)
			} else if next.tok != token.RParen {
				p.expectedError(next, "comma or right paren ')'")
				return result
			}
		}
		result.Rparen = p.eatOnly(token.RParen).pos

		return result
	} else {
		p.expectedError(left.Pos(), "block name")
	}
	return nil
}
