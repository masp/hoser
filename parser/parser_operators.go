package parser

import (
	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/token"
)

func (p *parser) parseEquals(left ast.Expr, eq tokenInfo) *ast.AssignExpr {
	right := p.parseExpression(eq.tok)
	return &ast.AssignExpr{Lhs: left, EqPos: eq.pos, Rhs: right}
}

func (p *parser) parseLParen(lparen tokenInfo) *ast.ParenExpr {
	expr := p.parseExpression(token.Invalid)
	p.eatOnly(token.RParen)
	return &ast.ParenExpr{X: expr}
}
