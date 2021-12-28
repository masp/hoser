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

func (p *parser) parseStubBlock() (stub ast.StubDecl) {
	stub.Name = p.parseIdentifier(p.eatOnly(token.Ident))
	stub.Inputs = p.parseArgs()

	// output is optional
	// e.g. `main() () {}`` is equivalent to `main() {}`
	next := p.peek()
	if next.tok == token.LParen {
		// parse output definition
		stub.Outputs = p.parseArgs()
	}
	return
}

func (p *parser) parsePipeBlock() (pipe ast.PipeDecl) {
	pipe.StubDecl = p.parseStubBlock()
	pipe.BegLBrack = p.eatOnly(token.LCurlyBrack).pos
	pipe.Body = p.parseFnBody()
	pipe.EndRBrack = p.eatOnly(token.RCurlyBrack).pos
	return
}

// parseArgs takes either the input or output arguments specification and converts it to a Map
// example:
// ([name: string, value: int]) -> Map{{Key: name, Val: string}, {Key: value, Val: int}}
func (p *parser) parseArgs() ast.FieldList {
	opener := p.eatOnly(token.LParen)
	return p.parseFieldList(opener)
}

func (p *parser) parseFnBody() []ast.Stmt {
	var exprs []ast.Stmt
	for {
		exprTok := p.peek()

		if exprTok.tok == token.RCurlyBrack || exprTok.tok == token.Eof {
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
